# P7 r01 — Claude crew seat 4/5 verbatim (★7 Security Posture, adversarial)

```yaml
round: r01
seat: 4/5
scope: "★7 Security Posture (adversarial)"
model: opus (cold, read-only)
manifest: planning-reviews/p7-input-r01.sha256
manifest_digest: 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e
task_id: a19b3e80e55f0fe20
transport: task-notification result (in-context), persisted same turn
verdict: "★7 FAIL (D1 shipment-raw PII untargeted; D2 resource format unowned; D3 source_ip derivation unpinned) + 5 advisories"
persisted: 2026-08-01
```

Verbatim result below (unedited):

---

# ★7 Security Posture — cold adversarial pass (seat 4 of 5)

**Seat:** Claude crew seat 4/5 — ADVERSARIAL, scope = ★7 Security Posture ONLY. No seven-★ verdict computed.
**Mission root:** `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync`
**Rubric:** `C:\Users\leandro.theodoro\.claude\plugins\cache\mnfs-harness\mnfs-workflow\0.2.0\skills\mission-planning\references\readiness-review-rubric.md` (★7 procedure applied verbatim)
**Manifest:** `planning-reviews/p7-input-r01.sha256`, 65 entries; declared trailer at `:66` — `# top-level digest (sha256 of the 65 sorted entries above): 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e` — matches the dispatch expectation. (Honesty note: I read the digest line; I did not recompute per-file hashes — read-only pass.)

---

## ★7 Security Posture — **FAIL**

**Rationale (one line):** the scope has three in-scope PII/auth surfaces; the webhook forgery surface and the `billing_info` fiscal surface ARE targeted (mitigation + Security-typed criteria), but the **buyer delivery PII carried in `order_shipments.raw` is permitted with zero mitigation and zero criterion — and it directly contradicts the mission's own Security-typed blocking rule** — while two further guards on the public endpoint (`resource` format, `source_ip` derivation) are owned by nobody.

**Cited excerpt (verbatim), primary:** `research/orders-persistence-interface-contract.md:55`

> `- \`raw jsonb\` NULL (shipment payload — PERMITIDO aqui, ADR-03) + marcador de truncamento.`

against `validation-contract.md:126-127` (mission, Security-typed MIS07-C5):

> `nenhuma migração da missão adiciona \`raw jsonb\` a tabela com dado` / `fiscal/comprador (ADR-03/R-6); colunas fiscais tipadas SÓ`

`order_shipments` IS a table with buyer data by its own field list — `orders-persistence-interface-contract.md:53`: `- Destino (o que a tela já mostra): \`receiver_name text\` NULL, \`dest_city text\` NULL,`.

---

### Defect locus D1 (load-bearing) — buyer delivery PII persisted raw, untargeted

- **Locus:** `research/orders-persistence-interface-contract.md:55` — offending token: **`raw jsonb` … `PERMITIDO aqui`** on `order_shipments` (migration **0088**, `M-02-sync-core-seam/F-01-core-ddl/feature.md:32`).
- **Amplifying locus:** `M-03-orders-shipment-persist/F-01-ml-ingest-readers/feature.md:62` — offending token: **`shipment raw inteiro permitido (IC-03)`**.
- **Why it is PII, from the artifacts + repo (not intuition):** the same feature brief states the shipment reader fetches the destination block — `M-03…/F-01…/feature.md:31`: `logistic_type, tracking_number, sla/status.limit, costs gross/seller, receiver_address`. In the repo, the shipment payload's destination carries the buyer's name and full street address: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipping_reader_test.go:30` — `"destination":{"receiver_id":74425755,"receiver_name":"João Silva","shipping_address":{…"zip_code":"20040-002","street_name":"Avenida Rio Branco","street_number":"1"…}}`; and `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipping_reader.go:88` — `// schema — the buyer delivery address lives under `destination.shipping_address`.` The repo **already classifies exactly these keys as PII**: `apps/server_core/cmd/mlprobe/main.go:41-43` — `"doc_number": true, "receiver_name": true, "receiver_phone": true,` / `"street_name": true, "street_number": true, "address_line": true,` / `"comment": true, "zip_code": true,`.
- **Why untargeted (silent omission, the exact ★7 failure mode):** the only PII mitigation in the spine is document-scoped — `mission.md:146-149`: `**Raw de \`billing_info\` NUNCA persiste** (documento + endereço fiscal — \`buyer_fiscal_reader.go:59-94\`); campos tipados sim. … Teste: fixture com documento → nenhuma coluna persistida o contém.` Every Security-typed criterion inherits that narrowing and therefore **passes vacuously over `order_shipments.raw`**: `M-02-sync-core-seam/validation-contract.md:153-154` (`nenhuma coluna \`raw\`/jsonb de billing_info em tabela com dado fiscal/comprador`); `M-03-orders-shipment-persist/F-01…/feature.md:84` (`Assert: raw não contém substring \`billing_info\``); `M-08-webhook-ingest/validation-contract.md:150` (`Fixture de forja com documento sintético`). A fixture containing a *document* proves nothing about a raw blob containing a *name + street + CEP*. This is the "guard parcial sob frase total" class: ADR-03's title asserts `PII nunca` (`mission.md:145`) while its body authorizes raw persistence in tables that hold buyer data (`mission.md:146-147`: `persistência de \`raw jsonb\` só em \`listings\`/\`orders\`/\`order_shipments\``) — note `orders` is authorized there too, while IC-03 defines no `raw` column for `orders`, leaving an implementer free to add one over a table whose 0089 block is exactly the buyer-fiscal columns.
- **Also self-contradicting:** R-6's own trigger fires on this design as written — `mission.md:373`: `migração adicionando \`raw jsonb\` a tabela com dado fiscal/comprador` — and MIS07-C5's grep command (`validation-contract.md:123`, `grep -rn "raw" apps/server_core/migrations/0086*..0095*`) would return 0088 and block the mission at QA.
- **Yes-if (minimal, inside approved scope / existing contracts):** EITHER (a) delete the `raw jsonb` column from `order_shipments` in `orders-persistence-interface-contract.md:55` (MIS07-C5 already forbids it — no new scope), OR (b) keep it and close the guard at three existing loci: (1) amend ADR-03 (`mission.md:145-149`) so the rule reads `raw de billing_info E do bloco destination/shipping_address NUNCA persiste`, naming the key set the repo already declares PII (`cmd/mlprobe/main.go:38-44`); (2) extend M02-C7 `Expected` (`M-02-sync-core-seam/validation-contract.md:153-154`) to name `order_shipments.raw` and require the scrub, keeping Type: Security; (3) change the assert at `M-03…/F-01…/feature.md:84` from `raw não contém substring \`billing_info\`` to a fixture asserting the persisted raw contains none of `receiver_name|street_name|street_number|zip_code|comment|receiver_phone`.

### Defect locus D2 — attacker-controlled `resource` reaches an authenticated outbound request with no pinned format, and nobody may decide it

- **Locus:** `research/webhook-inbox-interface-contract.md:51` — offending token: **`\`resource text\` NULL (ex.: \`/orders/2000012345\`)`** — an *example*, never a constraint; no CHECK, no regex, no `malformed` rule keyed on format (the only `malformed` trigger is `:71` `body não-JSON/sem resource`).
- **Consumer locus:** `M-08-webhook-ingest/F-02-worker-callback/feature.md:28` — `order id do \`resource\` (\`/orders/<id>\`), chama \`IngestOrder\``, i.e. an unauthenticated, attacker-supplied path segment is interpolated into a bearer-token-authenticated ML call (`IngestOrder(ctx, tenant, installation, providerOrderID)`, `research/sync-ingest-ports-interface-contract.md:32`).
- **Why unowned:** IC-04 forbids the executing feature from deciding it — `research/webhook-inbox-interface-contract.md:156`: `- Postura de segurança; enum de status; regra de topic ignorado; cap de body; não-consumo` under `## Must Not Decide In Feature Execution`. The contract does not decide it and the feature may not: the guard has no owner. The one Security criterion that touches forgery tests only *semantics*, not *shape* — `M-08…/validation-contract.md:45-46`: `POST forjado plausível (user_id válido, resource inventado)`.
- **Yes-if:** pin the accepted shape in `research/webhook-inbox-interface-contract.md:51` — `resource` accepted for processing only when it matches `^/orders/[0-9]+$`, anything else terminal `malformed` (200 preserved, always-200 untouched) — and add one case to M08-C1 `Evidence` (`M-08…/validation-contract.md:45`): a `resource` with traversal/encoded segments (e.g. `/orders/1/../../users/me`, `%2e%2e%2f`) → row `malformed`, ZERO outbound ML call. Grounded in ADR-11 as already ratified (`mission.md:200-201`: `Só \`resource\`+\`topic\` usados, só como ponteiro p/ refetch AUTENTICADO`) — no new scope.

### Defect locus D3 — `source_ip`/`ip_official` derivation unpinned; Q2's second validation criterion is unfalsifiable as written

- **Locus:** `research/webhook-inbox-interface-contract.md:56` — offending token: **`\`source_ip inet\` NOT NULL; \`ip_official boolean\` NOT NULL`** with no stated source for the value.
- **Why it matters:** the ratified production topology is a tunnel — `mission.md:61-62`: `webhook topic \`orders_v2\` SOMENTE (gate P1 2026-07-31); ngrok domínio fixo` / `como URL de produção`. Nothing in the 65-entry tree names how the IP is obtained (grep for `forwarded|proxy|XFF|RemoteAddr` across the mission tree returns zero hits in artifacts), and the repo has no client-IP helper at all (grep `X-Forwarded-For|RemoteAddr|realip|ClientIP` over `apps/server_core` → no matches). With the socket peer being the tunnel agent, `ip_official` is `false` for 100% of traffic and the mission-level Security assertion passes vacuously; with a header, the value is attacker-settable. Either way the Q2 criterion at `mission.md:344` — `notificação forjada não injeta dado (só fetch idempotente); IP não-oficial aparece em log` — cannot discriminate, and M08-C1's `IP off-allowlist … LOGADO e request processada` (`M-08…/validation-contract.md:48-49`) is satisfied by a constant.
- **Yes-if:** pin the derivation in `research/webhook-inbox-interface-contract.md:56` (socket peer, or a named trusted header under the ngrok hop, stating which) AND make M08-C1 discriminating: add the positive control — a request presenting an official IP (one of the 4 at `:57`) must yield `ip_official=true` while a non-official one yields `false`, both 200 — OR, if the operator's log-only decision (`mission.md:121`) is meant to accept a non-discriminating field, record it explicitly under `## Non-Functional Scope` (`mission.md:351-354`) as declined-with-reason so it stops being an asserted mitigation.

---

## What DOES pass (recorded so the fold does not double-count)

- Public endpoint authentication is **explicitly declined-with-reason by the operator**, not silently omitted: `mission.md:121` — `| 3 | build/runtime | Postura de segurança do POST /webhooks/{provider} público? | Hint não-confiável + refetch autenticado + IP log-only | Hint não-confiável + IP log (2026-07-31) |`; restated at `research/webhook-inbox-interface-contract.md:145` (`Sem cookie/CORS/auth — endpoint público por design; mitigação = conteúdo nunca confiável`).
- Body-never-becomes-domain-data (ADR-11) has mitigation + Security-typed criterion: `M-08…/validation-contract.md:38-42` (M08-C1, `Type: Security`) and `research/webhook-inbox-interface-contract.md:109-111`.
- Live ML write authorization gate exists and is Security-typed and blocking: `M-08…/validation-contract.md:132-144` (M08-C6, `Blocking failure: callback registrado sem autorização registrada`).
- Fiscal-document PII has mitigation (ADR-03) + Security-typed criteria (M02-C7 `M-02…/validation-contract.md:145-159`; MIS07-C5 `validation-contract.md:115-131`).
- Oversized-body and slow-body are targeted: `research/webhook-inbox-interface-contract.md:147` (`Cap de body 64KB ANTES do parse (LimitReader)`) + interactive 15s class proved by trickle (`M-08…/validation-contract.md:66-70`).
- Multi-tenant scoping is present on every new table: `channel-fees-interface-contract.md:46,103`; `divergences-interface-contract.md:44,82`; `orders-persistence-interface-contract.md:46`; `listings-sync-interface-contract.md:65`; `/sync/health` takes tenant from ctx under existing auth (`sync-health-interface-contract.md:68,175`).

---

## Advisory findings (never flip a ★ alone; not in the yes-if list)

1. **Unbounded write amplification from an unauthenticated caller.** No retention, pruning, or per-IP cap for `notifications_inbox`; each accepted POST persists up to 64KB in `raw_body` (`webhook-inbox-interface-contract.md:54`). The named mitigation — `M-08…/milestone.md:64`: `mitigação IC-04 (cap 64KB, dedupe, ponteiro-nunca-dado,` — bounds row *size* and *duplicates-with-`_id`*, never row *count* from a non-ML caller. Impact: DB growth / disk exhaustion.
2. **Reject-vs-truncate divergence at the security cap.** `M-08…/validation-contract.md:68`: `>64KB → rejeitado sem ler o resto` versus `webhook-inbox-interface-contract.md:54` and `M-08…/F-01…/feature.md:60`: `body >64KB→truncado com marcador`. "Rejeitado" also risks a non-200, colliding with always-200 (`:40`).
3. **False field name at a security-relevant decode point.** `M-03…/F-01…/feature.md:31` names `receiver_address`; the repo states that key does not exist — `shipping_reader_test.go:29`: `// NOT a flat \`receiver_address\` field (that key does not exist in the schema).` A scrub written against the wrong key would be a silent no-op.
4. **Inbox aggregation has no tenant predicate** while the port is tenant-parameterized: `sync-health-interface-contract.md:39` (`WebhookStats(ctx, tenantID)`) vs `:149-150` (`inbox agregados (\`status='received'\` count; …`) and `webhook-inbox-interface-contract.md:49` (`tenant_id` NULL). Deferred by the recorded assumption `mission.md:110-111` (`Installation/tenant única existente … tudo já tenant-scoped`) — but the assumption's second clause is not true of this table.
5. **Response-shape oracle.** The 500-on-INSERT-failure (`webhook-inbox-interface-contract.md:84`) is the only externally distinguishable response on an unauthenticated route, exposing a DB-health signal. Low impact; log-only posture is consistent.

(No further advisories; nothing dropped.)

---

## Manifested paths actually read (full-file reads)

Absolute, under `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\`:

- `...\mission.md`
- `...\validation-contract.md`
- `...\research\webhook-inbox-interface-contract.md`
- `...\research\orders-persistence-interface-contract.md`
- `...\research\sync-health-interface-contract.md`
- `...\research\sync-ingest-ports-interface-contract.md`
- `...\research\external-ml-api-facts.md`
- `...\M-01-ml-client-hardening\milestone.md`
- `...\M-02-sync-core-seam\F-01-core-ddl\feature.md`
- `...\M-02-sync-core-seam\validation-contract.md`
- `...\M-03-orders-shipment-persist\milestone.md`
- `...\M-03-orders-shipment-persist\F-01-ml-ingest-readers\feature.md`
- `...\M-03-orders-shipment-persist\F-02-ingest-order-v1\feature.md`
- `...\M-03-orders-shipment-persist\validation-contract.md`
- `...\M-06-orders-backfill-decomposition\F-02-decomposition-camada3\feature.md`
- `...\M-08-webhook-ingest\milestone.md`
- `...\M-08-webhook-ingest\F-01-inbox-endpoint\feature.md`
- `...\M-08-webhook-ingest\F-02-worker-callback\feature.md`
- `...\M-08-webhook-ingest\validation-contract.md`
- `...\M-09-sync-observability\F-01-sync-health-endpoint\feature.md`
- (input manifest, not one of the 65) `...\planning-reviews\p7-input-r01.sha256`

**Unread-path sweep (grep over ALL 65 manifested entries, content inspected at every hit):** patterns `receiver|billing_info|PII|scrub|fiscal|raw` (case-insensitive), `resource`, `tenant`, `ngrok|token|secret|credential|auth|CORS|cookie|público|allowlist|source_ip|spoof|rate.?limit|flood|retention|purge`, `forwarded|XFF|proxy|peer|RemoteAddr`. Every hit in M-04, M-05, M-06, M-07, M-09, IC-01/IC-02/IC-07, `research/codebase-*.md`, `research/p5-prerequisites.md` and `planning-reviews/*` was read in context; no additional auth/PII/multi-role surface was found there, and no artifact anywhere pins `source_ip` derivation, `resource` format, inbox retention, or a shipment-raw scrub.

**Repo files read for fact-check only (not manifested, read-only):** `C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\cmd\mlprobe\main.go`, `...\apps\server_core\internal\modules\connectors\adapters\mercado_livre\shipping_reader.go`, and grep over `...\apps\server_core\internal\modules\connectors\adapters\mercado_livre\shipping_reader_test.go`.
