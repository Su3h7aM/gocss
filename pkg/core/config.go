package core

// Preset é uma função que aplica uma configuração pré-definida.
type Preset func(config *ResolvedConfig)

// NewGenerator cria uma nova instância do UnoGenerator com a configuração resolvida.
func NewGenerator(config *ResolvedConfig) *UnoGenerator {
	// TODO: Implementar a lógica de mesclagem de presets e configuração do usuário.
	return &UnoGenerator{
		Config: config,
		Cache:  make(map[string][]*StringifiedUtil),
	}
}
