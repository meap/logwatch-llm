package presenter

import (
	"bytes"
	"text/template"

	_ "embed"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

//go:embed layout.html.tpl
var htmlLayout string

func ConvertMarkdownToHTML(markdown string) (string, error) {
	parser := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)

	var buf bytes.Buffer
	if err := parser.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}

	return renderWithLayout(buf.String())
}

func renderWithLayout(content string) (string, error) {
	tmpl, err := template.New("layout").Parse(htmlLayout)
	if err != nil {
		return "", err
	}

	var processedContent bytes.Buffer
	if err := tmpl.Execute(&processedContent, map[string]string{"Content": content}); err != nil {
		return "", err
	}

	return processedContent.String(), nil
}
