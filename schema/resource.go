// Package schema defines the implementation for relixy resources.
package schema

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/schema/base_schema"
	"github.com/relychan/relixy/schema/openapi"
	"github.com/relychan/rely-auth/auth"
	"go.yaml.in/yaml/v4"
)

var (
	ErrUnsupportedRelixyResourceKind    = errors.New("unsupported Relixy resource kind")
	ErrRelixyResourceDefinitionRequired = errors.New("require definition in Relixy resource")
)

// RelixyResource extends the Relixy resource interface to implement the JSON and YAML decoders.
type RelixyResource struct {
	base_schema.RelixyResource
}

type rawRelixyResourceJSON struct {
	base_schema.BaseResourceModel `yaml:",inline"`

	// Definition of the OpenAPI documentation.
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
	case base_schema.RelyAuthKind:
		var authDef auth.RelyAuthDefinition

		err = json.Unmarshal(rawValue.Definition, &authDef)
		if err != nil {
			return err
		}

		j.RelixyResource = &base_schema.RelyAuthResource{
			BaseResourceModel: rawValue.BaseResourceModel,
			Definition:        authDef,
		}
	case base_schema.OpenAPIKind:
		var oasDef openapi.RelixyOpenAPIResourceDefinition

		err = json.Unmarshal(rawValue.Definition, &oasDef)
		if err != nil {
			return err
		}

		j.RelixyResource = &openapi.RelixyOpenAPIResource{
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

	var baseModel base_schema.BaseResourceModel

	err = value.Decode(&baseModel)
	if err != nil {
		return err
	}

	switch baseModel.Kind {
	case base_schema.RelyAuthKind:
		var authDef auth.RelyAuthDefinition

		err = defNode.Decode(&authDef)
		if err != nil {
			return err
		}

		j.RelixyResource = &base_schema.RelyAuthResource{
			BaseResourceModel: baseModel,
			Definition:        authDef,
		}
	case base_schema.OpenAPIKind:
		var oasDef openapi.RelixyOpenAPIResourceDefinition

		err = defNode.Decode(&oasDef)
		if err != nil {
			return err
		}

		j.RelixyResource = &openapi.RelixyOpenAPIResource{
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
				Ref:         "https://raw.githubusercontent.com/relychan/rely-auth/refs/heads/main/jsonschema/auth.schema.json",
			},
		},
	}
}
