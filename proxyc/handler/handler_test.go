package handler

import (
	"testing"

	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/relychan/relixy/proxyc/handler/graphqlhandler"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"github.com/relychan/relixy/proxyc/handler/resthandler"
	"github.com/relychan/relixy/schema/base_schema"
	"github.com/relychan/relixy/schema/openapi"
	"go.yaml.in/yaml/v4"
	"gotest.tools/v3/assert"
)

func TestNewProxyHandler(t *testing.T) {
	testCases := []struct {
		name          string
		operation     *highv3.Operation
		options       *proxyhandler.NewRelixyHandlerOptions
		expectedType  base_schema.RelixyActionType
		expectError   bool
		errorContains string
	}{
		{
			name: "REST handler without proxy action",
			operation: &highv3.Operation{
				OperationId: "testOperation",
			},
			options: &proxyhandler.NewRelixyHandlerOptions{
				Method: "GET",
			},
			expectedType: base_schema.ProxyTypeREST,
			expectError:  false,
		},
		{
			name: "REST handler with explicit proxy action",
			operation: createOperationWithProxyAction(t, resthandler.RelixyRESTActionConfig{
				Type: base_schema.ProxyTypeREST,
				Request: &resthandler.RelixyRESTRequestConfig{
					Path: "/test",
				},
			}),
			options: &proxyhandler.NewRelixyHandlerOptions{
				Method: "POST",
			},
			expectedType: base_schema.ProxyTypeREST,
			expectError:  false,
		},
		{
			name: "GraphQL handler with valid query",
			operation: createOperationWithProxyAction(t, graphqlhandler.RelixyGraphQLActionConfig{
				Type: base_schema.ProxyTypeGraphQL,
				Request: &graphqlhandler.RelixyGraphQLRequestConfig{
					Query: "query { users { id name } }",
				},
			}),
			options: &proxyhandler.NewRelixyHandlerOptions{
				Method: "POST",
			},
			expectedType: base_schema.ProxyTypeGraphQL,
			expectError:  false,
		},
		{
			name: "GraphQL handler with invalid query",
			operation: createOperationWithProxyAction(t, graphqlhandler.RelixyGraphQLActionConfig{
				Type: base_schema.ProxyTypeGraphQL,
				Request: &graphqlhandler.RelixyGraphQLRequestConfig{
					Query: "invalid query {",
				},
			}),
			options: &proxyhandler.NewRelixyHandlerOptions{
				Method: "POST",
			},
			expectError:   true,
			errorContains: "Unexpected Name",
		},
		{
			name: "unsupported proxy type",
			operation: createOperationWithProxyAction(t, map[string]any{
				"type": "unsupported",
			}),
			options: &proxyhandler.NewRelixyHandlerOptions{
				Method: "GET",
			},
			expectError:   true,
			errorContains: "unsupported proxy type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := NewProxyHandler(tc.operation, tc.options)

			if tc.expectError {
				assert.Assert(t, err != nil, "expected error but got nil")
				if tc.errorContains != "" {
					assert.ErrorContains(t, err, tc.errorContains)
				}
			} else {
				assert.NilError(t, err)
				assert.Assert(t, handler != nil)
				assert.Equal(t, tc.expectedType, handler.Type())
			}
		})
	}
}

func TestRegisterProxyHandler(t *testing.T) {
	// Save original constructors
	originalConstructors := make(map[base_schema.RelixyActionType]proxyhandler.NewRelixyHandlerFunc)
	for k, v := range proxyHandlerConstructors {
		originalConstructors[k] = v
	}

	// Restore original constructors after test
	defer func() {
		proxyHandlerConstructors = originalConstructors
	}()

	customType := base_schema.RelixyActionType("custom")
	customConstructor := func(
		operation *highv3.Operation,
		proxyAction *yaml.Node,
		options *proxyhandler.NewRelixyHandlerOptions,
	) (proxyhandler.RelixyHandler, error) {
		return nil, nil
	}

	RegisterProxyHandler(customType, customConstructor)

	_, exists := proxyHandlerConstructors[customType]
	assert.Assert(t, exists, "custom handler should be registered")
}

// Helper function to create an operation with a proxy action extension
func createOperationWithProxyAction(t *testing.T, action any) *highv3.Operation {
	t.Helper()

	extensions := orderedmap.New[string, *yaml.Node]()

	actionData, err := yaml.Marshal(action)
	assert.NilError(t, err)

	var actionNode yaml.Node
	err = yaml.Unmarshal(actionData, &actionNode)
	assert.NilError(t, err)

	extensions.Set(openapi.XRelyProxyAction, &actionNode)

	return &highv3.Operation{
		OperationId: "testOperation",
		Extensions:  extensions,
	}
}
