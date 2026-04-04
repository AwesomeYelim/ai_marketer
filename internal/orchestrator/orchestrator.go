package orchestrator

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"

	"ai_marketer/internal/agent"
	"ai_marketer/internal/message"
)

// Orchestrator coordinates agent execution through pipelines.
type Orchestrator struct {
	registry *agent.Registry
}

// New creates a new Orchestrator.
func New(registry *agent.Registry) *Orchestrator {
	return &Orchestrator{registry: registry}
}

// Run executes a pipeline with the given user request.
func (o *Orchestrator) Run(ctx context.Context, pipeline *Pipeline, userRequest string) (*message.WorkflowContext, error) {
	wctx := message.NewWorkflowContext(userRequest)
	layers := resolveLayers(pipeline)

	log.Printf("[orchestrator] 파이프라인 '%s' 시작 (총 %d 레이어)\n", pipeline.Name, len(layers))

	for i, layer := range layers {
		log.Printf("[orchestrator] 레이어 %d: %d개 에이전트 실행\n", i+1, len(layer))

		if err := o.executeLayer(ctx, layer, wctx); err != nil {
			return wctx, fmt.Errorf("layer %d failed: %w", i+1, err)
		}

		// Trust manager revision loop
		if o.layerContains(layer, "trust_manager") {
			if err := o.handleTrustRevision(ctx, pipeline, wctx); err != nil {
				return wctx, err
			}
		}
	}

	log.Printf("[orchestrator] 파이프라인 '%s' 완료\n", pipeline.Name)
	return wctx, nil
}

// executeLayer runs all agents in a layer concurrently.
func (o *Orchestrator) executeLayer(ctx context.Context, layer []PipelineStep, wctx *message.WorkflowContext) error {
	if len(layer) == 1 {
		return o.executeAgent(ctx, layer[0].AgentName, wctx)
	}

	g, gctx := errgroup.WithContext(ctx)
	for _, step := range layer {
		agentName := step.AgentName
		g.Go(func() error {
			return o.executeAgent(gctx, agentName, wctx)
		})
	}
	return g.Wait()
}

// executeAgent runs a single agent.
func (o *Orchestrator) executeAgent(ctx context.Context, agentName string, wctx *message.WorkflowContext) error {
	a, err := o.registry.Get(agentName)
	if err != nil {
		return err
	}

	log.Printf("[orchestrator] 에이전트 '%s' 실행 중...\n", agentName)

	msg := message.NewMessage("orchestrator", agentName, message.TaskFullCampaign, wctx.UserRequest)

	resp, err := a.Process(ctx, msg, wctx)
	if err != nil {
		return fmt.Errorf("agent '%s' error: %w", agentName, err)
	}

	wctx.SetResult(agentName, resp)

	if resp.Success {
		log.Printf("[orchestrator] 에이전트 '%s' 완료 ✓\n", agentName)
	} else {
		log.Printf("[orchestrator] 에이전트 '%s' 실패: %s\n", agentName, resp.ErrorMessage)
	}

	return nil
}

// handleTrustRevision checks if the trust manager requested revisions and re-runs developer/executor if needed.
func (o *Orchestrator) handleTrustRevision(ctx context.Context, pipeline *Pipeline, wctx *message.WorkflowContext) error {
	trustResult, ok := wctx.GetResult("trust_manager")
	if !ok || !trustResult.Success {
		return nil
	}

	output, ok := trustResult.Output.(*agent.TrustManagerOutput)
	if !ok {
		return nil
	}

	if !output.NeedsRevision {
		return nil
	}

	if !wctx.IncrementRetry() {
		log.Printf("[orchestrator] 최대 재시도 횟수 도달 (%d회), 수정 없이 계속 진행\n", wctx.MaxRetries)
		return nil
	}

	log.Printf("[orchestrator] 신뢰도 담당이 수정을 요청했습니다 (재시도 %d/%d)\n", wctx.RetryCount, wctx.MaxRetries)

	// Store revision feedback in shared data for agents to reference
	wctx.Set("revision_feedback", trustResult.RawText)

	// Re-run developer and executor in parallel
	revisionLayer := []PipelineStep{
		{AgentName: "developer"},
		{AgentName: "executor"},
	}

	if err := o.executeLayer(ctx, revisionLayer, wctx); err != nil {
		return fmt.Errorf("revision layer failed: %w", err)
	}

	// Re-run trust manager
	if err := o.executeAgent(ctx, "trust_manager", wctx); err != nil {
		return fmt.Errorf("trust manager re-evaluation failed: %w", err)
	}

	// Recursively check if more revisions are needed
	return o.handleTrustRevision(ctx, pipeline, wctx)
}

func (o *Orchestrator) layerContains(layer []PipelineStep, name string) bool {
	for _, step := range layer {
		if step.AgentName == name {
			return true
		}
	}
	return false
}
