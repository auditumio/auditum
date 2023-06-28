#!/usr/bin/env bash
set \
  -o errexit \
  -o nounset \
  -o pipefail

log::info() {
  echo -e "$(date -u +"%Y-%m-%dT%H:%M:%SZ") \e[34m[INFO]\e[0m $*" >&2
}

log::success() {
  echo -e "$(date -u +"%Y-%m-%dT%H:%M:%SZ") \e[32m[SUCC]\e[0m $*" >&2
}

log::error() {
  echo -e "$(date -u +"%Y-%m-%dT%H:%M:%SZ") \e[31m[ERRO]\e[0m $*" >&2
}

export POSTGRES_HOST="127.0.0.1"
export POSTGRES_PORT="5432"
export POSTGRES_DBNAME="${INTEGRATION_TESTS_POSTGRES_DATABASE:-auditum_db}"
export POSTGRES_USERNAME="${INTEGRATION_TESTS_POSTGRES_USERNAME:-user}"
export POSTGRES_PASSWORD="${INTEGRATION_TESTS_POSTGRES_PASSWORD:-pass}"
export POSTGRES_SSL_MODE="disable"
export POSTGRES_LOG_QUERIES="${INTEGRATION_TESTS_POSTGRES_LOG_QUERIES:-}"

log::info "--> Start docker-compose ..."

docker-compose \
  --file test/docker-compose.yml \
  --project-directory . \
  up \
  --build \
  --force-recreate \
  --renew-anon-volumes \
  --detach

log::info "--> Run tests ..."

set +e
task test-integration -- ${*}
result=$?
set -e

log::info "--> Stop docker-compose ..."

docker-compose \
  --file test/docker-compose.yml \
  --project-directory . \
  down \
  --volumes

if [[ ${result} -eq 0 ]]; then
  log::success "Tests passed!"
  exit 0
else
  log::error "Tests failed!"
  exit 1
fi
