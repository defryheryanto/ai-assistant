package app

import (
	"context"

	"github.com/defryheryanto/ai-assistant/pkg/tools"
	calendartool "github.com/defryheryanto/ai-assistant/pkg/tools/calendar"
	timetool "github.com/defryheryanto/ai-assistant/pkg/tools/time"
	"github.com/tmc/langchaingo/llms/openai"
)

type SetupToolsParams struct {
	// CalendarService
	GoogleCredentialsFilePath string
	GoogleTokenFilePath       string

	// Registry
	OpenAIToken string
	OpenAIModel string
}

func SetupTools(ctx context.Context, params SetupToolsParams) (tools.Registry, error) {
	srv, err := setupServices(ctx, params)
	if err != nil {
		return nil, err
	}

	llm, err := openai.New(
		openai.WithToken(params.OpenAIToken),
		openai.WithModel(params.OpenAIModel),
	)
	if err != nil {
		return nil, err
	}

	toolRegistry := tools.NewRegistry(llm, true)
	registerTools(toolRegistry, srv)

	return toolRegistry, nil
}

func registerTools(registry tools.Registry, srv *services) {
	registry.Register(calendartool.NewCreateEventTool(srv.CalendarService))
	registry.Register(timetool.NewCurrentTimeTool())
}
