import { failureCodes, failureCopy } from "./failureCopy";

const expectedCopies = {
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
} as const;

describe("failureCopy", () => {
  it("returns the exact pt-BR copy for every IC-03 failure code", () => {
    for (const code of failureCodes) {
      const copy = failureCopy(code);

      expect(copy).toBe(expectedCopies[code]);
    }
  });

  it("returns the exact fallback for an unknown code", () => {
    expect(failureCopy("weird_code")).toBe("Falha desconhecida (weird_code)");
  });

  it("exports exactly the 12 known IC-03 failure codes", () => {
    expect(failureCodes).toEqual(Object.keys(expectedCopies));
    expect(failureCodes).toHaveLength(12);
    expect(new Set(failureCodes).size).toBe(12);
  });
});
