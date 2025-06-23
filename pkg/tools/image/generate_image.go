package image

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/tmc/langchaingo/llms"
)

// GenerateImageTool provides capability to generate image using OpenAI.
type GenerateImageTool struct {
	client openai.Client
}

// NewGenerateImageTool returns a new GenerateImageTool instance.
func NewGenerateImageTool(client openai.Client) *GenerateImageTool {
	return &GenerateImageTool{client: client}
}

// SystemPrompt returns additional prompt required by the tool.
func (t *GenerateImageTool) SystemPrompt() string {
	return ""
}

// Definition defines the function signature for this tool.
func (t *GenerateImageTool) Definition() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "GenerateImage",
			Description: "Generate an image using OpenAI based on the provided prompt",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"prompt": map[string]any{
						"type":        "string",
						"description": "Description of the desired image",
					},
				},
				"required": []string{"prompt"},
			},
		},
	}
}

// Execute calls the OpenAI API to generate an image.
func (t *GenerateImageTool) Execute(ctx context.Context, toolCall llms.ToolCall) (*llms.MessageContent, error) {
	var args struct {
		Prompt string `json:"prompt"`
	}
	if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
		return nil, err
	}

	resp, err := t.client.Images.Generations.New(ctx, openai.ImageGenerationNewParams{
		Prompt: args.Prompt,
		N:      1,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no image generated")
	}
	imageURL := resp.Data[0].URL

	return &llms.MessageContent{
		Role: llms.ChatMessageTypeTool,
		Parts: []llms.ContentPart{
			llms.ToolCallResponse{
				ToolCallID: toolCall.ID,
				Name:       toolCall.FunctionCall.Name,
				Content:    imageURL,
			},
		},
	}, nil
}
