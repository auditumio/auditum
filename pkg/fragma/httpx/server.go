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
	"net/http"
	"time"

	"go.uber.org/zap"
)

func NewServer(addr string, handler http.Handler, log *zap.Logger) *http.Server {
	// Can return error only on unknown level.
	// As we provide only known level, it's safe to panic.
	elog, err := zap.NewStdLogAt(log.Named("http_server"), zap.WarnLevel)
	if err != nil {
		panic(err)
	}

	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 60 * time.Second,
		ErrorLog:          elog,
	}
}
