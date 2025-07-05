package core

import (
	"fmt"
	"strings"
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
		// TODO: Ordenar e mesclar regras dentro da camada
		finalCSS.WriteString("}\n")
	}

	return finalCSS.String(), nil
}

// ParseToken é o coração do pipeline de resolução.
func (g *UnoGenerator) ParseToken(token string) ([]*StringifiedUtil, error) {
	// a. Verificar cache
	if cached, ok := g.Cache[token]; ok {
		return cached, nil
	}

	// TODO: Implementar o resto do pipeline de resolução

	return nil, nil
}

func (g *UnoGenerator) sortLayers(layers map[string][]*StringifiedUtil) []string {
	// TODO: Implementar a lógica de ordenação de camadas com base na configuração
	keys := make([]string, 0, len(layers))
	for k := range layers {
		keys = append(keys, k)
	}
	return keys
}

