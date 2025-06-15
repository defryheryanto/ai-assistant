package calendar_test

import (
	"context"
	"errors"
	"testing"

	"github.com/defryheryanto/ai-assistant/internal/calendar"
	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	whatsappmock "github.com/defryheryanto/ai-assistant/internal/whatsapp/mock"
	pkgcalendar "github.com/defryheryanto/ai-assistant/pkg/calendar"
	calMock "github.com/defryheryanto/ai-assistant/pkg/calendar/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGoogleCalendarService_CreateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCal := calMock.NewMockService(ctrl)
	mockUser := whatsappmock.NewMockUserService(ctrl)

	svc := calendar.New(mockCal, mockUser)

	ctx := context.Background()
	params := pkgcalendar.CreateEventParams{
		Summary: "Test Event",
	}
	baseResult := "base-event-id"

	t.Run("User not in context", func(t *testing.T) {
		mockCal.EXPECT().
			CreateEvent(ctx, params).
			Return(baseResult, nil).
			Times(1)

		result, err := svc.CreateEvent(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, baseResult, result)
	})

	t.Run("User in context, GetUserByWhatsAppJID error", func(t *testing.T) {
		usr := &contextgroup.UserContext{WhatsAppJID: "jid1"}
		ctxWithUser := contextgroup.SetUserContext(ctx, usr)

		mockUser.EXPECT().
			GetByJID(ctxWithUser, usr.WhatsAppJID).
			Return(nil, errors.New("user error")).
			Times(1)

		result, err := svc.CreateEvent(ctxWithUser, params)
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("User in context, not found in DB", func(t *testing.T) {
		usr := &contextgroup.UserContext{WhatsAppJID: "jid1"}
		ctxWithUser := contextgroup.SetUserContext(ctx, usr)

		mockUser.EXPECT().
			GetByJID(ctxWithUser, usr.WhatsAppJID).
			Return(nil, nil).
			Times(1)
		mockCal.EXPECT().
			CreateEvent(ctxWithUser, params).
			Return(baseResult, nil).
			Times(1)

		result, err := svc.CreateEvent(ctxWithUser, params)
		assert.NoError(t, err)
		assert.Equal(t, baseResult, result)
	})

	t.Run("User in context, found in DB", func(t *testing.T) {
		usr := &contextgroup.UserContext{WhatsAppJID: "jid1"}
		ctxWithUser := contextgroup.SetUserContext(ctx, usr)
		dbUser := &whatsapp.User{
			Email: "test@email.com",
		}

		// The expected params must include the attendee
		expectedParams := params
		expectedParams.Attendees = append(expectedParams.Attendees, dbUser.Email)

		mockUser.EXPECT().
			GetByJID(ctxWithUser, usr.WhatsAppJID).
			Return(dbUser, nil).
			Times(1)
		mockCal.EXPECT().
			CreateEvent(ctxWithUser, expectedParams).
			Return("final-event-id", nil).
			Times(1)

		result, err := svc.CreateEvent(ctxWithUser, params)
		assert.NoError(t, err)
		assert.Equal(t, "final-event-id", result)
	})
}
