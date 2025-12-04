// Package schema defines config schemas for the relyx service.
package schema

import (
	"github.com/hasura/goenvconf"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	orderedmap "github.com/pb33f/ordered-map/v2"
	"github.com/relychan/gohttpc/httpconfig"
)

// RelyProxyAPIDocument represents the structure of an API document.
type RelyProxyAPIDocument struct {
	SchemaRef string `json:"$schema,omitempty" yaml:"$schema,omitempty"`
	// Global settings of the proxy.
	Settings RelyProxySettings `json:"settings" yaml:"settings"`
	// Info represents a specification Info definitions
	// Provides metadata about the API. The metadata MAY be used by tooling as required.
	// - https://spec.openapis.org/oas/v3.1.0#info-object
	Info *RelyProxyAPIDocumentInfo `json:"info,omitempty" yaml:"info,omitempty"`
	// Servers is a slice of Server instances which provide connectivity information to a target server.
	Servers []RelyProxyServer `json:"servers" yaml:"servers"`
	// Security contains global security requirements/roles for the specification
	// A declaration of which security mechanisms can be used across the API. The list of values includes alternative
	// security requirement objects that can be used. Only one of the security requirement objects need to be satisfied
	// to authorize a request. Individual operations can override this definition. To make security optional,
	// an empty security requirement ({}) can be included in the array.
	Security []SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`

	// Tags is a slice of Tag instances defined by the specification
	// A list of tags used by the document with additional metadata. The order of the tags can be used to reflect on
	// their order by the parsing tools. Not all tags that are used by the Operation Object must be declared.
	// The tags that are not declared MAY be organized randomly or based on the tools’ logic.
	// Each tag name in the list MUST be unique.
	Tags []RelyProxyTag `json:"tags,omitempty" yaml:"tags,omitempty"`

	// Components is an element to hold various schemas for the document.
	// Components *RelyProxyComponents `json:"components,omitempty" yaml:"components,omitempty"`

	// Paths contains all the PathItem definitions for the specification.
	// The available paths and operations for the API, The most important part of ths spec.
	Paths orderedmap.OrderedMap[string, *RelyProxyPathItem] `json:"paths" yaml:"paths"`
}

// RelyProxyTag represents a high-level Tag instance for grouping and filtering.
type RelyProxyTag struct {
	// Name of the tag.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Summary of the tag.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// Description of the tag.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Parent tag.
	Parent string `json:"parent,omitempty" yaml:"parent,omitempty"`
	// Kind of the tag.
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`
}

// OAS converts the tag to OpenAPI model.
func (prt RelyProxyTag) OAS() *base.Tag {
	return &base.Tag{
		Name:        prt.Name,
		Summary:     prt.Summary,
		Description: prt.Description,
		Parent:      prt.Parent,
		Kind:        prt.Kind,
	}
}

// RelyProxyServer contains server configurations.
type RelyProxyServer struct {
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
	SecuritySchemes []RelyProxySecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	// Security contains global security requirements/roles for the target.
	// Use global securities if empty.
	Security SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
	// TLS configuration of the server. Used for mTLS authentication.
	TLS *httpconfig.TLSConfig `json:"tls,omitempty" mapstructure:"tls" yaml:"tls,omitempty"`
}

// OAS3 converts the tag to OpenAPI model.
func (prs RelyProxyServer) OAS3() *highv3.Server {
	serverURL, _ := prs.URL.GetOrDefault("/")

	return &highv3.Server{
		URL:         serverURL,
		Name:        prs.Name,
		Description: prs.Description,
	}
}

// RelyProxyComponents is an element to hold various schemas for the document.
type RelyProxyComponents struct {
	// Schemas         *orderedmap.Map[string, *highbase.SchemaProxy] `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	// Responses       *orderedmap.Map[string, *Response]             `json:"responses,omitempty" yaml:"responses,omitempty"`
	// Parameters      *orderedmap.Map[string, *Parameter]            `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// RequestBodies   *orderedmap.Map[string, *RequestBody]          `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	// Headers         *orderedmap.Map[string, *Header]               `json:"headers,omitempty" yaml:"headers,omitempty"`
	// Callbacks       *orderedmap.Map[string, *Callback]             `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	// PathItems       *orderedmap.Map[string, *PathItem]             `json:"pathItems,omitempty" yaml:"pathItems,omitempty"`
	// MediaTypes      *orderedmap.Map[string, *MediaType]            `json:"mediaTypes,omitempty" yaml:"mediaTypes,omitempty"` // OpenAPI 3.2+ mediaTypes section
}

// RelyProxyPathItem represents a definition object of an API path item.
type RelyProxyPathItem struct {
	// Description of the path item.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Summary of the path item.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// Defines the operation information of the GET method.
	Get *RelyProxyOperation `json:"get,omitempty" yaml:"get,omitempty"`
	// Defines the operation information of the PUT method.
	Put *RelyProxyOperation `json:"put,omitempty" yaml:"put,omitempty"`
	// Defines the operation information of the POST method.
	Post *RelyProxyOperation `json:"post,omitempty" yaml:"post,omitempty"`
	// Defines the operation information of the DELETE method.
	Delete *RelyProxyOperation `json:"delete,omitempty" yaml:"delete,omitempty"`
	// Defines the operation information of the PATCH method.
	Patch *RelyProxyOperation `json:"patch,omitempty" yaml:"patch,omitempty"`
	// Common parameters of the path item.
	Parameters []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// RelyProxyOperation represents the definition of an API operation.
type RelyProxyOperation struct {
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
	RequestBody *RelyProxyRequestBody `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	// Responses   *highv3.Responses    `json:"responses,omitempty" yaml:"responses,omitempty"`
	// Defines whether this operation is deprecated.
	Deprecated *bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	// Requires the security for this operation.
	Security SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
	// Defines action information to proxy request to the remote server.
	Proxy RelyProxyAction `json:"proxy,omitempty" yaml:"proxy,omitempty"`
}
