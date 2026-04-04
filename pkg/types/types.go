package types

// Keyword represents an SEO keyword with metadata.
type Keyword struct {
	Term           string `json:"term"`
	SearchVolume   int    `json:"search_volume"`
	Difficulty     int    `json:"difficulty"`
	Intent         string `json:"intent"`
	Relevance      string `json:"relevance"`
}

// ContentOutline represents a structured content outline.
type ContentOutline struct {
	Title       string           `json:"title"`
	Slug        string           `json:"slug"`
	Sections    []OutlineSection `json:"sections"`
	TargetWords int              `json:"target_words"`
}

// OutlineSection represents a section within a content outline.
type OutlineSection struct {
	Heading    string   `json:"heading"`
	KeyPoints  []string `json:"key_points"`
	Keywords   []string `json:"keywords"`
}

// MetaTag represents HTML meta tags for SEO.
type MetaTag struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Canonical   string `json:"canonical,omitempty"`
	OGTitle     string `json:"og_title,omitempty"`
	OGDesc      string `json:"og_description,omitempty"`
}

// SchemaMarkup represents structured data markup.
type SchemaMarkup struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
	RawJSON    string            `json:"raw_json"`
}

// TrustScore represents E-E-A-T evaluation scores.
type TrustScore struct {
	Experience    int    `json:"experience"`
	Expertise     int    `json:"expertise"`
	Authority     int    `json:"authority"`
	Trust         int    `json:"trust"`
	Overall       int    `json:"overall"`
	Feedback      string `json:"feedback"`
	NeedsRevision bool   `json:"needs_revision"`
}

// KPI represents a key performance indicator.
type KPI struct {
	Name       string `json:"name"`
	Current    string `json:"current"`
	Target     string `json:"target"`
	Suggestion string `json:"suggestion"`
}
