// Package gqlschema defines schemas for GraphQL resources
package gqlschema

import "github.com/relychan/relixy/schema/base_schema"

// RelixyGraphQLConfig contains configurations for GraphQL proxy.
type RelixyGraphQLConfig struct {
	// ScalarTypeMapping configures the custom type mapping between GraphQL scalar types and primitive types for data conversion.
	// Default scalar types are supported. Other types which aren't configured will be forwarded directly
	// from the request parameters and body without serialization.
	ScalarTypeMapping map[string]base_schema.PrimitiveType `json:"scalarTypeMapping,omitempty" yaml:"scalarTypeMapping,omitempty"`
}
