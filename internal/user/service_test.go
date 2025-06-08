package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/defryheryanto/ai-assistant/internal/user"
	"github.com/defryheryanto/ai-assistant/internal/user/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_GetUserByWhatsAppJID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	svc := user.NewService(mockRepo)

	ctx := context.Background()
	jid := "123@whatsapp.net"
	expectedUser := &user.User{
		ID:   1,
		Name: "Alice",
	}

	t.Run("found", func(t *testing.T) {
		mockRepo.EXPECT().
			FindUserByWhatsAppJID(ctx, jid).
			Return(expectedUser, nil).Times(1)

		result, err := svc.GetUserByWhatsAppJID(ctx, jid)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().
			FindUserByWhatsAppJID(ctx, jid).
			Return(nil, nil).Times(1)

		result, err := svc.GetUserByWhatsAppJID(ctx, jid)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("repo error", func(t *testing.T) {
		repoErr := errors.New("db error")
		mockRepo.EXPECT().
			FindUserByWhatsAppJID(ctx, jid).
			Return(nil, repoErr).Times(1)

		result, err := svc.GetUserByWhatsAppJID(ctx, jid)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, repoErr, err)
	})
}
