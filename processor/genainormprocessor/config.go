package genainormprocessor

// Config holds the configuration for the genai semantic normalizer processor.
type Config struct {
	// EnableDefaults enables the built-in vendor→gen_ai mapping table.
	EnableDefaults bool `mapstructure:"enable_defaults"`

	// Overwrite controls whether existing gen_ai.* attributes get overwritten.
	Overwrite bool `mapstructure:"overwrite"`

	// DropOriginal removes the vendor-specific attribute after normalization.
	DropOriginal bool `mapstructure:"drop_original"`

	// CustomMappings allows user-defined attribute mappings.
	// Key = vendor attribute, Value = gen_ai target attribute.
	CustomMappings map[string]string `mapstructure:"custom_mappings"`
}

func createDefaultConfig() *Config {
	return &Config{
		EnableDefaults: true,
		Overwrite:      false,
		DropOriginal:   false,
		CustomMappings: make(map[string]string),
	}
}