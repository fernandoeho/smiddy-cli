package claude

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// Client wraps the Anthropic SDK
type Client struct {
	sdk *anthropic.Client
}

// Message represents a single conversation turn
type Message struct {
	Role    string
	Content string
}

// NewClient creates a new Claude client using ANTHROPIC_API_KEY from env
func NewClient() (*Client, error) {
	sdk := anthropic.NewClient()
	return &Client{sdk: &sdk}, nil
}

// Complete sends a system prompt + messages and returns Claude's response
func (c *Client) Complete(ctx context.Context, systemPrompt string, messages []Message) (string, error) {
	var sdkMessages []anthropic.MessageParam

	for _, m := range messages {
		if m.Role == "user" {
			sdkMessages = append(sdkMessages, anthropic.NewUserMessage(
				anthropic.NewTextBlock(m.Content),
			))
		} else if m.Role == "assistant" {
			sdkMessages = append(sdkMessages, anthropic.NewAssistantMessage(
				anthropic.NewTextBlock(m.Content),
			))
		}
	}

	resp, err := c.sdk.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5,
		MaxTokens: 16000,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: sdkMessages,
	})
	if err != nil {
		return "", fmt.Errorf("claude API error: %w", err)
	}

	if len(resp.Content) == 0 {
		return "", fmt.Errorf("empty response from claude")
	}

	for _, block := range resp.Content {
		if text, ok := block.AsAny().(anthropic.TextBlock); ok {
			return text.Text, nil
		}
	}

	return "", fmt.Errorf("no text block in claude response")
}
