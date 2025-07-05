package main

import (
	"fmt"

	"github.com/su3h7am/gocss/pkg/core"
	"github.com/su3h7am/gocss/pkg/preset"
)

func main() {
	cfg := &core.Config{
		Presets: []core.Preset{
			preset.NewMini(),
		},
	}

	resolvedConfig := core.NewResolvedConfig(cfg)
	generator := core.NewGenerator(resolvedConfig)

	tokens := map[string]bool{
		"m-4":   true,
		"p-8":   true,
		"block": true,
	}

	css, err := generator.Generate(tokens)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	fmt.Println(css)
}
