# D5-B2 — Single OpenAPI Wire Authority / Tooling

> **Status:** ACCEPTED / CURRENT CONSOLIDATED AUTHORITY  
> **Canonical entrypoint:** `contracts/api/product/openapi.yaml`  
> **Current Product wire:** OpenAPI 3.1.2 / **106 operations / 31 ordinary Permissions / H-A-S**  
> **Implementation:** blocked until D9

## 1. Governing invariant

> **Marketplace Central has one machine-readable Product API wire authority. Source OpenAPI is authority; resolved bundles and generated TypeScript/Go are derived proof artifacts. Tooling may validate/project the contract but never become a second Product authority.**

W1–W4 own semantic wire laws. This file owns only the machine-description/tooling contract.

## 2. Canonical source

Exactly one Product OAD entrypoint exists:

```text
contracts/api/product/openapi.yaml
```

Repository-local relative `$ref` sources under `contracts/api/product/` are one logical authority. No independently complete per-domain OAD, remote refs, hand-maintained bundle, template/macro/code-first preprocessor or joined parallel API authority exists.

OAS baseline remains exactly:

```yaml
openapi: 3.1.2
jsonSchemaDialect: https://spec.openapis.org/oas/3.1/dialect/base
```

Change only for a concrete accepted wire requirement plus proven tooling support.

## 3. Artifact status

```text
Product source OAD                  AUTHORITY
resolved bundle                     DERIVED / TEMPORARY
generated TypeScript                DERIVED / TEMPORARY
generated Go/server projection      DERIVED / TEMPORARY
runtime handlers                    IMPLEMENTATION conforming to OAD after D9
```

Generated artifacts are not committed as a second wire authority unless a later distribution consumer proves that requirement explicitly.

## 4. Current proof obligations

Product-affecting changes prove proportionately:

- one exact Product OpenAPI entrypoint;
- OAS 3.1.2 and resolved local refs;
- no remote/non-local source refs;
- stable current 106 operation IDs and uniqueness;
- current 31 ordinary Permission vocabulary and H/A/S admission;
- current human-session + machine-bearer authentication profile;
- Organization/privacy and Product Problem invariants;
- idempotency/precondition/route/schema semantics protected by current verifier;
- deterministic resolved bundle;
- deterministic TypeScript generation + strict compilation;
- deterministic Go generation + compilation/tests and required minimum compatibility proof;
- generated colon-suffix custom-operation route expressibility;
- source-tree policy/reachability where current source hygiene requires it;
- no dirty working-tree side effect from proof.

The aggregate required CI remains `npm run gate`; targeted proofs may be invoked separately when their claim is materially in scope.

## 5. Source reachability / orphan policy

OpenAPI source files are authoring structure, not a dumping ground for unreachable historical definitions.

Current accepted source-hygiene rule:

> **A newly introduced source `pathItem` or schema intended as Product authority must be reachable from the canonical entrypoint. A pre-existing explicitly frozen unreachable definition may remain only under the exact allowlist while its historical compatibility/cleanup reason is still current; content drift or new orphan creation fails.**

`source-orphan-allowlist.json` is a bounded mechanical allowlist, not a second Product schema registry or permission to accumulate dead definitions.

`verify-oad-source-reachability.mjs` mechanically proves source reachability and also hosts the current PublicationRequirements parsed-bundle proof. If all allowed orphans are later safely removed, retire the allowlist rather than preserving empty ceremony.

## 6. Product OAD vs Technical Ingress

Technical provider/business-system ingress is not a second Product OAD and does not join the Product SDK merely because it uses HTTP. Technical ingress owns provider protocol semantics under D4 and may have separately proved route/mechanism contracts without entering the 106 Product operation census.

## 7. Authentication/profile encoding

Canonical OAD encodes the accepted Product auth split:

```text
human H → server-side session cookie + CSRF on unsafe methods
A/S     → confidential machine bearer
```

No operation shadows the root profile by convenience unless a future material exception is explicitly admitted. Browser never receives OIDC access/refresh tokens.

## 8. Stable Product Problem origin

Current Product Problem namespace/stable origin remains `https://conexus.fun` where the wire uses product problem URIs. Preview/ngrok/temporary hostnames never become canonical Product source authority.

## 9. Generator/tool pinning

Current proofs use the versions established by repository scripts, including Redocly, `openapi-typescript`, TypeScript and `oapi-codegen`. Version/tool changes require evidence that the generated semantic contract is not weakened or silently altered.

Do not introduce generator-specific `x-go-*`/similar semantic overrides into Product source merely to force one tool's desired shape when the accepted OpenAPI meaning is already expressible portably.

## 10. Verification proportionality

The heavy Product proof protects Product/proof-input changes. Repository/frontend planning changes that cannot alter the Product wire should not execute expensive generation merely as ceremony. Diff-aware selection belongs `scripts/gate.ps1`/repository rules and must fail safe to the stronger proof when reliable diff information is unavailable.

This proportionality changes **when** the proof runs, not what the current Product proof protects.

## 11. Forbidden tooling architectures

Do not introduce:

- second writable OpenAPI/SDK contract;
- hand-maintained generated bundle as authority;
- remote-ref dependency for core Product source;
- code-first schema authority parallel to OpenAPI;
- per-owner independently versioned Product APIs by symmetry;
- provider schema joined into Product OAD;
- custom preprocessing/template language without a concrete current requirement;
- CI checks for prose preferences/decision quality.

## 12. Reopen trigger

Reopen tooling only when a concrete Product wire/generator/runtime consumer cannot be represented/proved by the current single-OAD model, or when current proof cost becomes a material bottleneck that can be reduced with demonstrated assertion-level parity. Tool preference alone is not evidence.
