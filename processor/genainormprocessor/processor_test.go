package genainormprocessor

import (
	"context"
	"testing"

	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

func TestNormalizeOpenAIAttributes(t *testing.T) {
	cfg := createDefaultConfig()
	sink := new(consumertest.TracesSink)
	proc := newNormalizerProcessor(zap.NewNop(), cfg, sink)

	td := ptrace.NewTraces()
	span := td.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	span.SetName("chat")
	span.Attributes().PutStr("openai.model", "gpt-4o")
	span.Attributes().PutInt("openai.prompt_tokens", 100)
	span.Attributes().PutInt("openai.completion_tokens", 50)

	err := proc.ConsumeTraces(context.Background(), td)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := sink.AllTraces()
	if len(got) != 1 {
		t.Fatalf("expected 1 trace, got %d", len(got))
	}

	attrs := got[0].ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes()

	model, ok := attrs.Get("gen_ai.request.model")
	if !ok || model.Str() != "gpt-4o" {
		t.Errorf("expected gen_ai.request.model=gpt-4o, got %v", model)
	}

	inputTok, ok := attrs.Get("gen_ai.usage.input_tokens")
	if !ok || inputTok.Int() != 100 {
		t.Errorf("expected gen_ai.usage.input_tokens=100, got %v", inputTok)
	}

	system, ok := attrs.Get("gen_ai.system")
	if !ok || system.Str() != "openai" {
		t.Errorf("expected gen_ai.system=openai, got %v", system)
	}
}
func TestNormalizeAnthropicAttributes(t *testing.T) {
	cfg := createDefaultConfig()
	sink := new(consumertest.TracesSink)
	proc := newNormalizerProcessor(zap.NewNop(), cfg, sink)

	td := ptrace.NewTraces()
	span := td.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	span.SetName("complete")
	span.Attributes().PutStr("anthropic.model", "claude-3-opus")
	span.Attributes().PutInt("anthropic.input_tokens", 200)

	err := proc.ConsumeTraces(context.Background(), td)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	attrs := sink.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes()

	model, ok := attrs.Get("gen_ai.request.model")
	if !ok || model.Str() != "claude-3-opus" {
		t.Errorf("expected gen_ai.request.model=claude-3-opus, got %v", model)
	}

	system, ok := attrs.Get("gen_ai.system")
	if !ok || system.Str() != "anthropic" {
		t.Errorf("expected gen_ai.system=anthropic, got %v", system)
	}
}

func TestDropOriginal(t *testing.T) {
	cfg := createDefaultConfig()
	cfg.DropOriginal = true
	sink := new(consumertest.TracesSink)
	proc := newNormalizerProcessor(zap.NewNop(), cfg, sink)

	td := ptrace.NewTraces()
	span := td.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	span.Attributes().PutStr("openai.model", "gpt-4o")

	err := proc.ConsumeTraces(context.Background(), td)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	attrs := sink.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes()

	if _, ok := attrs.Get("openai.model"); ok {
		t.Error("expected openai.model to be removed when drop_original=true")
	}

	if _, ok := attrs.Get("gen_ai.request.model"); !ok {
		t.Error("expected gen_ai.request.model to exist")
	}
}
func TestNoOverwrite(t *testing.T) {
	cfg := createDefaultConfig()
	cfg.Overwrite = false
	sink := new(consumertest.TracesSink)
	proc := newNormalizerProcessor(zap.NewNop(), cfg, sink)

	td := ptrace.NewTraces()
	span := td.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	span.Attributes().PutStr("gen_ai.request.model", "existing-model")
	span.Attributes().PutStr("openai.model", "gpt-4o")

	err := proc.ConsumeTraces(context.Background(), td)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	attrs := sink.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes()

	model, _ := attrs.Get("gen_ai.request.model")
	if model.Str() != "existing-model" {
		t.Errorf("expected gen_ai.request.model=existing-model (no overwrite), got %v", model.Str())
	}
}

func TestCustomMappings(t *testing.T) {
	cfg := createDefaultConfig()
	cfg.CustomMappings = map[string]string{
		"my_custom.model_name": "gen_ai.request.model",
	}
	sink := new(consumertest.TracesSink)
	proc := newNormalizerProcessor(zap.NewNop(), cfg, sink)

	td := ptrace.NewTraces()
	span := td.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	span.Attributes().PutStr("my_custom.model_name", "custom-model-v2")

	err := proc.ConsumeTraces(context.Background(), td)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	attrs := sink.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes()

	model, ok := attrs.Get("gen_ai.request.model")
	if !ok || model.Str() != "custom-model-v2" {
		t.Errorf("expected gen_ai.request.model=custom-model-v2, got %v", model)
	}
}