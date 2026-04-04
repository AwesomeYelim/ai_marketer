package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"ai_marketer/internal/llm"
	"ai_marketer/internal/message"
)

// TrustManagerOutput is the structured output from the Trust Manager agent.
type TrustManagerOutput struct {
	EEATScore              EEATScore              `json:"eeat_score"`
	BrandAuthenticity      BrandAuthenticityScore `json:"brand_authenticity"`
	BrandConsistencyChecks []BrandConsistencyCheck `json:"brand_consistency_checks"`
	FactChecks             []FactCheck            `json:"fact_checks"`
	Revisions              []RevisionItem         `json:"revisions"`
	NeedsRevision          bool                   `json:"needs_revision"`
	OverallVerdict         string                 `json:"overall_verdict"`
}

// BrandAuthenticityScore evaluates brand authenticity across dimensions.
type BrandAuthenticityScore struct {
	VoiceConsistency     int `json:"voice_consistency" jsonschema:"description=Voice consistency score 1-10"`
	USPClarity           int `json:"usp_clarity" jsonschema:"description=USP clarity score 1-10"`
	PositioningAlignment int `json:"positioning_alignment" jsonschema:"description=Positioning alignment score 1-10"`
	Overall              int `json:"overall" jsonschema:"description=Average brand authenticity score 1-10"`
}

// BrandConsistencyCheck is a single brand consistency verification item.
type BrandConsistencyCheck struct {
	Aspect string `json:"aspect" jsonschema:"description=What was checked (voice_tone/message_hierarchy/competitor_terms)"`
	Status string `json:"status" jsonschema:"description=pass/fail/warning"`
	Detail string `json:"detail"`
}

type EEATScore struct {
	Experience int    `json:"experience" jsonschema:"description=Experience score 1-10"`
	Expertise  int    `json:"expertise" jsonschema:"description=Expertise score 1-10"`
	Authority  int    `json:"authority" jsonschema:"description=Authority score 1-10"`
	Trust      int    `json:"trust" jsonschema:"description=Trustworthiness score 1-10"`
	Overall    int    `json:"overall" jsonschema:"description=Overall E-E-A-T score 1-10"`
}

type FactCheck struct {
	Claim    string `json:"claim"`
	Status   string `json:"status" jsonschema:"description=verified/unverified/false"`
	Comment  string `json:"comment"`
}

type RevisionItem struct {
	Target      string `json:"target" jsonschema:"description=Which agent output needs revision (developer/executor)"`
	Section     string `json:"section"`
	Issue       string `json:"issue"`
	Suggestion  string `json:"suggestion"`
}

// TrustManager handles E-E-A-T evaluation and quality control.
type TrustManager struct {
	BaseAgent
}

// NewTrustManager creates a new Trust Manager agent.
func NewTrustManager(client *llm.Client) *TrustManager {
	return &TrustManager{
		BaseAgent: BaseAgent{
			AgentName:        "trust_manager",
			AgentDescription: "E-E-A-T 검증 및 품질 관리를 담당하는 신뢰도 담당",
			AgentTaskTypes:   []message.TaskType{message.TaskTrustEvaluation},
			LLM:              client,
			PromptFile:       "trust_manager.txt",
		},
	}
}

func (t *TrustManager) Process(ctx context.Context, msg *message.Message, wctx *message.WorkflowContext) (*message.AgentResponse, error) {
	systemPrompt, err := t.LoadSystemPrompt("prompts")
	if err != nil {
		return nil, fmt.Errorf("failed to load system prompt: %w", err)
	}

	// Gather outputs from developer and executor for review
	var prevContext string
	if devResult, ok := wctx.GetResult("developer"); ok && devResult.Success {
		prevContext += fmt.Sprintf("개발자 출력:\n%s\n\n", devResult.RawText)
	}
	if execResult, ok := wctx.GetResult("executor"); ok && execResult.Success {
		prevContext += fmt.Sprintf("수행자 출력:\n%s\n\n", execResult.RawText)
	}

	userPrompt := fmt.Sprintf("사용자 요청: %s\n\n%s\n위 출력물에 대해 E-E-A-T 기준으로 품질을 평가하고, 팩트체크를 수행하고, 수정이 필요한 부분을 지적해주세요.", wctx.UserRequest, prevContext)

	result, err := llm.CompleteStructured[TrustManagerOutput](t.LLM, ctx, systemPrompt, userPrompt)
	if err != nil {
		return &message.AgentResponse{
			AgentName:    t.Name(),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	output, _ := json.MarshalIndent(result, "", "  ")

	return &message.AgentResponse{
		AgentName: t.Name(),
		Success:   true,
		Output:    result,
		RawText:   string(output),
	}, nil
}
