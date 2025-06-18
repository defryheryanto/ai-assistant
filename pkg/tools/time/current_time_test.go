package time_test

import (
	"context"
	"testing"
	"time"

	tooltime "github.com/defryheryanto/ai-assistant/pkg/tools/time"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
)

func TestCurrentTimeTool_SystemPrompt(t *testing.T) {
	tool := tooltime.NewCurrentTimeTool()
	assert.Equal(t, "", tool.SystemPrompt())
}

func TestCurrentTimeTool_Definition(t *testing.T) {
	tool := tooltime.NewCurrentTimeTool()
	def := tool.Definition()

	assert.Equal(t, "function", def.Type)
	assert.NotNil(t, def.Function)
	if def.Function != nil {
		assert.Equal(t, "GetCurrentTime", def.Function.Name)
	}
}

func TestCurrentTimeTool_Execute(t *testing.T) {
	tool := tooltime.NewCurrentTimeTool()
	ctx := context.Background()

	resp, err := tool.Execute(ctx, llms.ToolCall{
		ID:           "123",
		FunctionCall: &llms.FunctionCall{Name: "GetCurrentTime"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, llms.ChatMessageTypeTool, resp.Role)

	res, ok := resp.Parts[0].(llms.ToolCallResponse)
	assert.True(t, ok)
	assert.Equal(t, "123", res.ToolCallID)
	assert.Equal(t, "GetCurrentTime", res.Name)

	_, parseErr := time.Parse(time.RFC3339Nano, res.Content)
	assert.NoError(t, parseErr)
}
