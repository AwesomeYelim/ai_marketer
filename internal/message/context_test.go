package message

import (
	"sync"
	"testing"
)

func TestNewWorkflowContext_Initializes(t *testing.T) {
	wctx := NewWorkflowContext("test request")

	if wctx.UserRequest != "test request" {
		t.Errorf("expected UserRequest 'test request', got %q", wctx.UserRequest)
	}
	if wctx.Results == nil {
		t.Error("expected Results map to be initialized")
	}
	if len(wctx.Results) != 0 {
		t.Errorf("expected empty Results, got %d entries", len(wctx.Results))
	}
	if wctx.SharedData == nil {
		t.Error("expected SharedData map to be initialized")
	}
	if len(wctx.SharedData) != 0 {
		t.Errorf("expected empty SharedData, got %d entries", len(wctx.SharedData))
	}
	if wctx.RetryCount != 0 {
		t.Errorf("expected RetryCount 0, got %d", wctx.RetryCount)
	}
	if wctx.MaxRetries != 2 {
		t.Errorf("expected MaxRetries 2, got %d", wctx.MaxRetries)
	}
}

func TestSetResult_GetResult(t *testing.T) {
	wctx := NewWorkflowContext("test")

	resp := &AgentResponse{
		AgentName: "knowledge_guide",
		Success:   true,
		RawText:   "result text",
	}

	wctx.SetResult("knowledge_guide", resp)

	got, ok := wctx.GetResult("knowledge_guide")
	if !ok {
		t.Fatal("expected GetResult to return true for existing agent")
	}
	if got.AgentName != "knowledge_guide" {
		t.Errorf("expected agent name 'knowledge_guide', got %q", got.AgentName)
	}
	if !got.Success {
		t.Error("expected Success to be true")
	}
	if got.RawText != "result text" {
		t.Errorf("expected RawText 'result text', got %q", got.RawText)
	}
}

func TestGetResult_NonExistent(t *testing.T) {
	wctx := NewWorkflowContext("test")

	_, ok := wctx.GetResult("nonexistent_agent")
	if ok {
		t.Error("expected GetResult to return false for non-existent agent")
	}
}

func TestSetGet_SharedData(t *testing.T) {
	wctx := NewWorkflowContext("test")

	wctx.Set("revision_feedback", "some feedback")

	val, ok := wctx.Get("revision_feedback")
	if !ok {
		t.Fatal("expected Get to return true for existing key")
	}
	if val != "some feedback" {
		t.Errorf("expected 'some feedback', got %v", val)
	}

	_, ok = wctx.Get("nonexistent_key")
	if ok {
		t.Error("expected Get to return false for non-existent key")
	}
}

func TestIncrementRetry_RespectsMaxRetries(t *testing.T) {
	wctx := NewWorkflowContext("test")
	wctx.MaxRetries = 2

	// First increment: RetryCount becomes 1, 1 <= 2 => true
	if !wctx.IncrementRetry() {
		t.Error("expected first IncrementRetry to return true")
	}
	if wctx.RetryCount != 1 {
		t.Errorf("expected RetryCount 1, got %d", wctx.RetryCount)
	}

	// Second increment: RetryCount becomes 2, 2 <= 2 => true
	if !wctx.IncrementRetry() {
		t.Error("expected second IncrementRetry to return true")
	}
	if wctx.RetryCount != 2 {
		t.Errorf("expected RetryCount 2, got %d", wctx.RetryCount)
	}

	// Third increment: RetryCount becomes 3, 3 <= 2 => false
	if wctx.IncrementRetry() {
		t.Error("expected third IncrementRetry to return false (exceeded max)")
	}
	if wctx.RetryCount != 3 {
		t.Errorf("expected RetryCount 3, got %d", wctx.RetryCount)
	}
}

func TestConcurrentAccess(t *testing.T) {
	wctx := NewWorkflowContext("concurrent test")

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent SetResult calls
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			resp := &AgentResponse{
				AgentName: "agent",
				Success:   true,
				RawText:   "output",
			}
			wctx.SetResult("agent", resp)
		}(i)
	}
	wg.Wait()

	// Verify the result was set (last write wins, but no race)
	got, ok := wctx.GetResult("agent")
	if !ok {
		t.Fatal("expected result to be set after concurrent writes")
	}
	if !got.Success {
		t.Error("expected Success to be true")
	}

	// Concurrent GetResult calls
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			wctx.GetResult("agent")
		}()
	}
	wg.Wait()

	// Concurrent Set/Get on SharedData
	wg.Add(numGoroutines * 2)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			wctx.Set("key", idx)
		}(i)
		go func() {
			defer wg.Done()
			wctx.Get("key")
		}()
	}
	wg.Wait()

	// Concurrent IncrementRetry
	wctx2 := NewWorkflowContext("retry test")
	wctx2.MaxRetries = numGoroutines + 10
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			wctx2.IncrementRetry()
		}()
	}
	wg.Wait()

	if wctx2.RetryCount != numGoroutines {
		t.Errorf("expected RetryCount %d after concurrent increments, got %d", numGoroutines, wctx2.RetryCount)
	}
}
