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

func jsonSchemaConfiguration() error {
	r := new(jsonschema.Reflector)

	err := r.AddGoComments(
		"github.com/relychan/relyx/schema",
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
		schema.RelyProxyGraphQLRequestConfig{},
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
		MinLength:   goutils.ToPtr(uint64(2)),
		Pattern:     `^(\d+(\.\d+)?h)?(\d+(\.\d+)?m)?(\d+(\.\d+)?s)?(\d+(\.\d+)?ms)?$`,
	}

	// delete unused types
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relyx/schema.RelyProxyEncoding]",
	)
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
		filepath.Join("..", "docs", "relyx.schema.json"),
		buffer.Bytes(), 0o644,
	)
}
