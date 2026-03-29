// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package baseschema define base schema resources.
package baseschema

import (
	"github.com/relychan/jsonschema"
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
	Name string `json:"name" yaml:"name" jsonschema:"default=default,description=Name of the resource"`
	// Description of the resource.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Description of the resource"`
	// System instruction for the current resource. It's important to the LLM.
	Instruction string `json:"instruction,omitempty" yaml:"instruction,omitempty" jsonschema:"description=System instruction for the current resource. It's important to the LLM"`
}

// RelixyResource abstracts an interface for Relixy resource.
type RelixyResource interface {
	// GetBaseResource returns the base resource information of the current resource.
	GetBaseResource() BaseResourceModel
}

// BaseResourceModel defines the base structure of a resource.
type BaseResourceModel struct {
	// Version of the resource.
	Version string `json:"version" yaml:"version" jsonschema:"enum=v1,description=Version of the resource"`
	// Kind of the resource.
	Kind RelixyResourceKind `json:"kind" yaml:"kind"`
	// Metadata of the resource.
	Metadata RelixyResourceMetadata `json:"metadata" yaml:"metadata" jsonschema:"description=Metadata of the resource"`
}

var _ RelixyResource = (*BaseResourceModel)(nil)

// GetBaseResource returns the base resource information of the current resource.
func (ror BaseResourceModel) GetBaseResource() BaseResourceModel {
	return ror
}

// RelyAuthResource represents a RelyAuth resource.
type RelyAuthResource struct {
	BaseResourceModel `yaml:",inline"`

	// Raw definition of a resource.
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
