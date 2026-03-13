// Package main generates the JSON schema for the relixy metadata.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/relychan/relixy/config"
	"github.com/relychan/relixy/schema"
)

func main() {
	err := genConfigurationSchema()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelixyAPIDocument: %w", err))
	}

	err = genServerConfigurationSchema()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelixyServerConfig: %w", err))
	}
}

func genConfigurationSchema() error {
	r := new(jsonschema.Reflector)

	for _, name := range []string{"/schema/openapi", "/schema/baseschema", "/schema/gqlschema"} {
		err := r.AddGoComments(
			"github.com/relychan/relixy"+name,
			".."+name,
			jsonschema.WithFullComment(),
		)
		if err != nil {
			return err
		}
	}

	reflectSchema := r.Reflect(schema.RelixyResource{})

	// custom schema types
	openapiSchema, err := genOpenAPIResourceSchema()
	if err != nil {
		return fmt.Errorf("failed to write jsonschema for RelixyOpenAPIResource: %w", err)
	}

	maps.Copy(reflectSchema.Definitions, openapiSchema.Definitions)

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

func genServerConfigurationSchema() error {
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

	// custom schema types
	reflectSchema.Definitions["ServerConfig"] = &jsonschema.Schema{
		Description: "Configurations for the HTTP server",
		Ref:         "https://raw.githubusercontent.com/relychan/gohttps/refs/heads/main/jsonschema/server.schema.json",
	}

	reflectSchema.Definitions["OTLPConfig"] = &jsonschema.Schema{
		Description: "Configurations for OpenTelemetry exporters",
		Ref:         "https://raw.githubusercontent.com/hasura/gotel/refs/heads/main/jsonschema/gotel.schema.json",
	}

	// delete unused definitions
	delete(reflectSchema.Definitions, "CORSConfig")

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")

	err = enc.Encode(reflectSchema)
	if err != nil {
		return err
	}

	return os.WriteFile( //nolint:gosec
		"relixy-server.schema.json",
		buffer.Bytes(), 0o644,
	)
}
