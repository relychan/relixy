// Package main generates the JSON schema for the relixy metadata.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/relychan/goutils"
	"github.com/relychan/relixy/config"
	"github.com/relychan/relixy/schema"
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
	openapiSchema := r.Reflect(openapi.RelixyOpenAPIResource{})
	maps.Copy(reflectSchema.Definitions, openapiSchema.Definitions)

	httpcSchema := downloadGoHTTPCSchema()
	maps.Copy(reflectSchema.Definitions, httpcSchema.Definitions)

	// delete unused definitions
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

func downloadGoHTTPCSchema() *jsonschema.Schema {
	rawResp, err := http.Get( //nolint:bodyclose,noctx
		"https://raw.githubusercontent.com/relychan/gohttpc/refs/heads/main/jsonschema/gohttpc.schema.json",
	)
	if err != nil {
		panic(err)
	}

	if rawResp != nil && rawResp.Body != nil {
		defer goutils.CatchWarnErrorFunc(rawResp.Body.Close)
	}

	if rawResp.StatusCode != http.StatusOK {
		rawBody, err := io.ReadAll(rawResp.Body)
		if err != nil {
			panic("failed to download gohttpc schema: " + rawResp.Status)
		}

		panic("failed to download gohttpc schema" + string(rawBody))
	}

	jsonSchema := new(jsonschema.Schema)

	err = json.NewDecoder(rawResp.Body).Decode(jsonSchema)
	if err != nil {
		panic(fmt.Errorf("failed to decode gohttpc schema: %w", err))
	}

	return jsonSchema
}
