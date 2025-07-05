# Plano de Implementação Detalhado: UnoCSS em Golang para Templ

Este documento descreve a arquitetura do UnoCSS e um plano detalhado para sua reimplementação em Golang, com foco na integração com o `templ`.

## 1. Visão Geral da Arquitetura UnoCSS

O UnoCSS é um motor de CSS atômico altamente configurável e extensível. Seu poder reside em um pipeline de processamento multifásico que transforma classes de utilitários em CSS de forma eficiente.

### 1.1. Fluxo de Processamento Principal

O processo, desde o código-fonte até o CSS final, segue estas etapas:

1.  **Extração (`Extraction`):** O código-fonte (HTML, componentes, etc.) é analisado por "Extratores" para encontrar todos os possíveis nomes de classe e utilitários (tokens).
2.  **Análise do Token (`Token Parsing`):** Cada token extraído passa por um pipeline de resolução:
    a.  **Cache:** O sistema verifica se o token já foi processado.
    b.  **Pré-processamento (`Preprocess`):** Funções de pré-processamento podem modificar a string do token.
    c.  **Correspondência de Variantes (`Variant Matching`):** O token é iterativamente correspondido contra uma lista de `Variants` (ex: `hover:`, `md:`, `dark:`). Cada variante correspondente remove seu prefixo e armazena um "handler" para ser aplicado posteriormente.
    d.  **Resolução de Atalhos (`Shortcut Resolution`):** O sistema verifica se o token (agora sem prefixos de variantes) corresponde a um `Shortcut`. Se sim, o atalho é expandido em um ou mais novos tokens, que são processados recursivamente.
    e.  **Resolução de Regras (`Rule Resolution`):** Se não for um atalho, o sistema tenta corresponder o token a uma `Rule`. As regras podem ser estáticas (correspondência exata) ou dinâmicas (usando expressões regulares). O handler da regra correspondente gera as propriedades CSS.
3.  **Geração do CSS (`CSS Generation`):**
    a.  **Aplicação de Variantes (`Variant Application`):** Os "handlers" das variantes, coletados anteriormente, são aplicados em ordem inversa, envolvendo o CSS gerado com media queries, pseudo-classes, etc.
    b.  **Pós-processamento (`Postprocess`):** Funções de pós-processamento podem modificar o objeto CSS final.
4.  **Montagem da Folha de Estilos (`Stylesheet Assembly`):**
    a.  **Preflights:** CSS base (resets) é adicionado.
    b.  **Camadas (`Layers`):** O CSS gerado é agrupado por camadas (`defaults`, `shortcuts`, `utilities`).
    c.  **Ordenação:** As camadas são ordenadas para garantir a precedência correta.
    d.  **Saída Final:** O CSS é serializado em uma string, potencialmente usando a diretiva `@layer` do CSS.

---

## 2. Estrutura de Implementação em Go

A implementação em Go deve espelhar essa arquitetura modular.

### 2.1. Estruturas de Dados Fundamentais (`pkg/core/types.go`)

As interfaces e structs em Go são a base para replicar a flexibilidade do UnoCSS.

```go
package core

import "regexp"

// UnoGenerator é a estrutura principal que orquestra todo o processo.
type UnoGenerator struct {
	Config *ResolvedConfig
	Cache  map[string][]*StringifiedUtil // Cache para tokens processados
	// ... outros campos de estado como Blocked, ParentOrders, etc.
}

// ResolvedConfig armazena a configuração final mesclada de presets e do usuário.
type ResolvedConfig struct {
	Theme         map[string]interface{}
	Rules         []Rule
	Variants      []Variant
	Shortcuts     map[string]Shortcut // Usar um mapa para acesso rápido a atalhos estáticos
	Preflights    []Preflight
	Extractors    []Extractor
	Layers        map[string]int
	Postprocess   []Postprocessor
	// ... outros campos de configuração
}

// Rule define como transformar um token em CSS.
type Rule struct {
	Matcher *regexp.Regexp // Para regras dinâmicas
	Static  string         // Para regras estáticas
	Handler func(match []string, ctx *RuleContext) *CSSEntry
	Meta    *RuleMeta
}

// RuleMeta contém metadados para uma regra.
type RuleMeta struct {
	Layer    string
	Internal bool
	// ... outros metadados
}

// Variant define como manipular prefixos como `hover:` ou `md:`.
type Variant struct {
	Matcher  func(token string, ctx *VariantContext) *VariantMatch
	Handler  func(entry *CSSEntry, match *VariantMatch) *CSSEntry
	MultiPass bool // Se a variante pode ser aplicada múltiplas vezes
}

// CSSEntry representa uma unidade de CSS gerada.
type CSSEntry struct {
	Properties map[string]string
	Selector   string
	Parent     string // Para media queries, etc.
	Layer      string
}

// RuleContext fornece contexto para os handlers de regras.
type RuleContext struct {
	RawSelector     string
	CurrentSelector string
	Theme           map[string]interface{}
	VariantHandlers []*VariantHandler // Handlers acumulados
}

// ... outras structs como Shortcut, Preflight, Extractor, etc.
```

### 2.2. O Motor de Geração (`pkg/core/generator.go`)

O `UnoGenerator` terá o método `Generate` como seu ponto de entrada principal.

```go
package core

// Generate processa um conjunto de tokens e retorna o CSS final.
func (g *UnoGenerator) Generate(tokens map[string]bool) (string, error) {
	// Mapa para agrupar CSS por camada
	layerCSS := make(map[string][]*StringifiedUtil)

	// 1. Processar cada token
	for token := range tokens {
		stringifiedUtils, err := g.ParseToken(token)
		if err != nil {
			// Lidar com o erro, talvez registrar e continuar
			continue
		}
		for _, util := range stringifiedUtils {
			layer := util.Layer
			if layer == "" {
				layer = "default"
			}
			layerCSS[layer] = append(layerCSS[layer], util)
		}
	}

	// 2. Adicionar Preflights
	// ...

	// 3. Ordenar camadas
	sortedLayers := g.sortLayers(layerCSS)

	// 4. Montar a string CSS final com @layer
	var finalCSS strings.Builder
	for _, layer := range sortedLayers {
		finalCSS.WriteString(fmt.Sprintf("@layer %s {
", layer))
		// Ordenar e mesclar regras dentro da camada
		// ...
		finalCSS.WriteString("}
")
	}

	return finalCSS.String(), nil
}

// ParseToken é o coração do pipeline de resolução.
func (g *UnoGenerator) ParseToken(token string) ([]*StringifiedUtil, error) {
	// a. Verificar cache
	if cached, ok := g.Cache[token]; ok {
		return cached, nil
	}

	// b. Pré-processamento (se houver)

	// c. Corresponder Variantes
	remainingToken, variantHandlers := g.matchVariants(token)

	// d. Expandir Atalhos (recursivamente)
	isShortcut, expandedTokens, err := g.expandShortcut(remainingToken)
	if err != nil { return nil, err }
	if isShortcut {
		// Processar os tokens expandidos e aplicar os variantHandlers
		// ...
	}

	// e. Corresponder Regras
	rule, match := g.matchRule(remainingToken)
	if rule == nil {
		// Token não correspondeu a nada
		g.Cache[token] = nil
		return nil, nil
	}

	// f. Gerar CSS a partir da regra
	ctx := &RuleContext{ /* ... */ }
	cssEntry := rule.Handler(match, ctx)

	// g. Aplicar Variantes
	finalEntry := g.applyVariants(cssEntry, variantHandlers)

	// h. Serializar e armazenar no cache
	// ...
}
```

### 2.3. Integração com Templ (`pkg/templ/integration.go`)

A integração com `templ` pode ser feita de duas maneiras principais:

1.  **Geração em Tempo de Compilação (Recomendado):**
    *   Criar uma ferramenta de CLI que observa os arquivos `.templ`.
    *   A CLI usa os `Extractors` para encontrar todas as classes usadas.
    *   Ela chama `core.UnoGenerator.Generate` para criar um único arquivo `unocss.css`.
    *   Este arquivo CSS é então servido estaticamente.

2.  **Geração em Tempo de Execução (Para Desenvolvimento):**
    *   Um middleware HTTP pode interceptar as requisições.
    *   Ele analisa o HTML renderizado pelo `templ` em cada requisição.
    *   Gera o CSS necessário dinamicamente. (Menos performático, mas útil para HMR).

#### CLI para Geração em Tempo de Compilação

```bash
# Exemplo de uso da CLI
unocss-go --watch --config unocss.config.go --output ./static/unocss.css "./**/*.templ"
```

O arquivo `unocss.config.go` permitiria a configuração usando Go:

```go
// unocss.config.go
package main

import (
	"github.com/su3h7am/unocss-go/core"
	"github.com/su3h7am/unocss-go/preset"
)

func GetConfig() *core.Config {
	return &core.Config{
		Presets: []core.Preset{
			preset.NewWind(), // Equivalente ao preset-wind
		},
		Safelist: []string{
			"bg-red-500",
			"text-white",
		},
	}
}
```

---

## 3. Plano de Implementação por Fases

1.  **Fase 1: Fundações do `core`**
    *   [ ] Definir todas as estruturas de dados em `pkg/core/types.go`.
    *   [ ] Implementar o `UnoGenerator` e o esqueleto do método `Generate`.
    *   [ ] Implementar a lógica de resolução de configuração (`resolveConfig`) que mescla presets.

2.  **Fase 2: Regras e `preset-mini`**
    *   [ ] Implementar a lógica de correspondência de regras (estáticas e dinâmicas).
    *   [ ] Portar as regras e utilitários do `preset-mini`. Este é o conjunto de regras mais fundamental e um ótimo ponto de partida.
    *   [ ] Garantir que os handlers de regras possam gerar `CSSEntry` corretamente.

3.  **Fase 3: Variantes**
    *   [ ] Implementar o sistema de correspondência de variantes (`matchVariants`).
    *   [ ] Implementar a lógica de aplicação de variantes (`applyVariants`).
    *   [ ] Portar as variantes mais comuns (pseudo-classes, media queries).

4.  **Fase 4: Atalhos e Camadas**
    *   [ ] Implementar a resolução de atalhos estáticos e dinâmicos, incluindo a recursão.
    *   [ ] Implementar o sistema de camadas, incluindo a ordenação e a geração com `@layer`.

5.  **Fase 5: Extratores e CLI**
    *   [ ] Implementar um extrator padrão que funciona com strings (similar ao `extractorSplit`).
    *   [ ] Implementar um extrator mais robusto para HTML/Templ.
    *   [ ] Construir a CLI para geração de arquivos, incluindo o modo de observação (`--watch`).

6.  **Fase 6: Presets Adicionais e Testes**
    *   [ ] Portar presets mais complexos como `preset-wind` e `preset-icons`.
    *   [ ] Escrever um conjunto abrangente de testes de unidade e de integração, comparando a saída com o UnoCSS original para garantir a paridade.

7.  **Fase 7: Integração com Templ e Documentação**
    *   [ ] Refinar a CLI e fornecer exemplos claros de como usá-la em um projeto `templ`.
    *   [ ] Escrever a documentação para o projeto Go, explicando como usá-lo e como estendê-lo.