package proxyhandler

import (
	"errors"

	"github.com/hasura/goenvconf"
	"github.com/relychan/gotransform"
)

var (
	errOAuth2ClientCredentialsRequired = errors.New("clientId and clientSecret must not be empty")
	errOAuth2TokenURLRequired          = errors.New(
		"tokenUrl is required in the OAuth2 Client Credentials flow",
	)
)

// ProxyActionType represents enums of proxy types.
type ProxyActionType string

// InsertRouteOptions represents options for inserting routes.
type InsertRouteOptions struct {
	GetEnv goenvconf.GetEnvFunc
}

// APIKeyCredentials holds apiKey credentials of the security scheme.
type APIKeyCredentials struct {
	APIKey *goenvconf.EnvString `json:"apiKey,omitempty" yaml:"apiKey,omitempty"`
}

// BasicCredentials holds basic credentials of the security scheme.
type BasicCredentials struct {
	Username *goenvconf.EnvString `json:"username,omitempty" yaml:"username,omitempty"`
	Password *goenvconf.EnvString `json:"password,omitempty" yaml:"password,omitempty"`
}

// OAuth2Credentials holds OAuth2 credentials of the security scheme.
type OAuth2Credentials struct {
	ClientID     *goenvconf.EnvString `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret *goenvconf.EnvString `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	// Optional query parameters for the token and refresh URLs.
	EndpointParams map[string]goenvconf.EnvString `json:"endpointParams,omitempty" yaml:"endpointParams,omitempty"`
}

// RelixyResponseRawConfig represents configurations for the proxy response.
type RelixyResponseRawConfig struct {
	// HTTP error code will be used if the response body has errors.
	// If not set, forward the HTTP status from the GraphQL response which is usually 200 OK.
	HTTPErrorCode *int `json:"httpErrorCode,omitempty" yaml:"httpErrorCode,omitempty" jsonschema:"minimum=400,maximum=599,default=400"`
	// Configurations for transforming response data.
	Transform *gotransform.TemplateTransformerConfig `json:"transform,omitempty" yaml:"transform,omitempty"`
}

// IsZero checks if the configuration is empty.
func (conf RelixyResponseRawConfig) IsZero() bool {
	return conf.HTTPErrorCode == nil &&
		(conf.Transform == nil || conf.Transform.IsZero())
}

// RelixyResponseConfig represents configurations for the proxy response.
type RelixyResponseConfig struct {
	// HTTP error code will be used if the response body has errors.
	// If not set, forward the HTTP status from the GraphQL response which is usually 200 OK.
	HTTPErrorCode *int
	// Configurations for transforming response data.
	Transform gotransform.TemplateTransformer
}

// NewRelixyResponseConfig creates a [RelixyResponseConfig] from raw configurations.
func NewRelixyResponseConfig(
	config *RelixyResponseRawConfig,
	getEnv goenvconf.GetEnvFunc,
) (RelixyResponseConfig, error) {
	result := RelixyResponseConfig{}

	if config == nil {
		return result, nil
	}

	result.HTTPErrorCode = config.HTTPErrorCode

	if config.Transform != nil {
		transformer, err := gotransform.NewTransformerFromConfig("", *config.Transform, getEnv)
		if err != nil {
			return result, err
		}

		result.Transform = transformer
	}

	return result, nil
}

// IsZero checks if the configuration is empty.
func (conf RelixyResponseConfig) IsZero() bool {
	return conf.HTTPErrorCode == nil &&
		(conf.Transform == nil || conf.Transform.IsZero())
}
