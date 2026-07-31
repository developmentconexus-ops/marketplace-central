# ADDENDUM-01 — sellable stock is company AND location scoped, and net of reservation

Date: 2026-07-29 · chip: CHIP-VENDAVEL · status: **CLOSED — `RATIFIED-BY-OPERATOR` + A7 closed by
the hub** (`local_99feb041`, pack/contract amended on main @ `aee1a222`).
Amends: `CONTEXT-PACK.md` rule `only_em_estoque` ("estoque disponível > 0, CODEMP IN (1,2)") and
the contract's VC-2 / VC-5.
Binds: the P2 plan and every implement worker. Fixed decision, not a degree of freedom.

## The ruling

Operator, verbatim (2026-07-29):

> "Uma coisa CODLOCAL é legal 10101 é estoque revenda, 10108 é show room não é pra contar, e
> 10102 é outlet esses que vendem"

Asked whether it binds the sellable FILTER only, or also the stock NUMBER already rendered on
the catalog screen: **"Filtro e número, os dois."**

## Ratified predicate (final authority — do not re-litigate)

```sql
p.ATIVO='S' AND p.CODPROD>0 AND p.USOPROD='R'
AND EXISTS (SELECT 1 FROM METALPRD.TGFEST e
  WHERE e.CODPROD=p.CODPROD AND e.CODPARC=0
    AND e.CODEMP IN (1,2) AND e.CODLOCAL IN (10101,10102)
  GROUP BY e.CODPROD HAVING SUM(NVL(e.ESTOQUE,0)-NVL(e.RESERVADO,0))>0)
```

- **`RESERVADO` is subtracted.** Not cosmetic: in scope there are 871 reserved lines / 37.428
  units, and **354 products go to zero by reservation alone**. The arithmetic stays isolated at
  ONE point per site.
- **`CODPARC = 0` stays**, even though it is redundant inside this cut today (it is the only
  CODPARC with stock there: 4.091 lines / 109.494 units). Third-party stock exists in other
  companies/locations, so the predicate is defensive, not dead. A reviewer asking to drop it
  "because it changes nothing" is answered by this line.
- **Whitelist, never blacklist** — a new internal location created tomorrow would start selling
  by itself under an exclusion rule.

| CODLOCAL | DESCRLOCAL (TGFLOC) | products | sells |
|---|---|---|---|
| 10101 | 1_REVENDA | 8.310 | **yes** |
| 10102 | 2_OUTLET | 595 | **yes** |
| 10108 | 8_SHOW ROMM *(sic in the ERP)* | 1.574 | no |
| 10106 | 6_PENDENCIA PRODUTO FORNECEDOR | 673 | no |
| 10107 | 7_QUEBRAS EM PALETES | 553 | no |
| 10105 | 5_ALMOXARIFADO | 273 | no |
| 10103 | 3_AGUARDANDO CONSERTO | 113 | no |
| 10109 | 9_USO E CONSUMO | 70 | no |

The selling codes land as a named commented constant carrying these names — no location-picker
UI, no per-company UI (YAGNI, out of scope).

## VC-2 acceptance number: **2.923 sellable products**

Locked. Reconciliation of every neighbouring number, so a wrong screen is diagnosable at a
glance:

| number | what it means |
|---|---|
| 3.822 | dead — gross ESTOQUE, no reservation, no location cut |
| 3.508 | forgot the LOCATION cut (585 products enter that must not: show room, almoxarifado, quebras…) |
| 3.277 | correct location cut, forgot the RESERVATION |
| **2.923** | ratified rule |

Decomposition: 10101 = 2.795 · 10102 = 187. The sum exceeds 2.923 because products appear in
both locations — **count DISTINCT CODPROD, never sum per location.**

The 585 delta is an ASSERTION, not a note: in the fixture, the correct filter and the
filter-without-CODLOCAL must produce **different** numbers, otherwise the test does not sustain
the pin.

## The three sites, and the consistency requirement

1. `internal_read/adapters/oracle/catalog_page.go:131-137` — the live `stock` CTE reads
   `est.CODEMP IN (1,2) AND est.CODLOCAL = 10101`; outlet is missing, so today's screen number
   is low for outlet-only products. Becomes `IN (10101, 10102)`. (`price_candidates` at `:146`
   also pins `e.CODLOCAL = 10101` — price selection is OUT of scope; do not "align" it while
   passing through.)
2. `internal_read/adapters/oracle/sync.go:252-256` — Q4 reads `WHERE CODPARC = 0` with no
   company and no location predicate, and sums gross `ESTOQUE - RESERVADO` across everything.
   The mirror is lying on two axes, not one. Keeps `GROUP BY CODPROD, CODLOCAL`, gains
   `CODEMP IN (1,2) AND CODLOCAL IN (10101,10102)`. VC-5's must-fail must fail if **either**
   predicate leaves the query — not only CODEMP.
3. The `only_em_estoque` toggle on both read paths (live Sankhya catalog query and
   MirrorMatcher).

**Consistency is mandatory:** mirror, live catalog and the N/M counter use the SAME definition
of available. If the three sites diverge, VC-2 fails by design.

## Mirror semantics change

`estoque_total` now means **available sellable stock**. Q4's comment at `sync.go:247-251`
asserts `total == sum of children`; filtering Q4 keeps that invariant true while removing
show-room rows from the per-location detail, and the comment gets rewritten to say what the
numbers now mean. Falsehood in a contract gets deleted, not softened. If any consumer of the
mirror depends on the old sense (gross, all companies), that is a `REQUEST` to the hub **before**
the change — not a discovery at QA.

## Provenance

Every figure above: **db-consult MNOS 2026-07-29 via hub** — read-only measurement
(COUNT/aggregation, zero raw rows). Cite it verbatim in the evidence pack.
