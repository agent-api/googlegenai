package googlegenai

import (
	"context"

	"github.com/go-logr/logr"

	"github.com/agent-api/core"
	"github.com/agent-api/googlegenai/client"
	"github.com/agent-api/googlegenai/models"
)

type Embedder struct {
	host string
	port int

	model *core.Model

	client *client.GoogleGenAIClient

	logger *logr.Logger
}

type EmbedderOpts struct {
	BaseURL string
	Port    int
	APIKey  string

	Logger *logr.Logger
}

func NewEmbedder(opts *EmbedderOpts) *Embedder {
	ctx := context.Background()
	opts.Logger.Info("creating a new google generative ai provider")

	client, err := client.NewClient(ctx, &client.GoogleGenAIClientOpts{
		Model:  models.GEMINI_1_5_FLASH,
		Logger: opts.Logger,
	})
	if err != nil {
		panic(err)
	}

	return &Embedder{
		client: client,
		logger: opts.Logger,
	}
}

func (e *Embedder) GenerateEmbedding(ctx context.Context, content string) (*core.Embedding, error) {
	e.logger.Info("generating vector embedding")

	return e.client.Vector(ctx, content)
}
