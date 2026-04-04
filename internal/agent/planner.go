package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"ai_marketer/internal/llm"
	"ai_marketer/internal/message"
)

// PlannerOutput is the structured output from the Planner agent.
type PlannerOutput struct {
	ContentOutline      ContentOutlinePlan  `json:"content_outline"`
	ActionItems         []ActionItem        `json:"action_items"`
	Schedule            []ScheduleItem      `json:"schedule"`
	Priority            string              `json:"priority"`
	BrandVoiceGuideline BrandVoiceGuideline `json:"brand_voice_guideline"`
	SectionMessageMap   []SectionMessage    `json:"section_message_map" jsonschema:"description=Maps brand messages to content sections"`
}

// BrandVoiceGuideline defines tone, style, and vocabulary for brand consistency.
type BrandVoiceGuideline struct {
	Tone           string   `json:"tone" jsonschema:"description=Overall tone e.g. 전문적이면서 친근한"`
	StyleRules     []string `json:"style_rules" jsonschema:"description=Writing style rules e.g. 능동태 사용"`
	PreferredTerms []string `json:"preferred_terms" jsonschema:"description=Preferred vocabulary"`
	AvoidTerms     []string `json:"avoid_terms" jsonschema:"description=Terms to avoid"`
	PlatformVoice  []PlatformVoiceVariant `json:"platform_voice" jsonschema:"description=Voice variations by platform"`
}

// PlatformVoiceVariant defines voice adjustments per platform.
type PlatformVoiceVariant struct {
	Platform string `json:"platform" jsonschema:"description=blog/instagram/twitter/linkedin"`
	ToneAdj  string `json:"tone_adjustment"`
}

// SectionMessage maps a content section to a brand message.
type SectionMessage struct {
	Section string `json:"section"`
	Message string `json:"brand_message"`
}

type ContentOutlinePlan struct {
	Title       string          `json:"title"`
	Slug        string          `json:"slug"`
	TargetWords int             `json:"target_words"`
	Sections    []PlanSection   `json:"sections"`
}

type PlanSection struct {
	Heading   string   `json:"heading"`
	KeyPoints []string `json:"key_points"`
	Keywords  []string `json:"keywords"`
	WordCount int      `json:"word_count"`
}

type ActionItem struct {
	Task        string `json:"task"`
	AssignedTo  string `json:"assigned_to" jsonschema:"description=Agent name (developer/executor)"`
	Description string `json:"description"`
}

type ScheduleItem struct {
	Phase       string `json:"phase"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

// Planner converts strategy into concrete action plans.
type Planner struct {
	BaseAgent
}

// NewPlanner creates a new Planner agent.
func NewPlanner(client *llm.Client) *Planner {
	return &Planner{
		BaseAgent: BaseAgent{
			AgentName:        "planner",
			AgentDescription: "전략을 구체적 액션플랜으로 변환하는 시행자",
			AgentTaskTypes:   []message.TaskType{message.TaskActionPlanning},
			LLM:              client,
			PromptFile:       "planner.txt",
		},
	}
}

func (p *Planner) Process(ctx context.Context, msg *message.Message, wctx *message.WorkflowContext) (*message.AgentResponse, error) {
	systemPrompt, err := p.LoadSystemPrompt("prompts")
	if err != nil {
		return nil, fmt.Errorf("failed to load system prompt: %w", err)
	}

	// Build context from previous agent results
	var prevContext string
	if kgResult, ok := wctx.GetResult("knowledge_guide"); ok && kgResult.Success {
		prevContext = fmt.Sprintf("지식 가이더 분석 결과:\n%s", kgResult.RawText)
	}

	userPrompt := fmt.Sprintf("사용자 요청: %s\n\n%s\n\n위 분석 결과를 기반으로 구체적인 액션플랜, 콘텐츠 아웃라인, 작업 일정을 수립해주세요.", wctx.UserRequest, prevContext)

	result, err := llm.CompleteStructured[PlannerOutput](p.LLM, ctx, systemPrompt, userPrompt)
	if err != nil {
		return &message.AgentResponse{
			AgentName:    p.Name(),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	output, _ := json.MarshalIndent(result, "", "  ")

	return &message.AgentResponse{
		AgentName: p.Name(),
		Success:   true,
		Output:    result,
		RawText:   string(output),
	}, nil
}
