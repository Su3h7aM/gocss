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
	}
}
