package tools

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

//go:generate mockgen -source registry.go -package mock -destination mock/mock.go

// ContextWindowManager is an interface to manage conversational history.
type ContextWindowManager interface {
	GetHistory(ctx context.Context, id string) ([]llms.MessageContent, error)
	SaveHistory(ctx context.Context, id string, history []llms.MessageContent) error
}

// Tool defines a callable function that can be registered with the registry.
type Tool interface {
	Definition() llms.Tool
	SystemPrompt() string
	Execute(ctx context.Context, toolCall llms.ToolCall) (*llms.MessageContent, error)
}

// Registry provides access to registered tools and executes them when needed.
type Registry interface {
	Register(tool Tool)
	GetTools() []llms.Tool
	Execute(ctx context.Context, contextID string, inquiry string) (string, error)
}

type registry struct {
	llm                  llms.Model
	toolFunctions        []Tool
	enableLog            bool
	systemPrompt         string
	contextWindowManager ContextWindowManager
}

func NewRegistry(llm llms.Model, options ...Option) *registry {
	r := &registry{
		llm:                  llm,
		toolFunctions:        []Tool{},
		enableLog:            false,
		systemPrompt:         defaultSystemPrompt,
		contextWindowManager: nil,
	}

	for _, opt := range options {
		opt(r)
	}

	return r
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

func (r *registry) generatePrompts() []llms.MessageContent {
	messageHistory := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, r.systemPrompt),
	}

	for _, tf := range r.toolFunctions {
		if tf.SystemPrompt() != "" {
			messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeSystem, tf.SystemPrompt()))
		}
	}

	return messageHistory
}

func (r *registry) Execute(ctx context.Context, contextID string, inquiry string) (string, error) {
	r.log(fmt.Sprintf("processing inquiry: %s", inquiry))
	basePrompts := r.generatePrompts()
	messageHistory := make([]llms.MessageContent, len(basePrompts))
	copy(messageHistory, basePrompts)

	var history []llms.MessageContent

	if r.contextWindowManager != nil && contextID != "" {
		history, err := r.contextWindowManager.GetHistory(ctx, contextID)
		if err != nil {
			return "", err
		}
		messageHistory = append(messageHistory, history...)
	}

	messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeHuman, inquiry))

	for {
		resp, err := r.llm.GenerateContent(ctx, messageHistory, llms.WithTools(r.GetTools()))
		if err != nil {
			return "", err
		}
		if len(resp.Choices[0].ToolCalls) == 0 {
			r.log("inquiry process done")
			messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content))
			if r.contextWindowManager != nil && contextID != "" {
				convHistory := messageHistory[len(basePrompts)+len(history):]
				_ = r.contextWindowManager.SaveHistory(ctx, contextID, convHistory)
			}
			return resp.Choices[0].Content, nil
		}

		messageHistory, err = r.executeTool(ctx, messageHistory, resp)
		if err != nil {
			return "", err
		}
	}
}
