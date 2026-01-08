package resthandler

import (
	"testing"

	"github.com/relychan/relixy/schema/base_schema"
	"gotest.tools/v3/assert"
)

func TestRESTHandler_Type(t *testing.T) {
	handler := &RESTHandler{}
	assert.Equal(t, base_schema.ProxyTypeREST, handler.Type())
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
				method:      "GET",
				requestPath: "",
			},
			expectedMethod: "GET",
			expectedPath:   "",
			hasCustomPath:  false,
		},
		{
			name: "handler with POST method and custom path",
			handler: &RESTHandler{
				method:      "POST",
				requestPath: "/custom/path",
			},
			expectedMethod: "POST",
			expectedPath:   "/custom/path",
			hasCustomPath:  true,
		},
		{
			name: "handler with PUT method",
			handler: &RESTHandler{
				method:      "PUT",
				requestPath: "/api/resource",
			},
			expectedMethod: "PUT",
			expectedPath:   "/api/resource",
			hasCustomPath:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedMethod, tc.handler.method)
			assert.Equal(t, tc.expectedPath, tc.handler.requestPath)
			assert.Equal(t, base_schema.ProxyTypeREST, tc.handler.Type())
		})
	}
}
