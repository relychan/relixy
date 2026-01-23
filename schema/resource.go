package schema

import (
	"encoding/json"
	"errors"
	"fmt"

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

type rawRelixyResourceYAML struct {
	base_schema.BaseResourceModel `yaml:",inline"`

	// Definition of the OpenAPI documentation.
	Definition yaml.Node `json:"definition"`
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
	var rawValue rawRelixyResourceYAML

	err := value.Decode(&rawValue)
	if err != nil {
		return err
	}

	if rawValue.Definition.IsZero() {
		return ErrRelixyResourceDefinitionRequired
	}

	switch rawValue.Kind {
	case base_schema.RelyAuthKind:
		var authDef auth.RelyAuthDefinition

		err = rawValue.Definition.Decode(&authDef)
		if err != nil {
			return err
		}

		j.RelixyResource = &base_schema.RelyAuthResource{
			BaseResourceModel: rawValue.BaseResourceModel,
			Definition:        authDef,
		}
	case base_schema.OpenAPIKind:
		var oasDef openapi.RelixyOpenAPIResourceDefinition

		err = rawValue.Definition.Decode(&oasDef)
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
