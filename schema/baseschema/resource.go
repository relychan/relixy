// Package baseschema define base schema resources.
package baseschema

import (
	"github.com/invopop/jsonschema"
	"github.com/relychan/rely-auth/auth"
)

// RelixyResourceKind represents a kind of the Relixy resource.
type RelixyResourceKind string

const (
	// RelyAuthKind represents a kind enum for the RelyAuth specification.
	RelyAuthKind RelixyResourceKind = "RelyAuth"
	// OpenAPIKind represents a kind enum for the OpenAPI specification.
	OpenAPIKind RelixyResourceKind = "OpenAPI"
)

// RelixyResourceMetadata represents common metadata of the resource.
type RelixyResourceMetadata struct {
	// Name of the resource.
	Name string `json:"name" yaml:"name" jsonschema:"default=default"`
	// Description of the resource.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// System instruction for the current resource. It's important to the LLM.
	Instruction string `json:"instruction,omitempty" yaml:"instruction,omitempty"`
}

// RelixyResource abstracts an interface for Relixy resource.
type RelixyResource interface {
	// GetBaseResource returns the base resource information of the current resource.
	GetBaseResource() BaseResourceModel
}

// BaseResourceModel defines the base structure of a resource.
type BaseResourceModel struct {
	// Version of the authentication config.
	Version string `json:"version" yaml:"version" jsonschema:"enum=v1"`
	// Kind of the resource.
	Kind RelixyResourceKind `json:"kind" yaml:"kind"`
	// Metadata of the resource.
	Metadata RelixyResourceMetadata `json:"metadata" yaml:"metadata"`
}

var _ RelixyResource = (*BaseResourceModel)(nil)

// GetBaseResource returns the base resource information of the current resource.
func (ror BaseResourceModel) GetBaseResource() BaseResourceModel {
	return ror
}

// RelyAuthResource represents a RelyAuth resource.
type RelyAuthResource struct {
	BaseResourceModel `yaml:",inline"`

	// Definition of the OpenAPI documentation.
	Definition auth.RelyAuthDefinition `json:"definition" yaml:"definition"`
}

var _ RelixyResource = (*RelyAuthResource)(nil)

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelyAuthResource) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("kind", &jsonschema.Schema{
			Description: "Kind of the resource which is always RelyAuth.",
			Type:        "string",
			Const:       RelyAuthKind,
		})

	schema.Properties.Delete("metadata")
}
