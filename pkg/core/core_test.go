package core

import (
	"reflect"
	"regexp"
	"testing"
)

func TestMatchRule(t *testing.T) {
	// Setup a dummy ResolvedConfig with some rules
	cfg := &ResolvedConfig{
		Rules: []Rule{
			{
				Static: "block",
				Handler: func(match []string, ctx *RuleContext) *CSSEntry {
					return &CSSEntry{Properties: map[string]string{"display": "block"}, Selector: ".block"}
				},
				Meta: &RuleMeta{Layer: "utilities"},
			},
			{
				Matcher: regexp.MustCompile(`^m-(\d+)$`),
				Handler: func(match []string, ctx *RuleContext) *CSSEntry {
					return &CSSEntry{Properties: map[string]string{"margin": match[1] + "px"}, Selector: ".m-" + match[1]}
				},
				Meta: &RuleMeta{Layer: "utilities"},
			},
		},
	}

	generator := &UnoGenerator{Config: cfg}

	// Test static rule match
	t.Run("static rule", func(t *testing.T) {
		rule, match := generator.matchRule("block")
		if rule == nil {
			t.Fatal("Expected rule for 'block', got nil")
		}
		if rule.Static != "block" {
			t.Errorf("Expected static rule 'block', got %s", rule.Static)
		}
		if !reflect.DeepEqual(match, []string{"block"}) {
			t.Errorf("Expected match [\"block\"] for 'block', got %v", match)
		}
	})

	// Test dynamic rule match
			t.Run("dynamic rule", func(t *testing.T) {
		rule, match := generator.matchRule("m-10")
		if rule == nil {
			t.Fatal("Expected rule for 'm-10', got nil")
		}
		if rule.Matcher == nil {
			t.Fatal("Expected dynamic rule, got static")
		}
		if !reflect.DeepEqual(match, []string{"m-10", "10"}) {
			t.Errorf("Expected match [\"m-10\", \"10\"] for 'm-10', got %v", match)
		}
	})

	// Test no match
	t.Run("no match", func(t *testing.T) {
		rule, match := generator.matchRule("non-existent")
		if rule != nil {
			t.Errorf("Expected no rule for 'non-existent', got %v", rule)
		}
		if match != nil {
			t.Errorf("Expected no match for 'non-existent', got %v", match)
		}
	})
}
