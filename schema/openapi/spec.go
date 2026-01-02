package openapi

import (
	"github.com/hasura/goenvconf"
	"github.com/invopop/jsonschema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	orderedmap "github.com/pb33f/ordered-map/v2"
	"github.com/relychan/gohttpc/httpconfig"
	"github.com/relychan/relixy/schema/base_schema"
)

// RelixyOpenAPI3ResourceSpecification represents the structure of an API document.
type RelixyOpenAPI3ResourceSpecification struct {
	// Version is the version of OpenAPI being used, extracted from the 'openapi: x.x.x' definition.
	// This is not a standard property of the OpenAPI model, it's a convenience mechanism only.
	Version string `json:"openapi,omitempty" yaml:"openapi,omitempty"`
	// OpenAPI document reference that this document depends on.
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	// Info represents a specification Info definitions
	// Provides metadata about the API. The metadata MAY be used by tooling as required.
	// - https://spec.openapis.org/oas/v3.1.0#info-object
	Info *base.Info `json:"info,omitempty" yaml:"info,omitempty"`
	// Servers is a slice of Server instances which provide connectivity information to a target server. If the servers
	// property is not provided, or is an empty array, the default value would be a Server Object with an url value of /.
	// - https://spec.openapis.org/oas/v3.1.0#server-object
	Servers []*RelixyServer `json:"servers" yaml:"servers"`
	// Security contains global security requirements/roles for the specification
	// A declaration of which security mechanisms can be used across the API. The list of values includes alternative
	// security requirement objects that can be used. Only one of the security requirement objects need to be satisfied
	// to authorize a request. Individual operations can override this definition. To make security optional,
	// an empty security requirement ({}) can be included in the array.
	Security []SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`

	// Components is an element to hold various schemas for the document.
	// Components *RelixyComponents `json:"components,omitempty" yaml:"components,omitempty"`

	// Paths contains all the PathItem definitions for the specification.
	// The available paths and operations for the API, The most important part of ths spec.
	// - https://spec.openapis.org/oas/v3.1.0#paths-object
	Paths orderedmap.OrderedMap[string, *RelixyOpenAPIv3PathItem] `json:"paths" yaml:"paths"`

	// Components is an element to hold various schemas for the document.
	// - https://spec.openapis.org/oas/v3.1.0#components-object
	Components *highv3.Components `json:"components,omitempty" yaml:"components,omitempty"`

	// Tags is a slice of base.Tag instances defined by the specification
	// A list of tags used by the document with additional metadata. The order of the tags can be used to reflect on
	// their order by the parsing tools. Not all tags that are used by the Operation Object must be declared.
	// The tags that are not declared MAY be organized randomly or based on the tools’ logic.
	// Each tag name in the list MUST be unique.
	// - https://spec.openapis.org/oas/v3.1.0#tag-object
	Tags []*base.Tag `json:"tags,omitempty" yaml:"tags,omitempty"`

	// ExternalDocs is an instance of base.ExternalDoc for.. well, obvious really, innit.
	// - https://spec.openapis.org/oas/v3.1.0#external-documentation-object
	ExternalDocs *base.ExternalDoc `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	// JsonSchemaDialect is a 3.1+ property that sets the dialect to use for validating *base.Schema definitions
	// The default value for the $schema keyword within Schema Objects contained within this OAS document.
	// This MUST be in the form of a URI.
	// - https://spec.openapis.org/oas/v3.1.0#schema-object
	JsonSchemaDialect string `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`

	// Self is a 3.2+ property that sets the base URI for the document for resolving relative references
	// - https://spec.openapis.org/oas/v3.2.0#openapi-object
	Self string `json:"$self,omitempty" yaml:"$self,omitempty"`

	// Webhooks is a 3.1+ property that is similar to callbacks, except, this defines incoming webhooks.
	// The incoming webhooks that MAY be received as part of this API and that the API consumer MAY choose to implement.
	// Closely related to the callbacks feature, this section describes requests initiated other than by an API call,
	// for example by an out-of-band registration. The key name is a unique string to refer to each webhook,
	// while the (optionally referenced) Path Item Object describes a request that may be initiated by the API provider
	// and the expected responses. An example is available.
	Webhooks *orderedmap.OrderedMap[string, *RelixyOpenAPIv3PathItem] `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyOpenAPI3ResourceSpecification) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("paths", &jsonschema.Schema{
			Type:        "object",
			Description: "Holds the relative paths to the individual endpoints and their operations. The path is appended to the URL from the Server Object in order to construct the full URL. The Paths Object MAY be empty, due to Access Control List (ACL) constraints.", //nolint:lll
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelixyOpenAPIv3PathItem",
			},
		})
}

// RelixyServer contains server configurations.
type RelixyServer struct {
	// Defines the server URL
	URL goenvconf.EnvString `json:"url" yaml:"url"`
	// An optional unique string to refer to the host designated by the URL.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Description of the server.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Defines custom headers to be injected to the remote server.
	// Merged with the global headers.
	Headers map[string]goenvconf.EnvString `json:"headers,omitempty" yaml:"headers,omitempty"`
	// Defines the weight of the server endpoint for load balancing.
	// Only take effect if there are many servers.
	Weight *int `json:"weight,omitempty" yaml:"weight,omitempty" jsonschema:"min=0,max=100"`
	// SecuritySchemes define security schemes that can be used by the operations.
	// Use global security schemes if empty.
	SecuritySchemes []RelixySecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	// Security contains global security requirements/roles for the target.
	// Use global securities if empty.
	Security SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
	// TLS configuration of the server. Used for mTLS authentication.
	TLS *httpconfig.TLSConfig `json:"tls,omitempty" mapstructure:"tls" yaml:"tls,omitempty"`
}

// OAS3 converts the tag to OpenAPI model.
func (prs RelixyServer) OAS3() *highv3.Server {
	serverURL, _ := prs.URL.GetOrDefault("/")

	return &highv3.Server{
		URL:         serverURL,
		Name:        prs.Name,
		Description: prs.Description,
	}
}

// RelixyOpenAPIv3PathItem represents a definition object of an API path item.
type RelixyOpenAPIv3PathItem struct {
	// Reference URI of the path item.
	Reference string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	// Description of the path item.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Summary of the path item.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// Defines the operation information of the GET method.
	Get *RelixyOpenAPIv3Operation `json:"get,omitempty" yaml:"get,omitempty"`
	// Defines the operation information of the PUT method.
	Put *RelixyOpenAPIv3Operation `json:"put,omitempty" yaml:"put,omitempty"`
	// Defines the operation information of the POST method.
	Post *RelixyOpenAPIv3Operation `json:"post,omitempty" yaml:"post,omitempty"`
	// Defines the operation information of the DELETE method.
	Delete *RelixyOpenAPIv3Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
	// Defines the operation information of the OPTIONS method.
	Options *RelixyOpenAPIv3Operation `json:"options,omitempty" yaml:"options,omitempty"`
	// Defines the operation information of the HEAD method.
	Head *RelixyOpenAPIv3Operation `json:"head,omitempty" yaml:"head,omitempty"`
	// Defines the operation information of the PATCH method.
	Patch *RelixyOpenAPIv3Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
	// Defines the operation information of the TRACE method.
	Trace *RelixyOpenAPIv3Operation `json:"trace,omitempty" yaml:"trace,omitempty"`
	// Defines the operation information of the QUERY method.
	Query *RelixyOpenAPIv3Operation `json:"query,omitempty" yaml:"query,omitempty"`
	// Common parameters of the path item.
	Parameters []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// Additional operations beside common operations.
	AdditionalOperations *orderedmap.OrderedMap[string, *RelixyOpenAPIv3Operation] `json:"additionalOperations,omitempty" yaml:"additionalOperations,omitempty"` // OpenAPI 3.2+ additional operations
	// Defines exclusive headers for this path item.
	Servers []*RelixyServer `json:"servers,omitempty" yaml:"servers,omitempty"`
}

// RelixyOpenAPIv3Operation represents the definition of an API operation.
type RelixyOpenAPIv3Operation struct {
	// List of tags for this operation.
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Summary information of the operation.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// Description of the operation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// ID of the operation. It should be unique.
	OperationID string `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	// Request parameters of the operation.
	Parameters []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// Request body information.
	RequestBody *RelixyRequestBody `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	// Responses   *highv3.Responses    `json:"responses,omitempty" yaml:"responses,omitempty"`
	// Defines whether this operation is deprecated.
	Deprecated *bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	// Requires the security for this operation.
	Security SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
	// Defines action information to proxy request to the remote server.
	Proxy        base_schema.RelixyAction `json:"proxy" yaml:"proxy"`
	ExternalDocs *base.ExternalDoc        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Responses    *highv3.Responses        `json:"responses,omitempty" yaml:"responses,omitempty"`
	Servers      []*RelixyServer          `json:"servers,omitempty" yaml:"servers,omitempty"`
	// Callbacks    *orderedmap.OrderedMap[string, *highv3.Callback]  `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

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
	Schema *RelixyOpenAPIv3Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
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
	Schema *RelixyOpenAPIv3Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	// A schema describing each item within a sequential media type.
	ItemSchema *RelixyOpenAPIv3Schema `json:"itemSchema,omitempty" yaml:"itemSchema,omitempty"`
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
	Schema *RelixyOpenAPIv3Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
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
