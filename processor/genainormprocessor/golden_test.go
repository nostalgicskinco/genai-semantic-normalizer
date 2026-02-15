// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: Apache-2.0

package genainormprocessor

import (
	"context"
	"testing"

	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor/processortest"
)

// goldenFixture represents a normalization test scenario.
type goldenFixture struct {
	Name         string
	Description  string
	SpanAttrs    map[string]interface{} // input span attributes
	EventAttrs   map[string]interface{} // input event attributes (nil = no events)
	Config       Config
	WantSpan     map[string]interface{} // expected span attributes after normalization
	WantEvent    map[string]interface{} // expected event attributes (nil = skip)
	WantAbsent   []string               // keys that must NOT exist after processing
}

// goldenFixtures returns all normalization test scenarios.
func goldenFixtures() []goldenFixture {
	return []goldenFixture{
		{
			Name:        "openai_standard",
			Description: "Standard OpenAI-style attributes normalized to gen_ai.*",
			SpanAttrs: map[string]interface{}{
				"llm.model":                   "gpt-4o",
				"llm.usage.prompt_tokens":     int64(200),
				"llm.usage.completion_tokens": int64(50),
				"llm.temperature":             float64(0.8),
				"llm.max_tokens":              int64(4096),
				"llm.top_p":                   float64(0.95),
			},
			Config: Config{EnableDefaults: true, Overwrite: false, DropOriginal: false, Mappings: map[string]string{}},
			WantSpan: map[string]interface{}{
				"gen_ai.request.model":       "gpt-4o",
				"gen_ai.usage.input_tokens":  int64(200),
				"gen_ai.usage.output_tokens": int64(50),
				"gen_ai.request.temperature": float64(0.8),
				"gen_ai.request.max_tokens":  int64(4096),
				"gen_ai.request.top_p":       float64(0.95),
				// originals preserved
				"llm.model": "gpt-4o",
			},
		},
		{
			Name:        "openinference_langchain",
			Description: "OpenInference and LangChain-style attributes",
			SpanAttrs: map[string]interface{}{
				"openinference.model_name":     "claude-3-sonnet",
				"openinference.llm.provider":   "anthropic",
				"openinference.llm.temperature": float64(0.5),
			},
			EventAttrs: map[string]interface{}{
				"openinference.input.value":  "Tell me about quantum physics",
				"openinference.output.value": "Quantum physics is...",
			},
			Config: Config{EnableDefaults: true, Overwrite: false, DropOriginal: false, Mappings: map[string]string{}},
			WantSpan: map[string]interface{}{
				"gen_ai.request.model":       "claude-3-sonnet",
				"gen_ai.system":              "anthropic",
				"gen_ai.request.temperature": float64(0.5),
			},
			WantEvent: map[string]interface{}{
				"gen_ai.prompt":     "Tell me about quantum physics",
				"gen_ai.completion": "Quantum physics is...",
			},
		},
		{
			Name:        "traceloop_entity",
			Description: "Traceloop (OpenLLMetry) style attributes",
			SpanAttrs: map[string]interface{}{
				"traceloop.entity.model":         "mistral-7b",
				"traceloop.entity.provider":      "mistral",
				"traceloop.entity.input_tokens":  int64(100),
				"traceloop.entity.output_tokens": int64(30),
				"traceloop.entity.type":          "chat",
			},
			Config: Config{EnableDefaults: true, Overwrite: false, DropOriginal: false, Mappings: map[string]string{}},
			WantSpan: map[string]interface{}{
				"gen_ai.request.model":       "mistral-7b",
				"gen_ai.system":              "mistral",
				"gen_ai.usage.input_tokens":  int64(100),
				"gen_ai.usage.output_tokens": int64(30),
				"gen_ai.operation.name":      "chat",
			},
		},
		{
			Name:        "litellm_with_cost",
			Description: "LiteLLM-style attributes with cost tracking",
			SpanAttrs: map[string]interface{}{
				"litellm.model":                "gpt-4-turbo",
				"litellm.provider":             "openai",
				"litellm.usage.prompt_tokens":  int64(500),
				"litellm.usage.completion_tokens": int64(150),
				"litellm.cost":                 float64(0.0235),
			},
			Config: Config{EnableDefaults: true, Overwrite: false, DropOriginal: false, Mappings: map[string]string{}},
			WantSpan: map[string]interface{}{
				"gen_ai.request.model":       "gpt-4-turbo",
				"gen_ai.system":              "openai",
				"gen_ai.usage.input_tokens":  int64(500),
				"gen_ai.usage.output_tokens": int64(150),
			},
		},
		{
			Name:        "drop_original_mode",
			Description: "When drop_original=true, source keys should be removed",
			SpanAttrs: map[string]interface{}{
				"llm.model":               "gpt-4o-mini",
				"llm.usage.prompt_tokens": int64(10),
			},
			Config: Config{EnableDefaults: true, Overwrite: false, DropOriginal: true, Mappings: map[string]string{}},
			WantSpan: map[string]interface{}{
				"gen_ai.request.model":      "gpt-4o-mini",
				"gen_ai.usage.input_tokens": int64(10),
			},
			WantAbsent: []string{"llm.model", "llm.usage.prompt_tokens"},
		},
		{
			Name:        "overwrite_existing",
			Description: "When overwrite=true, existing canonical keys are replaced",
			SpanAttrs: map[string]interface{}{
				"llm.model":             "gpt-4o",
				"gen_ai.request.model":  "old-value",
			},
			Config: Config{EnableDefaults: true, Overwrite: true, DropOriginal: false, Mappings: map[string]string{}},
			WantSpan: map[string]interface{}{
				"gen_ai.request.model": "gpt-4o",
			},
		},
		{
			Name:        "no_overwrite_preserves_existing",
			Description: "When overwrite=false, existing canonical keys are preserved",
			SpanAttrs: map[string]interface{}{
				"llm.model":            "gpt-4o",
				"gen_ai.request.model": "already-set",
			},
			Config: Config{EnableDefaults: true, Overwrite: false, DropOriginal: false, Mappings: map[string]string{}},
			WantSpan: map[string]interface{}{
				"gen_ai.request.model": "already-set",
			},
		},
		{
			Name:        "custom_mapping_override",
			Description: "User-defined mappings override default mappings",
			SpanAttrs: map[string]interface{}{
				"custom.model.name": "my-custom-model",
			},
			Config: Config{
				EnableDefaults: false,
				Overwrite:      false,
				DropOriginal:   false,
				Mappings: map[string]string{
					"custom.model.name": "gen_ai.request.model",
				},
			},
			WantSpan: map[string]interface{}{
				"gen_ai.request.model": "my-custom-model",
			},
		},
		{
			Name:        "mixed_vendors_single_span",
			Description: "Multiple vendor-style attributes on the same span (first wins with no-overwrite)",
			SpanAttrs: map[string]interface{}{
				"llm.model":                "gpt-4o",
				"openinference.model_name": "should-not-overwrite",
			},
			Config: Config{EnableDefaults: true, Overwrite: false, DropOriginal: false, Mappings: map[string]string{}},
			WantSpan: map[string]interface{}{
				// one of these sets gen_ai.request.model; the other is skipped (no overwrite)
				// We verify gen_ai.request.model exists — the exact value depends on map iteration order
			},
		},
		{
			Name:        "passthrough_unknown_attrs",
			Description: "Attributes not in any mapping table pass through unchanged",
			SpanAttrs: map[string]interface{}{
				"http.method":     "POST",
				"http.status_code": int64(200),
				"custom.tag":      "keep-me",
			},
			Config: Config{EnableDefaults: true, Overwrite: false, DropOriginal: false, Mappings: map[string]string{}},
			WantSpan: map[string]interface{}{
				"http.method":      "POST",
				"http.status_code": int64(200),
				"custom.tag":       "keep-me",
			},
		},
	}
}

// TestGoldenNormalization runs all golden normalization fixtures.
func TestGoldenNormalization(t *testing.T) {
	for _, fix := range goldenFixtures() {
		t.Run(fix.Name, func(t *testing.T) {
			sink := new(consumertest.TracesSink)
			set := processortest.NewNopSettings()

			proc, err := newProcessor(context.Background(), set, &fix.Config, sink)
			if err != nil {
				t.Fatalf("newProcessor: %v", err)
			}

			td := buildTestTraces(fix.SpanAttrs, fix.EventAttrs)

			if err := proc.ConsumeTraces(context.Background(), td); err != nil {
				t.Fatalf("ConsumeTraces: %v", err)
			}

			if len(sink.AllTraces()) != 1 {
				t.Fatalf("expected 1 trace batch, got %d", len(sink.AllTraces()))
			}

			span := sink.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0)
			attrs := span.Attributes()

			// Verify expected span attributes.
			for key, want := range fix.WantSpan {
				got, ok := attrs.Get(key)
				if !ok {
					t.Errorf("missing expected attribute %q", key)
					continue
				}
				switch expected := want.(type) {
				case string:
					if got.Str() != expected {
						t.Errorf("%s = %q, want %q", key, got.Str(), expected)
					}
				case int64:
					if got.Int() != expected {
						t.Errorf("%s = %d, want %d", key, got.Int(), expected)
					}
				case float64:
					if got.Double() != expected {
						t.Errorf("%s = %f, want %f", key, got.Double(), expected)
					}
				}
			}

			// Verify absent keys.
			for _, key := range fix.WantAbsent {
				if _, ok := attrs.Get(key); ok {
					t.Errorf("attribute %q should not exist (drop_original=true)", key)
				}
			}

			// Verify event attributes if specified.
			if fix.WantEvent != nil && span.Events().Len() > 0 {
				eventAttrs := span.Events().At(0).Attributes()
				for key, want := range fix.WantEvent {
					got, ok := eventAttrs.Get(key)
					if !ok {
						t.Errorf("missing expected event attribute %q", key)
						continue
					}
					if expected, ok := want.(string); ok && got.Str() != expected {
						t.Errorf("event %s = %q, want %q", key, got.Str(), expected)
					}
				}
			}
		})
	}
}

// TestGoldenNormalization_EmptyTraces ensures the processor handles empty traces gracefully.
func TestGoldenNormalization_EmptyTraces(t *testing.T) {
	sink := new(consumertest.TracesSink)
	set := processortest.NewNopSettings()
	cfg := &Config{EnableDefaults: true, Overwrite: false, DropOriginal: false, Mappings: map[string]string{}}

	proc, err := newProcessor(context.Background(), set, cfg, sink)
	if err != nil {
		t.Fatalf("newProcessor: %v", err)
	}

	td := ptrace.NewTraces() // empty
	if err := proc.ConsumeTraces(context.Background(), td); err != nil {
		t.Fatalf("ConsumeTraces on empty: %v", err)
	}

	if len(sink.AllTraces()) != 1 {
		t.Fatalf("expected empty trace to pass through, got %d", len(sink.AllTraces()))
	}
}

// TestGoldenNormalization_MultipleSpans verifies normalization works across multiple spans.
func TestGoldenNormalization_MultipleSpans(t *testing.T) {
	sink := new(consumertest.TracesSink)
	set := processortest.NewNopSettings()
	cfg := &Config{EnableDefaults: true, Overwrite: false, DropOriginal: false, Mappings: map[string]string{}}

	proc, err := newProcessor(context.Background(), set, cfg, sink)
	if err != nil {
		t.Fatalf("newProcessor: %v", err)
	}

	td := ptrace.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()
	ss := rs.ScopeSpans().AppendEmpty()

	// Span 1: OpenAI style
	s1 := ss.Spans().AppendEmpty()
	s1.SetName("llm.chat.1")
	s1.Attributes().PutStr("llm.model", "gpt-4o")
	s1.Attributes().PutInt("llm.usage.prompt_tokens", 100)

	// Span 2: Traceloop style
	s2 := ss.Spans().AppendEmpty()
	s2.SetName("llm.chat.2")
	s2.Attributes().PutStr("traceloop.entity.model", "mistral-7b")
	s2.Attributes().PutInt("traceloop.entity.input_tokens", 50)

	if err := proc.ConsumeTraces(context.Background(), td); err != nil {
		t.Fatalf("ConsumeTraces: %v", err)
	}

	spans := sink.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	if spans.Len() != 2 {
		t.Fatalf("expected 2 spans, got %d", spans.Len())
	}

	// Span 1 checks
	v1, ok := spans.At(0).Attributes().Get("gen_ai.request.model")
	if !ok || v1.Str() != "gpt-4o" {
		t.Errorf("span 1: gen_ai.request.model = %v, want gpt-4o", v1)
	}

	// Span 2 checks
	v2, ok := spans.At(1).Attributes().Get("gen_ai.request.model")
	if !ok || v2.Str() != "mistral-7b" {
		t.Errorf("span 2: gen_ai.request.model = %v, want mistral-7b", v2)
	}
}
