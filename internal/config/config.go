package config

import "flag"

type Config struct {
	OllamaURL string
	Model     string
	NoTools   bool
}

func Parse() Config {
	cfg := Config{}
	flag.StringVar(&cfg.OllamaURL, "url", "http://localhost:11434", "Ollama API URL")
	flag.StringVar(&cfg.Model, "model", "gemma4:e4b", "Model name")
	flag.BoolVar(&cfg.NoTools, "no-tools", false, "Disable tool use (plain chat mode)")
	flag.Parse()
	return cfg
}
