package openapi

import (
	"fmt"
)

var (
	// ErrInvalidParameterLocation occurs when the parameter location is invalid.
	ErrInvalidParameterLocation = fmt.Errorf(
		"invalid ParameterLocation. Expected one of %v",
		enumValueParameterLocations,
	)
	// ErrInvalidSecuritySchemeType occurs when the security scheme type is invalid.
	ErrInvalidSecuritySchemeType = fmt.Errorf(
		"invalid SecuritySchemeType. Expected one of %v",
		SupportedSecuritySchemeTypes(),
	)
)
