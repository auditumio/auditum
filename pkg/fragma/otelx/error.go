package otelx

import (
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func SetupErrorHandler(log *zap.Logger) {
	otel.SetErrorHandler(errorHandler{log: log.Named("otel")})
}

type errorHandler struct {
	log *zap.Logger
}

func (h errorHandler) Handle(err error) {
	h.log.Warn("OpenTelemetry error", zap.Error(err))
}
