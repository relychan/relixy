package openapi

import (
	"github.com/invopop/jsonschema"
	orderedmap "github.com/pb33f/ordered-map/v2"
	"github.com/relychan/relixy/schema/base_schema"
)

// RelixySchema represents a JSON Schema that support Swagger, OpenAPI 3 and OpenAPI 3.1
//
// Until 3.1 OpenAPI had a strange relationship with JSON Schema. It's been a super-set/sub-set
// mix, which has been confusing. So, instead of building a bunch of different models, we have compressed
// all variations into a single model that makes it easy to support multiple spec types.
//
//   - v2 schema: https://swagger.io/specification/v2/#schemaObject
//   - v3 schema: https://swagger.io/specification/#schema-object
//   - JSON Schema 2020-12: https://www.learnjsonschema.com/2020-12/
type RelixyOpenAPIv3Schema struct {
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
	Type []base_schema.PrimitiveType `json:"type,omitempty" yaml:"type,omitempty"`

	// The allOf keyword restricts instances to validate against every given subschema.
	// This keyword can be thought of as a logical conjunction (AND) operation,
	// as instances are valid if they satisfy every constraint of every subschema (the intersection of the constraints).
	AllOf []RelixyOpenAPIv3Schema `json:"allOf,omitempty" yaml:"allOf,omitempty"`

	// The oneOf keyword restricts instances to validate against exactly one (and only one) of the given subschemas and fail on the rest.
	// This keyword represents a logical exclusive disjunction (XOR) operation.
	// In practice, the vast majority of schemas don’t require exclusive disjunction semantics but a simple disjunction.
	// If you are not sure, the anyOf keyword is probably a better fit.
	OneOf []RelixyOpenAPIv3Schema `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	// The anyOf keyword restricts instances to validate against at least one (but potentially multiple) of the given subschemas.
	// This keyword represents a logical disjunction (OR) operation, as instances are valid if they satisfy the constraints of one or more subschemas (the union of the constraints).
	AnyOf []RelixyOpenAPIv3Schema `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	// When request bodies or response payloads may be one of a number of different schemas,
	// these should use the JSON Schema anyOf or oneOf keywords to describe the possible schemas (see [Composition and Inheritance]).
	//
	// [Composition and Inheritance]: https://spec.openapis.org/oas/v3.2.0.html#composition-and-inheritance-polymorphism
	Discriminator *Discriminator `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`

	// The prefixItems keyword restricts a number of items from the start of an array instance to validate against the given sequence of subschemas,
	// where the item at a given index in the array instance is evaluated against the subschema at the given index in the prefixItems array, if any.
	// Information about the number of subschemas that were evaluated against the array instance is reported using annotations.
	// Array items outside the range described by the prefixItems keyword is evaluated against the items keyword, if present.
	PrefixItems []RelixyOpenAPIv3Schema `json:"prefixItems,omitempty" yaml:"prefixItems,omitempty"`

	// 3.1 Specific properties
	Contains *RelixyOpenAPIv3Schema `json:"contains,omitempty" yaml:"contains,omitempty"`

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
	If *RelixyOpenAPIv3Schema `json:"if,omitempty" yaml:"if,omitempty"`
	// The else keyword restricts instances to validate against the given subschema if the if sibling keyword failed to validate against the instance.
	Else *RelixyOpenAPIv3Schema `json:"else,omitempty" yaml:"else,omitempty"`
	// The then keyword restricts instances to validate against the given subschema if the if sibling keyword successfully validated against the instance.
	Then *RelixyOpenAPIv3Schema `json:"then,omitempty" yaml:"then,omitempty"`
	// The dependentSchemas keyword restricts object instances to validate against one or more of the given subschemas if the corresponding properties are defined.
	// Note that the given subschemas are evaluated against the object that defines the property dependency.
	DependentSchemas *orderedmap.OrderedMap[string, *RelixyOpenAPIv3Schema] `json:"dependentSchemas,omitempty" yaml:"dependentSchemas,omitempty"`
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
	PatternProperties *orderedmap.OrderedMap[string, *RelixyOpenAPIv3Schema] `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
	// The propertyNames keyword restricts object instances to only define properties whose names match the given schema.
	// This keyword is evaluated against every property of the object instance, independently of keywords that indirectly introduce property names such as properties and patternProperties.
	// Annotations coming from inside this keyword are dropped.
	PropertyNames *RelixyOpenAPIv3Schema `json:"propertyNames,omitempty" yaml:"propertyNames,omitempty"`
	// The unevaluatedItems keyword is a generalisation of the items keyword that considers related keywords even when they are not direct siblings of this keyword.
	// More specifically, this keyword is affected by occurrences of prefixItems, items, contains, and unevaluatedItems itself,
	// as long as the evaluate path that led to unevaluatedItems is a prefix of the evaluate path of the others.
	UnevaluatedItems *RelixyOpenAPIv3Schema `json:"unevaluatedItems,omitempty" yaml:"unevaluatedItems,omitempty"`

	// Schema of elements if the type is array.
	Items *RelixyOpenAPIv3Schema `json:"items,omitempty" yaml:"items,omitempty"`

	// The $anchor keyword associates a subschema with the given URI fragment identifier, which can be referenced using the $ref keyword.
	// The fragment identifier is resolved against the URI of the schema resource.
	// Therefore, using this keyword to declare the same anchor more than once within the same schema resource results in an invalid schema.
	Anchor string `json:"$anchor,omitempty" yaml:"$anchor,omitempty"`
	// The not keyword restricts instances to fail validation against the given subschema.
	// This keyword represents a logical negation (NOT) operation.
	// In other words, the instance successfully validates against the schema only if it does not match the given subschema.
	Not *RelixyOpenAPIv3Schema `json:"not,omitempty" yaml:"not,omitempty"`
	// Definition of object properties if the type is object.
	Properties *orderedmap.OrderedMap[string, *RelixyOpenAPIv3Schema] `json:"properties,omitempty" yaml:"properties,omitempty"`
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
	AdditionalProperties *RelixyOpenAPIv3Schema `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
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
func (RelixyOpenAPIv3Schema) JSONSchemaExtend(schema *jsonschema.Schema) {
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
