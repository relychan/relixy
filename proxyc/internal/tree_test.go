package internal

import (
	"net/http"
	"testing"

	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"gotest.tools/v3/assert"
)

func TestTreeNodes(t *testing.T) {
	routes := []struct {
		Path     string
		Pattern  string
		Method   string
		Handlers *highv3.PathItem
		Route    Route
	}{
		{
			Pattern: "/posts",
			Path:    "/posts",
			Method:  http.MethodGet,
			Handlers: &highv3.PathItem{
				Get: &highv3.Operation{},
			},
			Route: Route{
				Node: &Node{
					typ: ntStatic,
				},
				ParamValues: map[string]string{},
			},
		},
		{
			Path:    "/posts/1",
			Pattern: "/posts/{id}",
			Method:  http.MethodPost,
			Handlers: &highv3.PathItem{
				Post: &highv3.Operation{},
			},
			Route: Route{
				Node: &Node{
					typ: ntParam,
				},
				ParamValues: map[string]string{
					"id": "1",
				},
			},
		},
		{
			Path:    "/posts/1/comments/abc",
			Pattern: "/posts/{id}/comments/{commentId}",
			Method:  http.MethodGet,
			Handlers: &highv3.PathItem{
				Get: &highv3.Operation{},
			},
			Route: Route{
				Node: &Node{
					typ: ntParam,
				},
				ParamValues: map[string]string{
					"id":        "1",
					"commentId": "abc",
				},
			},
		},
		{
			Path:    "/v1/random/route",
			Pattern: "/v1/*",
			Method:  http.MethodGet,
			Handlers: &highv3.PathItem{
				Get: &highv3.Operation{},
			},
			Route: Route{
				Node: &Node{
					typ: ntCatchAll,
				},
				ParamValues: map[string]string{},
			},
		},
	}

	node := new(Node)

	for _, route := range routes {
		_, err := node.InsertRoute(route.Pattern, route.Handlers, &proxyhandler.InsertRouteOptions{})
		assert.NilError(t, err, route.Pattern)
	}

	for _, route := range routes {
		t.Run(route.Path, func(t *testing.T) {
			postNode := node.FindRoute(route.Path, route.Method)
			assert.Assert(t, postNode != nil)
			assert.Equal(t, postNode.Node.typ, route.Route.Node.typ)
			assert.DeepEqual(t, postNode.Node.pattern, route.Pattern)
			assert.DeepEqual(t, postNode.ParamValues, route.Route.ParamValues)
			assert.Assert(t, postNode.Node.handlers != nil)
		})
	}

	notFoundNode := node.FindRoute("/posts/1/authors", http.MethodGet)
	assert.Assert(t, notFoundNode == nil)
}
