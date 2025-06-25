package calendar_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/defryheryanto/ai-assistant/pkg/calendar"
	mock "github.com/defryheryanto/ai-assistant/pkg/calendar/mock"
	toolpkg "github.com/defryheryanto/ai-assistant/pkg/tools/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"go.uber.org/mock/gomock"
)

func TestCreateEventTool_SystemPrompt(t *testing.T) {
	tool := toolpkg.NewCreateEventTool(nil, false)
	assert.Equal(t, "", tool.SystemPrompt())
}

func TestCreateEventTool_Definition(t *testing.T) {
	tool := toolpkg.NewCreateEventTool(nil, false)
	def := tool.Definition()
	assert.Equal(t, "function", def.Type)
	assert.NotNil(t, def.Function)
	assert.Equal(t, "CreateCalendarEvent", def.Function.Name)
}

func TestCreateEventTool_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockService := mock.NewMockService(ctrl)
		tool := toolpkg.NewCreateEventTool(mockService, true)
		args := map[string]interface{}{
			"summary":     "Test Event",
			"description": "desc",
			"location":    "loc",
			"start":       "2025-06-25T10:00:00+07:00",
			"end":         "2025-06-25T11:00:00+07:00",
			"attendees":   []string{"a@email.com"},
		}
		argBytes, _ := json.Marshal(args)
		toolCall := llms.ToolCall{
			ID: "1",
			FunctionCall: &llms.FunctionCall{
				Name:      "CreateCalendarEvent",
				Arguments: string(argBytes),
			},
		}
		mockService.EXPECT().CreateEvent(gomock.Any(), gomock.AssignableToTypeOf(calendar.CreateEventParams{})).DoAndReturn(
			func(ctx context.Context, params calendar.CreateEventParams) (string, error) {
				assert.True(t, params.IsCreatorAttendee)
				assert.Equal(t, "Test Event", params.Summary)
				return "http://event-link", nil
			},
		)
		msg, err := tool.Execute(context.Background(), toolCall)
		assert.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, llms.ChatMessageTypeTool, msg.Role)
		assert.Contains(t, fmt.Sprintf("%v", msg.Parts), "http://event-link")
	})

	t.Run("invalid json", func(t *testing.T) {
		mockService := mock.NewMockService(ctrl)
		tool := toolpkg.NewCreateEventTool(mockService, false)
		toolCall := llms.ToolCall{
			ID: "2",
			FunctionCall: &llms.FunctionCall{
				Name:      "CreateCalendarEvent",
				Arguments: "not-json",
			},
		}
		msg, err := tool.Execute(context.Background(), toolCall)
		assert.Error(t, err)
		assert.Nil(t, msg)
	})

	t.Run("calendar service error", func(t *testing.T) {
		mockService := mock.NewMockService(ctrl)
		tool := toolpkg.NewCreateEventTool(mockService, false)
		args := map[string]interface{}{
			"summary": "Test Event",
			"start":   "2025-06-25T10:00:00+07:00",
			"end":     "2025-06-25T11:00:00+07:00",
		}
		argBytes, _ := json.Marshal(args)
		toolCall := llms.ToolCall{
			ID: "3",
			FunctionCall: &llms.FunctionCall{
				Name:      "CreateCalendarEvent",
				Arguments: string(argBytes),
			},
		}
		mockService.EXPECT().CreateEvent(gomock.Any(), gomock.AssignableToTypeOf(calendar.CreateEventParams{})).Return("", errors.New("calendar error"))
		msg, err := tool.Execute(context.Background(), toolCall)
		assert.Error(t, err)
		assert.Nil(t, msg)
	})
}
