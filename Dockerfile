ARG GOLANG_VERSION="1.21"
ARG ALPINE_VERSION="3.18"

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION}

WORKDIR /opt/auditumio/auditum

COPY . .

RUN CGO_ENABLED=0 \
  go build \
  -trimpath \
  -mod=vendor \
  -o ./bin/auditum ./cmd/auditum

FROM alpine:${ALPINE_VERSION}

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /opt/auditumio/auditum

COPY --from=0 /opt/auditumio/auditum/bin/auditum /usr/local/bin/auditum
COPY --from=0 /opt/auditumio/auditum/config/auditum.yaml /opt/auditumio/auditum/auditum.yaml
COPY --from=0 /opt/auditumio/auditum/internal/sql/sqlite/migrations /opt/auditumio/auditum/migrations/sqlite
COPY --from=0 /opt/auditumio/auditum/internal/sql/postgres/migrations /opt/auditumio/auditum/migrations/postgres

ENV AUDITUM_store_sqlite_migrationsPath="/opt/auditumio/auditum/migrations/sqlite"
ENV AUDITUM_store_postgres_migrationsPath="/opt/auditumio/auditum/migrations/postgres"

USER nobody

EXPOSE 8080 9090

ENTRYPOINT ["/usr/local/bin/auditum"]
CMD ["serve", "--config", "/opt/auditumio/auditum/auditum.yaml"]
