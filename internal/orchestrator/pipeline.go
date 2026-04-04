package orchestrator

// PipelineStep defines a single step in an execution pipeline.
type PipelineStep struct {
	AgentName string   `json:"agent_name"`
	DependsOn []string `json:"depends_on,omitempty"`
}

// Pipeline defines a sequence of agent execution steps.
type Pipeline struct {
	Name  string         `json:"name"`
	Steps []PipelineStep `json:"steps"`
}

// FullCampaignPipeline is the default pipeline for a full SEO campaign.
// Steps 3 (developer) and 4 (executor) depend only on planner, so they run in parallel.
// trust_manager waits for both, then analyst runs last.
var FullCampaignPipeline = &Pipeline{
	Name: "full_campaign",
	Steps: []PipelineStep{
		{AgentName: "knowledge_guide", DependsOn: nil},
		{AgentName: "planner", DependsOn: []string{"knowledge_guide"}},
		{AgentName: "developer", DependsOn: []string{"planner"}},
		{AgentName: "executor", DependsOn: []string{"planner"}},
		{AgentName: "trust_manager", DependsOn: []string{"developer", "executor"}},
		{AgentName: "analyst", DependsOn: []string{"trust_manager"}},
	},
}

// resolveLayers groups pipeline steps into execution layers.
// Steps within the same layer can run in parallel.
func resolveLayers(p *Pipeline) [][]PipelineStep {
	completed := make(map[string]bool)
	remaining := make([]PipelineStep, len(p.Steps))
	copy(remaining, p.Steps)

	var layers [][]PipelineStep

	for len(remaining) > 0 {
		var layer []PipelineStep
		var next []PipelineStep

		for _, step := range remaining {
			ready := true
			for _, dep := range step.DependsOn {
				if !completed[dep] {
					ready = false
					break
				}
			}
			if ready {
				layer = append(layer, step)
			} else {
				next = append(next, step)
			}
		}

		if len(layer) == 0 {
			// Circular dependency or unresolvable — break to avoid infinite loop
			break
		}

		for _, step := range layer {
			completed[step.AgentName] = true
		}

		layers = append(layers, layer)
		remaining = next
	}

	return layers
}
