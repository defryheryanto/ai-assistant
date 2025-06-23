package tools_test

import (
	"context"
	"testing"

	"github.com/defryheryanto/ai-assistant/internal/whatsapp/tools"
	pkgmock "github.com/defryheryanto/ai-assistant/pkg/tools/mock"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"go.uber.org/mock/gomock"
)

func TestGenerateImageTool_SystemPrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	base := pkgmock.NewMockTool(ctrl)
	base.EXPECT().SystemPrompt().Return("base prompt").Times(1)

	dec := tools.NewGenerateImageTool(base, nil)
	assert.Equal(t, "base prompt", dec.SystemPrompt())
}

func TestGenerateImageTool_Definition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	base := pkgmock.NewMockTool(ctrl)
	base.EXPECT().Definition().Return(llms.Tool{Type: "function"}).Times(1)

	dec := tools.NewGenerateImageTool(base, nil)
	def := dec.Definition()
	assert.Equal(t, "function", def.Type)
}

func TestGenerateImageTool_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	base := pkgmock.NewMockTool(ctrl)
	call := llms.ToolCall{ID: "1", FunctionCall: &llms.FunctionCall{Name: "GenerateImage"}}
	resp := &llms.MessageContent{Role: llms.ChatMessageTypeTool, Parts: []llms.ContentPart{llms.ToolCallResponse{ToolCallID: "1", Name: "GenerateImage", Content: "http://img"}}}

	base.EXPECT().Execute(context.Background(), call).Return(resp, nil).Times(1)

	dec := tools.NewGenerateImageTool(base, nil)
	out, err := dec.Execute(context.Background(), call)
	assert.NoError(t, err)
	assert.Equal(t, resp, out)
}
