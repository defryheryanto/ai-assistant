package app

import (
	"context"

	"github.com/defryheryanto/ai-assistant/internal/user"
	"github.com/defryheryanto/ai-assistant/pkg/calendar"
	googlecalendar "github.com/defryheryanto/ai-assistant/pkg/calendar/google"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Services struct {
	UserService     user.Service
	CalendarService calendar.Service
	OpenAIClient    openai.Client
}

func SetupServices(ctx context.Context, params SetupToolsParams) (*Services, error) {
	repositories := SetupRepository(ctx, params.DB)

	userService := user.NewService(repositories.UserRepository)

	calendarService, err := googlecalendar.New(ctx, params.GoogleCredentialsFilePath, params.GoogleTokenFilePath)
	if err != nil {
		return nil, err
	}

	client := openai.NewClient(
		option.WithAPIKey(params.OpenAIToken),
	)

	return &Services{
		UserService:     userService,
		CalendarService: calendarService,
		OpenAIClient:    client,
	}, nil
}
