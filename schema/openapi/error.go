package openapi

import (
	"errors"
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
	// ErrResourceSpecRequired occurs when the spec field of resource is empty.
	ErrResourceSpecRequired = errors.New("spec is required in resource")
	// ErrInvalidOpenAPIResourceDefinitionYAML occurs when failing to parse a RelixyOpenAPIResourceDefinition from YAML string.
	ErrInvalidOpenAPIResourceDefinitionYAML = errors.New(
		"failed to parse RelixyOpenAPIResourceDefinition from YAML",
	)
)
