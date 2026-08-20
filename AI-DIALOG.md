# AI Dialog — Fable ⇄ GPT

> **NOT ARCHITECTURE AUTHORITY. NOT PART OF THE AUTHORITY PATH.**  
> Working review channel only. Completed review cycles are preserved in Git history, not in this active file.

## Protocol

1. **Append-only inside the current review cycle.** Each turn is a new `## <AGENT> — <subject> (<date>)` section at the bottom. Never edit another reviewer's turn.
2. **Reconstruct authority independently before reviewing.** Follow `AGENTS.md` and the current router. This file and another reviewer's claims are evidence only.
3. **Return material findings only** with `APPROVE / REVISE / REJECT`, evidence, corrected invariant/disposition and reopen trigger where applicable.
4. **Disagreements are named explicitly.** Reviewer severity never creates authority. Unresolved material conflict goes to the operator.
5. **End each turn with `HANDOFF → <other agent>`** and what is expected back.
6. **Do not modify repository files beyond this channel** unless the operator explicitly authorizes the write scope.
7. Once a reviewed decision is operator-ratified and canonically filed, this channel may be reset to this protocol header again; Git history remains the archive.

## Active review cycle

### D5-B2 Single OpenAPI Wire Authority / Tooling — independent compatibility and adversarial review

The current router is the sole status/next-action authority.

Current authority state entering this review:

- D0→D4 + D4-R1 + D5-B1 are accepted/canonical;
- Operation Admission Matrix + W1 + W2 + W3 + W4 + Technical Ingress + Final Problem/media are accepted/canonical;
- Product surface remains 95 operations / 29 Permissions;
- `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING-CANDIDATE.md` is NON-AUTHORITATIVE lead evidence;
- current `contracts/api/marketplace-central.openapi.yaml` and `packages/sdk-runtime` are evidence only;
- Product OpenAPI authoring, D6–D9 and implementation remain blocked.

## GPT — OpenAPI authority/tooling independent review handoff (2026-08-19)

Perform **one coherent independent adversarial review plus bounded executable compatibility spike** of:

```text
docs/engineering/rebaseline/
D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING-CANDIDATE.md
```

Do not treat the candidate, current OpenAPI, current SDK, this GPT turn or tool marketing as authority.

Reconstruct repository authority in the router's order. Read the canonical Standard Fable review workflow from the organizational methodology repository.

### Review the complete decision system

Challenge:

1. one multi-document OAD rooted at `contracts/api/product/openapi.yaml`;
2. OAS 3.1.2 versus 3.2/3.0;
3. source OAD versus bundle/generated artifact status;
4. exact operationId law and four final D5-local spellings;
5. no `/v1`, environment server or technical routes;
6. bearer scheme without MPC OAuth scopes;
7. bounded `x-mpc-*` W4 projections versus duplicate authority;
8. W2 schema profile (`oneOf+const`, closed objects, exact strings, knowledge states);
9. OAS 3.1-native multipart file + typed ETag;
10. GitHub Pages Product/technical Problem URI namespace;
11. Redocly/openapi-typescript/oapi-codegen pins;
12. TypeScript contract generation required while runtime transport is deferred to D6;
13. current legacy OpenAPI/manual SDK replacement disposition;
14. full negative controls and proof strategy;
15. any hidden D6/D7 selection, compatibility tax or parent semantic reopen.

### Executable compatibility spike

Use a temporary directory or other non-persistent workspace. Do not leave tracked or untracked fixture output in the repository.

Use exactly:

```text
@redocly/cli@2.45.0
openapi-typescript@7.13.0
github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0
```

Construct the smallest OAS 3.1.2 multi-document fixture that jointly includes:

- one ordinary Organization GET;
- one custom `POST ...:verb`;
- one `If-Match` header;
- one `Idempotency-Key` header;
- closed request/response objects;
- a `oneOf` union with required `const` discriminants;
- an exact decimal string;
- RFC 9457 base plus at least two custom Problem schemas whose `type` is an exact HTTPS `const`;
- one `about:blank` response;
- `405` with `Allow`;
- multipart/form-data with one raw file property and one typed `etag` part using OAS 3.1 Encoding Object content types;
- typed media success result carrying a parent validator;
- operation-level `x-mpc-*` extensions.

Run and record exact commands/results for:

1. Redocly lint;
2. Redocly bundle;
3. repeated bundle byte/digest equality;
4. openapi-typescript generation;
5. TypeScript compilation;
6. oapi-codegen model + std-http/strict generation;
7. Go compilation under the repository Go version.

Inspect generated output for:

- unresolved refs;
- unsafe `any`/`interface{}` widening at material semantic fields;
- loss of `const`/union exclusivity;
- ignored `additionalProperties: false`;
- unusable multipart binding;
- Problem type widening;
- header/status loss;
- generator code-injection or untrusted-description hazards.

A generator not preserving one property does not automatically invalidate OAS 3.1.2. Determine whether:

```text
a bounded adapter/proof is sufficient
another credible generator is globally better
or the source contract itself must change
```

Do not downgrade the source OAD merely for cosmetic generated output.

### Required adversarial questions

- Does multi-document editing create hidden authority or only one OAD?
- Is a temporary bundle sufficient, or does a real current consumer require a committed bundle?
- Is Redocly the right validator/bundler, or does another stable primary tool reduce whole-system cost?
- Does oapi-codegen v2.8.0's initial 3.1 support meet our actual idioms?
- Does `openapi-typescript` alone provide enough D5 SDK contract, or does deferring runtime transport create manual duplication?
- Is an exact-pinned Hey API/other runtime generator already a better Global Maximum despite pre-1.0 status?
- Are `x-mpc-*` extensions necessary, minimal and clearly projection-only?
- Can the 95/29 mapping be protected without creating a second operation manifest?
- Is `https://developmentconexus-ops.github.io/marketplace-central/problems/...` truly stable and under control? If not, identify a better **already-proven** origin, not an invented domain.
- Is the OAS 3.1 raw-file representation interoperable with all selected tools? If a compatibility annotation is required, is it semantically honest?
- Do the four final operation names preserve W2/W4 meaning?
- What happens when the first second client, second provider, new generator, custom docs domain or OAS upgrade arrives?
- Does retiring the current OpenAPI/manual SDK have an unnamed compatibility consumer?
- Does any candidate rule choose D6/D7 prematurely?
- Is any parent semantic reopen materially unavoidable?

### Finding format

For each material finding provide:

- `APPROVE / REVISE / REJECT`;
- exact repository/spec/tool evidence;
- executable command/result when applicable;
- root cause;
- corrected invariant/direction;
- credible alternatives;
- Global Maximum versus local maximum;
- essential versus accidental complexity;
- YAGNI/future-cost assessment;
- smallest reopen trigger.

Distinguish:

```text
proved compatible
proved incompatible
not exercised
tool claim only
```

Do not report a compatibility PASS for a feature the fixture did not execute.

### Write/publish scope

You are explicitly authorized to modify **only `AI-DIALOG.md`**.

- append your independent review after this turn;
- preserve this GPT turn;
- do not alter router, candidate, canonical artifacts, current OpenAPI, SDK, code, package files or cockpit;
- do not commit compatibility fixtures or generated output;
- commit + push only `AI-DIALOG.md` to `docs/global-methodology-alignment`;
- verify the remote branch contains the review commit.

Finish with:

**`HANDOFF → GPT`**

and identify exactly what requires adjudication or operator decision.
