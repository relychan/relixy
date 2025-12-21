// Package proxyc implements a client to proxy requests to external services.
package proxyc

import (
	"context"
	"fmt"

	"github.com/relychan/gohttpc"
	"github.com/relychan/gohttpc/loadbalancer"
	"github.com/relychan/gohttpc/loadbalancer/roundrobin"
	"github.com/relychan/relyx/proxyc/internal"
	"github.com/relychan/relyx/schema"
)

// ProxyClient helps manage and execute REST and GraphQL APIs from the API document.
type ProxyClient struct {
	clientOptions  *ProxyClientOptions
	lbClient       *loadbalancer.LoadBalancerClient
	metadata       *schema.RelyProxyAPIDocument
	node           *internal.Node
	defaultHeaders map[string]string
}

// ProxyClientOptions holds optional options to creates a proxy client.
type ProxyClientOptions struct {
	*gohttpc.ClientOptions
}

// NewProxyClient creates a proxy client from the API document.
func NewProxyClient(
	ctx context.Context,
	metadata *schema.RelyProxyAPIDocument,
	clientOptions *ProxyClientOptions,
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
	err := pc.initServers()
	if err != nil {
		return err
	}

	err = pc.initDefaultHeaders()
	if err != nil {
		return err
	}

	node, err := BuildMetadataTree(ctx, pc.metadata, pc.clientOptions)
	if err != nil {
		return err
	}

	pc.node = node

	return nil
}

func (pc *ProxyClient) initDefaultHeaders() error {
	for key, envValue := range pc.metadata.Settings.Headers {
		value, err := envValue.GetOrDefault("")
		if err != nil {
			return fmt.Errorf("failed to load header %s: %w", key, err)
		}

		if value != "" {
			pc.defaultHeaders[key] = value
		}
	}

	return nil
}

func (pc *ProxyClient) initServers() error {
	if len(pc.metadata.Servers) == 0 {
		return errServerURLRequired
	}

	var err error

	var healthCheckBuilder *loadbalancer.HTTPHealthCheckPolicyBuilder

	if pc.metadata.Settings.HealthCheck != nil && pc.metadata.Settings.HealthCheck.HTTP != nil {
		healthCheckBuilder, err = pc.metadata.Settings.HealthCheck.HTTP.ToPolicyBuilder()
		if err != nil {
			return err
		}
	} else {
		healthCheckBuilder = loadbalancer.NewHTTPHealthCheckPolicyBuilder()
	}

	switch len(pc.metadata.Servers) {
	case 0:
		return errServerURLRequired
	case 1:
		host, err := pc.initServer(&pc.metadata.Servers[0], healthCheckBuilder)
		if err != nil {
			return err
		}

		if host == nil {
			return ErrNoAvailableServer
		}

		wrr, err := roundrobin.NewWeightedRoundRobin([]*loadbalancer.Host{host})
		if err != nil {
			return err
		}

		pc.lbClient = loadbalancer.NewLoadBalancerClientWithOptions(wrr, pc.clientOptions)

		return nil
	default:
		hosts := make([]*loadbalancer.Host, 0, len(pc.metadata.Servers))

		for _, server := range pc.metadata.Servers {
			host, err := pc.initServer(&server, healthCheckBuilder)
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
	}

	return nil
}

func (pc *ProxyClient) initServer(
	server *schema.RelyProxyServer,
	healthCheckBuilder *loadbalancer.HTTPHealthCheckPolicyBuilder,
) (*loadbalancer.Host, error) {
	rawServerURL, err := server.URL.GetOrDefault("")
	if err != nil {
		return nil, err
	}

	if rawServerURL == "" {
		return nil, nil
	}

	host, err := loadbalancer.NewHost(
		pc.clientOptions.HTTPClient,
		rawServerURL,
		loadbalancer.WithHTTPHealthCheckPolicyBuilder(healthCheckBuilder),
	)
	if err != nil {
		return nil, err
	}

	if server.Name != "" {
		host.SetName(server.Name)
	}

	if server.Weight != nil && *server.Weight > 1 {
		host.SetWeight(*server.Weight)
	}

	if len(server.Headers) > 0 {
		headers := make(map[string]string)

		for key, header := range server.Headers {
			value, err := header.GetOrDefault("")
			if err != nil {
				return nil, fmt.Errorf("failed to get header %s: %w", key, err)
			}

			if value != "" {
				headers[key] = value
			}
		}

		host.SetHeaders(headers)
	}

	return host, nil
}
