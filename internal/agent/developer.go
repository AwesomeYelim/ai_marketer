package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"ai_marketer/internal/llm"
	"ai_marketer/internal/message"
)

// DeveloperOutput is the structured output from the Developer agent.
type DeveloperOutput struct {
	MetaTags           MetaTagsOutput       `json:"meta_tags"`
	SchemaMarkup       SchemaMarkupOutput   `json:"schema_markup"`
	OrganizationSchema OrganizationSchema   `json:"organization_schema"`
	BrandMeta          BrandMetaTags        `json:"brand_meta"`
	Sitemap            SitemapOutput        `json:"sitemap"`
	TechnicalSEO       []TechnicalItem      `json:"technical_seo"`
}

// OrganizationSchema is JSON-LD for the brand/organization.
type OrganizationSchema struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Logo         string   `json:"logo" jsonschema:"description=Logo URL placeholder"`
	URL          string   `json:"url"`
	SocialLinks  []string `json:"social_links"`
	RawJSON      string   `json:"raw_json" jsonschema:"description=Complete Organization JSON-LD"`
}

// BrandMetaTags holds brand-specific meta tag values.
type BrandMetaTags struct {
	BrandName string `json:"brand_name"`
	Tagline   string `json:"tagline"`
	OGSiteName string `json:"og_site_name"`
	ThemeColor string `json:"theme_color" jsonschema:"description=Brand primary color hex"`
}

type MetaTagsOutput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Canonical   string `json:"canonical"`
	OGTitle     string `json:"og_title"`
	OGDesc      string `json:"og_description"`
	OGType      string `json:"og_type"`
	Robots      string `json:"robots"`
}

type SchemaMarkupOutput struct {
	Type    string `json:"type"`
	RawJSON string `json:"raw_json" jsonschema:"description=Complete JSON-LD structured data"`
}

type SitemapOutput struct {
	URLs []SitemapURL `json:"urls"`
}

type SitemapURL struct {
	Loc        string `json:"loc"`
	Priority   string `json:"priority"`
	ChangeFreq string `json:"changefreq"`
}

type TechnicalItem struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

// Developer handles technical SEO implementation.
type Developer struct {
	BaseAgent
}

// NewDeveloper creates a new Developer agent.
func NewDeveloper(client *llm.Client) *Developer {
	return &Developer{
		BaseAgent: BaseAgent{
			AgentName:        "developer",
			AgentDescription: "기술적 SEO 구현을 담당하는 개발자",
			AgentTaskTypes:   []message.TaskType{message.TaskTechnicalSEO},
			LLM:              client,
			PromptFile:       "developer.txt",
		},
	}
}

func (d *Developer) Process(ctx context.Context, msg *message.Message, wctx *message.WorkflowContext) (*message.AgentResponse, error) {
	systemPrompt, err := d.LoadSystemPrompt("prompts")
	if err != nil {
		return nil, fmt.Errorf("failed to load system prompt: %w", err)
	}

	var prevContext string
	if planResult, ok := wctx.GetResult("planner"); ok && planResult.Success {
		prevContext = fmt.Sprintf("시행자 계획:\n%s", planResult.RawText)
	}

	userPrompt := fmt.Sprintf("사용자 요청: %s\n\n%s\n\n위 계획을 기반으로 메타태그, 스키마마크업, 사이트맵 등 기술적 SEO 요소를 구현해주세요.", wctx.UserRequest, prevContext)

	result, err := llm.CompleteStructured[DeveloperOutput](d.LLM, ctx, systemPrompt, userPrompt)
	if err != nil {
		return &message.AgentResponse{
			AgentName:    d.Name(),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	output, _ := json.MarshalIndent(result, "", "  ")

	return &message.AgentResponse{
		AgentName: d.Name(),
		Success:   true,
		Output:    result,
		RawText:   string(output),
	}, nil
}
