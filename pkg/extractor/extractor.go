package extractor

import (
	"regexp"
	"strings"
)

// ExtractorSplit implements core.Extractor by splitting the code by whitespace.
type ExtractorSplit struct{}

func (e *ExtractorSplit) Extract(code string, path string) []string {
	return strings.Fields(code)
}

// TemplExtractor implements core.Extractor by extracting class attributes from HTML/Templ code.
type TemplExtractor struct{}

// classRE matches class-like attributes and captures their content.
// It handles: class="...", className="...", :class="...", :className="..."
// It captures content within single or double quotes, including the quotes.
var classRE = regexp.MustCompile(`(?:class|className|:class|:className)\s*=\s*(["'][^"']*["'])`)

func (e *TemplExtractor) Extract(code string, path string) []string {
	var tokens []string
	matches := classRE.FindAllStringSubmatch(code, -1)
	for _, m := range matches {
		if len(m) > 1 && m[1] != "" {
			content := m[1]
			// Remove surrounding quotes
			if strings.HasPrefix(content, "\"") && strings.HasSuffix(content, "\"") {
				content = strings.TrimPrefix(content, "\"")
				content = strings.TrimSuffix(content, "\"")
			} else if strings.HasPrefix(content, "'") && strings.HasSuffix(content, "'") {
				content = strings.TrimPrefix(content, "'")
				content = strings.TrimSuffix(content, "'")
			}
			tokens = append(tokens, strings.Fields(content)...)
		}
	}
	return tokens
}
