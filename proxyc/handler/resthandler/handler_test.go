package resthandler

import (
	"testing"

	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"go.yaml.in/yaml/v4"
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

// TestNewRESTHandler tests the NewRESTHandler function
func TestNewRESTHandler(t *testing.T) {
	t.Run("nil_proxy_action", func(t *testing.T) {
		operation := &highv3.Operation{}
		options := &proxyhandler.NewRelixyHandlerOptions{
			Method: "GET",
		}

		handler, err := NewRESTHandler(operation, nil, options)
		assert.NilError(t, err)
		assert.Assert(t, handler != nil)
		assert.Equal(t, ProxyActionTypeREST, handler.Type())
	})

	t.Run("with_custom_request_path", func(t *testing.T) {
		operation := &highv3.Operation{}
		options := &proxyhandler.NewRelixyHandlerOptions{
			Method: "POST",
		}

		yamlConfig := `
type: rest
request:
  path: "/custom/path"
`
		var rawAction yaml.Node
		err := yaml.Unmarshal([]byte(yamlConfig), &rawAction)
		assert.NilError(t, err)

		handler, err := NewRESTHandler(operation, &rawAction, options)
		assert.NilError(t, err)
		assert.Assert(t, handler != nil)
	})

	t.Run("invalid_yaml", func(t *testing.T) {
		operation := &highv3.Operation{}
		options := &proxyhandler.NewRelixyHandlerOptions{
			Method: "GET",
		}

		yamlConfig := `
type: rest
request:
  headers:
    invalid: [1, 2, 3]
`
		var rawAction yaml.Node
		err := yaml.Unmarshal([]byte(yamlConfig), &rawAction)
		assert.NilError(t, err)

		_, err = NewRESTHandler(operation, &rawAction, options)
		assert.Assert(t, err != nil)
	})

	t.Run("invalid_response_config", func(t *testing.T) {
		operation := &highv3.Operation{}
		options := &proxyhandler.NewRelixyHandlerOptions{
			Method: "GET",
		}

		yamlConfig := `
type: rest
response:
  body:
    invalid: true
`
		var rawAction yaml.Node
		err := yaml.Unmarshal([]byte(yamlConfig), &rawAction)
		assert.NilError(t, err)

		_, err = NewRESTHandler(operation, &rawAction, options)
		assert.Assert(t, err != nil)
	})
}
