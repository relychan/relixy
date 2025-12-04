package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"
	"github.com/relychan/goutils"
	"github.com/relychan/relyx/schema"
)

func main() {
	err := jsonSchemaConfiguration()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelyProxyAPIDocument: %w", err))
	}
}

func jsonSchemaConfiguration() error { //nolint:funlen
	r := new(jsonschema.Reflector)

	err := r.AddGoComments(
		"github.com/relychan/restly-proxy/schema",
		"../schema",
		jsonschema.WithFullComment(),
	)
	if err != nil {
		return err
	}

	reflectSchema := r.Reflect(schema.RelyProxyAPIDocument{})

	for _, externalType := range []any{
		schema.OAuthFlow{},
		schema.RelyProxyPathItem{},
		schema.RelyProxyMediaType{},
		schema.RelyProxyEncoding{},
		schema.RelyProxyHeader{},
		schema.GraphQLVariableDefinition{},
		schema.RelyProxyGraphQLResponseConfig{},
	} {
		externalSchema := r.Reflect(externalType)

		for key, def := range externalSchema.Definitions {
			if _, ok := reflectSchema.Definitions[key]; !ok {
				reflectSchema.Definitions[key] = def
			}
		}
	}

	// custom schema types
	reflectSchema.Definitions["Duration"] = &jsonschema.Schema{
		Type:        "string",
		Description: "Duration string",
		Pattern:     "^((([0-9]+h)?([0-9]+m)?([0-9]+s))|(([0-9]+h)?([0-9]+m))|([0-9]+h))$",
	}
	reflectSchema.Definitions["RelyProxyGraphQLRequestConfig"] = newRelyProxyGraphQLRequestConfigSchema()
	reflectSchema.Definitions["RelyProxyAPIDocument"].Properties.Set("paths", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/RelyProxyPathItem",
		},
	})
	reflectSchema.Definitions["Parameter"].Properties.Set("in", &jsonschema.Schema{
		Type: "string",
		Enum: goutils.ToAnySlice(schema.SupportedParameterLocations()),
	})
	reflectSchema.Definitions["Parameter"].Properties.Set("content", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/RelyProxyMediaType",
		},
	})
	reflectSchema.Definitions["TemplateTransformerConfig"] = &jsonschema.Schema{
		Ref: "https://raw.githubusercontent.com/relychan/gotransform/refs/heads/main/jsonschema/gotransform.schema.json",
	}
	reflectSchema.Definitions["RelyProxyAction"] = newRelyProxyActionSchema()
	reflectSchema.Definitions["Discriminator"].Properties.Set("mapping", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Type: "string",
		},
	})
	reflectSchema.Definitions["RelyProxyEncoding"].Properties.Set("headers", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/RelyProxyHeader",
		},
	})
	reflectSchema.Definitions["RelyProxyMediaType"].Properties.Set("encoding", &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			Ref: "#/$defs/RelyProxyEncoding",
		},
	})
	reflectSchema.Definitions["RelyProxyMediaType"].Properties.
		Set("itemEncoding", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxyEncoding",
			},
		})
	reflectSchema.Definitions["RelyProxySchema"].Properties.
		Set("dependentSchemas", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxySchema",
			},
		})
	reflectSchema.Definitions["RelyProxySchema"].Properties.
		Set("patternProperties", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxySchema",
			},
		})
	reflectSchema.Definitions["RelyProxySchema"].Properties.
		Set("dependentRequired", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Type: "array",
				Items: &jsonschema.Schema{
					Type: "string",
				},
			},
		})
	reflectSchema.Definitions["RelyProxySchema"].Properties.
		Set("properties", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxySchema",
			},
		})
	reflectSchema.Definitions["RelyProxyHeader"].Properties.
		Set("content", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxyMediaType",
			},
		})
	reflectSchema.Definitions["RelyProxyRequestBody"].Properties.
		Set("content", &jsonschema.Schema{
			Type: "object",
			AdditionalProperties: &jsonschema.Schema{
				Ref: "#/$defs/RelyProxyMediaType",
			},
		})
	reflectSchema.Definitions["RelyProxySecurityScheme"] = newRelyProxySecuritySchemeSchema()

	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relyx/schema.RelyProxyPathItem]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relyx/schema.GraphQLVariableDefinition]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relyx/schema.RelyProxyEncoding]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relyx/schema.RelyProxyHeader]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relyx/schema.RelyProxySchema]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relyx/schema.RelyProxyMediaType]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,[]string]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,string]",
	)

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")

	err = enc.Encode(reflectSchema)
	if err != nil {
		return err
	}

	return os.WriteFile( //nolint:gosec
		filepath.Join("docs", "rest_proxy.schema.json"),
		buffer.Bytes(), 0o644,
	)
}
