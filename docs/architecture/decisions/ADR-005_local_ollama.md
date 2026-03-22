# ADR-005 — Local Ollama for AI Assistance

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

The `a → s` AI summary feature generates a natural-language summary of work done on an issue by analysing the issue description and linked git commit messages. This requires a language model.

Options include cloud LLM APIs (OpenAI, Anthropic, Google Gemini) and local LLM runtimes (Ollama, LM Studio, llama.cpp).

## Decision

Use a locally running [Ollama](https://ollama.com) instance at `http://localhost:11434` with the default model `llama3`.

The integration is a single `POST /api/generate` call with `stream: false`.

## Rationale

**Privacy:**
Issue descriptions and commit messages may contain sensitive internal information (product plans, security fixes, customer data). Sending this data to a cloud API raises privacy and compliance concerns that vary by organisation. A local LLM processes everything on the user's machine.

**No API key required:**
Cloud LLM APIs require account registration, billing setup, and API key management. A personal CLI tool should not impose this friction. Ollama runs locally with no authentication.

**No ongoing cost:**
Cloud LLM APIs charge per token. A local model has no per-request cost after the initial model download.

**Model flexibility:**
Ollama supports dozens of models (Llama 3, Mistral, Gemma, Phi, etc.). The default (`llama3`) can be overridden by changing the constant in `ollama/client.go` (see TD-03 for future configurability).

**Simple API:**
Ollama's HTTP API is a single endpoint (`POST /api/generate`) with a minimal JSON body. No SDK or complex integration needed.

## Alternatives Considered

| Alternative | Reason not chosen |
|-------------|------------------|
| **OpenAI API** | Sends issue content to cloud; requires API key and billing |
| **Anthropic Claude API** | Same concerns as OpenAI |
| **llama.cpp direct** | More complex integration; Ollama abstracts this cleanly |
| **No AI feature** | The feature is useful and differentiates the tool; local AI makes it viable |

## Consequences

- AI features require the user to install Ollama and pull the `llama3` model
- If Ollama is not running, the feature fails with an error — no silent degradation
- Error messages should guide users to start Ollama (current limitation — see R-03)
- The Ollama model name is hardcoded as `"llama3"` — see TD-03
- AI generation is slow on CPU; a spinner is shown during generation
- The summary quality depends entirely on the local model's capabilities
