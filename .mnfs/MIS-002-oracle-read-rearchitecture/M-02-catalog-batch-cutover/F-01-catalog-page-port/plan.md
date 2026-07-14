# F-01 catalog page port — implementation plan

1. Preserve the committed `internal_read/ports` interface and cursor types;
   add only the adapter-facing quality mapping needed by IC-01.
2. Add Oracle page SQL in `internal_read/adapters/oracle` using base `TGFPRO`,
   `TGFEST`, `TGFEXC`/`TGFTAB`, and `TGFCUS` joins. Use a single QueryContext,
   `limit+1` for list pagination, and bounded 1..50 search.
3. Scan rows into adapter-local SQL/null types, map nil facts and flags, and
   wrap all driver failures through the existing redaction helpers.
4. Add fake-queryer tests for one-query page sizes, three-page cursor behavior,
   invalid cursor short-circuiting, nullable facts, and duplicate prices.
5. Expose the page seam through the internal-read application/observability
   wrappers and switch catalog listing projection to it while retaining old
   entity-reader methods for other callers.
6. Run the required build/test ladder, capture verbatim output in
   `validation.md`, inspect the diff, and create the intentional second commit.
