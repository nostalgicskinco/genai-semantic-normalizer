package genainormalizerprocessor

import (
	"context"
	"testing"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
)

func TestApplyMappings_SpanAndEvents(t *testing.T) {
	next := consumertest.NewNop()
	settings := processor.CreateSettings{TelemetrySettings: component.TelemetrySettings{Logger: zap.NewNop()}}
	cfg := &Config{
		EnableDefaults: false,
		Overwrite:      false,
		DropOriginal:   false,
		Mappings: map[string]string{
			"llm.model_name": "gen_ai.request.model",
			"llm.provider":   "gen_ai.provider.name",
		},
	}

	p, err := newTracesProcessor(context.Background(), settings, cfg, next)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	td := ptrace.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()
	ss := rs.ScopeSpans().AppendEmpty()
	sp := ss.Spans().AppendEmpty()
	sp.Attributes().PutStr("llm.model_name", "gpt-4.1")
	sp.Attributes().PutStr("llm.provider", "openai")

	ev := sp.Events().AppendEmpty()
	ev.Attributes().PutStr("llm.model_name", "gpt-4.1")

	if err := p.ConsumeTraces(context.Background(), td); err != nil {
		t.Fatalf("consume err: %v", err)
	}

	attrs := sp.Attributes()
	assertHasStr(t, attrs, "gen_ai.request.model", "gpt-4.1")
	assertHasStr(t, attrs, "gen_ai.provider.name", "openai")

	eattrs := ev.Attributes()
	assertHasStr(t, eattrs, "gen_ai.request.model", "gpt-4.1")
}

func TestOverwriteAndDropOriginal(t *testing.T) {
	next := consumertest.NewNop()
	settings := processor.CreateSettings{TelemetrySettings: component.TelemetrySettings{Logger: zap.NewNop()}}
	cfg := &Config{
		EnableDefaults: false,
		Overwrite:      true,
		DropOriginal:   true,
		Mappings: map[string]string{
			"llm.model_name": "gen_ai.request.model",
		},
	}

	p, err := newTracesProcessor(context.Background(), settings, cfg, next)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	td := ptrace.NewTraces()
	sp := td.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	sp.Attributes().PutStr("llm.model_name", "gpt-4.1")
	sp.Attributes().PutStr("gen_ai.request.model", "already")

	if err := p.ConsumeTraces(context.Background(), td); err != nil {
		t.Fatalf("consume err: %v", err)
	}

	attrs := sp.Attributes()
	assertHasStr(t, attrs, "gen_ai.request.model", "gpt-4.1")
	if _, ok := attrs.Get("llm.model_name"); ok {
		t.Fatalf("expected original key dropped")
	}
}

func assertHasStr(t *testing.T, m pcommon.Map, key, want string) {
	v, ok := m.Get(key)
	if !ok {
		t.Fatalf("expected key %q", key)
	}
	if v.Str() != want {
		t.Fatalf("%q: got %q want %q", key, v.Str(), want)
	}
}

// Ensure factory compiles.
func TestFactoryType(t *testing.T) {
	f := NewFactory()
	if f.Type() != component.MustNewType(typeStr) {
		t.Fatalf("unexpected type")
	}
}
