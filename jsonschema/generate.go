// Package main generates the JSON schema for the relixy metadata.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"

	// "github.com/relychan/relixy/config"

	"github.com/relychan/goutils"
	"github.com/relychan/relixy/schema/openapi"
)

func main() {
	err := jsonSchemaConfiguration()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelixyAPIDocument: %w", err))
	}

	// err = jsonSchemaServerConfiguration()
	// if err != nil {
	// 	panic(fmt.Errorf("failed to write jsonschema for RelixyServerConfig: %w", err))
	// }
}

func jsonSchemaConfiguration() error {
	r := new(jsonschema.Reflector)

	for _, name := range []string{"/schema/openapi", "/schema/base_schema"} {
		err := r.AddGoComments(
			"github.com/relychan/relixy"+name,
			"../schema",
			jsonschema.WithFullComment(),
		)
		if err != nil {
			return err
		}
	}

	reflectSchema := r.Reflect(openapi.RelixyOpenAPIv3Resource{})

	// for _, externalType := range []any{
	// 	openapi.RelixyOpenAPIv3PathItem{},
	// 	openapi.RelixyMediaType{},
	// 	openapi.RelixyEncoding{},
	// 	openapi.RelixyHeader{},
	// 	base_schema.GraphQLVariableDefinition{},
	// 	base_schema.RelixyGraphQLRequestConfig{},
	// 	base_schema.RelixyGraphQLResponseConfig{},
	// 	openapi.RelixyAPIKeyAuthConfig{},
	// 	openapi.RelixyHTTPAuthConfig{},
	// 	openapi.RelixyBasicAuthConfig{},
	// 	openapi.RelixyOAuth2Config{},
	// 	openapi.RelixyOpenIDConnectConfig{},
	// 	openapi.RelixyCookieAuthConfig{},
	// 	openapi.RelixyMutualTLSAuthConfig{},
	// } {
	// 	externalSchema := r.Reflect(externalType)

	// 	for key, def := range externalSchema.Definitions {
	// 		if _, ok := reflectSchema.Definitions[key]; !ok {
	// 			reflectSchema.Definitions[key] = def
	// 		}
	// 	}
	// }

	// custom schema types
	reflectSchema.Definitions["Duration"] = &jsonschema.Schema{
		Type:        "string",
		Description: "Duration string",
		MinLength:   goutils.ToPtr(uint64(2)),
		Pattern:     `^(\d+(\.\d+)?h)?(\d+(\.\d+)?m)?(\d+(\.\d+)?s)?(\d+(\.\d+)?ms)?$`,
	}

	// reflectSchema.Definitions["TemplateTransformerConfig"] = &jsonschema.Schema{
	// 	Ref: "https://raw.githubusercontent.com/relychan/gotransform/refs/heads/main/jsonschema/gotransform.schema.json",
	// }

	// reflectSchema.Definitions["HTTPClientConfig"] = &jsonschema.Schema{
	// 	Ref: "https://raw.githubusercontent.com/relychan/gohttpc/refs/heads/main/jsonschema/gohttpc.schema.json",
	// }

	// reflectSchema.Definitions["TLSConfig"] = &jsonschema.Schema{
	// 	Ref: "https://raw.githubusercontent.com/relychan/gohttpc/refs/heads/main/jsonschema/gohttpc.schema.json#/$defs/TLSConfig",
	// }

	// reflectSchema.Definitions["HTTPHealthCheckConfig"] = &jsonschema.Schema{
	// 	Ref: "https://raw.githubusercontent.com/relychan/gohttpc/refs/heads/main/jsonschema/gohttpc.schema.json#/$defs/HTTPHealthCheckConfig",
	// }

	// delete unused types
	// delete(
	// 	reflectSchema.Definitions,
	// 	"OrderedMap[string,*github.com/relychan/relixy/schema.RelixyEncoding]",
	// )
	// delete(
	// 	reflectSchema.Definitions,
	// 	"OrderedMap[string,*github.com/relychan/relixy/schema.RelixyPathItem]",
	// )
	// delete(
	// 	reflectSchema.Definitions,
	// 	"OrderedMap[string,*github.com/relychan/relixy/schema.GraphQLVariableDefinition]",
	// )
	// delete(
	// 	reflectSchema.Definitions,
	// 	"OrderedMap[string,*github.com/relychan/relixy/schema.RelixyHeader]",
	// )
	// delete(
	// 	reflectSchema.Definitions,
	// 	"OrderedMap[string,*github.com/relychan/relixy/schema.RelixySchema]",
	// )
	// delete(
	// 	reflectSchema.Definitions,
	// 	"OrderedMap[string,*github.com/relychan/relixy/schema.RelixyMediaType]",
	// )
	// delete(reflectSchema.Definitions, "OrderedMap[string,[]string]")
	// delete(reflectSchema.Definitions, "OrderedMap[string,string]")
	// delete(reflectSchema.Definitions, "HTTPTransportConfig")
	// delete(reflectSchema.Definitions, "HTTPRetryConfig")
	// delete(reflectSchema.Definitions, "TLSClientCertificate")
	// delete(reflectSchema.Definitions, "HTTPDialerConfig")
	// delete(reflectSchema.Definitions, "HTTPClientAuthConfig")

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")

	err := enc.Encode(reflectSchema)
	if err != nil {
		return err
	}

	return os.WriteFile( //nolint:gosec
		"relixy-openapi.schema.json",
		buffer.Bytes(), 0o644,
	)
}

// func jsonSchemaServerConfiguration() error {
// 	r := new(jsonschema.Reflector)

// 	err := r.AddGoComments(
// 		"github.com/relychan/relixy/config",
// 		"../config",
// 		jsonschema.WithFullComment(),
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	reflectSchema := r.Reflect(config.RelixyServerConfig{})

// 	buffer := new(bytes.Buffer)
// 	enc := json.NewEncoder(buffer)
// 	enc.SetEscapeHTML(false)
// 	enc.SetIndent("", " ")

// 	err = enc.Encode(reflectSchema)
// 	if err != nil {
// 		return err
// 	}

// 	return os.WriteFile( //nolint:gosec
// 		"relixy-server.schema.json",
// 		buffer.Bytes(), 0o644,
// 	)
// }
