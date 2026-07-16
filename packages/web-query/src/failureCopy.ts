export type FailureCode =
  | "provider_validation"
  | "provider_rate_limited"
  | "provider_unavailable"
  | "provider_auth"
  | "listing_paused_remote"
  | "link_unresolved"
  | "policy_missing"
  | "sku_invariant_violation"
  | "stale_source"
  | "conflict_remote_changed"
  | "type_not_enabled"
  | "internal";

const failureCopies = {
  provider_validation: "Rejeitado pela validação do marketplace.",
  provider_rate_limited:
    "Limite de requisições do marketplace atingido. Tente novamente em instantes.",
  provider_unavailable: "Marketplace indisponível no momento.",
  provider_auth: "Falha de autenticação com o marketplace. Reconecte a instalação.",
  listing_paused_remote: "Anúncio pausado no marketplace.",
  link_unresolved: "Vínculo de produto não resolvido.",
  policy_missing: "Política de precificação ausente.",
  sku_invariant_violation: "Violação de invariante de SKU.",
  stale_source: "Dados de origem desatualizados.",
  conflict_remote_changed: "Conflito: anúncio alterado no marketplace.",
  type_not_enabled: "Tipo de operação não habilitado para esta instalação.",
  internal: "Erro interno.",
} as const satisfies Record<FailureCode, string>;

export const failureCodes = Object.keys(failureCopies) as readonly FailureCode[];

export function failureCopy(code: string): string {
  if (Object.prototype.hasOwnProperty.call(failureCopies, code)) {
    return failureCopies[code as FailureCode];
  }

  return `Falha desconhecida (${code})`;
}
