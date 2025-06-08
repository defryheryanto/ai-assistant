package google

import (
	"context"
	"os"

	"github.com/defryheryanto/ai-assistant/pkg/calendar"
	"golang.org/x/oauth2/google"
	gcalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

type GoogleCalendarService struct {
	calendarService *gcalendar.Service
	peopleService   *people.Service
}

func New(ctx context.Context, credentialsFilePath string, tokenFilePath string) (*GoogleCalendarService, error) {
	b, err := os.ReadFile(credentialsFilePath)
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(
		b,
		gcalendar.CalendarEventsScope,
		people.UserinfoProfileScope,
		people.UserinfoEmailScope,
	)
	if err != nil {
		return nil, err
	}

	c, err := getClient(config, tokenFilePath)
	if err != nil {
		return nil, err
	}

	calendarService, err := gcalendar.NewService(ctx, option.WithHTTPClient(c))
	if err != nil {
		return nil, err
	}

	peopleService, err := people.NewService(context.Background(), option.WithHTTPClient(c))
	if err != nil {
		return nil, err
	}

	return &GoogleCalendarService{
		calendarService: calendarService,
		peopleService:   peopleService,
	}, nil
}

func (s *GoogleCalendarService) getUserEmail() (string, error) {
	person, err := s.peopleService.People.Get("people/me").PersonFields("emailAddresses").Do()
	if err != nil {
		return "", err
	}
	if len(person.EmailAddresses) > 0 {
		return person.EmailAddresses[0].Value, nil
	}

	return "", nil
}

func (s *GoogleCalendarService) CreateEvent(ctx context.Context, params calendar.CreateEventParams) (string, error) {
	attendees := []*gcalendar.EventAttendee{}
	if params.IsCreatorAttendee {
		userEmail, err := s.getUserEmail()
		if err != nil {
			return "", err
		}

		attendees = append(attendees, &gcalendar.EventAttendee{Email: userEmail})
	}

	event := &gcalendar.Event{
		Summary:     params.Summary,
		Location:    params.Location,
		Description: params.Description,
		Start: &gcalendar.EventDateTime{
			DateTime: params.Start,
			TimeZone: "Asia/Jakarta",
		},
		End: &gcalendar.EventDateTime{
			DateTime: params.End,
			TimeZone: "Asia/Jakarta",
		},
		Attendees: attendees,
	}

	for _, email := range params.Attendees {
		event.Attendees = append(event.Attendees, &gcalendar.EventAttendee{Email: email})
	}

	createdEvent, err := s.calendarService.Events.Insert("primary", event).Do()
	if err != nil {
		return "", err
	}

	return createdEvent.HtmlLink, nil
}
