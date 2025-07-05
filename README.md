# GOCSS: Uma implementação do UnoCSS em Go

GOCSS é uma implementação nativa em Go do motor de CSS atômico [UnoCSS](https://unocss.dev/). O objetivo é fornecer uma ferramenta de alto desempenho, sem dependências de JavaScript, para gerar CSS utilitário sob demanda, com foco especial na integração com o ecossistema Go, incluindo o `templ`.

## Visão Geral

Este projeto reimplementa o pipeline de processamento do UnoCSS em Go, mantendo a arquitetura modular e extensível que o torna tão poderoso. O fluxo principal consiste em:

1.  **Extração:** Analisa arquivos-fonte (`.templ`, `.html`, etc.) para extrair classes de utilitários.
2.  **Análise e Resolução:** Processa cada classe (token) aplicando variantes (`hover:`, `md:`), expandindo atalhos e correspondendo a regras de geração de CSS.
3.  **Geração de CSS:** Monta a folha de estilos final, agrupando o CSS por camadas (`layers`), adicionando preflights (resets) e ordenando tudo para garantir a precedência correta.

## Plano de Implementação

O desenvolvimento do GOCSS seguirá um plano por fases para garantir uma base sólida e um progresso incremental.

### Fase 1: Fundações do `core`

-   [x] Definir todas as estruturas de dados e interfaces em `pkg/core/types.go`.
-   [x] Implementar o `UnoGenerator` e o esqueleto do método `Generate`.
-   [x] Implementar a lógica de resolução de configuração (`resolveConfig`) que mescla presets e configurações do usuário.

### Fase 2: Regras e `preset-wind`

-   [x] Implementar a lógica de correspondência de regras (estáticas e dinâmicas).
-   [x] Portar as regras e utilitários do `preset-wind` do UnoCSS (iniciado).
-   [x] Garantir que os handlers de regras possam gerar as entradas CSS (`CSSEntry`) corretamente.

### Fase 3: Variantes

-   [x] Implementar o sistema de correspondência de variantes (`matchVariants`).
-   [x] Implementar a lógica de aplicação de variantes (`applyVariants`) para envolver o CSS com pseudo-classes e media queries.
-   [x] Portar as variantes mais comuns (pseudo-classes, media queries).

### Fase 4: Atalhos e Camadas

-   [x] Implementar a resolução de atalhos (estáticos e dinâmicos), incluindo a expansão recursiva.
-   [x] Implementar o sistema de camadas, incluindo a ordenação e a geração com a diretiva `@layer`.

### Fase 5: Extratores e CLI

-   [x] Implementar um extrator padrão para strings e um mais robusto para HTML/Templ.
-   [x] Construir a ferramenta de linha de comando (`gocss`) para geração de arquivos, incluindo o modo de observação (`--watch`).

### Fase 6: Presets Adicionais e Testes

-   [ ] Portar presets mais complexos como `preset-wind` e `preset-icons`.
-   [x] Escrever um conjunto abrangente de testes de unidade e de integração (iniciado).

### Fase 7: Integração com Templ e Documentação

-   [ ] Refinar a CLI e fornecer exemplos claros de como usá-la em um projeto `templ`.
-   [ ] Escrever a documentação detalhada para o uso e extensibilidade do GOCSS.

## Como Usar

A ferramenta de linha de comando `gocss` permite gerar CSS a partir dos seus arquivos-fonte. Você pode especificar os arquivos de entrada usando padrões glob e o arquivo de saída.

```bash
# Gerar CSS a partir de arquivos HTML no diretório atual
gocss --input "./*.html" --output ./gocss.css

# Observar arquivos HTML e regenerar CSS automaticamente
gocss --input "./*.html" --output ./gocss.css --watch
```

### Configuração

A configuração do GOCSS é feita em um arquivo Go (por exemplo, `gocss.config.go`), que oferece total flexibilidade para definir regras, variantes, atalhos e presets.

```go
// gocss.config.go
package main

import (
	"github.com/su3h7am/gocss/pkg/core"
	"github.com/su3h7am/gocss/pkg/extractor"
	"github.com/su3h7am/gocss/pkg/preset"
)

func GetConfig() *core.Config {
	return &core.Config{
		Presets: []core.Preset{
			preset.NewWind(), // Equivalente ao preset-wind
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
		Safelist: []string{
			"bg-red-500",
			"text-white",
		},
	}
}
```

Para usar sua configuração personalizada, passe o caminho do arquivo para a CLI:

```bash
# Gerar CSS usando um arquivo de configuração personalizado
gocss --config ./gocss.config.go --input "./*.html" --output ./gocss.css
```

### Usando com Templ

Para integrar o GOCSS com projetos `templ`, você pode usar o `TemplExtractor` e especificar os arquivos `.templ` como entrada. O `TemplExtractor` é projetado para analisar a sintaxe `templ` e extrair as classes CSS.

Exemplo de uso com arquivos `templ`:

```bash
# Gerar CSS a partir de arquivos templ
gocss --input "./**/*.templ" --output ./static/gocss.css

# Observar arquivos templ e regenerar CSS automaticamente
gocss --input "./**/*.templ" --output ./static/gocss.css --watch
```

Certifique-se de que seu arquivo `gocss.config.go` inclua o `TemplExtractor` na sua lista de extractors, como mostrado no exemplo de configuração acima.