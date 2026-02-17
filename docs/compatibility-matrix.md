# Compatibility matrix (starter)

This matrix is intentionally minimal for the MVP.

| Source key (observed in the wild) | Canonical key (OTel GenAI semconv) |
|---|---|
| `llm.provider` | `gen_ai.provider.name` |
| `llm.model_name` | `gen_ai.request.model` |
| `llm.model` | `gen_ai.request.model` |
| `llm.operation` | `gen_ai.operation.name` |
| `llm.operation_name` | `gen_ai.operation.name` |
| `llm.usage.prompt_tokens` | `gen_ai.usage.input_tokens` |
| `llm.usage.completion_tokens` | `gen_ai.usage.output_tokens` |
| `llm.usage.total_tokens` | `gen_ai.usage.total_tokens` |
| `llm.token_count.prompt` | `gen_ai.usage.input_tokens` |
| `llm.token_count.completion` | `gen_ai.usage.output_tokens` |
| `llm.token_count.total` | `gen_ai.usage.total_tokens` |
| `llm.request.temperature` | `gen_ai.request.temperature` |
| `llm.request.top_p` | `gen_ai.request.top_p` |
| `llm.request.max_tokens` | `gen_ai.request.max_tokens` |
| `llm.response.model` | `gen_ai.response.model` |

Add your organization/framework keys via `processors.genai_semantic_normalizer.mappings`.
