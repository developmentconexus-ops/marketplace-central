# syntax=docker/dockerfile:1
# Production image for server_core: multi-stage, compiled binaries only —
# the runtime stage carries no Go toolchain and no source code.

FROM golang:1.25-bookworm AS build

ENV CGO_ENABLED=1

RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY go.work go.work.sum ./
COPY apps/server_core/go.mod apps/server_core/go.sum ./apps/server_core/
RUN go mod download -C apps/server_core

COPY apps/server_core ./apps/server_core
RUN mkdir -p /out \
    && go build -C apps/server_core -o /out/server ./cmd/server \
    && go build -C apps/server_core -o /out/migrate ./cmd/migrate


FROM debian:bookworm-slim

# godror/ODPI-C loads the Oracle client at runtime; basiclite keeps the image small.
ARG ORACLE_IC_BASE=https://download.oracle.com/otn_software/linux/instantclient/2390000
ARG ORACLE_IC_ZIP=instantclient-basiclite-linux.x64-23.9.0.25.07.zip

ENV MPC_ORACLE_LIB_DIR=/opt/oracle/instantclient \
    LD_LIBRARY_PATH=/opt/oracle/instantclient

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl libaio1 libnsl2 unzip \
    && mkdir -p /opt/oracle \
    && curl -fsSL "${ORACLE_IC_BASE}/${ORACLE_IC_ZIP}" -o /tmp/instantclient.zip \
    && unzip -q /tmp/instantclient.zip -d /opt/oracle \
    && mv /opt/oracle/instantclient_* /opt/oracle/instantclient \
    && rm /tmp/instantclient.zip \
    && apt-get purge -y unzip \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/server /out/migrate /app/
COPY docker/prod/backend-entrypoint.sh /app/entrypoint.sh

RUN useradd --system --uid 10001 --home /app mpc \
    && chmod +x /app/entrypoint.sh

USER mpc
WORKDIR /app
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=5 \
  CMD curl -fsS http://localhost:8080/healthz || exit 1

ENTRYPOINT ["/app/entrypoint.sh"]
