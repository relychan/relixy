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

// Package schema defines the implementation for relixy resources.
package schema

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/relychan/goutils"
	"github.com/relychan/jsonschema"
	"github.com/relychan/openapitools/oaschema"
	"github.com/relychan/relixy/schema/baseschema"
	"github.com/relychan/rely-auth/auth"
	"go.yaml.in/yaml/v4"
)

var (
	ErrUnsupportedRelixyResourceKind    = errors.New("unsupported Relixy resource kind")
	ErrRelixyResourceDefinitionRequired = errors.New("require definition in Relixy resource")
)

// RelixyResource extends the Relixy resource interface to implement the JSON and YAML decoders.
type RelixyResource struct {
	baseschema.RelixyResource
}

type rawRelixyResourceJSON struct {
	baseschema.BaseResourceModel `yaml:",inline"`

	Definition json.RawMessage `json:"definition"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelixyResource) UnmarshalJSON(b []byte) error {
	var rawValue rawRelixyResourceJSON

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	if len(rawValue.Definition) == 0 {
		return ErrRelixyResourceDefinitionRequired
	}

	switch rawValue.Kind {
	case baseschema.RelyAuthKind:
		var authDef auth.RelyAuthDefinition

		err = json.Unmarshal(rawValue.Definition, &authDef)
		if err != nil {
			return err
		}

		j.RelixyResource = &baseschema.RelyAuthResource{
			BaseResourceModel: rawValue.BaseResourceModel,
			Definition:        authDef,
		}
	case baseschema.OpenAPIKind:
		var oasDef oaschema.OpenAPIResourceDefinition

		err = json.Unmarshal(rawValue.Definition, &oasDef)
		if err != nil {
			return err
		}

		j.RelixyResource = &RelixyOpenAPIResource{
			BaseResourceModel: rawValue.BaseResourceModel,
			Definition:        oasDef,
		}
	default:
		return fmt.Errorf("%w `%s`", ErrUnsupportedRelixyResourceKind, rawValue.Kind)
	}

	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (j *RelixyResource) UnmarshalYAML(value *yaml.Node) error {
	defNode, err := goutils.GetNodeValueFromYAMLMap(value, "definition")
	if err != nil {
		return err
	}

	if defNode == nil {
		return ErrRelixyResourceDefinitionRequired
	}

	var baseModel baseschema.BaseResourceModel

	err = value.Decode(&baseModel)
	if err != nil {
		return err
	}

	switch baseModel.Kind {
	case baseschema.RelyAuthKind:
		var authDef auth.RelyAuthDefinition

		err = defNode.Decode(&authDef)
		if err != nil {
			return err
		}

		j.RelixyResource = &baseschema.RelyAuthResource{
			BaseResourceModel: baseModel,
			Definition:        authDef,
		}
	case baseschema.OpenAPIKind:
		var oasDef oaschema.OpenAPIResourceDefinition

		err = defNode.Decode(&oasDef)
		if err != nil {
			return err
		}

		j.RelixyResource = &RelixyOpenAPIResource{
			BaseResourceModel: baseModel,
			Definition:        oasDef,
		}
	default:
		return fmt.Errorf("%w `%s`", ErrUnsupportedRelixyResourceKind, baseModel.Kind)
	}

	return nil
}

// JSONSchema defines a custom definition for JSON schema.
func (RelixyResource) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Description: "Definition of an OpenAPI resource",
				Ref:         "#/$defs/RelixyOpenAPIResource",
			},
			{
				Description: "Definition of a RelyAuth resource",
				Ref:         "#/$defs/RelyAuthConfig",
			},
		},
	}
}
