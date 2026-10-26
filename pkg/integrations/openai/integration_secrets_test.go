package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/superplanehq/superplane/pkg/core"
	contexts "github.com/superplanehq/superplane/test/support/contexts"
)

func TestResolveSecrets(t *testing.T) {
	t.Parallel()

	secrets, err := (&OpenAI{}).ResolveSecrets(core.IntegrationSecretContext{
		Integration: &contexts.IntegrationContext{
			Configuration: map[string]any{
				"apiKey": "sk-test",
			},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, []byte("sk-test"), secrets.Values[integrationSecretOpenAIAPIKey])
	assert.NotContains(t, secrets.Values, integrationSecretOpenAIBaseURL)
	assert.Contains(t, secrets.Usage, "OPENAI_API_KEY")
	assert.NotContains(t, secrets.Usage, "sk-test")
	assert.Empty(t, secrets.Setup)
}

func TestResolveSecretsIncludesCustomBaseURL(t *testing.T) {
	t.Parallel()

	secrets, err := (&OpenAI{}).ResolveSecrets(core.IntegrationSecretContext{
		Integration: &contexts.IntegrationContext{
			Configuration: map[string]any{
				"apiKey":  "sk-test",
				"baseURL": " https://gateway.example.com/v1 ",
			},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, []byte("sk-test"), secrets.Values[integrationSecretOpenAIAPIKey])
	assert.Equal(t, []byte("https://gateway.example.com/v1"), secrets.Values[integrationSecretOpenAIBaseURL])
	assert.Contains(t, secrets.Usage, "OPENAI_BASE_URL")
	assert.NotContains(t, secrets.Usage, "sk-test")
}

func TestResolveSecretsSkipsBlankBaseURL(t *testing.T) {
	t.Parallel()

	secrets, err := (&OpenAI{}).ResolveSecrets(core.IntegrationSecretContext{
		Integration: &contexts.IntegrationContext{
			Configuration: map[string]any{
				"apiKey":  "sk-test",
				"baseURL": "   ",
			},
		},
	})
	require.NoError(t, err)
	assert.NotContains(t, secrets.Values, integrationSecretOpenAIBaseURL)
	assert.NotContains(t, secrets.Usage, "OPENAI_BASE_URL")
}

func TestResolveSecretsRequiresAPIKey(t *testing.T) {
	t.Parallel()

	_, err := (&OpenAI{}).ResolveSecrets(core.IntegrationSecretContext{
		Integration: &contexts.IntegrationContext{
			Configuration: map[string]any{
				"baseURL": "https://gateway.example.com/v1",
			},
		},
	})
	require.Error(t, err)
}
