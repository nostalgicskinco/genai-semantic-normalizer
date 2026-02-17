package genainormalizerprocessor

// defaultMappings provides a minimal, opinionated starting point.
// It is intentionally small: users should extend it as their stack dictates.
//
// The destination keys follow the OpenTelemetry GenAI semantic conventions (gen_ai.*).
// See: https://opentelemetry.io/docs/specs/semconv/registry/attributes/gen-ai/
var defaultMappings = map[string]string{
	// Common "LLM" / OpenInference-style keys.
	"llm.provider":   "gen_ai.provider.name",
	"llm.model_name": "gen_ai.request.model",
	"llm.model":      "gen_ai.request.model",

	// Operation names
	"llm.operation":   "gen_ai.operation.name",
	"llm.operation_name": "gen_ai.operation.name",

	// Token usage (frequently emitted by libraries with slightly different naming)
	"llm.usage.prompt_tokens":     "gen_ai.usage.input_tokens",
	"llm.usage.completion_tokens": "gen_ai.usage.output_tokens",
	"llm.usage.total_tokens":      "gen_ai.usage.total_tokens",

	"llm.token_count.prompt":     "gen_ai.usage.input_tokens",
	"llm.token_count.completion": "gen_ai.usage.output_tokens",
	"llm.token_count.total":      "gen_ai.usage.total_tokens",

	// Request parameters (common across providers)
	"llm.request.temperature": "gen_ai.request.temperature",
	"llm.request.top_p":       "gen_ai.request.top_p",
	"llm.request.max_tokens":  "gen_ai.request.max_tokens",

	// Response model (some SDKs distinguish request vs response model)
	"llm.response.model": "gen_ai.response.model",
}
