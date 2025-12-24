// Package main generates the JSON schema for the relixy metadata.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/config"
	"github.com/relychan/relixy/schema"
)

func main() {
	err := jsonSchemaConfiguration()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelyProxyAPIDocument: %w", err))
	}

	err = jsonSchemaServerConfiguration()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelixyServerConfig: %w", err))
	}
}

func jsonSchemaConfiguration() error {
	r := new(jsonschema.Reflector)

	err := r.AddGoComments(
		"github.com/relychan/relixy/schema",
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
		"OrderedMap[string,*github.com/relychan/relixy/schema.RelyProxyEncoding]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relixy/schema.RelyProxyPathItem]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relixy/schema.GraphQLVariableDefinition]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relixy/schema.RelyProxyHeader]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relixy/schema.RelyProxySchema]",
	)
	delete(
		reflectSchema.Definitions,
		"OrderedMap[string,*github.com/relychan/relixy/schema.RelyProxyMediaType]",
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
		"relixy.schema.json",
		buffer.Bytes(), 0o644,
	)
}

func jsonSchemaServerConfiguration() error {
	r := new(jsonschema.Reflector)

	err := r.AddGoComments(
		"github.com/relychan/relixy/config",
		"../config",
		jsonschema.WithFullComment(),
	)
	if err != nil {
		return err
	}

	reflectSchema := r.Reflect(config.RelixyServerConfig{})

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")

	err = enc.Encode(reflectSchema)
	if err != nil {
		return err
	}

	return os.WriteFile( //nolint:gosec
		"server.schema.json",
		buffer.Bytes(), 0o644,
	)
}
