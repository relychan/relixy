package internal

import (
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"
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
				ParamValues: map[string]string{
					"id": "1",
				},
			},
		},
		{
			Path:    "/posts/1/comments/abc",
			Pattern: "/posts/{id}/comments/{commentId:^[a-z]+$}",
			Method:  http.MethodGet,
			Handlers: &highv3.PathItem{
				Get: &highv3.Operation{},
			},
			Route: Route{
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
				ParamValues: map[string]string{},
			},
		},
	}

	node := new(Node)

	for _, route := range routes {
		_, err := node.InsertRoute(route.Pattern, route.Handlers, &proxyhandler.InsertRouteOptions{})
		assert.NilError(t, err, route.Pattern)
	}

	routeAsText := `
- / []
  - /posts [GET]
    - /{id} []
      - / [POST]
      - /comments []
        - /{commentId:^[a-z]+$} []
          - / [GET]
  - /v1 []
    - /* [GET]`

	assert.Equal(t, routeAsText, node.printDebug(0))
	for _, route := range routes {
		if !t.Run(route.Path, func(t *testing.T) {
			postNode := node.FindRoute(route.Path, route.Method)
			assert.Assert(t, postNode != nil)
			assert.Equal(t, postNode.Pattern, route.Pattern)
			assert.DeepEqual(t, postNode.ParamValues, route.Route.ParamValues)
		}) {
			break
		}
	}

	notFoundNode := node.FindRoute("/posts/1/authors", http.MethodGet)
	assert.Assert(t, notFoundNode == nil)
}

// TestRouteInsertionEdgeCases tests edge cases in route insertion
func TestRouteInsertionEdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		patterns    []string
		expectError bool
		errorType   error
	}{
		{
			name: "overlapping_static_routes",
			patterns: []string{
				"/posts",
				"/posts/new",
				"/posts/123",
			},
			expectError: false,
		},
		{
			name: "param_and_static_mix",
			patterns: []string{
				"/posts/{id}",
				"/posts/new",
			},
			expectError: false,
		},
		{
			name: "multiple_params_same_level",
			patterns: []string{
				"/posts/{id}",
				"/posts/{postId}",
			},
			expectError: false,
		},
		{
			name: "regexp_patterns",
			patterns: []string{
				"/posts/{id:[0-9]+}",
				"/posts/{slug:[a-z-]+}",
			},
			expectError: false,
		},
		{
			name: "nested_params",
			patterns: []string{
				"/users/{userId}/posts/{postId}",
				"/users/{userId}/comments/{commentId}",
			},
			expectError: false,
		},
		{
			name: "catchall_routes",
			patterns: []string{
				"/api/v1/*",
				"/api/v2/*",
			},
			expectError: false,
		},
		{
			name: "root_route",
			patterns: []string{
				"/",
			},
			expectError: false,
		},
		{
			name: "deep_nesting",
			patterns: []string{
				"/a/b/c/d/e/f/g",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := new(Node)

			for _, pattern := range tc.patterns {
				_, err := node.InsertRoute(pattern, &highv3.PathItem{
					Get: &highv3.Operation{},
				}, &proxyhandler.InsertRouteOptions{})

				if tc.expectError {
					assert.Assert(t, err != nil, "expected error for pattern: %s", pattern)
					if tc.errorType != nil {
						assert.ErrorIs(t, err, tc.errorType)
					}
					return
				}

				assert.NilError(t, err, "failed to insert pattern: %s", pattern)
			}
		})
	}
}

// TestRouteFindingEdgeCases tests edge cases in route finding
func TestRouteFindingEdgeCases(t *testing.T) {
	node := new(Node)

	// Setup routes
	routes := map[string]*highv3.PathItem{
		"/":                                  {Get: &highv3.Operation{}},
		"/posts":                             {Get: &highv3.Operation{}},
		"/posts/new":                         {Get: &highv3.Operation{}},
		"/posts/{id}":                        {Get: &highv3.Operation{}, Post: &highv3.Operation{}},
		"/posts/{id:[0-9]+}":                 {Put: &highv3.Operation{}},
		"/posts/{id}/comments":               {Get: &highv3.Operation{}},
		"/posts/{id}/comments/{commentId}":   {Get: &highv3.Operation{}},
		"/users/{userId}/posts/{postId}":     {Get: &highv3.Operation{}},
		"/api/v1/*":                          {Get: &highv3.Operation{}},
		"/products/{category}/{subcategory}": {Get: &highv3.Operation{}},
		"/products/{category}/{id:[0-9]+}":   {Get: &highv3.Operation{}},
	}

	for pattern, handlers := range routes {
		_, err := node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
		assert.NilError(t, err, "failed to insert route: %s", pattern)
	}

	testCases := []struct {
		name            string
		path            string
		method          string
		shouldFind      bool
		expectedParams  map[string]string
		expectedPattern string
	}{
		{
			name:            "root_path",
			path:            "/",
			method:          http.MethodGet,
			shouldFind:      true,
			expectedParams:  map[string]string{},
			expectedPattern: "/",
		},
		{
			name:            "static_exact_match",
			path:            "/posts",
			method:          http.MethodGet,
			shouldFind:      true,
			expectedParams:  map[string]string{},
			expectedPattern: "/posts",
		},
		{
			name:            "static_priority_over_param",
			path:            "/posts/new",
			method:          http.MethodGet,
			shouldFind:      true,
			expectedParams:  map[string]string{},
			expectedPattern: "/posts/new",
		},
		{
			name:       "regexp_match_numeric",
			path:       "/posts/456",
			method:     http.MethodPut,
			shouldFind: true,
			expectedParams: map[string]string{
				"id": "456",
			},
			expectedPattern: "/posts/{id:^[0-9]+$}",
		},
		{
			name:       "nested_params",
			path:       "/posts/123/comments/456",
			method:     http.MethodGet,
			shouldFind: true,
			expectedParams: map[string]string{
				"id":        "123",
				"commentId": "456",
			},
			expectedPattern: "/posts/{id}/comments/{commentId}",
		},
		{
			name:       "multiple_params_different_paths",
			path:       "/users/user123/posts/post456",
			method:     http.MethodGet,
			shouldFind: true,
			expectedParams: map[string]string{
				"userId": "user123",
				"postId": "post456",
			},
			expectedPattern: "/users/{userId}/posts/{postId}",
		},
		{
			name:            "catchall_match",
			path:            "/api/v1/anything/goes/here",
			method:          http.MethodGet,
			shouldFind:      true,
			expectedParams:  map[string]string{},
			expectedPattern: "/api/v1/*",
		},
		{
			name:       "multiple_params_same_segment",
			path:       "/products/electronics/smartphones",
			method:     http.MethodGet,
			shouldFind: true,
			expectedParams: map[string]string{
				"category":    "electronics",
				"subcategory": "smartphones",
			},
			expectedPattern: "/products/{category}/{subcategory}",
		},
		{
			name:       "regexp_priority",
			path:       "/products/electronics/12345",
			method:     http.MethodGet,
			shouldFind: true,
			expectedParams: map[string]string{
				"category": "electronics",
				"id":       "12345",
			},
			expectedPattern: "/products/{category}/{id:^[0-9]+$}",
		},
		{
			name:       "method_not_found",
			path:       "/posts/123",
			method:     http.MethodDelete,
			shouldFind: false,
		},
		{
			name:       "path_not_found",
			path:       "/nonexistent",
			method:     http.MethodGet,
			shouldFind: false,
		},
		{
			name:       "partial_path_not_found",
			path:       "/posts/123/nonexistent",
			method:     http.MethodGet,
			shouldFind: false,
		},
		{
			name:       "double_slash_path",
			path:       "/posts//comments",
			method:     http.MethodGet,
			shouldFind: true, // Router matches this with empty param
			expectedParams: map[string]string{
				"id": "", // Empty param value
			},
			expectedPattern: "/posts/{id}/comments",
		},
		{
			name:            "trailing_slash_mismatch",
			path:            "/posts/",
			method:          http.MethodGet,
			shouldFind:      true,
			expectedPattern: "/posts",
			expectedParams:  map[string]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			route := node.FindRoute(tc.path, tc.method)

			if tc.shouldFind {
				assert.Assert(t, route != nil, "expected to find route for path: %s", tc.path)
				assert.Assert(t, route.Method != nil)
				assert.Equal(t, tc.expectedPattern, route.Pattern)
				assert.DeepEqual(t, tc.expectedParams, route.ParamValues)
			} else {
				assert.Assert(t, route == nil, "expected not to find route for path: %s", tc.path)
			}
		})
	}
}

// TestComplexRoutingScenarios tests complex real-world routing scenarios
func TestComplexRoutingScenarios(t *testing.T) {
	t.Run("RESTful_API_with_versioning", func(t *testing.T) {
		node := new(Node)

		routes := []string{
			"/api/v1/users",
			"/api/v1/users/{id}",
			"/api/v1/users/{id}/posts",
			"/api/v1/users/{id}/posts/{postId}",
			"/api/v2/users",
			"/api/v2/users/{id}",
		}

		for _, route := range routes {
			_, err := node.InsertRoute(route, &highv3.PathItem{
				Get: &highv3.Operation{},
			}, &proxyhandler.InsertRouteOptions{})
			assert.NilError(t, err)
		}

		// Test v1 routes
		r := node.FindRoute("/api/v1/users", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "/api/v1/users", r.Pattern)

		r = node.FindRoute("/api/v1/users/123", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "123", r.ParamValues["id"])

		r = node.FindRoute("/api/v1/users/123/posts/456", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "123", r.ParamValues["id"])
		assert.Equal(t, "456", r.ParamValues["postId"])

		// Test v2 routes
		r = node.FindRoute("/api/v2/users/789", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "/api/v2/users/{id}", r.Pattern)
	})

	t.Run("mixed_static_and_dynamic_routes", func(t *testing.T) {
		node := new(Node)

		routes := map[string]*highv3.PathItem{
			"/posts":           {Get: &highv3.Operation{}},
			"/posts/new":       {Get: &highv3.Operation{}},
			"/posts/popular":   {Get: &highv3.Operation{}},
			"/posts/{id}":      {Get: &highv3.Operation{}},
			"/posts/{id}/edit": {Get: &highv3.Operation{}},
		}

		for pattern, handlers := range routes {
			_, err := node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
			assert.NilError(t, err)
		}

		// Static routes should take precedence
		r := node.FindRoute("/posts/new", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "/posts/new", r.Pattern)

		r = node.FindRoute("/posts/popular", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "/posts/popular", r.Pattern)

		// Dynamic route should match other IDs
		r = node.FindRoute("/posts/123", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "/posts/{id}", r.Pattern)
		assert.Equal(t, "123", r.ParamValues["id"])

		r = node.FindRoute("/posts/456/edit", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "/posts/{id}/edit", r.Pattern)
		assert.Equal(t, "456", r.ParamValues["id"])
	})

	t.Run("regexp_validation_routes", func(t *testing.T) {
		node := new(Node)

		routes := map[string]*highv3.PathItem{
			"/users/{id:[0-9]+}":                       {Get: &highv3.Operation{}},
			"/users/{username:[a-z]+}":                 {Post: &highv3.Operation{}},
			"/posts/{slug:[a-z0-9-]+}":                 {Get: &highv3.Operation{}},
			"/files/{filename:[a-zA-Z0-9._-]+}":        {Get: &highv3.Operation{}},
			"/dates/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}": {Get: &highv3.Operation{}},
		}

		for pattern, handlers := range routes {
			_, err := node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
			assert.NilError(t, err)
		}

		// Numeric ID should match
		r := node.FindRoute("/users/12345", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "12345", r.ParamValues["id"])

		// Alphabetic username should match
		r = node.FindRoute("/users/john", http.MethodPost)
		assert.Assert(t, r != nil)
		assert.Equal(t, "john", r.ParamValues["username"])

		// Slug with hyphens should match
		r = node.FindRoute("/posts/my-blog-post-123", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "my-blog-post-123", r.ParamValues["slug"])

		// Filename with dots and underscores should match
		r = node.FindRoute("/files/my_file.txt", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "my_file.txt", r.ParamValues["filename"])

		// Date format should match
		r = node.FindRoute("/dates/2024-01-15", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "2024-01-15", r.ParamValues["date"])

		// Invalid date format should not match
		r = node.FindRoute("/dates/2024-1-5", http.MethodGet)
		assert.Assert(t, r == nil)
	})

	t.Run("catchall_with_specific_routes", func(t *testing.T) {
		node := new(Node)

		routes := map[string]*highv3.PathItem{
			"/api/v1/users": {Get: &highv3.Operation{}},
			"/api/v1/posts": {Get: &highv3.Operation{}},
			"/api/v1/*":     {Get: &highv3.Operation{}},
			"/static/css/*": {Get: &highv3.Operation{}},
			"/static/js/*":  {Get: &highv3.Operation{}},
		}

		for pattern, handlers := range routes {
			_, err := node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
			assert.NilError(t, err)
		}

		// Specific routes should match first
		r := node.FindRoute("/api/v1/users", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "/api/v1/users", r.Pattern)

		// Catchall should match unspecified routes
		r = node.FindRoute("/api/v1/anything/else", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "/api/v1/*", r.Pattern)

		r = node.FindRoute("/static/css/main.css", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "/static/css/*", r.Pattern)
	})

	t.Run("multiple_HTTP_methods", func(t *testing.T) {
		node := new(Node)

		_, err := node.InsertRoute("/posts/{id}", &highv3.PathItem{
			Get:    &highv3.Operation{},
			Post:   &highv3.Operation{},
			Put:    &highv3.Operation{},
			Patch:  &highv3.Operation{},
			Delete: &highv3.Operation{},
		}, &proxyhandler.InsertRouteOptions{})
		assert.NilError(t, err)

		methods := []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		}

		for _, method := range methods {
			r := node.FindRoute("/posts/123", method)
			assert.Assert(t, r != nil, "method %s should be found", method)
			assert.Equal(t, "123", r.ParamValues["id"])
		}

		// Method not defined should not be found
		r := node.FindRoute("/posts/123", http.MethodHead)
		assert.Assert(t, r == nil)
	})

	t.Run("deeply_nested_resources", func(t *testing.T) {
		node := new(Node)

		pattern := "/orgs/{orgId}/teams/{teamId}/projects/{projectId}/tasks/{taskId}/comments/{commentId}"
		_, err := node.InsertRoute(pattern, &highv3.PathItem{
			Get: &highv3.Operation{},
		}, &proxyhandler.InsertRouteOptions{})
		assert.NilError(t, err)

		r := node.FindRoute("/orgs/org1/teams/team2/projects/proj3/tasks/task4/comments/comment5", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "org1", r.ParamValues["orgId"])
		assert.Equal(t, "team2", r.ParamValues["teamId"])
		assert.Equal(t, "proj3", r.ParamValues["projectId"])
		assert.Equal(t, "task4", r.ParamValues["taskId"])
		assert.Equal(t, "comment5", r.ParamValues["commentId"])
	})

	t.Run("special_characters_in_params", func(t *testing.T) {
		node := new(Node)

		routes := map[string]*highv3.PathItem{
			"/search/{query}":          {Get: &highv3.Operation{}},
			"/users/{email:[-.@\\w]+}": {Get: &highv3.Operation{}},
		}

		for pattern, handlers := range routes {
			_, err := node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
			assert.NilError(t, err)
		}

		// Query with special characters
		r := node.FindRoute("/search/hello-world", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "hello-world", r.ParamValues["query"])

		// Email-like parameter
		r = node.FindRoute("/users/user@example.com", http.MethodGet)
		assert.Assert(t, r != nil)
		assert.Equal(t, "user@example.com", r.ParamValues["email"])
	})
}

func (n Node) printDebug(indent int) string {
	var sb strings.Builder

	sb.WriteByte('\n')

	if indent > 0 {
		sb.WriteString(strings.Repeat(" ", indent))
	}

	sb.WriteString("- /")
	sb.WriteString(n.String())

	sb.WriteString(fmt.Sprintf(" %v", slices.Collect(maps.Keys(n.handlers))))

	for _, child := range n.children {
		for _, node := range child {
			sb.WriteString(node.printDebug(indent + 2))
		}
	}

	return sb.String()
}

// BenchmarkTree/insert_routes-11         	  234687	      5104 ns/op	   16232 B/op	     167 allocs/op
// BenchmarkTree/find_route-11            	 4148186	       288.3 ns/op	     456 B/op	       6 allocs/op
func BenchmarkTree(b *testing.B) {
	routes := map[string]*highv3.PathItem{
		"/posts":                   {Get: &highv3.Operation{}},
		"/posts/new":               {Get: &highv3.Operation{}},
		"/posts/popular":           {Get: &highv3.Operation{}},
		"/posts/{id}":              {Get: &highv3.Operation{}},
		"/posts/{id}/edit":         {Get: &highv3.Operation{}},
		"/users/{id:[0-9]+}":       {Get: &highv3.Operation{}},
		"/posts/{slug:[a-z0-9-]+}": {Get: &highv3.Operation{}},
	}

	b.Run("insert_routes", func(b *testing.B) {
		for b.Loop() {
			node := new(Node)

			for pattern, handlers := range routes {
				node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
			}
		}
	})

	b.Run("find_route", func(b *testing.B) {
		node := new(Node)

		for pattern, handlers := range routes {
			_, err := node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
			if err != nil {
				b.Fatal(err)
			}
		}

		for b.Loop() {
			node.FindRoute("/posts/hello", http.MethodGet)
		}
	})
}
