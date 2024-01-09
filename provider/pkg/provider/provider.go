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
	"google.golang.org/protobuf/types/known/emptypb"

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

// GetSchema fetches the schema for this resource provider.
func (p *xyzProvider) GetSchema(_ context.Context, _ *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	return &pulumirpc.GetSchemaResponse{schema: p.schema}, nil
}

// TODO: Implement the remaining functions from the `ResourceProviderServer`
// interface below. That is, Create, Read, Diff, Update, Delete and Inokve
// at a minimum.

// CheckConfig validates the configuration for this resource provider.
func (p *xyzProvider) CheckConfig(_ context.Context, _ *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return nil, nil
}

// DiffConfig checks the impact a hypothetical change to this provider's configuration will have on the provider.
func (p *xyzProvider) DiffConfig(_ context.Context, _ *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return nil, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (p *xyzProvider) Invoke(_ context.Context, _ *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	return nil, nil
}

// StreamInvoke dynamically executes a built-in function in the provider, which returns a stream
// of responses.
func (p *xyzProvider) StreamInvoke(_ *pulumirpc.InvokeRequest, _ pulumirpc.ResourceProvider_StreamInvokeServer) error {
	return nil
}

// Call dynamically executes a method in the provider associated with a component resource.
func (p *xyzProvider) Call(_ context.Context, _ *pulumirpc.CallRequest) (*pulumirpc.CallResponse, error) {
	return nil, nil
}

// Check validates that the given property bag is valid for a resource of the given type and returns the inputs
// that should be passed to successive calls to Diff, Create, or Update for this resource. As a rule, the provider
// inputs returned by a call to Check should preserve the original representation of the properties as present in
// the program inputs. Though this rule is not required for correctness, violations thereof can negatively impact
// the end-user experience, as the provider inputs are using for detecting and rendering diffs.
func (p *xyzProvider) Check(_ context.Context, _ *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return nil, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (p *xyzProvider) Diff(_ context.Context, _ *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return nil, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.  (The input ID
// must be blank.)  If this call fails, the resource must not have been created (i.e., it is "transactional").
func (p *xyzProvider) Create(_ context.Context, _ *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	return nil, nil
}

// Read the current live state associated with a resource.  Enough state must be include in the inputs to uniquely
// identify the resource; this is typically just the resource ID, but may also include some properties.
func (p *xyzProvider) Read(_ context.Context, _ *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	return nil, nil
}

// Update updates an existing resource with new values.
func (p *xyzProvider) Update(_ context.Context, _ *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	return nil, nil
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed to still exist.
func (p *xyzProvider) Delete(_ context.Context, _ *pulumirpc.DeleteRequest) (*emptypb.Empty, error) {
	return nil, nil
}

// Construct creates a new instance of the provided component resource and returns its state.
func (p *xyzProvider) Construct(_ context.Context, _ *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	return nil, nil
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
// Operations aborted in this way will return an error (e.g., `Update` and `Create` will either return a
// creation error or an initialization error). Since Cancel is advisory and non-blocking, it is up
// to the host to decide how long to wait after Cancel is called before (e.g.)
// hard-closing any gRPC connection.
func (p *xyzProvider) Cancel(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (p *xyzProvider) GetPluginInfo(_ context.Context, _ *emptypb.Empty) (*pulumirpc.PluginInfo, error) {
	return nil, nil
}

// Attach sends the engine address to an already running plugin.
func (p *xyzProvider) Attach(_ context.Context, _ *pulumirpc.PluginAttach) (*emptypb.Empty, error) {
	return nil, nil
}

// GetMapping fetches the mapping for this resource provider, if any. A provider should return an empty
// response (not an error) if it doesn't have a mapping for the given key.
func (p *xyzProvider) GetMapping(_ context.Context, _ *pulumirpc.GetMappingRequest) (*pulumirpc.GetMappingResponse, error) {
	return nil, nil
}

// GetMappings is an optional method that returns what mappings (if any) a provider supports. If a provider does not
// implement this method the engine falls back to the old behaviour of just calling GetMapping without a name.
// If this method is implemented than the engine will then call GetMapping only with the names returned from this method.
func (p *xyzProvider) GetMappings(_ context.Context, _ *pulumirpc.GetMappingsRequest) (*pulumirpc.GetMappingsResponse, error) {
	return nil, nil
}
