package contextgroup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseContext(t *testing.T) {
	t.Run("Nil Response", func(t *testing.T) {
		ctx := context.Background()
		ctxWith := SetResponseContext(ctx, nil)
		assert.NotNil(t, ctxWith)
		rc := GetResponseContext(ctxWith)
		assert.Nil(t, rc)
	})

	t.Run("Valid Response", func(t *testing.T) {
		ctx := context.Background()
		r := &ResponseContext{}
		ctx = SetResponseContext(ctx, r)
		rc := GetResponseContext(ctx)
		assert.NotNil(t, rc)
		assert.False(t, rc.MediaSent)
		MarkMediaSent(ctx)
		assert.True(t, r.MediaSent)
	})

	t.Run("Wrong Type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), responseKey, "not response")
		rc := GetResponseContext(ctx)
		assert.Nil(t, rc)
	})
}
