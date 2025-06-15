package tools

import (
	"context"
	"fmt"

	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	"github.com/tmc/langchaingo/llms"
)

type RegisterGroupTool struct {
	groupService whatsapp.GroupService
}

func NewRegisterGroupTool(groupService whatsapp.GroupService) *RegisterGroupTool {
	return &RegisterGroupTool{
		groupService: groupService,
	}
}

func (t *RegisterGroupTool) SystemPrompt() string {
	return ""
}

func (t *RegisterGroupTool) Definition() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "RegisterGroup",
			Description: "Register a group for access",
		},
	}
}

func (t *RegisterGroupTool) Execute(ctx context.Context, toolCall llms.ToolCall) (*llms.MessageContent, error) {
	whatsappContext := contextgroup.GetWhatsAppContext(ctx)
	if whatsappContext == nil {
		return nil, fmt.Errorf("whatsapp context is null")
	}

	_, err := t.groupService.Create(ctx, whatsapp.CreateGroupParams{
		GroupJID:     whatsappContext.CurrentChatJID,
		RegisteredBy: whatsappContext.SenderJID,
	})
	if err != nil {
		return nil, err
	}

	return &llms.MessageContent{
		Role: llms.ChatMessageTypeTool,
		Parts: []llms.ContentPart{
			llms.ToolCallResponse{
				ToolCallID: toolCall.ID,
				Name:       toolCall.FunctionCall.Name,
				Content:    "Group successfully created",
			},
		},
	}, nil
}
