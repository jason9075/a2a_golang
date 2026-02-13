package server

import (
	"a2a/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TaskUpdateFunc func(event interface{})

type HandlerFunc func(task *models.Task, msg *models.Message, update TaskUpdateFunc) (*models.Task, error)

type Server struct {
	card    models.AgentCard
	handler HandlerFunc
}

func NewA2AServer(card models.AgentCard, handler HandlerFunc) *Server {
	return &Server{
		card:    card,
		handler: handler,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read body first, as we might need it for debugging or later use
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// If GET request, return Agent Card directly
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.card)
		return
	}

	var rpcReq models.JSONRPCRequest
	if err := json.Unmarshal(body, &rpcReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Dispatch based on Method
	switch rpcReq.Method {
	case "message/send":
		s.handleSend(w, rpcReq, body)
	case "message/stream":
		s.handleStream(w, rpcReq, body)
	default:
		// Default to card info or error?
		// Usually /a2a endpoint handles JSON-RPC. If GET, maybe return card?
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(s.card)
			return
		}
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleSend(w http.ResponseWriter, rpcReq models.JSONRPCRequest, body []byte) {
	// Re-decode params specifically as TaskSendParams
	// Since rpcReq.Params is interface{}, and we know the structure for "message/send"
	// A cleaner way is to unmarshal directly into a struct that matches the method,
	// but here we just re-unmarshal or use map.
	// Let's use map/interface casting or re-unmarshal for safety.
	
	// Better: define a struct for the request with known Params type
	var specificReq struct {
		Params models.TaskSendParams `json:"params"`
	}
	if err := json.Unmarshal(body, &specificReq); err != nil {
		http.Error(w, "Invalid Params", http.StatusBadRequest)
		return
	}

	task := &models.Task{
		ID:     specificReq.Params.ID,
		Status: models.TaskStatus{State: models.TaskStateWorking},
	}
	msg := &specificReq.Params.Message

	// Dummy update function
	update := func(event interface{}) {}

	resultTask, err := s.handler(task, msg, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSONRPCResponse
	resp := models.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      rpcReq.JSONRPCMessageIdentifier.ID, // Extract ID string
		Result:  resultTask,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleStream(w http.ResponseWriter, rpcReq models.JSONRPCRequest, body []byte) {
	// Validate streaming support
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// Start stream
	
	// Parse params
	var specificReq struct {
		Params models.TaskSendParams `json:"params"`
	}
	if err := json.Unmarshal(body, &specificReq); err != nil {
		// In streaming, we might send an error event? Or just close.
		fmt.Fprintf(w, "data: {\"error\": \"Invalid Params\"}\n\n")
		return
	}

	task := &models.Task{
		ID:     specificReq.Params.ID,
		Status: models.TaskStatus{State: models.TaskStateWorking},
	}
	msg := &specificReq.Params.Message

	// Update function
	update := func(event interface{}) {
		// Wrap event in streaming response structure
		// The client expects SendTaskStreamingResponse
		resp := models.SendTaskStreamingResponse{
			Result: event,
		}
		data, _ := json.Marshal(resp)
		// SSE format: "data: <json>\n\n"? Or just raw JSON per line?
		// Client agent_a reads line by line and unmarshals.
		// Standard SSE uses "data: ...\n\n".
		// But agent_a code uses: `json.Unmarshal([]byte(line), &streamResp)`
		// This implies raw JSON lines (NDJSON), not standard SSE "data:" prefix.
		// Let's check agent_a loop:
		// line, err := reader.ReadString('\n')
		// json.Unmarshal([]byte(line), &streamResp)
		// So it is NDJSON!
		
		fmt.Fprintf(w, "%s\n", data)
		flusher.Flush()
	}

	// Call handler
	resultTask, err := s.handler(task, msg, update)
	if err != nil {
		// Log error?
		return
	}
	
	// Send final completion event
	final := true
	finalEvent := models.TaskArtifactUpdateEvent{
		ID:    resultTask.ID,
		Final: &final,
	}
	
	resp := models.SendTaskStreamingResponse{
		Result: finalEvent,
	}
	data, _ := json.Marshal(resp)
	fmt.Fprintf(w, "%s\n", data)
	flusher.Flush()
}
