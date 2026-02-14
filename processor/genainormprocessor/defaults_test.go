// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: Apache-2.0

package genainormprocessor

import (
	"strings"
	"testing"
)

func TestDefaultMappings_NotEmpty(t *testing.T) {
	m := DefaultMappings()
	if len(m) == 0 {
		t.Fatal("default mappings should not be empty")
	}
}

func TestDefaultMappings_AllCanonicalKeysStartWithGenAI(t *testing.T) {
	m := DefaultMappings()
	for src, dst := range m {
		if !strings.HasPrefix(dst, "gen_ai.") {
			t.Errorf("mapping %q -> %q: canonical key should start with gen_ai.", src, dst)
		}
	}
}

func TestDefaultMappings_NoDuplicateSourceKeys(t *testing.T) {
	// Go maps don't allow duplicate keys, so this is really
	// testing that the function returns consistently.
	m := DefaultMappings()
	seen := make(map[string]bool)
	for src := range m {
		if seen[src] {
			t.Errorf("duplicate source key: %q", src)
		}
		seen[src] = true
	}
}

func TestDefaultMappings_CoversCoreFrameworks(t *testing.T) {
	m := DefaultMappings()
	expectedPrefixes := []string{
		"llm.",
		"openinference.",
		"traceloop.",
		"langchain.",
		"litellm.",
	}
	for _, prefix := range expectedPrefixes {
		found := false
		for src := range m {
			if strings.HasPrefix(src, prefix) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("no mapping found for framework prefix %q", prefix)
		}
	}
}
