package openapi

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidSecuritySchemeType occurs when the security scheme type is invalid.
	ErrInvalidSecuritySchemeType = fmt.Errorf(
		"invalid SecuritySchemeType. Expected one of %v",
		SupportedSecuritySchemeTypes(),
	)
	// ErrFieldRequired occurs when a field is required.
	ErrFieldRequired = errors.New("field is required")
	// ErrInvalidOAuthFlowType occurs when the OAuth flow type is invalid.
	ErrInvalidOAuthFlowType = fmt.Errorf(
		"invalid OAuthFlowType. Expected %v",
		enumValuesOAuthFlows,
	)
	// ErrOAuth2FlowRequired occurs when the OAuth flow is null or empty.
	ErrOAuth2FlowRequired = errors.New("require at least 1 flow for oauth2 security")
)

var errSecuritySchemerRequired = errors.New("SecuritySchemer is required")
