package restrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/hasura/gotel"
	"github.com/relychan/relixy/config"
	"github.com/relychan/relixy/routes/ddnrouter"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

var testPlaceholderConfig = "../testdata/jsonplaceholder/config.yaml"

func TestRESTHandler_RESTServer(t *testing.T) {
	server, shutdown := initTestServer(t, testPlaceholderConfig)
	defer func() {
		server.Close()
		shutdown()
	}()

	testCases := []struct {
		Name         string
		Body         ddnrouter.PreRoutePluginRequestBody
		StatusCode   int
		ResponseBody any
	}{
		{
			Name: "getAlbums",
			Body: ddnrouter.PreRoutePluginRequestBody{
				Path:   server.URL + "/api/v1/albums",
				Method: "GET",
			},
			StatusCode: 200,
		},
		{
			Name: "getPostByID",
			Body: ddnrouter.PreRoutePluginRequestBody{
				Path:   server.URL + "/api/v1/posts/1",
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
			runSuccessRequest(t, tc.Body, tc.StatusCode, tc.ResponseBody)
			runUnauthorizedRequest(t, tc.Body, tc.StatusCode, tc.ResponseBody)
		})
	}
}

func TestRESTHandler_GraphQLServer(t *testing.T) {
	server, shutdown := initTestServer(t, "../testdata/rickandmortyapi/config.yaml")
	defer func() {
		server.Close()
		shutdown()
	}()

	testCases := []struct {
		Name         string
		Body         ddnrouter.PreRoutePluginRequestBody
		StatusCode   int
		ResponseBody any
	}{
		{
			Name: "getCharacters",
			Body: ddnrouter.PreRoutePluginRequestBody{
				Path:   server.URL + "/characters",
				Method: "GET",
			},
			StatusCode: 200,
		},
		{
			Name: "getCharacterByID",
			Body: ddnrouter.PreRoutePluginRequestBody{
				Path:   server.URL + "/characters/1",
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
			runSuccessRequest(t, tc.Body, tc.StatusCode, tc.ResponseBody)
			runUnauthorizedRequest(t, tc.Body, tc.StatusCode, tc.ResponseBody)
		})
	}
}

func TestRESTHandler_DDN(t *testing.T) {
	// os.Setenv("GRAPHQL_SERVER_URL", "http://localhost:3280/graphql")
	graphqlURL := os.Getenv("GRAPHQL_SERVER_URL")
	if graphqlURL == "" {
		return
	}

	os.Setenv("HASURA_M_AUTH", "test-secret")

	server, shutdown := initTestServer(t, "../../tests/ddn/config/plugin/config.yaml")
	defer func() {
		server.Close()
		shutdown()
	}()

	testCases := []struct {
		Name         string
		Body         ddnrouter.PreRoutePluginRequestBody
		StatusCode   int
		ResponseBody any
	}{
		{
			Name: "getArtistByName",
			Body: ddnrouter.PreRoutePluginRequestBody{
				Path:   server.URL + "/v1/api/rest/artistbyname/Queen",
				Method: "GET",
			},
			StatusCode: 200,
			ResponseBody: map[string]any{
				"data": map[string]any{
					"artist": []any{
						map[string]any{
							"name": "Queen",
						},
					},
				},
			},
		},
		{
			Name: "getArtists",
			Body: ddnrouter.PreRoutePluginRequestBody{
				Path:   server.URL + "/v1/api/rest/artists?limit=10&offset=20",
				Method: "GET",
			},
			StatusCode: 200,
			ResponseBody: map[string]any{
				"data": map[string]any{
					"artist": []any{
						map[string]any{
							"name": "Various Artists",
						},
						map[string]any{
							"name": "Led Zeppelin",
						},
						map[string]any{
							"name": "Frank Zappa & Captain Beefheart",
						},
						map[string]any{
							"name": "Marcos Valle",
						},
						map[string]any{
							"name": "Milton Nascimento & Bebeto",
						},
						map[string]any{
							"name": "Azymuth",
						},
						map[string]any{
							"name": "Gilberto Gil",
						},
						map[string]any{
							"name": "João Gilberto",
						},
						map[string]any{
							"name": "Bebel Gilberto",
						},
						map[string]any{
							"name": "Jorge Vercilo",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		runSuccessRequest(t, tc.Body, tc.StatusCode, tc.ResponseBody)
	}
}

func runSuccessRequest[T any](t *testing.T, r ddnrouter.PreRoutePluginRequestBody, statusCode int, responseBody T) {
	t.Helper()

	t.Run("success", func(t *testing.T) {
		var reader io.ReadSeeker
		if len(r.Body) > 0 {
			reader = bytes.NewReader(r.Body)
		}

		req, err := http.NewRequest(r.Method, r.Path, reader)
		require.NoError(t, err)

		req.Header.Set("hasura-m-auth", "test-secret")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode != statusCode {
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			t.Errorf("expected status code: %d; got: %d; response body: %s", statusCode, resp.StatusCode, string(respBody))
			t.FailNow()
		}

		var output, empty T

		err = json.NewDecoder(resp.Body).Decode(&output)
		require.NoError(t, err)

		// ignore empty expected response.
		if reflect.DeepEqual(responseBody, empty) {
			return
		}

		require.Equal(t, responseBody, output)
	})
}

func runUnauthorizedRequest[T any](t *testing.T, r ddnrouter.PreRoutePluginRequestBody, _ int, _ T) {
	t.Helper()
	// expected unauthorized status
	t.Run("unauthorized", func(t *testing.T) {
		var reader io.ReadSeeker
		if len(r.Body) > 0 {
			reader = bytes.NewReader(r.Body)
		}

		req, err := http.NewRequest(r.Method, r.Path, reader)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected unauthorized status; got: %d", resp.StatusCode)
			t.FailNow()
		}
	})
}

func TestSetupRouter_InvalidConfig(t *testing.T) {
	t.Setenv("RELIXY_CONFIG_PATH", "../testdata/invalid-config.yaml")

	_, _, err := config.LoadServerConfig(context.Background())
	require.ErrorIs(t, err, io.EOF)
}

func TestSetupRouter_ValidConfig(t *testing.T) {
	t.Setenv("RELIXY_CONFIG_PATH", testPlaceholderConfig)

	envVars, logger, err := config.LoadServerConfig(context.Background())
	require.NoError(t, err)

	otelExporters := &gotel.OTelExporters{
		Tracer: gotel.NewTracer("test"),
		Meter:  otel.Meter("test"),
		Logger: logger,
	}

	router, shutdown, err := SetupRouter(context.TODO(), envVars, otelExporters)
	require.NoError(t, err)
	require.True(t, router != nil)
	require.True(t, shutdown != nil)

	shutdown()
}

func TestRESTHandler_NotFoundPath(t *testing.T) {
	server, shutdown := initTestServer(t, testPlaceholderConfig)
	defer func() {
		server.Close()
		shutdown()
	}()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/v1/nonexistent", nil)
	require.NoError(t, err)

	req.Header.Set("hasura-m-auth", "test-secret")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRESTHandler_WithPathParams(t *testing.T) {
	server, shutdown := initTestServer(t, testPlaceholderConfig)
	defer func() {
		server.Close()
		shutdown()
	}()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/v1/posts/1", nil)
	require.NoError(t, err)

	req.Header.Set("hasura-m-auth", "test-secret")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	require.Equal(t, float64(1), result["id"])
}

func TestRESTHandler_GetAlbums(t *testing.T) {
	server, shutdown := initTestServer(t, testPlaceholderConfig)
	defer func() {
		server.Close()
		shutdown()
	}()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/v1/albums", nil)
	require.NoError(t, err)

	req.Header.Set("hasura-m-auth", "test-secret")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result []map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	require.True(t, len(result) > 0)
}

func initTestServer(t *testing.T, configPath string) (*httptest.Server, func()) {
	t.Setenv("RELIXY_CONFIG_PATH", configPath)

	envVars, logger, err := config.LoadServerConfig(context.Background())
	require.NoError(t, err)

	otelExporters := &gotel.OTelExporters{
		Tracer: gotel.NewTracer("test"),
		Meter:  otel.Meter("test"),
		Logger: logger,
	}

	router, shutdown, err := SetupRouter(context.TODO(), envVars, otelExporters)
	require.NoError(t, err)

	server := httptest.NewServer(router)

	return server, shutdown
}
