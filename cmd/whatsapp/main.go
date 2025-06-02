package main

import (
	"context"
	"fmt"

	"github.com/defryheryanto/ai-assistant/config"
	"github.com/defryheryanto/ai-assistant/internal/app"
)

func main() {
	config.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := setupWhatsmeowClient(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to setup whatsmeow client: %v", err))
	}

	toolRegistry, err := app.SetupTools(ctx, app.SetupToolsParams{
		GoogleCredentialsFilePath: config.GoogleCredentialsFilePath,
		GoogleTokenFilePath:       config.GoogleTokenFilePath,
		OpenAIToken:               config.OpenAIToken,
		OpenAIModel:               config.OpenAIModel,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to setup service: %v", err))
	}
	client.AddEventHandler(eventHandler(ctx, client, toolRegistry))

	connectWhatsmeowClient(client)

}
