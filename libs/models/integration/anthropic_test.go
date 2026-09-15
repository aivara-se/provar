//go:build integration

package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/thani-sh/provar/libs/domain"
	"github.com/thani-sh/provar/libs/models"
)

const (
	defaultAnthropicModel = "claude-sonnet-5"
	envAnthropicAPIKey    = "ANTHROPIC_API_KEY"
	envAnthropicModel     = "ANTHROPIC_API_MODEL"
	envAnthropicBaseURL   = "ANTHROPIC_API_URL"
	skipAnthropicMsg      = "skipping Anthropic integration test; credentials not set"
)

func TestAnthropicClient_Integration(t *testing.T) {
	apiKey, model, baseURL := resolveProviderConfig(domain.ProviderAnthropic, envAnthropicAPIKey, envAnthropicModel, envAnthropicBaseURL, defaultAnthropicModel)
	if apiKey == "" {
		t.Skip(skipAnthropicMsg)
	}
	client, err := models.NewClient(models.Anthropic, apiKey, baseURL, model)
	if err != nil {
		t.Fatalf("failed to create Anthropic client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()
	session, err := client.CreateSession(ctx, "")
	if err != nil {
		t.Fatalf("failed to create Anthropic session: %v", err)
	}
	err = session.Send(ctx, []models.Attachment{{Type: models.AttachmentTypeText, Text: testPrompt}})
	if err != nil {
		t.Fatalf("failed to send message to Anthropic: %v", err)
	}
	var response string
	for chunk := range session.Recv() {
		response += chunk
	}
	if !strings.Contains(response, expectedSubstring) {
		t.Errorf("expected response to contain %q, but got %q", expectedSubstring, response)
	}
}

func TestAnthropicClient_Integration_ToolCalling(t *testing.T) {
	apiKey, model, baseURL := resolveProviderConfig(domain.ProviderAnthropic, envAnthropicAPIKey, envAnthropicModel, envAnthropicBaseURL, defaultAnthropicModel)
	if apiKey == "" {
		t.Skip(skipAnthropicMsg)
	}
	client, err := models.NewClient(models.Anthropic, apiKey, baseURL, model)
	if err != nil {
		t.Fatalf("failed to create Anthropic client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()
	var toolCalled bool
	secretCode := randomSecretCode()
	tool := &secretTool{code: secretCode, called: &toolCalled}
	session, err := client.CreateSession(ctx, "", tool)
	if err != nil {
		t.Fatalf("failed to create Anthropic session: %v", err)
	}
	err = session.Send(ctx, []models.Attachment{{Type: models.AttachmentTypeText, Text: toolPromptSecret}})
	if err != nil {
		t.Fatalf("failed to send message to Anthropic: %v", err)
	}
	var response string
	for chunk := range session.Recv() {
		response += chunk
	}
	if !toolCalled {
		t.Error("expected tool to be executed by Anthropic session")
	}
	if !strings.Contains(response, secretCode) {
		t.Errorf("expected response to contain %q, but got %q", secretCode, response)
	}
}
