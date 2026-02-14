// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: Apache-2.0

package genainormprocessor

// DefaultMappings returns the built-in attribute mappings from common
// GenAI frameworks/vendors to the official OTel gen_ai.* semantic
// conventions. Sources: LangChain, OpenInference, OpenLLMetry,
// LiteLLM, Langfuse, Traceloop, Arize, custom SDKs.
func DefaultMappings() map[string]string {
	return map[string]string{
		// === Model ===
		"llm.model":                    "gen_ai.request.model",
		"llm.request.model":            "gen_ai.request.model",
		"openinference.model_name":     "gen_ai.request.model",
		"traceloop.entity.model":       "gen_ai.request.model",
		"langchain.llm.model_name":     "gen_ai.request.model",
		"litellm.model":                "gen_ai.request.model",
		"model_name":                   "gen_ai.request.model",
		"model":                        "gen_ai.request.model",

		// === Provider / System ===
		"llm.vendor":                   "gen_ai.system",
		"llm.provider":                 "gen_ai.system",
		"openinference.llm.provider":   "gen_ai.system",
		"traceloop.entity.provider":    "gen_ai.system",
		"langchain.llm.provider":       "gen_ai.system",
		"litellm.provider":             "gen_ai.system",

		// === Token usage: prompt ===
		"llm.usage.prompt_tokens":              "gen_ai.usage.input_tokens",
		"llm.token_count.prompt":               "gen_ai.usage.input_tokens",
		"openinference.llm.token_count.prompt": "gen_ai.usage.input_tokens",
		"traceloop.entity.input_tokens":        "gen_ai.usage.input_tokens",
		"langchain.tokens.prompt":              "gen_ai.usage.input_tokens",
		"litellm.usage.prompt_tokens":          "gen_ai.usage.input_tokens",
		"prompt_tokens":                        "gen_ai.usage.input_tokens",
		"input_tokens":                         "gen_ai.usage.input_tokens",

		// === Token usage: completion ===
		"llm.usage.completion_tokens":              "gen_ai.usage.output_tokens",
		"llm.token_count.completion":               "gen_ai.usage.output_tokens",
		"openinference.llm.token_count.completion": "gen_ai.usage.output_tokens",
		"traceloop.entity.output_tokens":           "gen_ai.usage.output_tokens",
		"langchain.tokens.completion":              "gen_ai.usage.output_tokens",
		"litellm.usage.completion_tokens":          "gen_ai.usage.output_tokens",
		"completion_tokens":                        "gen_ai.usage.output_tokens",
		"output_tokens":                            "gen_ai.usage.output_tokens",

		// === Token usage: total ===
		"llm.usage.total_tokens":       "gen_ai.usage.total_tokens",
		"total_tokens":                 "gen_ai.usage.total_tokens",

		// === Temperature ===
		"llm.temperature":                  "gen_ai.request.temperature",
		"llm.request.temperature":          "gen_ai.request.temperature",
		"openinference.llm.temperature":    "gen_ai.request.temperature",

		// === Max tokens ===
		"llm.max_tokens":                   "gen_ai.request.max_tokens",
		"llm.request.max_tokens":           "gen_ai.request.max_tokens",
		"openinference.llm.max_tokens":     "gen_ai.request.max_tokens",
		"max_tokens":                       "gen_ai.request.max_tokens",

		// === Top P ===
		"llm.top_p":                        "gen_ai.request.top_p",
		"llm.request.top_p":                "gen_ai.request.top_p",
		"top_p":                            "gen_ai.request.top_p",

		// === Stop sequences ===
		"llm.stop_sequences":               "gen_ai.request.stop_sequences",

		// === Frequency / Presence penalty ===
		"llm.frequency_penalty":            "gen_ai.request.frequency_penalty",
		"llm.presence_penalty":             "gen_ai.request.presence_penalty",

		// === Response: finish reason ===
		"llm.response.finish_reason":       "gen_ai.response.finish_reasons",
		"llm.finish_reason":                "gen_ai.response.finish_reasons",
		"finish_reason":                    "gen_ai.response.finish_reasons",

		// === Response: model (actual model used) ===
		"llm.response.model":               "gen_ai.response.model",

		// === Response: ID ===
		"llm.response.id":                  "gen_ai.response.id",

		// === Prompt / Completion content ===
		"llm.prompt":                       "gen_ai.prompt",
		"llm.completion":                   "gen_ai.completion",
		"openinference.input.value":        "gen_ai.prompt",
		"openinference.output.value":       "gen_ai.completion",
		"traceloop.entity.input":           "gen_ai.prompt",
		"traceloop.entity.output":          "gen_ai.completion",

		// === Operation name ===
		"llm.request.type":                 "gen_ai.operation.name",
		"traceloop.entity.type":            "gen_ai.operation.name",

		// === Cost ===
		"llm.usage.cost":                   "gen_ai.usage.cost",
		"llm.usage.cost_usd":               "gen_ai.usage.cost",
		"gen_ai.usage.cost_usd":            "gen_ai.usage.cost",
		"litellm.cost":                     "gen_ai.usage.cost",
	}
}
