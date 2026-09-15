//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/thani-sh/provar/libs/domain"
	"github.com/thani-sh/provar/libs/models"
)

const (
	testPrompt         = "Say 'hello' in reverse back to me in lower case"
	expectedSubstring  = "olleh"
	integrationTimeout = 60 * time.Second
	toolNameSecret     = "get_secret_code"
	toolDescSecret     = "Returns the secret authorization code for a specified username"
	toolParamsSecret   = `{"type":"object","properties":{"username":{"type":"string","description":"The username to query"}},"required":["username"]}`
	toolPromptSecret   = "What is the secret authorization code for user 'alice'? You must call get_secret_code to find it."
	secretCodeFormat   = "secret-token-%d"
	flagSettingsName   = "settings"
	flagSettingsUsage  = "load API keys and model settings from ~/.provar/settings.yml"
	envUseSettings     = "USE_SETTINGS"
	envTrue            = "true"
)

var useSettings = flag.Bool(flagSettingsName, false, flagSettingsUsage)

func resolveProviderConfig(p domain.Provider, envKey, envModel, envBaseURL, defaultModel string) (string, string, string) {
	var apiKey, model, baseURL string
	if *useSettings || os.Getenv(envUseSettings) == envTrue {
		if s, err := domain.LoadSettings(); err == nil {
			if cfg, ok := s.Providers[string(p)]; ok && cfg.APIKey != "" {
				apiKey = cfg.APIKey
				model = cfg.Model
				baseURL = cfg.BaseURL
			}
		}
	}
	if envK := os.Getenv(envKey); envK != "" {
		apiKey = envK
	}
	if envM := os.Getenv(envModel); envM != "" {
		model = envM
	}
	if envU := os.Getenv(envBaseURL); envU != "" {
		baseURL = envU
	}
	if model == "" {
		model = defaultModel
	}
	return apiKey, model, baseURL
}

func randomSecretCode() string {
	return fmt.Sprintf(secretCodeFormat, time.Now().UnixNano())
}

type secretTool struct {
	code   string
	called *bool
}

func (t *secretTool) Name() string {
	return toolNameSecret
}

func (t *secretTool) Description() string {
	return toolDescSecret
}

func (t *secretTool) Parameters() json.RawMessage {
	return json.RawMessage(toolParamsSecret)
}

func (t *secretTool) Execute(ctx context.Context, args json.RawMessage) (models.ToolResult, error) {
	if t.called != nil {
		*t.called = true
	}
	return models.ToolResult{
		Content: []models.Attachment{
			{Type: models.AttachmentTypeText, Text: t.code},
		},
	}, nil
}
