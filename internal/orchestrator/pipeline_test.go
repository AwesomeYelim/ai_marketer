package orchestrator

import (
	"testing"
)

func TestNewPipeline(t *testing.T) {
	p := &Pipeline{
		Name:  "test_pipeline",
		Steps: []PipelineStep{},
	}

	if p.Name != "test_pipeline" {
		t.Errorf("expected name 'test_pipeline', got %q", p.Name)
	}
	if len(p.Steps) != 0 {
		t.Errorf("expected 0 steps, got %d", len(p.Steps))
	}
}

func TestPipeline_AddSteps(t *testing.T) {
	p := &Pipeline{
		Name: "test",
	}

	p.Steps = append(p.Steps, PipelineStep{
		AgentName: "agent_a",
		DependsOn: nil,
	})
	p.Steps = append(p.Steps, PipelineStep{
		AgentName: "agent_b",
		DependsOn: []string{"agent_a"},
	})

	if len(p.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(p.Steps))
	}
	if p.Steps[0].AgentName != "agent_a" {
		t.Errorf("expected first step 'agent_a', got %q", p.Steps[0].AgentName)
	}
	if p.Steps[1].AgentName != "agent_b" {
		t.Errorf("expected second step 'agent_b', got %q", p.Steps[1].AgentName)
	}
}

func TestBuildLayers_SingleStep(t *testing.T) {
	p := &Pipeline{
		Name: "single",
		Steps: []PipelineStep{
			{AgentName: "agent_a", DependsOn: nil},
		},
	}

	layers := resolveLayers(p)
	if len(layers) != 1 {
		t.Fatalf("expected 1 layer, got %d", len(layers))
	}
	if len(layers[0]) != 1 {
		t.Fatalf("expected 1 step in layer 0, got %d", len(layers[0]))
	}
	if layers[0][0].AgentName != "agent_a" {
		t.Errorf("expected 'agent_a', got %q", layers[0][0].AgentName)
	}
}

func TestBuildLayers_LinearChain(t *testing.T) {
	p := &Pipeline{
		Name: "linear",
		Steps: []PipelineStep{
			{AgentName: "a", DependsOn: nil},
			{AgentName: "b", DependsOn: []string{"a"}},
			{AgentName: "c", DependsOn: []string{"b"}},
		},
	}

	layers := resolveLayers(p)
	if len(layers) != 3 {
		t.Fatalf("expected 3 layers, got %d", len(layers))
	}
	if layers[0][0].AgentName != "a" {
		t.Errorf("layer 0 should contain 'a', got %q", layers[0][0].AgentName)
	}
	if layers[1][0].AgentName != "b" {
		t.Errorf("layer 1 should contain 'b', got %q", layers[1][0].AgentName)
	}
	if layers[2][0].AgentName != "c" {
		t.Errorf("layer 2 should contain 'c', got %q", layers[2][0].AgentName)
	}
}

func TestBuildLayers_ParallelSteps(t *testing.T) {
	p := &Pipeline{
		Name: "parallel",
		Steps: []PipelineStep{
			{AgentName: "root", DependsOn: nil},
			{AgentName: "branch_a", DependsOn: []string{"root"}},
			{AgentName: "branch_b", DependsOn: []string{"root"}},
			{AgentName: "merge", DependsOn: []string{"branch_a", "branch_b"}},
		},
	}

	layers := resolveLayers(p)
	if len(layers) != 3 {
		t.Fatalf("expected 3 layers, got %d", len(layers))
	}

	// Layer 0: root
	if len(layers[0]) != 1 || layers[0][0].AgentName != "root" {
		t.Errorf("layer 0: expected [root], got %v", layerNames(layers[0]))
	}

	// Layer 1: branch_a and branch_b (parallel)
	if len(layers[1]) != 2 {
		t.Fatalf("layer 1: expected 2 steps, got %d", len(layers[1]))
	}
	names := layerNames(layers[1])
	if !containsAll(names, "branch_a", "branch_b") {
		t.Errorf("layer 1: expected [branch_a, branch_b], got %v", names)
	}

	// Layer 2: merge
	if len(layers[2]) != 1 || layers[2][0].AgentName != "merge" {
		t.Errorf("layer 2: expected [merge], got %v", layerNames(layers[2]))
	}
}

func TestBuildLayers_FullCampaignPipeline(t *testing.T) {
	layers := resolveLayers(FullCampaignPipeline)

	// Expected structure:
	// Layer 0: knowledge_guide
	// Layer 1: planner
	// Layer 2: developer, executor (parallel)
	// Layer 3: trust_manager
	// Layer 4: analyst
	if len(layers) != 5 {
		t.Fatalf("expected 5 layers for full campaign, got %d", len(layers))
	}

	// Layer 0
	if len(layers[0]) != 1 || layers[0][0].AgentName != "knowledge_guide" {
		t.Errorf("layer 0: expected [knowledge_guide], got %v", layerNames(layers[0]))
	}

	// Layer 1
	if len(layers[1]) != 1 || layers[1][0].AgentName != "planner" {
		t.Errorf("layer 1: expected [planner], got %v", layerNames(layers[1]))
	}

	// Layer 2: parallel
	if len(layers[2]) != 2 {
		t.Fatalf("layer 2: expected 2 parallel steps, got %d", len(layers[2]))
	}
	l2names := layerNames(layers[2])
	if !containsAll(l2names, "developer", "executor") {
		t.Errorf("layer 2: expected [developer, executor], got %v", l2names)
	}

	// Layer 3
	if len(layers[3]) != 1 || layers[3][0].AgentName != "trust_manager" {
		t.Errorf("layer 3: expected [trust_manager], got %v", layerNames(layers[3]))
	}

	// Layer 4
	if len(layers[4]) != 1 || layers[4][0].AgentName != "analyst" {
		t.Errorf("layer 4: expected [analyst], got %v", layerNames(layers[4]))
	}
}

func TestBuildLayers_CircularDependency(t *testing.T) {
	p := &Pipeline{
		Name: "circular",
		Steps: []PipelineStep{
			{AgentName: "a", DependsOn: []string{"c"}},
			{AgentName: "b", DependsOn: []string{"a"}},
			{AgentName: "c", DependsOn: []string{"b"}},
		},
	}

	// Should not hang or panic — just break out
	layers := resolveLayers(p)

	// With all steps forming a cycle, no step can be resolved.
	// resolveLayers should return empty layers without infinite loop.
	if len(layers) != 0 {
		t.Logf("circular dependency produced %d layers (acceptable as long as no hang)", len(layers))
	}
}

func TestBuildLayers_PartialCircular(t *testing.T) {
	p := &Pipeline{
		Name: "partial_circular",
		Steps: []PipelineStep{
			{AgentName: "root", DependsOn: nil},
			{AgentName: "a", DependsOn: []string{"root"}},
			{AgentName: "b", DependsOn: []string{"c"}}, // circular with c
			{AgentName: "c", DependsOn: []string{"b"}}, // circular with b
		},
	}

	layers := resolveLayers(p)

	// root and a should be resolved, b and c form a cycle
	if len(layers) < 2 {
		t.Fatalf("expected at least 2 layers (root and a), got %d", len(layers))
	}

	// Verify root is in first layer
	if layers[0][0].AgentName != "root" {
		t.Errorf("expected root in layer 0, got %q", layers[0][0].AgentName)
	}

	// Verify 'a' is in second layer
	if layers[1][0].AgentName != "a" {
		t.Errorf("expected 'a' in layer 1, got %q", layers[1][0].AgentName)
	}

	// Total resolved steps should be 2 (root + a), not 4
	totalSteps := 0
	for _, layer := range layers {
		totalSteps += len(layer)
	}
	if totalSteps != 2 {
		t.Errorf("expected 2 resolved steps, got %d", totalSteps)
	}
}

func TestBuildLayers_EmptyPipeline(t *testing.T) {
	p := &Pipeline{
		Name:  "empty",
		Steps: []PipelineStep{},
	}

	layers := resolveLayers(p)
	if len(layers) != 0 {
		t.Errorf("expected 0 layers for empty pipeline, got %d", len(layers))
	}
}

func TestBuildLayers_AllParallel(t *testing.T) {
	p := &Pipeline{
		Name: "all_parallel",
		Steps: []PipelineStep{
			{AgentName: "a", DependsOn: nil},
			{AgentName: "b", DependsOn: nil},
			{AgentName: "c", DependsOn: nil},
		},
	}

	layers := resolveLayers(p)
	if len(layers) != 1 {
		t.Fatalf("expected 1 layer for all-parallel pipeline, got %d", len(layers))
	}
	if len(layers[0]) != 3 {
		t.Errorf("expected 3 steps in layer, got %d", len(layers[0]))
	}
}

// Helper: extract agent names from a layer.
func layerNames(layer []PipelineStep) []string {
	names := make([]string, len(layer))
	for i, step := range layer {
		names[i] = step.AgentName
	}
	return names
}

// Helper: check if all expected names are present.
func containsAll(names []string, expected ...string) bool {
	set := make(map[string]bool)
	for _, n := range names {
		set[n] = true
	}
	for _, e := range expected {
		if !set[e] {
			return false
		}
	}
	return true
}
