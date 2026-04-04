package agent

import (
	"context"
	"os"
	"path/filepath"

	"ai_marketer/internal/llm"
	"ai_marketer/internal/message"
)

// Agent defines the interface all agents must implement.
type Agent interface {
	Name() string
	Description() string
	TaskTypes() []message.TaskType
	Process(ctx context.Context, msg *message.Message, wctx *message.WorkflowContext) (*message.AgentResponse, error)
}

// BaseAgent provides common functionality for all agents.
type BaseAgent struct {
	AgentName        string
	AgentDescription string
	AgentTaskTypes   []message.TaskType
	LLM              *llm.Client
	PromptFile       string
}

// Name returns the agent's name.
func (b *BaseAgent) Name() string { return b.AgentName }

// Description returns the agent's description.
func (b *BaseAgent) Description() string { return b.AgentDescription }

// TaskTypes returns the task types this agent handles.
func (b *BaseAgent) TaskTypes() []message.TaskType { return b.AgentTaskTypes }

// LoadSystemPrompt reads the agent's system prompt from the prompts directory.
func (b *BaseAgent) LoadSystemPrompt(promptsDir string) (string, error) {
	path := filepath.Join(promptsDir, b.PromptFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
