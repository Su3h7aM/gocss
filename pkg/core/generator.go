package core

import (
	"fmt"
	"strings"
	"sort"
)

// Generate processa um conjunto de tokens e retorna o CSS final.
func (g *UnoGenerator) Generate(tokens map[string]bool) (string, error) {
	layerCSS := make(map[string][]*StringifiedUtil)

	for token := range tokens {
		stringifiedUtils, err := g.ParseToken(token)
		if err != nil {
			// TODO: Lidar com o erro, talvez registrar e continuar
			continue
		}
		for _, util := range stringifiedUtils {
			// TODO: Obter a camada (layer) do util
			layer := "default"
			layerCSS[layer] = append(layerCSS[layer], util)
		}
	}

	// TODO: Adicionar Preflights

	// TODO: Ordenar camadas
	sortedLayers := g.sortLayers(layerCSS)

	var finalCSS strings.Builder
	for _, layer := range sortedLayers {
		finalCSS.WriteString(fmt.Sprintf("@layer %s {\n", layer))
		
		// Sort and merge rules within the layer (for now, just iterate)
		for _, util := range layerCSS[layer] {
			finalCSS.WriteString(fmt.Sprintf("  %s {\n", util.Selector))
			for prop, val := range util.Entries {
				finalCSS.WriteString(fmt.Sprintf("    %s: %s;\n", prop, val))
			}
			finalCSS.WriteString("  }\n")
		}
		finalCSS.WriteString("}\n")
	}

	return finalCSS.String(), nil
}

func (g *UnoGenerator) matchRule(token string) (*Rule, []string) {
	for _, rule := range g.Config.Rules {
		if rule.Static != "" {
			if rule.Static == token {
				return &rule, []string{token}
			}
		} else if rule.Matcher != nil {
			matches := rule.Matcher.FindStringSubmatch(token)
			if len(matches) > 0 {
				return &rule, matches
			}
		}
	}
	return nil, nil
}

// ParseToken é o coração do pipeline de resolução.
func (g *UnoGenerator) ParseToken(token string) ([]*StringifiedUtil, error) {
	// a. Verificar cache
	if cached, ok := g.Cache[token]; ok {
		return cached, nil
	}

	// e. Corresponder Regras
	rule, match := g.matchRule(token)
	if rule == nil {
		// Token não correspondeu a nada
		g.Cache[token] = nil
		return nil, nil
	}

	// f. Gerar CSS a partir da regra
	ctx := &RuleContext{ RawSelector: token, CurrentSelector: token } // Contexto simplificado por enquanto
	cssEntry := rule.Handler(match, ctx)

	if cssEntry == nil {
		return nil, nil
	}

	// g. Aplicar Variantes (a ser implementado)

	// h. Serializar e armazenar no cache
	util := &StringifiedUtil{
		Selector: cssEntry.Selector,
		Entries:  cssEntry.Properties,
		Layer:    rule.Meta.Layer,
	}
	g.Cache[token] = []*StringifiedUtil{util}

	return g.Cache[token], nil
}

func (g *UnoGenerator) sortLayers(layers map[string][]*StringifiedUtil) []string {
	keys := make([]string, 0, len(layers))
	for k := range layers {
		keys = append(keys, k)
	}
	
	// Sort alphabetically for now
		sort.Strings(keys)
	return keys
}