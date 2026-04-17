---
name: geno-cli
description: >-
  Agentic coding assistant TUI powered by Gemma 4 via Ollama.
  Claude Code-like experience running entirely locally.
  Use when user says /gt-geno-cli or asks about local agentic coding.
license: MIT
metadata:
  author: 42euge
  version: "0.1.0"
---

# geno-cli

A terminal-based agentic coding assistant powered by locally-hosted Gemma 4 models via Ollama. Provides a Claude Code-like experience with streaming chat, tool use (file read/write, shell exec, grep), and a rich TUI built on Go's Charmbracelet ecosystem.

## When to use

Activate this skill when the user asks to:
- Run geno-cli or launch a local coding assistant
- Use Gemma 4 for agentic coding tasks
- Work with a local LLM-powered coding TUI

## Quick start

```bash
# Ensure ollama is running with a Gemma 4 model
ollama pull gemma4:4b
ollama serve

# Run geno-cli
geno-cli --model gemma4:4b
```
