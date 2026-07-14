# ADR-007: Godror with OCI is the canonical Oracle runtime

**Date:** 2026-07-13
**Status:** accepted

## Context

MPC needs predictable read-only access to the Sankhya Oracle database from its
Go modular monolith. The first live lane proved the query contract, but its
execution wrapper could deadlock while building the Docker image and the
adapter left important pool and timeout behavior implicit.

Pure-Go `go-ora` and a separate Python `python-oracledb` service were evaluated.
Pure Go would reduce container weight, but its current release has unresolved
context-cancellation, RAC, and LOB risks that are material to a long-lived
operator service. Python would introduce a second runtime and service boundary
without a distinct ownership or scaling need.

## Decision

MPC keeps `godror` backed by ODPI-C/Oracle Instant Client as its single
canonical Oracle driver.

The runtime contract is:

- `database/sql` is opened through typed `godror.ConnectionParams`, never a
  hand-built or logged credential-bearing DSN;
- Oracle and Go pool bounds, wait timeout, idle timeout, maximum session
  lifetime, connect timeout, bootstrap timeout, and per-call timeout are
  explicit and validated;
- connection health is proven at bootstrap; ordinary reads do not issue an
  extra `Ping` before every query;
- adapter errors expose a stable read-error class and, when available, only an
  Oracle numeric code rather than raw driver text;
- the canonical live lane runs inside the same OCI-capable backend image, with
  a read-only checkout/root filesystem, dropped capabilities, bounded phases,
  concurrent output draining, and forced process-tree termination on timeout;
- Oracle validation issues no durable database write and requires an
  explicitly governed live-test opt-in. The governance lane classifies its
  side-effect envelope as `database-read` plus `session-temporary-write`: the
  default C05 lane issues SELECTs only, while the explicit `-EmitBaseline`
  mode additionally runs `EXPLAIN PLAN`, a session-private `PLAN_TABLE`
  global-temporary insert that leaves no durable database state. Each mode
  has its own test selector in `docker/live-oracle/profile.json`; the default
  lane never selects the baseline test.

## Consequences

- The backend image continues to carry CGO and Oracle Instant Client. That is a
  deliberate reliability tradeoff, not an accidental dependency.
- Docker build layers and Go modules must be cached; the live runner may never
  wait indefinitely on Docker or Oracle.
- A future driver change requires real compatibility evidence for every
  governed query shape, cancellation, LOB behavior, and the target Oracle
  topology. Container size alone is not sufficient reason to switch.
- No Python Oracle sidecar or dual-driver fallback is part of the target
  architecture.

