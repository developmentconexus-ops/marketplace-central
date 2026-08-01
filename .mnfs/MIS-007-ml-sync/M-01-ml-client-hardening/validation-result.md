# Milestone Validation Result — M-01-ml-client-hardening

```yaml
id: M-01-VR
type: milestone-validation-result
status: passed
owner: hub
parent: MIS-007
created: 2026-08-01
validation_level: QA-0
base_sha: 295e293f
merge_commit: e79c458e
chip_tip: adc7801c
```

## Verdict

**PASS.** 6/6 critérios de milestone + 2/2 de user-drive. Zero blocking failure.

## Evidência por critério

| ID | Verdito | Observável medido |
|----|---------|-------------------|
| M01-C1 | PASS | `TestResilienceDecoratorHonorsRetryAfterHeader` imprime `measured total elapsed = 2.0549565s, inter-call gap = 2.0538955s` — asserção sobre o relógio observado, não sobre config. `Retry-After: 2` honrado. |
| M01-C2 | PASS | `TestResilienceDecoratorTokenBucketThrottlesConcurrentRequests`: gaps observados 99.0724ms / 99.8369ms / 99.8846ms / 99.9549ms, total 400.3427ms sob N goroutines na mesma installation. Limite configurável provado por `TestResilienceDecoratorRateLimitPerMinuteDefaultsWhenUnset` (default existe porque o campo é config, não constante). |
| M01-C3 | PASS | `TestResilienceDecoratorBudgetExhaustedReturnsTypedError` — `RateLimitExhaustedError` tipado com contagem de tentativas. |
| M01-C4 | PASS | `TestGetItemsMultigetPartitionsInBatchesOf20` (45 ids → 20+20+5) + `TestGetItemsMultigetSuccessfulDTOsHaveNonEmptyRaw`. |
| M01-C5 | PASS | `TestResilienceDecoratorStockWriteDoesNotRetry`: `stock write 429 propagated after 67.006ms with 1 attempt(s)`. |
| M01-C6 | PASS | `go build`/`go vet`/`go test ./...` verdes na main integrada; 95 testes no pacote `mercado_livre` (`ok ... 11.649s`). `git diff --stat` do chip = 11 arquivos, todos sob `connectors/adapters/mercado_livre/` + `.mnfs` próprio. Zero migração, zero UI. |
| M01-U1 | PASS | Drive do hub na stack re-apontada pós-merge: /anuncios (34 anúncios reais), /pedidos, /precos (29 produtos com custo ERP), /integracoes. `read_console_messages(onlyErrors)` = vazio nas 4. |
| M01-U2 | PASS | Diff das DUAS capturas salvas (`_chip-m01/baseline-pre-merge.json` vs `capture-post-merge.json`): 20/20 campos idênticos, inclusive `Content-Length: 472`. Único delta = `tarifa.comissao.data`/`captured_at` (carimbo do momento da chamada). `fonte=COTACAO`, `degrau=3`, `estimativa=false` preservados — o round-trip vivo pelo adapter ML continua atravessando com o decorator no fio. |

## Emendas de contrato adjudicadas

Duas emendas auto-declaradas pelo chip. O hub NÃO aceitou por alegação — mediu o `git diff` no worktree antes do merge. Ambas FORTALECEM cobertura:

1. `pricing_reader_test.go`: removida a linha de tabela `{name: "rate limited", status: 429, want: ErrRateLimited}`, que travava o comportamento PRÉ-fix (nenhum retry). Substituída por `TestPriceToWinRateLimitedRetriesThenExhausts`, que assere o erro tipado `RateLimitExhaustedError`, `Attempts == 2` e a contagem real de chamadas ao servidor fake.
2. `items_multiget_reader_test.go`: o spec pedia cap de 1MB POR ITEM — impossível, porque o 1MB é o teto da resposta inteira dividida por até 20 itens. Implementado 256KiB/item e PINADO por asserção (`itemMultigetRawCap != 256*1024` falha nomeando a constante), de modo que drift futuro quebra alto. O teste de 429 saiu de `calls != 1` para `calls != 2` com budget real.

## Ressalvas registradas

- **tsc do FE vermelho (12 erros)** — pré-existente na main, catalogado no `HARNESS-PROFILE` (inventário de 2026-07-28, CHIP-VINC-NEUTRO). M-01 tocou ZERO arquivo de FE, logo não-atribuível por construção do write-set. Dívida da missão, não deste milestone.
- **Seam ML vivo** não exigido por este milestone (registrado no próprio contrato); certificação transitiva pelos live-drives de M-04/M-06.

## Handoff

- Current status: **passed / merged**.
- Merge: `e79c458e` (`--no-ff`); baseline B-6 em `b35ff536`.
- Next: worktree `laughing-jepsen-9f52ae` a remover; `capability_adapter.go` CONGELADO pós-close (endpoint novo = arquivo novo).
