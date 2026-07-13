#!/usr/bin/env bash
set -euo pipefail

strip_quotes() {
  local value="${1:-}"
  value="${value%$'\r'}"
  if [[ ${#value} -ge 2 ]]; then
    if [[ "${value:0:1}" == "'" && "${value: -1}" == "'" ]]; then
      printf "%s" "${value:1:${#value}-2}"
      return
    fi
    if [[ "${value:0:1}" == '"' && "${value: -1}" == '"' ]]; then
      printf "%s" "${value:1:${#value}-2}"
      return
    fi
  fi
  printf "%s" "$value"
}

export API_PORT="${API_PORT:-8080}"
export SERVER_ADDR="${SERVER_ADDR:-:${API_PORT}}"
export GOCACHE="${GOCACHE:-/root/.cache/go-build}"
export MPC_ORACLE_LIB_DIR="${MPC_ORACLE_LIB_DIR:-/opt/oracle/instantclient}"

if [[ -z "${MPC_ORACLE_USERNAME:-}" && -n "${SANKHYA_ORACLE_USER:-}" ]]; then
  export MPC_ORACLE_USERNAME="$(strip_quotes "$SANKHYA_ORACLE_USER")"
fi
if [[ -z "${MPC_ORACLE_PASSWORD:-}" && -n "${SANKHYA_ORACLE_PASSWORD:-}" ]]; then
  export MPC_ORACLE_PASSWORD="$(strip_quotes "$SANKHYA_ORACLE_PASSWORD")"
else
  export MPC_ORACLE_PASSWORD="$(strip_quotes "${MPC_ORACLE_PASSWORD:-}")"
fi
if [[ -z "${MPC_ORACLE_CONNECT_STRING:-}" && -n "${SANKHYA_ORACLE_HOST:-}" && -n "${SANKHYA_ORACLE_PORT:-}" && -n "${SANKHYA_ORACLE_SERVICE_NAME:-}" ]]; then
  export MPC_ORACLE_CONNECT_STRING="$(strip_quotes "$SANKHYA_ORACLE_HOST"):$(strip_quotes "$SANKHYA_ORACLE_PORT")/$(strip_quotes "$SANKHYA_ORACLE_SERVICE_NAME")"
fi


if [[ "${RUN_MIGRATIONS:-1}" == "1" ]]; then
  go run ./apps/server_core/cmd/migrate
fi

case "${1:-server}" in
  server)
    exec go run ./apps/server_core/cmd/server
    ;;
  test)
    shift || true
    exec go test "$@"
    ;;
  *)
    exec "$@"
    ;;
esac
