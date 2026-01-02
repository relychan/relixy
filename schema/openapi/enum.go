package openapi

import (
	"fmt"
	"slices"

	"github.com/invopop/jsonschema"
	"github.com/relychan/goutils"
)

const (
	HTTPAuthSchemeBearer = "bearer"
	AuthorizationHeader  = "Authorization"
)

// SecuritySchemeType represents the authentication scheme enum.
type SecuritySchemeType string

const (
	APIKeyScheme        SecuritySchemeType = "apiKey"
	BasicAuthScheme     SecuritySchemeType = "basic"
	CookieAuthScheme    SecuritySchemeType = "cookie"
	HTTPAuthScheme      SecuritySchemeType = "http"
	OAuth2Scheme        SecuritySchemeType = "oauth2"
	OpenIDConnectScheme SecuritySchemeType = "openIdConnect"
	MutualTLSScheme     SecuritySchemeType = "mutualTLS"
)

var enumSecuritySchemes = []SecuritySchemeType{
	APIKeyScheme,
	HTTPAuthScheme,
	BasicAuthScheme,
	CookieAuthScheme,
	OAuth2Scheme,
	OpenIDConnectScheme,
	MutualTLSScheme,
}

// SupportedSecuritySchemeTypes return all supported security scheme types.
func SupportedSecuritySchemeTypes() []SecuritySchemeType {
	return enumSecuritySchemes
}

// ParseSecuritySchemeType parses SecurityScheme from string.
func ParseSecuritySchemeType(value string) (SecuritySchemeType, error) {
	result := SecuritySchemeType(value)
	if !slices.Contains(enumSecuritySchemes, result) {
		return result, fmt.Errorf(
			"%w; got: %s",
			ErrInvalidSecuritySchemeType,
			value,
		)
	}

	return result, nil
}

// OAuthFlowType represents the OAuth flow type enum.
type OAuthFlowType string

const (
	// AuthorizationCodeFlow is the constant string for the OAuth Authorization Code flow.
	AuthorizationCodeFlow OAuthFlowType = "authorizationCode"
	// ImplicitFlow is the constant string for the OAuth Implicit flow.
	ImplicitFlow OAuthFlowType = "implicit"
	// PasswordFlow is the constant string for the OAuth Resource Owner Password flow.
	PasswordFlow OAuthFlowType = "password"
	// ClientCredentialsFlow is the constant string for the OAuth Client Credentials flow.
	ClientCredentialsFlow OAuthFlowType = "clientCredentials"
	// DeviceAuthorizationFlow is the constant string for the OAuth Device Authorization flow.
	DeviceAuthorizationFlow OAuthFlowType = "deviceAuthorization"
)

var enumValuesOAuthFlows = []OAuthFlowType{
	AuthorizationCodeFlow,
	ImplicitFlow,
	PasswordFlow,
	ClientCredentialsFlow,
	DeviceAuthorizationFlow,
}

// ParseOAuthFlowType parses OAuthFlowType from string.
func ParseOAuthFlowType(value string) (OAuthFlowType, error) {
	result := OAuthFlowType(value)
	if !slices.Contains(enumValuesOAuthFlows, result) {
		return result, fmt.Errorf(
			"%w; got <%s>",
			ErrInvalidOAuthFlowType,
			value,
		)
	}

	return result, nil
}

// ParameterLocation is [the location] of the parameter.
// Possible values are "query", "header", "path" or "cookie".
//
// [the location]: https://swagger.io/specification/#parameter-object
type ParameterLocation string

const (
	InQuery    ParameterLocation = "query"
	InHeader   ParameterLocation = "header"
	InPath     ParameterLocation = "path"
	InCookie   ParameterLocation = "cookie"
	InBody     ParameterLocation = "body"
	InFormData ParameterLocation = "formData"
)

var enumValueParameterLocations = []ParameterLocation{
	InQuery,
	InHeader,
	InPath,
	InCookie,
	InBody,
	InFormData,
}

// IsValid checks if the style enum is valid.
func (j ParameterLocation) IsValid() bool {
	return slices.Contains(enumValueParameterLocations, j)
}

// JSONSchema defines a custom definition for JSON schema.
func (ParameterLocation) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Description: "Location of the parameter which is defined in OpenAPI specifications",
		Enum:        goutils.ToAnySlice(SupportedParameterLocations()),
	}
}

// SupportedParameterLocations returns supported parameter locations.
func SupportedParameterLocations() []ParameterLocation {
	return enumValueParameterLocations
}

// ParseParameterLocation parses ParameterLocation from string.
func ParseParameterLocation(input string) (ParameterLocation, error) {
	result := ParameterLocation(input)
	if !result.IsValid() {
		return result, fmt.Errorf(
			"%w; got: %s",
			ErrInvalidParameterLocation,
			input,
		)
	}

	return result, nil
}
