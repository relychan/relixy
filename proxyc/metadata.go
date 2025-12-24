package proxyc

import (
	"context"
	"fmt"

	"github.com/relychan/relixy/proxyc/internal"
	"github.com/relychan/relixy/schema"
)

// BuildMetadataTree builds the metadata tree from the API document.
func BuildMetadataTree(
	ctx context.Context,
	document *schema.RelyProxyAPIDocument,
	clientOptions *ProxyClientOptions,
) (*internal.Node, error) {
	rootNode := new(internal.Node)
	options := &internal.InsertRouteOptions{}

	if clientOptions != nil && clientOptions.CustomEnvGetter != nil {
		options.GetEnv = clientOptions.CustomEnvGetter(ctx)
	}

	for pathItem := document.Paths.Oldest(); pathItem != nil; pathItem = pathItem.Next() {
		err := internal.ValidateOperations(pathItem)
		if err != nil {
			return nil, err
		}

		_, err = rootNode.InsertRoute(pathItem.Key, pathItem.Value, options)
		if err != nil {
			return nil, fmt.Errorf("failed to insert route %s: %w", pathItem.Key, err)
		}
	}

	return rootNode, nil
}
