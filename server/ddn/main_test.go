package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/hasura/gotel"
	"github.com/relychan/relyx/config"
	"github.com/relychan/relyx/routes/ddn"
	"go.opentelemetry.io/otel"
	"gotest.tools/v3/assert"
)

func TestRESTServer(t *testing.T) {
	server, shutdown := initTestServer(t, "./testdata/jsonplaceholder.yaml")
	defer func() {
		server.Close()
		shutdown()
	}()

	requestURL := server.URL + "/ddn/pre-route"
	testCases := []struct {
		Name         string
		Body         ddn.PreRoutePluginRequestBody
		StatusCode   int
		ResponseBody any
	}{
		{
			Name: "getAlbums",
			Body: ddn.PreRoutePluginRequestBody{
				Path:   "/albums",
				Method: "GET",
			},
			StatusCode: 200,
		},
		{
			Name: "getPostByID",
			Body: ddn.PreRoutePluginRequestBody{
				Path:   "/posts/1",
				Method: "GET",
			},
			StatusCode: 200,
			ResponseBody: map[string]any{
				"userId": float64(1),
				"id":     float64(1),
				"title":  "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
				"body":   "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			runPreRoute(t, requestURL, tc.Body, tc.StatusCode, tc.ResponseBody)
		})
	}
}

func TestGraphQLServer(t *testing.T) {
	server, shutdown := initTestServer(t, "./testdata/rickandmortyapi.yaml")
	defer func() {
		server.Close()
		shutdown()
	}()

	requestURL := server.URL + "/ddn/pre-route"
	testCases := []struct {
		Name         string
		Body         ddn.PreRoutePluginRequestBody
		StatusCode   int
		ResponseBody any
	}{
		{
			Name: "getCharacters",
			Body: ddn.PreRoutePluginRequestBody{
				Path:   "/characters",
				Method: "GET",
			},
			StatusCode: 200,
		},
		{
			Name: "getCharacterByID",
			Body: ddn.PreRoutePluginRequestBody{
				Path:   "/characters/1",
				Method: "GET",
			},
			StatusCode: 200,
			ResponseBody: map[string]any{
				"data": map[string]any{
					"character": map[string]any{
						"id":   "1",
						"name": "Rick Sanchez",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			runPreRoute(t, requestURL, tc.Body, tc.StatusCode, tc.ResponseBody)
		})
	}
}

func runPreRoute[T any](t *testing.T, requestURL string, body ddn.PreRoutePluginRequestBody, statusCode int, responseBody T) {
	t.Helper()

	bodyBytes, err := json.Marshal(body)
	assert.NilError(t, err)

	resp, err := http.Post(requestURL, "application/json", bytes.NewReader(bodyBytes))
	assert.NilError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != statusCode {
		respBody, err := io.ReadAll(resp.Body)
		assert.NilError(t, err)

		t.Errorf("expected status code: %d; got: %d; response body: %s", statusCode, resp.StatusCode, string(respBody))
		t.FailNow()
	}

	assert.Equal(t, resp.StatusCode, statusCode)

	var output, empty T

	err = json.NewDecoder(resp.Body).Decode(&output)
	assert.NilError(t, err)

	// ignore empty expected response.
	if reflect.DeepEqual(responseBody, empty) {
		return
	}

	assert.DeepEqual(t, responseBody, output)
}

func initTestServer(t *testing.T, configPath string) (*httptest.Server, func()) {
	t.Setenv("RELYX_CONFIG_PATH", configPath)

	envVars, err := config.LoadServerConfig()
	assert.NilError(t, err)

	otelExporters := &gotel.OTelExporters{
		Tracer: gotel.NewTracer("test"),
		Meter:  otel.Meter("test"),
		Logger: slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	}

	router, shutdown, err := setupRouter(envVars, otelExporters)
	assert.NilError(t, err)

	server := httptest.NewServer(router)

	return server, shutdown
}
