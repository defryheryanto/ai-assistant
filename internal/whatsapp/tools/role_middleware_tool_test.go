package tools_test

import (
	"context"
	"testing"

	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp/tools"
	"github.com/defryheryanto/ai-assistant/pkg/tools/mock"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"go.uber.org/mock/gomock"
)

func TestRoleMiddlewareTool_SystemPrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTool := mock.NewMockTool(ctrl)
	mockTool.EXPECT().SystemPrompt().Return("mock prompt").Times(1)

	wrapped := tools.NewRoleMiddlewareTool(mockTool, []whatsapp.UserRole{whatsapp.UserRoleAdmin})
	assert.Equal(t, "mock prompt", wrapped.SystemPrompt())
}

func TestRoleMiddlewareTool_Definition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTool := mock.NewMockTool(ctrl)
	mockTool.EXPECT().Definition().Return(llms.Tool{Type: "function"}).Times(1)

	wrapped := tools.NewRoleMiddlewareTool(mockTool, []whatsapp.UserRole{whatsapp.UserRoleAdmin})
	def := wrapped.Definition()
	assert.Equal(t, "function", def.Type)
}

func TestRoleMiddlewareTool_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTool := mock.NewMockTool(ctrl)
	toolCall := llms.ToolCall{
		ID: "123",
		FunctionCall: &llms.FunctionCall{
			Name: "SomeTool",
		},
	}
	resp := &llms.MessageContent{
		Role: llms.ChatMessageTypeTool,
		Parts: []llms.ContentPart{
			llms.ToolCallResponse{
				ToolCallID: "123",
				Name:       "SomeTool",
				Content:    "Executed successfully",
			},
		},
	}

	t.Run("allowed role", func(t *testing.T) {
		adminCtx := contextgroup.SetUserContext(context.Background(), &contextgroup.UserContext{
			ID:   1,
			Role: string(whatsapp.UserRoleAdmin),
		})
		mockTool.EXPECT().Execute(adminCtx, toolCall).Return(resp, nil).Times(1)

		wrapped := tools.NewRoleMiddlewareTool(mockTool, []whatsapp.UserRole{whatsapp.UserRoleAdmin})
		out, err := wrapped.Execute(adminCtx, toolCall)
		assert.NoError(t, err)
		assert.Equal(t, resp, out)
	})

	t.Run("disallowed role", func(t *testing.T) {
		userCtx := contextgroup.SetUserContext(context.Background(), &contextgroup.UserContext{
			ID:   2,
			Role: string(whatsapp.UserRoleUser),
		})

		wrapped := tools.NewRoleMiddlewareTool(mockTool, []whatsapp.UserRole{whatsapp.UserRoleAdmin})
		out, err := wrapped.Execute(userCtx, toolCall)
		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.Equal(t, "The requestor did not have permission to do this action", out.Parts[0].(llms.ToolCallResponse).Content)
	})

	t.Run("nil user context", func(t *testing.T) {
		wrapped := tools.NewRoleMiddlewareTool(mockTool, []whatsapp.UserRole{whatsapp.UserRoleAdmin})
		out, err := wrapped.Execute(context.Background(), toolCall)
		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.Equal(t, "The requestor did not have permission to do this action", out.Parts[0].(llms.ToolCallResponse).Content)
	})
}
