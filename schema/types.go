package schema

import (
	"github.com/invopop/jsonschema"
	orderedmap "github.com/pb33f/ordered-map/v2"
)

// Parameter represents a high-level OpenAPI 3+ Parameter object, that is backed by a low-level one.
// A unique parameter is defined by a combination of a name and location.
// https://spec.openapis.org/oas/v3.2.0#parameter-object
type Parameter struct {
	// The name of the parameter. Parameter names are case-sensitive.
	Name string `json:"name" yaml:"name"`
	// The location of the parameter. Possible values are query, header, path or cookie.
	In ParameterLocation `json:"in" yaml:"in"`
	// A brief description of the parameter. This could contain examples of use.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Determines whether this parameter is mandatory.
	// If the parameter location is path, this field is REQUIRED and its value MUST be true.
	// Otherwise, the field MAY be included and its default value is false.
	Required *bool `json:"required,omitempty" yaml:"required,omitempty"`
	// Specifies that a parameter is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	// If true, clients MAY pass a zero-length string value in place of parameters that would otherwise be omitted entirely,
	// which the server SHOULD interpret as the parameter being unused.
	// Default value is false. If style is used, and if behavior is n/a (cannot be serialized), the value of allowEmptyValue SHALL be ignored.
	// Interactions between this field and the parameter’s Schema Object are implementation-defined.
	// This field is valid only for query parameters.
	//
	// Deprecated: Use of this field is NOT RECOMMENDED, and it is likely to be removed in a later revision.
	AllowEmptyValue bool `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty" jsonschema:"deprecated"`
	// Describes how the parameter value will be serialized depending on the type of the parameter value.
	// Default values (based on value of in): for query - form; for path - simple; for header - simple;
	// for cookie - form (for compatibility reasons; note that style: cookie SHOULD be used with in: cookie.
	Style string `json:"style,omitempty" yaml:"style,omitempty"`
	// When this is true, parameter values of type array or object generate separate parameters for each value of the array or key-value pair of the map.
	// For other types of parameters, or when style is "deepObject", this field has no effect. When style is "form" or "cookie", the default value is true.
	// For all other styles, the default value is false.
	Explode *bool `json:"explode,omitempty" yaml:"explode,omitempty"`
	// When this is true, parameter values are serialized using reserved expansion, as defined by RFC6570 Section 3.2.3,
	// which allows RFC3986’s reserved character set, as well as percent-encoded triples, to pass through unchanged,
	// while still percent-encoding all other disallowed characters (including % outside of percent-encoded triples).
	// Applications are still responsible for percent-encoding reserved characters that are not allowed by the rules of the in destination or media type,
	// or are not allowed in the path by this specification; see URL Percent-Encoding for details.
	// The default value is false. This field only applies to in and style values that automatically percent-encode.
	AllowReserved bool `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	// The schema defining the type used for the parameter.
	Schema *RelixySchema `json:"schema,omitempty" yaml:"schema,omitempty"`
	// A map containing the representations for the parameter. The key is the media type and the value describes it. The map MUST only contain one entry.
	Content *orderedmap.OrderedMap[string, *RelixyMediaType] `json:"content,omitempty" yaml:"content,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (Parameter) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("content", &jsonschema.Schema{
			Type:        "object",
			Description: "A map containing the representations for the parameter. The key is the media type and the value describes it. The map MUST only contain one entry.",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelixyMediaType",
			},
		})
}

// RelixyRequestBody represents a high-level OpenAPI 3+ RequestBody object, backed by a low-level one.
//   - https://spec.openapis.org/oas/v3.2.0#request-body-object
type RelixyRequestBody struct {
	// Description of the request body.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// The content of the request body. The key is a media type or media type range, and the value describes it.
	// For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content *orderedmap.OrderedMap[string, *RelixyMediaType] `json:"content,omitempty" yaml:"content,omitempty"`
	// Determines if the request body is required in the request. Defaults to false.
	Required *bool `json:"required,omitempty" yaml:"required,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyRequestBody) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("content", &jsonschema.Schema{
			Type:        "object",
			Description: "The content of the request body. The key is a media type or media type range, and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*", //nolint:lll
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelixyMediaType",
			},
		})
}

// RelixyMediaType represents a high-level OpenAPI MediaType object that is backed by a low-level one.
// Each Media Type Object provides schema and examples for the media type identified by its key.
//   - https://spec.openapis.org/oas/v3.2.0#media-type-object
type RelixyMediaType struct {
	// The schema defining the content of the request, response, or parameter.
	Schema *RelixySchema `json:"schema,omitempty" yaml:"schema,omitempty"`
	// A schema describing each item within a sequential media type.
	ItemSchema *RelixySchema `json:"itemSchema,omitempty" yaml:"itemSchema,omitempty"`
	// A map between a property name and its encoding information, as defined under Encoding By Name.
	// The encoding field SHALL only apply when the media type is multipart or application/x-www-form-urlencoded.
	// If no Encoding Object is provided for a property, the behavior is determined by the default values documented for the Encoding Object.
	// This field MUST NOT be present if prefixEncoding or itemEncoding are present.
	Encoding *orderedmap.OrderedMap[string, *RelixyEncoding] `json:"encoding,omitempty" yaml:"encoding,omitempty"`
	// An array of positional encoding information, as defined under Encoding By Position.
	// The prefixEncoding field SHALL only apply when the media type is multipart.
	// If no Encoding Object is provided for a property, the behavior is determined by the default values documented for the Encoding Object.
	// This field MUST NOT be present if encoding is present.
	PrefixEncoding []*RelixyEncoding `json:"prefixEncoding,omitempty" yaml:"prefixEncoding,omitempty"`
	// A single Encoding Object that provides encoding information for multiple array items, as defined under Encoding By Position.
	// The itemEncoding field SHALL only apply when the media type is multipart.
	// If no Encoding Object is provided for a property, the behavior is determined by the default values documented for the Encoding Object.
	// This field MUST NOT be present if encoding is present.
	ItemEncoding *RelixyEncoding `json:"itemEncoding,omitempty" yaml:"itemEncoding,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyMediaType) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.Set("encoding", &jsonschema.Schema{
		Type:        "object",
		Description: "A map between a property name and its encoding information, as defined under Encoding By Name. The encoding field SHALL only apply when the media type is multipart or application/x-www-form-urlencoded. If no Encoding Object is provided for a property, the behavior is determined by the default values documented for the Encoding Object. This field MUST NOT be present if prefixEncoding or itemEncoding are present.", //nolint:lll
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/RelixyEncoding",
		},
	})
}

// RelixyEncoding represents an OpenAPI 3+ Encoding object
//   - https://spec.openapis.org/oas/v3.2.0#encoding-object
type RelixyEncoding struct {
	// The Content-Type for encoding a specific property.
	// The value is a comma-separated list, each element of which is either a specific media type (e.g. image/png) or a wildcard media type (e.g. image/*).
	ContentType string `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	// A map allowing additional information to be provided as headers.
	// Content-Type is described separately and SHALL be ignored in this section.
	// This field SHALL be ignored if the media type is not a multipart.
	Headers *orderedmap.OrderedMap[string, *RelixyHeader] `json:"headers,omitempty" yaml:"headers,omitempty"`
	// Describes how a specific property value will be serialized depending on its type.
	// See Parameter Object for details on the style field.
	// The behavior follows the same values as query parameters,
	// including the default value of "form" which applies only when contentType is not being used due to one or both of explode or allowReserved being explicitly specified.
	// Note that the initial ? used in query strings is not used in application/x-www-form-urlencoded message bodies, and MUST be removed
	// (if using an RFC6570 implementation) or simply not added (if constructing the string manually).
	// This field SHALL be ignored if the media type is not application/x-www-form-urlencoded or multipart/form-data.
	// If a value is explicitly defined, then the value of contentType (implicit or explicit) SHALL be ignored.
	Style string `json:"style,omitempty" yaml:"style,omitempty"`
	// When this is true, property values of type array or object generate separate parameters for each value of the array, or key-value-pair of the map.
	// For other types of properties, or when style is "deepObject", this field has no effect.
	// When style is "form", the default value is true. For all other styles, the default value is false.
	// This field SHALL be ignored if the media type is not application/x-www-form-urlencoded or multipart/form-data.
	// If a value is explicitly defined, then the value of contentType (implicit or explicit) SHALL be ignored.
	Explode *bool `json:"explode,omitempty" yaml:"explode,omitempty"`
	// When this is true, parameter values are serialized using reserved expansion, as defined by RFC6570 Section 3.2.3,
	// which allows RFC3986’s reserved character set, as well as percent-encoded triples, to pass through unchanged,
	// while still percent-encoding all other disallowed characters (including % outside of percent-encoded triples).
	// Applications are still responsible for percent-encoding reserved characters that are not allowed in the target media type;
	// see URL Percent-Encoding for details. The default value is false.
	// This field SHALL be ignored if the media type is not application/x-www-form-urlencoded or multipart/form-data.
	// If a value is explicitly defined, then the value of contentType (implicit or explicit) SHALL be ignored.
	AllowReserved bool `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyEncoding) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.Set("headers", &jsonschema.Schema{
		Type:        "object",
		Description: "A map allowing additional information to be provided as headers. Content-Type is described separately and SHALL be ignored in this section. This field SHALL be ignored if the media type is not a multipart.", //nolint:lll
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/RelixyHeader",
		},
	})
}

// RelixyHeader represents a high-level OpenAPI 3+ Header object that is backed by a low-level one.
//   - https://spec.openapis.org/oas/v3.2.0#header-object
type RelixyHeader struct {
	// A brief description of the header. This could contain examples of use.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Determines whether this header is mandatory. The default value is false.
	Required bool `json:"required,omitempty" yaml:"required,omitempty"`
	// Specifies that the header is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	// Describes how the header value will be serialized.
	// The default (and only legal value for headers) is simple.
	Style string `json:"style,omitempty" yaml:"style,omitempty"`
	// When this is true, header values of type array or object generate a single header whose value is a comma-separated list of the array items or key-value pairs of the map.
	// For other data types this field has no effect. The default value is false.
	Explode bool `json:"explode,omitempty" yaml:"explode,omitempty"`
	// The schema defining the type used for the header.
	Schema *RelixySchema `json:"schema,omitempty" yaml:"schema,omitempty"`
	// A map containing the representations for the header. The key is the media type and the value describes it. The map MUST only contain one entry.
	Content *orderedmap.OrderedMap[string, *RelixyMediaType] `json:"content,omitempty" yaml:"content,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyHeader) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("content", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelixyMediaType",
			},
		})
}

// Discriminator is only used by OpenAPI 3+ documents, it represents a polymorphic discriminator used for schemas
// When request bodies or response payloads may be one of a number of different schemas, a discriminator object can be used to aid in serialization, deserialization, and validation.
// The discriminator is a specific object in a schema which is used to inform the consumer of the document of an alternative schema based on the value associated with it.
// When using the discriminator, inline schemas will not be considered.
type Discriminator struct {
	// The name of the discriminating property in the payload that will hold the discriminating value.
	// The discriminating property MAY be defined as required or optional, but when defined as optional
	// the Discriminator Object MUST include a defaultMapping field that specifies which schema is expected to validate the structure of the model when the discriminating property is not present.
	PropertyName string `json:"propertyName" yaml:"propertyName"`
	// An object to hold mappings between payload values and schema names or URI references.
	Mapping *orderedmap.OrderedMap[string, string] `json:"mapping,omitempty" yaml:"mapping,omitempty"`
	// The schema name or URI reference to a schema that is expected to validate the structure of the model
	// when the discriminating property is not present in the payload or contains a value for which there is no explicit or implicit mapping.
	// OpenAPI 3.2+ defaultMapping for fallback schema
	DefaultMapping string `json:"defaultMapping,omitempty" yaml:"defaultMapping,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (Discriminator) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("mapping", &jsonschema.Schema{
			Type:        "object",
			Description: "An object to hold mappings between payload values and schema names or URI reference",
			AdditionalProperties: &jsonschema.Schema{
				Type: "string",
			},
		})
}

// RelixySchema represents a JSON Schema that support Swagger, OpenAPI 3 and OpenAPI 3.1
//
// Until 3.1 OpenAPI had a strange relationship with JSON Schema. It's been a super-set/sub-set
// mix, which has been confusing. So, instead of building a bunch of different models, we have compressed
// all variations into a single model that makes it easy to support multiple spec types.
//
//   - v2 schema: https://swagger.io/specification/v2/#schemaObject
//   - v3 schema: https://swagger.io/specification/#schema-object
//   - JSON Schema 2020-12: https://www.learnjsonschema.com/2020-12/
type RelixySchema struct {
	// 3.1 only, used to define a dialect for this schema, label is '$schema'.
	SchemaTypeRef string `json:"$schema,omitempty" yaml:"$schema,omitempty"`

	// Used to define a reference for this schema.
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`

	// The value of 'exclusiveMaximum' MUST be a number, representing an exclusive upper limit for a numeric instance.
	// If the instance is a number, then the instance is valid only if it has a value strictly less than (not equal to) 'exclusiveMaximum'.
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`

	// The value of 'exclusiveMinimum' MUST be a number, representing an exclusive lower limit for a numeric instance.
	// If the instance is a number, then the instance is valid only if it has a value strictly greater than (not equal to) 'exclusiveMinimum'.
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`

	// Types of the schema. Values MUST be one of the primitive types 'null', 'boolean', 'object', 'array', 'number', or 'string',
	// or 'integer' which matches any number with a zero fractional part.
	Type []PrimitiveType `json:"type,omitempty" yaml:"type,omitempty"`

	// The allOf keyword restricts instances to validate against every given subschema.
	// This keyword can be thought of as a logical conjunction (AND) operation,
	// as instances are valid if they satisfy every constraint of every subschema (the intersection of the constraints).
	AllOf []RelixySchema `json:"allOf,omitempty" yaml:"allOf,omitempty"`

	// The oneOf keyword restricts instances to validate against exactly one (and only one) of the given subschemas and fail on the rest.
	// This keyword represents a logical exclusive disjunction (XOR) operation.
	// In practice, the vast majority of schemas don’t require exclusive disjunction semantics but a simple disjunction.
	// If you are not sure, the anyOf keyword is probably a better fit.
	OneOf []RelixySchema `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	// The anyOf keyword restricts instances to validate against at least one (but potentially multiple) of the given subschemas.
	// This keyword represents a logical disjunction (OR) operation, as instances are valid if they satisfy the constraints of one or more subschemas (the union of the constraints).
	AnyOf []RelixySchema `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	// When request bodies or response payloads may be one of a number of different schemas,
	// these should use the JSON Schema anyOf or oneOf keywords to describe the possible schemas (see [Composition and Inheritance]).
	//
	// [Composition and Inheritance]: https://spec.openapis.org/oas/v3.2.0.html#composition-and-inheritance-polymorphism
	Discriminator *Discriminator `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`

	// The prefixItems keyword restricts a number of items from the start of an array instance to validate against the given sequence of subschemas,
	// where the item at a given index in the array instance is evaluated against the subschema at the given index in the prefixItems array, if any.
	// Information about the number of subschemas that were evaluated against the array instance is reported using annotations.
	// Array items outside the range described by the prefixItems keyword is evaluated against the items keyword, if present.
	PrefixItems []RelixySchema `json:"prefixItems,omitempty" yaml:"prefixItems,omitempty"`

	// 3.1 Specific properties
	Contains *RelixySchema `json:"contains,omitempty" yaml:"contains,omitempty"`

	// An instance array is valid against minContains in two ways, depending on the form of the annotation result of an adjacent 'contains' keyword.
	// The first way is if the annotation result is an array and the length of that array is greater than or equal to the 'minContains' value.
	// The second way is if the annotation result is a boolean "true" and the instance array length is greater than or equal to the 'minContains' value.
	//
	// If 'contains' is not present within the same schema object, then 'minContains' has no effect.
	// The value of 'minContains' MUST be a non-negative integer.
	// A value of 0 is allowed, but is only useful for setting a range of occurrences from 0 to the value of 'maxContains'.
	// A value of 0 with no 'maxContains' causes 'contains' to always pass validation.
	// Omitting this keyword has the same behavior as a value of 1.
	MinContains *int64 `json:"minContains,omitempty" yaml:"minContains,omitempty"`

	// An instance array is valid against 'maxContains' in two ways, depending on the form of the annotation result of an adjacent 'contains'.
	// The first way is if the annotation result is an array and the length of that array is less than or equal to the 'maxContains' value.
	// The second way is if the annotation result is a boolean 'true' and the instance array length is less than or equal to the 'maxContains' value.
	//
	// If 'contains' is not present within the same schema object, then 'maxContains' has no effect.
	// The value of 'maxContains' MUST be a non-negative integer.
	MaxContains *int64 `json:"maxContains,omitempty" yaml:"maxContains,omitempty"`
	// The if keyword introduces a subschema whose evaluation result restricts instances to validate against the then or else sibling subschemas (if present).
	// Note that the evaluation outcome of this subschema controls which other subschema to apply (if any) but has no direct effect on the overall validation result.
	If *RelixySchema `json:"if,omitempty" yaml:"if,omitempty"`
	// The else keyword restricts instances to validate against the given subschema if the if sibling keyword failed to validate against the instance.
	Else *RelixySchema `json:"else,omitempty" yaml:"else,omitempty"`
	// The then keyword restricts instances to validate against the given subschema if the if sibling keyword successfully validated against the instance.
	Then *RelixySchema `json:"then,omitempty" yaml:"then,omitempty"`
	// The dependentSchemas keyword restricts object instances to validate against one or more of the given subschemas if the corresponding properties are defined.
	// Note that the given subschemas are evaluated against the object that defines the property dependency.
	DependentSchemas *orderedmap.OrderedMap[string, *RelixySchema] `json:"dependentSchemas,omitempty" yaml:"dependentSchemas,omitempty"`
	// DependentRequired specifies properties that are required if a specific other property is present.
	// Their requirement is dependent on the presence of the other property.
	// The value of this field MUST be an object. Properties in this object, if any, MUST be arrays.
	// Elements in each array, if any, MUST be strings, and MUST be unique.
	// Validation succeeds if, for each name that appears in both the instance and as a name within this keyword's value,
	// every item in the corresponding array is also the name of a property in the instance.
	// Omitting this field has the same behavior as an empty object.
	DependentRequired *orderedmap.OrderedMap[string, []string] `json:"dependentRequired,omitempty" yaml:"dependentRequired,omitempty"`
	// The patternProperties keyword restricts properties of an object instance that match certain regular expressions to match their corresponding subschemas definitions.
	// Information about the properties that this keyword was evaluated for is reported using annotations.
	PatternProperties *orderedmap.OrderedMap[string, *RelixySchema] `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
	// The propertyNames keyword restricts object instances to only define properties whose names match the given schema.
	// This keyword is evaluated against every property of the object instance, independently of keywords that indirectly introduce property names such as properties and patternProperties.
	// Annotations coming from inside this keyword are dropped.
	PropertyNames *RelixySchema `json:"propertyNames,omitempty" yaml:"propertyNames,omitempty"`
	// The unevaluatedItems keyword is a generalisation of the items keyword that considers related keywords even when they are not direct siblings of this keyword.
	// More specifically, this keyword is affected by occurrences of prefixItems, items, contains, and unevaluatedItems itself,
	// as long as the evaluate path that led to unevaluatedItems is a prefix of the evaluate path of the others.
	UnevaluatedItems *RelixySchema `json:"unevaluatedItems,omitempty" yaml:"unevaluatedItems,omitempty"`

	// Schema of elements if the type is array.
	Items *RelixySchema `json:"items,omitempty" yaml:"items,omitempty"`

	// The $anchor keyword associates a subschema with the given URI fragment identifier, which can be referenced using the $ref keyword.
	// The fragment identifier is resolved against the URI of the schema resource.
	// Therefore, using this keyword to declare the same anchor more than once within the same schema resource results in an invalid schema.
	Anchor string `json:"$anchor,omitempty" yaml:"$anchor,omitempty"`
	// The not keyword restricts instances to fail validation against the given subschema.
	// This keyword represents a logical negation (NOT) operation.
	// In other words, the instance successfully validates against the schema only if it does not match the given subschema.
	Not *RelixySchema `json:"not,omitempty" yaml:"not,omitempty"`
	// Definition of object properties if the type is object.
	Properties *orderedmap.OrderedMap[string, *RelixySchema] `json:"properties,omitempty" yaml:"properties,omitempty"`
	// Title of the schema.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// A numeric instance is valid only if division by this field's value results in an integer.
	// The value of 'multipleOf' MUST be a number, strictly greater than 0.
	MultipleOf *float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty" jsonschema:"minimum=1"`
	// The value of 'maximum' MUST be a number, representing an inclusive upper limit for a numeric instance.
	// If the instance is a number, then this field validates only if the instance is less than or exactly equal to 'maximum'.
	Maximum *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	// The value of 'minimum' MUST be a number, representing an inclusive lower limit for a numeric instance.
	// If the instance is a number, then this field validates only if the instance is greater than or exactly equal to 'minimum'.
	Minimum *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	// A string instance is valid against this field if its length is less than, or equal to, the value of this field.
	// The length of a string instance is defined as the number of its characters as defined by RFC 8259.
	// The value of this field MUST be a non-negative integer.
	MaxLength *int64 `json:"maxLength,omitempty" yaml:"maxLength,omitempty" jsonschema:"minimum=0"`
	// A string instance is valid against this field if its length is greater than, or equal to, the value of this field.
	// The length of a string instance is defined as the number of its characters as defined by RFC 8259.
	// The value of this field MUST be a non-negative integer.
	// Omitting this field has the same behavior as a value of 0.
	MinLength *int64 `json:"minLength,omitempty" yaml:"minLength,omitempty" jsonschema:"minimum=0"`
	// A string instance is considered valid if the regular expression matches the instance successfully.
	// The value of this field MUST be a string. This string SHOULD be a valid regular expression, according to the ECMA-262 regular expression dialect.
	// Recall: regular expressions are not implicitly anchored.
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	// As defined by the JSON Schema Validation specification, data types can have an optional modifier keyword: format.
	// As described in that specification, format is treated as a non-validating annotation by default;
	// the ability to validate format varies across implementations.
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
	// An array instance is valid against 'maxItems' if its size is less than, or equal to, the value of this field.
	// The value of this keyword MUST be a non-negative integer.
	MaxItems *int64 `json:"maxItems,omitempty" yaml:"maxItems,omitempty" jsonschema:"minimum=0"`
	// An array instance is valid against 'minItems' if its size is greater than, or equal to, the value of this field.
	// The value of this field MUST be a non-negative integer.
	// Omitting this field has the same behavior as a value of 0.
	MinItems *int64 `json:"minItems,omitempty" yaml:"minItems,omitempty" jsonschema:"minimum=0"`
	// If this field has boolean value false, the instance validates successfully.
	// If it has boolean value true, the instance validates successfully if all of its elements are unique.
	// Omitting this field has the same behavior as a value of false.
	UniqueItems *bool `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	// An object instance is valid against 'maxProperties' if its number of properties is less than, or equal to, the value of this field.
	// The value of this keyword MUST be a non-negative integer.
	MaxProperties *int64 `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty" jsonschema:"minimum=0"`
	// An object instance is valid against 'minProperties' if its number of properties is greater than, or equal to, the value of this field.
	// The value of this field MUST be a non-negative integer.
	// Omitting this field has the same behavior as a value of 0.
	MinProperties *int64 `json:"minProperties,omitempty" yaml:"minProperties,omitempty" jsonschema:"minimum=0"`
	// The value of this field MUST be an array. Elements of this array, if any, MUST be strings, and MUST be unique.
	// An object instance is valid against this field if every item in the array is the name of a property in the instance.
	// Omitting this field has the same behavior as an empty array.
	Required []string `json:"required,omitempty" yaml:"required,omitempty"`
	// The value of this field MUST be an array. This array SHOULD have at least one element. Elements in the array SHOULD be unique.
	// An instance validates successfully against this keyword if its value is equal to one of the elements in this keyword's array value.
	// Elements in the array might be of any type, including null.
	Enum []any `json:"enum,omitempty" yaml:"enum,omitempty"`
	// Define a dictionary (also known as a map, hashmap or associative array) is a set of key/value pairs.
	AdditionalProperties *RelixySchema `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	// Description of the schema.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// The default keyword declares a default instance value for a schema or any of its subschemas,
	// typically to support specialised tooling like documentation and form generators.
	// This keyword does not affect validation, but the evaluator will collect its value as an annotation.
	Default any `json:"default,omitempty" yaml:"default,omitempty"`
	// Use of this field is functionally equivalent to an 'enum' with a single value.
	// An instance validates successfully against this keyword if its value is equal to the value of the field.
	// The value of this field MAY be of any type, including null.
	Const any `json:"const,omitempty" yaml:"const,omitempty"`
	// This keyword only takes effect if type is explicitly defined within the same Schema Object.
	// A true value indicates that both null values and values of the type specified by type are allowed.
	Nullable *bool `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	// Relevant only for Schema Object properties definitions.
	// Declares the property as "read only".
	// This means that it MAY be sent as part of a response but SHOULD NOT be sent as part of the request.
	// If the property is marked as readOnly being true and is in the required list, the required will take effect on the response only.
	// A property MUST NOT be marked as both readOnly and writeOnly being true. Default value is false.
	ReadOnly *bool `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	// Relevant only for Schema Object properties definitions. Declares the property as "write only".
	// Therefore, it MAY be sent as part of a request but SHOULD NOT be sent as part of the response.
	// If the property is marked as writeOnly being true and is in the required list, the required will take effect on the request only.
	// A property MUST NOT be marked as both readOnly and writeOnly being true. Default value is false.
	WriteOnly *bool `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	// Specifies that schema is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Deprecated *bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	// This MAY be used only on property schemas. It has no effect on root schemas.
	// Adds additional metadata to describe the XML representation of this property.
	// XML                  *XML                                `json:"xml,omitempty" yaml:"xml,omitempty"`
	// ExternalDocs         *ExternalDoc                        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// Example              *yaml.Node                          `json:"example,omitempty" yaml:"example,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixySchema) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("dependentSchemas", &jsonschema.Schema{
			Type:        "object",
			Description: "The dependentSchemas keyword restricts object instances to validate against one or more of the given subschemas if the corresponding properties are defined. Note that the given subschemas are evaluated against the object that defines the property dependency.", //nolint:lll
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelixySchema",
			},
		})
	schema.Properties.
		Set("patternProperties", &jsonschema.Schema{
			Type:        "object",
			Description: "The patternProperties keyword restricts properties of an object instance that match certain regular expressions to match their corresponding subschemas definitions. Information about the properties that this keyword was evaluated for is reported using annotations.", //nolint:lll
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelixySchema",
			},
		})
	schema.Properties.
		Set("dependentRequired", &jsonschema.Schema{
			Type:        "object",
			Description: "Specifies properties that are required if a specific other property is present. Their requirement is dependent on the presence of the other property. The value of this field MUST be an object. Properties in this object, if any, MUST be arrays. Elements in each array, if any, MUST be strings, and MUST be unique. Validation succeeds if, for each name that appears in both the instance and as a name within this keyword's value, every item in the corresponding array is also the name of a property in the instance. Omitting this field has the same behavior as an empty object.", //nolint:lll
			AdditionalProperties: &jsonschema.Schema{
				Type: "array",
				Items: &jsonschema.Schema{
					Type: "string",
				},
			},
		})
	schema.Properties.
		Set("properties", &jsonschema.Schema{
			Type:        "object",
			Description: "Definition of object properties if the type is object.",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelixySchema",
			},
		})
}

// RelixyAPIDocumentInfo represents a high-level Info object that provides metadata about the API.
// The metadata MAY be used by the clients if needed, and MAY be presented
// in editing or documentation generation tools for convenience.
type RelixyAPIDocumentInfo struct {
	// A short summary of the API.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// REQUIRED. The title of the API.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// A description of the API. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A URI for the Terms of Service for the API. This MUST be in the form of a URI.
	TermsOfService string `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	// The contact information for the exposed API.
	Contact *RelixyAPIDocumentInfoContact `json:"contact,omitempty" yaml:"contact,omitempty"`
	// The license information for the exposed API.
	License *RelixyAPIDocumentInfoLicense `json:"license,omitempty" yaml:"license,omitempty"`
	// REQUIRED. The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the version of the API being described or the version of the OpenAPI Description).
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

// RelixyAPIDocumentInfoContact represents a high-level representation of the Contact definitions.
type RelixyAPIDocumentInfoContact struct {
	// The identifying name of the contact person/organization.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// The URI for the contact information. This MUST be in the form of a URI.
	URL string `json:"url,omitempty" yaml:"url,omitempty"`
	// The email address of the contact person/organization. This MUST be in the form of an email address.
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// RelixyAPIDocumentInfoLicense represents a high-level representation of the License definitions.
type RelixyAPIDocumentInfoLicense struct {
	// REQUIRED. The license name used for the API.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// An SPDX-Licenses expression for the API. The identifier field is mutually exclusive of the url field.
	URL string `json:"url,omitempty" yaml:"url,omitempty"`
	// A URI for the license used for the API. This MUST be in the form of a URI.
	// The url field is mutually exclusive of the identifier field.
	Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
}
