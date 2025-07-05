package main

import (
	"github.com/su3h7am/gocss/pkg/core"
	"github.com/su3h7am/gocss/pkg/extractor"
	"github.com/su3h7am/gocss/pkg/preset"
)

func GetConfig() *core.Config {
	return &core.Config{
		Presets: []core.Preset{
			preset.NewWind(),
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
}
