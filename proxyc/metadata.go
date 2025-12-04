package proxyc

import (
	"fmt"

	"github.com/relychan/relyx/proxyc/internal"
	"github.com/relychan/relyx/schema"
)

// BuildMetadataTree builds the metadata tree from the API document.
func BuildMetadataTree(document *schema.RelyProxyAPIDocument) (*internal.Node, error) {
	rootNode := new(internal.Node)

	for pathItem := document.Paths.Oldest(); pathItem != nil; pathItem = pathItem.Next() {
		err := internal.ValidateOperations(pathItem)
		if err != nil {
			return nil, err
		}

		_, err = rootNode.InsertRoute(pathItem.Key, pathItem.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to insert route %s: %w", pathItem.Key, err)
		}
	}

	return rootNode, nil
}
