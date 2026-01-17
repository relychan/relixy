// Package main generates the JSON schema for the relixy metadata.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/relychan/relixy/config"
	"github.com/relychan/relixy/schema/openapi"
)

func main() {
	err := genConfigurationSchema()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelixyAPIDocument: %w", err))
	}

	err = genRelixyActionSchema()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelixyAction: %w", err))
	}

	err = genServerConfigurationSchema()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelixyServerConfig: %w", err))
	}
}

func genConfigurationSchema() error {
	r := new(jsonschema.Reflector)

	for _, name := range []string{"/schema/openapi", "/schema/base_schema"} {
		err := r.AddGoComments(
			"github.com/relychan/relixy"+name,
			".."+name,
			jsonschema.WithFullComment(),
		)
		if err != nil {
			return err
		}
	}

	reflectSchema := r.Reflect(openapi.RelixyOpenAPIResource{})

	// custom schema types
	reflectSchema.Definitions["HTTPClientConfig"] = &jsonschema.Schema{
		Ref: "https://raw.githubusercontent.com/relychan/gohttpc/refs/heads/main/jsonschema/gohttpc.schema.json",
	}

	// delete unused definitions
	delete(reflectSchema.Definitions, "HTTPTransportConfig")
	delete(reflectSchema.Definitions, "HTTPRetryConfig")
	delete(reflectSchema.Definitions, "HTTPClientAuthConfig")
	delete(reflectSchema.Definitions, "TLSClientCertificate")
	delete(reflectSchema.Definitions, "HTTPDialerConfig")
	delete(reflectSchema.Definitions, "Document")
	delete(reflectSchema.Definitions, "Contact")
	delete(reflectSchema.Definitions, "Components")
	delete(reflectSchema.Definitions, "ExternalDoc")
	delete(reflectSchema.Definitions, "Tag")
	delete(reflectSchema.Definitions, "SecurityRequirement")
	delete(reflectSchema.Definitions, "Server")
	delete(reflectSchema.Definitions, "Paths")
	delete(reflectSchema.Definitions, "Info")
	delete(reflectSchema.Definitions, "License")
	delete(reflectSchema.Definitions, "EnvInt")
	delete(reflectSchema.Definitions, "TLSConfig")

	for key := range reflectSchema.Definitions {
		if strings.HasPrefix(key, "Map[") {
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
