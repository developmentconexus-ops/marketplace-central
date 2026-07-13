# F-03 Catalog Compatibility and Cutover — Validation

## Verdict

Implemented, pending independent review/QA.

## Completed

- Active MSDB catalog composition/configuration and governance registry entries were removed.
- Catalog HTTP reads now require positive integer CODPROD and return the canonical endpoint model.

## Compatibility Outcome

The approved policy is explicit: all historic classification and pricing string
references are preserved as evidence and recorded `not_found` with a null
canonical identity. They are not coerced or attached. The deterministic resolver
only maps a positive decimal legacy value when the supplied canonical candidate
set proves exactly one equal CODPROD; zero candidates are `not_found`, and more
than one are `identity_conflict`.

## Proof

- Both F-02 and F-03 context packs compile successfully at the accepted base.
- Targeted catalog/internal-read/server/unit tests pass, including mapped,
  unmapped, and ambiguous compatibility cases.
- Active MSDB residue scan passes for active server/config/governance paths.
