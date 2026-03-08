package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/invopop/jsonschema"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/schema/baseschema"
	"go.yaml.in/yaml/v4"
)

const (
	nullTag = "!!null"
	strTag  = "!!str"
)

// RelixyOpenAPIResource represents an OpenAPI resource.
type RelixyOpenAPIResource struct {
	baseschema.BaseResourceModel `yaml:",inline"`

	// Definition of the OpenAPI documentation.
	Definition RelixyOpenAPIResourceDefinition `json:"definition" yaml:"definition"`
}

var _ baseschema.RelixyResource = (*RelixyOpenAPIResource)(nil)

// GetMetadata returns the metadata of the current resource.
func (ror RelixyOpenAPIResource) GetMetadata() baseschema.RelixyResourceMetadata {
	return ror.Metadata
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyOpenAPIResource) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("kind", &jsonschema.Schema{
			Description: "Kind of the resource which is always OpenAPI.",
			Type:        "string",
			Const:       baseschema.OpenAPIKind,
		})
}

// RelixyOpenAPIResourceDefinition defines fields of a relixy OpenAPI resource.
type RelixyOpenAPIResourceDefinition struct {
	// Settings of the OpenAPI v3 resource.
	Settings *RelixyOpenAPISettings `json:"settings,omitempty" yaml:"settings,omitempty"`
	// Path of URL of the referenced OpenAPI document.
	// Requires at least one of ref or spec.
	// If both fields are configured, the spec will be merged into the reference.
	Ref string `json:"ref,omitempty" yaml:"ref,omitempty"`
	// Specification of the OpenAPI v3 documentation.
	Spec *highv3.Document `json:"spec,omitempty" yaml:"spec,omitempty"`
}
type rawRelixyOpenAPIResourceDefinitionJSON struct {
	Settings *RelixyOpenAPISettings `json:"settings,omitempty"`
	Ref      string                 `json:"ref,omitempty"`
	Spec     json.RawMessage        `json:"spec"`
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyOpenAPIResourceDefinition) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("spec", &jsonschema.Schema{
			Title:       "OpenAPIv3Document",
			Description: "Specification of the OpenAPI v3 documentation.",
			Ref:         "https://raw.githubusercontent.com/relychan/relixy/refs/heads/main/jsonschema/openapi-3.json",
		})
}

// MarshalJSON implements json.Marshaler.
func (j RelixyOpenAPIResourceDefinition) MarshalJSON() ([]byte, error) {
	result := map[string]any{}

	if j.Ref != "" {
		result["ref"] = j.Ref
	}

	if j.Settings != nil {
		result["settings"] = j.Spec.Security
	}

	if j.Spec != nil {
		rawJSONBytes, err := j.Spec.RenderJSON("")
		if err != nil {
			return nil, err
		}

		result["spec"] = json.RawMessage(rawJSONBytes)
	}

	return json.Marshal(result)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelixyOpenAPIResourceDefinition) UnmarshalJSON(b []byte) error {
	var rawValue rawRelixyOpenAPIResourceDefinitionJSON

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	j.Ref = rawValue.Ref
	j.Settings = rawValue.Settings
	j.Spec = nil

	if len(rawValue.Spec) == 0 {
		return nil
	}

	doc, err := libopenapi.NewDocument(rawValue.Spec)
	if err != nil {
		return err
	}

	spec, err := doc.BuildV3Model()
	if err != nil {
		return err
	}

	j.Spec = &spec.Model

	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (j *RelixyOpenAPIResourceDefinition) UnmarshalYAML(value *yaml.Node) error {
	if value == nil || value.Kind != yaml.MappingNode {
		return fmt.Errorf(
			"%w. Expected an object, got %s",
			ErrInvalidOpenAPIResourceDefinitionYAML,
			value.Tag,
		)
	}

	contentLength := len(value.Content)

	for i := 0; i < contentLength; i++ {
		if i == contentLength-1 {
			break
		}

		keyNode := value.Content[i]
		if keyNode.Kind != yaml.ScalarNode && keyNode.Tag != "!!str" {
			return fmt.Errorf(
				"%w. Expected a key string, got %s",
				ErrInvalidOpenAPIResourceDefinitionYAML,
				keyNode.Tag,
			)
		}

		if keyNode.Value == "" {
			return fmt.Errorf("%w. Object key is empty", ErrInvalidOpenAPIResourceDefinitionYAML)
		}

		i++

		valueNode := value.Content[i]

		switch keyNode.Value {
		case "settings":
			err := valueNode.Decode(&j.Settings)
			if err != nil {
				return err
			}
		case "ref":
			switch valueNode.Tag {
			case strTag:
				j.Ref = valueNode.Value
			case nullTag:
			default:
				return fmt.Errorf(
					"%w. Expected ref is a string, got %s",
					ErrInvalidOpenAPIResourceDefinitionYAML,
					valueNode.Tag,
				)
			}
		case "spec":
			// Marshal the YAML node back to bytes for libopenapi
			specBytes, err := yaml.Dump(valueNode)
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
		default:
		}
	}

	return nil
}

// Build validates and merge the openapi specification with the reference if exist.
func (j *RelixyOpenAPIResourceDefinition) Build(ctx context.Context) (*highv3.Document, error) {
	if j.Ref == "" {
		if j.Spec == nil {
			return nil, ErrResourceSpecRequired
		}

		return j.Spec, nil
	}

	rawSourceReader, _, err := goutils.FileReaderFromPath(ctx, j.Ref)
	if err != nil {
		return nil, err
	}

	defer goutils.CatchWarnErrorFunc(rawSourceReader.Close)

	sourceBytes, err := io.ReadAll(rawSourceReader)
	if err != nil {
		return nil, err
	}

	sourceDoc, err := libopenapi.NewDocument(sourceBytes)
	if err != nil {
		return nil, err
	}

	var doc *highv3.Document

	if sourceDoc.GetSpecInfo().SpecFormat == datamodel.OAS2 {
		spec, err := sourceDoc.BuildV2Model()
		if err != nil {
			return nil, err
		}

		doc, err = convertSwaggerToOpenAPIv3Document(&spec.Model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert openapi spec v2 to v3: %w", err)
		}
	} else {
		sourceSpec, err := sourceDoc.BuildV3Model()
		if err != nil {
			return nil, err
		}

		doc = &sourceSpec.Model
	}

	if j.Spec != nil {
		mergeOpenAPIv3Documents(doc, j.Spec)
	}

	return doc, nil
}
