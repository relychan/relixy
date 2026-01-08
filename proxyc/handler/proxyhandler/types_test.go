package proxyhandler

import (
	"testing"

	"github.com/hasura/goenvconf"
	"gotest.tools/v3/assert"
)

func TestNewRelixyHandlerOptions_GetEnvFunc(t *testing.T) {
	testCases := []struct {
		name     string
		options  NewRelixyHandlerOptions
		expected bool // true if should return custom func, false if should return default
	}{
		{
			name: "with custom GetEnv function",
			options: NewRelixyHandlerOptions{
				GetEnv: func(key string) (string, error) {
					return "custom-value", nil
				},
			},
			expected: true,
		},
		{
			name: "without GetEnv function",
			options: NewRelixyHandlerOptions{
				Method: "GET",
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			getEnvFunc := tc.options.GetEnvFunc()
			assert.Assert(t, getEnvFunc != nil)

			if tc.expected {
				// Test custom function
				val, err := getEnvFunc("test")
				assert.NilError(t, err)
				assert.Equal(t, "custom-value", val)
			} else {
				// Test default function (should be goenvconf.GetOSEnv)
				assert.Assert(t, getEnvFunc != nil)
			}
		})
	}
}

func TestAPIKeyCredentials(t *testing.T) {
	apiKey := goenvconf.NewEnvStringValue("test-api-key")
	creds := APIKeyCredentials{
		APIKey: &apiKey,
	}

	assert.Assert(t, creds.APIKey != nil)
	val, err := creds.APIKey.GetCustom(goenvconf.GetOSEnv)
	assert.NilError(t, err)
	assert.Equal(t, "test-api-key", val)
}

func TestBasicCredentials(t *testing.T) {
	username := goenvconf.NewEnvStringValue("test-user")
	password := goenvconf.NewEnvStringValue("test-pass")

	creds := BasicCredentials{
		Username: &username,
		Password: &password,
	}

	assert.Assert(t, creds.Username != nil)
	assert.Assert(t, creds.Password != nil)

	user, err := creds.Username.GetCustom(goenvconf.GetOSEnv)
	assert.NilError(t, err)
	assert.Equal(t, "test-user", user)

	pass, err := creds.Password.GetCustom(goenvconf.GetOSEnv)
	assert.NilError(t, err)
	assert.Equal(t, "test-pass", pass)
}

func TestOAuth2Credentials(t *testing.T) {
	clientID := goenvconf.NewEnvStringValue("client-id")
	clientSecret := goenvconf.NewEnvStringValue("client-secret")

	creds := OAuth2Credentials{
		ClientID:     &clientID,
		ClientSecret: &clientSecret,
		EndpointParams: map[string]goenvconf.EnvString{
			"scope": goenvconf.NewEnvStringValue("read write"),
		},
	}

	assert.Assert(t, creds.ClientID != nil)
	assert.Assert(t, creds.ClientSecret != nil)
	assert.Equal(t, 1, len(creds.EndpointParams))

	id, err := creds.ClientID.GetCustom(goenvconf.GetOSEnv)
	assert.NilError(t, err)
	assert.Equal(t, "client-id", id)

	secret, err := creds.ClientSecret.GetCustom(goenvconf.GetOSEnv)
	assert.NilError(t, err)
	assert.Equal(t, "client-secret", secret)

	scope, err := creds.EndpointParams["scope"].GetCustom(goenvconf.GetOSEnv)
	assert.NilError(t, err)
	assert.Equal(t, "read write", scope)
}

func TestInsertRouteOptions(t *testing.T) {
	customGetEnv := func(key string) (string, error) {
		return "custom-value", nil
	}

	options := InsertRouteOptions{
		GetEnv: customGetEnv,
	}

	assert.Assert(t, options.GetEnv != nil)
	val, err := options.GetEnv("test")
	assert.NilError(t, err)
	assert.Equal(t, "custom-value", val)
}

func TestOAuth2CredentialsErrors(t *testing.T) {
	t.Run("errOAuth2ClientCredentialsRequired", func(t *testing.T) {
		err := errOAuth2ClientCredentialsRequired
		assert.ErrorContains(t, err, "clientId and clientSecret")
	})

	t.Run("errOAuth2TokenURLRequired", func(t *testing.T) {
		err := errOAuth2TokenURLRequired
		assert.ErrorContains(t, err, "tokenUrl")
	})
}
