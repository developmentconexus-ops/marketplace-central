# F-08 Sankhya Administrator Specification

## Deployment Contract

The administrator supplies the actual field name. MPC must never ship a
customer column name as a default.

| Setting | Required contract |
| --- | --- |
| `header_field_name` | Existing supported custom column on `TGFCAB`; exact identifier allowlisted from Oracle metadata, never interpolated from a request. |
| Type | Text equivalent to `VARCHAR2(160 CHAR)` or a supported larger text field; trim is forbidden, comparison is exact. |
| Value | `ml:v1:<installation-or-account-id>:<provider-order-id>`; immutable after confirmation. No buyer/customer data. |
| Applicability | TOP 313 only; blank is allowed for unrelated/legacy documents, but confirmation requires a nonblank exact key. |
| Uniqueness | Unique among nonblank values across the deployed MPC/account scope. The key itself contains installation/account identity. |
| MPC access | Service principal has metadata and SELECT access only. This phase grants MPC no Oracle INSERT/UPDATE/DDL permission. |
| Human access | Named Sankhya role may set the field through the sanctioned order-entry/admin surface before confirmation; ordinary users cannot alter it after confirmation. |

If the supported native extension cannot enforce uniqueness, the administrator
must deploy an approved normalized integration table/constraint through the
customer's sanctioned Sankhya customization process. MPC stays disabled until
that mechanism is named and a duplicate probe passes. Read-before-confirm by
itself is not sufficient concurrency protection.

## Line-Field Decision

No Sankhya item custom field is required for the assisted-only release. MPC
persists an immutable `mpc_line_id` and the explicitly confirmed TOP 313
(`NUNOTA`,`SEQUENCIA`) in its append-only ledger. This is sufficient because:

- the operator confirms exact lines rather than an inferred product match;
- duplicate identical items remain distinguishable by immutable MPC identity;
- partial invoices retain the same 313 origin and append every `TGFVAR` 306
  descendant with its `QTDATENDIDA`;
- one-to-many descendants are not collapsed; and
- custom-field copy to 306 is unproved and unnecessary for canonical lineage.

An item field may be introduced later only as a separately approved recovery/
operational aid. It must follow the same configured-name, immutable-key,
uniqueness, permission, and validation rules and remains secondary to the MPC
ledger plus `TGFVAR`.

## Pre-Enable Validation

An administrator and MPC operator must record values-free results for all
checks:

1. Field metadata exists on `TGFCAB`, has compatible text capacity, and its
   normalized identifier matches the configured identifier exactly.
2. The authorized entry surface exposes the field only for TOP 313 and does
   not truncate or transform the key.
3. A bounded aggregate SELECT reports no duplicate nonblank values and no
   overlength values. It outputs counts only, never key values or document rows.
4. MPC's Oracle principal can SELECT metadata/header/lineage but cannot mutate
   Oracle objects or rows.
5. Representative bounded tests prove effective TOP 313 headers and TOP 306
   descendants through exact `TGFVAR` origin predicates, including a partial
   or one-to-many case when such a safe fixture exists.
6. The application configuration validates at startup. Missing/invalid field,
   uniqueness attestation, or runtime probe keeps candidate viewing read-only
   and disables confirmation.

This Feature did not run those checks; their current results are `unknown`.

## Permissions And Operating Procedure

- Sankhya administrator: deploy/approve field and uniqueness mechanism; owns
  rollback and supported customization evidence.
- Order operator: create/select TOP 313, set the exact external key in the
  sanctioned UI, then explicitly confirm header and line mapping in MPC with a
  reason.
- MPC runtime: bounded SELECT and Postgres ledger append only. It never writes
  the Sankhya field in this release.
- Reviewer/auditor: reads mapping history, conflicts, actor/reason/source time,
  and lineage states; cannot edit history.

Configuration is secret-free: field identifier, feature-enabled flag,
uniqueness-mechanism ID/attestation revision, and validation timestamp. Do not
store Oracle credentials or provider IDs in logs.

## Disable And Rollback

Disabling the feature immediately blocks new confirmations and Oracle lineage
refreshes while preserving read-only history. Rollback does not delete MPC
ledger events, clear header values, reuse external keys, or drop fields/
constraints. The administrator may remove UI exposure or schema objects only
through a separate approved change after MPC is disabled and retention/audit
owners approve. Existing profitability remains missing/unknown wherever exact
validated lineage is unavailable.
