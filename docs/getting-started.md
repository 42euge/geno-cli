# Getting Started

## Prerequisites

- [Go](https://go.dev/) 1.23+
- [Ollama](https://ollama.com/) installed and running
- A Gemma 4 model pulled:

```bash
ollama pull gemma4:4b
```

## Installation

### Via geno-tools

From within an agent session:

```
/geno-tools install geno-cli
```

Or from the command line:

```bash
geno-tools install geno-cli
```

### Standalone

```bash
git clone https://github.com/42euge/geno-cli.git
cd geno-cli
make install
```

Or install directly with Go:

```bash
go install github.com/42euge/geno-cli/cmd/geno-cli@latest
```

## First use

Start Ollama if it is not already running:

```bash
ollama serve
```

Launch geno-cli:

```bash
geno-cli
```

### Options

```bash
# Use a specific model
geno-cli --model gemma4:26b

# Connect to a remote Ollama instance
geno-cli --url http://192.168.1.100:11434
```

## TUI controls

- Type your message and press Enter to send
- Use vi keys (`j`/`k`) to scroll the viewport
- The assistant can read/write files, run shell commands, and grep — all locally
