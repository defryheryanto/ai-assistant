package tools_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	whatsappmock "github.com/defryheryanto/ai-assistant/internal/whatsapp/mock"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp/tools"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"go.uber.org/mock/gomock"
)

func TestCreateUserTool_SystemPrompt(t *testing.T) {
	tool := tools.NewCreateUserTool(nil)
	expected := "When creating a new user, you must collect all required information such as name, phone, email, and role. For the role, only the following constant values are valid: 'user' and 'admin'. Do not allow any other values for the user’s role. If the user attempts to assign a different role, inform them that only user and admin are permitted."
	assert.Equal(t, expected, tool.SystemPrompt())
}

func TestCreateUserTool_Definition(t *testing.T) {
	tool := tools.NewCreateUserTool(nil)
	def := tool.Definition()

	assert.Equal(t, "function", def.Type)
	assert.NotNil(t, def.Function)
	if def.Function != nil {
		assert.Equal(t, "CreateUser", def.Function.Name)
		assert.Contains(t, def.Function.Description, "Create a new user")
		// Check params
		params, ok := def.Function.Parameters.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "object", params["type"])

		props, ok := params["properties"].(map[string]any)
		assert.True(t, ok)
		assert.Contains(t, props, "name")
		assert.Contains(t, props, "phone")
		assert.Contains(t, props, "role")
		assert.Contains(t, props, "email")

		// Optional: check required fields
		req, ok := params["required"].([]string)
		if !ok {
			// Sometimes JSON -> interface{} -> []any
			tmp, ok2 := params["required"].([]any)
			assert.True(t, ok2)
			assert.Len(t, tmp, 3)
		} else {
			assert.ElementsMatch(t, []string{"name", "phone", "email"}, req)
		}
	}
}

func TestCreateUserTool_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := whatsappmock.NewMockUserService(ctrl)
	tool := tools.NewCreateUserTool(mockUserService)

	validArgs := map[string]any{
		"name":  "Alice",
		"phone": "08123456789",
		"role":  "user",
		"email": "alice@email.com",
	}
	argBytes, _ := json.Marshal(validArgs)

	adminCtx := contextgroup.SetUserContext(context.Background(), &contextgroup.UserContext{
		ID:   1,
		Role: string(whatsapp.UserRoleAdmin),
	})
	nonAdminCtx := contextgroup.SetUserContext(context.Background(), &contextgroup.UserContext{
		ID:   2,
		Role: string(whatsapp.UserRoleUser),
	})
	nilUserCtx := context.Background()

	t.Run("permission denied - non-admin", func(t *testing.T) {
		resp, err := tool.Execute(nonAdminCtx, llms.ToolCall{
			ID: "noadmin",
			FunctionCall: &llms.FunctionCall{
				Name:      "CreateUser",
				Arguments: string(argBytes),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, llms.ChatMessageTypeTool, resp.Role)
		if assert.Len(t, resp.Parts, 1) {
			res, ok := resp.Parts[0].(llms.ToolCallResponse)
			assert.True(t, ok)
			assert.Equal(t, "noadmin", res.ToolCallID)
			assert.Contains(t, res.Content, "did not have permission")
		}
	})

	t.Run("permission denied - nil user", func(t *testing.T) {
		resp, err := tool.Execute(nilUserCtx, llms.ToolCall{
			ID: "niluser",
			FunctionCall: &llms.FunctionCall{
				Name:      "CreateUser",
				Arguments: string(argBytes),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, llms.ChatMessageTypeTool, resp.Role)
		if assert.Len(t, resp.Parts, 1) {
			res, ok := resp.Parts[0].(llms.ToolCallResponse)
			assert.True(t, ok)
			assert.Equal(t, "niluser", res.ToolCallID)
			assert.Contains(t, res.Content, "did not have permission")
		}
	})

	t.Run("success - admin context", func(t *testing.T) {
		mockUserService.EXPECT().
			Create(adminCtx, whatsapp.CreateUserParams{
				Name:  "Alice",
				Phone: "08123456789",
				Role:  whatsapp.UserRole("user"),
				Email: "alice@email.com",
			}).
			Return(int64(1), nil).
			Times(1)

		resp, err := tool.Execute(adminCtx, llms.ToolCall{
			ID: "123",
			FunctionCall: &llms.FunctionCall{
				Name:      "CreateUser",
				Arguments: string(argBytes),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, llms.ChatMessageTypeTool, resp.Role)
		found := false
		for _, part := range resp.Parts {
			if res, ok := part.(llms.ToolCallResponse); ok {
				found = true
				assert.Equal(t, "123", res.ToolCallID)
				assert.Equal(t, "CreateUser", res.Name)
				assert.Equal(t, "User successfully created", res.Content)
			}
		}
		assert.True(t, found, "ToolCallResponse not found in response parts")
	})

	t.Run("invalid json", func(t *testing.T) {
		resp, err := tool.Execute(adminCtx, llms.ToolCall{
			ID: "456",
			FunctionCall: &llms.FunctionCall{
				Name:      "CreateUser",
				Arguments: "{invalid-json}",
			},
		})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("userService error", func(t *testing.T) {
		mockUserService.EXPECT().
			Create(adminCtx, whatsapp.CreateUserParams{
				Name:  "Alice",
				Phone: "08123456789",
				Role:  whatsapp.UserRole("user"),
				Email: "alice@email.com",
			}).
			Return(int64(0), errors.New("create error")).
			Times(1)

		resp, err := tool.Execute(adminCtx, llms.ToolCall{
			ID: "789",
			FunctionCall: &llms.FunctionCall{
				Name:      "CreateUser",
				Arguments: string(argBytes),
			},
		})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
