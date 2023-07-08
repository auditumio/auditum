ARG GOLANG_VERSION="1.20"
ARG ALPINE_VERSION="3.18"

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION}

WORKDIR /opt/infragmo/auditum

COPY . .

RUN CGO_ENABLED=0 \
  go build \
  -trimpath \
  -mod=vendor \
  -o ./bin/auditum ./cmd/auditum

FROM alpine:${ALPINE_VERSION}

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /opt/infragmo/auditum

COPY --from=0 /opt/infragmo/auditum/bin/auditum /usr/local/bin/auditum
COPY --from=0 /opt/infragmo/auditum/config/auditum.yaml /opt/infragmo/auditum/auditum.yaml
COPY --from=0 /opt/infragmo/auditum/internal/sql/sqlite/migrations /opt/infragmo/auditum/migrations/sqlite
COPY --from=0 /opt/infragmo/auditum/internal/sql/postgres/migrations /opt/infragmo/auditum/migrations/postgres

ENV AUDITUM_store_sqlite_migrationsPath="/opt/infragmo/auditum/migrations/sqlite"
ENV AUDITUM_store_postgres_migrationsPath="/opt/infragmo/auditum/migrations/postgres"

USER nobody

EXPOSE 8080 9090

ENTRYPOINT ["/usr/local/bin/auditum"]
CMD ["serve", "--config", "/opt/infragmo/auditum/auditum.yaml"]
