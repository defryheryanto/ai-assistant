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

	hist, err := cw.GetHistory(ctx, "abc")
	assert.NoError(t, err)
	assert.Equal(t, conv, hist)

	// returned history is a copy
	conv[0] = llms.TextParts(llms.ChatMessageTypeHuman, "changed")
	hist2, err := cw.GetHistory(ctx, "abc")
	assert.NoError(t, err)
	assert.NotEqual(t, conv, hist2)

	// limit enforcement
	conv = append(conv, llms.TextParts(llms.ChatMessageTypeHuman, "third"))
	err = cw.SaveHistory(ctx, "abc", conv)
	assert.NoError(t, err)
	hist3, err := cw.GetHistory(ctx, "abc")
	assert.NoError(t, err)
	expected := conv[len(conv)-2:]
	assert.Equal(t, expected, hist3)
}
