package base_schema

import (
	"github.com/invopop/jsonschema"
	"github.com/relychan/goutils"
)

// RelixyActionType represents enums of proxy types.
type RelixyActionType string

const (
	ProxyTypeGraphQL RelixyActionType = "graphql"
	ProxyTypeREST    RelixyActionType = "rest"
)

// PrimitiveType represents primitive data types.
type PrimitiveType string

const (
	// String represents text data types.
	String PrimitiveType = "string"
	// Number represents floating-point number data types.
	Number PrimitiveType = "number"
	// Integer represents integer number data types.
	Integer PrimitiveType = "integer"
	// Boolean represents boolean number data types.
	Boolean PrimitiveType = "boolean"
	// Array represents array data types.
	Array PrimitiveType = "array"
	// Object represents object data types.
	Object PrimitiveType = "object"
	// Null represents the null data type.
	Null PrimitiveType = "null"
)

var enumValuePrimitiveTypes = []PrimitiveType{
	String,
	Number,
	Integer,
	Boolean,
	Array,
	Object,
	Null,
}

// JSONSchema defines a custom definition for JSON schema.
func (PrimitiveType) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Description: "A primitive type in OpenAPI specification",
		Enum:        goutils.ToAnySlice(enumValuePrimitiveTypes),
	}
}
