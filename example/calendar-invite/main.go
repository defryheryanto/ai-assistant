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
	toolRegistry, err := app.SetupTools(ctx, app.SetupToolsParams{
		GoogleCredentialsFilePath: config.GoogleCredentialsFilePath,
		GoogleTokenFilePath:       config.GoogleTokenFilePath,
		OpenAIToken:               config.OpenAIToken,
		OpenAIModel:               config.OpenAIModel,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to setup service: %v", err))
	}

	inquiry := `
		Create a calendar invite tomorrow for Badminton at 8PM GMT+7 until 11PM at Kharisma Badminton Hall.
		The players are considered intermediate, so prepare yourself!
	`
	resp, err := toolRegistry.Execute(ctx, inquiry)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
