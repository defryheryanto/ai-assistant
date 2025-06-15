package contextgroup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhatsAppContext(t *testing.T) {
	t.Run("Nil WhatsApp Context", func(t *testing.T) {
		ctx := context.Background()
		ctxWithUser := SetWhatsAppContext(ctx, nil)
		assert.NotNil(t, ctxWithUser)

		result := GetWhatsAppContext(ctxWithUser)
		assert.Nil(t, result)
	})

	t.Run("Valid WhatsApp Context", func(t *testing.T) {
		ctx := context.Background()
		u := &WhatsAppContext{
			CurrentChatJID: "jid1",
			SenderJID:      "jid1",
		}
		ctxWithUser := SetWhatsAppContext(ctx, u)
		assert.NotNil(t, ctxWithUser)

		result := GetWhatsAppContext(ctxWithUser)
		assert.NotNil(t, result)
		assert.EqualValues(t, u, result)
	})

	t.Run("No WhatsApp Context In Context", func(t *testing.T) {
		ctx := context.Background()
		result := GetWhatsAppContext(ctx)
		assert.Nil(t, result)
	})

	t.Run("Wrong Type In Context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userKey, "not a user")
		result := GetWhatsAppContext(ctx)
		assert.Nil(t, result)
	})
}
