package tools

import (
	"context"
	"slices"

	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	"github.com/defryheryanto/ai-assistant/pkg/tools"
	"github.com/tmc/langchaingo/llms"
)

type RoleMiddlewareTool struct {
	base         tools.Tool
	allowedRoles []whatsapp.UserRole
}

func NewRoleMiddlewareTool(base tools.Tool, allowedRoles []whatsapp.UserRole) *RoleMiddlewareTool {
	return &RoleMiddlewareTool{
		base:         base,
		allowedRoles: allowedRoles,
	}
}

func (t *RoleMiddlewareTool) SystemPrompt() string {
	return t.base.SystemPrompt()
}

func (t *RoleMiddlewareTool) Definition() llms.Tool {
	return t.base.Definition()
}

func (t *RoleMiddlewareTool) Execute(ctx context.Context, toolCall llms.ToolCall) (*llms.MessageContent, error) {
	usr := contextgroup.GetUserContext(ctx)
	if usr == nil || !slices.Contains(t.allowedRoles, whatsapp.UserRole(usr.Role)) {
		return &llms.MessageContent{
			Role: llms.ChatMessageTypeTool,
			Parts: []llms.ContentPart{
				llms.ToolCallResponse{
					ToolCallID: toolCall.ID,
					Name:       toolCall.FunctionCall.Name,
					Content:    "The requestor did not have permission to do this action",
				},
			},
		}, nil
	}

	return t.base.Execute(ctx, toolCall)
}
