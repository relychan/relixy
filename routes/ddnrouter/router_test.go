package ddnrouter

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
	"github.com/relychan/goutils/httpheader"
	"github.com/relychan/relixy/config"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func TestRestifiedPlugin_RESTServer(t *testing.T) {
	server, shutdown := initTestServer(t, "../testdata/jsonplaceholder/config.yaml")
	defer func() {
		server.Close()
		shutdown()
	}()

	requestURL := server.URL + "/ddn/pre-route"
	testCases := []struct {
		Name         string
		Body         PreRoutePluginRequestBody
		StatusCode   int
		ResponseBody any
	}{
		{
			Name: "getAlbums",
			Body: PreRoutePluginRequestBody{
				Path:   "/api/v1/albums",
				Method: "GET",
			},
			StatusCode: 200,
		},
		{
			Name: "countAlbums",
			Body: PreRoutePluginRequestBody{
				Path:   "/api/v1/albums-count",
				Method: "POST",
			},
			StatusCode: 200,
			ResponseBody: map[string]any{
				"count": float64(100),
			},
		},
		{
			Name: "getPostByID",
			Body: PreRoutePluginRequestBody{
				Path:   "/api/v1/posts/1",
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

func TestRestifiedPlugin_GraphQLServer(t *testing.T) {
	server, shutdown := initTestServer(t, "../testdata/rickandmortyapi/config.yaml")
	defer func() {
		server.Close()
		shutdown()
	}()

	requestURL := server.URL + "/ddn/pre-route"
	testCases := []struct {
		Name         string
		Body         PreRoutePluginRequestBody
		StatusCode   int
		ResponseBody any
	}{
		{
			Name: "getCharacters",
			Body: PreRoutePluginRequestBody{
				Path:   "/characters",
				Method: "GET",
			},
			StatusCode: 200,
		},
		{
			Name: "getCharacterByID",
			Body: PreRoutePluginRequestBody{
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

func runPreRoute[T any](t *testing.T, requestURL string, body PreRoutePluginRequestBody, statusCode int, responseBody T) {
	t.Helper()

	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(bodyBytes))
		require.NoError(t, err)

		req.Header.Set(httpheader.ContentType, httpheader.ContentTypeJSON)
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

	// expected unauthorized status
	t.Run("unauthorized", func(t *testing.T) {
		resp, err := http.DefaultClient.Post(requestURL, httpheader.ContentTypeJSON, bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected unauthorized status; got: %d", resp.StatusCode)
			t.FailNow()
		}
	})
}

func TestRestifiedPlugin_DDN(t *testing.T) {
	// os.Setenv("DDN_ENGINE_HOST", "http://localhost:3280")
	engineHost := os.Getenv("DDN_ENGINE_HOST")
	if engineHost == "" {
		return
	}

	testCases := []struct {
		Path     string
		Method   string
		Body     any
		Expected any
	}{
		{
			Path:   "/v1/api/rest/artistbyname/Queen",
			Method: http.MethodGet,
			Expected: map[string]any{
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
			Path:   "/v1/api/rest/artists?limit=10&offset=20",
			Method: http.MethodGet,
			Expected: map[string]any{
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
		req, err := http.NewRequest(tc.Method, engineHost+tc.Path, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		defer resp.Body.Close()

		var respBody any

		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		require.Equal(t, tc.Expected, respBody)
	}
}

func TestSetupRouter_InvalidConfig(t *testing.T) {
	t.Setenv("RELIXY_CONFIG_PATH", "../testdata/invalid-config.yaml")

	_, _, err := config.LoadServerConfig(context.Background())
	require.ErrorIs(t, err, io.EOF)
}

func TestSetupRouter_ValidConfig(t *testing.T) {
	t.Setenv("RELIXY_CONFIG_PATH", "../testdata/jsonplaceholder/config.yaml")

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

func TestPreRoutePlugin_InvalidJSON(t *testing.T) {
	server, shutdown := initTestServer(t, "../testdata/jsonplaceholder/config.yaml")
	defer func() {
		server.Close()
		shutdown()
	}()

	requestURL := server.URL + "/ddn/pre-route"

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBufferString("invalid json"))
	require.NoError(t, err)

	req.Header.Set(httpheader.ContentType, httpheader.ContentTypeJSON)
	req.Header.Set("hasura-m-auth", "test-secret")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPreRoutePlugin_InvalidContentType(t *testing.T) {
	server, shutdown := initTestServer(t, "../testdata/jsonplaceholder/config.yaml")
	defer func() {
		server.Close()
		shutdown()
	}()

	requestURL := server.URL + "/ddn/pre-route"
	body := PreRoutePluginRequestBody{
		Path:   "/albums",
		Method: "GET",
	}

	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(bodyBytes))
	require.NoError(t, err)

	req.Header.Set(httpheader.ContentType, "text/plain")
	req.Header.Set("hasura-m-auth", "test-secret")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
}

func TestPreRoutePlugin_NotFoundPath(t *testing.T) {
	server, shutdown := initTestServer(t, "../testdata/jsonplaceholder/config.yaml")
	defer func() {
		server.Close()
		shutdown()
	}()

	requestURL := server.URL + "/ddn/pre-route"
	body := PreRoutePluginRequestBody{
		Path:   "/nonexistent",
		Method: "GET",
	}

	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(bodyBytes))
	require.NoError(t, err)

	req.Header.Set(httpheader.ContentType, httpheader.ContentTypeJSON)
	req.Header.Set("hasura-m-auth", "test-secret")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestPreRoutePlugin_GetAlbums(t *testing.T) {
	server, shutdown := initTestServer(t, "../testdata/jsonplaceholder/config.yaml")
	defer func() {
		server.Close()
		shutdown()
	}()

	requestURL := server.URL + "/ddn/pre-route"
	body := PreRoutePluginRequestBody{
		Path:   "/api/v1/albums",
		Method: "GET",
	}

	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(bodyBytes))
	require.NoError(t, err)

	req.Header.Set(httpheader.ContentType, httpheader.ContentTypeJSON)
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
