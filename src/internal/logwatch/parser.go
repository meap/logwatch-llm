package logwatch

import (
	"bufio"
	"regexp"
	"strings"
)

// Section represents a named section in the logwatch output.
type LogwatchSection struct {
	Name    string
	Content string
}

type LogwatchReport struct {
	Host string
}

type LogwatchResult struct {
	Report   LogwatchReport
	Sections []LogwatchSection
}

func parseHeaders(content string) LogwatchReport {
	scanner := bufio.NewScanner(strings.NewReader(content))
	hostRe := regexp.MustCompile(`^\s*Logfiles for Host: (.+?)\s*$`)
	var host string
	for scanner.Scan() {
		line := scanner.Text()
		if matches := hostRe.FindStringSubmatch(line); matches != nil {
			host = matches[1]
			break
		}
	}
	return LogwatchReport{Host: host}
}

// ParseSectionsFromString parses the input content and returns a slice of Section structs.
func ParseSectionsFromString(content string) LogwatchResult {
	var sections []LogwatchSection
	var currentSection string
	var contentBuilder strings.Builder
	inSection := false

	startRe := regexp.MustCompile(`-- (.*) Begin\s*-+\s*$`)
	endRe := regexp.MustCompile(`-- (.*) End\s*-+\s*$`)

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if matches := startRe.FindStringSubmatch(line); matches != nil {
			currentSection = matches[1]
			contentBuilder.Reset()
			inSection = true

			continue
		}
		if matches := endRe.FindStringSubmatch(line); matches != nil {
			if inSection && matches[1] == currentSection {
				sections = append(sections, LogwatchSection{Name: currentSection, Content: contentBuilder.String()})
				inSection = false
				currentSection = ""
			}
			continue
		}
		if inSection {
			contentBuilder.WriteString(line)
			contentBuilder.WriteByte('\n')
		}
	}

	headers := parseHeaders(content)

	return LogwatchResult{
		Sections: sections,
		Report:   headers,
	}
}
