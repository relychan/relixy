package proxyc

import (
	"context"

	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/relixy/proxyc/internal"
)

// BuildMetadataTree builds the metadata tree from the API document.
func BuildMetadataTree(
	ctx context.Context,
	document *highv3.Document,
	clientOptions *ProxyClientOptions,
) (*internal.Node, error) {
	rootNode := new(internal.Node)

	if document.Paths.PathItems == nil {
		return rootNode, nil
	}

	options := &internal.InsertRouteOptions{}

	if clientOptions != nil && clientOptions.CustomEnvGetter != nil {
		options.GetEnv = clientOptions.CustomEnvGetter(ctx)
	}

	for pathItem := document.Paths.PathItems.Oldest(); pathItem != nil; pathItem = pathItem.Next() {
		// TODO: use native path item
		// err := internal.ValidateOperations(pathItem)
		// if err != nil {
		// 	return nil, err
		// }

		// _, err = rootNode.InsertRoute(pathItem.Key, pathItem.Value, options)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to insert route %s: %w", pathItem.Key, err)
		// }
	}

	return rootNode, nil
}
