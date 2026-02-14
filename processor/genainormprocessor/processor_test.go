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

func buildTestTraces(spanAttrs map[string]interface{}, eventAttrs map[string]interface{}) ptrace.Traces {
	td := ptrace.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()
	ss := rs.ScopeSpans().AppendEmpty()
	s := ss.Spans().AppendEmpty()
	s.SetName("llm.chat")
	for k, v := range spanAttrs {
		switch val := v.(type) {
		case string:
			s.Attributes().PutStr(k, val)
		case int64:
			s.Attributes().PutInt(k, val)
		case float64:
			s.Attributes().PutDouble(k, val)
		}
	}
	if eventAttrs != nil {
		evt := s.Events().AppendEmpty()
		evt.SetName("gen_ai.prompt")
		for k, v := range eventAttrs {
			switch val := v.(type) {
			case string:
				evt.Attributes().PutStr(k, val)
			case int64:
				evt.Attributes().PutInt(k, val)
			}
		}
	}
	return td
}

func TestProcessor_EndToEnd(t *testing.T) {
	sink := new(consumertest.TracesSink)
	set := processortest.NewNopSettings()

	cfg := &Config{
		EnableDefaults: true,
		Overwrite:      false,
		DropOriginal:   false,
		Mappings:       map[string]string{},
	}

	proc, err := newProcessor(context.Background(), set, cfg, sink)
	if err != nil {
		t.Fatalf("failed to create processor: %v", err)
	}

	td := buildTestTraces(
		map[string]interface{}{
			"llm.model":                    "gpt-4o",
			"llm.usage.prompt_tokens":      int64(200),
			"llm.usage.completion_tokens":  int64(50),
			"llm.temperature":              float64(0.8),
		},
		nil,
	)

	err = proc.ConsumeTraces(context.Background(), td)
	if err != nil {
		t.Fatalf("ConsumeTraces failed: %v", err)
	}

	// Verify the sink received the traces
	if len(sink.AllTraces()) != 1 {
		t.Fatalf("expected 1 trace batch, got %d", len(sink.AllTraces()))
	}

	spans := sink.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	attrs := spans.At(0).Attributes()

	// Check canonical keys were created
	v, ok := attrs.Get("gen_ai.request.model")
	if !ok || v.Str() != "gpt-4o" {
		t.Errorf("expected gen_ai.request.model=gpt-4o, got %v", v)
	}

	v, ok = attrs.Get("gen_ai.usage.input_tokens")
	if !ok || v.Int() != 200 {
		t.Errorf("expected gen_ai.usage.input_tokens=200, got %v", v)
	}

	v, ok = attrs.Get("gen_ai.usage.output_tokens")
	if !ok || v.Int() != 50 {
		t.Errorf("expected gen_ai.usage.output_tokens=50, got %v", v)
	}

	v, ok = attrs.Get("gen_ai.request.temperature")
	if !ok || v.Double() != 0.8 {
		t.Errorf("expected gen_ai.request.temperature=0.8, got %v", v)
	}

	// Originals should still exist (drop_original=false)
	_, ok = attrs.Get("llm.model")
	if !ok {
		t.Error("original llm.model should still exist")
	}
}

func TestProcessor_SpanEvents(t *testing.T) {
	sink := new(consumertest.TracesSink)
	set := processortest.NewNopSettings()

	cfg := &Config{
		EnableDefaults: false,
		Overwrite:      false,
		DropOriginal:   true,
		Mappings: map[string]string{
			"openinference.input.value": "gen_ai.prompt",
		},
	}

	proc, err := newProcessor(context.Background(), set, cfg, sink)
	if err != nil {
		t.Fatalf("failed to create processor: %v", err)
	}

	td := buildTestTraces(
		map[string]interface{}{},
		map[string]interface{}{
			"openinference.input.value": "Tell me about dogs",
		},
	)

	err = proc.ConsumeTraces(context.Background(), td)
	if err != nil {
		t.Fatalf("ConsumeTraces failed: %v", err)
	}

	event := sink.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Events().At(0)
	eAttrs := event.Attributes()

	v, ok := eAttrs.Get("gen_ai.prompt")
	if !ok || v.Str() != "Tell me about dogs" {
		t.Errorf("expected event attr gen_ai.prompt, got %v", v)
	}

	// Original should be dropped
	_, ok = eAttrs.Get("openinference.input.value")
	if ok {
		t.Error("original event attr should be dropped with drop_original=true")
	}
}

func TestProcessor_Capabilities(t *testing.T) {
	sink := new(consumertest.TracesSink)
	set := processortest.NewNopSettings()
	cfg := createDefaultConfig().(*Config)

	proc, _ := newProcessor(context.Background(), set, cfg, sink)
	if !proc.Capabilities().MutatesData {
		t.Error("processor should report MutatesData=true")
	}
}
