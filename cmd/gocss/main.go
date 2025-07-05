package main

import (
	"fmt"

	"github.com/su3h7am/gocss/pkg/core"
	"github.com/su3h7am/gocss/pkg/extractor"
	"github.com/su3h7am/gocss/pkg/preset"
)

func main() {
	cfg := &core.Config{
		Presets: []core.Preset{
			preset.NewMini(),
		},
		Layers: map[string]int{
			"base":       0,
			"components": 1,
			"utilities":  2,
		},
		Extractors: []core.Extractor{
			&extractor.ExtractorSplit{},
			&extractor.TemplExtractor{},
		},
	}

	resolvedConfig := core.NewResolvedConfig(cfg)
	generator := core.NewGenerator(resolvedConfig)

	files := map[string]string{
		"test.html": `<div class="m-4 p-8 block text-red-500 bg-blue-500 text-lg font-bold w-full h-screen hover:text-green-500 sm:p-16 btn btn-red"></div>`,
		"test2.html": `<span class="text-white rounded"></span>`,
	}

	css, err := generator.Generate(files)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	fmt.Println(css)
}
