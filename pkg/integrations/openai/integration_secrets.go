package openai

import (
	"fmt"
	"strings"

	"github.com/superplanehq/superplane/pkg/core"
)

const (
	integrationSecretOpenAIAPIKey  = "OPENAI_API_KEY"
	integrationSecretOpenAIBaseURL = "OPENAI_BASE_URL"
	openAISecretUsage              = `An OpenAI API key is available in the OPENAI_API_KEY environment variable.
Do not print the key.`
	openAIBaseURLUsage = `The integration base URL is available in the OPENAI_BASE_URL environment variable.
Use that URL for OpenAI-compatible endpoints.`
)

func (o *OpenAI) ResolveSecrets(ctx core.IntegrationSecretContext) (core.IntegrationSecrets, error) {
	apiKey, err := ctx.Integration.GetConfig("apiKey")
	if err != nil {
		return core.IntegrationSecrets{}, fmt.Errorf("failed to get API key: %w", err)
	}

	key := strings.TrimSpace(string(apiKey))
	if key == "" {
		return core.IntegrationSecrets{}, fmt.Errorf("apiKey is required")
	}

	values := map[string][]byte{
		integrationSecretOpenAIAPIKey: []byte(key),
	}
	usage := openAISecretUsage

	// A custom base URL says the integration talks to an OpenAI-compatible
	// endpoint. Agent CLIs on runners read OPENAI_BASE_URL, so the value must
	// reach them or they call api.openai.com with a foreign key.
	if baseURL := customBaseURL(ctx); baseURL != "" {
		values[integrationSecretOpenAIBaseURL] = []byte(baseURL)
		usage += "\n" + openAIBaseURLUsage
	}

	return core.IntegrationSecrets{
		Values: values,
		Usage:  usage,
	}, nil
}

func customBaseURL(ctx core.IntegrationSecretContext) string {
	baseURL, err := ctx.Integration.GetConfig("baseURL")
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(baseURL))
}
