package whatsapp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGroupService_GetByJID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockGroupRepository(ctrl)
	svc := whatsapp.NewGroupService(mockRepo)

	ctx := context.Background()
	jid := "12345@g.us"
	expectedGroup := &whatsapp.Group{
		ID:           1,
		GroupJID:     jid,
		IsActive:     true,
		RegisteredBy: "user123",
	}

	t.Run("found", func(t *testing.T) {
		mockRepo.EXPECT().
			FindByJID(ctx, jid).
			Return(expectedGroup, nil).
			Times(1)

		result, err := svc.GetByJID(ctx, jid)
		assert.NoError(t, err)
		assert.Equal(t, expectedGroup, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().
			FindByJID(ctx, jid).
			Return(nil, nil).
			Times(1)

		result, err := svc.GetByJID(ctx, jid)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("repo error", func(t *testing.T) {
		repoErr := errors.New("db error")
		mockRepo.EXPECT().
			FindByJID(ctx, jid).
			Return(nil, repoErr).
			Times(1)

		result, err := svc.GetByJID(ctx, jid)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, repoErr, err)
	})
}

func TestGroupService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockGroupRepository(ctrl)
	svc := whatsapp.NewGroupService(mockRepo)
	ctx := context.Background()

	params := whatsapp.CreateGroupParams{
		GroupJID:     "12345@g.us",
		RegisteredBy: "user123",
	}
	expectedGroup := &whatsapp.Group{
		GroupJID:     params.GroupJID,
		IsActive:     true,
		RegisteredBy: params.RegisteredBy,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			Insert(ctx, expectedGroup).
			Return(int64(1), nil).
			Times(1)

		id, err := svc.Create(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().
			Insert(ctx, expectedGroup).
			Return(int64(0), errors.New("insert error")).
			Times(1)

		id, err := svc.Create(ctx, params)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})
}
