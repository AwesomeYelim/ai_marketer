package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds application configuration.
type Config struct {
	LLM       LLMConfig       `yaml:"llm"`
	Pipeline  PipelineConfig  `yaml:"pipeline"`
	Prompts   PromptsConfig   `yaml:"prompts"`
}

// LLMConfig holds LLM-related configuration.
type LLMConfig struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

// PipelineConfig holds pipeline-related configuration.
type PipelineConfig struct {
	MaxRetries int `yaml:"max_retries"`
}

// PromptsConfig holds prompt directory configuration.
type PromptsConfig struct {
	Dir string `yaml:"dir"`
}

// Load reads configuration from a YAML file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig(), nil
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Override API key from environment variable if set
	if envKey := os.Getenv("ANTHROPIC_API_KEY"); envKey != "" {
		cfg.LLM.APIKey = envKey
	}

	return cfg, nil
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		LLM: LLMConfig{
			Model: "claude-sonnet-4-6",
		},
		Pipeline: PipelineConfig{
			MaxRetries: 2,
		},
		Prompts: PromptsConfig{
			Dir: "prompts",
		},
	}
}
