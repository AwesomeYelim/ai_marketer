package message

import (
	"strings"
	"testing"
	"time"
)

func TestNewMessage_CreatesCorrectFields(t *testing.T) {
	msg := NewMessage("orchestrator", "knowledge_guide", TaskKeywordResearch, "test content")

	if msg.From != "orchestrator" {
		t.Errorf("expected From 'orchestrator', got %q", msg.From)
	}
	if msg.To != "knowledge_guide" {
		t.Errorf("expected To 'knowledge_guide', got %q", msg.To)
	}
	if msg.TaskType != TaskKeywordResearch {
		t.Errorf("expected TaskType %q, got %q", TaskKeywordResearch, msg.TaskType)
	}
	if msg.Content != "test content" {
		t.Errorf("expected Content 'test content', got %q", msg.Content)
	}
	if msg.ID == "" {
		t.Error("expected non-empty ID")
	}
	if !strings.HasPrefix(msg.ID, "msg_") {
		t.Errorf("expected ID to start with 'msg_', got %q", msg.ID)
	}
	if msg.Metadata == nil {
		t.Error("expected Metadata to be initialized (non-nil)")
	}
	if msg.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	// CreatedAt should be very recent
	if time.Since(msg.CreatedAt) > 5*time.Second {
		t.Error("CreatedAt is not recent")
	}
}

func TestNewMessage_UniqueIDs(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		msg := NewMessage("a", "b", TaskFullCampaign, "content")
		if ids[msg.ID] {
			t.Fatalf("duplicate ID generated: %q", msg.ID)
		}
		ids[msg.ID] = true
	}
}

func TestTaskType_Constants(t *testing.T) {
	tests := []struct {
		taskType TaskType
		expected string
	}{
		{TaskKeywordResearch, "keyword_research"},
		{TaskContentStrategy, "content_strategy"},
		{TaskActionPlanning, "action_planning"},
		{TaskTechnicalSEO, "technical_seo"},
		{TaskContentGeneration, "content_generation"},
		{TaskTrustEvaluation, "trust_evaluation"},
		{TaskPerformanceAnalysis, "performance_analysis"},
		{TaskFullCampaign, "full_campaign"},
	}

	for _, tt := range tests {
		if string(tt.taskType) != tt.expected {
			t.Errorf("TaskType %q != expected %q", tt.taskType, tt.expected)
		}
	}
}
