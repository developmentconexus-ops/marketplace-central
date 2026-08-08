import type { JSX } from "react";
import { MarginChip, UnknownValue } from "@marketplace-central/ui";
import type { PricingCalcProfile, PricingDecomposition } from "@marketplace-central/sdk-runtime";
import { TariffCarimbo } from "./tariffBadge";
import type { TariffBlock } from "./tariffBadge";

export interface DecompositionPanelProps {
  decomposition: PricingDecomposition;
  profile: PricingCalcProfile;
  blockingState: string | null;
  /** Destination UF for the DIFAL line label (from the analysis CEP). */
  difalUf?: string | null;
  /** Tarifa carimbo (comissão + frete provenance) — SolverPanel's badge, reused here. */
  tarifa?: TariffBlock | null;
}

/** A decimal-string value, or UnknownValue (—) when absent — never a fabricated 0. */
function Value({ amount, hint }: { amount: string | null; hint?: string }): JSX.Element {
  if (amount === null) return <UnknownValue hint={hint} />;
  return <span className="font-mono text-ink">{amount}</span>;
}

/**
 * Fiscal components of the D-41 per-cell path. When any of these is listed in
 * componentes_desconhecidos, TaxesForItem did not resolve the item's cell, so
 * imposto=null means "unknown", not "apurado elsewhere". Kept as a literal
 * list rather than "any unknown at all": custo or frete being unknown says
 * nothing about whether the tax was apportioned.
 */
const FISCAL_COMPONENTS = ["icms_saida", "difal", "pis_cofins", "fcp", "restituicao_st"];

function fiscalUnresolved(d: PricingDecomposition): boolean {
  return d.componentes_desconhecidos.some((c) => FISCAL_COMPONENTS.includes(c));
}

function parsePct(value: string | null): number | null {
  if (value === null) return null;
  const n = Number(value);
  return Number.isFinite(n) ? n : null;
}

/**
 * DecompositionPanel renders the per-component IC-04 breakdown of a price:
 * comissão, taxa fixa, frete, imposto, DIFAL, tarifa Full, custo ERP and the
 * resulting retorno. Unknown components (nullable, listed in
 * componentes_desconhecidos) render as "—" and NEVER as 0 (ADR-17); when the
 * cost is unresolved the SEM_CUSTO blocking banner is shown. Margin health is
 * classified by the operator's CalcProfile thresholds, not hard-coded bands.
 */
export function DecompositionPanel({
  decomposition: d,
  profile,
  blockingState,
  difalUf,
  tarifa,
}: DecompositionPanelProps): JSX.Element {
  const thresholds = {
    healthy: Number(profile.limiar_verde_pct),
    tight: Number(profile.limiar_amarelo_pct),
  };

  const semCusto = blockingState === "SEM_CUSTO";

  return (
    <div data-testid="decomposition-panel" className="flex flex-col gap-1 font-mono text-sm">
      {semCusto ? (
        <p
          role="alert"
          data-testid="blocking-sem-custo"
          className="mb-2 rounded-md bg-warn-soft px-3 py-2 text-warn"
        >
          Sem custo do ERP para este produto — a margem não pode ser calculada (SEM_CUSTO).
        </p>
      ) : null}

      {!profile.difal_enabled ? (
        <p
          data-testid="difal-off-warning"
          className="mb-2 rounded-md bg-amber-soft px-3 py-2 text-xs text-amber"
        >
          Margem sem DIFAL — não use para decisão de venda interestadual.
        </p>
      ) : null}

      <Row label="Preço">
        <Value amount={d.preco} />
      </Row>
      <Row label="(−) Comissão">
        <span className="flex items-center gap-1.5">
          <Value amount={d.comissao} />
          <TariffCarimbo comp={tarifa?.comissao} testId="decomp-tarifa-comissao" />
        </span>
      </Row>
      <Row label="(−) Taxa fixa">
        <Value amount={d.taxa_fixa} />
      </Row>
      <Row label="(−) Frete">
        <span className="flex items-center gap-1.5">
          <Value amount={d.frete} hint="frete calculado por peso e CEP" />
          {/* Carimbo only beside a shown value: never a provenance stamp next to "—",
              even if the backend's tarifa.frete disagrees with a null decomposition.frete
              (ADR-17). Comissão needs no such guard — decomposition.comissao is never null. */}
          {d.frete !== null ? (
            <TariffCarimbo comp={tarifa?.frete} testId="decomp-tarifa-frete" />
          ) : null}
        </span>
      </Row>
      {/* imposto=null carries TWO different states and they need two different
          affordances (dual gate, then re-gate — both blocking, in opposite
          directions).

          Decompose sets Imposto=nil purely on ICMSCell != nil, independent of
          whether TaxesForItem resolved anything. So:

          - cell resolved  ⇒ the tax WAS apurado, per fiscal cell; only this
            legacy aggregate field does not carry it. Rendering <UnknownValue/>
            there says "a fact is MISSING" (this panel's affordance for
            componentes_desconhecidos) about a fact that is not missing.
          - cell NOT resolved (ambíguo, CodTrib/AliquotaInterna/UFDestino
            absent) ⇒ icms_saida/difal/pis_cofins land in
            componentes_desconhecidos and margem comes back nil. Saying "por
            célula fiscal" there declares an apportionment that never happened.

          componentes_desconhecidos is the discriminator, and it is the same
          field the rest of the panel already trusts. */}
      <Row label="(−) Imposto">
        {d.imposto === null ? (
          fiscalUnresolved(d) ? (
            <UnknownValue hint="a apuração por célula fiscal deste item não pôde ser resolvida" />
          ) : (
            <span
              data-testid="imposto-por-celula"
              className="text-faint"
              title="apurado por célula fiscal do item; os componentes ainda não são publicados por esta API"
            >
              por célula fiscal
            </span>
          )
        ) : (
          <Value amount={d.imposto} />
        )}
      </Row>
      <Row label={`(−) DIFAL${difalUf ? ` ${difalUf}` : ""}`}>
        <Value amount={d.difal} hint="diferencial de alíquota da UF de destino" />
      </Row>
      <Row label="(−) Tarifa Full">
        <Value amount={d.tarifa_full} />
      </Row>
      <Row label="(−) Custo ERP">
        <Value amount={d.custo} hint="custo do produto vem do ERP" />
      </Row>

      <div className="mt-1 flex items-center justify-between border-t border-border pt-2">
        <span className="text-muted">Retorno /un</span>
        <span className="flex items-center gap-2">
          <Value amount={d.margem_valor} />
          <MarginChip marginPct={parsePct(d.margem_pct)} thresholds={thresholds} />
        </span>
      </div>
    </div>
  );
}

function Row({ label, children }: { label: string; children: JSX.Element }): JSX.Element {
  return (
    <div className="flex items-center justify-between">
      <span className="text-muted">{label}</span>
      {children}
    </div>
  );
}
