package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"go.yaml.in/yaml/v4"
	"gotest.tools/v3/assert"
)

func TestRelixyOpenAPIResourceDefinition_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyOpenAPIResourceDefinition)
	}{
		{
			name: "valid minimal spec",
			jsonData: `{
				"spec": {
					"openapi": "3.0.0",
					"info": {
						"title": "Test API",
						"version": "1.0.0"
					},
					"paths": {}
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, def *RelixyOpenAPIResourceDefinition) {
				assert.Assert(t, def.Spec != nil)
				assert.Equal(t, "Test API", def.Spec.Info.Title)
				assert.Equal(t, "1.0.0", def.Spec.Info.Version)
			},
		},
		{
			name: "valid spec with settings",
			jsonData: `{
				"settings": {
					"basePath": "/api/v1"
				},
				"spec": {
					"openapi": "3.0.0",
					"info": {
						"title": "Test API",
						"version": "1.0.0"
					},
					"paths": {}
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, def *RelixyOpenAPIResourceDefinition) {
				assert.Assert(t, def.Spec != nil)
				assert.Equal(t, "/api/v1", def.Settings.BasePath)
				assert.Equal(t, "Test API", def.Spec.Info.Title)
			},
		},
		{
			name: "valid spec with servers",
			jsonData: `{
				"spec": {
					"openapi": "3.0.0",
					"info": {
						"title": "Test API",
						"version": "1.0.0"
					},
					"servers": [
						{
							"url": "https://api.example.com",
							"x-rely-url-env": "SERVER_URL"
						}
					],
					"paths": {}
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, def *RelixyOpenAPIResourceDefinition) {
				assert.Assert(t, def.Spec != nil)
				assert.Assert(t, def.Spec.Servers != nil)
				assert.Equal(t, "https://api.example.com", def.Spec.Servers[0].URL)
				urlFromEnv, exist := def.Spec.Servers[0].Extensions.Get(XRelyURLEnv)
				assert.Assert(t, exist)
				assert.Equal(t, "SERVER_URL", urlFromEnv.Value)
			},
		},
		{
			name: "missing spec",
			jsonData: `{
				"settings": {
					"basePath": "/api/v1"
				}
			}`,
			expectError: false,
			checkFunc:   nil,
		},
		{
			name:        "null spec",
			jsonData:    `{"spec": null}`,
			expectError: true,
			checkFunc:   nil,
		},
		{
			name:        "empty object",
			jsonData:    `{}`,
			expectError: false,
			checkFunc:   nil,
		},
		{
			name: "invalid spec format",
			jsonData: `{
				"spec": {
					"invalid": "data"
				}
			}`,
			expectError: true,
			checkFunc:   nil,
		},
		{
			name:        "invalid json",
			jsonData:    `{"spec": invalid}`,
			expectError: true,
			checkFunc:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var def RelixyOpenAPIResourceDefinition
			err := json.Unmarshal([]byte(tc.jsonData), &def)
			if tc.expectError {
				assert.Assert(t, err != nil, "expected error but got nil")
			} else {
				assert.NilError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &def)
				}
			}
		})
	}
}

func TestRelixyOpenAPIResourceDefinition_UnmarshalYAML(t *testing.T) {
	testCases := []struct {
		name        string
		yamlData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyOpenAPIResourceDefinition)
	}{
		{
			name: "valid minimal spec with servers and paths",
			yamlData: `spec:
  openapi: "3.0.0"
  info:
    title: Test API
    version: "1.0.0"
  servers:
    - url: https://api.example.com
      x-rely-url-env: SERVER_URL
  paths:
    /users:
      get:
        operationId: getUsers`,
			expectError: false,
			checkFunc: func(t *testing.T, def *RelixyOpenAPIResourceDefinition) {
				assert.Assert(t, def.Spec != nil)
				assert.Assert(t, def.Spec.Servers != nil)
				assert.Assert(t, def.Spec.Paths != nil)
				assert.Equal(t, "https://api.example.com", def.Spec.Servers[0].URL)
				urlFromEnv, exist := def.Spec.Servers[0].Extensions.Get(XRelyURLEnv)
				assert.Assert(t, exist)
				assert.Equal(t, "SERVER_URL", urlFromEnv.Value)
			},
		},
		{
			name: "valid spec with settings",
			yamlData: `settings:
  basePath: /api/v1
spec:
  openapi: "3.0.0"
  info:
    title: Test API
    version: "1.0.0"
  servers:
    - url: https://api.example.com
  paths: {}`,
			expectError: false,
			checkFunc: func(t *testing.T, def *RelixyOpenAPIResourceDefinition) {
				assert.Assert(t, def.Spec != nil)
				assert.Equal(t, "/api/v1", def.Settings.BasePath)
			},
		},
		{
			name: "missing spec",
			yamlData: `settings:
  basePath: /api/v1`,
			expectError: false,
			checkFunc:   nil,
		},
		{
			name:        "null spec",
			yamlData:    `spec: null`,
			expectError: true,
			checkFunc:   nil,
		},
		{
			name:        "empty object",
			yamlData:    `{}`,
			expectError: false,
			checkFunc:   nil,
		},
		{
			name: "invalid spec format",
			yamlData: `spec:
  invalid: data`,
			expectError: true,
			checkFunc:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var def RelixyOpenAPIResourceDefinition
			err := yaml.Unmarshal([]byte(tc.yamlData), &def)

			if tc.expectError {
				assert.Assert(t, err != nil, "expected error but got nil")
			} else {
				assert.NilError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &def)
				}
			}
		})
	}
}

func TestRelixyOpenAPIResource_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyOpenAPIResource)
	}{
		{
			name: "valid complete resource",
			jsonData: `{
				"version": "v1",
				"kind": "OpenAPI",
				"metadata": {
					"name": "test-api",
					"description": "Test API description"
				},
				"definition": {
					"spec": {
						"openapi": "3.0.0",
						"info": {
							"title": "Test API",
							"version": "1.0.0"
						},
						"paths": {}
					}
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, res *RelixyOpenAPIResource) {
				assert.Equal(t, "v1", res.Version)
				assert.Equal(t, "OpenAPI", res.Kind)
				assert.Equal(t, "test-api", res.Metadata.Name)
				assert.Equal(t, "Test API description", res.Metadata.Description)
				assert.Assert(t, res.Definition.Spec != nil)
				assert.Equal(t, "Test API", res.Definition.Spec.Info.Title)
			},
		},
		{
			name: "valid resource with settings",
			jsonData: `{
				"version": "v1",
				"kind": "OpenAPI",
				"metadata": {
					"name": "test-api"
				},
				"definition": {
					"settings": {
						"basePath": "/api/v1"
					},
					"spec": {
						"openapi": "3.0.0",
						"info": {
							"title": "Test API",
							"version": "1.0.0"
						},
						"paths": {}
					}
				}
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, res *RelixyOpenAPIResource) {
				assert.Equal(t, "/api/v1", res.Definition.Settings.BasePath)
				assert.Assert(t, res.Definition.Spec != nil)
			},
		},
		{
			name: "missing definition spec",
			jsonData: `{
				"version": "v1",
				"kind": "OpenAPI",
				"metadata": {
					"name": "test-api"
				},
				"definition": {
					"settings": {
						"basePath": "/api/v1"
					}
				}
			}`,
			expectError: false,
			checkFunc:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res RelixyOpenAPIResource
			err := json.Unmarshal([]byte(tc.jsonData), &res)
			if tc.expectError {
				assert.Assert(t, err != nil, "expected error but got nil")
			} else {
				assert.NilError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &res)
				}
			}
		})
	}
}

func TestRelixyOpenAPIResource_UnmarshalYAML(t *testing.T) {
	testCases := []struct {
		name        string
		yamlData    string
		expectError bool
		checkFunc   func(*testing.T, *RelixyOpenAPIResource)
	}{
		{
			name: "valid complete resource",
			yamlData: `version: v1
kind: OpenAPI
metadata:
  name: test-api
  description: Test API description
definition:
  spec:
    openapi: "3.0.0"
    info:
      title: Test API
      version: "1.0.0"
    servers:
      - url: https://api.example.com
    paths:
      /users:
        get:
          operationId: getUsers`,
			expectError: false,
			checkFunc: func(t *testing.T, res *RelixyOpenAPIResource) {
				assert.Equal(t, "v1", res.Version)
				assert.Equal(t, "OpenAPI", res.Kind)
				assert.Equal(t, "test-api", res.Metadata.Name)
				assert.Equal(t, "Test API description", res.Metadata.Description)
				assert.Assert(t, res.Definition.Spec != nil)
			},
		},
		{
			name: "valid resource with settings",
			yamlData: `version: v1
kind: OpenAPI
metadata:
  name: test-api
definition:
  settings:
    basePath: /api/v1
  spec:
    openapi: "3.0.0"
    info:
      title: Test API
      version: "1.0.0"
    servers:
      - url: https://api.example.com
    paths: {}`,
			expectError: false,
			checkFunc: func(t *testing.T, res *RelixyOpenAPIResource) {
				assert.Equal(t, "/api/v1", res.Definition.Settings.BasePath)
				assert.Assert(t, res.Definition.Spec != nil)
			},
		},
		{
			name: "missing definition spec",
			yamlData: `version: v1
kind: OpenAPI
metadata:
  name: test-api
definition:
  settings:
    basePath: /api/v1`,
			expectError: false,
			checkFunc:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res RelixyOpenAPIResource
			err := yaml.Unmarshal([]byte(tc.yamlData), &res)

			if tc.expectError {
				assert.Assert(t, err != nil, "expected error but got nil")
			} else {
				assert.NilError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, &res)
				}
			}
		})
	}
}

func TestRelixyOpenAPIResource_BuildJSON(t *testing.T) {
	rawFileSpec := `{
  "openapi": "3.0.0",
  "info": {
    "title": "Test API",
    "version": "1.0.0"
  },
  "paths": {
    "/users": {
      "get": {
        "x-rely-proxy-action": {
          "type": "graphql",
          "request": {
            "query": "{users { id }}"
          }
        }
      }
    },
    "/projects": {
      "get": {
        "operationId": "getProjects"
      }
    }
  }
}`

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "openapi.json")

	err := os.WriteFile(tempFile, []byte(rawFileSpec), 0664)
	assert.NilError(t, err)

	jsonData := fmt.Sprintf(`{
  "version": "v1",
  "kind": "OpenAPI",
  "metadata": {
    "name": "test-api",
    "description": "Test API description"
  },
  "definition": {
    "ref": "%s",
    "spec": {
      "openapi": "3.0.0",
      "info": {
        "title": "Test API",
        "version": "1.0.0"
      },
      "servers": [
        {
          "url": "https://api.example.com"
        }
      ],
      "paths": {
        "/users": {
          "get": {
            "operationId": "getUsers"
          }
        }
      }
    }
  }
}`, tempFile)

	expectedData := `{
	"spec": {
		"openapi": "3.0.0",
		"info": {
			"title": "Test API",
			"version": "1.0.0"
		},
		"servers": [
			{
			"url": "https://api.example.com"
			}
		],
		"paths": {
			"/users": {
			"get": {
				"operationId": "getUsers",
				"x-rely-proxy-action": {
				"type": "graphql",
				"request": {
					"query": "{users { id }}"
				}
				}
			}
			},
			"/projects": {
			"get": {
				"operationId": "getProjects"
			}
			}
		}
		}
}`

	var rawResource RelixyOpenAPIResource
	err = json.Unmarshal([]byte(jsonData), &rawResource)
	assert.NilError(t, err)

	result, err := rawResource.Definition.Build(context.TODO())
	assert.NilError(t, err)

	resultJsonBytes, err := json.Marshal(RelixyOpenAPIResourceDefinition{
		Spec: result,
	})
	assert.NilError(t, err)

	log.Println("result", string(resultJsonBytes))

	var resultJson any
	assert.NilError(t, json.Unmarshal(resultJsonBytes, &resultJson))

	var expectedJson any
	assert.NilError(t, json.Unmarshal([]byte(expectedData), &expectedJson))

	assert.DeepEqual(t, resultJson, expectedJson)
}

func TestRelixyOpenAPIResource_BuildYAML(t *testing.T) {
	rawFileSpec := `openapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
paths:
  /users:
    get:
      x-rely-proxy-action:
        type: graphql
        request:
          query: '{users { id }}'
  /projects:
    get:
      operationId: getProjects`

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "openapi.yaml")

	err := os.WriteFile(tempFile, []byte(rawFileSpec), 0664)
	assert.NilError(t, err)

	yamlData := fmt.Sprintf(`version: v1
kind: OpenAPI
metadata:
  name: test-api
  description: Test API description
definition:
  ref: '%s'
  spec:
    openapi: "3.0.0"
    info:
      title: Test API
      version: "1.0.0"
    servers:
      - url: https://api.example.com
    paths:
      /users:
        get:
          operationId: getUsers`, tempFile)

	expectedData := `openapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
servers:
  - url: https://api.example.com
paths:
  /users:
    get:
      operationId: getUsers
      x-rely-proxy-action:
        type: graphql
        request:
          query: '{users { id }}'
  /projects:
    get:
      operationId: getProjects`

	var rawResource RelixyOpenAPIResource
	err = yaml.Unmarshal([]byte(yamlData), &rawResource)
	assert.NilError(t, err)

	result, err := rawResource.Definition.Build(context.TODO())
	assert.NilError(t, err)

	resultYamlBytes, err := yaml.Dump(result)
	assert.NilError(t, err)

	var resultYaml any
	assert.NilError(t, yaml.Load(resultYamlBytes, &resultYaml))

	var expectedYaml any
	assert.NilError(t, yaml.Load([]byte(expectedData), &expectedYaml))

	assert.DeepEqual(t, resultYaml, expectedYaml)
}
