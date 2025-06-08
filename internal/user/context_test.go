package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserContext(t *testing.T) {
	t.Run("Nil User", func(t *testing.T) {
		ctx := context.Background()
		ctxWithUser := SetUserToContext(ctx, nil)
		assert.NotNil(t, ctxWithUser)

		result := GetUserFromContext(ctxWithUser)
		assert.Nil(t, result)
	})

	t.Run("Valid User", func(t *testing.T) {
		ctx := context.Background()
		u := &User{
			ID:   123,
			Name: "Alice",
		}
		ctxWithUser := SetUserToContext(ctx, u)
		assert.NotNil(t, ctxWithUser)

		result := GetUserFromContext(ctxWithUser)
		assert.NotNil(t, result)
		assert.EqualValues(t, u, result)
	})

	t.Run("No User In Context", func(t *testing.T) {
		ctx := context.Background()
		result := GetUserFromContext(ctx)
		assert.Nil(t, result)
	})

	t.Run("Wrong Type In Context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), key, "not a user")
		result := GetUserFromContext(ctx)
		assert.Nil(t, result)
	})
}
