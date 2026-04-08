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

package schema

import (
	"github.com/relychan/jsonschema"
	"github.com/relychan/openapitools/oaschema"
	"github.com/relychan/openapitools/openapiclient"
	"github.com/relychan/relixy/schema/baseschema"
)

// OpenAPIClient represents the wrapper of an OpenAPI client.
type OpenAPIClient struct {
	*openapiclient.ProxyClient

	ResourceMetadata baseschema.RelixyResourceMetadata
}

// RelixyOpenAPIResource represents an OpenAPI resource.
type RelixyOpenAPIResource struct {
	baseschema.BaseResourceModel `yaml:",inline"`

	// Definition of the OpenAPI documentation.
	Definition oaschema.OpenAPIResourceDefinition `json:"definition" yaml:"definition"`
}

var _ baseschema.RelixyResource = (*RelixyOpenAPIResource)(nil)

// GetMetadata returns the metadata of the current resource.
func (ror RelixyOpenAPIResource) GetMetadata() baseschema.RelixyResourceMetadata {
	return ror.Metadata
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyOpenAPIResource) JSONSchemaExtend(schema *jsonschema.Schema) {
	defSchema, _ := schema.Properties.Get("definition")
	defSchema.Description = "Definition of the OpenAPI documentation"
	schema.Properties.Set("definition", defSchema)

	schema.Properties.
		Set("kind", &jsonschema.Schema{
			Description: "Kind of the resource which is always OpenAPI.",
			Type:        "string",
			Const:       baseschema.OpenAPIKind,
		})
}
