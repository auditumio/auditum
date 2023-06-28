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
