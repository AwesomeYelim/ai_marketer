package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"ai_marketer/internal/llm"
	"ai_marketer/internal/message"
)

// ExecutorOutput is the structured output from the Executor agent.
type ExecutorOutput struct {
	BlogPost          BlogPostOutput   `json:"blog_post"`
	SocialCopy        []SocialCopyItem `json:"social_copy"`
	Summary           string           `json:"summary"`
	BrandVoiceApplied bool             `json:"brand_voice_applied" jsonschema:"description=Whether brand voice guidelines were applied"`
	USPMentionCount   int              `json:"usp_mention_count" jsonschema:"description=Number of times USP appears in content"`
	BrandHashtags     []string         `json:"brand_hashtags" jsonschema:"description=Brand-specific hashtags used"`
	BrandStorytelling string           `json:"brand_storytelling" jsonschema:"description=Brand story element woven into content"`
}

type BlogPostOutput struct {
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	Content  string `json:"content" jsonschema:"description=Full blog post content in markdown format"`
	WordCount int   `json:"word_count"`
	Excerpt  string `json:"excerpt"`
}

type SocialCopyItem struct {
	Platform string `json:"platform" jsonschema:"description=Target platform (twitter/linkedin/instagram)"`
	Content  string `json:"content"`
	Hashtags string `json:"hashtags"`
}

// Executor generates content (blog posts, social media copy).
type Executor struct {
	BaseAgent
}

// NewExecutor creates a new Executor agent.
func NewExecutor(client *llm.Client) *Executor {
	return &Executor{
		BaseAgent: BaseAgent{
			AgentName:        "executor",
			AgentDescription: "콘텐츠 작성 및 생성을 담당하는 수행자",
			AgentTaskTypes:   []message.TaskType{message.TaskContentGeneration},
			LLM:              client,
			PromptFile:       "executor.txt",
		},
	}
}

func (e *Executor) Process(ctx context.Context, msg *message.Message, wctx *message.WorkflowContext) (*message.AgentResponse, error) {
	systemPrompt, err := e.LoadSystemPrompt("prompts")
	if err != nil {
		return nil, fmt.Errorf("failed to load system prompt: %w", err)
	}

	var prevContext string
	if planResult, ok := wctx.GetResult("planner"); ok && planResult.Success {
		prevContext = fmt.Sprintf("시행자 계획:\n%s", planResult.RawText)
	}
	if kgResult, ok := wctx.GetResult("knowledge_guide"); ok && kgResult.Success {
		prevContext += fmt.Sprintf("\n\n지식 가이더 분석:\n%s", kgResult.RawText)
	}

	userPrompt := fmt.Sprintf("사용자 요청: %s\n\n%s\n\n위 계획과 분석을 기반으로 블로그 글과 소셜미디어 카피를 작성해주세요.", wctx.UserRequest, prevContext)

	result, err := llm.CompleteStructured[ExecutorOutput](e.LLM, ctx, systemPrompt, userPrompt)
	if err != nil {
		return &message.AgentResponse{
			AgentName:    e.Name(),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	output, _ := json.MarshalIndent(result, "", "  ")

	return &message.AgentResponse{
		AgentName: e.Name(),
		Success:   true,
		Output:    result,
		RawText:   string(output),
	}, nil
}
