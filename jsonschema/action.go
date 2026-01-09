package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/relychan/relixy/proxyc/handler/graphqlhandler"
	"github.com/relychan/relixy/proxyc/handler/resthandler"
)

type RelixyActionConfig struct{}

// JSONSchema defines a custom definition for JSON schema.
func (RelixyActionConfig) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Description: "Proxy configuration to the remote REST service",
				Ref:         "#/$defs/RelixyRESTActionConfig",
			},
			{
				Description: "Configurations for proxying request to the remote GraphQL server",
				Ref:         "#/$defs/RelixyGraphQLActionConfig",
			},
		},
	}
}

func genRelixyActionSchema() error {
	r := new(jsonschema.Reflector)

	for _, name := range []string{
		"proxyc/handler/graphqlhandler",
		"proxyc/handler/proxyhandler",
		"proxyc/handler/resthandler",
	} {
		err := r.AddGoComments(
			"github.com/relychan/relixy/"+name,
			"../"+name,
			jsonschema.WithFullComment(),
		)
		if err != nil {
			return err
		}
	}

	reflectSchema := r.Reflect(RelixyActionConfig{})

	for _, externalType := range []any{
		graphqlhandler.RelixyGraphQLActionConfig{},
		resthandler.RelixyRESTActionConfig{},
	} {
		externalSchema := r.Reflect(externalType)

		for key, def := range externalSchema.Definitions {
			if _, ok := reflectSchema.Definitions[key]; !ok {
				reflectSchema.Definitions[key] = def
			}
		}
	}

	reflectSchema.Definitions["TemplateTransformerConfig"] = &jsonschema.Schema{
		Ref: "https://raw.githubusercontent.com/relychan/gotransform/refs/heads/main/jsonschema/gotransform.schema.json",
	}

	for key := range reflectSchema.Definitions {
		if strings.HasPrefix(key, "OrderedMap[") {
			delete(reflectSchema.Definitions, key)
		}
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")

	err := enc.Encode(reflectSchema)
	if err != nil {
		return err
	}

	return os.WriteFile( //nolint:gosec
		"relixy-action.schema.json",
		buffer.Bytes(), 0o644,
	)
}
