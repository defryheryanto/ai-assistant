package tools

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type Tool interface {
	Definition() llms.Tool
	Execute(ctx context.Context, toolCall llms.ToolCall) (*llms.MessageContent, error)
}

type Registry interface {
	Register(tool Tool)
	GetTools() []llms.Tool
	Execute(ctx context.Context, inquiry string) (string, error)
}

type registry struct {
	llm           llms.Model
	toolFunctions []Tool
	enableLog     bool
}

func NewRegistry(llm llms.Model, enableLog bool) *registry {
	return &registry{
		llm:           llm,
		toolFunctions: []Tool{},
		enableLog:     enableLog,
	}
}

func (r *registry) Register(tool Tool) {
	r.toolFunctions = append(r.toolFunctions, tool)
}

func (r *registry) GetTools() []llms.Tool {
	toolDefinitions := make([]llms.Tool, 0, len(r.toolFunctions))

	for _, t := range r.toolFunctions {
		toolDefinitions = append(toolDefinitions, t.Definition())
	}

	return toolDefinitions
}

func (r *registry) executeTool(ctx context.Context, messageHistory []llms.MessageContent, resp *llms.ContentResponse) ([]llms.MessageContent, error) {
	for _, choice := range resp.Choices {
		for _, toolCall := range choice.ToolCalls {
			assistantResponse := llms.MessageContent{
				Role: llms.ChatMessageTypeAI,
				Parts: []llms.ContentPart{
					llms.ToolCall{
						ID:   toolCall.ID,
						Type: toolCall.Type,
						FunctionCall: &llms.FunctionCall{
							Name:      toolCall.FunctionCall.Name,
							Arguments: toolCall.FunctionCall.Arguments,
						},
					},
				},
			}
			messageHistory = append(messageHistory, assistantResponse)

			isToolFound := false
			for _, toolFunc := range r.toolFunctions {
				if toolCall.FunctionCall.Name == toolFunc.Definition().Function.Name {
					r.log(fmt.Sprintf("calling tool %s", toolCall.FunctionCall.Name))
					toolResponse, err := toolFunc.Execute(ctx, toolCall)
					if err != nil {
						return nil, err
					}
					r.log(fmt.Sprintf("tool %s called with response %s", toolCall.FunctionCall.Name, toolResponse))

					messageHistory = append(messageHistory, *toolResponse)
					isToolFound = true
				}
			}
			if !isToolFound {
				return nil, fmt.Errorf("unsupported tool: %s", toolCall.FunctionCall.Name)
			}
		}
	}

	return messageHistory, nil
}

func (r *registry) Execute(ctx context.Context, inquiry string) (string, error) {
	r.log(fmt.Sprintf("processing inquiry: %s", inquiry))
	messageHistory := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a helpful AI assistant. Whenever you use a tool and the tool returns a link (such as a calendar invite or external resource), you must always clearly include that link in your response to the user. If no link is returned, answer as usual."),
		llms.TextParts(llms.ChatMessageTypeHuman, inquiry),
	}

	for {
		resp, err := r.llm.GenerateContent(ctx, messageHistory, llms.WithTools(r.GetTools()))
		if err != nil {
			return "", err
		}
		if len(resp.Choices[0].ToolCalls) == 0 {
			r.log("inquiry process done")
			return resp.Choices[0].Content, nil
		}

		messageHistory, err = r.executeTool(ctx, messageHistory, resp)
		if err != nil {
			return "", err
		}
	}
}
