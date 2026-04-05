package agent

import (
	"testing"
)

func TestNewRegistry_Empty(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}

	all := r.All()
	if len(all) != 0 {
		t.Errorf("expected empty registry, got %d agents", len(all))
	}
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()

	kg := NewKnowledgeGuide(nil)
	r.Register(kg)

	got, err := r.Get("knowledge_guide")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got.Name() != "knowledge_guide" {
		t.Errorf("expected agent name 'knowledge_guide', got %q", got.Name())
	}
}

func TestRegistry_GetUnregistered(t *testing.T) {
	r := NewRegistry()

	_, err := r.Get("nonexistent")
	if err == nil {
		t.Error("expected error for unregistered agent, got nil")
	}
}

func TestRegistry_ListAll(t *testing.T) {
	r := NewRegistry()

	agents := []Agent{
		NewKnowledgeGuide(nil),
		NewPlanner(nil),
		NewDeveloper(nil),
		NewExecutor(nil),
		NewTrustManager(nil),
		NewAnalyst(nil),
	}

	for _, a := range agents {
		r.Register(a)
	}

	all := r.All()
	if len(all) != 6 {
		t.Errorf("expected 6 agents, got %d", len(all))
	}

	expectedNames := []string{
		"knowledge_guide", "planner", "developer",
		"executor", "trust_manager", "analyst",
	}
	for _, name := range expectedNames {
		if _, ok := all[name]; !ok {
			t.Errorf("expected agent %q in All() result", name)
		}
	}
}

func TestRegistry_RegisterOverwrite(t *testing.T) {
	r := NewRegistry()

	// Register an agent
	kg1 := NewKnowledgeGuide(nil)
	r.Register(kg1)

	// Register another agent with the same name (overwrites)
	kg2 := NewKnowledgeGuide(nil)
	r.Register(kg2)

	all := r.All()
	if len(all) != 1 {
		t.Errorf("expected 1 agent after overwrite, got %d", len(all))
	}
}

func TestRegistry_AllReturnsCopy(t *testing.T) {
	r := NewRegistry()
	r.Register(NewKnowledgeGuide(nil))

	all := r.All()
	// Modify the returned map
	all["injected"] = NewPlanner(nil)

	// Original registry should be unaffected
	_, err := r.Get("injected")
	if err == nil {
		t.Error("expected Get('injected') to fail — All() should return a copy")
	}
}
