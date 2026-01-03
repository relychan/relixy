package openapi

import (
	"fmt"
	"slices"
)

const (
	// XRelyURLEnv is the extension name enum of the server URL.
	XRelyURLEnv                = "x-rely-url-env"
	XRelyServerWeight          = "x-rely-server-weight"
	XRelyServerHeaders         = "x-rely-server-headers"
	XRelyServerSecuritySchemes = "x-rely-server-security-schemes"
	XRelyServerSecurity        = "x-rely-server-security"
	XRelyServerTLS             = "x-rely-server-tls"
	XRelyProxyAction           = "x-rely-proxy-action"
	// XRelySecurityCredentials is the extension name enum of security credentials.
	XRelySecurityCredentials = "x-rely-security-credentials"
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
	// InQuery is the constant enum that indicates the parameter location in query.
	InQuery = "query"
	// InHeader is the constant enum that indicates the parameter location in header.
	InHeader = "header"
	// InPath is the constant enum that indicates the parameter location in path.
	InPath = "path"
	// InCookie is the constant enum that indicates the parameter location in cookie.
	InCookie = "cookie"
	// InBody is the constant enum that indicates the parameter location in body.
	InBody = "body"
	// InFormData is the constant enum that indicates the parameter location in formData.
	InFormData = "formData"
)

var enumValueParameterLocations = []string{
	InQuery,
	InHeader,
	InPath,
	InCookie,
	InBody,
	InFormData,
}

// IsParameterLocationValid checks if the input string is a valid parameter location.
func IsParameterLocationValid(input string) bool {
	return slices.Contains(enumValueParameterLocations, input)
}
