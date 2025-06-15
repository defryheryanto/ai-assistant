package tools

import (
	"context"
	"encoding/json"

	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/internal/user"
	"github.com/tmc/langchaingo/llms"
)

type CreateUserTool struct {
	userService user.Service
}

func NewCreateUserTool(userService user.Service) *CreateUserTool {
	return &CreateUserTool{
		userService: userService,
	}
}

func (t *CreateUserTool) SystemPrompt() string {
	return "When creating a new user, you must collect all required information such as name, phone, email, and role. For the role, only the following constant values are valid: 'user' and 'admin'. Do not allow any other values for the user’s role. If the user attempts to assign a different role, inform them that only user and admin are permitted."
}

func (t *CreateUserTool) Definition() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "CreateUser",
			Description: "Create a new user based on the given information",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "The name of the new user",
					},
					"phone": map[string]any{
						"type":        "string",
						"description": "The phone number of the new user",
					},
					"role": map[string]any{
						"type":        "string",
						"description": "The role of the new user. A valid value will always be 'admin' or 'user'",
					},
					"email": map[string]any{
						"type":        "string",
						"description": "The email of the new user",
					},
				},
				"required": []string{"name", "phone", "email"},
			},
		},
	}
}

func (t *CreateUserTool) Execute(ctx context.Context, toolCall llms.ToolCall) (*llms.MessageContent, error) {
	// TODO(defryheryanto): Move this to a middleware-like func
	usr := contextgroup.GetUserContext(ctx)
	if usr == nil || usr.Role != string(user.RoleAdmin) {
		return &llms.MessageContent{
			Role: llms.ChatMessageTypeTool,
			Parts: []llms.ContentPart{
				llms.ToolCallResponse{
					ToolCallID: toolCall.ID,
					Name:       toolCall.FunctionCall.Name,
					Content:    "The requestor did not have permission to create user",
				},
			},
		}, nil
	}

	var args struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Role  string `json:"role"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
		return nil, err
	}

	_, err := t.userService.Create(ctx, user.CreateUserParams{
		Name:  args.Name,
		Phone: args.Phone,
		Role:  user.Role(args.Role),
		Email: args.Email,
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
				Content:    "User successfully created",
			},
		},
	}, nil
}
