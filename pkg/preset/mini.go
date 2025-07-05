package preset

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/su3h7am/gocss/pkg/core"
)

// NewMini retorna um preset com regras básicas, similar ao preset-mini.
func NewMini() core.Preset {
	return func(config *core.ResolvedConfig) {
		config.Rules = append(config.Rules, getMiniRules()...)
	}
}

func getMiniRules() []core.Rule {
	return []core.Rule{
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
	}
}