# P7 Claude-side Readiness — MIS-007-ml-sync — round r01 (computed fold)

```yaml
id: P7-CLAUDE-R01
type: readiness-review-claude-fold
round: r01
manifest: planning-reviews/p7-input-r01.sha256
manifest_digest: 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e
manifest_reverified: "2026-08-01 — 65/65 hashes recomputed post-crew, DRIFT: NONE"
crew: 5 cold mission-reviewer seats (Opus), parallel, read-only
seat_artifacts:
  - planning-reviews/p7-seat1-star1-star5-r01.md   # ★1 FAIL, ★5 PASS
  - planning-reviews/p7-seat2-star2-star3-r01.md   # ★2 FAIL, ★3 FAIL
  - planning-reviews/p7-seat3-star4-star6-r01.md   # ★4 FAIL, ★6 PASS
  - planning-reviews/p7-seat4-star7-r01.md          # ★7 FAIL (adversarial)
  - planning-reviews/p7-seat5-doublepass-r01.md     # ★2 FAIL, ★7 PASS (double-pass)
fold_rule: "computed, not chosen — a ★ FAILS when ANY covering seat returns a valid FAIL; FAIL never downgraded"
verdict: "Needs revision"
sol_dispatched: false   # procedure: Sol HIGH only after Claude-side Ready
round_counts_toward_cap: true   # round 1 of max 3
```

## Folded per-criterion result

| ★ | Seats covering | Fold | Loci |
|---|---|---|---|
| ★1 Completeness | 1 | **FAIL** | B-1 |
| ★2 Consistency | 2, 5 | **FAIL** | B-2, B-3 |
| ★3 Seam Ownership | 2 | **FAIL** | B-4 |
| ★4 Verifiability | 3 | **FAIL** | B-5, B-6 |
| ★5 Traceability | 1 | PASS | — |
| ★6 Evidence Honesty | 3 | PASS | — |
| ★7 Security Posture | 4, 5 | **FAIL** | B-7, B-8, B-9 (seat 4 valid FAIL prevails over seat 5 PASS — fold never downgrades) |

## Union of blocking findings (9 loci)

- **B-1 (★1, seat 1)** — `mission.md:28` Outcome promete `comissão/frete camada 2` em listings; frete camada 2 = honesto-desconhecido por contrato (IC-01:65,87), sem produtor, sem critério, sem linha de Non-Scope. Yes-if: remover `frete` do Outcome + linha de Non-Scope nomeando frete-camada-2 declinado com razão.
- **B-2 (★2, seats 2+5, achado independente)** — fork do DTO `divergences[]`: IC-02 declara `{kind, detected_at}`; M-05/F-03:30,70 cita "IC-02 §DTO — kind/expected/observed/timestamps" (seção §DTO não existe; shape mora em §Required Outputs); F-03:38 exige "os dois lados + timestamps", insatisfazível com 2 chaves; M-06/F-03 consome sem shape. Yes-if: IC-02 §Required Outputs declara a projeção completa que o panel precisa (vocabulário de coluna já ratificado) + F-03 cita a seção certa com a lista exata de chaves.
- **B-3 (★2, seat 5)** — ADR-09 "TODO fee em tela carrega (camada, origem, coletado_em)" violado pelo consumidor aprovado: `decomposicao` de IC-03:104-108 tem `comissao_total`/`frete_seller` sem origem/coletado_em, renderizados no PedidoDrawer (M-06/F-03:36-39) sem asserção de proveniência; MIS07-C7 Command não cobre /pedidos. Yes-if: (a) adicionar origem+coletado_em aos termos de fee da decomposicao + asserção no drive M-06 (valores já existem no ledger camada 3 — sem escopo novo).
- **B-4 (★3, seat 2)** — write-set do M-05 fora do grant registrado: F-01/F-02 declaram fiação no `listings/composition/**` (package que M-04 cria e possui) e porta de leitura mirror em `listings/ports/**`; o lock aditivo (mission.md:319-323, célula :303) enumera só application/transport/repository.go; três artefatos enunciam o lock de três formas (advisory A-1). Yes-if: estender o lock a composition/** e ports/** (additive-only) nas TRÊS enunciações, idênticas.
- **B-5 (★4, seat 3)** — MIS07-C2 usa duas estatísticas (Expected p95, Blocking mediana) sobre n=3 onde p95 é indefinido; não-decidível binariamente. Yes-if: UMA estatística nos dois lugares ("as 3 amostras <2s").
- **B-6 (★4, seat 3)** — M01-U2 exige "resposta idêntica à pré-merge" sem obrigar captura de baseline pré-merge (irrecuperável depois); M-07 já carrega a cláusula-modelo. Yes-if: cláusula de baseline before-merge em M-01 Evidence Requirements + prova de M01-U2 = comparação das duas capturas salvas.
- **B-7 (★7, seat 4)** — `order_shipments.raw jsonb` "PERMITIDO" (IC-03:55) carrega PII de entrega (receiver_name, street, CEP — chaves que o repo já classifica PII em cmd/mlprobe/main.go:41-43) sem mitigação nem critério; contradiz MIS07-C5/R-6 ("nenhuma migração adiciona raw jsonb a tabela com dado fiscal/comprador") e o próprio título de ADR-03 "PII nunca"; lista raw-permitido de ADR-03 inclui `orders` sem coluna raw definida. Yes-if aplicado: arm (a) — DELETAR a coluna raw de order_shipments (direção do critério de segurança mais estrito já ratificado), consertar ADR-03 e o assert do M-03/F-01.
- **B-8 (★7, seat 4)** — `resource` (attacker-controlled) vira chamada ML autenticada sem formato pinado; IC-04 proíbe a feature de decidir; guard sem dono. Yes-if: pinar `^/orders/[0-9]+$` no IC-04 (fora do formato → terminal `malformed`, always-200 intacto) + caso de traversal no M08-C1.
- **B-9 (★7, seat 4)** — derivação de `source_ip`/`ip_official` não pinada; sob túnel ngrok o socket peer é o agente → critério Q2 vira constante (passe vácuo) ou header attacker-settable. Yes-if: pinar derivação (header do túnel, nomeado) + controle positivo/negativo no M08-C1 + registrar caráter informativo/log-only (decisão P1 do operador) em Non-Functional Scope.

## Advisories (union; nunca flipam ★; fold-after-verdict a critério do autor)

Seat 1: A1 R-2 cita `missed_feeds 2 dias` (perna fictícia — missed_feeds não-consumido); A2 decisões finas fora do heading Accepted assumptions; A3 fraseado da linha Q1; A4 Handoff cita artefatos fora do manifesto; A5 M-09/F-02 refetch decidido no spec (ok); A6 M09-C6 diferido — MIS07-C8 deve nomeá-lo; A7 R-7 sem critério (processo de hub).
Seat 2: A-1 três enunciações do lock M-05 (absorvida em B-4); A-2 AnunciosPage.tsx fora dos owned paths F-03; A-3 chaves sem installation_id em IC-02/03/07 (FK inexpressível — advisory); A-4 fraseado "assinatura verbatim" M-03/F-02; A-5 "IC-03 §DTO" seção inexistente em M-06/F-03:73; A-6 Error Matrix IC-01 placeholder vs falha típada M-07; A-7 reserva 0094-0095 antes de 0093 (runner por filename — ok factual); A-8 root.go `—` do M-05 não-provável; A-9 âncora exata (positivo).
Seat 3: 1 placeholder de nome de teste M01-C1 (spec-time); 2 fato #11 assumed não espelhado em Accepted assumptions; 3 Evidence Requirements da missão sem arquivo-por-critério; 4 M-05 VC cita memória de sessão p/ ?context=; 5 r07 cita §InboxHealth (arquivo-morto, tratado).
Seat 4: 1 amplificação de escrita sem retenção/cap no inbox; 2 reject-vs-truncate divergem no cap 64KB (M08-C2 vs IC-04); 3 `receiver_address` nome de campo falso em M-03/F-01:31 (repo: chave não existe); 4 agregação de inbox sem predicado de tenant (tenant_id NULL); 5 oracle de shape 500-on-INSERT.
Seat 5: 1 sem critério de flood/rate-limit; 2 `{provider}` slug sem constraint (forka chave de dedupe); 3 retenção source_ip/raw_body; 4 rows `processing` invisíveis na saúde; 5 M09-C6 diferido fácil de perder.

## Disposition

Claude-side verdict **Needs revision** → Sol HIGH NÃO despachado nesta rodada (procedimento §3). Repairs dentro da autoridade de reparo (locais, preservam outcome/escopo/arquitetura/fronteiras/aceite de risco — cada um resolve contradição interna na direção do artefato ratificado mais estrito, sem escopo novo) serão aplicados; advisories baratas folded na mesma passada. Depois: pre-check do autor, manifesto novo r02, crew fria nova. `status: needs_revision` gravado em mission.md.
