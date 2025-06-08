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

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	svc := user.NewService(mockRepo)
	ctx := context.Background()

	// base params
	baseParams := user.CreateUserParams{
		Name:  "Alice",
		Phone: "12345678",
		Email: "alice@email.com",
	}

	t.Run("role is empty, should default to RoleUser", func(t *testing.T) {
		params := baseParams
		params.Role = ""
		expectedUser := &user.User{
			Name:        params.Name,
			WhatsAppJID: params.Phone + "@s.whatsapp.net",
			Role:        user.RoleUser,
			Email:       params.Email,
		}
		mockRepo.EXPECT().
			Insert(ctx, expectedUser).
			Return(int64(100), nil).
			Times(1)

		id, err := svc.Create(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), id)
	})

	t.Run("role is RoleUser", func(t *testing.T) {
		params := baseParams
		params.Role = user.RoleUser
		expectedUser := &user.User{
			Name:        params.Name,
			WhatsAppJID: params.Phone + "@s.whatsapp.net",
			Role:        user.RoleUser,
			Email:       params.Email,
		}
		mockRepo.EXPECT().
			Insert(ctx, expectedUser).
			Return(int64(101), nil).
			Times(1)

		id, err := svc.Create(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, int64(101), id)
	})

	t.Run("role is RoleAdmin", func(t *testing.T) {
		params := baseParams
		params.Role = user.RoleAdmin
		expectedUser := &user.User{
			Name:        params.Name,
			WhatsAppJID: params.Phone + "@s.whatsapp.net",
			Role:        user.RoleAdmin,
			Email:       params.Email,
		}
		mockRepo.EXPECT().
			Insert(ctx, expectedUser).
			Return(int64(102), nil).
			Times(1)

		id, err := svc.Create(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, int64(102), id)
	})

	t.Run("invalid role, should default to RoleUser", func(t *testing.T) {
		params := baseParams
		params.Role = "invalid"
		expectedUser := &user.User{
			Name:        params.Name,
			WhatsAppJID: params.Phone + "@s.whatsapp.net",
			Role:        user.RoleUser,
			Email:       params.Email,
		}
		mockRepo.EXPECT().
			Insert(ctx, expectedUser).
			Return(int64(103), nil).
			Times(1)

		id, err := svc.Create(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, int64(103), id)
	})

	t.Run("repo error", func(t *testing.T) {
		params := baseParams
		params.Role = user.RoleUser
		expectedUser := &user.User{
			Name:        params.Name,
			WhatsAppJID: params.Phone + "@s.whatsapp.net",
			Role:        user.RoleUser,
			Email:       params.Email,
		}
		mockRepo.EXPECT().
			Insert(ctx, expectedUser).
			Return(int64(0), errors.New("insert error")).
			Times(1)

		id, err := svc.Create(ctx, params)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})
}
