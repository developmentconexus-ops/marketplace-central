# Registro de renumeração — ADRs órfãos → numeração global

**Data:** 2026-08-05 · **Base:** `d57ea44b` + restauração de 001/002/003
**Origem:** Onda 0 do plano `2026-08-05-arquitetura-protocolo-de-modulo-plan.md`

## O problema medido

Não existia registro global de ADR. **Cada missão renumerou de ~ADR-01 por conta
própria.** O mesmo número identifica decisões diferentes e não relacionadas conforme a
missão que cita. Nove números (09–17) tinham 1.365+ citações e nenhum documento; por trás
deles há **19 decisões distintas**.

`ADR-017` é a única sem colisão: 1.378 citações, um só sentido em todas as missões.
**Mantém o número 017.**

## Decisão de numeração

1. Uma decisão distinta = um número global. Números novos alocados em 009–022.
2. `017` fica onde está (nenhuma das 1.378 citações precisa mudar).
3. Decisões de **processo de harness** saem da série de arquitetura e viram
   "Decisões de processo de registro" em `docs/HARNESS-PROFILE.md` — foi a mistura de
   arquitetura com processo que produziu a colisão.
4. Esquema de 3 dígitos (`ADR-017`). `ADR-17` e `ADR-04` são grafias mortas.

## Série de arquitetura — `docs/architecture/decisions/`

| Novo | Era | Missão | Decisão | Cites |
|---|---|---|---|---|
| 009 | ADR-09 | MIS-007 | Todo valor de fee carrega proveniência (layer, origem, coletado_em) | 26 |
| 010 | ADR-09 | MIS-004 | Polling/GET only contra ML; dado "vivo" exige refresh visível | 6 |
| 011 | ADR-10 | MIS-007 | Divergences: uma linha aberta por (entity,kind), detectada no ingest | 12 |
| 012 | ADR-10 | MIS-004 | DIFAL com fonte única dentro de `pricing` | 9 |
| 013 | ADR-11 | MIS-007 | Webhook é ponteiro, nunca dado; sempre 200; dedupe emendado | 22 |
| 014 | ADR-11 | MIS-004 | Coleta de referência de mercado é on-demand; runtime docker local | 5 |
| 015 | ADR-12 | MIS-003 | Módulo canônico `listings` é read-only | 15 |
| 016 | ADR-12 | MIS-004 | `sdk-runtime` manual: OpenAPI + SDK no mesmo commit | 12 |
| **017** | **ADR-17** | **todas** | **Unknown nunca vira zero — ausência honesta ponta a ponta** | **1378** |
| 018 | ADR-13 | MIS-003 | Mutation envelope = tabela-protocolo + poller in-process | 10 |
| 019 | ADR-13 | MIS-007 | Ingestão de listings alimenta o observer de snapshot, 1 linha por item | 28 |
| 020 | ADR-14 | MIS-003 | Dado de mercado só via `CollectorPort`; sem scraping; honest-empty | 3 |
| 021 | ADR-15 | MIS-003 | Seam de FE: um shell, TanStack Query exclusivo para server state | 14 |
| 022 | ADR-16 | MIS-007 | Invariante de SKU: `SELLER_SKU == CODPROD` em write de provider | 8 |

## Série de processo — `docs/HARNESS-PROFILE.md`

Não são decisões de arquitetura do produto. Saem da série ADR.

| Ref | Era | Missão | Decisão | Cites |
|---|---|---|---|---|
| P-1 | ADR-09 | MIS-001/003 | Segurança proporcional: writes off por padrão, sem PII/segredo em evidência | 1 |
| P-2 | ADR-10 | MIS-001/003 | Mock nunca alega integração viva | 4 |
| P-3 | ADR-11 | MIS-001/003 | Gate histórico preservado: QA reprovado fica fixo no SHA, nunca reescrito | 4 |
| P-4 | ADR-12 | MIS-007 | Migrações 0086+ pré-alocadas pelo hub em ranges disjuntos | 4 |
| P-5 | ADR-14 | MIS-007 | `root.go` e o par OpenAPI+SDK são seams serializados pelo hub; ≤1 commit de contrato FE em voo | 35 |

## A segunda colisão — `ADR-01`…`ADR-08` de dois dígitos

Medido depois da primeira rodada, e é pior que a dos órfãos. Os documentos de três dígitos
`001`–`008` são uma série antiga (2026-04). As citações de **dois** dígitos no código vêm das
missões MIS-002/004/006/007, que renumeraram de ADR-01 cada uma por conta própria.

| Medida | Valor |
|---|---|
| Citações de dois dígitos (01–08) | **399** |
| Delas, em código vivo (Go/SQL/OpenAPI/FE) | **44** |
| Delas, só em registro de missão `.mnfs` | 355 |
| Decisões distintas por trás | **~30** (4 missões × ~4 números) |
| Números cujo documento de mesmo número **bate** | **0 de 8** |

`ADR-04` sozinho tem 94 citações e quatro sentidos vivos. `ADR-04` no código diz "escritor
único por seam"; o documento `004-integration-catalog-plugin-framework.md` é outro assunto.
**Citação que aponta para um documento existente e errado é pior que citação órfã** — a órfã
avisa que falta algo, esta não.

### Decisão de escopo

Recebem documento **só as decisões que governam código vivo** — 7 delas, 100% das 44
citações de código. As outras ~23 vivem apenas em registro de missão encerrada e ficam
resolvidas por este crosswalk.

| Novo | Era | Missão | Decisão | Cites vivos |
|---|---|---|---|---|
| 024 | ADR-04 | MIS-007 | Escritor único de ingest: `IngestOrder` é o único caminho de escrita | 13 |
| 025 | ADR-03 | MIS-007 | Raw seletivo: payload cru guardado, PII/fiscal nunca | 10 |
| 026 | ADR-07 | MIS-007 | Vocabulário de fase do scheduler `backfill\|incremental\|sweep` (emendado 2026-08-01) | 6 |
| 027 | ADR-06 | MIS-007 | MASS-CLOSURE: ausente de um pull parcial ≠ fechado | 6 |
| 028 | ADR-05 | MIS-007 | Substituição antes de deleção; allowlist de read-guard só encolhe | 4 |
| 029 | ADR-02 | MIS-007 | Decorator de resiliência; opt-out no-retry para writes | 3 |
| 030 | ADR-08 | MIS-007 | Segunda instância do Scheduler por instalação | 2 |

Os documentos `001`–`008` **mantêm seus números**. São uma série real e anterior; quem colidiu
foi a numeração local das missões.

### Não recebem documento (só crosswalk)

Decisões que existiram, foram ratificadas na sua missão, e hoje não têm citação em código
vivo: MIS-002 (portas Oracle paginadas, envelope de cursor, cache TTL + singleflight,
tetos de batch/pool, gate de execução), MIS-004 (base=main pós-merge, fonte ERP dual,
identidade CODPROD/EAN/REFFORN, evidência de preço em market, dono único do adapter ML,
identidade≠veredito, shell FE retheme-first, zero-writes-ML), MIS-006 (`products_mirror` =
estado corrente, porta `SourceKind`/`ProductSourceAdapter`, active-source por tenant,
upsert keep-absent, auto-aprovação por CODPROD+EAN, 3 lacunas de adapter adiadas).
Harvests completos em `adr-0N-twodigit-citations.md`.

## A terceira rodada — as 6 citações vivas que sobraram

O escopo "7 decisões = 100% das 44 citações vivas" da seção anterior era subcontagem: são
**10** decisões vivas, não 7. As 3 restantes foram fechadas depois:

| Novo | Era | Missão | Decisão | Cites vivos |
|---|---|---|---|---|
| 031 | ADR-04 | MIS-006 | Upsert keep-absent no `products_mirror`: rebuild nunca deleta fisicamente | 4 |
| 032 | ADR-05 | MIS-004 | Leitura de catalog-offers do ML atrás de flag que nasce desligada | 2 |
| — | ADR-02 | MIS-006 | `sourcekind` sem dependências → **redirecionada para ADR-023 §5**, sem documento novo | 1 |

## Correções ao plano da Onda 0

Duas afirmações do plano `2026-08-05-arquitetura-protocolo-de-modulo-plan.md` foram medidas
e são falsas. Ficam registradas aqui porque o plano é histórico:

- **§1.1 — "9 números nunca existiram" trata número como decisão.** Por trás dos 9 números
  havia 19 decisões distintas, e a colisão pior (`ADR-01`…`08`, 0 de 8 batendo com o
  documento de mesmo número) o plano nem viu.
- **§1.2 — "`tenant_config` só tem `transport/`" é parcial.** Tem `transport/` como única
  camada, mas também código de domínio e de adapter solto na raiz do módulo
  (`active_source.go`, `context.go`, `repository.go` com SQL cru). Disposição em ADR-023 §8.

## Regra de reescrita de citação

Uma citação `ADR-NN` só pode ser reescrita quando a missão dona do arquivo é conhecida:

- `.mnfs/MIS-00X/**` → a missão é o diretório. Sem ambiguidade.
- Código Go, testes, OpenAPI, FE → resolver pelo **assunto** da citação contra a tabela
  acima, nunca pelo número sozinho. Se o assunto não bater com nenhuma linha, a citação
  vira pendência escrita, não um chute.
- `ADR-17`/`ADR-017` → **não se toca.** 1.378 citações já corretas.
- Grafia de 2 dígitos vira 3 dígitos na mesma passada.

## Contradições registradas (vão escritas nos ADRs que as herdam)

- **ADR-11/MIS-007 (→013):** a tupla de dedupe foi estreitada por emenda para só `_id`;
  citações posteriores largam o qualificador "emendado". A lacuna de flood/retenção sem
  `_id` foi apontada e nunca fechada.
- **ADR-14/MIS-007 (→P-5):** a redação original dizia "≤1 milestone com contrato FE em
  voo", o que contradizia o plano ratificado de 3 lanes paralelas da própria missão. Foi
  emendada para "≤1 **commit** de contrato em voo".
- **ADR-13/MIS-007 (→019):** a PK de `listing_variations` aparece como 5 colunas e como 4
  colunas sob a mesma citação; o próprio corpo se corrigiu depois.
- **ADR-15 (→021):** carve-out nomeado — páginas legadas (Classifications, Marketplaces)
  mantêm fetch próprio até serem refeitas.
- **MIS-004 P3:** as mesmas asserções foram renumeradas entre candidatos de revisor antes
  da reconciliação (ex.: "zero writes ML" saiu de ADR-08 para ADR-09). A numeração já era
  instável dentro de uma única missão.
- **MIS-003:** `DISPATCH-LEDGER.md:19` registra ESCALATION não resolvida de que ADR-12 e
  ADR-17 nunca tiveram documento. Esta onda é a resolução dessa escalação.
