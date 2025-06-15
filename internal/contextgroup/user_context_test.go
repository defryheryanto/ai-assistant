package contextgroup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserContext(t *testing.T) {
	t.Run("Nil User", func(t *testing.T) {
		ctx := context.Background()
		ctxWithUser := SetUserContext(ctx, nil)
		assert.NotNil(t, ctxWithUser)

		result := GetUserContext(ctxWithUser)
		assert.Nil(t, result)
	})

	t.Run("Valid User", func(t *testing.T) {
		ctx := context.Background()
		u := &UserContext{
			ID:   123,
			Name: "Alice",
		}
		ctxWithUser := SetUserContext(ctx, u)
		assert.NotNil(t, ctxWithUser)

		result := GetUserContext(ctxWithUser)
		assert.NotNil(t, result)
		assert.EqualValues(t, u, result)
	})

	t.Run("No User In Context", func(t *testing.T) {
		ctx := context.Background()
		result := GetUserContext(ctx)
		assert.Nil(t, result)
	})

	t.Run("Wrong Type In Context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userKey, "not a user")
		result := GetUserContext(ctx)
		assert.Nil(t, result)
	})
}
