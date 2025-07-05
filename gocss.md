Perfeito. Vou aprofundar a arquitetura e funcionamento interno do UnoCSS e elaborar um plano técnico atualizado e expandido para a implementação em Go, priorizando as funcionalidades essenciais, compatibilidade com o Templ e facilidade de extensão com presets. Retornarei em breve com o resultado.


# Introdução

UnoCSS é um mecanismo *atomic CSS* baseado em extração estática que escaneia o código-fonte em busca de classes utilitárias e gera apenas o CSS necessário. Internamente, o pipeline do UnoCSS faz a extração de tokens (classes/utilitários), aplica variantes (como `hover:`), expande atalhos, resolve regras (estáticas e dinâmicas), agrupa por camadas (layers) e emite CSS ordenado. A implementação em Go seguirá esse fluxo, sem dependência de JS, com foco em modularidade e desempenho. O plano a seguir detalha cada etapa (pipeline), os principais módulos/componentes e como cobrir conceitos como regras, variantes, atalhos, layers, preflights, extractors, cache e presets, visando extensibilidade e uso em produção.

## Visão Geral do Pipeline UnoCSS

O processamento do UnoCSS começa com **extração de classes** do código-fonte (HTML, templates, etc.) e termina com a **síntese do CSS final**. De forma resumida, o pipeline é:

* **Extração de tokens**: ler arquivos fonte (p. ex. `.templ`) e identificar strings de classes/utilitários.
* **Análise de tokens**: dividir cada sequência (token) em operador primário e variantes, aplicando transformações de variante (como `hover:`) progressivamente.
* **Resolução de regras**: para cada token processado, buscar correspondência em regras estáticas ou dinâmicas definidas no config. Regras estáticas mapeiam nomes fixos a declarações CSS; regras dinâmicas usam regex e função geradora.
* **Expansão de shortcuts**: detectar tokens correspondentes a atalhos (shortcuts) e substituir por múltiplos tokens básicos que serão resolvidos também como regras.
* **Aplicação de variantes**: variantes (como `hover:`, `sm:`, etc.) são aplicadas no seletor CSS gerado, tipicamente como pseudo-classes ou wrappers de media query.
* **Agrupamento por camadas (layers)**: cada utilitário CSS resultante é marcado com a camada configurada; por fim, as regras são agrupadas e ordenadas por camada conforme a precedência definida.
* **Preflight e CSS global**: injetar CSS global (resets ou preflight) definido no config antes das regras utilitárias.
* **Pós-processamento**: opcionais transformações finais, como minificação ou transformadores customizados (variant groups, diretrizes `@apply`, etc.).
* **Saída**: emitir o CSS concatenado em arquivo, respeitando ordem de camadas e combinações de seletores (a título de exemplo, regras idênticas podem ser mescladas).

Cada etapa será implementada em módulos separados. O diagrama abaixo ilustra o fluxo geral:



```
Templ/HTML (.templ) --(Extractor)--> Tokens raw 
    --> Expansão de atalhos 
        --> Aplicação de variantes 
            --> Parser de regras (static/dinâmico) 
                --> Regras correspondentes + seletor CSS (com variantes) 
                    --> Agrupamento por layer 
                        --> CSS final (préflight + utilitários)
```

A seguir detalhamos os componentes e o design dos módulos correspondentes.

## Módulos Principais

### 1. Extractors (Extração de Tokens)

Os *extractors* escaneiam arquivos-fonte (`.templ`, `.html`, etc.) e extraem potenciais classes/utilitários. Por padrão, um extractor *split* simples divide o conteúdo em palavras-chave, mas podemos escrever um extractor específico para Templ (por exemplo, usando regex para `class="..."`). O config permite múltiplos extractors: cada arquivo é processado por todos, juntando os tokens encontrados.

Em Go, teríamos algo como:

```go
type Extractor interface {
    Extract(code string, path string) []string
}
```

E um extractor default (`ExtractorSplit`) que faz `strings.Fields` no código. Um extractor Templ custom pode usar regex:

```go
var classRE = regexp.MustCompile(`class\s*=\s*"([^"]+)"`)
func (e *TemplExtractor) Extract(code, path string) []string {
    matches := classRE.FindAllStringSubmatch(code, -1)
    var tokens []string
    for _, m := range matches {
        tokens = append(tokens, strings.Fields(m[1])...)
    }
    return tokens
}
```

Todo token extraído é limpo (por exemplo, removemos filtros `!important`, ou valores arbitrários já processados por transformadores), e enviado à etapa seguinte.

### 2. Parser de Tokens e Regras

Após a extração, cada token bruto é analisado para **aplicar variantes** e **expandir shortcuts**, antes de resolver regras.

* **Variants**: implementamos variantes como funções que verificam prefixos. Exemplo em Go:

  ```go
  type Variant struct {
      Prefix string
      Selector func(base string) string
  }
  ```

  No processamento, iteramos sobre variantes configuradas: para um token `"hover:bg-red"`, encontramos `Prefix="hover:"`, então removemos o prefixo (`base="bg-red"`) e guardamos a transformação de seletor (`":hover"`). Variantes podem ser encadeadas (e.g. `"sm:hover:bg-red"` aplica `sm:` e depois `hover:`). O core chama successive vezes algo como `matchVariants(rawToken)` e acumula transformações.

* **Shortcuts**: atalhos são mapeamentos de um token para múltiplos utilitários. Exemplo estático: `"btn" → "py-2 px-4 ..."`; dinâmico: regex e função que retorna uma string ou slice. Em Go:

  ```go
  type Shortcut struct {
      Pattern *regexp.Regexp
      Expand func([]string) []string  // recebe grupos regex, retorna tokens expandidos
  }
  ```

  Ao processar `"btn-red"`, um shortcut dinâmico pode corresponder `/^btn-(.*)$/` e expandir em `[]string{"bg-"+color+"-400", "text-"+color+"-100", ...}`. Depois, cada token expandido é tratado normalmente pelo parser de regras. É importante prevenir loops infinitos: aplicamos shortcuts recursivamente até não encontrar mais correspondências, limitando profundidade.

* **Regras Estáticas vs Dinâmicas**: o core do resolvedor mantém duas estruturas: um mapa de regras estáticas (string → CSS) e uma lista de regras dinâmicas (regex + função). Em Go:

  ```go
  type Rule struct {
      Pattern *regexp.Regexp     // nil para static
      StaticClass string         // usado se for estática (Pattern=nil)
      Handler func([]string) CSS // função para regras dinâmicas
      Layer string              // camada alvo (opcional)
  }
  ```

  Regras estáticas são registradas no mapa `map[string]CSS` para lookup O(1) em tokens exatos. Regras dinâmicas são verificadas sequencialmente: para cada token, testar regex.Match; se casar, chamar o handler. Por exemplo:

  ```go
  // Static rule
  rules = append(rules, Rule{StaticClass: "m-1", Handler: func([]string) CSS{
      return CSS{"margin": "0.25rem"}
  }})
  // Dynamic rule
  rules = append(rules, Rule{Pattern: regexp.MustCompile(`^m-(\d+)$`), Handler: func(groups []string) CSS{
      d, _ := strconv.Atoi(groups[1])
      return CSS{"margin": fmt.Sprintf("%dpx", d*4)}
  }})
  ```

  Assim, o token `"m-3"` ativa a segunda regra. O exemplo acima geraria `.m-3 { margin: 12px; }`.

Cada correspondência produz uma estrutura interna que inclui o seletor CSS (lembre-se de escapar caracteres especiais) e o objeto CSS (mapa de propriedades→valores). Se variantes estiverem ativas, o seletor incorporará pseudo-classes ou wrappers. Todo token processado é armazenado (ou marcado) para evitar duplicatas e para ordenação final. A engine registra quais regras foram aplicadas, útil para cache e detalhes de debugging.

### 3. Agrupamento por Camadas (Layers)

Cada regra (ou shortcut/preflight) pode especificar uma *layer* opcional. Implementamos isso carregando a informação no metadado do `Rule`. Na hora de gerar CSS, coletamos as declarações em fatias por camada. Por exemplo, temos um mapa `map[string][]CSSRule` onde a chave é a layer. Ao resolver `"m-2"` numa regra com layer `"utilities"`, adicionamos essa declaração em `layers["utilities"]`.

A precedência de camadas é definida no config (`layers: map[string]int`). Regras sem layer ficam na `"default"`. Camadas são ordenadas conforme o valor (default = 0, negativo fica antes). Na concatenação final, iteramos camadas em ordem configurada. O código Go pode usar:

```go
layers := make(map[string]int) // configurado pelo usuário
// após agrupar: sort layers by value then name
// output CSS in that order
```

Isso garante que, por exemplo, utilitários (layer maior) sobrescrevam componentes (layer menor). Se `outputToCssLayers` for ativado, também podemos inserir a diretiva `@layer` correspondente no CSS gerado.

### 4. Preflights (CSS Global)

O UnoCSS permite injetar CSS arbitrário antes dos utilitários via **preflights**. No Go, definimos:

```go
type Preflight struct {
    Layer string            // opcional
    CSS   string            // conteúdo CSS bruto ou função de geração
}
```

Podemos ter fontes estáticas de CSS (style-reset) ou gerar dinamicamente (p. ex., via tema). Na geração, concatenamos todos `Preflight.CSS` na camada especificada (normalmente antes das regras). Exemplo de config:

```go
preflights := []Preflight{{
    Layer: "base",
    CSS: "* { margin:0; padding:0; }",
}}
```

Na saída final, eles aparecem no topo do CSS (geralmente layer `base`).

### 5. Cache e Desempenho

Para uso em produção, implementaremos caching local de tokens já processados. A engine mantém um cache (por token ou por arquivo) com o resultado de `parseToken`/`parseUtil`, evitando retrabalho em builds incrementais. Por exemplo, um `map[string]CSSRule` de tokens resolvidos. Esse cache deve invalidar entradas se a configuração mudar ou no watch mode entre reinicializações.

Além disso, durante geração podemos mesclar regras idênticas (múltiplos seletores para mesmas declarações) para reduzir tamanho do CSS. Exemplo: se `.m-2` e `.hover\:m-2:hover` produzem `{ margin:0.5rem }`, podemos uni-los numa regra combinada. No Go, podemos usar chaves de string (JSON) para agrupar.

### 6. CLI e Integração com Templ

O projeto Go fornecerá um *CLI* similar ao `@unocss/cli`. Este CLI lerá configurações, receberá padrões (glob) de arquivos `.templ` e observação (`--watch`) para regenerar CSS em mudança. Em Go:

```go
func main() {
    // Exemplo simplificado de flag
    patterns := flag.String("glob", "*.templ", "Glob de arquivos Templ")
    watch := flag.Bool("watch", false, "Ativa watch")
    outFile := flag.String("o", "uno.css", "Arquivo de saída")
    flag.Parse()

    // Carrega config (Go struct, YAML/JSON ou código Go)
    cfg := LoadConfig("unocss.yaml")
    generator := NewUnoGenerator(cfg)

    // Função de build
    build := func() {
        files := GlobFiles(*patterns)
        tokens := ExtractTokens(files, cfg.Extractors)
        css := generator.GenerateCSS(tokens)
        ioutil.WriteFile(*outFile, []byte(css), 0644)
    }
    build()
    if *watch {
        WatchFiles(*patterns, build)
    }
}
```

O modo *watch* fica escutando alterações nos arquivos `.templ`, disparando rebuilds automáticos. Cada build reutiliza o cache para acelerar reprocessamento.

## Configuração e Extensibilidade

### Configuração em Go

O arquivo de configuração do UnoCSS será representado em Go como uma struct. Por exemplo:

```go
type Config struct {
    Rules      []Rule
    Shortcuts  []Shortcut
    Variants   []Variant
    Extractors []Extractor
    Preflights []Preflight
    Layers     map[string]int
    Presets    []Preset // ver abaixo
    // outras opções: Safelist, Blocklist, Theme, etc.
}
```

Usuário pode criar essa configuração em código Go ou carregar de YAML/JSON. Ao inicializar o generator, aplicamos *presets* e mesclamos com config do usuário. A mesclagem prioriza configurações do usuário, enquanto presets fornecem defaults. Em Go, podemos fazer merges customizados (concatenação de slices, combinação profunda de maps).

### Presets

Presets são conjuntos pré-definidos de regras/variants/shortcuts. Implementamos presets como funções em Go, por exemplo:

```go
type Preset func(cfg *Config)

func PresetUnoBasic(cfg *Config) {
    cfg.Rules = append(cfg.Rules, Rule{ StaticClass: "m-2", Handler: func([]string) CSS{ return CSS{"margin":"0.5rem"} }})
    cfg.Variants = append(cfg.Variants, Variant{Prefix:"hover:", Selector: func(s string) string { return s + ":hover" }})
    // ...
}
```

O usuário carrega presets na config antes das regras custom. No startup, iteramos `cfg.Presets` invocando cada função para popular o config. A arquitetura modular permite criar novos presets facilmente (basta definir a função e adicioná-la à lista). Também é possível criar *plugins* de terceiros seguindo a mesma interface.

### API de Extensibilidade

Projetaremos a API de modo que seja simples registrar novas regras, variantes e extractors em tempo de código. Por exemplo:

```go
generator.AddRule(rule Rule)
generator.AddVariant(variant Variant)
```

ou via config:

```go
cfg.Rules = append(cfg.Rules, myRule)
```

A separação entre módulos garante que alguém pode implementar seu próprio extractor ou pós-processador (Postprocessor) e registrá-lo. Implementar um novo preset seria escrever uma função que preenche `Config`.

## Exemplo em Go: Definição de Regras e Variante

```go
// Estrutura de regra atômica
type Rule struct {
    // Se StaticClass != "", trata-se de regra estática exata
    StaticClass string
    // Regexp para regra dinâmica
    Pattern *regexp.Regexp
    // Handler para gerar CSS a partir de grupos regex
    Handler func(groups []string) map[string]string
    // Camada da regra (opcional)
    Layer string
}

// Exemplo: regra dinâmica para padding
pRule := Rule{
    Pattern: regexp.MustCompile(`^p-(\d+)$`),
    Handler: func(groups []string) map[string]string {
        n := groups[1]
        // Converte string para número
        val, _ := strconv.Atoi(n)
        return map[string]string{"padding": fmt.Sprintf("%dpx", val*4)}
    },
    Layer: "utilities",
}

// Variante hover
hoverVariant := Variant{
    Prefix: "hover:",
    Selector: func(sel string) string {
        return sel + ":hover"
    },
}

// Shortcut estático
cfg.Shortcuts = append(cfg.Shortcuts, Shortcut{
    Pattern: nil, // padrão para atalho estático pelo próprio map
    Expand: func(_ []string) []string {
        return []string{"py-2", "px-4", "bg-blue-500", "text-white"}
    },
})
```

## Conclusão

Este plano técnico detalha um design modular para replicar o funcionamento do UnoCSS em Go, cobrindo desde a extração de tokens até a geração do CSS final organizado por camadas. Cada conceito principal do UnoCSS (regras estáticas/dinâmicas, variantes, shortcuts, camadas, preflights, cache, presets) é endereçado com componentes dedicados. A arquitetura em Go favorecerá extensibilidade (presets e plugins customizados) e desempenho (caching local e geração incremental). Exemplos de código em Go acima ilustram como estruturar regras, variantes e o CLI observador, garantindo que o projeto possa evoluir para uso em produção com boa performance e capacidade de configuração.

**Referências:** Documentação oficial do UnoCSS (regras, variantes, extractors, etc) e análise do código-fonte, além de explicações do pipeline em posts de desenvolvedores.

