package core

import "regexp"

// UnoGenerator é a estrutura principal que orquestra todo o processo.
type UnoGenerator struct {
	Config *ResolvedConfig
	Cache  map[string][]*StringifiedUtil // Cache para tokens processados
}

// ResolvedConfig armazena a configuração final mesclada de presets e do usuário.
type ResolvedConfig struct {
	Theme         map[string]interface{}
	Rules         []Rule
	Variants      []Variant
	Shortcuts     []Shortcut
	Preflights    []Preflight
	Extractors    []Extractor
	Layers        map[string]int
	Postprocess   []Postprocessor
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

// Outras structs a serem definidas
// StringifiedUtil representa uma regra de CSS processada e pronta para ser escrita.
type StringifiedUtil struct {
	Selector string
	Entries  map[string]string
	Layer    string
	Parent   string // For media queries, e.g., "@media (min-width: 640px)"
}
type Shortcut struct {
	Pattern *regexp.Regexp
	Static  string
	Expand  func(match []string) []string
}
type Preflight struct{}
type Extractor interface {
	Extract(code string, path string) []string
}
type Postprocessor interface{}
type VariantContext struct{}

type VariantMatch struct {
	Matcher string
	// Add other fields as needed for variant context, e.g., value for `sm:`
}

type VariantHandler struct {
	Variant *Variant
	Match   *VariantMatch
}
