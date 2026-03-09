// Package proxyc implements a client to proxy requests to external services.
package proxyc

import (
	"context"
	"fmt"

	"github.com/hasura/goenvconf"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/relychan/gohttpc"
	"github.com/relychan/gohttpc/httpconfig"
	"github.com/relychan/gohttpc/loadbalancer"
	"github.com/relychan/gohttpc/loadbalancer/roundrobin"
	"github.com/relychan/relixy/proxyc/handler/proxyhandler"
	"github.com/relychan/relixy/proxyc/internal"
	"github.com/relychan/relixy/schema/openapi"
)

// ProxyClient helps manage and execute REST and GraphQL APIs from the API document.
type ProxyClient struct {
	clientOptions  *gohttpc.ClientOptions
	lbClient       *loadbalancer.LoadBalancerClient
	metadata       *openapi.RelixyOpenAPIResource
	node           *internal.Node
	defaultHeaders map[string]string
	authenticators *proxyhandler.OpenAPIAuthenticator
}

// NewProxyClient creates a proxy client from the API document.
func NewProxyClient(
	ctx context.Context,
	metadata *openapi.RelixyOpenAPIResource,
	clientOptions *gohttpc.ClientOptions,
) (*ProxyClient, error) {
	client := &ProxyClient{
		metadata:       metadata,
		clientOptions:  clientOptions,
		defaultHeaders: map[string]string{},
	}

	err := client.init(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Metadata returns the metadata of the proxy client.
func (pc *ProxyClient) Metadata() *openapi.RelixyOpenAPIResource {
	return pc.metadata
}

// Close method performs cleanup and closure activities on the client instance.
func (pc *ProxyClient) Close() error {
	if pc.clientOptions != nil && pc.clientOptions.HTTPClient != nil {
		pc.clientOptions.HTTPClient.CloseIdleConnections()
	}

	if pc.lbClient != nil {
		return pc.lbClient.Close()
	}

	return nil
}

func (pc *ProxyClient) init(ctx context.Context) error {
	spec, err := pc.metadata.Definition.Build(ctx)
	if err != nil {
		return err
	}

	err = pc.initHTTPClient()
	if err != nil {
		return err
	}

	err = pc.initServers(spec)
	if err != nil {
		return err
	}

	err = pc.initDefaultHeaders()
	if err != nil {
		return err
	}

	pc.authenticators, err = proxyhandler.NewOpenAPIv3Authenticator(
		spec,
		pc.clientOptions.GetEnvFunc(),
	)
	if err != nil {
		return err
	}

	node, err := BuildMetadataTree(spec, pc.clientOptions)
	if err != nil {
		return err
	}

	pc.node = node

	return nil
}

func (pc *ProxyClient) initDefaultHeaders() error {
	getEnv := pc.clientOptions.GetEnvFunc()

	for key, envValue := range pc.metadata.Definition.Settings.Headers {
		value, err := envValue.GetCustom(getEnv)
		if err != nil {
			return fmt.Errorf("failed to load header %s: %w", key, err)
		}

		if value != "" {
			pc.defaultHeaders[key] = value
		}
	}

	return nil
}

func (pc *ProxyClient) initServers(spec *highv3.Document) error {
	if len(spec.Servers) == 0 {
		return errServerURLRequired
	}

	var err error

	var healthCheckBuilder *loadbalancer.HTTPHealthCheckPolicyBuilder

	if pc.metadata.Definition.Settings.HealthCheck != nil &&
		pc.metadata.Definition.Settings.HealthCheck.HTTP != nil {
		healthCheckBuilder, err = pc.metadata.Definition.Settings.HealthCheck.HTTP.ToPolicyBuilder()
		if err != nil {
			return err
		}
	} else {
		healthCheckBuilder = loadbalancer.NewHTTPHealthCheckPolicyBuilder()
	}

	hosts := make([]*loadbalancer.Host, 0, len(spec.Servers))

	for _, server := range spec.Servers {
		host, err := pc.initServer(server, healthCheckBuilder)
		if err != nil {
			return err
		}

		if host != nil {
			hosts = append(hosts, host)
		}
	}

	if len(hosts) == 0 {
		return ErrNoAvailableServer
	}

	wrr, err := roundrobin.NewWeightedRoundRobin(hosts)
	if err != nil {
		return err
	}

	pc.lbClient = loadbalancer.NewLoadBalancerClientWithOptions(wrr, pc.clientOptions)

	return nil
}

func (pc *ProxyClient) initServer(
	server *highv3.Server,
	healthCheckBuilder *loadbalancer.HTTPHealthCheckPolicyBuilder,
) (*loadbalancer.Host, error) {
	getEnv := pc.clientOptions.GetEnvFunc()

	serverURL, err := parseServerURL(server, getEnv)
	if err != nil {
		return nil, err
	}

	if serverURL == "" {
		return nil, nil
	}

	host, err := loadbalancer.NewHost(
		pc.clientOptions.HTTPClient,
		serverURL,
		loadbalancer.WithHTTPHealthCheckPolicyBuilder(healthCheckBuilder),
	)
	if err != nil {
		return nil, err
	}

	if server.Name != "" {
		host.SetName(server.Name)
	}

	rawWeight, exist := server.Extensions.Get(openapi.XRelyServerWeight)
	if exist && rawWeight != nil {
		var weight int

		err := rawWeight.Decode(&weight)
		if err != nil {
			return nil, fmt.Errorf("failed to decode weight from server: %w", err)
		}

		if weight > 1 {
			host.SetWeight(weight)
		}
	}

	rawHeaders, exist := server.Extensions.Get(openapi.XRelyServerHeaders)
	if exist && rawHeaders != nil {
		headerEnvs := map[string]goenvconf.EnvString{}

		err := rawHeaders.Decode(&headerEnvs)
		if err != nil {
			return nil, fmt.Errorf("failed to decode headers from server: %w", err)
		}

		if len(headerEnvs) > 0 {
			headers := make(map[string]string)

			for key, header := range headerEnvs {
				value, err := header.GetCustom(getEnv)
				if err != nil {
					return nil, fmt.Errorf("failed to get header %s: %w", key, err)
				}

				if value != "" {
					headers[key] = value
				}
			}

			host.SetHeaders(headers)
		}
	}

	return host, nil
}

func (pc *ProxyClient) initHTTPClient() error {
	var httpConfig *httpconfig.HTTPClientConfig

	if pc.metadata.Definition.Settings != nil && pc.metadata.Definition.Settings.HTTP != nil {
		httpConfig = pc.metadata.Definition.Settings.HTTP
	} else if pc.clientOptions.HTTPClient == nil {
		httpConfig = new(httpconfig.HTTPClientConfig)
	}

	if httpConfig != nil {
		httpClient, err := httpconfig.NewHTTPClientFromConfig(
			httpConfig,
			pc.clientOptions,
		)
		if err != nil {
			return err
		}

		pc.clientOptions.HTTPClient = httpClient
	}

	return nil
}
