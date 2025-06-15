package tools_test

import (
	"context"
	"errors"
	"testing"

	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp/mock"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp/tools"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"go.uber.org/mock/gomock"
)

func TestRegisterGroupTool_SystemPrompt(t *testing.T) {
	tool := tools.NewRegisterGroupTool(nil)
	assert.Equal(t, "", tool.SystemPrompt())
}

func TestRegisterGroupTool_Definition(t *testing.T) {
	tool := tools.NewRegisterGroupTool(nil)
	def := tool.Definition()

	assert.Equal(t, "function", def.Type)
	assert.NotNil(t, def.Function)
	if def.Function != nil {
		assert.Equal(t, "RegisterGroup", def.Function.Name)
		assert.Contains(t, def.Function.Description, "Register a group")
	}
}

func TestRegisterGroupTool_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGroupService := mock.NewMockGroupService(ctrl)
	tool := tools.NewRegisterGroupTool(mockGroupService)

	validCtx := contextgroup.SetWhatsAppContext(context.Background(), &contextgroup.WhatsAppContext{
		CurrentChatJID: "12345@g.us",
		SenderJID:      "6281234567890@s.whatsapp.net",
	})
	nilCtx := context.Background()

	t.Run("missing whatsapp context", func(t *testing.T) {
		resp, err := tool.Execute(nilCtx, llms.ToolCall{
			ID: "missingCtx",
			FunctionCall: &llms.FunctionCall{
				Name: "RegisterGroup",
			},
		})
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "whatsapp context is null")
	})

	t.Run("success", func(t *testing.T) {
		mockGroupService.EXPECT().
			Create(validCtx, whatsapp.CreateGroupParams{
				GroupJID:     "12345@g.us",
				RegisteredBy: "6281234567890@s.whatsapp.net",
			}).
			Return(int64(42), nil).
			Times(1)

		resp, err := tool.Execute(validCtx, llms.ToolCall{
			ID: "success",
			FunctionCall: &llms.FunctionCall{
				Name: "RegisterGroup",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, llms.ChatMessageTypeTool, resp.Role)

		res, ok := resp.Parts[0].(llms.ToolCallResponse)
		assert.True(t, ok)
		assert.Equal(t, "RegisterGroup", res.Name)
		assert.Equal(t, "Group successfully created", res.Content)
	})

	t.Run("groupService error", func(t *testing.T) {
		mockGroupService.EXPECT().
			Create(validCtx, whatsapp.CreateGroupParams{
				GroupJID:     "12345@g.us",
				RegisteredBy: "6281234567890@s.whatsapp.net",
			}).
			Return(int64(0), errors.New("create failed")).
			Times(1)

		resp, err := tool.Execute(validCtx, llms.ToolCall{
			ID: "failure",
			FunctionCall: &llms.FunctionCall{
				Name: "RegisterGroup",
			},
		})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
