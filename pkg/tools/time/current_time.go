package time

import (
	"context"
	"time"

	"github.com/tmc/langchaingo/llms"
)

type CurrentTimeTool struct{}

func NewCurrentTimeTool() *CurrentTimeTool {
	return &CurrentTimeTool{}
}

func (t *CurrentTimeTool) Definition() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "GetCurrentTime",
			Description: "Get the current time for any time calculating",
		},
	}
}

func (t *CurrentTimeTool) Execute(ctx context.Context, toolCall llms.ToolCall) (*llms.MessageContent, error) {
	return &llms.MessageContent{
		Role: llms.ChatMessageTypeTool,
		Parts: []llms.ContentPart{
			llms.ToolCallResponse{
				ToolCallID: toolCall.ID,
				Name:       toolCall.FunctionCall.Name,
				Content:    time.Now().Format(time.RFC3339Nano),
			},
		},
	}, nil
}
