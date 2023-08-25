package grpcgateway

import (
	"context"
	"fmt"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/auditumio/auditum/pkg/fragma/grpcx"
)

// RegistrableService describes a service implementation that can be registered
// in the gRPC-Gateway server.
// gRPC service server implementations must implement this interface in order
// to be registered in gRPC-Gateway.
type RegistrableService interface {
	RegisterGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}

type Gateway struct {
	log      *zap.Logger
	services map[string][]RegistrableService
	handler  *http.ServeMux
}

func NewGateway(log *zap.Logger, opts ...GatewayOption) *Gateway {
	log = log.Named("grpc_gateway")

	rootHandler := http.NewServeMux()
	rootHandler.Handle("/metrics", promhttp.Handler())

	gateway := &Gateway{
		log:      log,
		services: make(map[string][]RegistrableService),
		handler:  rootHandler,
	}

	for _, opt := range opts {
		opt(gateway)
	}

	return gateway
}

type GatewayOption func(*Gateway)

func WithRegistrableServices(basePath string, services ...RegistrableService) GatewayOption {
	return func(g *Gateway) {
		basePath = strings.TrimSuffix(basePath, "/")
		g.services[basePath] = append(g.services[basePath], services...)
	}
}

func (g *Gateway) Handler() http.Handler {
	return g.handler
}

func (g *Gateway) ConnectAndRegister(addr string) error {
	g.log.Info("Connecting to gRPC server...")

	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}

	conn, err := grpcx.NewClientConnection(addr)
	if err != nil {
		return fmt.Errorf("connect to gRPC server: %v", err)
	}

	g.log.Info("Connected to gRPC server")

	ctx := context.Background()
	for basePath, services := range g.services {
		mux := newGatewayMux(g.log)

		for _, service := range services {
			if err := service.RegisterGateway(ctx, mux, conn); err != nil {
				return fmt.Errorf("register %T: %v", service, err)
			}
		}

		basePathPattern := basePath + "/"

		var basePathHandler http.Handler
		basePathHandler = mux
		basePathHandler = prettyMiddleware(basePathHandler)
		basePathHandler = otelhttp.NewHandler(
			basePathHandler,
			"grpc-gateway "+basePath,
		)

		g.handler.Handle(basePathPattern, http.StripPrefix(basePath, basePathHandler))
	}

	g.log.Info("Registered gRPC services in gRPC-Gateway server")

	return nil
}

func newGatewayMux(log *zap.Logger) *runtime.ServeMux {
	muxOpts := []runtime.ServeMuxOption{
		// We use wildcard as fallback, so users are not forced to specify
		// "Accept: application/json" header.
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			getMarshaler(false),
		),
		runtime.WithMarshalerOption(
			"application/json+pretty",
			getMarshaler(true),
		),
		runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher()),
		runtime.WithOutgoingHeaderMatcher(outgoingHeaderMatcher(log)),
		runtime.WithUnescapingMode(runtime.UnescapingModeAllExceptReserved),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			pattern, ok := runtime.HTTPPathPattern(ctx)
			if ok {
				return metadata.Pairs("http-path-pattern", pattern)
			}

			return nil
		}),
	}

	return runtime.NewServeMux(muxOpts...)
}

func getMarshaler(pretty bool) runtime.Marshaler {
	multiline := false
	indent := ""
	if pretty {
		multiline = true
		indent = "  "
	}

	return &runtime.HTTPBodyMarshaler{
		Marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				Multiline:    multiline,
				Indent:       indent,
				AllowPartial: false,
				// Use snake_case for field names instead of camelCase.
				UseProtoNames:  true,
				UseEnumNumbers: false,
				// Do not emit fields with zero values.
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				AllowPartial: false,
				// At the moment it is less useful for users but more practical
				// for us to discard unknown fields. At least until we can return
				// a pretty message.
				// Currently, the returned error looks like this:
				// {"code":3, "message":"proto: (line 1:2): unknown field \"prosject\""}
				DiscardUnknown: true,
			},
		},
	}
}

func incomingHeaderMatcher() runtime.HeaderMatcherFunc {
	return func(key string) (string, bool) {
		key = textproto.CanonicalMIMEHeaderKey(key)
		if key == "X-Request-Id" {
			return key, true
		}

		return runtime.DefaultHeaderMatcher(key)
	}
}

func outgoingHeaderMatcher(log *zap.Logger) runtime.HeaderMatcherFunc {
	return func(key string) (string, bool) {
		key = textproto.CanonicalMIMEHeaderKey(key)
		if key == "Content-Type" {
			// Drop "Grpc-Metadata-Content-Type: application/grpc".
			return "", false
		}

		log.Warn("Drop unknown outgoing header",
			zap.String("header_name", key),
		)
		return "", false
	}
}

func prettyMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Checking Values as map[string][]string also catches ?pretty and ?pretty=
		// r.URL.Query().Get("pretty") would not.
		if _, ok := r.URL.Query()["pretty"]; ok {
			r.Header.Set("Accept", "application/json+pretty")
		}
		handler.ServeHTTP(w, r)
	})
}
