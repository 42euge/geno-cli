# geno-cli

Agentic coding assistant TUI powered by Gemma 4 via Ollama. A local coding agent experience — private, free, no API keys.

## Install via geno-tools

From within an agent session:

```
/geno-tools install geno-cli
```

Or from the command line:

```bash
geno-tools install geno-cli
```

## Install standalone

```bash
git clone https://github.com/42euge/geno-cli.git
cd geno-cli
make install
```

Or with Go:

```bash
go install github.com/42euge/geno-cli/cmd/geno-cli@latest
```

## Prerequisites

- [Go](https://go.dev/) 1.23+
- [Ollama](https://ollama.com/) with a Gemma 4 model pulled

```bash
ollama pull gemma4:4b
```

## Usage

```bash
# Start with default model (gemma4:4b)
geno-cli

# Use a specific model
geno-cli --model gemma4:26b

# Connect to a remote Ollama instance
geno-cli --url http://192.168.1.100:11434
```

## Features

- Streaming chat with markdown rendering
- Agentic tool use: file read/write, shell commands, grep
- Vi-key scrollable viewport
- Syntax-highlighted code blocks
- Runs entirely locally on consumer hardware

## License

MIT
