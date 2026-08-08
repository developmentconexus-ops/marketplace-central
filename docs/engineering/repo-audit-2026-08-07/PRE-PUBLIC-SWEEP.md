# PRE-PUBLIC SWEEP — marketplace-central

**Date:** 2026-08-07
**Target repo:** `origin` = `github.com/developmentconexus-ops/marketplace-central` (currently private)
**Also configured:** `legacy` = `github.com/leandrotcawork/marketplace-central` (separate repo — **not** affected by flipping `origin`, but see §5)
**Scope swept:** all 2109 commits reachable from every local branch, every remote-tracking ref (`origin/*`, `legacy/*`) and every tag; 24 761 reachable objects; 4326 distinct historical paths; 3122 files at HEAD.

> **Redaction discipline.** No secret or personal-data value appears anywhere in this document. Findings are identified by path, line, commit SHA, and a class label. Where a match had to be proved real rather than a placeholder, it is described by structural properties only (length, character classes, token shape).

---

## 1. VERDICT

# `SAFE AFTER LISTED REMEDIATIONS`

Two blocking findings, both the **same single credential** (a PostgreSQL role password), reachable at HEAD **and** already pushed to `origin/main`. Everything else is clean or non-blocking.

The good news, stated plainly so it is not lost: **no `.env`, `.pem`, `.key`, `.p12`, `.pfx`, `.jks`, `id_rsa*`, `.ovpn`, `.tfstate`, `.npmrc`, `.netrc`, `.pgpass`, or `.aws` file has EVER been tracked in this repository, on any branch, in any commit.** No AWS key, no GitHub PAT, no OpenAI/Anthropic key, no Slack token, no Google API key, no JWT, no PEM private key, and no Mercado Livre OAuth `client_secret`, access token or refresh token exists in the reachable history. No real Oracle/Sankhya host, service name, username or password exists anywhere. No CPF or CNPJ that survives check-digit validation as real data exists anywhere.

---

## 2. BLOCKING FINDINGS

### B-1 — PostgreSQL role password, inline in a tracked plan document (CURRENT TREE + HISTORY)

| | |
|---|---|
| **Path** | `docs/superpowers/plans/2026-04-06-pricing-simulator-v2.md` |
| **Lines** | **77** and **88** |
| **Class** | PostgreSQL role password, assigned to `PGPASSWORD=` in a runnable `psql` command |
| **Introduced** | commit `ab8dd49dff294c83686a0fa81f6883288bd30e6b` (2026-04-06, subject `Commit`) |
| **Present at** | local `HEAD`, local `main`, **and `origin/main`** — confirmed via `git cat-file -e origin/main:<path>` |
| **Tree or history** | **Both.** Still live in the working tree. Deleting the file does not remove it from history. |

**Why this is not a placeholder.** The value is a 23-character all-lowercase string of the shape `word_word_word` (three word segments joined by separators, no uppercase, no digits). It contains none of the placeholder markers (`change*`, `example`, `your`, `seu`/`sua`, `xxx`, `placeholder`, `todo`, `replace`, `dummy`, `fake`, `sample`, `openssl rand`, `${...}`, `<...>`). It is **not** self-referential (unlike the safe `postgres:postgres` / `marketplace:marketplace` dev defaults found elsewhere). It is **not** the hermetic-lane generated password — `scripts/harness/Postgres.psm1:121-123` generates a 48-character hex string from 24 random bytes, which this is not. It is paired with a **17-character named application DB role** (not `postgres`, not `marketplace`) against `127.0.0.1:5432` — the real local cluster, not an ephemeral test container port. Both lines are inside ```bash fences presented as commands the operator is meant to run.

It is the **only distinct `PGPASSWORD` value in the entire reachable history** (verified: 1 blob, 1 distinct value across all 24 761 objects).

### B-2 — The same password, inside a `postgres://` DSN, in `.claude/settings.local.json` (HISTORY ONLY)

| | |
|---|---|
| **Path** | `.claude/settings.local.json` |
| **Line** | 47 and 49 of the historical blob (inside `permissions.allow[]`, as `Bash(psql 'postgres://<role>:<password>@127.0.0.1:<port>/...')`) |
| **Class** | Postgres DSN with inline password — **same credential as B-1**, byte-identical (proved by exact-substring match, value never printed) |
| **Blobs** | `37319a7dca07455b56b0301018a6e5b6b297f694`, `30103464071fd22de8ceb21acf1a5a8554e97a15` |
| **Commits** | added at `4743f778e3714ab8b729735b61d0088d3f5c8078` (2026-04-04); also in `b9cd863f6000c3f982cb47a7c9ebbfdcf7b6597e`; **removed** at `2fe477085035c02978f5e24276c29a69ed96baae` ("chore: untrack .claude/settings.local.json") |
| **Tree or history** | **History only.** Not in the HEAD tree, not in `origin/main`'s HEAD tree. |
| **Already pushed?** | **Yes.** `git merge-base --is-ancestor 4743f778 origin/main` → true. All three commits are ancestors of `origin/main`. |

This is the textbook case the mission brief named: a file added before the ignore rule existed, later untracked, and permanently in history.

### Remediation for B-1 + B-2

> **Status 2026-08-07 — partial, and the partial part is the cheap part.**
>
> **B-1's working-tree copy is scrubbed** at commit `1ec2d081`. Both `PGPASSWORD=` assignments in
> `docs/superpowers/plans/2026-04-06-pricing-simulator-v2.md` now read the value from the
> environment, so the commands stay runnable and the credential is out of the live tree. The value
> was never read, printed, or placed in a transcript — the edit was a blind `sed` over
> `PGPASSWORD=[^[:space:]]*`, verified by counting the redacted form (2) and residual matches (0).
>
> **This does not close either finding.** B-2 is the same credential in a historical blob of
> `.claude/settings.local.json`, and all three of its commits are already ancestors of
> `origin/main`. Nothing done to the tree reaches a pushed copy.
>
> **Operator decision: not rotating now, and it is not needed before P4.** The push is not the
> disclosure event — `origin` is private, so the flip is. Rotation gates **P6** only.
>
> **The one fact that decides whether rotation is ever required, and only the operator has it:**
> is that PostgreSQL role valid anywhere other than the local machine? Local-only puts it with N-1
> and N-2 in the accepted column, and this sweep over-rated it. Valid on any reachable host —
> staging, a prod mirror, anything not `127.0.0.1` — and the flip publishes a working credential.
> **Answer before P6.**

> ### B-1 + B-2 CLOSED — 2026-08-07, by operator answer
>
> **Operator states the role is valid only on their local `127.0.0.1` cluster.** That is the fact
> this sweep said it did not have and would not guess at, and it settles both findings: a published
> string that authenticates against nothing reachable is not a disclosed credential. B-1 and B-2
> **move to the accepted column alongside N-1 and N-2**, and this sweep over-rated them — correctly,
> given what it could see, because credential shape is not credential reach and it declined to
> assume the difference.
>
> **No rotation. No history rewrite.** Rewriting would invalidate every SHA from 2026-04-04 onward
> for a dead string, which is the worse trade even when it is free.
>
> Two things still hold, and neither blocks P6:
>
> - B-1's working-tree copy stays scrubbed (`1ec2d081`). Not because the value is dangerous, but so
>   it is not re-copied into a new file by someone reading the plan later.
> - Remediation step 4 — the pre-commit guard against `PGPASSWORD=`, against `postgres://*:*@` with
>   a non-`${...}` password, and against staging `.claude/settings.local.json` — **is still worth
>   building, and is now the only part of this remediation that is.** It is level 5 as a hook, so it
>   belongs in the gate as a check, not in `.githooks/`. Carried to issue #2's design as a candidate;
>   it costs one grep and it closes the class rather than this instance.
>
> **P6 is unblocked.**

Both findings are **one credential**. One rotation closes both.

1. **ROTATE THE PASSWORD.** Change the password of that PostgreSQL role on every cluster where it is valid (local dev, any staging, any prod mirror). Update whatever consumes it — `.env` (untracked), Docker Compose overrides, the operator's own shell profile.

   > **Rotation is the reliable remediation, not history rewrite.** Once a value has been pushed, it must be treated as disclosed. A `git filter-repo` rewrite does not reach: existing clones on other machines, forks, GitHub's own dangling-object cache (blobs stay fetchable by SHA at `/<owner>/<repo>/blob/<sha>` long after the rewrite), PR/issue attachments, CI caches, or anyone's local reflog. Rewriting history without rotating leaves the credential live and merely harder to find. Rotating without rewriting leaves a dead string in history, which is harmless.

2. **Scrub the current tree** so the value is not the first thing a reader sees: edit `docs/superpowers/plans/2026-04-06-pricing-simulator-v2.md:77` and `:88` to use `PGPASSWORD="$PGPASSWORD"` or an explicit `<set-in-your-shell>` marker. This is cosmetic once step 1 is done, but it stops the same value being re-copied into a new file.

3. **Optional, only after step 1:** if the operator wants the string gone from history as well, `git filter-repo --replace-text` over the three commits plus `ab8dd49d`. This rewrites every SHA from 2026-04-04 onward, invalidates every open branch and worktree, and requires a force-push. **Do not do this instead of rotating. Do not do it while another session holds this checkout.** Given rotation is already sufficient, the recommendation is to skip it.

4. Add a pre-commit guard so this class does not recur: block `PGPASSWORD=`, `postgres://*:*@` with a non-`${...}` password, and `.claude/settings.local.json` from ever being staged.

---

## 3. NON-BLOCKING FINDINGS

| # | Path / location | Class | Assessment |
|---|---|---|---|
| N-1 | `docker-compose.yml:5-7, 27` | Postgres dev password, literal | **Self-referential dev default** (user == password, and the DB name is the project name). Local-only service on the compose network. Standard practice; conventionally accepted as public. No action. |
| N-2 | `apps/server_core/tests/unit/postgres_config_test.go:10, 29` | Postgres DSN, literal | `postgres:postgres@localhost` — canonical test default. No action. |
| N-3 | `.superpowers/sdd/m06-f01-correction-report.md:22`, `m06-f02-correction-report.md` | Postgres DSN with inline password | 3-char role, 8-char two-word hyphenated password, `localhost:5436` (hermetic container port). Almost certainly a throwaway lane credential. **Marked for operator eyeball** — if that string is reused anywhere real, rotate it too. |
| N-4 | `docs/superpowers/specs/2026-04-02-product-integration-rework-design.md` | Postgres DSN, literal | 4-char role / 4-char password (`user:pass` shape). Placeholder. No action. |
| N-5 | `.agents/skills/mastra/references/common-errors.md` | Postgres DSN, literal | Vendored third-party documentation example. No action. |
| N-6 | `apps/server_core/internal/composition/market_adapters_test.go:236, 238, 241` | Mercado Livre `APP_USR-` token shape, bearer header, `TG-`/`refresh-` strings | **Synthetic fixture.** The `APP_USR-` string is 38 chars with 3 separators; a real ML access token is ~74 chars with 4 segments (app id 16 digits, date 6, hex 32, user id). The variable is literally named `secrets` and the test asserts these are scrubbed from logs — i.e. this is the *redaction test*. No action. |
| N-7 | Commit metadata, all 2109 commits | Personal email address | **Exactly 1 distinct author and 1 distinct committer identity across the entire history**, a personal `gmail.com` address (not a role address). It will be world-readable on every commit. The same address also appears in 7 blobs (agent logs, a commit-stat file, a validation doc). Nothing to fix technically — this is the operator's own identity and their call. If undesired, it requires the same history rewrite caveat as B-2. |
| N-8 | `.claude/settings.local.json` (history), `additionalDirectories` | Local filesystem path + machine hostname | Exposes the operator's Windows user directory and laptop name. Same identity already public via N-7. No action. |
| N-9 | `templates/~$catalogo-exemplo.xlsx` (history only, deleted) | Excel lock file | Contains the Excel-registered user name. Same identity as N-7. No action. |
| N-10 | `docker/dev/env.container.example:9`, `docker/dev/README.md`, `docs/superpowers/plans/2026-07-08-*.md`, `.mnfs/MIS-008-.../TASK5-LIVE-DRIVE.md` | ngrok tunnel hostname (`*.ngrok-free.dev`) | A specific free-tier ngrok subdomain used as the OAuth redirect host during dev. Free-tier subdomains are reassigned per session, so this is almost certainly dead. **Not a credential** — `NGROK_AUTHTOKEN` is `${NGROK_AUTHTOKEN:-}` everywhere (verified `docker-compose.yml:78`, `env.container.example:5` empty). Worth deleting for tidiness only. |
| N-11 | ~6 files with populated `nickname` / `receiver_name` / `last_name` / `street_name` / `zip_code` / `city_name` JSON values | Possible real ML buyer data | **See §3.1 — operator eyeball required.** |
| N-12 | Repository root | No `LICENSE`, no `SECURITY.md` | Not a leak, but a public repo with no licence is "all rights reserved" by default, which is usually not the intent. Add both before or shortly after the flip. |
| N-13 | `.git/objects` | 48 garbage `tmp_obj_*` files; ~25 large unreachable blobs (up to 9.2 MB) including Go/x-crypto `testdata` private keys and a 3.9 MB locale word list | **These are unreachable from every ref and every reflog** (verified) — they were never committed, therefore never pushed, therefore **cannot** be exposed by the flip. They are local repo litter from an aborted `git add`. `git gc --prune=now` would clear them; that is the operator's call and is not required for the flip. Listed only so nobody re-discovers them and panics. |

### 3.1 — PII items requiring the operator's own eyes (N-11)

The sweep found **zero** formatted CPFs (`000.000.000-00`) and **zero** formatted CNPJs (`00.000.000/0000-00`) anywhere in the reachable history. Every 11-digit run adjacent to a `cpf`/`doc_number` label was check-digit validated: every one either **failed the check digit** (synthetic) or **passed but is the canonical ascending test CPF** whose first nine digits are `1..9` — the textbook fake. Conclusion: **no real CPF is committed.**

What remains is a small, bounded set of populated name/address fields in ML API fixtures. Structural analysis says "plausibly synthetic" but cannot prove it, and per the mission brief I am not guessing. There are **9 distinct values** across these locations:

- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/buyer_fiscal_reader_test.go` — `nickname` ×2, `last_name`, `street_name`, `city_name`, `zip_code`
- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipment_ingest_reader_test.go` — `receiver_name`, `zip_code`
- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/order_ingest_reader_test.go`, `capability_adapter_test.go`, `resilience_decorator_test.go` — `nickname`
- `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter_test.go` — `nickname`
- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/testdata/order_body.json` — `nickname`
- `.mnfs/MIS-004-mvp-demo/M-08-pedidos/_chip-m08-buyer/fix.diff` — same values as `buyer_fiscal_reader_test.go`
- `.mnfs/MIS-004-mvp-demo/M-08-pedidos/_chip-m08-shipfix/fix-round4.diff` and `.mnfs/MIS-007-ml-sync/planning-reviews/p7-seat4-star7-r01.md` — a shared `receiver_name` + `street_name` + `zip_code` triple
- `docs/superpowers/plans/2026-07-08-marketplace-integration-foundation.md`, `docs/superpowers/plans/2026-08-02-p2-dinheiro-real-pedidos.md` — `nickname`

**Evidence that they are synthetic:** every `nickname` is a pure-uppercase-letter token with **no digits** — real ML nicknames are letters+digits. Values are reused verbatim across 2–4 unrelated files, which is fixture behaviour, not capture behaviour. One `receiver_name` carries an explicit placeholder word. Volume is 9 values, not 38 orders' worth.

**Evidence that they might not be:** one `receiver_name` (2 tokens, accented Latin characters) plus a 3-token `street_name` and a `zip_code` appear together as a coherent shipping address in a *planning review document* and a *fix diff* — exactly where a copy-pasted live API response would land. The dev stack did run against ~38 real orders.

> ### N-11 CLOSED — 2026-08-07, by direct inspection of the values
>
> **Verdict: synthetic. No real buyer PII is committed.** The sweep correctly refused to guess from
> structure; the content settles it. Named here because naming them is what makes the finding
> verifiable — none of these identifies a real person:
>
> - `João Silva` — the Brazilian equivalent of "John Smith", with `nickname` `JOAOSILVA` derived from it
> - `Avenida Rio Branco`, `street_number` `1`, `city_name` `Rio de Janeiro` — a landmark avenue at
>   number 1. `zip_code` `20040-002` is that public street's real CEP, which is a property of the
>   street, not of a person
> - CPF `12345678909` — the canonical ascending test CPF this document already identified as the
>   textbook fake
> - CNPJ `11222333000181` — the textbook test CNPJ
> - `ACME`, `Ana`, `NB`, and `testdata/order_body.json:19` which reads `"nickname": "REDACTED"`
>
> The "might not be" hypothesis — a coherent address triple landing in a planning review and a fix
> diff, exactly where copy-pasted live data would — is refuted by the content: a landmark address at
> number 1 paired with a placeholder name. Note also that
> `.mnfs/MIS-007-ml-sync/planning-reviews/p7-seat4-star7-r01.md:49` is itself an argument that this
> exact key set must never be persisted — it used a fixture address as its example.
>
> **No action required. This finding does not block the flip.**

> ### §5.2 item 3 CLOSED — 2026-08-07, all committed images inspected
>
> The six uninspected screenshots were opened. `git log --all --diff-filter=A` confirms **exactly 8
> raster images were ever committed** (all still in tree) plus 5 `public/*.svg` files that are
> leftover Next.js scaffolding. All 8 are now visually inspected and clean: synthetic `T029`
> installations, credential fields reading `Not connected` / `No active credential`, no PII, no
> tokens. The marketplace grid exposes baseline commission percentages per provider — business data,
> already covered by the accepted inventory in §4.
>
> **No action required. This finding does not block the flip.**

**Action (superseded by the closure box above — retained for the record):** open the ~6 files above, look at those specific fields, and confirm. If any is a real buyer, replace it in the tree and treat it as an LGPD disclosure the same way as B-1 (the history caveat applies identically — but for PII, scrubbing the tree plus a rewrite is more defensible than for a credential, because there is nothing to "rotate").

### 3.2 — Binary evidence (no text grep can reach these)

Ten images and two spreadsheets are tracked at HEAD; twelve images/spreadsheets total have ever been committed. This is a small enough set to be certain about.

- **Visually inspected (4 of 10):** `04-madeira-manual-credentials-panel.png` (the highest-risk one — the API Token field is **empty**), `02-amazon-panel-lwa.png`, `2026-04-13-t029-ui-mercado-livre.png`, `2026-04-13-t029-integrations-overview.png`. All show synthetic `T029` test installations, `Not connected` credentials, no PII, no tokens.
- **Not visually inspected (6):** `01-marketplaces-grid.png`, `03-shopee-blocked-panel.png`, `2026-04-13-t029-ui-magalu.png`, `2026-04-13-t029-ui-shopee.png` — same two capture sessions, same UI family, same synthetic data set as the four inspected. Low risk, but see §5.
- **Spreadsheets:** `example-erp.xlsx` (4 KB) and `identity-rejections.xlsx` (2 KB) contain no shared strings, no CPF/CNPJ, no emails — synthetic fixtures. The deleted `templates/catalogo-exemplo.xlsx` (25 KB, 26 rows, from commits `db6c0764`/`c41fbfef`) contains zero email/CPF/CNPJ matches across every part of the OOXML package — a product catalogue template.

---

## 4. INVENTORY — what a competitor or attacker learns (not objecting, just listing)

The operator has already accepted business-logic exposure. This is what that means concretely.

**ERP / Sankhya internals**
- 23 distinct `TGF*` Oracle table names and 16 distinct `AD_*`/`TSI*` table names, with column names, join logic and filter predicates. Files mentioning Sankhya: 423.
- The vendability rule (`USOPROD`, `CODEMP`, `CODLOCAL` whitelists), the product/order linkage two-hop identity, the mirror-first sync direction.
- ICMS/DIFAL/ST/FCP tax formulas including the CST-keyed single-formula derivation, the `TGFICM` two-axis predicate, `CUSSEMICM`'s two branches, and hardcoded PIS/COFINS at 9.25%. 54 ADRs and 85 SQL migrations document the schema in full.
- The named application DB role (17 chars) appears in `docs/architecture/decisions/001-metalshopping-direct-read.md` and `002-mpc-schema-same-cluster.md` without a password. Role-name disclosure is a minor credential-stuffing aid; not a control.

**Mercado Livre integration**
- 79 distinct `MLB########` listing IDs. These are **public marketplace listings** — anyone can browse them. What they reveal is *which* listings are yours, i.e. the seller's catalogue footprint.
- 4 ML order IDs. Not resolvable without seller authentication.
- Full capability map, rate-limit handling, pagination behaviour, the `sale_fee`-per-unit finding, the `?context=channel_marketplace` pricing quirk, and the shape of every request the system makes.
- Baseline commission percentages per marketplace, visible in committed screenshots (Amazon 12.0%, Leroy Merlin 18.0%, Madeira Madeira 15.0%).

**Hosts, IPs, ports**
- Public IPs found (`54.88.218.97`, `18.215.140.160`, `18.213.114.129`, `18.206.34.84`, `35.245.*`, `35.236.*`, `35.186.*`) are **Mercado Livre's published webhook source ranges**, documented in `.mnfs/MIS-007-ml-sync/research/*` and `docs/design/MIS-007-ML-SYNC-DESIGN.md`. Public vendor data, not ours.
- All `10.x` matches are semantic-version strings, not addresses. No RFC1918 address of any real internal host is committed.
- Internal hostnames are `marketplace-central.local`, `postgres`, `backend`, `frontend` — Docker-network names, not routable. `app.test` / `provider.test` / `*.invalid` are reserved test domains.
- **No genuinely reachable internal hostname and no exposed port is disclosed.** Nothing in this section is a security control.

**Identity / process**
- 1 distinct commit author identity (personal gmail address) across all 2109 commits — see N-7.
- 1345 `.mnfs/` files and 45 wiki pages expose the entire agent-harness working method, every dispatch ledger, every gate verdict, every debt row, and roughly a dozen agent transcript logs up to 1.3 MB. This is a candid record of how the software was built, including its failures. Nothing secret; consider whether it is the intended public face.
- Internal identifiers exposed: mission IDs (`MIS-001`..`MIS-008`), milestone/chip IDs, debt IDs (`D-01`..`D-122`), ADR numbers. No external ticket system, no customer identifiers.
- Legacy product names `MetalShopping` (36 files) and `MetalDocs` (8 files) appear, plus `conexus` (5 files) — early-lineage naming that reveals adjacent internal projects.

---

## 5. COVERAGE AND LIMITS

### What was actually searched

- **Object-level, not path-level.** The primary sweep streamed **every blob reachable from every ref** (`git rev-list --objects --all` → 24 761 objects → `git cat-file --batch`) through pattern matchers. This is complete with respect to content: it does not depend on knowing filenames, does not miss renamed files, and does not miss files deleted long ago. Where a finding needed attribution, blobs were mapped back to paths and commits.
- **Path-level completeness** was established independently from the same object listing (4326 distinct paths, all refs) rather than from `--diff-filter=A`, which rename detection can hide entries from.
- **Token shapes swept across all history:** `AKIA*`, `ASIA*`, `ghp_/gho_/ghs_/ghu_/ghr_`, `github_pat_*`, `sk-ant-*`, `sk-*`(40+), `xoxb/xoxp/xoxa/xoxr/xoxs-*`, `AIza*`, `-----BEGIN * PRIVATE KEY-----`, JWT triples, `APP_USR-*`, `TG-*`. **All zero in reachable history** except the two synthetic families documented in N-6, and the PEM/`sk-` hits which were traced to **unreachable local-only blobs** (Go stdlib testdata and a Danish word list where `sk-` occurs inside words) — see N-13.
- **Connection strings:** every `scheme://user:pass@` in every reachable blob was extracted and classified by scheme, user length, password length, password character class, self-referentiality, and variable-reference-vs-literal. 19 distinct DSN signatures; all resolved. `jdbc:`, `mongodb+srv://`, `oracle://`, `Data Source=`, `User Id=` → zero hits.
- **Credential-assignment sweep** covering both quoted and unquoted forms of `PASSWORD/PASSWD/PWD/SECRET/TOKEN/APIKEY/API_KEY/ACCESS_KEY/PRIVATE_KEY/PGPASSWORD`, plus project-specific `MPC_ENCRYPTION_KEY`, `MC_DATABASE_URL`, `SANKHYA_ORACLE_*`, `MPC_PROVIDER_MERCADOLIVRE_CLIENT_*`. Every hit was resolved to variable-reference, placeholder, self-referential dev default, synthetic test fixture, or (once) a real credential.
- **All three `*example*` env files** were key-by-key value-length audited: `deploy/env.production.example` has **empty values** for every Sankhya Oracle field and both ML OAuth fields, and placeholder-marked values for `POSTGRES_PASSWORD` and `MPC_ENCRYPTION_KEY` (32 chars, shape `WORD_word_9_word_word`, carries a placeholder marker, preceded by the comment `Generate: openssl rand -hex 16`). `docker/dev/env.container.example` has an empty `NGROK_AUTHTOKEN`. `apps/web/.env.example` has one non-secret URL.
- **PII:** formatted and unformatted CPF/CNPJ, `cpf`/`doc_number`/`documento`/`identification` JSON fields, `buyer`/`nickname`/`first_name`/`last_name`/`receiver_name`, address fields, `zip_code`, phone shapes, and personal email domains — all across full history, with check-digit validation on CPF candidates and cross-file hash-reuse analysis on name values.
- **Binary:** all 12 committed images/spreadsheets enumerated; 4 images read visually; 3 spreadsheets unpacked and pattern-scanned.
- **Infrastructure:** the single workflow `.github/workflows/release-images.yml` (see §5.1), both Compose files, all four Dockerfiles, `.github/CODEOWNERS` (empty at HEAD), and the tracked `.claude/settings.json` (permissions only, no `env` block).

### 5.1 — GitHub Actions behaviour after the flip

`.github/workflows/release-images.yml` is the only workflow at HEAD (a `wiki-lint.yml` existed historically and is gone).

- Triggers: `push` to `main`, `push` on `v*` tags, and `workflow_dispatch`. **There is no `pull_request` or `pull_request_target` trigger**, so a fork's PR cannot cause it to run. This is the single most important property and it is correct.
- Runner: `ubuntu-latest` — GitHub-hosted. **No self-hosted runner**, so no risk of a fork's code executing on your infrastructure.
- Secrets: only `secrets.GITHUB_TOKEN`, the auto-provisioned per-run token, scoped by an explicit `permissions:` block (`contents: read`, `packages: write`). **No custom secret names are referenced**, so nothing to leak via logs and nothing an attacker can enumerate.
- Registry: `ghcr.io`. **Action item, not blocking:** confirm the visibility of the published GHCR packages separately. Package visibility on GHCR is configured independently of repository visibility; making the repo public does not automatically publish the images, but it does make the package names and the exact build recipe public. Decide deliberately.
- Log exposure: build logs on a public repo are world-readable. This workflow logs only Docker build output. Nothing secret is echoed.

`docker-compose.yml` and `deploy/docker-compose.prod.yml` were audited line by line. The prod file parameterises every secret (`${POSTGRES_PASSWORD:?set in .env}`, `${MPC_DOMAIN:?set in .env}`). The dev file has the self-referential dev password (N-1). No Dockerfile bakes a secret; the only `ENV`/`ARG` values are the public Oracle Instant Client download URL and `CGO_ENABLED`.

### 5.2 — What I could not do, and honest false-negative risk

**A sweep that finds nothing is not proof of nothing.** Specifically:

1. **Pattern-based detection cannot see an unpatterned secret.** A password that is an ordinary Portuguese phrase, an API key with no vendor prefix, or a secret embedded in prose without a `KEY=` marker will not match any regex I ran. B-1 was caught only because it sat next to `PGPASSWORD=`; had it been written as "a senha é ..." in a sentence, I would have missed it. No automated tool would have found it either. **This is the main residual risk.**
2. **No dedicated scanner was used** (gitleaks/trufflehog were out of scope per the constraints, correctly — they must not be installed here). Those tools carry hundreds of vendor-specific rules I do not. If the operator wants a second opinion, GitHub's own **secret scanning with push protection** becomes available on public repos at no cost and should be enabled *at the moment of the flip*; it will retro-scan the history and will find anything vendor-shaped that I missed.
3. **Six of ten committed images were not opened.** They are from the same two capture sessions as the four that were, and those four were clean synthetic data — but a screenshot is opaque to every text method used here. Ten images is a five-minute manual review; do it.
4. **`.mnfs/` agent transcript logs are large and were pattern-scanned, not read.** `planner-sol-medium.log` (1.3 MB), `agent__gate-r6-sol.log` (776 KB) and several others contain raw model output that could quote anything the agent saw. They passed every credential and PII pattern, but prose-form disclosure (see limit 1) would survive.
5. **The `legacy` remote is a different repository** (`leandrotcawork/marketplace-central`). Its refs were included in this sweep because they are fetched locally, so the *content* is covered — but flipping `origin` does not change `legacy`'s visibility, and conversely, if `legacy` is already public, then everything in `legacy/master`, `legacy/main` and `legacy/feat/llm-wiki-m1` is *already* exposed regardless of what you do to `origin`. **Check `legacy`'s current visibility before assuming this decision is still yours to make.**
6. **Local `main` is 98 commits ahead of `origin/main`.** The sweep covered the superset (all local refs), so no unswept commit can reach GitHub. But note that B-1 and B-2 are already on `origin/main` today — they are not something the next push would introduce.
7. **Another session holds write access to this checkout.** `HEAD` moved during the sweep (from `1473e863` to `5909f960`). Commits made after `5909f960` were not swept. Re-run the credential-assignment sweep against the final tip before flipping.

### 5.3 — What would make this certain

Enable GitHub secret scanning + push protection at the instant of the flip; open the six unreviewed screenshots; eyeball the nine PII values in §3.1; confirm the §3.1 N-3 lane password is disposable; and re-run the sweep against the final `main` tip after the other session finishes.

---

## 6. COMMANDS RUN

Reproducible; all read-only. Paths shown relative to the repo root. `$SHAS` is a file containing every reachable object SHA, built once:

```bash
git rev-list --objects --all > objmap.txt
awk '{print $1}' objmap.txt | sort -u > reachable_shas.txt
```

**Scope and shape**
```bash
git rev-list --all --count                       # 2109
git rev-list --count origin/main                 # 1981
git ls-files | wc -l                             # 3122
awk '{print $2}' objmap.txt | sort -u | wc -l     # 4326 distinct paths, all refs
git branch -a ; git remote -v
git count-objects -vH
```

**Path-level sensitive-file sweep (authoritative, all refs)**
```bash
awk '{print $2}' objmap.txt | sort -u > allreachablepaths.txt
grep -Ei '(^|/)\.env' allreachablepaths.txt
grep -Ei '\.(pem|key|p12|pfx|jks|ovpn|tfstate|keystore|ppk|crt|cer|kdbx|gpg|asc)$' allreachablepaths.txt
grep -Ei 'id_rsa|id_ed25519|\.ssh/|\.npmrc|\.netrc|\.pgpass|(^|/)\.aws|docker/config\.json' allreachablepaths.txt
grep -Ei 'settings(\.local)?\.json|\.mcp\.json' allreachablepaths.txt
git log --all --diff-filter=A --name-only -- '*.env' '*.pem' '*.key' '*.p12' '*.pfx' '*.jks' '*id_rsa*' '*.ovpn' '*credentials*' '*secrets*' '*.tfstate'
```

**Content sweep across all history** — one streaming pass per matcher family, blob-attributed:
```bash
git cat-file --batch --buffer < reachable_shas.txt | tr -d '\000' | awk '
  /^[0-9a-f]{40} blob [0-9]+$/            { cur=$1; next }
  /^[0-9a-f]{40} (tree|commit|tag) [0-9]+$/ { cur=""; next }
  function hit(k){ if (cur!="" && !(k SUBSEP cur in seen)) { seen[k SUBSEP cur]=1; print k" "cur } }
  { if ($0 ~ /AKIA[0-9A-Z]{16}/)                                   hit("AWS_AKIA");
    if ($0 ~ /gh[pousr]_[A-Za-z0-9]{36}/)                          hit("GITHUB_PAT");
    if ($0 ~ /sk-ant-[A-Za-z0-9_-]{20}/)                           hit("ANTHROPIC");
    if ($0 ~ /xox[baprs]-[A-Za-z0-9-]{12}/)                        hit("SLACK");
    if ($0 ~ /AIza[0-9A-Za-z_-]{35}/)                              hit("GOOGLE");
    if ($0 ~ /BEGIN [A-Z ]*PRIVATE KEY/)                           hit("PRIVKEY");
    if ($0 ~ /eyJ[A-Za-z0-9_-]{15,}\.eyJ[A-Za-z0-9_-]{15,}\./)     hit("JWT");
    if ($0 ~ /APP_USR-[0-9]{10}/)                                  hit("ML_APPUSR"); }'
# blob -> path: awk 'NR==FNR{p[$1]=$2;next}{print $1" "p[$2]}' objmap.txt <hits>
```
Same harness re-run for: `client_secret` literals, `MPC_ENCRYPTION_KEY`, `SANKHYA_ORACLE_*`, quoted and unquoted `PASSWORD/SECRET/TOKEN/APIKEY/PGPASSWORD` assignments, `Bearer` headers, Oracle TNS `(DESCRIPTION=`/`SERVICE_NAME=`/`SID=`, every `scheme://user:pass@` DSN (classified by user/password length + charset + self-referentiality + literal-vs-varref), CPF/CNPJ formatted and unformatted, `nickname`/`first_name`/`last_name`/`receiver_name`/`street_name`/`zip_code`/`email` populated values, and personal email domains.

**Placeholder-vs-real discrimination** (values never printed)
```bash
# CPF check digits — reports only VALID / INVALID / REPEATED / ASCENDING
grep -oE '[0-9]{11}' <file> | sort -u | awk -f cpf.awk

# env example files — key names and value LENGTHS only
awk -F= 'NF>1 && $0!~/^#/ {k=$1; v=substr($0,index($0,"=")+1);
  printf "key=%s vallen=%d marker=%s\n", k, length(v),
  (tolower(v) ~ /change|example|your|xxx|placeholder|todo|replace|dummy|fake|sample|\$\{|</ ? "PLACEHOLDER":"NONE")}' deploy/env.production.example

# token/password shape — length + character classes, never content
sed -n '34p' deploy/env.production.example | sed 's/^[^=]*=//' \
  | sed -E 's/[a-z]+/w/g; s/[A-Z]+/W/g; s/[0-9]+/9/g'

# cross-file fixture reuse — md5 of value, value itself discarded
... | while IFS='|' read k v; do echo "$k $(printf '%s' "$v" | md5sum | cut -c1-8) $f"; done
```

**Attribution of the blocking findings**
```bash
git log --all --format='%H %ad %s' --date=short -- .claude/settings.local.json
git log --all --find-object=37319a7dca07455b56b0301018a6e5b6b297f694
git log --all --format='%H %ad %s' --date=short -- docs/superpowers/plans/2026-04-06-pricing-simulator-v2.md
git cat-file -e origin/main:docs/superpowers/plans/2026-04-06-pricing-simulator-v2.md   # exit 0 => already public-bound
git merge-base --is-ancestor 4743f778e3714ab8b729735b61d0088d3f5c8078 origin/main       # true
# same-credential proof, value never emitted:
git cat-file blob 37319a7d... | grep -oE 'postgres://[^:@" ]+:[^@" ]+@' | sed -E 's|postgres://[^:]+:||; s|@$||' > needle
git grep -l -F -f needle -- .        # -> docs/superpowers/plans/2026-04-06-pricing-simulator-v2.md
rm -f needle
```

**Unreachable-object proof (N-13)**
```bash
git rev-list --objects --all          | grep -c '^<sha> '   # 0
git rev-list --objects --all --reflog | grep -cE '^(<sha>|...) '  # 0
```

**Commit metadata**
```bash
git log --all --format='%an <%ae>' | sort -u | wc -l   # 1
git log --all --format='%ae' | sort -u | sed -E 's/.*@//' | sort | uniq -c
```

**Infrastructure**
```bash
sed -n '1,20p' .github/workflows/release-images.yml
grep -nE 'secrets\.|runs-on|registry|permissions:|pull_request' .github/workflows/*.yml
grep -nE '^\s*(- )?[A-Z_]+:' docker-compose.yml deploy/docker-compose.prod.yml
grep -nE '^(ENV|ARG)' docker/{dev,prod}/*.Dockerfile
unzip -l / unzip -p  <each committed .xlsx>            # into scratchpad, never the repo
```

**Not run, deliberately:** `git push`, `reset`, `revert`, `stash`, `clean`, `gc`, `prune`, `filter-repo`, any history rewrite, `gh repo edit`, any visibility change, `docker inspect`, bare `printenv`, `Get-ChildItem Env:`, sourcing any `.env`, any test suite, any harness or gate command, and any contact whatsoever with Oracle/Sankhya. No file in the repository was created or modified except this report. All temporary artefacts were written to the session scratchpad outside the repo.
