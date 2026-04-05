package agent

import (
	"os"
	"path/filepath"
	"testing"
)

// agentTestCase holds test expectations for each agent constructor.
type agentTestCase struct {
	name        string
	constructor func() Agent
	expectedName string
}

func allAgentCases() []agentTestCase {
	return []agentTestCase{
		{
			name:         "KnowledgeGuide",
			constructor:  func() Agent { return NewKnowledgeGuide(nil) },
			expectedName: "knowledge_guide",
		},
		{
			name:         "Planner",
			constructor:  func() Agent { return NewPlanner(nil) },
			expectedName: "planner",
		},
		{
			name:         "Developer",
			constructor:  func() Agent { return NewDeveloper(nil) },
			expectedName: "developer",
		},
		{
			name:         "Executor",
			constructor:  func() Agent { return NewExecutor(nil) },
			expectedName: "executor",
		},
		{
			name:         "TrustManager",
			constructor:  func() Agent { return NewTrustManager(nil) },
			expectedName: "trust_manager",
		},
		{
			name:         "Analyst",
			constructor:  func() Agent { return NewAnalyst(nil) },
			expectedName: "analyst",
		},
	}
}

func TestAgent_Name(t *testing.T) {
	for _, tc := range allAgentCases() {
		t.Run(tc.name, func(t *testing.T) {
			a := tc.constructor()
			if a.Name() != tc.expectedName {
				t.Errorf("expected Name() = %q, got %q", tc.expectedName, a.Name())
			}
		})
	}
}

func TestAgent_Description_NonEmpty(t *testing.T) {
	for _, tc := range allAgentCases() {
		t.Run(tc.name, func(t *testing.T) {
			a := tc.constructor()
			if a.Description() == "" {
				t.Errorf("expected non-empty Description() for %s", tc.name)
			}
		})
	}
}

func TestAgent_TaskTypes_NonEmpty(t *testing.T) {
	for _, tc := range allAgentCases() {
		t.Run(tc.name, func(t *testing.T) {
			a := tc.constructor()
			if len(a.TaskTypes()) == 0 {
				t.Errorf("expected non-empty TaskTypes() for %s", tc.name)
			}
		})
	}
}

func TestLoadSystemPrompt_Success(t *testing.T) {
	// Create a temp directory with a prompt file
	dir := t.TempDir()
	promptContent := "You are a test agent. Respond in JSON."
	if err := os.WriteFile(filepath.Join(dir, "test_agent.txt"), []byte(promptContent), 0644); err != nil {
		t.Fatalf("failed to write prompt file: %v", err)
	}

	base := &BaseAgent{
		AgentName:  "test_agent",
		PromptFile: "test_agent.txt",
	}

	content, err := base.LoadSystemPrompt(dir)
	if err != nil {
		t.Fatalf("LoadSystemPrompt() returned error: %v", err)
	}
	if content != promptContent {
		t.Errorf("expected prompt content %q, got %q", promptContent, content)
	}
}

func TestLoadSystemPrompt_FromPromptsDir(t *testing.T) {
	// Test with the actual prompts directory if it exists
	promptsDir := filepath.Join("..", "..", "prompts")
	if _, err := os.Stat(promptsDir); os.IsNotExist(err) {
		t.Skip("prompts directory not found, skipping")
	}

	agents := allAgentCases()
	for _, tc := range agents {
		t.Run(tc.name, func(t *testing.T) {
			// Access the BaseAgent's PromptFile through the concrete type
			var base *BaseAgent
			switch a := tc.constructor().(type) {
			case *KnowledgeGuide:
				base = &a.BaseAgent
			case *Planner:
				base = &a.BaseAgent
			case *Developer:
				base = &a.BaseAgent
			case *Executor:
				base = &a.BaseAgent
			case *TrustManager:
				base = &a.BaseAgent
			case *Analyst:
				base = &a.BaseAgent
			default:
				t.Fatalf("unknown agent type for %s", tc.name)
			}

			content, err := base.LoadSystemPrompt(promptsDir)
			if err != nil {
				t.Fatalf("LoadSystemPrompt() returned error: %v", err)
			}
			if len(content) == 0 {
				t.Error("expected non-empty system prompt")
			}
		})
	}
}

func TestLoadSystemPrompt_NonExistentDir(t *testing.T) {
	base := &BaseAgent{
		AgentName:  "test_agent",
		PromptFile: "test_agent.txt",
	}

	_, err := base.LoadSystemPrompt("/nonexistent/directory")
	if err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}
}

func TestLoadSystemPrompt_NonExistentFile(t *testing.T) {
	dir := t.TempDir()

	base := &BaseAgent{
		AgentName:  "test_agent",
		PromptFile: "nonexistent.txt",
	}

	_, err := base.LoadSystemPrompt(dir)
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}
