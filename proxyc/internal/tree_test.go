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

// TestPatNextSegment tests the pattern segment parsing with various edge cases
// func TestPatNextSegment(t *testing.T) {
// 	testCases := []struct {
// 		name          string
// 		pattern       string
// 		expectedTyp   nodeTyp
// 		expectedKey   string
// 		expectedRegex string
// 		expectedTail  byte
// 		expectedStart int
// 		expectedEnd   int
// 		expectError   bool
// 		errorType     error
// 	}{
// 		{
// 			name:          "static_path",
// 			pattern:       "/posts",
// 			expectedTyp:   ntStatic,
// 			expectedKey:   "",
// 			expectedRegex: "",
// 			expectedTail:  0,
// 			expectedStart: 0,
// 			expectedEnd:   6,
// 			expectError:   false,
// 		},
// 		{
// 			name:          "param_simple",
// 			pattern:       "{id}",
// 			expectedTyp:   ntParam,
// 			expectedKey:   "id",
// 			expectedRegex: "",
// 			expectedTail:  '/', // default tail is '/'
// 			expectedStart: 0,
// 			expectedEnd:   4,
// 			expectError:   false,
// 		},
// 		{
// 			name:          "param_with_tail",
// 			pattern:       "{id}/comments",
// 			expectedTyp:   ntParam,
// 			expectedKey:   "id",
// 			expectedRegex: "",
// 			expectedTail:  '/',
// 			expectedStart: 0,
// 			expectedEnd:   4,
// 			expectError:   false,
// 		},
// 		{
// 			name:          "regexp_param",
// 			pattern:       "{id:[0-9]+}",
// 			expectedTyp:   ntRegexp,
// 			expectedKey:   "id",
// 			expectedRegex: "^[0-9]+$",
// 			expectedTail:  '/', // default tail is '/'
// 			expectedStart: 0,
// 			expectedEnd:   11,
// 			expectError:   false,
// 		},
// 		{
// 			name:          "regexp_with_anchors",
// 			pattern:       "{id:^[a-z]+$}",
// 			expectedTyp:   ntRegexp,
// 			expectedKey:   "id",
// 			expectedRegex: "^[a-z]+$",
// 			expectedTail:  '/', // default tail is '/'
// 			expectedStart: 0,
// 			expectedEnd:   13,
// 			expectError:   false,
// 		},
// 		{
// 			name:          "catchall",
// 			pattern:       "*",
// 			expectedTyp:   ntCatchAll,
// 			expectedKey:   "*",
// 			expectedRegex: "",
// 			expectedTail:  0,
// 			expectedStart: 0,
// 			expectedEnd:   1,
// 			expectError:   false,
// 		},
// 		{
// 			name:        "wildcard_not_last",
// 			pattern:     "*/something",
// 			expectError: true,
// 			errorType:   ErrWildcardMustBeLast,
// 		},
// 		{
// 			name:        "missing_closing_bracket",
// 			pattern:     "{id",
// 			expectError: true,
// 			errorType:   ErrMissingClosingBracket,
// 		},
// 		{
// 			name:          "nested_braces",
// 			pattern:       "{id:{nested}}",
// 			expectedTyp:   ntRegexp,
// 			expectedKey:   "id",
// 			expectedRegex: "^{nested}$",
// 			expectedTail:  '/', // default tail is '/'
// 			expectedStart: 0,
// 			expectedEnd:   13,
// 			expectError:   false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			typ, key, regex, tail, start, end, err := patNextSegment(tc.pattern)

// 			if tc.expectError {
// 				assert.Assert(t, err != nil, "expected error but got nil")
// 				if tc.errorType != nil {
// 					assert.ErrorIs(t, err, tc.errorType)
// 				}
// 			} else {
// 				assert.NilError(t, err)
// 				assert.Equal(t, tc.expectedTyp, typ, "node type mismatch")
// 				assert.Equal(t, tc.expectedKey, key, "param key mismatch")
// 				assert.Equal(t, tc.expectedRegex, regex, "regex pattern mismatch")
// 				assert.Equal(t, tc.expectedTail, tail, "tail byte mismatch")
// 				assert.Equal(t, tc.expectedStart, start, "start index mismatch")
// 				assert.Equal(t, tc.expectedEnd, end, "end index mismatch")
// 			}
// 		})
// 	}
// }

// TestPatParamKeys tests parameter key extraction from patterns
// func TestPatParamKeys(t *testing.T) {
// 	testCases := []struct {
// 		name         string
// 		pattern      string
// 		expectedKeys []string
// 		expectError  bool
// 		errorType    error
// 	}{
// 		{
// 			name:         "no_params",
// 			pattern:      "/posts",
// 			expectedKeys: []string{},
// 			expectError:  false,
// 		},
// 		{
// 			name:         "single_param",
// 			pattern:      "/posts/{id}",
// 			expectedKeys: []string{"id"},
// 			expectError:  false,
// 		},
// 		{
// 			name:         "multiple_params",
// 			pattern:      "/posts/{id}/comments/{commentId}",
// 			expectedKeys: []string{"id", "commentId"},
// 			expectError:  false,
// 		},
// 		{
// 			name:         "regexp_param",
// 			pattern:      "/posts/{id:[0-9]+}",
// 			expectedKeys: []string{"id"},
// 			expectError:  false,
// 		},
// 		{
// 			name:        "duplicate_params",
// 			pattern:     "/posts/{id}/comments/{id}",
// 			expectError: true,
// 			errorType:   ErrDuplicatedParamKey,
// 		},
// 		{
// 			name:         "catchall",
// 			pattern:      "/api/*",
// 			expectedKeys: []string{"*"},
// 			expectError:  false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			keys, err := patParamKeys(tc.pattern)

// 			if tc.expectError {
// 				assert.Assert(t, err != nil, "expected error but got nil")
// 				if tc.errorType != nil {
// 					assert.ErrorIs(t, err, tc.errorType)
// 				}
// 			} else {
// 				assert.NilError(t, err)
// 				assert.DeepEqual(t, tc.expectedKeys, keys)
// 			}
// 		})
// 	}
// }

// // TestLongestPrefix tests the longest common prefix function
// func TestLongestPrefix(t *testing.T) {
// 	testCases := []struct {
// 		name     string
// 		k1       string
// 		k2       string
// 		expected int
// 	}{
// 		{
// 			name:     "identical_strings",
// 			k1:       "/posts",
// 			k2:       "/posts",
// 			expected: 6,
// 		},
// 		{
// 			name:     "common_prefix",
// 			k1:       "/posts/123",
// 			k2:       "/posts/456",
// 			expected: 7,
// 		},
// 		{
// 			name:     "no_common_prefix",
// 			k1:       "/posts",
// 			k2:       "/users",
// 			expected: 1,
// 		},
// 		{
// 			name:     "empty_strings",
// 			k1:       "",
// 			k2:       "",
// 			expected: 0,
// 		},
// 		{
// 			name:     "one_empty",
// 			k1:       "/posts",
// 			k2:       "",
// 			expected: 0,
// 		},
// 		{
// 			name:     "substring",
// 			k1:       "/posts",
// 			k2:       "/posts/123",
// 			expected: 6,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			result := longestPrefix(tc.k1, tc.k2)
// 			assert.Equal(t, tc.expected, result)
// 		})
// 	}
// }

// // TestNodesSorting tests the nodes sorting functionality
// func TestNodesSorting(t *testing.T) {
// 	testCases := []struct {
// 		name     string
// 		nodes    nodes
// 		expected []byte // expected labels after sorting
// 	}{
// 		{
// 			name: "already_sorted",
// 			nodes: nodes{
// 				{label: 'a'},
// 				{label: 'b'},
// 				{label: 'c'},
// 			},
// 			expected: []byte{'a', 'b', 'c'},
// 		},
// 		{
// 			name: "reverse_order",
// 			nodes: nodes{
// 				{label: 'c'},
// 				{label: 'b'},
// 				{label: 'a'},
// 			},
// 			expected: []byte{'a', 'b', 'c'},
// 		},
// 		{
// 			name: "random_order",
// 			nodes: nodes{
// 				{label: 'm'},
// 				{label: 'a'},
// 				{label: 'z'},
// 				{label: 'b'},
// 			},
// 			expected: []byte{'a', 'b', 'm', 'z'},
// 		},
// 		{
// 			name: "tail_sort_param_nodes",
// 			nodes: nodes{
// 				{label: 'a', typ: ntParam, tail: '/'},
// 				{label: 'b', typ: ntParam, tail: '-'},
// 			},
// 			expected: []byte{'b', 'a'}, // '/' tail should be last
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.nodes.Sort()
// 			for i, node := range tc.nodes {
// 				assert.Equal(t, tc.expected[i], node.label)
// 			}
// 		})
// 	}
// }

// // TestNodesFindEdge tests binary search for finding edges
// func TestNodesFindEdge(t *testing.T) {
// 	nodes := nodes{
// 		{label: 'a'},
// 		{label: 'c'},
// 		{label: 'e'},
// 		{label: 'g'},
// 		{label: 'z'},
// 	}

// 	testCases := []struct {
// 		name      string
// 		label     byte
// 		expectNil bool
// 	}{
// 		{name: "find_first", label: 'a', expectNil: false},
// 		{name: "find_middle", label: 'e', expectNil: false},
// 		{name: "find_last", label: 'z', expectNil: false},
// 		{name: "not_found_before", label: 'A', expectNil: true},
// 		{name: "not_found_between", label: 'd', expectNil: true},
// 		{name: "not_found_after", label: '~', expectNil: true},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			result := nodes.findEdge(tc.label)
// 			if tc.expectNil {
// 				assert.Assert(t, result == nil)
// 			} else {
// 				assert.Assert(t, result != nil)
// 				assert.Equal(t, tc.label, result.label)
// 			}
// 		})
// 	}
// }

// // TestComplexRoutingScenarios tests complex real-world routing scenarios
// func TestComplexRoutingScenarios(t *testing.T) {
// 	t.Run("RESTful_API_with_versioning", func(t *testing.T) {
// 		node := new(Node)

// 		routes := []string{
// 			"/api/v1/users",
// 			"/api/v1/users/{id}",
// 			"/api/v1/users/{id}/posts",
// 			"/api/v1/users/{id}/posts/{postId}",
// 			"/api/v2/users",
// 			"/api/v2/users/{id}",
// 		}

// 		for _, route := range routes {
// 			_, err := node.InsertRoute(route, &highv3.PathItem{
// 				Get: &highv3.Operation{},
// 			}, &proxyhandler.InsertRouteOptions{})
// 			assert.NilError(t, err)
// 		}

// 		// Test v1 routes
// 		r := node.FindRoute("/api/v1/users", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "/api/v1/users", r.Node.pattern)

// 		r = node.FindRoute("/api/v1/users/123", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "123", r.ParamValues["id"])

// 		r = node.FindRoute("/api/v1/users/123/posts/456", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "123", r.ParamValues["id"])
// 		assert.Equal(t, "456", r.ParamValues["postId"])

// 		// Test v2 routes
// 		r = node.FindRoute("/api/v2/users/789", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "/api/v2/users/{id}", r.Node.pattern)
// 	})

// 	t.Run("mixed_static_and_dynamic_routes", func(t *testing.T) {
// 		node := new(Node)

// 		routes := map[string]*highv3.PathItem{
// 			"/posts":           {Get: &highv3.Operation{}},
// 			"/posts/new":       {Get: &highv3.Operation{}},
// 			"/posts/popular":   {Get: &highv3.Operation{}},
// 			"/posts/{id}":      {Get: &highv3.Operation{}},
// 			"/posts/{id}/edit": {Get: &highv3.Operation{}},
// 		}

// 		for pattern, handlers := range routes {
// 			_, err := node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
// 			assert.NilError(t, err)
// 		}

// 		// Static routes should take precedence
// 		r := node.FindRoute("/posts/new", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "/posts/new", r.Node.pattern)

// 		r = node.FindRoute("/posts/popular", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "/posts/popular", r.Node.pattern)

// 		// Dynamic route should match other IDs
// 		r = node.FindRoute("/posts/123", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "/posts/{id}", r.Node.pattern)
// 		assert.Equal(t, "123", r.ParamValues["id"])

// 		r = node.FindRoute("/posts/456/edit", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "/posts/{id}/edit", r.Node.pattern)
// 		assert.Equal(t, "456", r.ParamValues["id"])
// 	})

// 	t.Run("regexp_validation_routes", func(t *testing.T) {
// 		node := new(Node)

// 		routes := map[string]*highv3.PathItem{
// 			"/users/{id:[0-9]+}":                       {Get: &highv3.Operation{}},
// 			"/users/{username:[a-z]+}":                 {Post: &highv3.Operation{}},
// 			"/posts/{slug:[a-z0-9-]+}":                 {Get: &highv3.Operation{}},
// 			"/files/{filename:[a-zA-Z0-9._-]+}":        {Get: &highv3.Operation{}},
// 			"/dates/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}": {Get: &highv3.Operation{}},
// 		}

// 		for pattern, handlers := range routes {
// 			_, err := node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
// 			assert.NilError(t, err)
// 		}

// 		// Numeric ID should match
// 		r := node.FindRoute("/users/12345", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "12345", r.ParamValues["id"])

// 		// Alphabetic username should match
// 		r = node.FindRoute("/users/john", http.MethodPost)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "john", r.ParamValues["username"])

// 		// Slug with hyphens should match
// 		r = node.FindRoute("/posts/my-blog-post-123", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "my-blog-post-123", r.ParamValues["slug"])

// 		// Filename with dots and underscores should match
// 		r = node.FindRoute("/files/my_file.txt", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "my_file.txt", r.ParamValues["filename"])

// 		// Date format should match
// 		r = node.FindRoute("/dates/2024-01-15", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "2024-01-15", r.ParamValues["date"])

// 		// Invalid date format should not match
// 		r = node.FindRoute("/dates/2024-1-5", http.MethodGet)
// 		assert.Assert(t, r == nil)
// 	})

// 	t.Run("catchall_with_specific_routes", func(t *testing.T) {
// 		node := new(Node)

// 		routes := map[string]*highv3.PathItem{
// 			"/api/v1/users": {Get: &highv3.Operation{}},
// 			"/api/v1/posts": {Get: &highv3.Operation{}},
// 			"/api/v1/*":     {Get: &highv3.Operation{}},
// 			"/static/css/*": {Get: &highv3.Operation{}},
// 			"/static/js/*":  {Get: &highv3.Operation{}},
// 		}

// 		for pattern, handlers := range routes {
// 			_, err := node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
// 			assert.NilError(t, err)
// 		}

// 		// Specific routes should match first
// 		r := node.FindRoute("/api/v1/users", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "/api/v1/users", r.Node.pattern)

// 		// Catchall should match unspecified routes
// 		r = node.FindRoute("/api/v1/anything/else", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "/api/v1/*", r.Node.pattern)

// 		r = node.FindRoute("/static/css/main.css", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "/static/css/*", r.Node.pattern)
// 	})

// 	t.Run("multiple_HTTP_methods", func(t *testing.T) {
// 		node := new(Node)

// 		_, err := node.InsertRoute("/posts/{id}", &highv3.PathItem{
// 			Get:    &highv3.Operation{},
// 			Post:   &highv3.Operation{},
// 			Put:    &highv3.Operation{},
// 			Patch:  &highv3.Operation{},
// 			Delete: &highv3.Operation{},
// 		}, &proxyhandler.InsertRouteOptions{})
// 		assert.NilError(t, err)

// 		methods := []string{
// 			http.MethodGet,
// 			http.MethodPost,
// 			http.MethodPut,
// 			http.MethodPatch,
// 			http.MethodDelete,
// 		}

// 		for _, method := range methods {
// 			r := node.FindRoute("/posts/123", method)
// 			assert.Assert(t, r != nil, "method %s should be found", method)
// 			assert.Equal(t, "123", r.ParamValues["id"])
// 		}

// 		// Method not defined should not be found
// 		r := node.FindRoute("/posts/123", http.MethodHead)
// 		assert.Assert(t, r == nil)
// 	})

// 	t.Run("deeply_nested_resources", func(t *testing.T) {
// 		node := new(Node)

// 		pattern := "/orgs/{orgId}/teams/{teamId}/projects/{projectId}/tasks/{taskId}/comments/{commentId}"
// 		_, err := node.InsertRoute(pattern, &highv3.PathItem{
// 			Get: &highv3.Operation{},
// 		}, &proxyhandler.InsertRouteOptions{})
// 		assert.NilError(t, err)

// 		r := node.FindRoute("/orgs/org1/teams/team2/projects/proj3/tasks/task4/comments/comment5", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "org1", r.ParamValues["orgId"])
// 		assert.Equal(t, "team2", r.ParamValues["teamId"])
// 		assert.Equal(t, "proj3", r.ParamValues["projectId"])
// 		assert.Equal(t, "task4", r.ParamValues["taskId"])
// 		assert.Equal(t, "comment5", r.ParamValues["commentId"])
// 	})

// 	t.Run("special_characters_in_params", func(t *testing.T) {
// 		node := new(Node)

// 		routes := map[string]*highv3.PathItem{
// 			"/search/{query}":      {Get: &highv3.Operation{}},
// 			"/users/{email:[^/]+}": {Get: &highv3.Operation{}},
// 		}

// 		for pattern, handlers := range routes {
// 			_, err := node.InsertRoute(pattern, handlers, &proxyhandler.InsertRouteOptions{})
// 			assert.NilError(t, err)
// 		}

// 		// Query with special characters
// 		r := node.FindRoute("/search/hello-world", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "hello-world", r.ParamValues["query"])

// 		// Email-like parameter
// 		r = node.FindRoute("/users/user@example.com", http.MethodGet)
// 		assert.Assert(t, r != nil)
// 		assert.Equal(t, "user@example.com", r.ParamValues["email"])
// 	})
// }

// // TestNodeIsLeaf tests the isLeaf method
// func TestNodeIsLeaf(t *testing.T) {
// 	testCases := []struct {
// 		name     string
// 		node     *Node
// 		expected bool
// 	}{
// 		{
// 			name: "leaf_node_with_handlers",
// 			node: &Node{
// 				handlers: map[string]MethodHandler{
// 					http.MethodGet: {},
// 				},
// 			},
// 			expected: true,
// 		},
// 		{
// 			name: "non_leaf_node",
// 			node: &Node{
// 				handlers: nil,
// 			},
// 			expected: false,
// 		},
// 		{
// 			name: "empty_handlers_map",
// 			node: &Node{
// 				handlers: map[string]MethodHandler{},
// 			},
// 			expected: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			result := tc.node.isLeaf()
// 			assert.Equal(t, tc.expected, result)
// 		})
// 	}
// }

// // TestNodeGetEdge tests the getEdge method
// func TestNodeGetEdge(t *testing.T) {
// 	node := &Node{
// 		children: [ntCatchAll + 1]nodes{
// 			ntStatic: {
// 				{label: 'a', tail: 0},
// 				{label: 'b', tail: 0},
// 			},
// 			ntParam: {
// 				{label: '{', tail: '/', paramName: "id"},
// 			},
// 			ntRegexp: {
// 				{label: '{', tail: '/', prefix: "^[0-9]+$"},
// 				{label: '{', tail: '/', prefix: "^[a-z]+$"},
// 			},
// 		},
// 	}

// 	testCases := []struct {
// 		name      string
// 		ntyp      nodeTyp
// 		label     byte
// 		tail      byte
// 		prefix    string
// 		expectNil bool
// 		checkFunc func(t *testing.T, n *Node)
// 	}{
// 		{
// 			name:      "find_static_node",
// 			ntyp:      ntStatic,
// 			label:     'a',
// 			tail:      0,
// 			expectNil: false,
// 			checkFunc: func(t *testing.T, n *Node) {
// 				assert.Equal(t, byte('a'), n.label)
// 			},
// 		},
// 		{
// 			name:      "find_param_node",
// 			ntyp:      ntParam,
// 			label:     '{',
// 			tail:      '/',
// 			expectNil: false,
// 			checkFunc: func(t *testing.T, n *Node) {
// 				assert.Equal(t, "id", n.paramName)
// 			},
// 		},
// 		{
// 			name:      "find_regexp_node_first",
// 			ntyp:      ntRegexp,
// 			label:     '{',
// 			tail:      '/',
// 			prefix:    "^[0-9]+$",
// 			expectNil: false,
// 			checkFunc: func(t *testing.T, n *Node) {
// 				assert.Equal(t, "^[0-9]+$", n.prefix)
// 			},
// 		},
// 		{
// 			name:      "find_regexp_node_second",
// 			ntyp:      ntRegexp,
// 			label:     '{',
// 			tail:      '/',
// 			prefix:    "^[a-z]+$",
// 			expectNil: false,
// 			checkFunc: func(t *testing.T, n *Node) {
// 				assert.Equal(t, "^[a-z]+$", n.prefix)
// 			},
// 		},
// 		{
// 			name:      "not_found_wrong_label",
// 			ntyp:      ntStatic,
// 			label:     'z',
// 			tail:      0,
// 			expectNil: true,
// 		},
// 		{
// 			name:      "not_found_wrong_prefix",
// 			ntyp:      ntRegexp,
// 			label:     '{',
// 			tail:      '/',
// 			prefix:    "^[0-9a-z]+$",
// 			expectNil: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			result := node.getEdge(tc.ntyp, tc.label, tc.tail, tc.prefix)

// 			if tc.expectNil {
// 				assert.Assert(t, result == nil)
// 			} else {
// 				assert.Assert(t, result != nil)
// 				if tc.checkFunc != nil {
// 					tc.checkFunc(t, result)
// 				}
// 			}
// 		})
// 	}
// }

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
