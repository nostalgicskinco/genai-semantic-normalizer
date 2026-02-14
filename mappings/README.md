# Vendor Mapping Packs

Per-vendor/framework YAML files documenting the exact attribute keys each GenAI framework emits and their canonical `gen_ai.*` mappings.

## Available packs

| File | Vendor | Framework | Key count |
|---|---|---|---|
| [openai.yaml](openai.yaml) | OpenAI | OpenLLMetry | 18 |
| [anthropic.yaml](anthropic.yaml) | Anthropic | OpenLLMetry | 20 |
| [langchain.yaml](langchain.yaml) | LangChain | LangChain | 17 |
| [llamaindex.yaml](llamaindex.yaml) | LlamaIndex | OpenInference | 14 |
| [vertexai.yaml](vertexai.yaml) | Google | Vertex AI | 19 |
| [bedrock.yaml](bedrock.yaml) | AWS | Bedrock | 20 |
| [traceloop.yaml](traceloop.yaml) | Traceloop | OpenLLMetry | 8 |
| [litellm.yaml](litellm.yaml) | LiteLLM | LiteLLM | 16 |

## How to use

These YAML files serve two purposes:

1. **Documentation** — understand what attributes each framework emits
2. **Custom config** — copy mappings into your collector config's `mappings:` field for framework-specific overrides

The processor's built-in defaults already cover all of these. The YAML packs are useful when you need to understand or customize specific vendor behavior.

## Contributing a new pack

1. Create a new YAML file following the format of existing packs
2. Include: vendor, framework, description, mappings, and known span names
3. Open a PR (see [CONTRIBUTING.md](../CONTRIBUTING.md))
