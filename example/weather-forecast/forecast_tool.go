package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

type ForecastParams struct {
	Location string `json:"location"`
}

type ForecastTool struct {
}

func NewForecastTool() *ForecastTool {
	return &ForecastTool{}
}

func (t *ForecastTool) Definition() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "GetWeatherForecast",
			Description: "Get weather forecast based on the location and unit",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"location": map[string]any{
						"type":        "string",
						"description": "The city and state, e.g. San Francisco, CA",
					},
					"unit": map[string]any{
						"type": "string",
						"enum": []string{"fahrenheit", "celsius"},
					},
				},
				"required": []string{"location"},
			},
		},
	}
}

func (t *ForecastTool) Execute(toolCall llms.ToolCall) (*llms.MessageContent, error) {
	var args ForecastParams
	if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
		return nil, err
	}

	weather, err := t.getCurrentWeather(args)
	if err != nil {
		return nil, err
	}

	return &llms.MessageContent{
		Role: llms.ChatMessageTypeTool,
		Parts: []llms.ContentPart{
			llms.ToolCallResponse{
				ToolCallID: toolCall.ID,
				Name:       toolCall.FunctionCall.Name,
				Content:    weather,
			},
		},
	}, nil
}

func (t *ForecastTool) getCurrentWeather(params ForecastParams) (string, error) {
	weatherResponses := map[string]string{
		"boston":  "72 and sunny",
		"chicago": "65 and windy",
	}

	loweredLocation := strings.ToLower(params.Location)

	var weatherInfo string
	found := false
	for key, value := range weatherResponses {
		if strings.Contains(loweredLocation, key) {
			weatherInfo = value
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("no weather info for %q", params.Location)
	}

	b, err := json.Marshal(weatherInfo)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
