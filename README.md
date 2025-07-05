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

-   [ ] Definir todas as estruturas de dados e interfaces em `pkg/core/types.go`.
-   [ ] Implementar o `UnoGenerator` e o esqueleto do método `Generate`.
-   [ ] Implementar a lógica de resolução de configuração (`resolveConfig`) que mescla presets e configurações do usuário.

### Fase 2: Regras e `preset-mini`

-   [ ] Implementar a lógica de correspondência de regras (estáticas e dinâmicas).
-   [ ] Portar as regras e utilitários do `preset-mini` do UnoCSS.
-   [ ] Garantir que os handlers de regras possam gerar as entradas CSS (`CSSEntry`) corretamente.

### Fase 3: Variantes

-   [ ] Implementar o sistema de correspondência de variantes (`matchVariants`).
-   [ ] Implementar a lógica de aplicação de variantes (`applyVariants`) para envolver o CSS com pseudo-classes e media queries.
-   [ ] Portar as variantes mais comuns (pseudo-classes, media queries).

### Fase 4: Atalhos e Camadas

-   [ ] Implementar a resolução de atalhos (estáticos e dinâmicos), incluindo a expansão recursiva.
-   [ ] Implementar o sistema de camadas, incluindo a ordenação e a geração com a diretiva `@layer`.

### Fase 5: Extratores e CLI

-   [ ] Implementar um extrator padrão para strings e um mais robusto para HTML/Templ.
-   [ ] Construir a ferramenta de linha de comando (`gocss`) para geração de arquivos, incluindo o modo de observação (`--watch`).

### Fase 6: Presets Adicionais e Testes

-   [ ] Portar presets mais complexos como `preset-wind` e `preset-icons`.
-   [ ] Escrever um conjunto abrangente de testes de unidade e de integração.

### Fase 7: Integração com Templ e Documentação

-   [ ] Refinar a CLI e fornecer exemplos claros de como usá-la em um projeto `templ`.
-   [ ] Escrever a documentação detalhada para o uso e extensibilidade do GOCSS.

## Como Usar (Futuro)

A ferramenta de linha de comando permitirá a geração de CSS da seguinte forma:

```bash
# Observar arquivos .templ e gerar o CSS em static/gocss.css
gocss --watch --config gocss.config.go --output ./static/gocss.css "./**/*.templ"
```

A configuração será feita em um arquivo Go, permitindo total flexibilidade:

```go
// gocss.config.go
package main

import (
	"github.com/su3h7am/gocss/core"
	"github.com/su3h7am/gocss/preset"
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
