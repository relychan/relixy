package schema

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// RelixyDocumentToOAS3 converts a proxy document to OAS 3.
func RelixyDocumentToOAS3(doc RelixyAPIDocument) *highv3.Document {
	result := &highv3.Document{
		Servers: make([]*highv3.Server, len(doc.Servers)),
		Tags:    make([]*base.Tag, len(doc.Tags)),
	}

	for i, server := range doc.Servers {
		result.Servers[i] = server.OAS3()
	}

	for i, tag := range doc.Tags {
		result.Tags[i] = tag.OAS()
	}

	return result
}
