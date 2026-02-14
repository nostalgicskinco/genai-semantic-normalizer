# genai-semantic-normalizer

**OTel Collector processor that normalizes vendor/framework LLM attributes to canonical `gen_ai.*` semantic conventions.**

Every GenAI framework (LangChain, OpenInference, OpenLLMetry, LiteLLM, Traceloop) uses different attribute names for the same thing. This processor maps them all to the official [OTel GenAI semantic conventions](https://opentelemetry.io/docs/specs/semconv/gen-ai/) so you can build one Grafana dashboard and one set of alerts — regardless of which SDKs your teams use.

---

## The problem

Your traces look like this depending on the SDK:

```
Team A (OpenLLMetry):    llm.usage.prompt_tokens = 150
Team B (OpenInference):  openinference.llm.token_count.prompt = 150
Team C (LangChain):      langchain.tokens.prompt = 150
Team D (LiteLLM):        litellm.usage.prompt_tokens = 150
```

After this processor, they all become:

```
gen_ai.usage.input_tokens = 150
```

One key. One dashboard. One alert.

---

## Quick start

Add `genai_semantic_normalizer` to your collector pipeline:

```yaml
processors:
  genai_semantic_normalizer:
    enable_defaults: true    # load 90+ built-in mappings
    overwrite: false         # don't overwrite existing canonical keys
    drop_original: false     # keep source attrs alongside canonical
    mappings:
      # add your org-specific keys here
      my.custom.model: gen_ai.request.model

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [genai_semantic_normalizer]
      exporters: [otlp, debug]
```

### Build with OCB (recommended)

```bash
# Install the OTel Collector Builder
go install go.opentelemetry.io/collector/cmd/builder@latest

# Build a custom collector with the normalizer
ocb --config dist/ocb-config.yaml
./dist/otelcol-genai --config examples/collector-config.yaml
```

### Build from source

```bash
go build -o otelcol-genai ./cmd/...
./otelcol-genai --config examples/collector-config.yaml
```

---

## Configuration

| Field | Type | Default | Description |
|---|---|---|---|
| `enable_defaults` | bool | `true` | Load 90+ built-in mappings for LangChain, OpenInference, OpenLLMetry, LiteLLM, Traceloop, and generic keys |
| `overwrite` | bool | `false` | If canonical key already exists, overwrite it with the mapped value |
| `drop_original` | bool | `false` | Remove the source attribute after mapping |
| `mappings` | map | `{}` | Custom source→canonical mappings (merged on top of defaults) |

---

## What gets mapped

The processor normalizes attributes across these categories:

| Category | Canonical Key | Example Sources |
|---|---|---|
| Model | `gen_ai.request.model` | `llm.model`, `openinference.model_name`, `langchain.llm.model_name` |
| Provider | `gen_ai.system` | `llm.vendor`, `llm.provider`, `traceloop.entity.provider` |
| Input tokens | `gen_ai.usage.input_tokens` | `llm.usage.prompt_tokens`, `prompt_tokens`, `input_tokens` |
| Output tokens | `gen_ai.usage.output_tokens` | `llm.usage.completion_tokens`, `completion_tokens` |
| Temperature | `gen_ai.request.temperature` | `llm.temperature`, `openinference.llm.temperature` |
| Max tokens | `gen_ai.request.max_tokens` | `llm.max_tokens`, `max_tokens` |
| Finish reason | `gen_ai.response.finish_reasons` | `llm.finish_reason`, `finish_reason` |
| Cost | `gen_ai.usage.cost` | `llm.usage.cost_usd`, `litellm.cost` |
| Prompt text | `gen_ai.prompt` | `llm.prompt`, `openinference.input.value` |
| Completion text | `gen_ai.completion` | `llm.completion`, `openinference.output.value` |

Full list: [docs/compatibility-matrix.md](docs/compatibility-matrix.md)

---

## Before / After

**Before** (raw span from OpenLLMetry + LiteLLM):
```
llm.model: "gpt-4o"
llm.vendor: "openai"
llm.usage.prompt_tokens: 200
llm.usage.completion_tokens: 50
llm.temperature: 0.8
llm.response.finish_reason: "stop"
litellm.cost: 0.0035
```

**After** (canonical keys added by normalizer):
```
gen_ai.request.model: "gpt-4o"
gen_ai.system: "openai"
gen_ai.usage.input_tokens: 200
gen_ai.usage.output_tokens: 50
gen_ai.request.temperature: 0.8
gen_ai.response.finish_reasons: "stop"
gen_ai.usage.cost: 0.0035

# originals preserved (drop_original=false)
llm.model: "gpt-4o"
llm.vendor: "openai"
...
```

---

## Grafana dashboard

Import the included dashboard from [grafana/genai-overview.dashboard.json](grafana/genai-overview.dashboard.json).

Panels include: LLM calls over time, token usage by model, models/providers in use, average latency, cost over time, and finish reasons.

Works with Tempo as a datasource (TraceQL queries on normalized `gen_ai.*` attributes).

---

## How it works

```
Your Apps (mixed SDKs)
    │
    │  llm.model, openinference.model_name,
    │  langchain.llm.model_name, etc.
    ▼
┌────────────────────────────────────────┐
│  OTel Collector                        │
│  ┌──────────────────────────────────┐  │
│  │  genai_semantic_normalizer       │  │
│  │                                  │  │
│  │  90+ built-in mappings           │  │
│  │  + your custom mappings          │  │
│  │  → gen_ai.request.model          │  │
│  │  → gen_ai.usage.input_tokens     │  │
│  │  → gen_ai.system                 │  │
│  │  → ...                           │  │
│  └──────────────────────────────────┘  │
└────────────────────────────────────────┘
    │
    │  Canonical gen_ai.* attributes
    ▼
One Grafana Dashboard / One Alert Set
```

Mapping applies to both **span attributes** and **span event attributes** (prompt/completion events).

---

## Pair with genaisafe processor

This processor pairs well with [otelcol-genai-safe](https://github.com/nostalgicskinco/opentelemetry-collector-processor-genai) for a complete GenAI observability pipeline:

```yaml
processors:
  genai_semantic_normalizer:    # normalize keys first
    enable_defaults: true
  genaisafe:                     # then redact/detect
    redact:
      mode: hash_and_preview
      keys: ["gen_ai.prompt", "gen_ai.completion"]

service:
  pipelines:
    traces:
      processors: [genai_semantic_normalizer, genaisafe]
```

---

## Contributing

PRs welcome. To add mappings for a new framework, update `defaults.go` and `docs/compatibility-matrix.md`.

```bash
go test -v ./...
go vet ./...
```

---

## License

Apache 2.0
