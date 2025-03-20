package client

import "github.com/agent-api/core"

// ChatRequest represents a request to the chat endpoint
type ChatRequest struct {
	Model    string
	Messages []*core.Message
	Tools    []*core.Tool
}

// ChatResponse represents a response from the chat endpoint
type ChatResponse struct {
	Message core.Message
	Model   string
}
