package resthandler

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestRESTHandler_Type(t *testing.T) {
	handler := &RESTHandler{}
	assert.Equal(t, ProxyActionTypeREST, handler.Type())
}

func TestRESTHandler_Properties(t *testing.T) {
	testCases := []struct {
		name           string
		handler        *RESTHandler
		expectedMethod string
		expectedPath   string
		hasCustomPath  bool
	}{
		{
			name: "handler with GET method",
			handler: &RESTHandler{
				method: "GET",
			},
			expectedMethod: "GET",
			expectedPath:   "",
			hasCustomPath:  false,
		},
		{
			name: "handler with POST method and custom path",
			handler: &RESTHandler{
				method: "POST",
				customRequest: &customRESTRequest{
					Path: "/custom/path",
				},
			},
			expectedMethod: "POST",
			expectedPath:   "/custom/path",
			hasCustomPath:  true,
		},
		{
			name: "handler with PUT method",
			handler: &RESTHandler{
				method: "PUT",
				customRequest: &customRESTRequest{
					Path: "/api/resource",
				},
			},
			expectedMethod: "PUT",
			expectedPath:   "/api/resource",
			hasCustomPath:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedMethod, tc.handler.method)
			assert.Equal(t, ProxyActionTypeREST, tc.handler.Type())
			if tc.expectedPath != "" {
				assert.Equal(t, tc.expectedPath, tc.handler.customRequest.Path)
			}
		})
	}
}
