package proxyc

import (
	"fmt"

	"github.com/hasura/goenvconf"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/gohttpc"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"github.com/relychan/relixy/proxyc/internal"
)

// BuildMetadataTree builds the metadata tree from the API document.
func BuildMetadataTree(
	document *highv3.Document,
	clientOptions *gohttpc.ClientOptions,
) (*internal.Node, error) {
	rootNode := new(internal.Node)

	if document.Paths.PathItems == nil {
		return rootNode, nil
	}

	options := &proxyhandler.InsertRouteOptions{
		GetEnv: goenvconf.GetOSEnv,
	}

	if clientOptions != nil && clientOptions.GetEnv != nil {
		options.GetEnv = clientOptions.GetEnv
	}

	for pathItem := document.Paths.PathItems.Oldest(); pathItem != nil; pathItem = pathItem.Next() {
		_, err := rootNode.InsertRoute(pathItem.Key, pathItem.Value, options)
		if err != nil {
			return nil, fmt.Errorf("failed to insert route %s: %w", pathItem.Key, err)
		}
	}

	return rootNode, nil
}
