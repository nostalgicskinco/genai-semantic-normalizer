package genainormalizerprocessor

// Config defines configuration for the GenAI semantic normalizer.
//
// The processor maps framework/vendor-specific GenAI attribute names to the
// OpenTelemetry GenAI semantic convention keys (gen_ai.*).
//
// Example:
//   mappings:
//     llm.model_name: gen_ai.request.model
//     llm.provider:   gen_ai.provider.name
//   overwrite: false
//   drop_original: false
//
// Notes:
// - Mapping is applied to span attributes AND span event attributes.
// - If overwrite is false and destination already exists, the destination is left untouched.
// - If drop_original is true, the source key is removed when it differs from the destination.

type Config struct {
	// Mappings is a map of source_attribute_key -> destination_attribute_key.
	Mappings map[string]string `mapstructure:"mappings"`

	// Overwrite controls whether existing destination keys may be overwritten.
	Overwrite bool `mapstructure:"overwrite"`

	// DropOriginal controls whether to remove the original (source) attribute when mapped.
	DropOriginal bool `mapstructure:"drop_original"`

	// EnableDefaults toggles the built-in default mapping set.
	// If false, only user-provided mappings are used.
	EnableDefaults bool `mapstructure:"enable_defaults"`
}

func createDefaultConfig() *Config {
	return &Config{
		Mappings:       map[string]string{},
		Overwrite:      false,
		DropOriginal:   false,
		EnableDefaults: true,
	}
}
