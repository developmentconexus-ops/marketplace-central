# F-03 Catalog Compatibility and Cutover — Specification

Legacy MPC product references may migrate only when their stored value is a
positive decimal CODPROD. The forward migration records any non-numeric or
non-positive value as `not_found`; it never guesses from EAN, reference, title,
or seller SKU. Runtime readers use the canonical integer reference only.

Multiple upstream candidates remain `identity_conflict`; neither conflict nor
not-found attaches an enrichment, classification, pricing record, or product
link to a product.

### F03-AC01

Active runtime has no MSDB composition or configuration reader.

### F03-AC02

Compatibility resolves only exact positive CODPROD equality; mapped, not_found,
and identity_conflict states are deterministic.
