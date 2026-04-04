package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"ai_marketer/internal/llm"
	"ai_marketer/internal/message"
)

// AnalystOutput is the structured output from the Analyst agent.
type AnalystOutput struct {
	KPIs              []KPIItem        `json:"kpis"`
	BrandKPIs         BrandKPIs        `json:"brand_kpis"`
	CompetitorCompare []CompetitorNote `json:"competitor_compare"`
	Improvements      []Improvement    `json:"improvements"`
	OverallSummary    string           `json:"overall_summary"`
	ExpectedImpact    string           `json:"expected_impact"`
}

// BrandKPIs holds brand-specific performance indicators.
type BrandKPIs struct {
	AwarenessProjection     string `json:"awareness_projection" jsonschema:"description=Projected brand awareness improvement"`
	PositioningClarity      int    `json:"positioning_clarity" jsonschema:"description=Positioning clarity score 1-10"`
	VoiceConsistencyScore   int    `json:"voice_consistency_score" jsonschema:"description=Voice consistency across outputs 1-10"`
	DifferentiationLevel    string `json:"differentiation_level" jsonschema:"description=Differentiation vs competitors (high/medium/low)"`
	BrandCompetitorAnalysis string `json:"brand_competitor_analysis" jsonschema:"description=Brand positioning vs competitor brands"`
}

type KPIItem struct {
	Name        string `json:"name"`
	Current     string `json:"current"`
	Projected   string `json:"projected"`
	Suggestion  string `json:"suggestion"`
}

type CompetitorNote struct {
	Competitor  string `json:"competitor"`
	Strength    string `json:"strength"`
	Weakness    string `json:"weakness"`
	Opportunity string `json:"opportunity"`
}

type Improvement struct {
	Area        string `json:"area"`
	Current     string `json:"current_state"`
	Recommended string `json:"recommended"`
	Priority    string `json:"priority" jsonschema:"description=high/medium/low"`
	Impact      string `json:"impact"`
}

// Analyst handles performance prediction and optimization suggestions.
type Analyst struct {
	BaseAgent
}

// NewAnalyst creates a new Analyst agent.
func NewAnalyst(client *llm.Client) *Analyst {
	return &Analyst{
		BaseAgent: BaseAgent{
			AgentName:        "analyst",
			AgentDescription: "성과 예측 및 최적화 제안을 담당하는 분석가",
			AgentTaskTypes:   []message.TaskType{message.TaskPerformanceAnalysis},
			LLM:              client,
			PromptFile:       "analyst.txt",
		},
	}
}

func (a *Analyst) Process(ctx context.Context, msg *message.Message, wctx *message.WorkflowContext) (*message.AgentResponse, error) {
	systemPrompt, err := a.LoadSystemPrompt("prompts")
	if err != nil {
		return nil, fmt.Errorf("failed to load system prompt: %w", err)
	}

	// Gather all previous agent outputs
	var prevContext string
	for _, name := range []string{"knowledge_guide", "planner", "developer", "executor", "trust_manager"} {
		if result, ok := wctx.GetResult(name); ok && result.Success {
			prevContext += fmt.Sprintf("[%s 출력]:\n%s\n\n", name, result.RawText)
		}
	}

	userPrompt := fmt.Sprintf("사용자 요청: %s\n\n%s\n위 전체 파이프라인 결과를 분석하여 KPI 예측, 경쟁사 비교, 개선점을 제시해주세요.", wctx.UserRequest, prevContext)

	result, err := llm.CompleteStructured[AnalystOutput](a.LLM, ctx, systemPrompt, userPrompt)
	if err != nil {
		return &message.AgentResponse{
			AgentName:    a.Name(),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	output, _ := json.MarshalIndent(result, "", "  ")

	return &message.AgentResponse{
		AgentName: a.Name(),
		Success:   true,
		Output:    result,
		RawText:   string(output),
	}, nil
}
