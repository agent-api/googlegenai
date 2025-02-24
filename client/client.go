package client

import (
	"context"
	"log/slog"
	"os"

	"github.com/agent-api/core/types"

	googGenAI "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	//opts []option.RequestOption

	client *googGenAI.Client

	model string

	logger *slog.Logger
}

type GeminiClientOpts struct {
	Logger *slog.Logger
	Model  *types.Model
}

func NewClient(ctx context.Context, opts *GeminiClientOpts) (*GeminiClient, error) {
	client, err := googGenAI.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}

	return &GeminiClient{
		client: client,
		model:  opts.Model.ID,
		logger: opts.Logger,
	}, nil
}

func (c *GeminiClient) Done() {
	c.client.Close()
}

// Convert your Tool to OpenAI's ChatCompletionToolParam
//func ToOpenAIToolParam(t *types.Tool) (*openai.ChatCompletionToolParam, error) {
//var schemaMap map[string]interface{}
//if err := json.Unmarshal(t.JSONSchema, &schemaMap); err != nil {
//return nil, err
//}

//return &openai.ChatCompletionToolParam{
//Type: openai.F(openai.ChatCompletionToolTypeFunction),
//Function: openai.F(openai.FunctionDefinitionParam{
//Name:        openai.String(t.Name),
//Description: openai.String(t.Description),
//Parameters:  openai.F(openai.FunctionParameters(schemaMap)),
//}),
//}, nil
//}

func (c *GeminiClient) Chat(ctx context.Context, req *ChatRequest) (ChatResponse, error) {
	// TODO - need to handle adding to Gemini history
	//geminiMessages := []*googGenAI.Content{
	//{
	//Parts: []googGenAI.Part{
	//googGenAI.Text(req.Messages[0].Content),
	//},
	//Role: "user",
	//},
	//}

	// TODO - need to handle multiple messages

	// TODO - need to handle tools

	model := c.client.GenerativeModel(req.Model)
	cs := model.StartChat()
	res, err := cs.SendMessage(ctx, googGenAI.Text(req.Messages[0].Content))
	if err != nil {
		panic(err)
	}

	// big hack, gross
	content, ok := res.Candidates[0].Content.Parts[0].(googGenAI.Text)
	if !ok {
		panic("not ok")
	}

	return ChatResponse{
		Message: types.Message{
			Content: string(content),
		},
	}, nil

}
