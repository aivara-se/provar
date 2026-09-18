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
	defaultOpenAIModel = "gpt-5.6-terra"
	envOpenAIAPIKey    = "OPENAI_API_KEY"
	envOpenAIModel     = "OPENAI_API_MODEL"
	envOpenAIBaseURL   = "OPENAI_API_URL"
	skipOpenAIMsg      = "skipping OpenAI integration test; credentials not set"
)

func TestOpenAIClient_Integration(t *testing.T) {
	apiKey, model, baseURL := resolveProviderConfig(domain.ProviderOpenAI, envOpenAIAPIKey, envOpenAIModel, envOpenAIBaseURL, defaultOpenAIModel)
	if apiKey == "" {
		t.Skip(skipOpenAIMsg)
	}
	client, err := models.NewClient(models.OpenAI, apiKey, baseURL, model)
	if err != nil {
		t.Fatalf("failed to create OpenAI client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()
	session, err := client.CreateSession(ctx, "")
	if err != nil {
		t.Fatalf("failed to create OpenAI session: %v", err)
	}
	err = session.Send(ctx, []models.Attachment{{Type: models.AttachmentTypeText, Text: testPrompt}})
	if err != nil {
		t.Fatalf("failed to send message to OpenAI: %v", err)
	}
	var response string
	for chunk := range session.Recv() {
		response += chunk
	}
	if !strings.Contains(response, expectedSubstring) {
		t.Errorf("expected response to contain %q, but got %q", expectedSubstring, response)
	}
}

func TestOpenAIClient_Integration_ToolCalling(t *testing.T) {
	apiKey, model, baseURL := resolveProviderConfig(domain.ProviderOpenAI, envOpenAIAPIKey, envOpenAIModel, envOpenAIBaseURL, defaultOpenAIModel)
	if apiKey == "" {
		t.Skip(skipOpenAIMsg)
	}
	client, err := models.NewClient(models.OpenAI, apiKey, baseURL, model)
	if err != nil {
		t.Fatalf("failed to create OpenAI client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()
	var toolCalled bool
	secretCode := randomSecretCode()
	tool := &secretTool{code: secretCode, called: &toolCalled}
	session, err := client.CreateSession(ctx, "", tool)
	if err != nil {
		t.Fatalf("failed to create OpenAI session: %v", err)
	}
	err = session.Send(ctx, []models.Attachment{{Type: models.AttachmentTypeText, Text: toolPromptSecret}})
	if err != nil {
		t.Fatalf("failed to send message to OpenAI: %v", err)
	}
	var response string
	for chunk := range session.Recv() {
		response += chunk
	}
	if !toolCalled {
		t.Error("expected tool to be executed by OpenAI session")
	}
	if !strings.Contains(response, secretCode) {
		t.Errorf("expected response to contain %q, but got %q", secretCode, response)
	}
}
