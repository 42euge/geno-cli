# geno-cli

Agentic coding assistant TUI powered by Gemma 4 via Ollama. Runs entirely locally — private, free, no API keys.

## Overview

geno-cli is a terminal-based coding assistant that connects to locally-hosted Gemma 4 models through Ollama. It provides streaming chat with tool use (file read/write, shell commands, grep) and a rich TUI built on Go's Charmbracelet ecosystem.

## Key features

- **Streaming chat** with Glamour markdown rendering
- **Agentic tool use**: read/write files, execute shell commands, grep
- **Vi-key scrollable viewport** for navigating long outputs
- **Syntax-highlighted code blocks**
- **Fully local** — no API keys, no cloud dependencies

## Navigation

- [Getting Started](getting-started.md) — install prerequisites, build, and run
