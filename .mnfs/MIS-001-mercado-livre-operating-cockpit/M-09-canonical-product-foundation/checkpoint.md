# M-09 Terminal Checkpoint

- status: `implementation_complete_pending_qa`
- milestone task: `019f5d00-0b82-7b61-9920-32c7bd490333`
- accepted base: `3087e39fd2dd410763c3a186ffdb75ea8e67a64d`
- current integrated SHA: `3c39f9f01004ea8720048f0772dd3c97b75b91ab`
- completed: harness criterion-format compatibility (`d69490a0aef4ba2c3198e3e39fb20477e811e276`); F-01 canonical product identity (`3c39f9f01004ea8720048f0772dd3c97b75b91ab`).
- completed cutover work: F-02 active catalog reads now use the internal_read/Oracle adapter; the governed live smoke reran successfully and MSDB composition/config was removed.
- compatibility policy: legacy classification/pricing strings are retained only as
  `not_found` evidence; no value is coerced or attached. The resolver accepts
  only exact positive CODPROD equality and exposes ambiguous input as
  `identity_conflict`.
- evidence:
  - `F-01-canonical-product-identity/validation.md`
  - `F-02-oracle-catalog-cutover/validation.md`
  - `F-03-catalog-compatibility-cutover/validation.md`
- next: freeze the implementation SHA and request fixed-SHA review then QA.
