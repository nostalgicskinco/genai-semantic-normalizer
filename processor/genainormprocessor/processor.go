package genainormprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

// defaultMappings maps common vendor-specific attributes to gen_ai.* conventions.
var defaultMappings = map[string]string{
	// OpenAI
	"openai.model":              "gen_ai.request.model",
	"openai.api_base":           "gen_ai.system",
	"openai.max_tokens":         "gen_ai.request.max_tokens",
	"openai.temperature":        "gen_ai.request.temperature",
	"openai.top_p":              "gen_ai.request.top_p",
	"openai.prompt_tokens":      "gen_ai.usage.input_tokens",
	"openai.completion_tokens":  "gen_ai.usage.output_tokens",
	"openai.total_tokens":       "gen_ai.usage.total_tokens",
	"openai.finish_reason":      "gen_ai.response.finish_reasons",

	// Anthropic
	"anthropic.model":           "gen_ai.request.model",
	"anthropic.max_tokens":      "gen_ai.request.max_tokens",
	"anthropic.input_tokens":    "gen_ai.usage.input_tokens",
	"anthropic.output_tokens":   "gen_ai.usage.output_tokens",
	"anthropic.stop_reason":     "gen_ai.response.finish_reasons",

	// Cohere
	"cohere.model_id":           "gen_ai.request.model",
	"cohere.prompt_tokens":      "gen_ai.usage.input_tokens",
	"cohere.response_tokens":    "gen_ai.usage.output_tokens",
}

func init() {
	// Azure OpenAI
	defaultMappings["az.ai.model"] = "gen_ai.request.model"
	defaultMappings["az.ai.prompt_tokens"] = "gen_ai.usage.input_tokens"
	defaultMappings["az.ai.completion_tokens"] = "gen_ai.usage.output_tokens"

	// Google / Vertex AI
	defaultMappings["google.model"] = "gen_ai.request.model"
	defaultMappings["google.prompt_token_count"] = "gen_ai.usage.input_tokens"
	defaultMappings["google.candidates_token_count"] = "gen_ai.usage.output_tokens"

	// Generic LLM attributes
	defaultMappings["llm.model"] = "gen_ai.request.model"
	defaultMappings["llm.prompt"] = "gen_ai.prompt"
	defaultMappings["llm.completion"] = "gen_ai.completion"
	defaultMappings["llm.token_count.prompt"] = "gen_ai.usage.input_tokens"
	defaultMappings["llm.token_count.completion"] = "gen_ai.usage.output_tokens"
}
type normalizerProcessor struct {
	logger       *zap.Logger
	config       *Config
	nextConsumer consumer.Traces
	mappings     map[string]string
}

func newNormalizerProcessor(
	logger *zap.Logger,
	cfg *Config,
	next consumer.Traces,
) *normalizerProcessor {
	mappings := make(map[string]string)

	if cfg.EnableDefaults {
		for k, v := range defaultMappings {
			mappings[k] = v
		}
	}

	// Custom mappings override defaults
	for k, v := range cfg.CustomMappings {
		mappings[k] = v
	}

	return &normalizerProcessor{
		logger:       logger,
		config:       cfg,
		nextConsumer: next,
		mappings:     mappings,
	}
}

func (p *normalizerProcessor) Start(_ context.Context, _ component.Host) error {
	p.logger.Info("genai_semantic_normalizer started",
		zap.Int("mapping_count", len(p.mappings)),
		zap.Bool("overwrite", p.config.Overwrite),
		zap.Bool("drop_original", p.config.DropOriginal),
	)
	return nil
}

func (p *normalizerProcessor) Shutdown(_ context.Context) error {
	return nil
}

func (p *normalizerProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}
func (p *normalizerProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		ilss := rss.At(i).ScopeSpans()
		for j := 0; j < ilss.Len(); j++ {
			spans := ilss.At(j).Spans()
			for k := 0; k < spans.Len(); k++ {
				p.normalizeSpan(spans.At(k))
			}
		}
	}
	return p.nextConsumer.ConsumeTraces(ctx, td)
}

func (p *normalizerProcessor) normalizeSpan(span ptrace.Span) {
	attrs := span.Attributes()
	for vendorKey, genaiKey := range p.mappings {
		val, exists := attrs.Get(vendorKey)
		if !exists {
			continue
		}

		// Check if target already exists
		if !p.config.Overwrite {
			if _, targetExists := attrs.Get(genaiKey); targetExists {
				continue
			}
		}

		// Copy value to normalized key
		val.CopyTo(attrs.PutEmpty(genaiKey))

		// Optionally remove the vendor-specific key
		if p.config.DropOriginal {
			attrs.Remove(vendorKey)
		}
	}

	// Infer gen_ai.system from known vendor prefixes if not set
	if _, exists := attrs.Get("gen_ai.system"); !exists {
		system := inferSystem(attrs)
		if system != "" {
			attrs.PutStr("gen_ai.system", system)
		}
	}
}

func inferSystem(attrs ptrace.Map) string {
	prefixes := map[string]string{
		"openai.":    "openai",
		"anthropic.": "anthropic",
		"cohere.":    "cohere",
		"az.ai.":     "az.ai.openai",
		"google.":    "vertex_ai",
	}

	found := ""
	attrs.Range(func(k string, _ ptrace.Value) bool {
		for prefix, system := range prefixes {
			if len(k) > len(prefix) && k[:len(prefix)] == prefix {
				found = system
				return false
			}
		}
		return true
	})
	return found
}