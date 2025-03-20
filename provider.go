package googlegenai

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/agent-api/core"
	"github.com/agent-api/googlegenai/client"
	"github.com/agent-api/googlegenai/models"
)

// Provider implements the LLMProvider interface for Google Gen AI
type Provider struct {
	host string
	port int

	model *core.Model

	// client is the internal Ollama HTTP client
	client *client.GoogleGenAIClient

	logger *logr.Logger
}

type ProviderOpts struct {
	BaseURL string
	Port    int
	APIKey  string

	Logger *logr.Logger
}

// NewProvider creates a new Google gen AI provider
func NewProvider(opts *ProviderOpts) *Provider {
	ctx := context.Background()
	opts.Logger.Info("creating a new google generative ai provider")

	client, err := client.NewClient(ctx, &client.GoogleGenAIClientOpts{
		Model:  models.GEMINI_1_5_FLASH,
		Logger: opts.Logger,
	})
	if err != nil {
		panic(err)
	}

	return &Provider{
		client: client,
		logger: opts.Logger,
	}
}

func (p *Provider) GetCapabilities(ctx context.Context) (*core.Capabilities, error) {
	p.logger.Info("Fetching capabilities")

	// Placeholder for future implementation
	p.logger.Info("GetCapabilities method is not implemented yet")

	return nil, nil
}

func (p *Provider) UseModel(ctx context.Context, model *core.Model) error {
	p.logger.Info("Setting model", "modelID", model.ID)

	p.model = model

	return nil
}

// Generate implements the LLMProvider interface for basic responses
func (p *Provider) Generate(ctx context.Context, opts *core.GenerateOptions) (*core.Message, error) {
	p.logger.Info("Generate request received", "modelID", p.model.ID)

	resp, err := p.client.Chat(ctx, &client.ChatRequest{
		Model:    p.model.ID,
		Messages: opts.Messages,
		Tools:    opts.Tools,
	})

	if err != nil {
		p.logger.Error(err, "Error calling client chat method", err)
		return nil, fmt.Errorf("error calling client chat method: %w", err)
	}

	return &core.Message{
		Role:      core.AssistantMessageRole,
		Content:   resp.Message.Content,
		ToolCalls: resp.Message.ToolCalls,
	}, nil
}

// GenerateStream streams the response token by token
func (p *Provider) GenerateStream(ctx context.Context, opts *core.GenerateOptions) (<-chan *core.Message, <-chan string, <-chan error) {
	p.logger.Info("starting stream generation", "modelID", p.model.ID)

	return p.client.ChatStream(ctx, &client.ChatRequest{
		Model:    p.model.ID,
		Messages: opts.Messages,
		Tools:    opts.Tools,
	})
}
