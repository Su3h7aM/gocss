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
		
		// Group by parent (e.g., media queries)
		parentCSS := make(map[string][]*StringifiedUtil)
		for _, util := range layerCSS[layer] {
			parent := util.Parent
			if parent == "" {
				parent = "_default"
			}
			parentCSS[parent] = append(parentCSS[parent], util)
		}

		// Sort parents (e.g., media queries) for consistent output
		parentKeys := make([]string, 0, len(parentCSS))
		for k := range parentCSS {
			parentKeys = append(parentKeys, k)
		}
		sort.Strings(parentKeys)

		for _, parent := range parentKeys {
			if parent != "_default" {
				finalCSS.WriteString(fmt.Sprintf("  %s {\n", parent))
			}

			for _, util := range parentCSS[parent] {
				finalCSS.WriteString(fmt.Sprintf("    %s {\n", util.Selector))
				for prop, val := range util.Entries {
					finalCSS.WriteString(fmt.Sprintf("      %s: %s;\n", prop, val))
				}
				finalCSS.WriteString("    }\n")
			}

			if parent != "_default" {
				finalCSS.WriteString("  }\n")
			}
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

func (g *UnoGenerator) applyVariants(entry *CSSEntry, handlers []*VariantHandler) *CSSEntry {
	// Apply variants in reverse order
	for i := len(handlers) - 1; i >= 0; i-- {
		handler := handlers[i]
		entry = handler.Variant.Handler(entry, handler.Match)
	}
	return entry
}

func (g *UnoGenerator) matchVariants(token string) (string, []*VariantHandler) {
	var handlers []*VariantHandler
	current := token

	for {
		matched := false
		for _, variant := range g.Config.Variants {
			ctx := &VariantContext{ /* Add relevant context here if needed */ }
			if m := variant.Matcher(current, ctx); m != nil {
				// The matcher should return the remaining token and the match details
				// For now, let's assume m.Matcher contains the prefix that was matched
				// and the remaining token is current after trimming the prefix.
				
				if strings.HasPrefix(current, m.Matcher) {
					current = strings.TrimPrefix(current, m.Matcher)
					handlers = append(handlers, &VariantHandler{
						Variant: &variant,
						Match:   m,
					})
					matched = true
					break
				}
			}
		}
		if !matched {
			break
		}
	}
	return current, handlers
}

// ParseToken é o coração do pipeline de resolução.
func (g *UnoGenerator) ParseToken(token string) ([]*StringifiedUtil, error) {
	// a. Verificar cache
	if cached, ok := g.Cache[token]; ok {
		return cached, nil
	}

	// c. Corresponder Variantes
	remainingToken, variantHandlers := g.matchVariants(token)

	// e. Corresponder Regras
	rule, match := g.matchRule(remainingToken)
	if rule == nil {
		// Token não correspondeu a nada
		g.Cache[token] = nil
		return nil, nil
	}

	// f. Gerar CSS a partir da regra
	ctx := &RuleContext{ RawSelector: token, CurrentSelector: remainingToken } // Contexto simplificado por enquanto
	cssEntry := rule.Handler(match, ctx)

	if cssEntry == nil {
		return nil, nil
	}

	// g. Aplicar Variantes
	finalEntry := g.applyVariants(cssEntry, variantHandlers)

	// h. Serializar e armazenar no cache
	util := &StringifiedUtil{
		Selector: finalEntry.Selector,
		Entries:  finalEntry.Properties,
		Layer:    rule.Meta.Layer,
		Parent:   finalEntry.Parent,
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