package internal

import (
	"net/http"
	"testing"

	"github.com/relychan/relixy/schema"
	"gotest.tools/v3/assert"
)

func TestTreeNodes(t *testing.T) {
	routes := []struct {
		Path     string
		Pattern  string
		Method   string
		Handlers *schema.RelixyPathItem
		Route    Route
	}{
		{
			Pattern: "/posts",
			Path:    "/posts",
			Method:  http.MethodGet,
			Handlers: &schema.RelixyPathItem{
				Get: &schema.RelixyOperation{},
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
			Handlers: &schema.RelixyPathItem{
				Post: &schema.RelixyOperation{},
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
			Handlers: &schema.RelixyPathItem{
				Get: &schema.RelixyOperation{},
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
			Handlers: &schema.RelixyPathItem{
				Get: &schema.RelixyOperation{},
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
		_, err := node.InsertRoute(route.Pattern, route.Handlers, &InsertRouteOptions{})
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
