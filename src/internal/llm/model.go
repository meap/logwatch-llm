package llm

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/openai"
)

type ModelInfo struct {
	Provider  string
	ModelName string
}

var SupportedModels = map[string]ModelInfo{
	"gpt-4o":                  {Provider: "openai", ModelName: "gpt-4o"},
	"gpt-4.1-mini":            {Provider: "openai", ModelName: "gpt-4.1-mini"},
	"claude-3-5-haiku-latest": {Provider: "anthropic", ModelName: "claude-3-5-haiku-latest"},
	// ...add more as needed
}

func NewLLMFromConfig(modelName string) (llms.Model, error) {
	info, ok := SupportedModels[modelName]
	if !ok {
		return nil, fmt.Errorf("unsupported model: %s", modelName)
	}

	logrus.Infof("Creating LLM with provider: %s, model: %s", info.Provider, info.ModelName)

	switch info.Provider {
	case "openai":
		return openai.New(openai.WithModel(info.ModelName))
	case "anthropic":
		return anthropic.New(anthropic.WithModel(info.ModelName))
	default:
		return nil, fmt.Errorf("unsupported provider: %s", info.Provider)
	}
}
