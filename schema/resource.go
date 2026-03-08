// Package schema defines the implementation for relixy resources.
package schema

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/schema/baseschema"
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
