package core

import (
	"reflect"
	"regexp"
	"strings"
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

func TestApplyVariants(t *testing.T) {
	cfg := &ResolvedConfig{
		Variants: []Variant{
			{
				Matcher: func(token string, ctx *VariantContext) *VariantMatch {
					if strings.HasPrefix(token, "hover:") {
						return &VariantMatch{Matcher: "hover:"}
					}
					return nil
				},
				Handler: func(entry *CSSEntry, match *VariantMatch) *CSSEntry {
					entry.Selector = entry.Selector + ":hover"
					return entry
				},
			},
			{
				Matcher: func(token string, ctx *VariantContext) *VariantMatch {
					if strings.HasPrefix(token, "sm:") {
						return &VariantMatch{Matcher: "sm:"}
					}
					return nil
				},
				Handler: func(entry *CSSEntry, match *VariantMatch) *CSSEntry {
					entry.Parent = "@media (min-width: 640px)"
					return entry
				},
			},
		},
	}
	generator := &UnoGenerator{Config: cfg}

	// Test single variant
	t.Run("single variant", func(t *testing.T) {
		entry := &CSSEntry{Selector: ".text-red", Properties: map[string]string{"color": "red"}}
		hoverVariant := &VariantHandler{Variant: &cfg.Variants[0], Match: &VariantMatch{Matcher: "hover:"}}
		
		finalEntry := generator.applyVariants(entry, []*VariantHandler{hoverVariant})
		
		if finalEntry.Selector != ".text-red:hover" {
			t.Errorf("Expected selector .text-red:hover, got %s", finalEntry.Selector)
		}
		if finalEntry.Parent != "" {
			t.Errorf("Expected empty parent, got %s", finalEntry.Parent)
		}
	})

	// Test multiple variants
	t.Run("multiple variants", func(t *testing.T) {
		entry := &CSSEntry{Selector: ".text-blue", Properties: map[string]string{"color": "blue"}}
		smVariant := &VariantHandler{Variant: &cfg.Variants[1], Match: &VariantMatch{Matcher: "sm:"}}
		hoverVariant := &VariantHandler{Variant: &cfg.Variants[0], Match: &VariantMatch{Matcher: "hover:"}}

		finalEntry := generator.applyVariants(entry, []*VariantHandler{smVariant, hoverVariant})

		if finalEntry.Selector != ".text-blue:hover" {
			t.Errorf("Expected selector .text-blue:hover, got %s", finalEntry.Selector)
		}
		if finalEntry.Parent != "@media (min-width: 640px)" {
			t.Errorf("Expected parent @media (min-width: 640px), got %s", finalEntry.Parent)
		}
	})
}

func TestExpandShortcut(t *testing.T) {
	cfg := &ResolvedConfig{
		Rules: []Rule{
			{
				Static: "py-2",
				Handler: func(match []string, ctx *RuleContext) *CSSEntry { return &CSSEntry{Properties: map[string]string{"padding-top": "0.5rem"}} },
				Meta: &RuleMeta{Layer: "utilities"},
			},
			{
				Static: "px-4",
				Handler: func(match []string, ctx *RuleContext) *CSSEntry { return &CSSEntry{Properties: map[string]string{"padding-left": "1rem"}} },
				Meta: &RuleMeta{Layer: "utilities"},
			},
			{
				Static: "bg-blue-500",
				Handler: func(match []string, ctx *RuleContext) *CSSEntry { return &CSSEntry{Properties: map[string]string{"background-color": "blue"}} },
				Meta: &RuleMeta{Layer: "utilities"},
			},
			{
				Static: "text-white",
				Handler: func(match []string, ctx *RuleContext) *CSSEntry { return &CSSEntry{Properties: map[string]string{"color": "white"}} },
				Meta: &RuleMeta{Layer: "utilities"},
			},
		},
		Shortcuts: []Shortcut{
			{
				Static: "btn",
				Expand: func(match []string) []string { return []string{"py-2", "px-4", "bg-blue-500", "text-white"} },
			},
		},
}
	generator := &UnoGenerator{Config: cfg}

	// Test static shortcut expansion
	t.Run("static shortcut", func(t *testing.T) {
		isShortcut, expandedTokens, err := generator.expandShortcut("btn")
		if err != nil { t.Fatal(err) }
		if !isShortcut { t.Fatal("Expected btn to be a shortcut") }
		if !reflect.DeepEqual(expandedTokens, []string{"py-2", "px-4", "bg-blue-500", "text-white"}) {
			t.Errorf("Expected expanded tokens %v, got %v", []string{"py-2", "px-4", "bg-blue-500", "text-white"}, expandedTokens)
		}
	})
}

func TestSortLayers(t *testing.T) {
	cfg := &ResolvedConfig{
		Layers: map[string]int{
			"base":       0,
			"components": 1,
			"utilities":  2,
		},
	}
	generator := &UnoGenerator{Config: cfg}

	// Dummy layerCSS map
	layerCSS := map[string][]*StringifiedUtil{
		"utilities":  {},
		"base":       {},
		"components": {},
		"custom":     {},
	}

	// Expected order
	expectedOrder := []string{"base", "components", "utilities", "custom"}

	// Get sorted layers
		sorted := generator.sortLayers(layerCSS)

	// Compare
		if !reflect.DeepEqual(sorted, expectedOrder) {
			t.Errorf("Expected sorted layers %v, got %v", expectedOrder, sorted)
	}
}