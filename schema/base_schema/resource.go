// Package base_schema define base schema resources.
package base_schema

// RelixyResourceKind represents a kind of the Relixy resource.
type RelixyResourceKind string

const (
	// OpenAPI3Kind represents a kind enum for the OpenAPI 3 specification.
	OpenAPI3Kind RelixyResourceKind = "OpenAPI3"
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
	Kind string `json:"kind" yaml:"kind"`
	// Metadata of the resource.
	Metadata RelixyResourceMetadata `json:"metadata" yaml:"metadata"`
}

var _ RelixyResource = (*BaseResourceModel)(nil)

// GetBaseResource returns the base resource information of the current resource.
func (ror BaseResourceModel) GetBaseResource() BaseResourceModel {
	return ror
}
