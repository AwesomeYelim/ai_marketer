package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Client wraps the Anthropic SDK client.
type Client struct {
	inner *anthropic.Client
	model anthropic.Model
}

// NewClient creates a new LLM client.
// If apiKey is empty, it uses the ANTHROPIC_API_KEY environment variable.
func NewClient(apiKey string, model string) *Client {
	var opts []option.RequestOption
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}

	client := anthropic.NewClient(opts...)

	m := anthropic.Model(model)
	if model == "" {
		m = anthropic.ModelClaudeOpus4_6
	}

	return &Client{
		inner: &client,
		model: m,
	}
}

// Complete sends a simple text completion request.
func (c *Client) Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	params := anthropic.MessageNewParams{
		Model:     c.model,
		MaxTokens: 8192,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
	}

	if systemPrompt != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: systemPrompt},
		}
	}

	resp, err := c.inner.Messages.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("claude API error: %w", err)
	}

	if len(resp.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return resp.Content[0].Text, nil
}

// CompleteStructured sends a request and parses the response into the given type T.
// It instructs Claude to respond in JSON matching the schema of T.
func CompleteStructured[T any](c *Client, ctx context.Context, systemPrompt, userPrompt string) (*T, error) {
	schema := SchemaToString[T]()

	jsonSystemPrompt := systemPrompt + "\n\n" +
		"IMPORTANT: You MUST respond with ONLY a valid JSON object matching this schema. " +
		"Do NOT include any text before or after the JSON. Do NOT use markdown code fences.\n\n" +
		"JSON Schema:\n" + schema

	params := anthropic.MessageNewParams{
		Model:     c.model,
		MaxTokens: 8192,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
		System: []anthropic.TextBlockParam{
			{Text: jsonSystemPrompt},
		},
	}

	resp, err := c.inner.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	if len(resp.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude")
	}

	raw := resp.Content[0].Text

	var result T
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("failed to parse structured response: %w\nraw response: %s", err, raw)
	}

	return &result, nil
}
