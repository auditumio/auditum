package httpx

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ServerController struct {
	server *http.Server
	log    *zap.Logger
	errs   chan error
}

func NewServerController(server *http.Server, log *zap.Logger) *ServerController {
	return &ServerController{
		server: server,
		log: log.Named("http_server_controller").
			With(zap.String("addr", server.Addr)),
		errs: make(chan error, 1),
	}
}

func (sm *ServerController) Start() {
	sm.log.Info("Starting HTTP server...")

	go sm.start()

	// TODO: wait for server to be healthy.
}

func (sm *ServerController) Wait() <-chan error {
	return sm.errs
}

func (sm *ServerController) Stop(ctx context.Context) error {
	sm.log.Info("Stopping HTTP server...")

	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := sm.server.Shutdown(ctx); err != nil {
		return err
	}

	return <-sm.errs
}

func (sm *ServerController) start() {
	defer close(sm.errs)

	if err := sm.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}

		sm.errs <- err
		return
	}
}
