package app

import (
	"context"

	"github.com/defryheryanto/ai-assistant/pkg/calendar"
	googlecalendar "github.com/defryheryanto/ai-assistant/pkg/calendar/google"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Services struct {
	CalendarService calendar.Service
	OpenAIClient    openai.Client
}

func SetupServices(ctx context.Context, params SetupToolsParams) (*Services, error) {
	calendarService, err := googlecalendar.New(ctx, params.GoogleCredentialsFilePath, params.GoogleTokenFilePath)
	if err != nil {
		return nil, err
	}

	client := openai.NewClient(
		option.WithAPIKey(params.OpenAIToken),
	)

	return &Services{
		CalendarService: calendarService,
		OpenAIClient:    client,
	}, nil
}
