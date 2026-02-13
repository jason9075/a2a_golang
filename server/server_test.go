package server

import (
	"a2a/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestA2AServer_Handshake(t *testing.T) {
	// 1. Setup
	card := models.AgentCard{
		Name:    "TestAgent",
		Version: "0.0.1",
		URL:     "http://test.local",
	}
	handler := func(task *models.Task, msg *models.Message, update TaskUpdateFunc) (*models.Task, error) {
		return task, nil
	}
	srv := NewA2AServer(card, handler)

	// 2. Execute (GET request for handshake/card)
	req := httptest.NewRequest(http.MethodGet, "/a2a", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	// 3. Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var respCard models.AgentCard
	if err := json.NewDecoder(w.Body).Decode(&respCard); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if respCard.Name != "TestAgent" {
		t.Errorf("Expected agent name 'TestAgent', got '%s'", respCard.Name)
	}
}

func TestA2AServer_MessageSend(t *testing.T) {
	// 1. Setup
	card := models.AgentCard{Name: "TestAgent"}
	handler := func(task *models.Task, msg *models.Message, update TaskUpdateFunc) (*models.Task, error) {
		task.Metadata = map[string]interface{}{"reply": "pong"}
		return task, nil
	}
	srv := NewA2AServer(card, handler)

	// 2. Execute (POST request for message/send)
	text := "ping"
	params := models.TaskSendParams{
		ID: "task-1",
		Message: models.Message{
			Role: "user",
			Parts: []models.Part{{Text: &text}},
		},
	}
	rpcReq := models.JSONRPCRequest{
		JSONRPCMessage: models.JSONRPCMessage{
			JSONRPC: "2.0",
			JSONRPCMessageIdentifier: models.JSONRPCMessageIdentifier{ID: "req-1"},
		},
		Method: "message/send",
		Params: params,
	}
	body, _ := json.Marshal(rpcReq)

	req := httptest.NewRequest(http.MethodPost, "/a2a", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	// 3. Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var rpcResp models.JSONRPCResponse
	if err := json.NewDecoder(w.Body).Decode(&rpcResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if rpcResp.ID != "req-1" {
		t.Errorf("Expected ID 'req-1', got '%s'", rpcResp.ID)
	}
}
