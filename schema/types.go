package schema

import (
	"github.com/invopop/jsonschema"
	orderedmap "github.com/pb33f/ordered-map/v2"
)

// Parameter represents a high-level OpenAPI 3+ Parameter object, that is backed by a low-level one.
// A unique parameter is defined by a combination of a name and location.
// https://spec.openapis.org/oas/v3.1.0#parameter-object
type Parameter struct {
	Name            string                                              `json:"name,omitempty" yaml:"name,omitempty"`
	In              ParameterLocation                                   `json:"in,omitempty" yaml:"in,omitempty"`
	Description     string                                              `json:"description,omitempty" yaml:"description,omitempty"`
	Required        *bool                                               `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool                                                `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                                                `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string                                              `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         *bool                                               `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool                                                `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *RelyProxySchema                                    `json:"schema,omitempty" yaml:"schema,omitempty"`
	Content         *orderedmap.OrderedMap[string, *RelyProxyMediaType] `json:"content,omitempty" yaml:"content,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (Parameter) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("content", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxyMediaType",
			},
		})
}

// RelyProxyRequestBody represents a high-level OpenAPI 3+ RequestBody object, backed by a low-level one.
//   - https://spec.openapis.org/oas/v3.1.0#request-body-object
type RelyProxyRequestBody struct {
	Description string                                              `json:"description,omitempty" yaml:"description,omitempty"`
	Content     *orderedmap.OrderedMap[string, *RelyProxyMediaType] `json:"content,omitempty" yaml:"content,omitempty"`
	Required    *bool                                               `json:"required,omitempty" yaml:"required,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelyProxyRequestBody) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("content", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxyMediaType",
			},
		})
}

// RelyProxyMediaType represents a high-level OpenAPI MediaType object that is backed by a low-level one.
//
// Each Media Type Object provides schema and examples for the media type identified by its key.
//   - https://spec.openapis.org/oas/v3.1.0#media-type-object
type RelyProxyMediaType struct {
	Schema     *RelyProxySchema `json:"schema,omitempty" yaml:"schema,omitempty"`
	ItemSchema *RelyProxySchema `json:"itemSchema,omitempty" yaml:"itemSchema,omitempty"`
	// Example      *yaml.Node                             `json:"example,omitempty" yaml:"example,omitempty"`
	// Examples     *orderedmap.Map[string, *base.Example] `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding     *orderedmap.OrderedMap[string, *RelyProxyEncoding] `json:"encoding,omitempty" yaml:"encoding,omitempty"`
	ItemEncoding *orderedmap.OrderedMap[string, *RelyProxyEncoding] `json:"itemEncoding,omitempty" yaml:"itemEncoding,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelyProxyMediaType) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.Set("encoding", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/RelyProxyEncoding",
		},
	})
	schema.Properties.
		Set("itemEncoding", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxyEncoding",
			},
		})
}

// RelyProxyEncoding represents an OpenAPI 3+ Encoding object
//   - https://spec.openapis.org/oas/v3.1.0#encoding-object
type RelyProxyEncoding struct {
	ContentType   string                                           `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers       *orderedmap.OrderedMap[string, *RelyProxyHeader] `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string                                           `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       *bool                                            `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool                                             `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelyProxyEncoding) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.Set("headers", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/RelyProxyHeader",
		},
	})
}

// RelyProxyHeader represents a high-level OpenAPI 3+ Header object that is backed by a low-level one.
//   - https://spec.openapis.org/oas/v3.1.0#header-object
type RelyProxyHeader struct {
	Description     string                                              `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                                                `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool                                                `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                                                `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string                                              `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         bool                                                `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool                                                `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *RelyProxySchema                                    `json:"schema,omitempty" yaml:"schema,omitempty"`
	Content         *orderedmap.OrderedMap[string, *RelyProxyMediaType] `json:"content,omitempty" yaml:"content,omitempty"`
	// Example         *yaml.Node                                 `json:"example,omitempty" yaml:"example,omitempty"`
	// Examples        *orderedmap.Map[string, *highbase.Example] `json:"examples,omitempty" yaml:"examples,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelyProxyHeader) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("content", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxyMediaType",
			},
		})
}

// Discriminator is only used by OpenAPI 3+ documents, it represents a polymorphic discriminator used for schemas
// When request bodies or response payloads may be one of a number of different schemas, a discriminator object can be used to aid in serialization, deserialization, and validation.
// The discriminator is a specific object in a schema which is used to inform the consumer of the document of an alternative schema based on the value associated with it.
// When using the discriminator, inline schemas will not be considered.
type Discriminator struct {
	PropertyName   string                                 `json:"propertyName,omitempty" yaml:"propertyName,omitempty"`
	Mapping        *orderedmap.OrderedMap[string, string] `json:"mapping,omitempty" yaml:"mapping,omitempty"`
	DefaultMapping string                                 `json:"defaultMapping,omitempty" yaml:"defaultMapping,omitempty"` // OpenAPI 3.2+ defaultMapping for fallback schema
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (Discriminator) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("mapping", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Type: "string",
			},
		})
}

// RelyProxySchema represents a JSON Schema that support Swagger, OpenAPI 3 and OpenAPI 3.1
//
// Until 3.1 OpenAPI had a strange relationship with JSON Schema. It's been a super-set/sub-set
// mix, which has been confusing. So, instead of building a bunch of different models, we have compressed
// all variations into a single model that makes it easy to support multiple spec types.
//
//   - v2 schema: https://swagger.io/specification/v2/#schemaObject
//   - v3 schema: https://swagger.io/specification/#schema-object
//   - v3.1 schema: https://spec.openapis.org/oas/v3.1.0#schema-object
type RelyProxySchema struct {
	// 3.1 only, used to define a dialect for this schema, label is '$schema'.
	SchemaTypeRef string `json:"$schema,omitempty" yaml:"$schema,omitempty"`

	// Used to define a reference for this schema.
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`

	// In versions 2 and 3.0, this ExclusiveMaximum can only be a boolean.
	// In version 3.1, ExclusiveMaximum is a number.
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`

	// In versions 2 and 3.0, this ExclusiveMinimum can only be a boolean.
	// In version 3.1, ExclusiveMinimum is a number.
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`

	// In versions 2 and 3.0, this Type is a single value, so array will only ever have one value
	// in version 3.1, Type can be multiple values
	Type []string `json:"type,omitempty" yaml:"type,omitempty"`

	// Schemas are resolved on demand using a SchemaProxy
	AllOf []RelyProxySchema `json:"allOf,omitempty" yaml:"allOf,omitempty"`

	// Polymorphic Schemas are only available in version 3+
	OneOf         []RelyProxySchema `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf         []RelyProxySchema `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Discriminator *Discriminator    `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`

	// in 3.1 examples can be an array (which is recommended)
	// Examples []*yaml.Node `json:"examples,omitempty" yaml:"examples,omitempty"`

	// in 3.1 prefixItems provides tuple validation support.
	PrefixItems []RelyProxySchema `json:"prefixItems,omitempty" yaml:"prefixItems,omitempty"`

	// 3.1 Specific properties
	Contains          *RelyProxySchema                                 `json:"contains,omitempty" yaml:"contains,omitempty"`
	MinContains       *int64                                           `json:"minContains,omitempty" yaml:"minContains,omitempty"`
	MaxContains       *int64                                           `json:"maxContains,omitempty" yaml:"maxContains,omitempty"`
	If                *RelyProxySchema                                 `json:"if,omitempty" yaml:"if,omitempty"`
	Else              *RelyProxySchema                                 `json:"else,omitempty" yaml:"else,omitempty"`
	Then              *RelyProxySchema                                 `json:"then,omitempty" yaml:"then,omitempty"`
	DependentSchemas  *orderedmap.OrderedMap[string, *RelyProxySchema] `json:"dependentSchemas,omitempty" yaml:"dependentSchemas,omitempty"`
	DependentRequired *orderedmap.OrderedMap[string, []string]         `json:"dependentRequired,omitempty" yaml:"dependentRequired,omitempty"`
	PatternProperties *orderedmap.OrderedMap[string, *RelyProxySchema] `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
	PropertyNames     *RelyProxySchema                                 `json:"propertyNames,omitempty" yaml:"propertyNames,omitempty"`
	UnevaluatedItems  *RelyProxySchema                                 `json:"unevaluatedItems,omitempty" yaml:"unevaluatedItems,omitempty"`

	// in 3.1 Items can be a Schema or a boolean
	Items *RelyProxySchema `json:"items,omitempty" yaml:"items,omitempty"`

	// 3.1 only, part of the JSON Schema spec provides a way to identify a sub-schema
	Anchor string `json:"$anchor,omitempty" yaml:"$anchor,omitempty"`

	// Compatible with all versions
	Not                  *RelyProxySchema                                 `json:"not,omitempty" yaml:"not,omitempty"`
	Properties           *orderedmap.OrderedMap[string, *RelyProxySchema] `json:"properties,omitempty" yaml:"properties,omitempty"`
	Title                string                                           `json:"title,omitempty" yaml:"title,omitempty"`
	MultipleOf           *float64                                         `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Maximum              *float64                                         `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	Minimum              *float64                                         `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	MaxLength            *int64                                           `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength            *int64                                           `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Pattern              string                                           `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Format               string                                           `json:"format,omitempty" yaml:"format,omitempty"`
	MaxItems             *int64                                           `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems             *int64                                           `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems          *bool                                            `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	MaxProperties        *int64                                           `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        *int64                                           `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required             []string                                         `json:"required,omitempty" yaml:"required,omitempty"`
	Enum                 []any                                            `json:"enum,omitempty" yaml:"enum,omitempty"`
	AdditionalProperties *RelyProxySchema                                 `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Description          string                                           `json:"description,omitempty" yaml:"description,omitempty"`
	Default              any                                              `json:"default,omitempty" yaml:"default,omitempty"`
	Const                any                                              `json:"const,omitempty" yaml:"const,omitempty"`
	Nullable             *bool                                            `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	ReadOnly             *bool                                            `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly            *bool                                            `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	// XML                  *XML                                `json:"xml,omitempty" yaml:"xml,omitempty"`
	// ExternalDocs         *ExternalDoc                        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// Example              *yaml.Node                          `json:"example,omitempty" yaml:"example,omitempty"`
	Deprecated *bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelyProxySchema) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("dependentSchemas", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxySchema",
			},
		})
	schema.Properties.
		Set("patternProperties", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxySchema",
			},
		})
	schema.Properties.
		Set("dependentRequired", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Type: "array",
				Items: &jsonschema.Schema{
					Type: "string",
				},
			},
		})
	schema.Properties.
		Set("properties", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxySchema",
			},
		})
}

// RelyProxyAPIDocumentInfo represents a high-level Info object that provides metadata about the API.
// The metadata MAY be used by the clients if needed, and MAY be presented
// in editing or documentation generation tools for convenience.
type RelyProxyAPIDocumentInfo struct {
	Summary        string                           `json:"summary,omitempty" yaml:"summary,omitempty"`
	Title          string                           `json:"title,omitempty" yaml:"title,omitempty"`
	Description    string                           `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string                           `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *RelyProxyAPIDocumentInfoContact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *RelyProxyAPIDocumentInfoLicense `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string                           `json:"version,omitempty" yaml:"version,omitempty"`
}

// RelyProxyAPIDocumentInfoContact represents a high-level representation of the Contact definitions.
type RelyProxyAPIDocumentInfoContact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// RelyProxyAPIDocumentInfoLicense represents a high-level representation of the License definitions.
type RelyProxyAPIDocumentInfoLicense struct {
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	URL        string `json:"url,omitempty" yaml:"url,omitempty"`
	Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
}
