package openapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/invopop/jsonschema"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	orderedmap "github.com/pb33f/ordered-map/v2"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/schema/base_schema"
	"go.yaml.in/yaml/v4"
)

// ErrResourceSpecRequired occurs when the spec field of resource is empty.
var ErrResourceSpecRequired = errors.New("spec is required in resource")

// RelixyOpenAPIResource represents an OpenAPI resource.
type RelixyOpenAPIResource struct {
	base_schema.BaseResourceModel `yaml:",inline"`

	// Definition of the OpenAPI documentation.
	Definition RelixyOpenAPIResourceDefinition `json:"definition" yaml:"definition"`
}

var _ base_schema.RelixyResource = (*RelixyOpenAPIResource)(nil)

// GetMetadata returns the metadata of the current resource.
func (ror RelixyOpenAPIResource) GetMetadata() base_schema.RelixyResourceMetadata {
	return ror.Metadata
}

// JSONSchemaExtend modifies the JSON schema afterwards.
func (RelixyOpenAPIResource) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Properties.
		Set("kind", &jsonschema.Schema{
			Description: "Kind of the resource which is always OpenAPI.",
			Type:        "string",
			Const:       "OpenAPI",
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
	Spec *highv3.Document `json:"spec" yaml:"spec"`
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

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelixyOpenAPIResourceDefinition) UnmarshalJSON(b []byte) error {
	rawValue := map[string]json.RawMessage{}

	err := json.Unmarshal(b, &rawValue)
	if err != nil {
		return err
	}

	rawSettings, ok := rawValue["settings"]
	if ok && rawSettings != nil {
		err = json.Unmarshal(rawSettings, &j.Settings)
		if err != nil {
			return err
		}
	}

	rawSpec, ok := rawValue["spec"]
	if !ok || rawSpec == nil {
		return nil
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

	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (j *RelixyOpenAPIResourceDefinition) UnmarshalYAML(value *yaml.Node) error {
	rawValue := map[string]yaml.Node{}

	err := value.Decode(&rawValue)
	if err != nil {
		return err
	}

	rawSettings, ok := rawValue["settings"]
	if ok {
		err = rawSettings.Decode(&j.Settings)
		if err != nil {
			return err
		}
	}

	rawSpec, ok := rawValue["spec"]
	if !ok {
		return nil
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

	sourceSpec, err := sourceDoc.BuildV3Model()
	if err != nil {
		return nil, err
	}

	if j.Spec != nil {
		mergeOpenAPIv3Documents(&sourceSpec.Model, j.Spec)
	}

	return &sourceSpec.Model, nil
}

func mergeOpenAPIv3Documents(dest, src *highv3.Document) {
	dest.Extensions = mergeOrderedMaps(dest.Extensions, src.Extensions)
	dest.Webhooks = mergeOrderedMaps(dest.Webhooks, src.Webhooks)
	dest.ExternalDocs = mergeExternalDocs(dest.ExternalDocs, src.ExternalDocs)

	if len(src.Security) > 0 {
		dest.Security = src.Security
	}

	if len(src.Servers) > 0 {
		dest.Servers = src.Servers
	}

	if len(src.Tags) > 0 {
		dest.Tags = src.Tags
	}

	mergeOpenAPIv3Paths(dest, src)

	if src.Components == nil {
		return
	}

	if dest.Components == nil {
		dest.Components = &highv3.Components{}
	}

	dest.Components.Callbacks = mergeOrderedMaps(
		dest.Components.Callbacks,
		src.Components.Callbacks,
	)

	dest.Components.Examples = mergeOrderedMaps(
		dest.Components.Examples,
		src.Components.Examples,
	)

	dest.Components.Extensions = mergeOrderedMaps(
		dest.Components.Extensions,
		src.Components.Extensions,
	)

	dest.Components.Headers = mergeOrderedMaps(
		dest.Components.Headers,
		src.Components.Headers,
	)

	dest.Components.Links = mergeOrderedMaps(
		dest.Components.Links,
		src.Components.Links,
	)

	dest.Components.MediaTypes = mergeOrderedMaps(
		dest.Components.MediaTypes,
		src.Components.MediaTypes,
	)

	dest.Components.Parameters = mergeOrderedMaps(
		dest.Components.Parameters,
		src.Components.Parameters,
	)

	dest.Components.PathItems = mergeOrderedMaps(
		dest.Components.PathItems,
		src.Components.PathItems,
	)

	dest.Components.RequestBodies = mergeOrderedMaps(
		dest.Components.RequestBodies,
		src.Components.RequestBodies,
	)

	dest.Components.Responses = mergeOrderedMaps(
		dest.Components.Responses,
		src.Components.Responses,
	)

	dest.Components.Schemas = mergeOrderedMaps(
		dest.Components.Schemas,
		src.Components.Schemas,
	)

	dest.Components.SecuritySchemes = mergeOrderedMaps(
		dest.Components.SecuritySchemes,
		src.Components.SecuritySchemes,
	)
}

func mergeOpenAPIv3Paths(dest, src *highv3.Document) {
	if src.Paths == nil {
		return
	}

	if dest.Paths == nil {
		dest.Paths = src.Paths

		return
	}

	dest.Paths.Extensions = mergeOrderedMaps(
		dest.Paths.Extensions,
		src.Paths.Extensions,
	)

	if src.Paths.PathItems == nil || src.Paths.PathItems.Len() == 0 {
		return
	}

	for iter := src.Paths.PathItems.Oldest(); iter != nil; iter = iter.Next() {
		mergeOpenAPIv3PathItem(dest, iter)
	}
}

func mergeOpenAPIv3PathItem(
	dest *highv3.Document,
	iter *orderedmap.Pair[string, *highv3.PathItem],
) {
	pathItem, present := dest.Paths.PathItems.Get(iter.Key)
	if !present || pathItem == nil {
		dest.Paths.PathItems.Set(iter.Key, iter.Value)

		return
	}

	pathItem.AdditionalOperations = mergeOrderedMaps(
		pathItem.AdditionalOperations,
		iter.Value.AdditionalOperations,
	)

	pathItem.Extensions = mergeOrderedMaps(
		pathItem.Extensions,
		iter.Value.Extensions,
	)

	if iter.Value.Description != "" {
		pathItem.Description = iter.Value.Description
	}

	if len(iter.Value.Servers) > 0 {
		pathItem.Servers = iter.Value.Servers
	}

	pathItem.Delete = mergeOpenAPIv3Operation(
		pathItem.Delete,
		iter.Value.Delete,
	)

	pathItem.Get = mergeOpenAPIv3Operation(
		pathItem.Get,
		iter.Value.Get,
	)

	pathItem.Head = mergeOpenAPIv3Operation(
		pathItem.Head,
		iter.Value.Head,
	)

	pathItem.Options = mergeOpenAPIv3Operation(
		pathItem.Options,
		iter.Value.Options,
	)

	pathItem.Patch = mergeOpenAPIv3Operation(
		pathItem.Patch,
		iter.Value.Patch,
	)

	pathItem.Post = mergeOpenAPIv3Operation(
		pathItem.Post,
		iter.Value.Post,
	)

	pathItem.Put = mergeOpenAPIv3Operation(
		pathItem.Put,
		iter.Value.Put,
	)

	pathItem.Query = mergeOpenAPIv3Operation(
		pathItem.Query,
		iter.Value.Query,
	)

	pathItem.Trace = mergeOpenAPIv3Operation(
		pathItem.Trace,
		iter.Value.Trace,
	)

	if len(iter.Value.Parameters) > 0 {
		pathItem.Parameters = iter.Value.Parameters
	}

	if len(iter.Value.Reference) > 0 {
		pathItem.Reference = iter.Value.Reference
	}

	if iter.Value.Summary != "" {
		pathItem.Summary = iter.Value.Summary
	}
}

func mergeOpenAPIv3Operation(dest, src *highv3.Operation) *highv3.Operation {
	if src == nil {
		return dest
	}

	if dest == nil {
		return src
	}

	dest.Callbacks = mergeOrderedMaps(dest.Callbacks, src.Callbacks)
	dest.Extensions = mergeOrderedMaps(dest.Extensions, src.Extensions)
	dest.ExternalDocs = mergeExternalDocs(dest.ExternalDocs, src.ExternalDocs)

	if src.Deprecated != nil {
		dest.Deprecated = src.Deprecated
	}

	if src.Description != "" {
		dest.Description = src.Description
	}

	if src.OperationId != "" {
		dest.OperationId = src.OperationId
	}

	if len(src.Parameters) > 0 {
		dest.Parameters = src.Parameters
	}

	if len(src.Security) > 0 {
		dest.Security = src.Security
	}

	if len(src.Servers) > 0 {
		dest.Servers = src.Servers
	}

	if src.Summary != "" {
		dest.Summary = src.Summary
	}

	if len(src.Tags) > 0 {
		dest.Tags = src.Tags
	}

	dest.RequestBody = mergeOpenAPIv3RequestBody(dest.RequestBody, src.RequestBody)
	dest.Responses = mergeOpenAPIv3Responses(dest.Responses, src.Responses)

	return dest
}

func mergeOpenAPIv3RequestBody(dest, src *highv3.RequestBody) *highv3.RequestBody {
	if src == nil {
		return dest
	}

	if dest == nil {
		return src
	}

	dest.Content = mergeOrderedMaps(dest.Content, src.Content)
	dest.Extensions = mergeOrderedMaps(dest.Extensions, src.Extensions)

	if src.Description != "" {
		dest.Description = src.Description
	}

	if src.Reference != "" {
		dest.Reference = src.Reference
	}

	if src.Required != nil {
		dest.Required = src.Required
	}

	return dest
}

func mergeOpenAPIv3Responses(dest, src *highv3.Responses) *highv3.Responses {
	if src == nil {
		return dest
	}

	if dest == nil {
		return src
	}

	dest.Codes = mergeOrderedMaps(dest.Codes, src.Codes)
	dest.Extensions = mergeOrderedMaps(dest.Extensions, src.Extensions)

	if src.Default != nil {
		dest.Default = src.Default
	}

	return dest
}

func mergeExternalDocs(dest, src *base.ExternalDoc) *base.ExternalDoc {
	if src == nil {
		return dest
	}

	if dest == nil {
		return src
	}

	if src.Description != "" {
		dest.Description = src.Description
	}

	if src.URL != "" {
		dest.URL = src.URL
	}

	dest.Extensions = mergeOrderedMaps(
		dest.Extensions,
		src.Extensions,
	)

	return dest
}
