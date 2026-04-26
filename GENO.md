# geno-cli — local agentic coding assistant TUI

Terminal-based coding assistant powered by Gemma 4 via Ollama. Provides streaming chat with tool use (file read/write, shell exec, grep) and a rich TUI built on Go's Charmbracelet ecosystem. Runs entirely locally — private, free, no API keys.

## Skills

| Skill | Sub-skillset | Slash command |
|-------|-------------|---------------|
| geno-cli | — | — (umbrella) |

## Repo structure

```
geno-cli/
├── GENO.md              # agent instructions (this file)
├── SKILL.md             # umbrella skill manifest
├── genotools.yaml       # geno-tools manifest
├── skills/              # skill definitions
│   └── geno-cli/        #   umbrella
├── docs/                # MkDocs Material site
├── internal/            # Go source
│   ├── agent/           #   LLM agent loop and system prompt
│   ├── app/             #   Bubble Tea TUI (model/update/view)
│   ├── chat/            #   chat history and message types
│   ├── config/          #   CLI flags and configuration
│   ├── ollama/          #   Ollama HTTP client and streaming
│   ├── render/          #   Glamour markdown rendering
│   └── tools/           #   tool registry (bash, grep, read, write, list)
├── go.mod               # Go module
├── Makefile             # build/install/run targets
└── LICENSE              # MIT
```

## Architecture

### Entry point

`cmd/geno-cli/main.go` (referenced by `go.mod` module `github.com/42euge/geno-cli`) — parses CLI flags via `internal/config`, initializes the Bubble Tea app.

### Key modules

- **`internal/agent`** — orchestrates the LLM agent loop: sends messages to Ollama, parses tool-call responses, executes tools, feeds results back. `system_prompt.go` defines the system prompt.
- **`internal/app`** — Bubble Tea model with `app.go` (state), `update.go` (message handling), `view.go` (rendering). Vi-key scrollable viewport.
- **`internal/ollama`** — HTTP client for the Ollama `/api/chat` endpoint with streaming support. Types in `types.go`.
- **`internal/tools`** — tool registry (`registry.go`) with implementations: `bash.go`, `grep.go`, `read_file.go`, `write_file.go`, `list_files.go`.
- **`internal/render`** — Glamour-based markdown rendering for chat output.

### Data flow

1. User types a message in the TUI
2. `app.Update` sends it to the agent loop
3. `agent.Loop` streams the response from Ollama
4. If the LLM emits a tool call, the agent executes it and feeds the result back
5. Final response is rendered via Glamour and displayed in the viewport

## Dependencies and runtime

- **Go** 1.23+ (build)
- **Ollama** running locally with a Gemma 4 model pulled (`ollama pull gemma4:4b`)
- No Python, no venv, no API keys

## Conventions

- Single umbrella skill — no sub-skills currently
- Go source lives under `internal/` following standard Go project layout
- Build and install via `Makefile` targets (`make build`, `make install`)
