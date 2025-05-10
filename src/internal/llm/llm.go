package llm

import (
	"context"
	"strings"

	_ "embed"

	"github.com/meap/logwatch-llm/internal/logwatch"
	"github.com/meap/logwatch-llm/internal/system"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/llms"
)

//go:embed system_prompt.txt.tpl
var systemPrompt string

type AnalyzeResult struct {
	Content    string
	StopReason string
}

func populateSystemPrompt(tmpl string) string {
	info := system.GetSystemInfo()

	data := map[string]string{
		"OS":              info.OS,
		"Arch":            info.Arch,
		"KernelVersion":   info.KernelVersion,
		"Platform":        info.Platform,
		"PlatformVersion": info.PlatformVersion,
	}
	// Add any extra info from Other
	for k, v := range info.Other {
		data[k] = v
	}

	// Serialize as comma-separated key=value pairs
	var sb strings.Builder
	first := true
	for k, v := range data {
		if v == "" {
			continue
		}
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(k + "=" + v)
	}

	// Use the template, replacing {{.SystemInfo}} with the serialized string
	result := strings.ReplaceAll(tmpl, "{{.SystemInfo}}", sb.String())

	return result
}

func sanitizeMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) >= 2 && strings.HasPrefix(lines[0], "```") && strings.HasPrefix(lines[len(lines)-1], "```") {
		lines = lines[1 : len(lines)-1]
	}

	return strings.Join(lines, "\n")
}

func AnalyzeLogwatchSections(modelName string, sections []logwatch.LogwatchSection) (AnalyzeResult, error) {
	llm, err := NewLLMFromConfig(modelName)
	if err != nil {
		return AnalyzeResult{}, err
	}

	ctx := context.Background()

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, populateSystemPrompt(systemPrompt)),
	}

	for _, section := range sections {
		msg := "SECTION: " + section.Name + "\nCONTENT:\n" + section.Content
		content = append(content, llms.TextParts(llms.ChatMessageTypeHuman, msg))
	}

	response, err := llm.GenerateContent(ctx, content, llms.WithTemperature(0.2), llms.WithMaxTokens(8192))

	if err != nil {
		return AnalyzeResult{}, err
	}

	logrus.Info("LLM successfully generated response with stop reason: ", response.Choices[0].StopReason)

	return AnalyzeResult{
		Content:    sanitizeMarkdown(response.Choices[0].Content),
		StopReason: response.Choices[0].StopReason,
	}, nil
}
