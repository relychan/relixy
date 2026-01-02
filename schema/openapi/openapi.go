// Package openapi defines metadata schema for OpenAPI handlers.
package openapi

import (
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// RelixyDocumentToOAS3 converts a proxy document to OAS 3.
func RelixyDocumentToOAS3(doc RelixyOpenAPI3ResourceSpecification) *highv3.Document {
	result := &highv3.Document{
		Servers:           make([]*highv3.Server, len(doc.Servers)),
		Tags:              doc.Tags,
		Version:           doc.Version,
		Info:              doc.Info,
		Components:        doc.Components,
		ExternalDocs:      doc.ExternalDocs,
		JsonSchemaDialect: doc.JsonSchemaDialect,
		// Security: doc.Security,
	}

	for i, server := range doc.Servers {
		result.Servers[i] = server.OAS3()
	}

	return result
}
