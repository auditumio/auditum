// Copyright 2023 Igor Zibarev
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
