package image_test

import (
	"context"
	"testing"

	imagetool "github.com/defryheryanto/ai-assistant/pkg/tools/image"
	"github.com/openai/openai-go"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
)

func TestGenerateImageTool_SystemPrompt(t *testing.T) {
	tool := imagetool.NewGenerateImageTool(openai.Client{})
	assert.Equal(t, "", tool.SystemPrompt())
}

func TestGenerateImageTool_Definition(t *testing.T) {
	tool := imagetool.NewGenerateImageTool(openai.Client{})
	def := tool.Definition()

	assert.Equal(t, "function", def.Type)
	assert.NotNil(t, def.Function)
	if def.Function != nil {
		assert.Equal(t, "GenerateImage", def.Function.Name)
	}
}

func TestGenerateImageTool_Execute(t *testing.T) {
	tool := imagetool.NewGenerateImageTool(openai.Client{})
	ctx := context.Background()

	_, err := tool.Execute(ctx, llms.ToolCall{
		ID:           "123",
		FunctionCall: &llms.FunctionCall{Name: "GenerateImage", Arguments: `{"prompt":"a cat"}`},
	})

	// Since actual OpenAI call requires network, expect an error
	assert.Error(t, err)
}
