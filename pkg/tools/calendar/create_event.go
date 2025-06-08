package calendar

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/defryheryanto/ai-assistant/pkg/calendar"
	"github.com/tmc/langchaingo/llms"
)

type CreateEventTool struct {
	calendarService calendar.Service
}

func NewCreateEventTool(calendarService calendar.Service) *CreateEventTool {
	return &CreateEventTool{
		calendarService: calendarService,
	}
}

func (t *CreateEventTool) Definition() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "CreateCalendarEvent",
			Description: "Create a Google Calendar event with specified details and attendees.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"summary": map[string]any{
						"type":        "string",
						"description": "Title or summary of the event.",
					},
					"description": map[string]any{
						"type":        "string",
						"description": "A longer description for the event.",
					},
					"location": map[string]any{
						"type":        "string",
						"description": "Location of the event.",
					},
					"start": map[string]any{
						"type":        "string",
						"format":      "date-time",
						"description": "Start time in ISO 8601 format, e.g. 2024-05-29T09:00:00+07:00",
					},
					"end": map[string]any{
						"type":        "string",
						"format":      "date-time",
						"description": "End time in ISO 8601 format, e.g. 2024-05-29T10:00:00+07:00",
					},
					"attendees": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type":        "string",
							"description": "Email address of an attendee",
						},
						"description": "A list of emails for event attendees.",
					},
				},
				"required": []string{"summary", "start", "end"},
			},
		},
	}
}

func (t *CreateEventTool) Execute(ctx context.Context, toolCall llms.ToolCall) (*llms.MessageContent, error) {
	var args calendar.CreateEventParams
	if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
		return nil, err
	}

	eventLink, err := t.calendarService.CreateEvent(ctx, args)
	if err != nil {
		return nil, err
	}

	return &llms.MessageContent{
		Role: llms.ChatMessageTypeTool,
		Parts: []llms.ContentPart{
			llms.ToolCallResponse{
				ToolCallID: toolCall.ID,
				Name:       toolCall.FunctionCall.Name,
				Content:    fmt.Sprintf(`{"message": "Calendar invite successfully created.", "link": "%s"}`, eventLink),
			},
		},
	}, nil
}
