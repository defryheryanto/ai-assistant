package main

import (
	"context"
	"fmt"

	"github.com/defryheryanto/ai-assistant/config"
	"github.com/defryheryanto/ai-assistant/pkg/calendar/google"
	"github.com/defryheryanto/ai-assistant/pkg/tools"
	"github.com/defryheryanto/ai-assistant/pkg/tools/calendar"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	config.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	llm, err := openai.New(
		openai.WithToken(config.OpenAIToken),
		openai.WithModel(config.OpenAIModel),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to setup llm: %v", err))
	}

	toolRegistry := tools.NewRegistry(llm, tools.WithLoggerOption())

	calendarService, err := google.New(ctx, config.GoogleCredentialsFilePath, config.GoogleTokenFilePath)
	if err != nil {
		panic(fmt.Sprintf("failed to setup google calendar service: %v", err))
	}
	toolRegistry.Register(calendar.NewCreateEventTool(calendarService, true))

	inquiry := "I have a meeting with John Doe on October 2 2025, from 9 AM to 10 AM at Central Park. Please create a calendar event for this meeting and invite"
	resp, err := toolRegistry.Execute(ctx, "default", inquiry)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
