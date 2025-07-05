package extractor

import (
	"reflect"
	"testing"
)

func TestExtractorSplit(t *testing.T) {
	e := &ExtractorSplit{}

	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "basic split",
			code:     "m-4 p-8 block",
			expected: []string{"m-4", "p-8", "block"},
		},
		{
			name:     "with newlines and tabs",
			code:     "m-4\np-8\tblock",
			expected: []string{"m-4", "p-8", "block"},
		},
		{
			name:     "empty string",
			code:     "",
			expected: []string{},
		},
		{
			name:     "multiple spaces",
			code:     "m-4  p-8   block",
			expected: []string{"m-4", "p-8", "block"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := e.Extract(tt.code, "")
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ExtractorSplit.Extract() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTemplExtractor(t *testing.T) {
	e := &TemplExtractor{}

	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "basic class attribute",
			code:     `<div class="m-4 p-8"></div>`,
			expected: []string{"m-4", "p-8"},
		},
		{
			name:     "className attribute",
			code:     `<div className="flex items-center"></div>`,
			expected: []string{"flex", "items-center"},
		},
		{
			name:     "mixed attributes",
			code:     `<div class="m-4" className="p-8"></div>`,
			expected: []string{"m-4", "p-8"},
		},
		{
			name:     "no class attribute",
			code:     `<div></div>`,
			expected: []string{},
		},
		{
			name:     "empty class attribute",
			code:     `<div class=""></div>`,
			expected: []string{},
		},
		{
			name:     "class with newlines",
			code:     "<div class=\"m-4\np-8\"></div>", // Use actual newline
			expected: []string{"m-4", "p-8"},
		},
		{
			name:     "class with tabs",
			code:     "<div class=\"m-4\tp-8\"></div>", // Use actual tab
			expected: []string{"m-4", "p-8"},
		},
		{
			name:     "class with extra spaces",
			code:     `<div class=" m-4  p-8 "></div>`,
			expected: []string{"m-4", "p-8"},
		},
		{
			name:     "multiple class attributes in one tag",
			code:     `<div class="one two" className="three four"></div>`,
			expected: []string{"one", "two", "three", "four"},
		},
		{
			name:     "class in templ syntax",
			code:     `@div(templ.Class("foo bar"))`,
			expected: []string{},
		},
		{
			name:     "colon class attribute",
			code:     `<div :class="'foo bar'"></div>`,
			expected: []string{"foo", "bar"},
		},
		{
			name:     "colon className attribute",
			code:     `<div :className="'baz qux'"></div>`,
			expected: []string{"baz", "qux"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := e.Extract(tt.code, "")
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("TemplExtractor.Extract() got = %v, want %v", got, tt.expected)
			}
		})
	}
}