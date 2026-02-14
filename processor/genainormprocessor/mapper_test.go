// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: Apache-2.0

package genainormprocessor

import (
	"testing"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

func newAttrs(kv map[string]interface{}) pcommon.Map {
	m := pcommon.NewMap()
	for k, v := range kv {
		switch val := v.(type) {
		case string:
			m.PutStr(k, val)
		case int64:
			m.PutInt(k, val)
		case float64:
			m.PutDouble(k, val)
		case bool:
			m.PutBool(k, val)
		}
	}
	return m
}

func TestMapper_BasicMapping(t *testing.T) {
	cfg := &Config{
		EnableDefaults: false,
		Overwrite:      false,
		DropOriginal:   false,
		Mappings: map[string]string{
			"llm.model": "gen_ai.request.model",
		},
	}
	m := newMapper(cfg)

	attrs := newAttrs(map[string]interface{}{
		"llm.model": "gpt-4o",
	})

	m.apply(attrs)

	v, ok := attrs.Get("gen_ai.request.model")
	if !ok {
		t.Fatal("canonical key gen_ai.request.model not created")
	}
	if v.Str() != "gpt-4o" {
		t.Errorf("expected gpt-4o, got %q", v.Str())
	}

	// Original should still exist (drop_original=false)
	_, ok = attrs.Get("llm.model")
	if !ok {
		t.Error("original key llm.model should still exist when drop_original=false")
	}
}

func TestMapper_DropOriginal(t *testing.T) {
	cfg := &Config{
		EnableDefaults: false,
		Overwrite:      false,
		DropOriginal:   true,
		Mappings: map[string]string{
			"llm.model": "gen_ai.request.model",
		},
	}
	m := newMapper(cfg)

	attrs := newAttrs(map[string]interface{}{
		"llm.model": "claude-3",
	})

	m.apply(attrs)

	// Canonical key should exist
	v, ok := attrs.Get("gen_ai.request.model")
	if !ok || v.Str() != "claude-3" {
		t.Error("canonical key should be set to claude-3")
	}

	// Original should be gone
	_, ok = attrs.Get("llm.model")
	if ok {
		t.Error("original key should be removed when drop_original=true")
	}
}

func TestMapper_NoOverwrite(t *testing.T) {
	cfg := &Config{
		EnableDefaults: false,
		Overwrite:      false,
		DropOriginal:   false,
		Mappings: map[string]string{
			"llm.model": "gen_ai.request.model",
		},
	}
	m := newMapper(cfg)

	attrs := newAttrs(map[string]interface{}{
		"llm.model":            "gpt-4o",
		"gen_ai.request.model": "already-set",
	})

	m.apply(attrs)

	// Should NOT overwrite existing canonical key
	v, _ := attrs.Get("gen_ai.request.model")
	if v.Str() != "already-set" {
		t.Errorf("should not overwrite, got %q", v.Str())
	}
}

func TestMapper_Overwrite(t *testing.T) {
	cfg := &Config{
		EnableDefaults: false,
		Overwrite:      true,
		DropOriginal:   false,
		Mappings: map[string]string{
			"llm.model": "gen_ai.request.model",
		},
	}
	m := newMapper(cfg)

	attrs := newAttrs(map[string]interface{}{
		"llm.model":            "gpt-4o",
		"gen_ai.request.model": "old-value",
	})

	m.apply(attrs)

	// SHOULD overwrite
	v, _ := attrs.Get("gen_ai.request.model")
	if v.Str() != "gpt-4o" {
		t.Errorf("should overwrite, got %q", v.Str())
	}
}

func TestMapper_IntValues(t *testing.T) {
	cfg := &Config{
		EnableDefaults: false,
		Overwrite:      false,
		DropOriginal:   false,
		Mappings: map[string]string{
			"llm.usage.prompt_tokens":     "gen_ai.usage.input_tokens",
			"llm.usage.completion_tokens": "gen_ai.usage.output_tokens",
		},
	}
	m := newMapper(cfg)

	attrs := newAttrs(map[string]interface{}{
		"llm.usage.prompt_tokens":     int64(150),
		"llm.usage.completion_tokens": int64(89),
	})

	m.apply(attrs)

	v, ok := attrs.Get("gen_ai.usage.input_tokens")
	if !ok || v.Int() != 150 {
		t.Errorf("expected input_tokens=150, got %v", v)
	}

	v, ok = attrs.Get("gen_ai.usage.output_tokens")
	if !ok || v.Int() != 89 {
		t.Errorf("expected output_tokens=89, got %v", v)
	}
}

func TestMapper_DoubleValues(t *testing.T) {
	cfg := &Config{
		EnableDefaults: false,
		Overwrite:      false,
		DropOriginal:   false,
		Mappings: map[string]string{
			"llm.temperature": "gen_ai.request.temperature",
		},
	}
	m := newMapper(cfg)

	attrs := newAttrs(map[string]interface{}{
		"llm.temperature": float64(0.7),
	})

	m.apply(attrs)

	v, ok := attrs.Get("gen_ai.request.temperature")
	if !ok || v.Double() != 0.7 {
		t.Errorf("expected temperature=0.7, got %v", v)
	}
}

func TestMapper_DefaultsPlusUserOverride(t *testing.T) {
	cfg := &Config{
		EnableDefaults: true,
		Overwrite:      false,
		DropOriginal:   false,
		Mappings: map[string]string{
			// Override a default mapping
			"llm.model": "gen_ai.response.model",
			// Add a custom one
			"my.custom.key": "gen_ai.request.model",
		},
	}
	m := newMapper(cfg)

	// User override should win over default
	if m.table["llm.model"] != "gen_ai.response.model" {
		t.Errorf("user mapping should override default, got %q", m.table["llm.model"])
	}

	// Custom key should exist
	if m.table["my.custom.key"] != "gen_ai.request.model" {
		t.Error("custom user mapping should be present")
	}

	// A default that wasn't overridden should still exist
	if m.table["openinference.model_name"] != "gen_ai.request.model" {
		t.Error("non-overridden default mapping should still exist")
	}
}

func TestMapper_NoMappingsNoChange(t *testing.T) {
	cfg := &Config{
		EnableDefaults: false,
		Overwrite:      false,
		DropOriginal:   false,
		Mappings:       map[string]string{},
	}
	m := newMapper(cfg)

	attrs := newAttrs(map[string]interface{}{
		"some.random.key": "value",
	})

	m.apply(attrs)

	if attrs.Len() != 1 {
		t.Errorf("expected 1 attribute, got %d", attrs.Len())
	}
}
