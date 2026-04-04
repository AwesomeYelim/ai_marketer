package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"ai_marketer/internal/llm"
	"ai_marketer/internal/message"
)

// KnowledgeGuideOutput is the structured output from the Knowledge Guide agent.
type KnowledgeGuideOutput struct {
	Keywords         []KeywordAnalysis `json:"keywords"`
	TargetAudience   TargetAudience    `json:"target_audience"`
	ContentStrategy  ContentStrategy   `json:"content_strategy"`
	CompetitorNotes  string            `json:"competitor_notes"`
	BrandPositioning BrandPositioning  `json:"brand_positioning"`
}

// BrandPositioning defines the brand's strategic positioning.
type BrandPositioning struct {
	PositioningStatement string            `json:"positioning_statement" jsonschema:"description=We are X for Y who Z format"`
	USP                  string            `json:"usp" jsonschema:"description=Unique value proposition"`
	Personality          []string          `json:"personality" jsonschema:"description=Brand personality keywords"`
	PrimaryMessage       string            `json:"primary_message" jsonschema:"description=Core brand promise"`
	Differentiators      []string          `json:"differentiators" jsonschema:"description=Top 3 differentiators"`
	CompetitorPositions  []CompetitorBrand `json:"competitor_positions"`
}

// CompetitorBrand represents a competitor's brand positioning.
type CompetitorBrand struct {
	Name        string `json:"name"`
	Positioning string `json:"positioning"`
	Strength    string `json:"strength"`
	Gap         string `json:"gap" jsonschema:"description=Gap we can exploit"`
}

type KeywordAnalysis struct {
	Term         string `json:"term"`
	SearchVolume string `json:"search_volume" jsonschema:"description=Estimated search volume (high/medium/low)"`
	Difficulty   string `json:"difficulty" jsonschema:"description=Competition difficulty (high/medium/low)"`
	Intent       string `json:"intent" jsonschema:"description=Search intent (informational/transactional/navigational)"`
}

type TargetAudience struct {
	Description string   `json:"description"`
	PainPoints  []string `json:"pain_points"`
	SearchHabits string  `json:"search_habits"`
}

type ContentStrategy struct {
	Approach     string   `json:"approach"`
	ContentTypes []string `json:"content_types"`
	Tone         string   `json:"tone"`
	KeyMessages  []string `json:"key_messages"`
}

// KnowledgeGuide analyzes SEO strategy, keywords, and competitors.
type KnowledgeGuide struct {
	BaseAgent
}

// NewKnowledgeGuide creates a new Knowledge Guide agent.
func NewKnowledgeGuide(client *llm.Client) *KnowledgeGuide {
	return &KnowledgeGuide{
		BaseAgent: BaseAgent{
			AgentName:        "knowledge_guide",
			AgentDescription: "SEO 전략, 키워드, 경쟁사 분석을 담당하는 지식 가이더",
			AgentTaskTypes:   []message.TaskType{message.TaskKeywordResearch, message.TaskContentStrategy},
			LLM:              client,
			PromptFile:       "knowledge_guide.txt",
		},
	}
}

func (k *KnowledgeGuide) Process(ctx context.Context, msg *message.Message, wctx *message.WorkflowContext) (*message.AgentResponse, error) {
	systemPrompt, err := k.LoadSystemPrompt("prompts")
	if err != nil {
		return nil, fmt.Errorf("failed to load system prompt: %w", err)
	}

	userPrompt := fmt.Sprintf("사용자 요청: %s\n\n위 요청에 대해 SEO 키워드 분석, 타겟 오디언스 분석, 콘텐츠 전략을 수립해주세요.", wctx.UserRequest)

	result, err := llm.CompleteStructured[KnowledgeGuideOutput](k.LLM, ctx, systemPrompt, userPrompt)
	if err != nil {
		return &message.AgentResponse{
			AgentName:    k.Name(),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	output, _ := json.MarshalIndent(result, "", "  ")

	return &message.AgentResponse{
		AgentName: k.Name(),
		Success:   true,
		Output:    result,
		RawText:   string(output),
	}, nil
}
