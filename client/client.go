package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/agent-api/core"
	"github.com/go-logr/logr"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GoogleGenAIClient struct {
	//opts []option.RequestOption

	client *genai.Client

	model string

	logger *logr.Logger
}

type GoogleGenAIClientOpts struct {
	Logger *logr.Logger
	Model  *core.Model
}

func NewClient(ctx context.Context, opts *GoogleGenAIClientOpts) (*GoogleGenAIClient, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}

	return &GoogleGenAIClient{
		client: client,
		model:  opts.Model.ID,
		logger: opts.Logger,
	}, nil
}

func (c *GoogleGenAIClient) Done() {
	c.client.Close()
}

func (c *GoogleGenAIClient) Chat(ctx context.Context, req *ChatRequest) (ChatResponse, error) {
	tools := []*genai.Tool{}
	for _, tool := range req.Tools {
		// since the agent-api core's tool params are a slice of bytes and ASSUMED
		// to be valid JSON schema: attempt to unmarshal wholesale
		schema := &genai.Schema{}
		err := json.Unmarshal(tool.JSONSchema, schema)
		if err != nil {
			panic(err)
		}

		tools = append(tools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        tool.Name,
					Description: tool.Description,
					Parameters:  schema,
				},
			},
		})
	}

	model := c.client.GenerativeModel(req.Model)
	model.Tools = tools

	history := []*genai.Content{}
	message := ""
	for i, m := range req.Messages {
		if i == len(req.Messages)-1 {
			message = m.Content
			continue
		}

		role := ""
		switch m.Role {
		case core.UserMessageRole:
			role = "user"

		case core.AssistantMessageRole:
			role = "model"
		}

		history = append(history, &genai.Content{
			Parts: []genai.Part{
				genai.Text(m.Content),
			},
			Role: role,
		})
	}

	cs := model.StartChat()
	cs.History = history

	res, err := cs.SendMessage(ctx, genai.Text(message))
	if err != nil {
		panic(err)
	}

	// big hack, gross
	content, ok := res.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		panic("not ok")
	}

	return ChatResponse{
		Message: core.Message{
			Content: string(content),
		},
	}, nil
}

func (c *GoogleGenAIClient) ChatStream(ctx context.Context, req *ChatRequest) (<-chan *core.Message, <-chan string, <-chan error) {
	c.logger.V(1).Info("received chat stream message request")

	// TODO - need to handle more messages

	// TODO - handle tools

	msgChan := make(chan *core.Message)
	deltaChan := make(chan string)
	errChan := make(chan error, 1)

	c.logger.V(1).Info("kicking async go func for chat stream")

	model := c.client.GenerativeModel(req.Model)
	cs := model.StartChat()

	go func() {
		defer close(msgChan)
		defer close(deltaChan)
		defer close(errChan)

		iter := cs.SendMessageStream(ctx, genai.Text(req.Messages[0].Content))
		for {
			resp, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				errChan <- fmt.Errorf("google gen ai streaming iterator error: %w", err)
			}

			delta, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
			if !ok {
				panic("not ok")
			}

			deltaChan <- string(delta)

			// TODO - need to accumulate deltas for message
		}
	}()

	return msgChan, deltaChan, errChan
}

func (c *GoogleGenAIClient) Vector(ctx context.Context, content string) (*core.Embedding, error) {
	em := c.client.EmbeddingModel("embedding-001")
	res, err := em.EmbedContent(ctx, genai.Text(content))

	if err != nil {
		return nil, err
	}

	return &core.Embedding{
		ID:      "woof",
		Vector:  res.Embedding.Values,
		Content: content,
	}, nil
}
