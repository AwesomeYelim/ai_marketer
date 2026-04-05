package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig_ReturnsValidValues(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	if cfg.LLM.Model != "claude-sonnet-4-6" {
		t.Errorf("expected model 'claude-sonnet-4-6', got %q", cfg.LLM.Model)
	}
	if cfg.Pipeline.MaxRetries != 2 {
		t.Errorf("expected max_retries 2, got %d", cfg.Pipeline.MaxRetries)
	}
	if cfg.Prompts.Dir != "prompts" {
		t.Errorf("expected prompts dir 'prompts', got %q", cfg.Prompts.Dir)
	}
	if cfg.LLM.APIKey != "" {
		t.Errorf("expected empty API key in default config, got %q", cfg.LLM.APIKey)
	}
}

func TestLoad_NonExistentPath_ReturnsDefault(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}
	// Should return default values
	if cfg.LLM.Model != "claude-sonnet-4-6" {
		t.Errorf("expected default model, got %q", cfg.LLM.Model)
	}
	if cfg.Pipeline.MaxRetries != 2 {
		t.Errorf("expected default max_retries 2, got %d", cfg.Pipeline.MaxRetries)
	}
}

func TestLoad_ValidYAML_ParsesCorrectly(t *testing.T) {
	dir := t.TempDir()
	yamlPath := filepath.Join(dir, "test_config.yaml")

	yamlContent := `llm:
  api_key: "test-key-123"
  model: "claude-opus-4-6"
pipeline:
  max_retries: 5
prompts:
  dir: "custom_prompts"
`
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp YAML: %v", err)
	}

	// Clear env var to avoid interference
	origKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	defer func() {
		if origKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", origKey)
		}
	}()

	cfg, err := Load(yamlPath)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.LLM.APIKey != "test-key-123" {
		t.Errorf("expected api_key 'test-key-123', got %q", cfg.LLM.APIKey)
	}
	if cfg.LLM.Model != "claude-opus-4-6" {
		t.Errorf("expected model 'claude-opus-4-6', got %q", cfg.LLM.Model)
	}
	if cfg.Pipeline.MaxRetries != 5 {
		t.Errorf("expected max_retries 5, got %d", cfg.Pipeline.MaxRetries)
	}
	if cfg.Prompts.Dir != "custom_prompts" {
		t.Errorf("expected prompts dir 'custom_prompts', got %q", cfg.Prompts.Dir)
	}
}

func TestLoad_EnvVarOverride(t *testing.T) {
	dir := t.TempDir()
	yamlPath := filepath.Join(dir, "test_config.yaml")

	yamlContent := `llm:
  api_key: "yaml-key"
  model: "claude-sonnet-4-6"
`
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp YAML: %v", err)
	}

	// Set env var — should override YAML value
	origKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Setenv("ANTHROPIC_API_KEY", "env-key-override")
	defer func() {
		if origKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", origKey)
		} else {
			os.Unsetenv("ANTHROPIC_API_KEY")
		}
	}()

	cfg, err := Load(yamlPath)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.LLM.APIKey != "env-key-override" {
		t.Errorf("expected env var override 'env-key-override', got %q", cfg.LLM.APIKey)
	}
}

func TestLoad_TempYAML_VerifyAllFields(t *testing.T) {
	dir := t.TempDir()
	yamlPath := filepath.Join(dir, "full_config.yaml")

	yamlContent := `llm:
  api_key: "full-test-key"
  model: "claude-haiku-4-6"
pipeline:
  max_retries: 10
prompts:
  dir: "/opt/custom/prompts"
`
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp YAML: %v", err)
	}

	origKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	defer func() {
		if origKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", origKey)
		}
	}()

	cfg, err := Load(yamlPath)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.LLM.APIKey != "full-test-key" {
		t.Errorf("api_key mismatch: got %q", cfg.LLM.APIKey)
	}
	if cfg.LLM.Model != "claude-haiku-4-6" {
		t.Errorf("model mismatch: got %q", cfg.LLM.Model)
	}
	if cfg.Pipeline.MaxRetries != 10 {
		t.Errorf("max_retries mismatch: got %d", cfg.Pipeline.MaxRetries)
	}
	if cfg.Prompts.Dir != "/opt/custom/prompts" {
		t.Errorf("prompts dir mismatch: got %q", cfg.Prompts.Dir)
	}
}
