package llm

import (
	"testing"
)

// We re-define minimal structs matching the agent output types here
// because the llm package cannot import the agent package (circular dependency).
// These mirror the structures in internal/agent/*.go.

type testKnowledgeGuideOutput struct {
	Keywords         []testKeywordAnalysis `json:"keywords"`
	TargetAudience   testTargetAudience    `json:"target_audience"`
	ContentStrategy  testContentStrategy   `json:"content_strategy"`
	CompetitorNotes  string                `json:"competitor_notes"`
	BrandPositioning testBrandPositioning  `json:"brand_positioning"`
}

type testBrandPositioning struct {
	PositioningStatement string               `json:"positioning_statement"`
	USP                  string               `json:"usp"`
	Personality          []string             `json:"personality"`
	PrimaryMessage       string               `json:"primary_message"`
	Differentiators      []string             `json:"differentiators"`
	CompetitorPositions  []testCompetitorBrand `json:"competitor_positions"`
}

type testCompetitorBrand struct {
	Name        string `json:"name"`
	Positioning string `json:"positioning"`
	Strength    string `json:"strength"`
	Gap         string `json:"gap"`
}

type testKeywordAnalysis struct {
	Term         string `json:"term"`
	SearchVolume string `json:"search_volume"`
	Difficulty   string `json:"difficulty"`
	Intent       string `json:"intent"`
}

type testTargetAudience struct {
	Description  string   `json:"description"`
	PainPoints   []string `json:"pain_points"`
	SearchHabits string   `json:"search_habits"`
}

type testContentStrategy struct {
	Approach     string   `json:"approach"`
	ContentTypes []string `json:"content_types"`
	Tone         string   `json:"tone"`
	KeyMessages  []string `json:"key_messages"`
}

type testPlannerOutput struct {
	ContentOutline      testContentOutlinePlan  `json:"content_outline"`
	ActionItems         []testActionItem        `json:"action_items"`
	Schedule            []testScheduleItem      `json:"schedule"`
	Priority            string                  `json:"priority"`
	BrandVoiceGuideline testBrandVoiceGuideline `json:"brand_voice_guideline"`
	SectionMessageMap   []testSectionMessage    `json:"section_message_map"`
}

type testBrandVoiceGuideline struct {
	Tone           string                    `json:"tone"`
	StyleRules     []string                  `json:"style_rules"`
	PreferredTerms []string                  `json:"preferred_terms"`
	AvoidTerms     []string                  `json:"avoid_terms"`
	PlatformVoice  []testPlatformVoiceVariant `json:"platform_voice"`
}

type testPlatformVoiceVariant struct {
	Platform string `json:"platform"`
	ToneAdj  string `json:"tone_adjustment"`
}

type testSectionMessage struct {
	Section string `json:"section"`
	Message string `json:"brand_message"`
}

type testContentOutlinePlan struct {
	Title       string            `json:"title"`
	Slug        string            `json:"slug"`
	TargetWords int               `json:"target_words"`
	Sections    []testPlanSection `json:"sections"`
}

type testPlanSection struct {
	Heading   string   `json:"heading"`
	KeyPoints []string `json:"key_points"`
	Keywords  []string `json:"keywords"`
	WordCount int      `json:"word_count"`
}

type testActionItem struct {
	Task        string `json:"task"`
	AssignedTo  string `json:"assigned_to"`
	Description string `json:"description"`
}

type testScheduleItem struct {
	Phase       string `json:"phase"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

type testDeveloperOutput struct {
	MetaTags           testMetaTagsOutput     `json:"meta_tags"`
	SchemaMarkup       testSchemaMarkupOutput `json:"schema_markup"`
	OrganizationSchema testOrganizationSchema `json:"organization_schema"`
	BrandMeta          testBrandMetaTags      `json:"brand_meta"`
	Sitemap            testSitemapOutput      `json:"sitemap"`
	TechnicalSEO       []testTechnicalItem    `json:"technical_seo"`
}

type testOrganizationSchema struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	URL         string   `json:"url"`
	SocialLinks []string `json:"social_links"`
	RawJSON     string   `json:"raw_json"`
}

type testBrandMetaTags struct {
	BrandName  string `json:"brand_name"`
	Tagline    string `json:"tagline"`
	OGSiteName string `json:"og_site_name"`
	ThemeColor string `json:"theme_color"`
}

type testMetaTagsOutput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Canonical   string `json:"canonical"`
	OGTitle     string `json:"og_title"`
	OGDesc      string `json:"og_description"`
	OGType      string `json:"og_type"`
	Robots      string `json:"robots"`
}

type testSchemaMarkupOutput struct {
	Type    string `json:"type"`
	RawJSON string `json:"raw_json"`
}

type testSitemapOutput struct {
	URLs []testSitemapURL `json:"urls"`
}

type testSitemapURL struct {
	Loc        string `json:"loc"`
	Priority   string `json:"priority"`
	ChangeFreq string `json:"changefreq"`
}

type testTechnicalItem struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

type testExecutorOutput struct {
	BlogPost          testBlogPostOutput   `json:"blog_post"`
	SocialCopy        []testSocialCopyItem `json:"social_copy"`
	Summary           string               `json:"summary"`
	BrandVoiceApplied bool                 `json:"brand_voice_applied"`
	USPMentionCount   int                  `json:"usp_mention_count"`
	BrandHashtags     []string             `json:"brand_hashtags"`
	BrandStorytelling string               `json:"brand_storytelling"`
}

type testBlogPostOutput struct {
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Content   string `json:"content"`
	WordCount int    `json:"word_count"`
	Excerpt   string `json:"excerpt"`
}

type testSocialCopyItem struct {
	Platform string `json:"platform"`
	Content  string `json:"content"`
	Hashtags string `json:"hashtags"`
}

type testTrustManagerOutput struct {
	EEATScore              testEEATScore                `json:"eeat_score"`
	BrandAuthenticity      testBrandAuthenticityScore   `json:"brand_authenticity"`
	BrandConsistencyChecks []testBrandConsistencyCheck  `json:"brand_consistency_checks"`
	FactChecks             []testFactCheck              `json:"fact_checks"`
	Revisions              []testRevisionItem           `json:"revisions"`
	NeedsRevision          bool                         `json:"needs_revision"`
	OverallVerdict         string                       `json:"overall_verdict"`
}

type testBrandAuthenticityScore struct {
	VoiceConsistency     int `json:"voice_consistency"`
	USPClarity           int `json:"usp_clarity"`
	PositioningAlignment int `json:"positioning_alignment"`
	Overall              int `json:"overall"`
}

type testBrandConsistencyCheck struct {
	Aspect string `json:"aspect"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

type testEEATScore struct {
	Experience int `json:"experience"`
	Expertise  int `json:"expertise"`
	Authority  int `json:"authority"`
	Trust      int `json:"trust"`
	Overall    int `json:"overall"`
}

type testFactCheck struct {
	Claim   string `json:"claim"`
	Status  string `json:"status"`
	Comment string `json:"comment"`
}

type testRevisionItem struct {
	Target     string `json:"target"`
	Section    string `json:"section"`
	Issue      string `json:"issue"`
	Suggestion string `json:"suggestion"`
}

type testAnalystOutput struct {
	KPIs              []testKPIItem        `json:"kpis"`
	BrandKPIs         testBrandKPIs        `json:"brand_kpis"`
	CompetitorCompare []testCompetitorNote `json:"competitor_compare"`
	Improvements      []testImprovement    `json:"improvements"`
	OverallSummary    string               `json:"overall_summary"`
	ExpectedImpact    string               `json:"expected_impact"`
}

type testBrandKPIs struct {
	AwarenessProjection     string `json:"awareness_projection"`
	PositioningClarity      int    `json:"positioning_clarity"`
	VoiceConsistencyScore   int    `json:"voice_consistency_score"`
	DifferentiationLevel    string `json:"differentiation_level"`
	BrandCompetitorAnalysis string `json:"brand_competitor_analysis"`
}

type testKPIItem struct {
	Name       string `json:"name"`
	Current    string `json:"current"`
	Projected  string `json:"projected"`
	Suggestion string `json:"suggestion"`
}

type testCompetitorNote struct {
	Competitor  string `json:"competitor"`
	Strength    string `json:"strength"`
	Weakness    string `json:"weakness"`
	Opportunity string `json:"opportunity"`
}

type testImprovement struct {
	Area        string `json:"area"`
	Current     string `json:"current_state"`
	Recommended string `json:"recommended"`
	Priority    string `json:"priority"`
	Impact      string `json:"impact"`
}

func TestGenerateSchema_KnowledgeGuideOutput(t *testing.T) {
	schema := GenerateSchema[testKnowledgeGuideOutput]()
	assertSchemaValid(t, schema, "KnowledgeGuideOutput")

	// Check branding field
	assertPropertyExists(t, schema, "brand_positioning", "KnowledgeGuideOutput")
}

func TestGenerateSchema_PlannerOutput(t *testing.T) {
	schema := GenerateSchema[testPlannerOutput]()
	assertSchemaValid(t, schema, "PlannerOutput")

	// Check branding fields
	assertPropertyExists(t, schema, "brand_voice_guideline", "PlannerOutput")
	assertPropertyExists(t, schema, "section_message_map", "PlannerOutput")
}

func TestGenerateSchema_DeveloperOutput(t *testing.T) {
	schema := GenerateSchema[testDeveloperOutput]()
	assertSchemaValid(t, schema, "DeveloperOutput")

	// Check branding fields
	assertPropertyExists(t, schema, "organization_schema", "DeveloperOutput")
	assertPropertyExists(t, schema, "brand_meta", "DeveloperOutput")
}

func TestGenerateSchema_ExecutorOutput(t *testing.T) {
	schema := GenerateSchema[testExecutorOutput]()
	assertSchemaValid(t, schema, "ExecutorOutput")

	// Check branding fields
	assertPropertyExists(t, schema, "brand_voice_applied", "ExecutorOutput")
	assertPropertyExists(t, schema, "usp_mention_count", "ExecutorOutput")
	assertPropertyExists(t, schema, "brand_hashtags", "ExecutorOutput")
	assertPropertyExists(t, schema, "brand_storytelling", "ExecutorOutput")
}

func TestGenerateSchema_TrustManagerOutput(t *testing.T) {
	schema := GenerateSchema[testTrustManagerOutput]()
	assertSchemaValid(t, schema, "TrustManagerOutput")

	// Check branding fields
	assertPropertyExists(t, schema, "brand_authenticity", "TrustManagerOutput")
	assertPropertyExists(t, schema, "brand_consistency_checks", "TrustManagerOutput")
}

func TestGenerateSchema_AnalystOutput(t *testing.T) {
	schema := GenerateSchema[testAnalystOutput]()
	assertSchemaValid(t, schema, "AnalystOutput")

	// Check branding fields
	assertPropertyExists(t, schema, "brand_kpis", "AnalystOutput")
}

func TestSchemaToString_ReturnsValidJSON(t *testing.T) {
	s := SchemaToString[testKnowledgeGuideOutput]()
	if s == "" || s == "{}" {
		t.Error("SchemaToString returned empty or trivial result")
	}
	// Basic JSON check: starts with { and ends with }
	if s[0] != '{' || s[len(s)-1] != '}' {
		t.Errorf("SchemaToString result does not look like JSON: %s", s[:50])
	}
}

// assertSchemaValid checks that the schema map has the expected top-level keys.
func assertSchemaValid(t *testing.T, schema map[string]interface{}, name string) {
	t.Helper()

	if schema == nil {
		t.Fatalf("GenerateSchema[%s]() returned nil", name)
	}
	if _, ok := schema["type"]; !ok {
		t.Errorf("GenerateSchema[%s]() missing 'type' key", name)
	}
	if _, ok := schema["properties"]; !ok {
		t.Errorf("GenerateSchema[%s]() missing 'properties' key", name)
	}

	// Verify type is "object"
	if tp, ok := schema["type"].(string); ok {
		if tp != "object" {
			t.Errorf("GenerateSchema[%s]() type is %q, expected 'object'", name, tp)
		}
	}
}

// assertPropertyExists checks that a given property key is in the schema's properties map.
func assertPropertyExists(t *testing.T, schema map[string]interface{}, key, schemaName string) {
	t.Helper()

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("GenerateSchema[%s]() properties is not a map", schemaName)
	}
	if _, exists := props[key]; !exists {
		t.Errorf("GenerateSchema[%s]() missing expected property %q", schemaName, key)
	}
}
