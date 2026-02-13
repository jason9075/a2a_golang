package models

type AgentCard struct {
	Name         string            `json:"name"`
	Description  *string           `json:"description,omitempty"`
	Version      string            `json:"version"`
	URL          string            `json:"url"`
	Capabilities AgentCapabilities `json:"capabilities"`
	Skills       []AgentSkill      `json:"skills,omitempty"`
}

type AgentCapabilities struct {
	Streaming *bool `json:"streaming,omitempty"`
}

type AgentSkill struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type TaskState string

const (
	TaskStateWorking   TaskState = "working"
	TaskStateCompleted TaskState = "completed"
	TaskStateFailed    TaskState = "failed"
)

type Task struct {
	ID       string                 `json:"id"`
	Status   TaskStatus             `json:"status"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type TaskStatus struct {
	State TaskState `json:"state"`
}

type Message struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type Part struct {
	Text *string `json:"text,omitempty"`
}

type Artifact struct {
	Parts []Part `json:"parts"`
}

type TaskArtifactUpdateEvent struct {
	ID       string   `json:"id"`
	Artifact Artifact `json:"artifact"`
	Final    *bool    `json:"final,omitempty"`
}

// JSON-RPC Types
type JSONRPCMessageIdentifier struct {
	ID string `json:"id"`
}

type JSONRPCMessage struct {
	JSONRPC                  string                   `json:"jsonrpc"`
	JSONRPCMessageIdentifier JSONRPCMessageIdentifier `json:"id"`
}

type JSONRPCRequest struct {
	JSONRPCMessage
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type TaskSendParams struct {
	ID      string  `json:"id"`
	Message Message `json:"message"`
}

type SendTaskStreamingResponse struct {
	Result interface{} `json:"result"`
}
