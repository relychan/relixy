package schema

import (
	"encoding/json"
	"fmt"

	"github.com/hasura/goenvconf"
	"github.com/invopop/jsonschema"
	"github.com/relychan/gohttpc/authc/authscheme"
	"github.com/relychan/goutils"
	wk8orderedmap "github.com/wk8/go-ordered-map/v2"
	"go.yaml.in/yaml/v4"
)

// SecuritySchemer abstracts an interface of SecurityScheme.
type SecuritySchemer interface {
	GetType() SecuritySchemeType
	Validate() error
}

// RelixySecurityScheme contains authentication configurations.
// The schema follows [OpenAPI 3] specification
//
// [OpenAPI 3]: https://swagger.io/docs/specification/authentication
type RelixySecurityScheme struct {
	SecuritySchemer
}

type rawSecurityScheme struct {
	Type SecuritySchemeType `json:"type" yaml:"type"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelixySecurityScheme) UnmarshalJSON(b []byte) error {
	var rawScheme rawSecurityScheme

	err := json.Unmarshal(b, &rawScheme)
	if err != nil {
		return err
	}

	switch rawScheme.Type {
	case APIKeyScheme:
		var config RelixyAPIKeyAuthConfig

		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case BasicAuthScheme:
		var config RelixyBasicAuthConfig

		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case HTTPAuthScheme:
		var config RelixyHTTPAuthConfig

		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case OAuth2Scheme:
		var config RelixyOAuth2Config

		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case OpenIDConnectScheme:
		var config RelixyOpenIDConnectConfig

		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case CookieAuthScheme:
		j.SecuritySchemer = &RelixyCookieAuthConfig{
			Type: rawScheme.Type,
		}
	case MutualTLSScheme:
		j.SecuritySchemer = &RelixyMutualTLSAuthConfig{
			Type: rawScheme.Type,
		}
	default:
		return fmt.Errorf("%w; got: %s", ErrInvalidSecuritySchemeType, rawScheme.Type)
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (j RelixySecurityScheme) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.SecuritySchemer)
}

// MarshalYAML implements the yaml.Marshaler interface.
func (j RelixySecurityScheme) MarshalYAML() (any, error) {
	return yaml.Marshal(j.SecuritySchemer)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (j *RelixySecurityScheme) UnmarshalYAML(value *yaml.Node) error {
	var rawScheme rawSecurityScheme

	err := value.Decode(&rawScheme)
	if err != nil {
		return err
	}

	switch rawScheme.Type {
	case APIKeyScheme:
		var config RelixyAPIKeyAuthConfig

		err := value.Decode(&config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case BasicAuthScheme:
		var config RelixyBasicAuthConfig

		err := value.Decode(&config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case HTTPAuthScheme:
		var config RelixyHTTPAuthConfig

		err := value.Decode(&config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case OAuth2Scheme:
		var config RelixyOAuth2Config

		err := value.Decode(&config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case OpenIDConnectScheme:
		var config RelixyOpenIDConnectConfig

		err := value.Decode(&config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case CookieAuthScheme:
		j.SecuritySchemer = &RelixyCookieAuthConfig{
			Type: rawScheme.Type,
		}
	case MutualTLSScheme:
		j.SecuritySchemer = &RelixyMutualTLSAuthConfig{
			Type: rawScheme.Type,
		}
	default:
		return fmt.Errorf("%w; got: %s", ErrInvalidSecuritySchemeType, rawScheme.Type)
	}

	return nil
}

// Validate if the current instance is valid.
func (j *RelixySecurityScheme) Validate() error {
	if j.SecuritySchemer == nil {
		return errSecuritySchemerRequired
	}

	return j.SecuritySchemer.Validate()
}

// JSONSchema defines a custom definition for JSON schema.
func (RelixySecurityScheme) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Description: "Configurations for the apiKey authentication",
				Ref:         "#/$defs/RelixyAPIKeyAuthConfig",
			},
			{
				Description: "Configurations for the http authentication",
				Ref:         "#/$defs/RelixyHTTPAuthConfig",
			},
			{
				Description: "Configurations for the basic authentication",
				Ref:         "#/$defs/RelixyBasicAuthConfig",
			},
			{
				Description: "Configurations for OAuth 2.0 authentication",
				Ref:         "#/$defs/RelixyOAuth2Config",
			},
			{
				Description: "Configurations for OpenID Connect authentication",
				Ref:         "#/$defs/RelixyOpenIDConnectConfig",
			},
			{
				Description: "Configurations for Cookie authentication",
				Ref:         "#/$defs/RelixyCookieAuthConfig",
			},
			{
				Description: "Configurations for mutualTLS authentication",
				Ref:         "#/$defs/RelixyMutualTLSAuthConfig",
			},
		},
	}
}

// RelixyAPIKeyAuthConfig contains configurations for [apiKey authentication]
//
// [apiKey authentication]: https://swagger.io/docs/specification/authentication/api-keys/
type RelixyAPIKeyAuthConfig struct {
	// The type of the security scheme that is always 'apiKey'.
	Type SecuritySchemeType `json:"type" yaml:"type" jsonschema:"enum=apiKey"`
	// The location of the API key. Valid values are 'query', 'header', or 'cookie'.
	In authscheme.AuthLocation `json:"in" yaml:"in" jsonschema:"enum=query,enum=header,enum=cookie"`
	// The name of the header, query or cookie parameter to be used.
	Name string `json:"name" yaml:"name"`
	// The value of the API key that is either a literal value of an environment variable.
	Value goenvconf.EnvString `json:"value" yaml:"value"`
	// Declares this security scheme to be deprecated.
	// Consumers SHOULD refrain from usage of the declared scheme. Default value is false.
	Deprecated bool `json:"deprecated,omitempty" mapstructure:"deprecated" yaml:"deprecated,omitempty"`
}

var _ SecuritySchemer = (*RelixyAPIKeyAuthConfig)(nil)

// NewAPIKeyAuthConfig creates a new APIKeyAuthConfig instance.
func NewAPIKeyAuthConfig(
	name string,
	in authscheme.AuthLocation,
	value goenvconf.EnvString,
) *RelixyAPIKeyAuthConfig {
	return &RelixyAPIKeyAuthConfig{
		Name:  name,
		In:    in,
		Value: value,
	}
}

// Validate if the current instance is valid.
func (j *RelixyAPIKeyAuthConfig) Validate() error {
	if j.Name == "" {
		return fmt.Errorf("name %w for apiKey security", ErrFieldRequired)
	}

	err := j.In.Validate()
	if err != nil {
		return err
	}

	return nil
}

// GetType get the type of security scheme.
func (*RelixyAPIKeyAuthConfig) GetType() SecuritySchemeType {
	return APIKeyScheme
}

// RelixyHTTPAuthConfig contains configurations for http authentication
// If the scheme is [bearer], the authenticator follows OpenAPI 3 specification.
//
// [bearer]: https://swagger.io/docs/specification/authentication/bearer-authentication
type RelixyHTTPAuthConfig struct {
	// The type of the security scheme that is always 'http'.
	Type SecuritySchemeType `json:"type" yaml:"type" jsonschema:"enum=http"`
	// The name of the header to be used.
	Name string `json:"name" yaml:"name"`
	// The name of the HTTP Authentication scheme to be used in the Authorization header as defined in RFC9110.
	Scheme string `json:"scheme" yaml:"scheme"`
	// The value of the HTTP credential.
	Value goenvconf.EnvString `json:"value" yaml:"value"`
	// A hint to the client to identify how the bearer token is formatted.
	// Bearer tokens are usually generated by an authorization server,
	// so this information is primarily for documentation purposes.
	BearerFormat string `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	// Declares this security scheme to be deprecated.
	// Consumers SHOULD refrain from usage of the declared scheme. Default value is false.
	Deprecated bool `json:"deprecated,omitempty" mapstructure:"deprecated" yaml:"deprecated,omitempty"`
}

var _ SecuritySchemer = (*RelixyHTTPAuthConfig)(nil)

// NewHTTPAuthConfig creates a new HTTPAuthConfig instance.
func NewHTTPAuthConfig(
	scheme string,
	name string,
	value goenvconf.EnvString,
) *RelixyHTTPAuthConfig {
	return &RelixyHTTPAuthConfig{
		Type:   HTTPAuthScheme,
		Name:   name,
		Scheme: scheme,
		Value:  value,
	}
}

// Validate if the current instance is valid.
func (ss *RelixyHTTPAuthConfig) Validate() error {
	if ss.Scheme == "" {
		return fmt.Errorf("scheme %w for http security", ErrFieldRequired)
	}

	return nil
}

// GetType get the type of security scheme.
func (*RelixyHTTPAuthConfig) GetType() SecuritySchemeType {
	return HTTPAuthScheme
}

// RelixyBasicAuthConfig contains configurations for the [basic] authentication.
//
// [basic]: https://swagger.io/docs/specification/authentication/basic-authentication
type RelixyBasicAuthConfig struct {
	// The type of the security scheme that is always 'basic'.
	Type SecuritySchemeType `json:"type" yaml:"type" jsonschema:"enum=basic"`
	// The name of the header to be used.
	Name string `json:"name" yaml:"name"`
	// The username of this basic authentication.
	Username goenvconf.EnvString `json:"username" yaml:"username"`
	// The password of this basic authentication.
	Password goenvconf.EnvString `json:"password" yaml:"password"`
	// Declares this security scheme to be deprecated.
	// Consumers SHOULD refrain from usage of the declared scheme. Default value is false.
	Deprecated bool `json:"deprecated,omitempty" mapstructure:"deprecated" yaml:"deprecated,omitempty"`
}

// NewBasicAuthConfig creates a new BasicAuthConfig instance.
func NewBasicAuthConfig(username, password goenvconf.EnvString) *RelixyBasicAuthConfig {
	return &RelixyBasicAuthConfig{
		Type:     BasicAuthScheme,
		Username: username,
		Password: password,
	}
}

// Validate if the current instance is valid.
func (*RelixyBasicAuthConfig) Validate() error {
	return nil
}

// GetType get the type of security scheme.
func (*RelixyBasicAuthConfig) GetType() SecuritySchemeType {
	return BasicAuthScheme
}

// OAuthFlow contains flow configurations for [OAuth 2.0] authentication.
//
// [OAuth 2.0]: https://swagger.io/docs/specification/authentication/oauth2
type OAuthFlow struct {
	// The authorization URL to be used for this flow. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	AuthorizationURL string `json:"authorizationUrl,omitempty" mapstructure:"authorizationUrl" yaml:"authorizationUrl,omitempty"`
	// The token URL to be used for this flow. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	TokenURL *goenvconf.EnvString `json:"tokenUrl,omitempty"         mapstructure:"tokenUrl"         yaml:"tokenUrl,omitempty"`
	// The URL to be used for obtaining refresh tokens. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	RefreshURL string `json:"refreshUrl,omitempty"       mapstructure:"refreshUrl"       yaml:"refreshUrl,omitempty"`
	// The available scopes for the OAuth2 security scheme. A map between the scope name and a short description for it. The map MAY be empty.
	Scopes map[string]string `json:"scopes,omitempty"           mapstructure:"scopes"           yaml:"scopes,omitempty"`
	// Client ID of the authentication.
	ClientID *goenvconf.EnvString `json:"clientId,omitempty"         mapstructure:"clientId"         yaml:"clientId,omitempty"`
	// Client Secret of the authentication.
	ClientSecret *goenvconf.EnvString `json:"clientSecret,omitempty"     mapstructure:"clientSecret"     yaml:"clientSecret,omitempty"`
	// Optional query parameters that are added to the request endpoint.
	EndpointParams map[string]goenvconf.EnvString `json:"endpointParams,omitempty"   mapstructure:"endpointParams"   yaml:"endpointParams,omitempty"`
	// The device authorization URL to be used for this flow. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	DeviceAuthorizationURL string `json:"deviceAuthorizationUrl,omitempty" mapstructure:"deviceAuthorizationUrl" yaml:"deviceAuthorizationUrl,omitempty"`
}

// Validate if the current instance is valid.
func (ss OAuthFlow) Validate(flowType OAuthFlowType) error {
	if (ss.TokenURL == nil && (flowType == PasswordFlow ||
		flowType == ClientCredentialsFlow ||
		flowType == AuthorizationCodeFlow)) ||
		ss.TokenURL.IsZero() {
		return fmt.Errorf("tokenUrl %w for oauth2 %s flow", ErrFieldRequired, flowType)
	}

	if flowType != ClientCredentialsFlow {
		return nil
	}

	if ss.ClientID == nil || ss.ClientID.IsZero() {
		return fmt.Errorf("clientId %w for oauth2 %s flow", ErrFieldRequired, flowType)
	}

	if ss.ClientSecret == nil || ss.ClientSecret.IsZero() {
		return fmt.Errorf("clientSecret %w for oauth2 %s flow", ErrFieldRequired, flowType)
	}

	return nil
}

// RelixyOAuth2Config contains configurations for [OAuth 2.0] API specification
//
// [OAuth 2.0]: https://swagger.io/docs/specification/authentication/oauth2
type RelixyOAuth2Config struct {
	// The type of the security scheme that is always 'oauth2'.
	Type SecuritySchemeType `json:"type"  mapstructure:"type"  yaml:"type" jsonschema:"enum=oauth2"`
	// An object containing configuration information for the flow types supported.
	Flows map[OAuthFlowType]OAuthFlow `json:"flows" mapstructure:"flows" yaml:"flows"`
	// URL to the OAuth2 authorization server metadata RFC8414. TLS is required.
	OAuth2MetadataURL string `json:"oauth2MetadataUrl,omitempty" mapstructure:"oauth2MetadataUrl" yaml:"oauth2MetadataUrl,omitempty"`
	// Declares this security scheme to be deprecated.
	// Consumers SHOULD refrain from usage of the declared scheme. Default value is false.
	Deprecated bool `json:"deprecated,omitempty" mapstructure:"deprecated" yaml:"deprecated,omitempty"`
}

var _ SecuritySchemer = (*RelixyOAuth2Config)(nil)

// NewOAuth2Config creates a new OAuth2Config instance.
func NewOAuth2Config(flows map[OAuthFlowType]OAuthFlow) *RelixyOAuth2Config {
	return &RelixyOAuth2Config{
		Type:  OAuth2Scheme,
		Flows: flows,
	}
}

// GetType get the type of security scheme.
func (ss RelixyOAuth2Config) GetType() SecuritySchemeType {
	return ss.Type
}

// Validate if the current instance is valid.
func (ss RelixyOAuth2Config) Validate() error {
	if len(ss.Flows) == 0 {
		return ErrOAuth2FlowRequired
	}

	for key, flow := range ss.Flows {
		err := flow.Validate(key)
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
	}

	return nil
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyOAuth2Config) JSONSchemaExtend(schema *jsonschema.Schema) {
	oauthFlowRef := &jsonschema.Schema{
		Description: "The flow configurations for OAuth 2.0 authentication",
		Ref:         "#/$defs/OAuthFlow",
	}
	implicitFlow := wk8orderedmap.New[string, *jsonschema.Schema]()
	implicitFlow.Set(string(ImplicitFlow), oauthFlowRef)

	passwordFlow := wk8orderedmap.New[string, *jsonschema.Schema]()
	passwordFlow.Set(string(PasswordFlow), oauthFlowRef)

	ccFlow := wk8orderedmap.New[string, *jsonschema.Schema]()
	ccFlow.Set(string(ClientCredentialsFlow), oauthFlowRef)

	acFlow := wk8orderedmap.New[string, *jsonschema.Schema]()
	acFlow.Set(string(AuthorizationCodeFlow), oauthFlowRef)

	daFlow := wk8orderedmap.New[string, *jsonschema.Schema]()
	daFlow.Set(string(DeviceAuthorizationFlow), oauthFlowRef)

	schema.Properties.Set("flows", &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "object",
				Title:       "OAuthPasswordFlow",
				Description: "Configuration for the OAuth Resource Owner Password flow",
				Required:    []string{string(PasswordFlow)},
				Properties:  passwordFlow,
			},
			{
				Type:        "object",
				Title:       "OAuthImplicitFlow",
				Description: "Configuration for the OAuth Implicit flow",
				Required:    []string{string(ImplicitFlow)},
				Properties:  implicitFlow,
			},
			{
				Type:        "object",
				Title:       "OAuthClientCredentialsFlow",
				Description: "Configuration for the OAuth Client Credentials flow. Previously called application in OpenAPI 2.0.",
				Required: []string{
					string(ClientCredentialsFlow),
					"clientId",
					"clientSecret",
					"tokenUrl",
				},
				Properties: ccFlow,
			},
			{
				Type:        "object",
				Title:       "OAuthAuthorizationCodeFlow",
				Description: "Configuration for the OAuth Authorization Code flow. Previously called accessCode in OpenAPI 2.0.",
				Required:    []string{string(AuthorizationCodeFlow)},
				Properties:  acFlow,
			},
			{
				Type:        "object",
				Title:       "OAuthDeviceAuthorizationFlow",
				Description: "Configuration for the OAuth Device Authorization flow.",
				Required:    []string{string(DeviceAuthorizationFlow), "deviceAuthorizationUrl"},
				Properties:  daFlow,
			},
		},
	})
}

// RelixyOpenIDConnectConfig contains configurations for [OpenID Connect] API specification
//
// [OpenID Connect]: https://swagger.io/docs/specification/authentication/openid-connect-discovery
type RelixyOpenIDConnectConfig struct {
	// The type of the security scheme that is always 'openIdConnect'.
	Type SecuritySchemeType `json:"type" mapstructure:"type" yaml:"type" jsonschema:"enum=openIdConnect"`
	// Well-known URL to discover the OpenID-Connect-Discovery provider metadata.
	OpenIDConnectURL string `json:"openIdConnectUrl" mapstructure:"openIdConnectUrl" yaml:"openIdConnectUrl"`
	// Declares this security scheme to be deprecated.
	// Consumers SHOULD refrain from usage of the declared scheme. Default value is false.
	Deprecated bool `json:"deprecated,omitempty" mapstructure:"deprecated" yaml:"deprecated,omitempty"`
}

var _ SecuritySchemer = (*RelixyOpenIDConnectConfig)(nil)

// NewOpenIDConnectConfig creates a new OpenIDConnectConfig instance.
func NewOpenIDConnectConfig(oidcURL string) *RelixyOpenIDConnectConfig {
	return &RelixyOpenIDConnectConfig{
		Type:             OpenIDConnectScheme,
		OpenIDConnectURL: oidcURL,
	}
}

// GetType get the type of security scheme.
func (RelixyOpenIDConnectConfig) GetType() SecuritySchemeType {
	return OpenIDConnectScheme
}

// Validate if the current instance is valid.
func (ss RelixyOpenIDConnectConfig) Validate() error {
	if ss.OpenIDConnectURL == "" {
		return fmt.Errorf("openIdConnectUrl %w for oidc security", ErrFieldRequired)
	}

	_, err := goutils.ParseRelativeOrHTTPURL(ss.OpenIDConnectURL)
	if err != nil {
		return fmt.Errorf("openIdConnectUrl: %w", err)
	}

	return nil
}

// RelixyCookieAuthConfig represents a cookie authentication configuration.
type RelixyCookieAuthConfig struct {
	// The type of the security scheme that is always 'cookie'.
	Type SecuritySchemeType `json:"type" yaml:"type" jsonschema:"enum=cookie"`
	// Declares this security scheme to be deprecated.
	// Consumers SHOULD refrain from usage of the declared scheme. Default value is false.
	Deprecated bool `json:"deprecated,omitempty" mapstructure:"deprecated" yaml:"deprecated,omitempty"`
}

var _ SecuritySchemer = (*RelixyCookieAuthConfig)(nil)

// NewCookieAuthConfig creates a new CookieAuthConfig instance.
func NewCookieAuthConfig() *RelixyCookieAuthConfig {
	return &RelixyCookieAuthConfig{
		Type: CookieAuthScheme,
	}
}

// GetType get the type of security scheme.
func (RelixyCookieAuthConfig) GetType() SecuritySchemeType {
	return CookieAuthScheme
}

// Validate if the current instance is valid.
func (RelixyCookieAuthConfig) Validate() error {
	return nil
}

// RelixyMutualTLSAuthConfig represents a mutualTLS authentication configuration.
type RelixyMutualTLSAuthConfig struct {
	// The type of the security scheme that is always 'mutualTLS'.
	Type SecuritySchemeType `json:"type" yaml:"type" jsonschema:"enum=mutualTLS"`
	// Declares this security scheme to be deprecated.
	// Consumers SHOULD refrain from usage of the declared scheme. Default value is false.
	Deprecated bool `json:"deprecated,omitempty" mapstructure:"deprecated" yaml:"deprecated,omitempty"`
}

var _ SecuritySchemer = (*RelixyMutualTLSAuthConfig)(nil)

// NewMutualTLSAuthConfig creates a new MutualTLSAuthConfig instance.
func NewMutualTLSAuthConfig() *RelixyMutualTLSAuthConfig {
	return &RelixyMutualTLSAuthConfig{
		Type: MutualTLSScheme,
	}
}

// GetType get the type of security scheme.
func (RelixyMutualTLSAuthConfig) GetType() SecuritySchemeType {
	return MutualTLSScheme
}

// Validate if the current instance is valid.
func (RelixyMutualTLSAuthConfig) Validate() error {
	return nil
}
