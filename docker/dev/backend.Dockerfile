# syntax=docker/dockerfile:1
FROM golang:1.25-bookworm

ARG ORACLE_IC_BASE=https://download.oracle.com/otn_software/linux/instantclient/2390000
ARG ORACLE_IC_ZIP=instantclient-basiclite-linux.x64-23.9.0.25.07.zip

ENV CGO_ENABLED=1 \
    ORACLE_HOME=/opt/oracle/instantclient \
    LD_LIBRARY_PATH=/opt/oracle/instantclient \
    PATH=/opt/oracle/instantclient:$PATH

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
      build-essential \
      ca-certificates \
      curl \
      libaio1 \
      libnsl2 \
      pkg-config \
      postgresql-client \
      unzip \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /opt/oracle \
    && curl -fsSL "${ORACLE_IC_BASE}/${ORACLE_IC_ZIP}" -o /tmp/instantclient.zip \
    && unzip -q /tmp/instantclient.zip -d /opt/oracle \
    && mv /opt/oracle/instantclient_* /opt/oracle/instantclient \
    && rm /tmp/instantclient.zip \
    && ldconfig

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY apps/server_core/go.mod apps/server_core/go.sum ./apps/server_core/
RUN go mod download -C apps/server_core

COPY apps/server_core ./apps/server_core
RUN mkdir -p /opt/mpc/bin \
    && go test -c -o /opt/mpc/bin/oracle-live.test ./apps/server_core/internal/modules/internal_read/adapters/oracle

EXPOSE 8080
