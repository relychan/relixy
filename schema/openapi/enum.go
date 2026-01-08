package openapi

import (
	"fmt"
	"slices"
)

const (
	// XRelyURLEnv is the extension name enum of the server URL.
	XRelyURLEnv = "x-rely-url-env"
	// XRelyServerWeight is the extension name enum for the weight of server if the load balancer is configured.
	XRelyServerWeight = "x-rely-server-weight"
	// XRelyServerHeaders is the extension name enum for custom headers for the server.
	XRelyServerHeaders = "x-rely-server-headers"
	// XRelyServerSecuritySchemes is the extension name enum for server security schemes.
	XRelyServerSecuritySchemes = "x-rely-server-security-schemes"
	// XRelyServerSecurity is the extension name enum for a server security.
	XRelyServerSecurity = "x-rely-server-security"
	// XRelyServerTLS is the extension name enum for a server TLS config.
	XRelyServerTLS = "x-rely-server-tls"
	// XRelyProxyAction is the extension name enum for a proxy action.
	XRelyProxyAction = "x-rely-proxy-action"
	// XRelySecurityCredentials is the extension name enum for security credentials.
	XRelySecurityCredentials = "x-rely-security-credentials"
	// XRelyOAuth2TokenURLEnv is the extension name enum of a custom environment variable for OAuth2 token URL.
	XRelyOAuth2TokenURLEnv = "x-rely-oauth2-token-url-env" //nolint:gosec
	// XRelyOAuth2RefreshURLEnv is the extension name enum of a custom environment variable for OAuth2 refresh URL.
	XRelyOAuth2RefreshURLEnv = "x-rely-oauth2-refresh-url-env"
)

const (
	HTTPAuthSchemeBearer = "bearer"
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

// ParameterLocation is [the location] of the parameter.
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
