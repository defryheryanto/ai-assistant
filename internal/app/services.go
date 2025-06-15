package app

import (
	"context"

	"github.com/defryheryanto/ai-assistant/internal/calendar"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	pkgcalendar "github.com/defryheryanto/ai-assistant/pkg/calendar"
	googlecalendar "github.com/defryheryanto/ai-assistant/pkg/calendar/google"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Services struct {
	UserService          whatsapp.UserService
	WhatsAppGroupService whatsapp.GroupService
	CalendarService      pkgcalendar.Service
	OpenAIClient         openai.Client
}

func SetupServices(ctx context.Context, params SetupToolsParams) (*Services, error) {
	repositories := SetupRepository(ctx, params.DB)

	userService := whatsapp.NewUserService(repositories.UserRepository)

	whatsappGroupService := whatsapp.NewGroupService(repositories.WhatsAppGroupRepository)

	var calendarService pkgcalendar.Service
	var err error
	calendarService, err = googlecalendar.New(ctx, params.GoogleCredentialsFilePath, params.GoogleTokenFilePath)
	if err != nil {
		return nil, err
	}

	calendarService = calendar.New(calendarService, userService)

	client := openai.NewClient(
		option.WithAPIKey(params.OpenAIToken),
	)

	return &Services{
		UserService:          userService,
		WhatsAppGroupService: whatsappGroupService,
		CalendarService:      calendarService,
		OpenAIClient:         client,
	}, nil
}
