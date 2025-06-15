package app

import (
	"context"
	"database/sql"

	"github.com/defryheryanto/ai-assistant/config"
	whatsapptool "github.com/defryheryanto/ai-assistant/internal/whatsapp/tools"
	"github.com/defryheryanto/ai-assistant/pkg/tools"
	calendartool "github.com/defryheryanto/ai-assistant/pkg/tools/calendar"
	timetool "github.com/defryheryanto/ai-assistant/pkg/tools/time"
	"github.com/tmc/langchaingo/llms/openai"
)

type SetupToolsParams struct {
	DB *sql.DB

	// CalendarService
	GoogleCredentialsFilePath string
	GoogleTokenFilePath       string

	// Registry
	OpenAIToken string
	OpenAIModel string
}

func SetupTools(ctx context.Context, params SetupToolsParams) (tools.Registry, *Services, error) {
	srv, err := SetupServices(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	llm, err := openai.New(
		openai.WithToken(params.OpenAIToken),
		openai.WithModel(params.OpenAIModel),
	)
	if err != nil {
		return nil, nil, err
	}

	toolRegistry := tools.NewRegistry(
		llm,
		tools.WithLoggerOption(),
		tools.WithSystemPromptOption(config.AssistantSystemPrompt),
	)
	registerTools(toolRegistry, srv)

	return toolRegistry, srv, nil
}

func registerTools(registry tools.Registry, srv *Services) {
	registry.Register(calendartool.NewCreateEventTool(srv.CalendarService, false))
	registry.Register(timetool.NewCurrentTimeTool())
	registry.Register(whatsapptool.NewCreateUserTool(srv.UserService))
	registry.Register(whatsapptool.NewRegisterGroupTool(srv.WhatsAppGroupService))
}
