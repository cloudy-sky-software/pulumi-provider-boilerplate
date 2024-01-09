package provider

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cloudy-sky-software/pulumi-xyz/provider/pkg/openapi"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/pkg/errors"

	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	providerGen "github.com/cloudy-sky-software/pulschema/pkg"
)

// Provider implements Pulumi's `ResourceProviderServer` interface.
// The implemented methods assume that the cloud provider supports RESTful
// APIs that accept a content-type of `application/json`.
type xyzProvider struct {
	pulumirpc.UnimplementedResourceProviderServer

	host    *provider.HostClient
	name    string
	version string

	metadata providerGen.ProviderMetadata

	baseURL    string
	httpClient *http.Client
	openAPIDoc openapi3.T
	schema     pschema.PackageSpec

	// TODO: Customize the auth requirements to suit your needs.
	apiKey string
}

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}

// makeProvider returns an instance of the REST-based resource provider handler.
func makeProvider(host *provider.HostClient, name, version string, pulumiSchemaBytes, openapiDocBytes, metadataBytes []byte) (pulumirpc.ResourceProviderServer, error) {
	openapiDoc := openapi.GetOpenAPISpec(openapiDocBytes)

	var metadata providerGen.ProviderMetadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil, errors.Wrap(err, "unmarshaling the metadata bytes to json")
	}

	httpClient := &http.Client{
		// The transport is mostly a copy of the http.DefaultTransport
		// with the exception of ForceAttemptHTTP2 set to false.
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: defaultTransportDialContext(&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}),
			ForceAttemptHTTP2:     false,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("unable to handle redirects")
		},
	}

	var pulumiSchema pschema.PackageSpec
	if err := json.Unmarshal(pulumiSchemaBytes, &pulumiSchema); err != nil {
		return nil, errors.Wrap(err, "unmarshaling pulumi schema into its package spec form")
	}

	// Return the new provider
	return &xyzProvider{
		host:       host,
		name:       name,
		version:    version,
		schema:     pulumiSchema,
		baseURL:    openapiDoc.Servers[0].URL,
		openAPIDoc: *openapiDoc,
		metadata:   metadata,
		httpClient: httpClient,
	}, nil
}

// Configure is called by the engine to give you a chance to prepare the
// provider to handle resource requests, such as configuring auth etc.
//
// TODO: Implement the remaining functions from the `ResourceProviderServer` interface.
// That is, Create, Read, Diff, Update, Delete at a very minimum.
func (p *xyzProvider) Configure(_ context.Context, req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	apiKey, ok := req.GetVariables()["xyz:config:apiKey"]
	if !ok {
		// Check if it's set as an env var.
		envVarNames := p.schema.Provider.InputProperties["apiKey"].DefaultInfo.Environment
		for _, n := range envVarNames {
			v := os.Getenv(n)
			if v != "" {
				apiKey = v
			}
		}

		// Return an error if the API key is still empty.
		if apiKey == "" {
			return nil, errors.New("api key is required")
		}
	}

	logging.V(3).Info("Configuring XYZ API key")
	p.apiKey = apiKey

	return &pulumirpc.ConfigureResponse{
		AcceptSecrets: true,
	}, nil
}
