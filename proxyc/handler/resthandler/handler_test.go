package resthandler

import (
	"testing"

	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"github.com/relychan/relixy/schema/base_schema"
	"gotest.tools/v3/assert"
)

func TestNewRESTHandler(t *testing.T) {
	testCases := []struct {
		name        string
		operation   *highv3.Operation
		proxyAction *base_schema.RelixyAction
		options     *proxyhandler.NewRelixyHandlerOptions
		expectError bool
	}{
		{
			name:      "basic REST handler",
			operation: &highv3.Operation{},
			proxyAction: &base_schema.RelixyAction{
				Type: base_schema.ProxyTypeREST,
			},
			options: &proxyhandler.NewRelixyHandlerOptions{
				Method: "GET",
			},
			expectError: false,
		},
		{
			name:      "REST handler with custom path",
			operation: &highv3.Operation{},
			proxyAction: &base_schema.RelixyAction{
				Type: base_schema.ProxyTypeREST,
				Path: "/custom/path",
			},
			options: &proxyhandler.NewRelixyHandlerOptions{
				Method: "POST",
			},
			expectError: false,
		},
		{
			name:        "REST handler with nil proxy action",
			operation:   &highv3.Operation{},
			proxyAction: nil,
			options: &proxyhandler.NewRelixyHandlerOptions{
				Method: "GET",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := NewRESTHandler(tc.operation, tc.proxyAction, tc.options)

			if tc.expectError {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				assert.Assert(t, handler != nil)
				assert.Equal(t, base_schema.ProxyTypeREST, handler.Type())
			}
		})
	}
}

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
