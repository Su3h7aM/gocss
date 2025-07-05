package preset

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/su3h7am/gocss/pkg/core"
)

// NewMini retorna um preset com regras básicas, similar ao preset-mini.
func NewMini() core.Preset {
	return func(config *core.ResolvedConfig) {
		config.Rules = append(config.Rules, getMiniRules()...)
		config.Variants = append(config.Variants, getMiniVariants()...)
		config.Shortcuts = append(config.Shortcuts, getMiniShortcuts()...)
	}
}

func getMiniRules() []core.Rule {
	return []core.Rule{
		// Base rule for testing layers
		{
			Static: "html",
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				return &core.CSSEntry{
					Properties: map[string]string{"box-sizing": "border-box"},
					Selector:   "html",
				}
			},
			Meta: &core.RuleMeta{Layer: "base"},
		},
		// Margin
		{
			Matcher: regexp.MustCompile(`^m-(\d+)$`),
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				val, _ := strconv.Atoi(match[1])
				return &core.CSSEntry{
					Properties: map[string]string{"margin": fmt.Sprintf("%dpx", val*4)},
					Selector:   fmt.Sprintf(".%s", ctx.RawSelector),
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		// Padding
		{
			Matcher: regexp.MustCompile(`^p-(\d+)$`),
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				val, _ := strconv.Atoi(match[1])
				return &core.CSSEntry{
					Properties: map[string]string{"padding": fmt.Sprintf("%dpx", val*4)},
					Selector:   fmt.Sprintf(".%s", ctx.RawSelector),
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		{
			Static: "py-2",
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				return &core.CSSEntry{
					Properties: map[string]string{"padding-top": "0.5rem", "padding-bottom": "0.5rem"},
					Selector:   ".py-2",
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		{
			Static: "px-4",
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				return &core.CSSEntry{
					Properties: map[string]string{"padding-left": "1rem", "padding-right": "1rem"},
					Selector:   ".px-4",
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		// Display
		{
			Static: "block",
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				return &core.CSSEntry{
					Properties: map[string]string{"display": "block"},
					Selector:   ".block",
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		// Colors
		{
			Matcher: regexp.MustCompile(`^text-(red|blue|green)-(\d+)$`),
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				color := match[1]
				shade := match[2]
				// Simplified color mapping for now
				colors := map[string]map[string]string{
					"red":   {"500": "#ef4444"},
					"blue":  {"500": "#3b82f6"},
					"green": {"500": "#22c55e"},
				}
				if hex, ok := colors[color][shade]; ok {
					return &core.CSSEntry{
						Properties: map[string]string{"color": hex},
						Selector:   fmt.Sprintf(".%s", ctx.RawSelector),
					}
				}
				return nil
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		{
			Static: "text-white",
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				return &core.CSSEntry{
					Properties: map[string]string{"color": "#fff"},
					Selector:   ".text-white",
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		{
			Matcher: regexp.MustCompile(`^bg-(red|blue|green)-(\d+)$`),
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				color := match[1]
				shade := match[2]
				// Simplified color mapping for now
				colors := map[string]map[string]string{
					"red":   {"500": "#ef4444"},
					"blue":  {"500": "#3b82f6"},
					"green": {"500": "#22c55e"},
				}
				if hex, ok := colors[color][shade]; ok {
					return &core.CSSEntry{
						Properties: map[string]string{"background-color": hex},
						Selector:   fmt.Sprintf(".%s", ctx.RawSelector),
					}
				}
				return nil
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		// Typography
		{
			Static: "text-lg",
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				return &core.CSSEntry{
					Properties: map[string]string{"font-size": "1.125rem", "line-height": "1.75rem"},
					Selector:   ".text-lg",
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		{
			Static: "font-bold",
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				return &core.CSSEntry{
					Properties: map[string]string{"font-weight": "700"},
					Selector:   ".font-bold",
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		// Sizing
		{
			Static: "w-full",
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				return &core.CSSEntry{
					Properties: map[string]string{"width": "100%"},
					Selector:   ".w-full",
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		{
			Static: "h-screen",
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				return &core.CSSEntry{
					Properties: map[string]string{"height": "100vh"},
					Selector:   ".h-screen",
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
		// Border Radius
		{
			Static: "rounded",
			Handler: func(match []string, ctx *core.RuleContext) *core.CSSEntry {
				return &core.CSSEntry{
					Properties: map[string]string{"border-radius": "0.25rem"},
					Selector:   ".rounded",
				}
			},
			Meta: &core.RuleMeta{Layer: "utilities"},
		},
	}
}

func getMiniVariants() []core.Variant {
	return []core.Variant{
		{
			Matcher: func(token string, ctx *core.VariantContext) *core.VariantMatch {
				if strings.HasPrefix(token, "hover:") {
					return &core.VariantMatch{Matcher: "hover:"}
				}
				return nil
			},
			Handler: func(entry *core.CSSEntry, match *core.VariantMatch) *core.CSSEntry {
				entry.Selector = entry.Selector + ":hover"
				return entry
			},
		},
		{
			Matcher: func(token string, ctx *core.VariantContext) *core.VariantMatch {
				if strings.HasPrefix(token, "sm:") {
					return &core.VariantMatch{Matcher: "sm:"}
				}
				return nil
			},
			Handler: func(entry *core.CSSEntry, match *core.VariantMatch) *core.CSSEntry {
				entry.Parent = "@media (min-width: 640px)"
				return entry
			},
		},
	}
}

func getMiniShortcuts() []core.Shortcut {
	return []core.Shortcut{
		{
			Static: "btn",
			Expand: func(match []string) []string {
				return []string{"py-2", "px-4", "bg-blue-500", "text-white", "font-bold", "rounded"}
			},
		},
		{
			Pattern: regexp.MustCompile(`^btn-(red|blue|green)$`),
			Expand: func(match []string) []string {
				color := match[1]
				return []string{fmt.Sprintf("bg-%s-500", color), "text-white", "font-bold", "rounded"}
			},
		},
	}
}