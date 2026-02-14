// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: Apache-2.0

package genainormprocessor

// Config holds the processor configuration.
type Config struct {
	// EnableDefaults loads the built-in mapping set for common
	// frameworks (LangChain, OpenInference, OpenLLMetry, etc.).
	EnableDefaults bool `mapstructure:"enable_defaults"`

	// Overwrite controls whether an existing canonical attribute
	// is overwritten when the source attribute is also present.
	Overwrite bool `mapstructure:"overwrite"`

	// DropOriginal removes the source attribute after mapping.
	DropOriginal bool `mapstructure:"drop_original"`

	// Mappings is a user-supplied map of source_key → canonical_key.
	// These are merged on top of defaults (if enabled).
	Mappings map[string]string `mapstructure:"mappings"`
}
