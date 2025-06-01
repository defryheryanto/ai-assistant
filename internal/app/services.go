package app

import (
	"context"

	"github.com/defryheryanto/ai-assistant/pkg/calendar"
	googlecalendar "github.com/defryheryanto/ai-assistant/pkg/calendar/google"
)

type services struct {
	CalendarService calendar.Service
}

func setupServices(ctx context.Context, params SetupToolsParams) (*services, error) {
	calendarService, err := googlecalendar.New(ctx, params.GoogleCredentialsFilePath, params.GoogleTokenFilePath)
	if err != nil {
		return nil, err
	}

	return &services{
		CalendarService: calendarService,
	}, nil
}
