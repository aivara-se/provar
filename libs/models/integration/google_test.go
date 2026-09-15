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
	defaultGoogleModel = "gemini-3.7-flash"
	envGoogleAPIKey    = "GEMINI_API_KEY"
	envGoogleModel     = "GEMINI_API_MODEL"
	skipGoogleMsg      = "skipping Google integration test; credentials not set"
)

func TestGoogleClient_Integration(t *testing.T) {
	apiKey, model, baseURL := resolveProviderConfig(domain.ProviderGoogle, envGoogleAPIKey, envGoogleModel, "", defaultGoogleModel)
	if apiKey == "" {
		t.Skip(skipGoogleMsg)
	}
	client, err := models.NewClient(models.Google, apiKey, baseURL, model)
	if err != nil {
		t.Fatalf("failed to create Google client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()
	session, err := client.CreateSession(ctx, "")
	if err != nil {
		t.Fatalf("failed to create Google session: %v", err)
	}
	err = session.Send(ctx, []models.Attachment{{Type: models.AttachmentTypeText, Text: testPrompt}})
	if err != nil {
		t.Fatalf("failed to send message to Google: %v", err)
	}
	var response string
	for chunk := range session.Recv() {
		response += chunk
	}
	if !strings.Contains(response, expectedSubstring) {
		t.Errorf("expected response to contain %q, but got %q", expectedSubstring, response)
	}
}

func TestGoogleClient_Integration_ToolCalling(t *testing.T) {
	apiKey, model, baseURL := resolveProviderConfig(domain.ProviderGoogle, envGoogleAPIKey, envGoogleModel, "", defaultGoogleModel)
	if apiKey == "" {
		t.Skip(skipGoogleMsg)
	}
	client, err := models.NewClient(models.Google, apiKey, baseURL, model)
	if err != nil {
		t.Fatalf("failed to create Google client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout)
	defer cancel()
	var toolCalled bool
	secretCode := randomSecretCode()
	tool := &secretTool{code: secretCode, called: &toolCalled}
	session, err := client.CreateSession(ctx, "", tool)
	if err != nil {
		t.Fatalf("failed to create Google session: %v", err)
	}
	err = session.Send(ctx, []models.Attachment{{Type: models.AttachmentTypeText, Text: toolPromptSecret}})
	if err != nil {
		t.Fatalf("failed to send message to Google: %v", err)
	}
	var response string
	for chunk := range session.Recv() {
		response += chunk
	}
	if !toolCalled {
		t.Error("expected tool to be executed by Google session")
	}
	if !strings.Contains(response, secretCode) {
		t.Errorf("expected response to contain %q, but got %q", secretCode, response)
	}
}
