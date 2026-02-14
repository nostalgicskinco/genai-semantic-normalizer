// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: Apache-2.0

package genainormprocessor

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// mapper holds the compiled mapping table and applies it to attribute maps.
type mapper struct {
	table        map[string]string
	overwrite    bool
	dropOriginal bool
}

// newMapper builds a mapper from config. If enableDefaults is true,
// the built-in mappings are loaded first, then user mappings override.
func newMapper(cfg *Config) *mapper {
	table := make(map[string]string)

	if cfg.EnableDefaults {
		for src, dst := range DefaultMappings() {
			table[src] = dst
		}
	}

	// User mappings override defaults
	for src, dst := range cfg.Mappings {
		table[src] = dst
	}

	return &mapper{
		table:        table,
		overwrite:    cfg.Overwrite,
		dropOriginal: cfg.DropOriginal,
	}
}

// apply normalizes attributes in-place. For each source key found
// in the mapping table, it copies the value to the canonical key
// (respecting overwrite policy) and optionally drops the original.
func (m *mapper) apply(attrs pcommon.Map) {
	// Collect keys to process first to avoid mutating during iteration
	type pending struct {
		srcKey string
		dstKey string
		val    pcommon.Value
	}

	var ops []pending

	attrs.Range(func(k string, v pcommon.Value) bool {
		if dst, ok := m.table[k]; ok {
			ops = append(ops, pending{srcKey: k, dstKey: dst, val: v})
		}
		return true
	})

	for _, op := range ops {
		// Skip if canonical key already exists and overwrite is off
		if _, exists := attrs.Get(op.dstKey); exists && !m.overwrite {
			continue
		}

		// Copy value to canonical key
		copyValue(attrs, op.dstKey, op.val)

		// Optionally remove source key
		if m.dropOriginal {
			attrs.Remove(op.srcKey)
		}
	}
}

// copyValue writes a pcommon.Value to a map under the given key,
// preserving the original type.
func copyValue(attrs pcommon.Map, key string, v pcommon.Value) {
	switch v.Type() {
	case pcommon.ValueTypeStr:
		attrs.PutStr(key, v.Str())
	case pcommon.ValueTypeInt:
		attrs.PutInt(key, v.Int())
	case pcommon.ValueTypeDouble:
		attrs.PutDouble(key, v.Double())
	case pcommon.ValueTypeBool:
		attrs.PutBool(key, v.Bool())
	default:
		// For slices, maps, bytes — copy as string repr
		attrs.PutStr(key, v.AsString())
	}
}
