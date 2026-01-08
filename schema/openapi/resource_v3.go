package openapi

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/pb33f/libopenapi"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/relixy/schema/base_schema"
	"go.yaml.in/yaml/v4"
)

// ErrResourceSpecRequired occurs when the spec field of resource is empty.
var ErrResourceSpecRequired = errors.New("spec is required in resource")

// RelixyOpenAPIv3Resource represents an OpenAPI v3 resource.
type RelixyOpenAPIv3Resource struct {
	base_schema.BaseResourceModel `yaml:",inline"`

	// Definition of the OpenAPI v3 documentation.
	Definition RelixyOpenAPIv3ResourceDefinition `json:"definition" yaml:"definition"`
}

var _ base_schema.RelixyResource = (*RelixyOpenAPIv3Resource)(nil)

// GetMetadata returns the metadata of the current resource.
func (ror RelixyOpenAPIv3Resource) GetMetadata() base_schema.RelixyResourceMetadata {
	return ror.Metadata
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyOpenAPIv3Resource) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("kind", &jsonschema.Schema{
			Description: "Kind of the resource which is always OpenAPI3.",
			Type:        "string",
			Const:       "OpenAPI3",
		})
}

type RelixyOpenAPIv3ResourceDefinition struct {
	// Settings of the OpenAPI v3 resource.
	Settings RelixyOpenAPISettings `json:"settings" yaml:"settings"`
	// Specification of the OpenAPI v3 documentation.
	Spec *highv3.Document `json:"spec" yaml:"spec"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyOpenAPIv3ResourceDefinition) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("spec", &jsonschema.Schema{
			Description: "Specification of the OpenAPI v3 documentation.",
			Ref:         "openapi-3.json",
		})
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelixyOpenAPIv3ResourceDefinition) UnmarshalJSON(b []byte) error {
	rawValue := map[string]json.RawMessage{}

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	rawSpec, ok := rawValue["spec"]
	if !ok || rawSpec == nil {
		return fmt.Errorf("%w OpenAPI v3", ErrResourceSpecRequired)
	}

	doc, err := libopenapi.NewDocument(rawSpec)
	if err != nil {
		return err
	}

	spec, err := doc.BuildV3Model()
	if err != nil {
		return err
	}

	j.Spec = &spec.Model

	rawSettings, ok := rawValue["settings"]
	if ok && rawSettings != nil {
		err = json.Unmarshal(rawSettings, &j.Settings)
		if err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (j *RelixyOpenAPIv3ResourceDefinition) UnmarshalYAML(value *yaml.Node) error {
	rawValue := map[string]yaml.Node{}

	err := value.Decode(&rawValue)
	if err != nil {
		return err
	}

	rawSpec, ok := rawValue["spec"]
	if !ok {
		return fmt.Errorf("%w OpenAPI v3", ErrResourceSpecRequired)
	}

	// Marshal the YAML node back to bytes for libopenapi
	specBytes, err := yaml.Marshal(rawSpec)
	if err != nil {
		return err
	}

	doc, err := libopenapi.NewDocument(specBytes)
	if err != nil {
		return err
	}

	spec, err := doc.BuildV3Model()
	if err != nil {
		return err
	}

	j.Spec = &spec.Model

	rawSettings, ok := rawValue["settings"]
	if ok {
		err = rawSettings.Decode(&j.Settings)
		if err != nil {
			return err
		}
	}

	return nil
}
