package core

// Config é a estrutura de configuração que o usuário define.
type Config struct {
	Rules      []Rule
	Shortcuts  []Shortcut
	Variants   []Variant
	Extractors []Extractor
	Preflights []Preflight
	Layers     map[string]int
	Presets    []Preset
}

// Preset é uma função que aplica uma configuração pré-definida.
type Preset func(config *ResolvedConfig)

// NewResolvedConfig cria uma nova instância de ResolvedConfig aplicando presets e configurações do usuário.
func NewResolvedConfig(cfg *Config) *ResolvedConfig {
	resolved := &ResolvedConfig{
		Theme:         make(map[string]interface{}),
		Rules:         []Rule{},
		Variants:      []Variant{},
		Shortcuts:     []Shortcut{},
		Preflights:    []Preflight{},
		Extractors:    []Extractor{},
		Layers:        make(map[string]int),
		Postprocess:   []Postprocessor{},
	}

	// Apply presets first
	for _, p := range cfg.Presets {
		p(resolved)
	}

	// Merge user's config (user config overrides presets)
	resolved.Rules = append(resolved.Rules, cfg.Rules...)
	resolved.Variants = append(resolved.Variants, cfg.Variants...)
	resolved.Preflights = append(resolved.Preflights, cfg.Preflights...)
	resolved.Extractors = append(resolved.Extractors, cfg.Extractors...)

	// Merge layers (user layers override preset layers)
	for k, v := range cfg.Layers {
		resolved.Layers[k] = v
	}

	// TODO: Merge shortcuts, theme, postprocess, etc.

	return resolved
}

// NewGenerator cria uma nova instância do UnoGenerator com a configuração resolvida.
func NewGenerator(config *ResolvedConfig) *UnoGenerator {
	return &UnoGenerator{
		Config: config,
		Cache:  make(map[string][]*StringifiedUtil),
	}
}