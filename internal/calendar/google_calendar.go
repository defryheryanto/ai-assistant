package calendar

import (
	"context"

	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	"github.com/defryheryanto/ai-assistant/pkg/calendar"
)

type GoogleCalendarService struct {
	baseService calendar.Service
	userService whatsapp.UserService
}

func New(baseService calendar.Service, userService whatsapp.UserService) *GoogleCalendarService {
	return &GoogleCalendarService{
		baseService: baseService,
		userService: userService,
	}
}

func (s *GoogleCalendarService) CreateEvent(ctx context.Context, params calendar.CreateEventParams) (string, error) {
	usr := contextgroup.GetUserContext(ctx)
	if usr == nil {
		return s.baseService.CreateEvent(ctx, params)
	}

	res, err := s.userService.GetByJID(ctx, usr.WhatsAppJID)
	if err != nil {
		return "", err
	}
	if res == nil {
		return s.baseService.CreateEvent(ctx, params)
	}

	params.Attendees = append(params.Attendees, res.Email)

	return s.baseService.CreateEvent(ctx, params)
}
