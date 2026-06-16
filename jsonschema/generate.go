// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main generates the JSON schema for the rely metadata.
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
	err := genSchema()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelyResource: %w", err))
	}

	err = genServerConfigurationSchema()
	if err != nil {
		panic(fmt.Errorf("failed to write jsonschema for RelyServerConfig: %w", err))
	}
}

func genSchema() error {
	r := new(jsonschema.Reflector)

	for _, name := range []string{"/schema"} {
		err := r.AddGoComments(
			"github.com/relychan/relixy"+name,
			".."+name,
			jsonschema.WithFullComment(),
		)
		if err != nil {
			return err
		}
	}

	reflectSchema := r.Reflect(schema.RelyResource{})

	// custom schema types
	openapiSchema, err := genOpenAPIResourceSchema()
	if err != nil {
		return fmt.Errorf("failed to write jsonschema for RelyResource: %w", err)
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
		"rely-resources.schema.json",
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

	reflectSchema := r.Reflect(config.RelyServerConfig{})

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
		"rely-server.schema.json",
		buffer.Bytes(), 0o644,
	)
}
