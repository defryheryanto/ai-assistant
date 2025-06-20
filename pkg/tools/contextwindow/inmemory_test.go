package contextwindow_test

import (
	"context"
	"testing"

	"github.com/defryheryanto/ai-assistant/pkg/tools/contextwindow"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
)

func TestInMemoryContextWindow(t *testing.T) {
	cw := contextwindow.NewInMemoryContextWindow()
	ctx := context.Background()

	t.Run("returns empty history for new id", func(t *testing.T) {
		hist, err := cw.GetHistory(ctx, "abc")
		assert.NoError(t, err)
		assert.Empty(t, hist)
	})

	conv := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, "hi"),
		llms.TextParts(llms.ChatMessageTypeAI, "hello"),
	}

	t.Run("saves and retrieves history", func(t *testing.T) {
		err := cw.SaveHistory(ctx, "abc", conv)
		assert.NoError(t, err)

		hist, err := cw.GetHistory(ctx, "abc")
		assert.NoError(t, err)
		assert.Equal(t, conv, hist)
	})

	t.Run("returned history is a copy", func(t *testing.T) {
		conv[0] = llms.TextParts(llms.ChatMessageTypeHuman, "changed")
		hist, err := cw.GetHistory(ctx, "abc")
		assert.NoError(t, err)
		assert.NotEqual(t, conv, hist)
		expected := []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, "hi"),
			llms.TextParts(llms.ChatMessageTypeAI, "hello"),
		}
		assert.Equal(t, expected, hist)
	})
}
