package main

import (
	"context"
	"fmt"

	"github.com/defryheryanto/ai-assistant/config"
	"github.com/defryheryanto/ai-assistant/pkg/tools"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	config.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	llm, err := openai.New(
		openai.WithToken("your-open-ai-key"),
		openai.WithModel("gpt-4.1-mini"),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to setup llm: %v", err))
	}

	toolRegistry := tools.NewRegistry(llm, tools.WithLoggerOption())
	toolRegistry.Register(NewForecastTool())

	inquiry := "What is the weather like in Boston?"
	resp, err := toolRegistry.Execute(ctx, "default", inquiry)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
