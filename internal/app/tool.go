package app

import (
	"context"
	"database/sql"

	"github.com/defryheryanto/ai-assistant/config"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	whatsapptool "github.com/defryheryanto/ai-assistant/internal/whatsapp/tools"
	"github.com/defryheryanto/ai-assistant/pkg/tools"
	calendartool "github.com/defryheryanto/ai-assistant/pkg/tools/calendar"
	"github.com/defryheryanto/ai-assistant/pkg/tools/contextwindow"
	imagetool "github.com/defryheryanto/ai-assistant/pkg/tools/image"
	timetool "github.com/defryheryanto/ai-assistant/pkg/tools/time"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mau.fi/whatsmeow"
)

type SetupToolsParams struct {
	DB *sql.DB

	// CalendarService
	GoogleCredentialsFilePath string
	GoogleTokenFilePath       string

	// Registry
	OpenAIToken    string
	OpenAIModel    string
	WhatsAppClient *whatsmeow.Client
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

	contextWindowManager := contextwindow.NewInMemoryContextWindow(contextwindow.WithLimit(15))
	toolRegistry := tools.NewRegistry(
		llm,
		tools.WithLoggerOption(),
		tools.WithSystemPromptOption(config.AssistantSystemPrompt),
		tools.WithContextWindowManager(contextWindowManager),
	)
	registerTools(toolRegistry, srv, params.WhatsAppClient)

	return toolRegistry, srv, nil
}

func registerTools(registry tools.Registry, srv *Services, waClient *whatsmeow.Client) {
	registry.Register(calendartool.NewCreateEventTool(srv.CalendarService, false))
	registry.Register(timetool.NewCurrentTimeTool())

	allowedAdminRole := []whatsapp.UserRole{whatsapp.UserRoleAdmin}
	registry.Register(whatsapptool.NewRoleMiddlewareTool(whatsapptool.NewCreateUserTool(srv.UserService), allowedAdminRole))
	registry.Register(whatsapptool.NewRoleMiddlewareTool(whatsapptool.NewRegisterGroupTool(srv.WhatsAppGroupService), allowedAdminRole))
	imgTool := imagetool.NewGenerateImageTool(srv.OpenAIClient)
	registry.Register(whatsapptool.NewGenerateImageTool(imgTool, waClient))
}
