package grpcx

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServerController struct {
	addr     string
	unixAddr string
	server   *grpc.Server
	log      *zap.Logger
	errs     chan error
}

func NewServerController(
	addr string,
	server *grpc.Server,
	log *zap.Logger,
	opts ...ServerControllerOption,
) *ServerController {
	ctrl := &ServerController{
		addr:     addr,
		unixAddr: "",
		server:   server,
		log: log.Named("grpc_server_controller").
			With(zap.String("addr", addr)),
		errs: make(chan error, 1),
	}

	for _, opt := range opts {
		opt(ctrl)
	}

	return ctrl
}

type ServerControllerOption func(*ServerController)

func ServerControllerWithUnixSocket(path string) ServerControllerOption {
	return func(c *ServerController) {
		c.unixAddr = path
	}
}

func (sm *ServerController) Start(ctx context.Context) error {
	sm.log.Info("Starting gRPC server...")

	go sm.start()

	return sm.waitHealthy(ctx)
}

func (sm *ServerController) Wait() <-chan error {
	return sm.errs
}

func (sm *ServerController) Stop(ctx context.Context) error {
	sm.log.Info("Stopping gRPC server...")

	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// GracefulStop does not support context, so we have to implement basic
	// context pattern.

	stop := make(chan struct{})
	go func() {
		sm.server.GracefulStop()
		close(stop)
	}()

	select {
	case <-ctx.Done():
		sm.server.Stop()
		return ctx.Err()
	case <-stop:
		return <-sm.errs
	}
}

func (sm *ServerController) start() {
	defer close(sm.errs)

	lis, err := net.Listen("tcp", sm.addr)
	if err != nil {
		sm.errs <- err
		return
	}

	if sm.unixAddr != "" {
		lis, err := net.Listen("unix", sm.unixAddr)
		if err != nil {
			sm.errs <- err
			return
		}

		var wg sync.WaitGroup
		wg.Add(1)
		defer wg.Wait()

		go func() {
			defer wg.Done()

			if err := sm.server.Serve(lis); err != nil {
				sm.errs <- err
			}
		}()
	}

	if err := sm.server.Serve(lis); err != nil {
		sm.errs <- err
		return
	}
}

func (sm *ServerController) waitHealthy(ctx context.Context) error {
	conn, err := sm.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	sm.log.Info("Waiting for gRPC server to be healthy...")

	client := grpc_health_v1.NewHealthClient(conn)

	// Do first check immediately
	if checkHealthy(ctx, client) {
		sm.log.Info("gRPC server is healthy")
		return nil
	}

	const interval = 1 * time.Millisecond
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sm.errs:
			return err
		default:
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sm.errs:
			return err
		case <-t.C:
			if checkHealthy(ctx, client) {
				sm.log.Info("gRPC server is healthy")
				return nil
			}

			sm.log.Debug("gRPC server is not healthy")
		}
	}
}

func (sm *ServerController) dial(ctx context.Context) (*grpc.ClientConn, error) {
	type dial struct {
		conn *grpc.ClientConn
		err  error
	}

	dials := make(chan dial)

	go func() {
		const timeout = 10 * time.Second
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		conn, err := grpc.DialContext(
			ctx,
			sm.dialAddr(),
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		dials <- dial{conn: conn, err: err}
	}()

	select {
	case d := <-dials:
		return d.conn, d.err
	case err := <-sm.errs:
		return nil, err
	}
}

func (sm *ServerController) dialAddr() string {
	addr := sm.addr
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}

	return addr
}

func checkHealthy(ctx context.Context, client grpc_health_v1.HealthClient) bool {
	out, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return false
	}

	return out.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING
}
