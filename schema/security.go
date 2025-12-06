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

// RelyProxySecurityScheme contains authentication configurations.
// The schema follows [OpenAPI 3] specification
//
// [OpenAPI 3]: https://swagger.io/docs/specification/authentication
type RelyProxySecurityScheme struct {
	SecuritySchemer
}

type rawSecurityScheme struct {
	Type SecuritySchemeType `json:"type" yaml:"type"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelyProxySecurityScheme) UnmarshalJSON(b []byte) error {
	var rawScheme rawSecurityScheme

	err := json.Unmarshal(b, &rawScheme)
	if err != nil {
		return err
	}

	switch rawScheme.Type {
	case APIKeyScheme:
		var config APIKeyAuthConfig

		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case BasicAuthScheme:
		var config BasicAuthConfig

		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case HTTPAuthScheme:
		var config HTTPAuthConfig

		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case OAuth2Scheme:
		var config OAuth2Config

		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case OpenIDConnectScheme:
		var config OpenIDConnectConfig

		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case CookieAuthScheme:
		j.SecuritySchemer = &CookieAuthConfig{
			Type: rawScheme.Type,
		}
	case MutualTLSScheme:
		j.SecuritySchemer = &MutualTLSAuthConfig{
			Type: rawScheme.Type,
		}
	default:
		return fmt.Errorf("%w; got: %s", ErrInvalidSecuritySchemeType, rawScheme.Type)
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (j RelyProxySecurityScheme) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.SecuritySchemer)
}

// MarshalYAML implements the yaml.Marshaler interface.
func (j RelyProxySecurityScheme) MarshalYAML() (any, error) {
	return yaml.Marshal(j.SecuritySchemer)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (j *RelyProxySecurityScheme) UnmarshalYAML(value *yaml.Node) error {
	var rawScheme rawSecurityScheme

	err := value.Decode(&rawScheme)
	if err != nil {
		return err
	}

	switch rawScheme.Type {
	case APIKeyScheme:
		var config APIKeyAuthConfig

		err := value.Decode(&config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case BasicAuthScheme:
		var config BasicAuthConfig

		err := value.Decode(&config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case HTTPAuthScheme:
		var config HTTPAuthConfig

		err := value.Decode(&config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case OAuth2Scheme:
		var config OAuth2Config

		err := value.Decode(&config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case OpenIDConnectScheme:
		var config OpenIDConnectConfig

		err := value.Decode(&config)
		if err != nil {
			return err
		}

		j.SecuritySchemer = &config
	case CookieAuthScheme:
		j.SecuritySchemer = &CookieAuthConfig{
			Type: rawScheme.Type,
		}
	case MutualTLSScheme:
		j.SecuritySchemer = &MutualTLSAuthConfig{
			Type: rawScheme.Type,
		}
	default:
		return fmt.Errorf("%w; got: %s", ErrInvalidSecuritySchemeType, rawScheme.Type)
	}

	return nil
}

// Validate if the current instance is valid.
func (j *RelyProxySecurityScheme) Validate() error {
	if j.SecuritySchemer == nil {
		return errSecuritySchemerRequired
	}

	return j.SecuritySchemer.Validate()
}

// JSONSchema defines a custom definition for JSON schema.
func (RelyProxySecurityScheme) JSONSchema() *jsonschema.Schema {
	envStringRef := &jsonschema.Schema{
		Ref: "#/$defs/EnvString",
	}
	apiKeySchema := wk8orderedmap.New[string, *jsonschema.Schema]()
	apiKeySchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{APIKeyScheme},
	})
	apiKeySchema.Set("value", envStringRef)
	apiKeySchema.Set("in", &jsonschema.Schema{
		Type: "string",
		Enum: goutils.ToAnySlice(authscheme.GetSupportedAuthLocations()),
	})
	apiKeySchema.Set("name", &jsonschema.Schema{
		Type: "string",
	})

	httpAuthSchema := wk8orderedmap.New[string, *jsonschema.Schema]()
	httpAuthSchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{HTTPAuthScheme},
	})
	httpAuthSchema.Set("value", envStringRef)
	httpAuthSchema.Set("header", &jsonschema.Schema{
		Type: "string",
	})
	httpAuthSchema.Set("scheme", &jsonschema.Schema{
		Type: "string",
	})

	basicAuthSchema := wk8orderedmap.New[string, *jsonschema.Schema]()
	basicAuthSchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{BasicAuthScheme},
	})
	basicAuthSchema.Set("username", envStringRef)
	basicAuthSchema.Set("password", envStringRef)
	httpAuthSchema.Set("header", &jsonschema.Schema{
		Description: "Request contains a header field in the form of Authorization: Basic [credentials]",
		OneOf: []*jsonschema.Schema{
			{Type: "string"},
			{Type: "null"},
		},
	})

	oidcSchema := wk8orderedmap.New[string, *jsonschema.Schema]()
	oidcSchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{OpenIDConnectScheme},
	})
	oidcSchema.Set("openIdConnectUrl", &jsonschema.Schema{
		Type: "string",
	})

	cookieSchema := wk8orderedmap.New[string, *jsonschema.Schema]()
	cookieSchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{CookieAuthScheme},
	})

	mutualTLSSchema := wk8orderedmap.New[string, *jsonschema.Schema]()
	mutualTLSSchema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{MutualTLSScheme},
	})

	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:       "object",
				Required:   []string{"type", "value", "in", "name"},
				Properties: apiKeySchema,
			},
			{
				Type:       "object",
				Properties: basicAuthSchema,
				Required:   []string{"type", "username", "password"},
			},
			{
				Type:       "object",
				Properties: httpAuthSchema,
				Required:   []string{"type", "value", "header", "scheme"},
			},
			newOAuth2ConfigSchema(),
			{
				Type:       "object",
				Properties: oidcSchema,
				Required:   []string{"type", "openIdConnectUrl"},
			},
			{
				Type:       "object",
				Properties: cookieSchema,
				Required:   []string{"type"},
			},
			{
				Type:       "object",
				Properties: mutualTLSSchema,
				Required:   []string{"type"},
			},
		},
	}
}

// APIKeyAuthConfig contains configurations for [apiKey authentication]
//
// [apiKey authentication]: https://swagger.io/docs/specification/authentication/api-keys/
type APIKeyAuthConfig struct {
	Type  SecuritySchemeType      `json:"type" yaml:"type"`
	In    authscheme.AuthLocation `json:"in" yaml:"in"`
	Name  string                  `json:"name" yaml:"name"`
	Value goenvconf.EnvString     `json:"value" yaml:"value"`
}

var _ SecuritySchemer = (*APIKeyAuthConfig)(nil)

// NewAPIKeyAuthConfig creates a new APIKeyAuthConfig instance.
func NewAPIKeyAuthConfig(
	name string,
	in authscheme.AuthLocation,
	value goenvconf.EnvString,
) *APIKeyAuthConfig {
	return &APIKeyAuthConfig{
		Name:  name,
		In:    in,
		Value: value,
	}
}

// Validate if the current instance is valid.
func (j *APIKeyAuthConfig) Validate() error {
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
func (*APIKeyAuthConfig) GetType() SecuritySchemeType {
	return APIKeyScheme
}

// HTTPAuthConfig contains configurations for http authentication
// If the scheme is [bearer], the authenticator follows OpenAPI 3 specification.
//
// [bearer]: https://swagger.io/docs/specification/authentication/bearer-authentication
type HTTPAuthConfig struct {
	Type   SecuritySchemeType  `json:"type" yaml:"type"`
	Header string              `json:"header" yaml:"header"`
	Scheme string              `json:"scheme" yaml:"scheme"`
	Value  goenvconf.EnvString `json:"value" yaml:"value"`
}

var _ SecuritySchemer = (*HTTPAuthConfig)(nil)

// NewHTTPAuthConfig creates a new HTTPAuthConfig instance.
func NewHTTPAuthConfig(scheme string, header string, value goenvconf.EnvString) *HTTPAuthConfig {
	return &HTTPAuthConfig{
		Type:   HTTPAuthScheme,
		Header: header,
		Scheme: scheme,
		Value:  value,
	}
}

// Validate if the current instance is valid.
func (ss *HTTPAuthConfig) Validate() error {
	if ss.Scheme == "" {
		return fmt.Errorf("scheme %w for http security", ErrFieldRequired)
	}

	return nil
}

// GetType get the type of security scheme.
func (*HTTPAuthConfig) GetType() SecuritySchemeType {
	return HTTPAuthScheme
}

// BasicAuthConfig contains configurations for the [basic] authentication.
//
// [basic]: https://swagger.io/docs/specification/authentication/basic-authentication
type BasicAuthConfig struct {
	Type     SecuritySchemeType  `json:"type" yaml:"type"`
	Header   string              `json:"header" yaml:"header"`
	Username goenvconf.EnvString `json:"username" yaml:"username"`
	Password goenvconf.EnvString `json:"password" yaml:"password"`
}

// NewBasicAuthConfig creates a new BasicAuthConfig instance.
func NewBasicAuthConfig(username, password goenvconf.EnvString) *BasicAuthConfig {
	return &BasicAuthConfig{
		Type:     BasicAuthScheme,
		Username: username,
		Password: password,
	}
}

// Validate if the current instance is valid.
func (*BasicAuthConfig) Validate() error {
	return nil
}

// GetType get the type of security scheme.
func (*BasicAuthConfig) GetType() SecuritySchemeType {
	return BasicAuthScheme
}

// OAuthFlow contains flow configurations for [OAuth 2.0] API specification
//
// [OAuth 2.0]: https://swagger.io/docs/specification/authentication/oauth2
type OAuthFlow struct {
	AuthorizationURL string                         `json:"authorizationUrl,omitempty" mapstructure:"authorizationUrl" yaml:"authorizationUrl,omitempty"`
	TokenURL         *goenvconf.EnvString           `json:"tokenUrl,omitempty"         mapstructure:"tokenUrl"         yaml:"tokenUrl,omitempty"`
	RefreshURL       string                         `json:"refreshUrl,omitempty"       mapstructure:"refreshUrl"       yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string              `json:"scopes,omitempty"           mapstructure:"scopes"           yaml:"scopes,omitempty"`
	ClientID         *goenvconf.EnvString           `json:"clientId,omitempty"         mapstructure:"clientId"         yaml:"clientId,omitempty"`
	ClientSecret     *goenvconf.EnvString           `json:"clientSecret,omitempty"     mapstructure:"clientSecret"     yaml:"clientSecret,omitempty"`
	EndpointParams   map[string]goenvconf.EnvString `json:"endpointParams,omitempty"   mapstructure:"endpointParams"   yaml:"endpointParams,omitempty"`
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

// OAuth2Config contains configurations for [OAuth 2.0] API specification
//
// [OAuth 2.0]: https://swagger.io/docs/specification/authentication/oauth2
type OAuth2Config struct {
	Type  SecuritySchemeType          `json:"type"  mapstructure:"type"  yaml:"type"`
	Flows map[OAuthFlowType]OAuthFlow `json:"flows" mapstructure:"flows" yaml:"flows"`
}

var _ SecuritySchemer = (*OAuth2Config)(nil)

// NewOAuth2Config creates a new OAuth2Config instance.
func NewOAuth2Config(flows map[OAuthFlowType]OAuthFlow) *OAuth2Config {
	return &OAuth2Config{
		Type:  OAuth2Scheme,
		Flows: flows,
	}
}

// GetType get the type of security scheme.
func (ss OAuth2Config) GetType() SecuritySchemeType {
	return ss.Type
}

// Validate if the current instance is valid.
func (ss OAuth2Config) Validate() error {
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

// OpenIDConnectConfig contains configurations for [OpenID Connect] API specification
//
// [OpenID Connect]: https://swagger.io/docs/specification/authentication/openid-connect-discovery
type OpenIDConnectConfig struct {
	Type             SecuritySchemeType `json:"type"             mapstructure:"type"             yaml:"type"`
	OpenIDConnectURL string             `json:"openIdConnectUrl" mapstructure:"openIdConnectUrl" yaml:"openIdConnectUrl"`
}

var _ SecuritySchemer = (*OpenIDConnectConfig)(nil)

// NewOpenIDConnectConfig creates a new OpenIDConnectConfig instance.
func NewOpenIDConnectConfig(oidcURL string) *OpenIDConnectConfig {
	return &OpenIDConnectConfig{
		Type:             OpenIDConnectScheme,
		OpenIDConnectURL: oidcURL,
	}
}

// GetType get the type of security scheme.
func (ss OpenIDConnectConfig) GetType() SecuritySchemeType {
	return ss.Type
}

// Validate if the current instance is valid.
func (ss OpenIDConnectConfig) Validate() error {
	if ss.OpenIDConnectURL == "" {
		return fmt.Errorf("openIdConnectUrl %w for oidc security", ErrFieldRequired)
	}

	_, err := goutils.ParseRelativeOrHTTPURL(ss.OpenIDConnectURL)
	if err != nil {
		return fmt.Errorf("openIdConnectUrl: %w", err)
	}

	return nil
}

// CookieAuthConfig represents a cookie authentication configuration.
type CookieAuthConfig struct {
	Type SecuritySchemeType `json:"type" yaml:"type"`
}

var _ SecuritySchemer = (*CookieAuthConfig)(nil)

// NewCookieAuthConfig creates a new CookieAuthConfig instance.
func NewCookieAuthConfig() *CookieAuthConfig {
	return &CookieAuthConfig{
		Type: CookieAuthScheme,
	}
}

// GetType get the type of security scheme.
func (CookieAuthConfig) GetType() SecuritySchemeType {
	return CookieAuthScheme
}

// Validate if the current instance is valid.
func (CookieAuthConfig) Validate() error {
	return nil
}

// MutualTLSAuthConfig represents a mutualTLS authentication configuration.
type MutualTLSAuthConfig struct {
	Type SecuritySchemeType `json:"type" yaml:"type"`
}

var _ SecuritySchemer = (*MutualTLSAuthConfig)(nil)

// NewMutualTLSAuthConfig creates a new MutualTLSAuthConfig instance.
func NewMutualTLSAuthConfig() *MutualTLSAuthConfig {
	return &MutualTLSAuthConfig{
		Type: MutualTLSScheme,
	}
}

// GetType get the type of security scheme.
func (MutualTLSAuthConfig) GetType() SecuritySchemeType {
	return MutualTLSScheme
}

// Validate if the current instance is valid.
func (MutualTLSAuthConfig) Validate() error {
	return nil
}

func newOAuth2ConfigSchema() *jsonschema.Schema {
	oauth2Schema := wk8orderedmap.New[string, *jsonschema.Schema]()
	oauth2Schema.Set("type", &jsonschema.Schema{
		Type: "string",
		Enum: []any{OAuth2Scheme},
	})

	oauthFlowRef := &jsonschema.Schema{
		Ref: "#/$defs/OAuthFlow",
	}
	implicitFlow := wk8orderedmap.New[string, *jsonschema.Schema]()
	implicitFlow.Set(string(ImplicitFlow), oauthFlowRef)

	passwordFlow := wk8orderedmap.New[string, *jsonschema.Schema]()
	passwordFlow.Set(string(PasswordFlow), oauthFlowRef)

	ccFlow := wk8orderedmap.New[string, *jsonschema.Schema]()
	ccFlow.Set(string(ClientCredentialsFlow), oauthFlowRef)

	acFlow := wk8orderedmap.New[string, *jsonschema.Schema]()
	acFlow.Set(string(AuthorizationCodeFlow), oauthFlowRef)

	oauth2Schema.Set("flows", &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:       "object",
				Required:   []string{string(PasswordFlow)},
				Properties: passwordFlow,
			},
			{
				Type:       "object",
				Required:   []string{string(ImplicitFlow)},
				Properties: implicitFlow,
			},
			{
				Type:       "object",
				Required:   []string{string(ClientCredentialsFlow)},
				Properties: ccFlow,
			},
			{
				Type:       "object",
				Required:   []string{string(AuthorizationCodeFlow)},
				Properties: acFlow,
			},
		},
	})

	return &jsonschema.Schema{
		Type:       "object",
		Properties: oauth2Schema,
		Required:   []string{"type", "flows"},
	}
}
