package message

import (
	"fmt"
	"time"
)

// TaskType defines the type of task an agent can handle.
type TaskType string

const (
	TaskKeywordResearch   TaskType = "keyword_research"
	TaskContentStrategy   TaskType = "content_strategy"
	TaskActionPlanning    TaskType = "action_planning"
	TaskTechnicalSEO      TaskType = "technical_seo"
	TaskContentGeneration TaskType = "content_generation"
	TaskTrustEvaluation   TaskType = "trust_evaluation"
	TaskPerformanceAnalysis TaskType = "performance_analysis"
	TaskFullCampaign      TaskType = "full_campaign"
)

// Message represents a message passed between the orchestrator and agents.
type Message struct {
	ID        string            `json:"id"`
	From      string            `json:"from"`
	To        string            `json:"to"`
	TaskType  TaskType          `json:"task_type"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// NewMessage creates a new message.
func NewMessage(from, to string, taskType TaskType, content string) *Message {
	return &Message{
		ID:        generateID(),
		From:      from,
		To:        to,
		TaskType:  taskType,
		Content:   content,
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}
}

// AgentResponse represents the output from an agent.
type AgentResponse struct {
	AgentName    string      `json:"agent_name"`
	Success      bool        `json:"success"`
	Output       interface{} `json:"output"`
	RawText      string      `json:"raw_text"`
	ErrorMessage string      `json:"error_message,omitempty"`
}

var idCounter int

func generateID() string {
	idCounter++
	return fmt.Sprintf("msg_%d_%d", time.Now().UnixNano(), idCounter)
}
