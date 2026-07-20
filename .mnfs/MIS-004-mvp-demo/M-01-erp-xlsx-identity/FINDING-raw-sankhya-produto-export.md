# FINDING — raw Sankhya "Produto" export shape (operator's real catalog file)

D-119 2026-07-20. Operator file: `Downloads/PRODUTOS MERCADO LIVRE.xlsx` (1.6 MB).
This is the REAL customer catalog for the demo. The current xlsx parser CANNOT ingest it.

## Actual shape (verified via openpyxl/zip inspection)
- **3 sheets = product categories**: `FACAS` (138), `FERRAMENTAS` (1837), `MAQUINAS` (37) → **2012 produtos**.
  Matches memory `#004-E 2012 produtos`. The parser reads only `GetSheetList()[0]` → drops 1874.
- **Preamble rows** (every sheet, identical layout):
  - row 1: `A1="Produto"` (title banner, rest empty)
  - row 2: `A2="Emissão:17/07/2026 14:47:16"`, `B2="Total de registros:138"`, `C2="Usuário: 63 - SAUL"`
  - row 3: **the real header row**
  - row 4+: data
  The parser assumes `rows[0]` is the header → takes "Produto" → required-column rejection.
- **Header row (row 3), Portuguese Sankhya labels** (14 cols):
  `Código | Descrição | Marca | Marca(dup) | Código de Barra | NCM | Referência do Fornecedor |
   Cód. de Benefício Fiscal na UF | CA | Estoque Mínimo | Cor do Fundo | Cor da Fonte | Ativo | Data de Alteração`
  Parser expects English keys `CODPROD/DESCRPROD/CUSTO/ESTOQUE_FISICO/EAN/REFFORN/MARCA/NCM/GRUPO` → none match.

## Alias map (PT header → NormalizedRow standard key)
- `Código` → CODPROD
- `Descrição` → DESCRPROD
- `Código de Barra` → EAN
- `Referência do Fornecedor` → REFFORN
- `Marca` → MARCA (two "Marca" columns C/D — take first non-empty)
- `NCM` → NCM
- sheet name (FACAS/FERRAMENTAS/MAQUINAS) → GRUPO / DESCRGRUPO (category)
- **No CUSTO, no current ESTOQUE** ("Estoque Mínimo" ≠ physical stock) → lenient path, custo/estoque
  honest-unknown (ADR-17), `source=catalogo_cliente`. This is exactly the two-source case
  (migration 0072). Keep the STRICT path (English CODPROD/CUSTO/ESTOQUE_FISICO for the Sankhya
  cost/stock file #003-E) UNCHANGED — do not regress it.

## Confirmed hub-side normalizer (proof the mapping is complete)
A python normalizer over the 3 sheets (skip preamble, alias PT headers, union, GRUPO=sheet)
produced exactly **138+1837+37 = 2012** standard rows. Operator ELECTED (Q1) NOT to use this
shortcut — the importer must ingest the RAW file natively, imported live through the UI as the
customer would (fidelity). Normalizer kept only as the mapping oracle.

## UI gaps (operator: "botão de configuração da plataforma todo quebrado")
- `POST /erp/imports` (multipart `file`, `source=catalogo_cliente` → lenient) EXISTS and works.
- NO upload control anywhere in the web app (`ImportacaoSection` is read-only history).
- `/integracoes` and `/marketplaces` routes are `WorkspacePlaceholder` stubs (AppRouter.tsx:48,51).
- Operator ELECTED (Q2) a REAL Configuração-da-plataforma screen: drop xlsx + pick source +
  import → POST /erp/imports, show result + history. → CHIP-IMPORT-FIX.
