# F-A3 — Idade honesta e coleta de mercado periódica: Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fazer o operador conseguir distinguir dado de 15 minutos de dado de 15 dias em toda tela que mostra idade, e dar à coleta de mercado um job periódico que renove o que envelheceu sem estourar o rate limit do ML.

**Architecture:** Duas fatias independentes. A **Fatia A** conserta a causa raiz da idade invisível numa função compartilhada e propaga para `/mercado`, que hoje descarta o dado. A **Fatia B** liga um scheduler ao serviço de coleta que já existe, colhendo só o que já envelheceu — nunca o catálogo inteiro. A Fatia A entrega valor sozinha e deve ser mergeada antes da B.

**Tech Stack:** React + TypeScript + Vitest + Testing Library (`apps/web`, `packages/ui`, `packages/web-query`, `packages/feature-*`), Go 1.2x (`apps/server_core`), Postgres.

---

## Contexto: por que este plano existe

A decisão do operador em 2026-08-03 foi **as duas coisas, não uma ou outra** — job periódico resolve o dado velho, idade visível resolve a confiança quando o job falha (`docs/superpowers/specs/2026-08-03-mis008-operacao-diaria-design.md:236`).

Ao medir a metade "idade visível", a causa raiz apareceu fora de `market/`:

```ts
// packages/web-query/src/index.ts:104-109
return `dados de ${date.toLocaleTimeString("pt-BR", {
  hour: "2-digit", minute: "2-digit", second: "2-digit", hour12: false,
})}`;
```

`formatAsOf` renderiza **hora do dia**. Sem data, sem idade. Um dado colhido há três dias renderiza `dados de 06:00:00`, byte-idêntico a um colhido agora. O TTL de STALE do domínio é **uma hora** (`apps/server_core/internal/modules/listings/domain/signal.go:26`) — o indicador não consegue expressar nem o próprio limiar que existe para sinalizar.

O comentário em `apps/web/src/pages/ListingDetailPanel.tsx:156` afirma que o indicador existe *"so old numbers read as old"*. Não lê. É prosa falsa sobre comportamento observável.

**O operador escolheu a opção A em 2026-08-03: consertar na raiz.** Colisão foi medida, não suposta — nenhum branch vivo (`worktree-fa1b-reautorizacao`, `p2b-imposto-ex-ante`) toca `market/`, `web-query/`, `ui/`, `precos/` ou `feature-*`.

## Medição (§12.3) — oito respostas com `file:line`

| # | Pergunta | Resposta |
|---|---|---|
| 1 | Quem já faz isso | `formatAsOf` (`packages/web-query/src/index.ts:96`) e **três cópias** de uma escada de idade relativa: `apps/web/src/pages/dashboard/DashboardPage.tsx:48-54`, `apps/web/src/pages/integracoes/SyncHealthCard.tsx:41-51`, `apps/web/src/pages/ListingsRefreshControl.tsx:37-38` |
| 2 | Quem consome | 10 sítios de `FreshnessIndicator`: `AnunciosTable.tsx:122`, `ListingDetailPanel.tsx:204`, `ListingsSummary.tsx:49`, `precos/tariffBadge.tsx:66`, `produto/AnunciosVinculadosTab.tsx:49`, `produto/EstoqueTab.tsx:57`, `produto/ProdutoHeader.tsx:22`, `packages/feature-inventory/src/StockSeguroPage.tsx:332`, `packages/feature-products/src/CatalogPage.tsx:114`, mais o caption em `AnunciosTable.tsx:266` |
| 3 | Qual contrato muda | **Nenhum na Fatia A** — só formatação de tela. A Fatia B não muda contrato HTTP; adiciona um método de porta interna |
| 4 | Produtor real do valor | `MarketPriceIntelAggregate.fetched_at` (`packages/sdk-runtime/src/market.ts:104`, não-nulo) e `ListingSignalEvidence.fetched_at` (`packages/sdk-runtime/src/index.ts:333`). Os dois já chegam ao FE |
| 5 | Teste que pega regressão | **Nenhum pega o defeito.** `AnunciosTable.test.tsx:190` assere `getByLabelText("Data freshness")` — **presença, nunca o valor**. Quatro testes fixam o formato defeituoso: `AnunciosTable.test.tsx:72`, `packages/feature-products/src/CatalogPage.test.tsx:184,189`, `packages/ui/src/FactStates.test.tsx:52,61` |
| 6 | Dívida registrada que toca | Nenhuma. Este plano nomeia D-48…D-51 |
| 7 | O que o provider entrega | Coleta custa até **6 chamadas ML por produto** (`collection_pipeline_service.go:163,194,204,230,280,301`). Bucket é 900/min compartilhado (`connectors/adapters/mercado_livre/resilience_decorator.go:48`) |
| 8 | Local maximum | Idade já é colhida, gravada (`migrations/0053_market_signals_aggregates.sql:12,33`), contratada e entregue ao FE. `buildOppRows` (`apps/web/src/pages/mercado/oportunidades.ts:71-80`) monta a linha sem ela |

## Anti-redundância (§12.4)

| O que a task faz | O que já existe | Por que não serve sozinho |
|---|---|---|
| T1 escada de idade compartilhada | Três cópias idênticas (`DashboardPage.tsx:48`, `SyncHealthCard.tsx:41`, `ListingsRefreshControl.tsx:37`) | Nenhuma é exportada; nenhuma é usada por `formatAsOf`. T2 apaga as cópias em vez de criar uma quarta |
| T3 componente de idade | **Dois** `FreshnessIndicator` (`packages/ui/src/FreshnessIndicator.tsx:3` e `packages/web-query/src/index.ts:168`) | São componentes diferentes com `aria-label` diferentes; os testes ficam presos a qual cópia o arquivo importou |
| T4 idade em `/mercado` | `MarketPriceIntelAggregate.fetched_at` já entregue | Nenhum componente de `/mercado` o lê |
| T6 listar agregado velho | `LatestMarketAggregates` (`market/adapters/postgres/aggregate_repository.go:35`), `LatestMarketAggregatesBySource` (`:79`), `EvidenceReader.Aggregates` (`ports/evidence_reader.go:14`) | **Todos os três exigem a lista de IDs por parâmetro.** Nada enumera. O scheduler precisa perguntar "quem envelheceu", que nenhum deles responde |
| T7 scheduler | `NewProductsScheduler` (`root.go:694`), `NewFeeSyncScheduler` (`root.go:685`) | São schedulers de outros módulos; o padrão é copiado, o job não |

## Quatro defeitos

- **D-48 — `formatAsOf` não expressa idade.** `web-query/src/index.ts:104`. Hora do dia. 15 min e 15 dias renderizam igual em 10 telas.
- **D-49 — dois `FreshnessIndicator`.** `packages/ui/src/FreshnessIndicator.tsx:3` (estilizado, `aria-label="Atualização dos dados"`) e `packages/web-query/src/index.ts:168` (sem estilo, `aria-label="Data freshness"`).
- **D-50 — `/mercado` descarta a idade.** `oportunidades.ts:71-80` e as 3 abas.
- **D-51 — coleta de mercado não tem job periódico.** `root.go:707` cria o serviço; `root.go:718` liga só ao clique do operador. Único módulo com tabela sem job (`design:§5.2`).

## Fake não é evidência

Ordem do operador: *"não quero mock, fallback, validação sempre real para realmente validar, não fazer teste apenas pra passar"*.

- **Quatro testes existentes vão ficar vermelhos na T1, e isso é o vermelho certo** — eles fixam exatamente o defeito. Eles são **reescritos, não apagados**: a intenção (a tela mostra a idade) é preservada e reforçada (a idade passa a ser distinguível). Anote a adjudicação no pack.
- Todo controle negativo abaixo tem must-fail por injeção.
- A Fatia B não fecha com teste de unidade: fecha com o job rodando contra Postgres real e ML real.

## File Structure

**Fatia A**

| Arquivo | Ação | Responsabilidade |
|---|---|---|
| `packages/web-query/src/index.ts` | Modificar | `formatRelativeAge` compartilhado; `formatAsOf` passa a usá-lo; remover o `FreshnessIndicator` duplicado |
| `packages/web-query/src/index.test.ts` | Criar/Modificar | Escada de idade, incluindo o degrau que hoje não existe |
| `packages/ui/src/FreshnessIndicator.tsx` | Modificar | Único componente; ganha `title` com o instante absoluto |
| `apps/web/src/pages/dashboard/DashboardPage.tsx` | Modificar | Apagar a cópia local |
| `apps/web/src/pages/integracoes/SyncHealthCard.tsx` | Modificar | Apagar a cópia local |
| `apps/web/src/pages/ListingsRefreshControl.tsx` | Modificar | Apagar a cópia local |
| `apps/web/src/pages/AnunciosTable.tsx` | Modificar | Importar o componente unificado |
| `apps/web/src/pages/produto/AnunciosVinculadosTab.tsx` | Modificar | Importar o componente unificado |
| `apps/web/src/pages/mercado/oportunidades.ts` | Modificar | `OppRow` carrega `fetchedAt` |
| `apps/web/src/pages/mercado/OportunidadesTable.tsx` | Modificar | Coluna de idade |
| `apps/web/src/pages/mercado/RepricingTable.tsx` | Modificar | Idade por linha via `market_signal.evidence.fetched_at` |
| `apps/web/src/pages/mercado/MonitoradosTab.tsx` | Modificar | Idade por linha |

**Fatia B**

| Arquivo | Ação | Responsabilidade |
|---|---|---|
| `apps/server_core/internal/modules/market/ports/evidence_reader.go` | Modificar | Porta `StaleAggregateProductIDs` |
| `apps/server_core/internal/modules/market/adapters/postgres/aggregate_repository.go` | Modificar | SQL que enumera o que envelheceu |
| `apps/server_core/internal/modules/market/application/collection_scheduler.go` | Criar | Job periódico com teto |
| `apps/server_core/internal/modules/market/application/collection_scheduler_test.go` | Criar | Teto, cadência, erro não para o lote |
| `apps/server_core/internal/composition/root.go` | Modificar | Ligar o scheduler |
| `apps/server_core/tests/integration/market_stale_aggregates_test.go` | Criar | Prova do SQL contra Postgres real |

## Comandos

```bash
cd apps/web && npx vitest run
```

```bash
cd packages/web-query && npx vitest run
```

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/market/... -count=1
```

```bash
pwsh scripts/harness.ps1 -Command integration
```

> **Vermelho pré-existente:** `TestListingsReadContractEndToEnd` (`tests/integration/listings_read_test.go:181`) já falha na main e faz a lane reportar `status=blocked`. Não é seu. Seu teste está verde quando o `failure_token=test=` dele não aparece.

---

# FATIA A — Idade honesta

### Task 1: `formatRelativeAge` compartilhado e `formatAsOf` honesto

**Files:**
- Modify: `packages/web-query/src/index.ts`
- Test: `packages/web-query/src/index.test.ts`

- [ ] **Step 1: Escreva os testes que falham**

Se `packages/web-query/src/index.test.ts` não existir, crie com os imports que os vizinhos usam. Se existir, acrescente ao fim.

```ts
import { describe, expect, it } from "vitest";
import { formatAsOf, formatRelativeAge } from "./index";

// now fixo: sem isto os degraus da escada dependem do relógio da máquina e o
// teste vira flaky em vez de falsificável.
const NOW = new Date("2026-08-03T12:00:00Z").getTime();
const ago = (ms: number) => new Date(NOW - ms).toISOString();

const SEC = 1000;
const MIN = 60 * SEC;
const HOUR = 60 * MIN;
const DAY = 24 * HOUR;

describe("formatRelativeAge", () => {
  it("percorre a escada de idade", () => {
    expect(formatRelativeAge(ago(30 * SEC), NOW)).toBe("há menos de 1 min");
    expect(formatRelativeAge(ago(5 * MIN), NOW)).toBe("há 5 min");
    expect(formatRelativeAge(ago(3 * HOUR), NOW)).toBe("há 3 h");
    expect(formatRelativeAge(ago(5 * DAY), NOW)).toBe("há 5 d");
  });

  // O DEFEITO. Estes dois pares renderizavam a MESMA string antes desta task,
  // porque formatAsOf mostrava só a hora do dia. 15 min e 15 dias ficavam
  // byte-idênticos sempre que caíam no mesmo horário.
  it("distingue idades que caem no mesmo horário do dia", () => {
    expect(formatRelativeAge(ago(15 * MIN), NOW)).not.toBe(
      formatRelativeAge(ago(15 * DAY), NOW),
    );
    expect(formatRelativeAge(ago(2 * HOUR), NOW)).not.toBe(
      formatRelativeAge(ago(2 * HOUR + 3 * DAY), NOW),
    );
  });

  // ADR-17: ausência nunca vira zero nem um instante plausível.
  it("é honesto sobre idade desconhecida", () => {
    expect(formatRelativeAge(null, NOW)).toBe("idade desconhecida");
    expect(formatRelativeAge(undefined, NOW)).toBe("idade desconhecida");
    expect(formatRelativeAge("não é uma data", NOW)).toBe("idade desconhecida");
  });

  // Relógio do servidor adiantado em relação ao do browser produz idade
  // negativa. "há -3 min" seria pior que inútil.
  it("não inventa idade negativa", () => {
    expect(formatRelativeAge(new Date(NOW + 5 * MIN).toISOString(), NOW)).toBe("agora");
  });
});

describe("formatAsOf", () => {
  it("delega para a escada de idade", () => {
    expect(formatAsOf(ago(5 * MIN), NOW)).toBe("há 5 min");
    expect(formatAsOf(ago(5 * DAY), NOW)).toBe("há 5 d");
  });

  it("é honesto sobre idade desconhecida", () => {
    expect(formatAsOf(null, NOW)).toBe("idade desconhecida");
  });
});
```

- [ ] **Step 2: Rode e veja falhar**

```bash
cd packages/web-query && npx vitest run src/index.test.ts
```

Esperado: falha de compilação, `formatRelativeAge is not exported`. Esse é o vermelho da assinatura, não do comportamento — os dois vêm no Step 4.

- [ ] **Step 3: Implemente**

Em `packages/web-query/src/index.ts`, substitua o corpo de `formatAsOf` (`:96-110`) por:

```ts
const SECOND_MS = 1000;
const MINUTE_MS = 60 * SECOND_MS;
const HOUR_MS = 60 * MINUTE_MS;
const DAY_MS = 24 * HOUR_MS;

/**
 * Idade de um fato, em pt-BR, relativa a agora.
 *
 * Substitui a formatação anterior, que renderizava só a hora do dia
 * (`dados de 06:00:00`) e portanto não distinguia um dado de 15 minutos de um
 * de 15 dias — o TTL de STALE do domínio é de UMA HORA
 * (apps/server_core/internal/modules/listings/domain/signal.go:26), então o
 * indicador não conseguia expressar nem o próprio limiar que existe para
 * sinalizar.
 *
 * `now` é parâmetro para o teste poder fixar o relógio; a escada é a mesma que
 * estava copiada em DashboardPage, SyncHealthCard e ListingsRefreshControl,
 * apagadas na Task 2.
 */
export function formatRelativeAge(
  asOf: string | null | undefined,
  now: number = Date.now(),
): string {
  if (!asOf) return "idade desconhecida";
  const date = new Date(asOf);
  if (Number.isNaN(date.getTime())) return "idade desconhecida";

  const elapsed = now - date.getTime();
  // Relógio do servidor à frente do browser. Um fato do futuro não tem idade;
  // "há -3 min" seria pior que inútil (ADR-17: não inventamos o que não sabemos).
  if (elapsed < 0) return "agora";
  if (elapsed < MINUTE_MS) return "há menos de 1 min";
  if (elapsed < HOUR_MS) return `há ${Math.floor(elapsed / MINUTE_MS)} min`;
  if (elapsed < DAY_MS) return `há ${Math.floor(elapsed / HOUR_MS)} h`;
  return `há ${Math.floor(elapsed / DAY_MS)} d`;
}

/**
 * Rótulo de frescor exibido ao lado de um fato. Alias fino sobre
 * formatRelativeAge, mantido porque dez sítios já o importam por este nome.
 */
export function formatAsOf(
  asOf: string | null | undefined,
  now: number = Date.now(),
): string {
  return formatRelativeAge(asOf, now);
}
```

- [ ] **Step 4: Rode e veja passar — e veja o resto do repo ficar vermelho**

```bash
cd packages/web-query && npx vitest run src/index.test.ts
```

Esperado: verde.

```bash
cd apps/web && npx vitest run
```

Esperado: **vermelho**, e os arquivos são exatamente estes:

- `apps/web/src/pages/AnunciosTable.test.tsx:72` — `/^Anúncios, dados de \d{2}:\d{2}:\d{2}$/`
- `packages/feature-products/src/CatalogPage.test.tsx:184,189` — `/^dados de \d{2}:\d{2}:\d{2}$/`
- `packages/ui/src/FactStates.test.tsx:52,61` — `dados de ${expectedTime}` e `"dados de desconhecido"`

**Se algum outro arquivo aparecer, pare e leia antes de mexer** — significa que existe um consumidor que a medição não achou.

- [ ] **Step 5: Reescreva os quatro testes, não os apague**

A intenção de cada um é *"a tela mostra a idade do dado"*. Ela sobrevive e fica mais forte. Trocar a expectativa é reafirmar a intenção contra o comportamento novo; **apagar a asserção não é opção**.

`apps/web/src/pages/AnunciosTable.test.tsx:72`:

```tsx
    // Antes: /^Anúncios, dados de \d{2}:\d{2}:\d{2}$/ — hora do dia, que não
    // distinguia 15 min de 15 dias (D-48). A intenção (o caption datar a
    // tabela) é a mesma; a asserção agora exige idade, que é o que o operador lê.
    expect(caption.textContent).toMatch(/^Anúncios, (agora|há .+)$/);
```

`packages/feature-products/src/CatalogPage.test.tsx:184` e `:189`:

```tsx
    await waitFor(() => expect(indicator.textContent).toMatch(/^(agora|há .+)$/));
```

`packages/ui/src/FactStates.test.tsx:52,61`: leia as duas asserções no contexto delas. A de `:52` compara com um horário calculado — troque por `formatRelativeAge` importado, passando o mesmo `now` do fixture, para o teste não reimplementar a escada. A de `:61` vira:

```tsx
      "idade desconhecida",
```

- [ ] **Step 6: Verde em todo lugar**

```bash
cd apps/web && npx vitest run
```

```bash
cd packages/web-query && npx vitest run
```

- [ ] **Step 7: Must-fail por injeção**

Volte `formatRelativeAge` para o defeito, colando isto logo depois da guarda de `NaN`:

```ts
  return `dados de ${date.toLocaleTimeString("pt-BR", { hour: "2-digit", minute: "2-digit", second: "2-digit", hour12: false })}`;
```

Rode `cd packages/web-query && npx vitest run`. Esperado: `distingue idades que caem no mesmo horário do dia` FALHA. Remova a injeção e confirme verde. **Sem este passo o teste da escada não distingue "mede a idade" de "mede alguma string".**

- [ ] **Step 8: Commit**

```bash
git add packages/web-query/src apps/web/src/pages/AnunciosTable.test.tsx packages/feature-products/src/CatalogPage.test.tsx packages/ui/src/FactStates.test.tsx
git commit -m "fix(web-query): make the freshness label express age, not time of day"
```

---

### Task 2: Apagar as três cópias da escada

**Files:**
- Modify: `apps/web/src/pages/dashboard/DashboardPage.tsx`
- Modify: `apps/web/src/pages/integracoes/SyncHealthCard.tsx`
- Modify: `apps/web/src/pages/ListingsRefreshControl.tsx`

> **Checagem de colisão antes de começar.** `SyncHealthCard.tsx` fica em `apps/web/src/pages/integracoes/`, diretório onde a fatia F-A1b trabalha. Rode:
>
> ```bash
> git log --oneline -5 -- apps/web/src/pages/integracoes/SyncHealthCard.tsx
> ```
>
> A F-A1b **não** lista este arquivo (ela toca `ConnectionHealthCard.tsx` e `IntegracoesPage.tsx`). Se aparecer commit recente dela aqui, pule este arquivo e registre como pendência em vez de disputar o seam.

- [ ] **Step 1: Rode os testes das três telas antes de tocar**

```bash
cd apps/web && npx vitest run src/pages/dashboard src/pages/integracoes/SyncHealthCard.test.tsx src/pages/ListingsRefreshControl.test.tsx
```

Anote quais passam. Verde antes e verde depois é o que prova que a troca foi comportamentalmente neutra. Se algum arquivo de teste não existir, anote — é uma tela sem rede, e isso vira observação no pack, não motivo para pular.

- [ ] **Step 2: Troque as três**

Em `apps/web/src/pages/integracoes/SyncHealthCard.tsx`, apague a função local `formatRelative` (`:41-51`) e importe a compartilhada:

```tsx
import { formatRelativeAge } from "@marketplace-central/web-query";
```

Troque os dois call-sites (`:76` e `:118`) de `formatRelative(x)` para `formatRelativeAge(x)`.

Em `apps/web/src/pages/dashboard/DashboardPage.tsx`, apague a função local que contém a escada de `:48-54` e substitua suas chamadas por `formatRelativeAge`. **Leia a assinatura da função local antes:** ela recebe `ageSeconds` (segundos, não ISO). Converta no call-site em vez de mudar a compartilhada — ou, se o chamador já tiver o ISO à mão, passe o ISO direto.

Em `apps/web/src/pages/ListingsRefreshControl.tsx`, apague a escada de `:37-38` e use `formatRelativeAge`. Note que esta cópia é mais curta (`há 3s` para segundos); a compartilhada diz `há menos de 1 min`. Se algum teste fixar o formato de segundos, **reescreva-o** como na Task 1 Step 5 e anote a mudança de string no pack.

- [ ] **Step 3: Verde**

```bash
cd apps/web && npx vitest run
```

- [ ] **Step 4: Prove que não sobrou cópia**

```bash
cd apps/web && grep -rn "há menos de 1 min" src/ | grep -v node_modules
```

Esperado: **nenhuma linha**. A única ocorrência da escada agora vive em `packages/web-query/src/index.ts`. Se aparecer alguma, é cópia que a medição não achou — apague também.

- [ ] **Step 5: Commit**

```bash
git add apps/web/src/pages/dashboard apps/web/src/pages/integracoes/SyncHealthCard.tsx apps/web/src/pages/ListingsRefreshControl.tsx
git commit -m "refactor(web): collapse three copied age ladders onto the shared one"
```

---

### Task 3: Um só `FreshnessIndicator`

Hoje existem dois componentes com o mesmo nome e `aria-label` diferentes, e cada teste que busca por label está preso a qual cópia o arquivo importou. Sobrevive o de `packages/ui` (estilizado, label em pt-BR).

**Files:**
- Modify: `packages/ui/src/FreshnessIndicator.tsx`
- Modify: `packages/web-query/src/index.ts`
- Modify: `apps/web/src/pages/AnunciosTable.tsx`
- Modify: `apps/web/src/pages/produto/AnunciosVinculadosTab.tsx`

- [ ] **Step 1: Escreva o teste que falha**

Crie ou acrescente a `packages/ui/src/FreshnessIndicator.test.tsx`:

```tsx
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { FreshnessIndicator } from "./FreshnessIndicator";

describe("FreshnessIndicator", () => {
  it("mostra a idade e guarda o instante absoluto no title", () => {
    const iso = new Date(Date.now() - 3 * 60 * 60 * 1000).toISOString();
    render(<FreshnessIndicator asOf={iso} />);

    const el = screen.getByLabelText("Atualização dos dados");
    expect(el).toHaveTextContent("há 3 h");
    // A idade relativa é o que o operador lê de relance; o instante exato tem
    // que continuar alcançável para ele conseguir cruzar com um log.
    expect(el).toHaveAttribute("title", expect.stringContaining(":"));
  });

  it("é honesto quando não há instante", () => {
    render(<FreshnessIndicator asOf={null} />);
    const el = screen.getByLabelText("Atualização dos dados");
    expect(el).toHaveTextContent("idade desconhecida");
    // Sem instante não há title: um tooltip vazio ou "Invalid Date" mentiria.
    expect(el).not.toHaveAttribute("title");
  });
});
```

- [ ] **Step 2: Rode e veja falhar**

```bash
cd packages/ui && npx vitest run src/FreshnessIndicator.test.tsx
```

Esperado: o primeiro falha por não haver `title`.

- [ ] **Step 3: Implemente**

Substitua `packages/ui/src/FreshnessIndicator.tsx` inteiro:

```tsx
import { formatAsOf, formatDateTime } from "@marketplace-central/web-query";

/**
 * Rótulo de frescor de um fato. Único componente de idade do produto — a cópia
 * que vivia em @marketplace-central/web-query foi removida (D-49), porque duas
 * cópias com aria-label diferentes prendiam cada teste a qual delas o arquivo
 * tinha importado.
 *
 * Mostra idade relativa ("há 3 h"), que é o que se lê de relance, e guarda o
 * instante absoluto no title para o operador cruzar com um log.
 */
export function FreshnessIndicator({ asOf }: { asOf: string | null | undefined }) {
  const absolute = formatDateTime(asOf);
  return (
    <span
      className="text-muted text-xs font-mono"
      aria-label="Atualização dos dados"
      {...(absolute === null ? {} : { title: absolute })}
    >
      {formatAsOf(asOf)}
    </span>
  );
}
```

`formatDateTime` já existe em `packages/web-query/src/index.ts:117` e devolve `null` para entrada ausente ou inválida — é por isso que a guarda acima é `=== null` e não um `try`.

Em `packages/web-query/src/index.ts`, **apague** a função `FreshnessIndicator` (`:168-170`). Se `createElement` e `ReactNode` ficarem sem uso, apague os imports também — o `tsc` acusa.

Em `apps/web/src/pages/AnunciosTable.tsx:11`, separe os imports:

```tsx
import { formatAsOf } from "@marketplace-central/web-query";
import { FreshnessIndicator } from "@marketplace-central/ui";
```

Cuidado: o arquivo pode já importar outras coisas de `@marketplace-central/ui`. Junte na linha existente em vez de criar um segundo import do mesmo pacote.

Em `apps/web/src/pages/produto/AnunciosVinculadosTab.tsx:6`, faça o mesmo: `FreshnessIndicator` passa a vir de `@marketplace-central/ui`, e `listingsQueryKeys` continua vindo de `@marketplace-central/web-query`.

- [ ] **Step 4: Conserte os testes que buscavam o label antigo**

```bash
cd apps/web && npx vitest run
```

Esperado: falham as buscas por `getByLabelText("Data freshness")` — entre elas `AnunciosTable.test.tsx:190`. Troque a string para `"Atualização dos dados"`. É renomeação de seletor, não remoção de asserção.

**Enquanto estiver em `AnunciosTable.test.tsx:190`, conserte a asserção vazia.** Ela hoje é:

```tsx
      expect(screen.getByLabelText("Data freshness")).toBeInTheDocument();
```

Presença passa com qualquer string, inclusive com a errada — que é por isso que D-48 sobreviveu a esta rede. Troque por valor:

```tsx
      // Presença sozinha passava mesmo com o formato defeituoso (D-48). O teste
      // se chama "freshness age marker", então ele tem que assertar a IDADE.
      expect(screen.getByLabelText("Atualização dos dados")).toHaveTextContent(/^(agora|há .+)$/);
```

- [ ] **Step 5: Verde e prova de que não sobrou o segundo componente**

```bash
cd apps/web && npx vitest run
```

```bash
grep -rn "FreshnessIndicator" packages/web-query/src/
```

Esperado: **nenhuma linha**.

```bash
grep -rn "Data freshness" apps/web/src packages/
```

Esperado: **nenhuma linha**.

- [ ] **Step 6: Commit**

```bash
git add packages/ui/src packages/web-query/src apps/web/src/pages/AnunciosTable.tsx apps/web/src/pages/AnunciosTable.test.tsx apps/web/src/pages/produto/AnunciosVinculadosTab.tsx
git commit -m "refactor(ui): one FreshnessIndicator, with the absolute instant in the title"
```

---

### Task 4: Idade nas três abas de `/mercado`

**Files:**
- Modify: `apps/web/src/pages/mercado/oportunidades.ts`
- Modify: `apps/web/src/pages/mercado/OportunidadesTable.tsx`
- Modify: `apps/web/src/pages/mercado/RepricingTable.tsx`
- Modify: `apps/web/src/pages/mercado/MonitoradosTab.tsx`
- Test: `apps/web/src/pages/mercado/oportunidades.test.tsx`

- [ ] **Step 1: Escreva os testes que falham**

Acrescente a `apps/web/src/pages/mercado/oportunidades.test.tsx`. O arquivo já tem um fixture de agregado com `fetched_at: "2026-07-19T06:00:00Z"` em `:40` — **reuse-o**.

```tsx
  it("carrega a idade do agregado para a linha", () => {
    const rows = buildOppRows(facts, aggregates, verdicts);
    expect(rows[0].fetchedAt).toBe("2026-07-19T06:00:00Z");
  });

  it("renderiza a idade na linha de Oportunidades", () => {
    const rows = buildOppRows(facts, aggregates, verdicts);
    render(<OportunidadesTable rows={rows} />);
    // O agregado é de 2026-07-19; qualquer "agora" seria idade fabricada.
    expect(screen.getAllByLabelText("Atualização dos dados")[0]).toHaveTextContent(/^há \d+ d$/);
  });
```

Os nomes `facts`, `aggregates`, `verdicts` são os do fixture que já existe no arquivo — **leia o topo e use os que estiverem lá**, não crie um segundo fixture.

- [ ] **Step 2: Rode e veja falhar**

```bash
cd apps/web && npx vitest run src/pages/mercado/oportunidades.test.tsx
```

Esperado: erro de tipo, `Property 'fetchedAt' does not exist on type 'OppRow'`.

- [ ] **Step 3: Implemente**

Em `apps/web/src/pages/mercado/oportunidades.ts`, acrescente ao `OppRow` (depois de `evidenceState`, `:20`):

```ts
  /**
   * Quando a evidência de mercado desta linha foi colhida. Não-nulo no contrato
   * (sdk-runtime/src/market.ts:104) e já entregue ao FE desde sempre — a linha é
   * que o descartava (D-50). Sem ele o operador não distingue mediana de hoje de
   * mediana de três semanas atrás.
   */
  fetchedAt: string;
```

e no `rows.push` (`:71-80`):

```ts
      fetchedAt: agg.fetched_at,
```

Em `apps/web/src/pages/mercado/OportunidadesTable.tsx`, acrescente `"ATUALIZADO"` ao array de cabeçalhos e a célula correspondente na linha:

```tsx
<FreshnessIndicator asOf={r.fetchedAt} />
```

com `import { FreshnessIndicator } from "@marketplace-central/ui";`. **Leia a definição de grid do arquivo antes** — se ele usar um `GRID_COLS` como o `RepricingTable.tsx:11`, acrescente uma faixa (`110px`) na mesma posição da coluna nova, senão o cabeçalho desalinha das células silenciosamente.

Em `apps/web/src/pages/mercado/RepricingTable.tsx`, a idade vem de `r.market_signal?.evidence?.fetched_at ?? null` (o tipo está em `packages/sdk-runtime/src/index.ts:351`). Acrescente `"ATUALIZADO"` ao `HEAD` (`:15-25`), uma faixa `110px` ao `GRID_COLS` (`:11`) **na mesma posição**, e a célula:

```tsx
<FreshnessIndicator asOf={r.market_signal?.evidence?.fetched_at ?? null} />
```

Quando não há sinal, `FreshnessIndicator` já renderiza `idade desconhecida` — não fabrique um instante e não esconda a coluna.

Em `apps/web/src/pages/mercado/MonitoradosTab.tsx`, leia o componente e ache o campo de instante do modelo que ele renderiza. Se **não houver nenhum** campo de instante no que a aba recebe, **pare e registre** como pendência nomeada no pack em vez de inventar um: a aba estaria sem produtor, e isso é ADR-17, não um bug de renderização.

- [ ] **Step 4: Verde**

```bash
cd apps/web && npx vitest run src/pages/mercado/
```

- [ ] **Step 5: Must-fail por injeção**

Troque a célula de Oportunidades por `<FreshnessIndicator asOf={new Date().toISOString()} />`. Rode de novo. Esperado: `renderiza a idade na linha de Oportunidades` FALHA — `há 0 d` não casa `/^há \d+ d$/`… **e se casar, o teste é fraco.** Se passar, aperte a asserção para o número de dias que separa a data do fixture de hoje, e só então restaure.

- [ ] **Step 6: Commit**

```bash
git add apps/web/src/pages/mercado/
git commit -m "feat(mercado): show how old the market evidence on each row is"
```

---

### Task 5: Drive ao vivo da Fatia A

Unidade prova ramo. Só o drive prova que o operador vê a idade certa nos dados reais dele.

- [ ] **Step 1: Confirme que o stack roda o teu código**

```bash
docker ps --format '{{.Names}}\t{{.CreatedAt}}' | grep -E 'frontend|backend'
git log -1 --format='%H %ci'
```

Container mais velho que o último commit → **pare e recrie**, senão o drive lê ausência de observável como ausência de defeito:

```bash
docker compose up -d --force-recreate frontend
```

- [ ] **Step 2: Percorra as telas**

Abra e capture screenshot de cada uma, confirmando que a idade aparece **e é plausível contra o `fetched_at` do banco**:

| tela | onde olhar |
|---|---|
| `/mercado` → Reprecificação | coluna ATUALIZADO por linha |
| `/mercado` → Oportunidades | coluna ATUALIZADO por linha |
| `/anuncios` | marcador de frescor nas linhas STALE |
| `/produto/<codprod>` | header e aba Estoque |
| `/integracoes` | SyncHealthCard |

Cruze uma linha com o banco:

```sql
SELECT product_id, fetched_at, now() - fetched_at AS idade
  FROM market_price_intel_aggregates
 ORDER BY fetched_at DESC LIMIT 5;
```

A idade da tela tem que bater com a coluna `idade`. **Divergência aqui reprova a fatia** — é a diferença entre mostrar idade e mostrar um número.

- [ ] **Step 3: Controle positivo — o dado velho tem que PARECER velho**

Envelheça uma linha e recarregue:

```sql
UPDATE market_price_intel_aggregates
   SET fetched_at = now() - interval '9 days'
 WHERE product_id = '<codprod-de-teste>';
```

Esperado na tela: `há 9 d`. Antes desta fatia, essa mesma linha renderizaria a hora do dia e passaria por fresca. Screenshot dos dois estados. Depois **restaure** o valor original que você anotou no Step 2.

- [ ] **Step 4: Grave a evidência**

Screenshots e saídas de SQL em `.mnfs/MIS-008-operacao-diaria/FA3-idade-honesta/_evidence/`. Não escrito = não aconteceu.

**A Fatia A fecha aqui e pode ser mergeada sozinha.**

---

# FATIA B — Coleta de mercado periódica

> A Fatia A tem que estar mergeada antes desta começar. A idade visível é o instrumento que mostra se o job está funcionando; sem ela, um scheduler quebrado é invisível — que é exatamente a classe de defeito que a F-A1 acabou de fechar em `integrations/`.

### Task 6: Enumerar o que envelheceu

Os três leitores de agregado exigem a lista de IDs por parâmetro (`aggregate_repository.go:35,79`, `ports/evidence_reader.go:14`). Nada enumera. O scheduler precisa perguntar *"quem envelheceu"*.

**Files:**
- Modify: `apps/server_core/internal/modules/market/ports/evidence_reader.go`
- Modify: `apps/server_core/internal/modules/market/adapters/postgres/aggregate_repository.go`
- Create: `apps/server_core/tests/integration/market_stale_aggregates_test.go`

- [ ] **Step 1: Escreva o teste de integração**

Copie a montagem de fixture do vizinho `apps/server_core/internal/modules/market/adapters/postgres/aggregate_repository_integration_test.go`, que já compila e usa este mesmo repositório.

```go
//go:build integration

package integration

// Prova contra Postgres real que a enumeração devolve só o que passou do TTL,
// em ordem de mais velho primeiro, respeitando o teto. Fake nenhum conhece
// ORDER BY nem LIMIT, e é justamente a ordem que decide o que o job colhe
// primeiro quando o teto corta.
func TestStaleAggregateProductIDsReturnsOldestFirstUnderTheCeiling(t *testing.T) {
	// Três agregados: 3 dias, 2 horas, 5 minutos.
	// TTL de 1 hora  -> os dois primeiros são elegíveis, o terceiro não.
	// Teto de 1      -> só o de 3 dias volta.
	//
	// O de 5 minutos é o CONTROLE NEGATIVO: sem ele, uma query sem cláusula de
	// idade devolveria os três e o teste passaria pelo motivo errado.
}
```

Escreva o corpo completo com o fixture do vizinho. Cubra, com asserção separada cada:

1. TTL de 1 hora, teto alto → volta exatamente os IDs de 3 dias e de 2 horas, **nunca** o de 5 minutos.
2. Mesma chamada, teto 1 → volta **só** o de 3 dias (mais velho primeiro).
3. TTL maior que a idade de todos → devolve lista **vazia**, não erro e não `nil` tratado como falha.

Limpe o que criou num `t.Cleanup` que **reporta** os erros com `t.Errorf` — um `_, _ =` silencioso esconderia nome de tabela errado e vazaria linha para toda execução seguinte da lane.

- [ ] **Step 2: Rode a lane e veja falhar**

```bash
pwsh scripts/harness.ps1 -Command integration
```

Esperado: falha de compilação, método inexistente.

- [ ] **Step 3: Implemente**

Em `ports/evidence_reader.go`, acrescente à interface `EvidenceReader`:

```go
	// StaleAggregateProductIDs enumera os codprods cujo agregado mais recente
	// passou de olderThan, mais velho primeiro, no máximo limit.
	//
	// Os outros três leitores (Aggregates, LatestMarketAggregates,
	// LatestMarketAggregatesBySource) exigem a lista de IDs por parâmetro; nenhum
	// responde "quem envelheceu", que é a única pergunta que o job periódico faz.
	StaleAggregateProductIDs(ctx context.Context, olderThan time.Duration, limit int) ([]string, error)
```

Em `aggregate_repository.go`, ao lado de `LatestMarketAggregates` (`:35`):

```go
func (r *Repository) StaleAggregateProductIDs(ctx context.Context, olderThan time.Duration, limit int) ([]string, error) {
	// DISTINCT ON pega a linha mais recente por produto (a tabela é append-only,
	// PK (tenant_id, listing_id, fetched_at) em migrations/0053:13), e o filtro de
	// idade é aplicado DEPOIS, sobre essa linha. Filtrar antes devolveria produto
	// que TEM colheita recente só porque também tem uma velha no histórico.
	const query = `
		SELECT product_id
		  FROM (
		       SELECT DISTINCT ON (product_id) product_id, fetched_at
		         FROM market_price_intel_aggregates
		        WHERE tenant_id = $1
		        ORDER BY product_id, fetched_at DESC
		       ) latest
		 WHERE latest.fetched_at < now() - $2::interval
		 ORDER BY latest.fetched_at ASC
		 LIMIT $3
	`
	rows, err := r.pool.Query(ctx, query, r.tenantID, olderThan.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]string, 0, limit)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
```

**Confirme o nome real da tabela e da coluna de produto antes de rodar** — leia `LatestMarketAggregates` em `:35-78` e use exatamente o que ele usa. Se divergir, corrija a query, não o teste.

- [ ] **Step 4: Verde**

```bash
pwsh scripts/harness.ps1 -Command integration
```

Verde = `failure_token=test=TestStaleAggregateProductIDsReturnsOldestFirstUnderTheCeiling` **não** aparece.

- [ ] **Step 5: Must-fail por injeção**

Apague a linha `WHERE latest.fetched_at < now() - $2::interval`. Rode de novo. Esperado: o teste FALHA, porque o agregado de 5 minutos volta. Restaure e confirme verde.

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/modules/market apps/server_core/tests/integration/market_stale_aggregates_test.go
git commit -m "feat(market): enumerate the aggregates that went stale"
```

---

### Task 7: O scheduler

**Orçamento medido, e ele manda no desenho.** `Collect` faz até 6 chamadas ML por produto (`collection_pipeline_service.go:163,194,204,230,280,301`). O bucket é 900/min **compartilhado com todo o resto** (`resilience_decorator.go:48`). O catálogo vendável tem 2.923 produtos — varrer tudo custaria ~17.500 chamadas e saturaria o bucket por ~20 minutos, matando o sync de anúncios e de pedidos.

Por isso o job colhe **só o que já foi colhido e envelheceu**, com teto por ciclo. Ele renova; não descobre. Descoberta continua sendo o clique do operador, que é o comportamento ratificado (D-F4-p, `useMarketCollection.ts:9-10`).

**Files:**
- Create: `apps/server_core/internal/modules/market/application/collection_scheduler.go`
- Test: `apps/server_core/internal/modules/market/application/collection_scheduler_test.go`

- [ ] **Step 1: Escreva os testes que falham**

```go
package application

// O job renova evidência que envelheceu; ele NÃO varre o catálogo. Os testes
// abaixo fixam as três propriedades que separam as duas coisas.

// 1. O teto por ciclo é respeitado — é o que impede o job de estourar o bucket
//    de 900/min que o sync de anúncios e pedidos também usa.
func TestCollectionSchedulerRespectsThePerCycleCeiling(t *testing.T) {}

// 2. Um produto que falha NÃO aborta o lote. Mesma classe que a F-A2 fechou em
//    integrations/background/refresh_ticker.go:37-51, onde RunOnce retornava no
//    primeiro erro e as sessões restantes do lote nunca eram tentadas.
func TestCollectionSchedulerContinuesAfterOneProductFails(t *testing.T) {}

// 3. Toda falha é LOGADA. Controle negativo do defeito exato da F-A1:
//    `case <-ticker.C: _ = t.RunOnce(ctx)` descartava o erro sem log e a falha
//    ficava invisível em todas as camadas.
func TestCollectionSchedulerLogsEveryFailure(t *testing.T) {}
```

Escreva os três corpos completos. Para o (3), capture com `slog.New(slog.NewJSONHandler(&buf, nil))` e assere sobre `buf.String()` — o mesmo padrão que `apps/server_core/internal/modules/integrations/background/refresh_ticker_test.go` usa; **leia esse arquivo e copie a montagem** em vez de inventar outra.

Injete o relógio e o intervalo por parâmetro. **Não** use `time.Sleep` para esperar tick: o teste chama `RunOnce` direto, e o `Start` é uma casca fina sobre o ticker.

- [ ] **Step 2: Rode e veja falhar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/market/... -count=1
```

- [ ] **Step 3: Implemente**

```go
package application

// CollectionScheduler renova periodicamente a evidência de mercado que passou
// do TTL.
//
// Ele NÃO varre o catálogo. Collect custa até 6 chamadas ML por produto
// (collection_pipeline_service.go:163,194,204,230,280,301) e o bucket é de
// 900/min COMPARTILHADO com o sync de anúncios e de pedidos
// (connectors/adapters/mercado_livre/resilience_decorator.go:48); os 2.923
// produtos vendáveis custariam ~17.500 chamadas e travariam todo o resto.
// Descobrir mercado novo continua sendo o clique do operador (D-F4-p).
type CollectionScheduler struct {
	reader   staleAggregateReader
	collect  productCollector
	ttl      time.Duration
	ceiling  int
	interval time.Duration
	logger   *slog.Logger
}

// RunOnce colhe um ciclo. Erro por produto é logado e o lote CONTINUA — parar
// no primeiro erro foi o defeito que a F-A2 consertou no ticker de refresh
// (integrations/background/refresh_ticker.go:37-51).
func (s *CollectionScheduler) RunOnce(ctx context.Context) error {
	ids, err := s.reader.StaleAggregateProductIDs(ctx, s.ttl, s.ceiling)
	if err != nil {
		s.logger.Error("market collection scheduler: enumeração falhou", "err", err)
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	var failures int
	for _, id := range ids {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if _, err := s.collect.Collect(ctx, id); err != nil {
			failures++
			// Logar SEMPRE. O defeito da F-A1 era exatamente isto descartado:
			// `case <-ticker.C: _ = t.RunOnce(ctx)`, erro sem log, falha
			// invisível em todas as camadas e a tela verde.
			s.logger.Error("market collection scheduler: produto falhou",
				"codprod", id, "err", err)
		}
	}
	s.logger.Info("market collection scheduler: ciclo concluído",
		"colhidos", len(ids)-failures, "falhas", failures, "teto", s.ceiling)
	return nil
}

func (s *CollectionScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.RunOnce(ctx); err != nil {
				s.logger.Error("market collection scheduler: ciclo falhou", "err", err)
			}
		}
	}
}
```

Escreva `staleAggregateReader` e `productCollector` como interfaces locais do pacote (porta do lado do consumidor), e o construtor `NewCollectionScheduler` com todos os campos por parâmetro. `productCollector` casa com `CollectionPipelineService.Collect(ctx, codprod) (CollectionSummary, error)` (`collection_pipeline_service.go:134`).

- [ ] **Step 4: Verde e must-fail**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/market/... -count=1
```

Depois injete os dois defeitos, um de cada vez:

- Troque o `s.logger.Error` de dentro do laço por `_ = err`. Esperado: `TestCollectionSchedulerLogsEveryFailure` FALHA.
- Troque o corpo do `if err != nil` do laço por `return err`. Esperado: `TestCollectionSchedulerContinuesAfterOneProductFails` FALHA.

Restaure os dois e confirme verde. **Sem os dois must-fail, os testes não distinguem "o job roda" de "o job avisa quando quebra".**

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/market/application/collection_scheduler.go apps/server_core/internal/modules/market/application/collection_scheduler_test.go
git commit -m "feat(market): periodic scheduler that renews stale market evidence"
```

---

### Task 8: Ligar o scheduler

**Files:**
- Modify: `apps/server_core/internal/composition/root.go`

- [ ] **Step 1: Ligue**

Logo depois de `markettransport.NewHandlerWithCollections(...).Register(mux)` (`root.go:718`):

```go
	// MIS-006 ensinou que um scheduler construído e não ligado tica um NO-OP e o
	// dado nunca atualiza (o comentário em root.go:686-689 documenta o caso do
	// products scheduler). Este fica ligado aqui, ao lado do handler que usa o
	// MESMO serviço de coleta.
	//
	// TTL de 1h espelha o signalStaleTTL do domínio de listings
	// (listings/domain/signal.go:26): é o mesmo limiar que já marca o sinal como
	// STALE na tela, então o job renova exatamente o que a tela chama de velho.
	// Teto de 50 por ciclo a cada 30min = no máximo ~300 chamadas ML por ciclo
	// (6 por produto), folgado dentro do bucket de 900/min compartilhado.
	go marketapp.NewCollectionScheduler(
		marketModuleRepo, marketCollectionSvc,
		time.Hour, 50, 30*time.Minute, slog.Default(),
	).Start(context.Background())
```

- [ ] **Step 2: Compile e suba**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go build ./...
```

```bash
docker compose up -d --force-recreate backend
```

- [ ] **Step 3: Prove que está de pé DENTRO do container**

Compilar na tua máquina não prova que o binário em execução tem o job. Compare a idade do container com o commit e confirme o código dentro:

```bash
docker ps --format '{{.Names}}\t{{.CreatedAt}}' | grep backend
git log -1 --format='%H %ci'
```

```bash
MSYS_NO_PATHCONV=1 docker exec marketplace-central-backend-1 sh -c 'grep -c NewCollectionScheduler /workspace/apps/server_core/internal/composition/root.go'
```

Esperado: `1`.

- [ ] **Step 4: Commit**

```bash
git add apps/server_core/internal/composition/root.go
git commit -m "feat(market): start the collection scheduler alongside its handler"
```

---

### Task 9: Drive ao vivo da Fatia B

- [ ] **Step 1: Envelheça de propósito e espere um ciclo**

Anote o estado antes:

```sql
SELECT product_id, fetched_at FROM market_price_intel_aggregates
 ORDER BY fetched_at DESC LIMIT 10;
```

Escolha um `product_id` que **já tenha sido colhido com sucesso** (o job renova, não descobre — um produto sem agregado nunca entra) e envelheça:

```sql
UPDATE market_price_intel_aggregates
   SET fetched_at = now() - interval '3 hours'
 WHERE product_id = '<codprod>';
```

- [ ] **Step 2: Confirme o ciclo no log**

```bash
docker compose logs -f backend | grep "market collection scheduler"
```

Esperado, em até 30 minutos: `ciclo concluído colhidos=N falhas=M teto=50`.

- [ ] **Step 3: Confirme que a evidência de fato renovou**

```sql
SELECT product_id, fetched_at, now() - fetched_at AS idade
  FROM market_price_intel_aggregates
 WHERE product_id = '<codprod>' ORDER BY fetched_at DESC LIMIT 3;
```

Esperado: **uma linha nova**, com `idade` de minutos. A tabela é append-only — a linha velha de 3 horas continua lá, e é isso que prova que houve colheita nova em vez de um `UPDATE`.

- [ ] **Step 4: Fecha o laço com a Fatia A**

Abra `/mercado` e confirme que a coluna ATUALIZADO daquele produto mudou de `há 3 h` para `há N min`. **É este passo que prova que as duas fatias se enxergam** — o job renova e a tela conta a verdade sobre isso.

Screenshot dos dois estados.

- [ ] **Step 5: Controle negativo do orçamento de rate limit**

O job não pode faminto o resto. Durante um ciclo, confirme que o sync de anúncios continuou:

```sql
SELECT entity, last_success_at, now() - last_success_at AS idade
  FROM sync_state ORDER BY last_success_at DESC;
```

Esperado: `last_success_at` de listings **avançando** durante e depois do ciclo. Se estagnar, o teto de 50 está alto demais para este tenant — baixe e registre o número medido, não o palpite.

- [ ] **Step 6: Grave a evidência**

`.mnfs/MIS-008-operacao-diaria/FA3-coleta-periodica/_evidence/`.

---

### Task 10: Fechamento

- [ ] **Step 1: Registre as dívidas**

No pack de evidência e no ledger da missão:

> **D-48 — RESOLVIDA** nesta fatia (`formatRelativeAge`).
> **D-49 — RESOLVIDA** nesta fatia (um só `FreshnessIndicator`).
> **D-50 — RESOLVIDA** nesta fatia (idade nas abas de `/mercado`).
> **D-51 — RESOLVIDA** nesta fatia (scheduler de coleta).
> **D-52 — asserção de presença em rede de teste de tela.** `AnunciosTable.test.tsx:190` assertava `toBeInTheDocument()` sobre o marcador de frescor e por isso não pegou D-48 por toda a vida do defeito. Corrigida aqui pontualmente. **A classe não foi varrida:** ninguém mediu quantas outras asserções `toBeInTheDocument()` sobre células de valor existem no repo. Varredura é fatia própria.
> **D-53 — o job renova, não descobre.** Produto sem agregado nunca entra no ciclo periódico; só o clique do operador cria a primeira evidência. Consequência aceita e medida (orçamento de 6 chamadas ML × 2.923 produtos contra bucket de 900/min), não esquecida.

- [ ] **Step 2: Atualize o spec**

Em `docs/superpowers/specs/2026-08-03-mis008-operacao-diaria-design.md`, marque `F-A3` (`:236`) como fechada com o SHA. Corrija a linha `:367`, que declara a Onda 0 disjunta com `F-A3` confinado a `market/` — a medição mostrou que a causa raiz da idade invisível vive em `packages/web-query`, e o operador aprovou alargar em 2026-08-03. **Deixar a linha como está seria prosa falsa sobre o repo.**

- [ ] **Step 3: Commit**

```bash
git add .mnfs/MIS-008-operacao-diaria docs/superpowers/specs/2026-08-03-mis008-operacao-diaria-design.md
git commit -m "docs(F-A3): evidence pack, debts D-48..D-53, close F-A3 in the design"
```

---

## Fora de escopo deste plano

- **Varrer a classe do D-52** (asserções de presença sobre células de valor em toda a rede de testes de tela). Fatia própria.
- **Descoberta periódica de mercado novo** (D-53). Exigiria orçamento de rate limit que o bucket compartilhado não tem hoje; entra quando houver fila com prioridade, não antes.
- **Cadência configurável.** Todo scheduler do produto tem cadência hardcoded (`design:§5.2`). Este segue o padrão de propósito; tornar todos configuráveis é mudança de infraestrutura, não desta fatia.
- **`F-00`** (scheduler periódico de pedidos). Bloqueado por `D-16` (`.mnfs/HARNESS-DEBTS.md:299-311`), plano próprio.
