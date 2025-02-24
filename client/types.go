package client

import "github.com/agent-api/core/types"

// ChatRequest represents a request to the chat endpoint
type ChatRequest struct {
	Model    string
	Messages []*types.Message
	Tools    []*types.Tool
}

// ChatResponse represents a response from the chat endpoint
type ChatResponse struct {
	Message types.Message
	Model   string
}
