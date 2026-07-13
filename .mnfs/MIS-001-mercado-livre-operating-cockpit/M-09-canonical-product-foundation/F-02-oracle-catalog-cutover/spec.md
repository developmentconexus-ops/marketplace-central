# F-02 Oracle Catalog Cutover — Specification

Catalog list/get/search must replace the active MetalShopping/MSDB reader with
MPC-owned internal-read facts: positive CODPROD identity, separate EAN/reference
search keys, nullable cost/price/stock facts, and source quality metadata. The
cutover may proceed only with a governed read-only Oracle proof; no fixture or
MSDB fallback is acceptable.

## Blocking Prerequisite

The configured environment must provide the canonical Oracle configuration and
explicit `MPC_ORACLE_LIVE_TEST=1` opt-in so the existing read-only smoke can
prove a positive CODPROD and source timestamps without exposing a row,
credential, or PII.

### F02-AC01

Missing Oracle facts are null with explicit unknown quality rather than zero.

### F02-AC02

Active catalog composition has no MSDB dependency.

### F02-AC03

The governed read-only Oracle lane proves a positive CODPROD read.
