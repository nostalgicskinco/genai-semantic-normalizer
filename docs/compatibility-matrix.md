# Compatibility Matrix

Which framework/vendor attributes map to which `gen_ai.*` canonical keys.

## Model

| Source Key | Canonical Key | Framework |
|---|---|---|
| `llm.model` | `gen_ai.request.model` | OpenLLMetry |
| `llm.request.model` | `gen_ai.request.model` | OpenLLMetry |
| `openinference.model_name` | `gen_ai.request.model` | OpenInference / Arize |
| `traceloop.entity.model` | `gen_ai.request.model` | Traceloop |
| `langchain.llm.model_name` | `gen_ai.request.model` | LangChain |
| `litellm.model` | `gen_ai.request.model` | LiteLLM |
| `model_name` | `gen_ai.request.model` | Generic |
| `model` | `gen_ai.request.model` | Generic |

## Provider / System

| Source Key | Canonical Key | Framework |
|---|---|---|
| `llm.vendor` | `gen_ai.system` | OpenLLMetry |
| `llm.provider` | `gen_ai.system` | OpenLLMetry |
| `openinference.llm.provider` | `gen_ai.system` | OpenInference |
| `traceloop.entity.provider` | `gen_ai.system` | Traceloop |
| `langchain.llm.provider` | `gen_ai.system` | LangChain |
| `litellm.provider` | `gen_ai.system` | LiteLLM |

## Token Usage: Input (Prompt)

| Source Key | Canonical Key | Framework |
|---|---|---|
| `llm.usage.prompt_tokens` | `gen_ai.usage.input_tokens` | OpenLLMetry |
| `llm.token_count.prompt` | `gen_ai.usage.input_tokens` | OpenLLMetry |
| `openinference.llm.token_count.prompt` | `gen_ai.usage.input_tokens` | OpenInference |
| `traceloop.entity.input_tokens` | `gen_ai.usage.input_tokens` | Traceloop |
| `langchain.tokens.prompt` | `gen_ai.usage.input_tokens` | LangChain |
| `litellm.usage.prompt_tokens` | `gen_ai.usage.input_tokens` | LiteLLM |
| `prompt_tokens` | `gen_ai.usage.input_tokens` | Generic |
| `input_tokens` | `gen_ai.usage.input_tokens` | Generic |

## Token Usage: Output (Completion)

| Source Key | Canonical Key | Framework |
|---|---|---|
| `llm.usage.completion_tokens` | `gen_ai.usage.output_tokens` | OpenLLMetry |
| `llm.token_count.completion` | `gen_ai.usage.output_tokens` | OpenLLMetry |
| `openinference.llm.token_count.completion` | `gen_ai.usage.output_tokens` | OpenInference |
| `traceloop.entity.output_tokens` | `gen_ai.usage.output_tokens` | Traceloop |
| `langchain.tokens.completion` | `gen_ai.usage.output_tokens` | LangChain |
| `litellm.usage.completion_tokens` | `gen_ai.usage.output_tokens` | LiteLLM |
| `completion_tokens` | `gen_ai.usage.output_tokens` | Generic |
| `output_tokens` | `gen_ai.usage.output_tokens` | Generic |

## Request Parameters

| Source Key | Canonical Key | Framework |
|---|---|---|
| `llm.temperature` | `gen_ai.request.temperature` | OpenLLMetry |
| `llm.request.temperature` | `gen_ai.request.temperature` | OpenLLMetry |
| `openinference.llm.temperature` | `gen_ai.request.temperature` | OpenInference |
| `llm.max_tokens` | `gen_ai.request.max_tokens` | OpenLLMetry |
| `llm.request.max_tokens` | `gen_ai.request.max_tokens` | OpenLLMetry |
| `openinference.llm.max_tokens` | `gen_ai.request.max_tokens` | OpenInference |
| `max_tokens` | `gen_ai.request.max_tokens` | Generic |
| `llm.top_p` | `gen_ai.request.top_p` | OpenLLMetry |
| `llm.request.top_p` | `gen_ai.request.top_p` | OpenLLMetry |
| `top_p` | `gen_ai.request.top_p` | Generic |
| `llm.frequency_penalty` | `gen_ai.request.frequency_penalty` | OpenLLMetry |
| `llm.presence_penalty` | `gen_ai.request.presence_penalty` | OpenLLMetry |

## Response

| Source Key | Canonical Key | Framework |
|---|---|---|
| `llm.response.finish_reason` | `gen_ai.response.finish_reasons` | OpenLLMetry |
| `llm.finish_reason` | `gen_ai.response.finish_reasons` | OpenLLMetry |
| `finish_reason` | `gen_ai.response.finish_reasons` | Generic |
| `llm.response.model` | `gen_ai.response.model` | OpenLLMetry |
| `llm.response.id` | `gen_ai.response.id` | OpenLLMetry |

## Prompt / Completion Content

| Source Key | Canonical Key | Framework |
|---|---|---|
| `llm.prompt` | `gen_ai.prompt` | OpenLLMetry |
| `llm.completion` | `gen_ai.completion` | OpenLLMetry |
| `openinference.input.value` | `gen_ai.prompt` | OpenInference |
| `openinference.output.value` | `gen_ai.completion` | OpenInference |
| `traceloop.entity.input` | `gen_ai.prompt` | Traceloop |
| `traceloop.entity.output` | `gen_ai.completion` | Traceloop |

## Cost

| Source Key | Canonical Key | Framework |
|---|---|---|
| `llm.usage.cost` | `gen_ai.usage.cost` | OpenLLMetry |
| `llm.usage.cost_usd` | `gen_ai.usage.cost` | OpenLLMetry |
| `gen_ai.usage.cost_usd` | `gen_ai.usage.cost` | Custom |
| `litellm.cost` | `gen_ai.usage.cost` | LiteLLM |

---

**Reference**: [OTel GenAI Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/gen-ai/)
