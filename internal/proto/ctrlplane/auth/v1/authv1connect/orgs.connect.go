// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: ctrlplane/auth/v1/orgs.proto

package authv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// OrgServiceName is the fully-qualified name of the OrgService service.
	OrgServiceName = "ctrlplane.auth.v1.OrgService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// OrgServiceCreateOrgProcedure is the fully-qualified name of the OrgService's CreateOrg RPC.
	OrgServiceCreateOrgProcedure = "/ctrlplane.auth.v1.OrgService/CreateOrg"
	// OrgServiceGetOrgByIDProcedure is the fully-qualified name of the OrgService's GetOrgByID RPC.
	OrgServiceGetOrgByIDProcedure = "/ctrlplane.auth.v1.OrgService/GetOrgByID"
	// OrgServiceSetOrgHooksProcedure is the fully-qualified name of the OrgService's SetOrgHooks RPC.
	OrgServiceSetOrgHooksProcedure = "/ctrlplane.auth.v1.OrgService/SetOrgHooks"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	orgServiceServiceDescriptor           = v1.File_ctrlplane_auth_v1_orgs_proto.Services().ByName("OrgService")
	orgServiceCreateOrgMethodDescriptor   = orgServiceServiceDescriptor.Methods().ByName("CreateOrg")
	orgServiceGetOrgByIDMethodDescriptor  = orgServiceServiceDescriptor.Methods().ByName("GetOrgByID")
	orgServiceSetOrgHooksMethodDescriptor = orgServiceServiceDescriptor.Methods().ByName("SetOrgHooks")
)

// OrgServiceClient is a client for the ctrlplane.auth.v1.OrgService service.
type OrgServiceClient interface {
	// CreateOrg creates a new organization.
	CreateOrg(context.Context, *connect.Request[v1.CreateOrgRequest]) (*connect.Response[v1.CreateOrgResponse], error)
	// GetOrgByID retrieves an organization by its globally unique identifier.
	GetOrgByID(context.Context, *connect.Request[v1.GetOrgByIDRequest]) (*connect.Response[v1.GetOrgByIDResponse], error)
	// SetOrgHooks sets the hooks for an organization.
	SetOrgHooks(context.Context, *connect.Request[v1.SetOrgHooksRequest]) (*connect.Response[emptypb.Empty], error)
}

// NewOrgServiceClient constructs a client for the ctrlplane.auth.v1.OrgService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewOrgServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) OrgServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &orgServiceClient{
		createOrg: connect.NewClient[v1.CreateOrgRequest, v1.CreateOrgResponse](
			httpClient,
			baseURL+OrgServiceCreateOrgProcedure,
			connect.WithSchema(orgServiceCreateOrgMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getOrgByID: connect.NewClient[v1.GetOrgByIDRequest, v1.GetOrgByIDResponse](
			httpClient,
			baseURL+OrgServiceGetOrgByIDProcedure,
			connect.WithSchema(orgServiceGetOrgByIDMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		setOrgHooks: connect.NewClient[v1.SetOrgHooksRequest, emptypb.Empty](
			httpClient,
			baseURL+OrgServiceSetOrgHooksProcedure,
			connect.WithSchema(orgServiceSetOrgHooksMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// orgServiceClient implements OrgServiceClient.
type orgServiceClient struct {
	createOrg   *connect.Client[v1.CreateOrgRequest, v1.CreateOrgResponse]
	getOrgByID  *connect.Client[v1.GetOrgByIDRequest, v1.GetOrgByIDResponse]
	setOrgHooks *connect.Client[v1.SetOrgHooksRequest, emptypb.Empty]
}

// CreateOrg calls ctrlplane.auth.v1.OrgService.CreateOrg.
func (c *orgServiceClient) CreateOrg(ctx context.Context, req *connect.Request[v1.CreateOrgRequest]) (*connect.Response[v1.CreateOrgResponse], error) {
	return c.createOrg.CallUnary(ctx, req)
}

// GetOrgByID calls ctrlplane.auth.v1.OrgService.GetOrgByID.
func (c *orgServiceClient) GetOrgByID(ctx context.Context, req *connect.Request[v1.GetOrgByIDRequest]) (*connect.Response[v1.GetOrgByIDResponse], error) {
	return c.getOrgByID.CallUnary(ctx, req)
}

// SetOrgHooks calls ctrlplane.auth.v1.OrgService.SetOrgHooks.
func (c *orgServiceClient) SetOrgHooks(ctx context.Context, req *connect.Request[v1.SetOrgHooksRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.setOrgHooks.CallUnary(ctx, req)
}

// OrgServiceHandler is an implementation of the ctrlplane.auth.v1.OrgService service.
type OrgServiceHandler interface {
	// CreateOrg creates a new organization.
	CreateOrg(context.Context, *connect.Request[v1.CreateOrgRequest]) (*connect.Response[v1.CreateOrgResponse], error)
	// GetOrgByID retrieves an organization by its globally unique identifier.
	GetOrgByID(context.Context, *connect.Request[v1.GetOrgByIDRequest]) (*connect.Response[v1.GetOrgByIDResponse], error)
	// SetOrgHooks sets the hooks for an organization.
	SetOrgHooks(context.Context, *connect.Request[v1.SetOrgHooksRequest]) (*connect.Response[emptypb.Empty], error)
}

// NewOrgServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewOrgServiceHandler(svc OrgServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	orgServiceCreateOrgHandler := connect.NewUnaryHandler(
		OrgServiceCreateOrgProcedure,
		svc.CreateOrg,
		connect.WithSchema(orgServiceCreateOrgMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	orgServiceGetOrgByIDHandler := connect.NewUnaryHandler(
		OrgServiceGetOrgByIDProcedure,
		svc.GetOrgByID,
		connect.WithSchema(orgServiceGetOrgByIDMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	orgServiceSetOrgHooksHandler := connect.NewUnaryHandler(
		OrgServiceSetOrgHooksProcedure,
		svc.SetOrgHooks,
		connect.WithSchema(orgServiceSetOrgHooksMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/ctrlplane.auth.v1.OrgService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case OrgServiceCreateOrgProcedure:
			orgServiceCreateOrgHandler.ServeHTTP(w, r)
		case OrgServiceGetOrgByIDProcedure:
			orgServiceGetOrgByIDHandler.ServeHTTP(w, r)
		case OrgServiceSetOrgHooksProcedure:
			orgServiceSetOrgHooksHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedOrgServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedOrgServiceHandler struct{}

func (UnimplementedOrgServiceHandler) CreateOrg(context.Context, *connect.Request[v1.CreateOrgRequest]) (*connect.Response[v1.CreateOrgResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ctrlplane.auth.v1.OrgService.CreateOrg is not implemented"))
}

func (UnimplementedOrgServiceHandler) GetOrgByID(context.Context, *connect.Request[v1.GetOrgByIDRequest]) (*connect.Response[v1.GetOrgByIDResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ctrlplane.auth.v1.OrgService.GetOrgByID is not implemented"))
}

func (UnimplementedOrgServiceHandler) SetOrgHooks(context.Context, *connect.Request[v1.SetOrgHooksRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ctrlplane.auth.v1.OrgService.SetOrgHooks is not implemented"))
}
