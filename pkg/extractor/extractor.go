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

var classRE = regexp.MustCompile(`class\s*=\s*"([^"]+)"|className\s*=\s*"([^"]+)"|:class="([^"]+)"|:className="([^"]+)"`) // Add more patterns as needed

func (e *TemplExtractor) Extract(code string, path string) []string {
	var tokens []string
	matches := classRE.FindAllStringSubmatch(code, -1)
	for _, m := range matches {
		for i, group := range m {
			if i == 0 || group == "" { // Skip full match and empty groups
				continue
			}
			tokens = append(tokens, strings.Fields(group)...)
		}
	}
	return tokens
}
