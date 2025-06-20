package contextwindow_test

import (
	"context"
	"testing"

	"github.com/defryheryanto/ai-assistant/pkg/tools/contextwindow"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
)

func TestInMemoryContextWindow_GetHistory(t *testing.T) {
	cw := contextwindow.NewInMemoryContextWindow()
	ctx := context.Background()

	hist, err := cw.GetHistory(ctx, "abc")
	assert.NoError(t, err)
	assert.Empty(t, hist)
}

func TestInMemoryContextWindow_SaveHistory(t *testing.T) {
	cw := contextwindow.NewInMemoryContextWindow(contextwindow.WithLimit(2))
	ctx := context.Background()

	conv := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, "hi"),
		llms.TextParts(llms.ChatMessageTypeAI, "hello"),
	}

	err := cw.SaveHistory(ctx, "abc", conv)
	assert.NoError(t, err)

	t.Run("history stored", func(t *testing.T) {
		hist, err := cw.GetHistory(ctx, "abc")
		assert.NoError(t, err)
		assert.Equal(t, conv, hist)
	})

	t.Run("limit enforcement", func(t *testing.T) {
		err := cw.SaveHistory(ctx, "abc", []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, "third"),
		})
		assert.NoError(t, err)

		hist, err := cw.GetHistory(ctx, "abc")
		assert.NoError(t, err)
		expected := []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeAI, "hello"),
			llms.TextParts(llms.ChatMessageTypeHuman, "third"),
		}
		assert.Equal(t, expected, hist)
	})
}
