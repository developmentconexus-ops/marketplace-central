export type ActiveSourceName = "sankhya" | "xlsx" | "catalogo_cliente";

export type SourceKindName = "upload_snapshot" | "live_read_through";

export interface ActiveSourceConfig {
  active_source: ActiveSourceName;
  /** Always server-derived from active_source (tenant_config.DefaultKind) — never client-supplied. */
  source_kind: SourceKindName;
  set_at: string;
  set_by: string | null;
}

/**
 * source_kind is intentionally NOT a field here: the server always derives it
 * from active_source and ignores any client-supplied value (ADR-17).
 */
export interface SetActiveSourceRequest {
  active_source: ActiveSourceName;
  set_by?: string | null;
}

export interface SellableAssortmentConfig {
  only_revenda: boolean;
  only_em_estoque: boolean;
  only_ecommerce_eligible: boolean;
}

export interface SetSellableAssortmentRequest {
  only_revenda: boolean;
  only_em_estoque: boolean;
  only_ecommerce_eligible: boolean;
}

export interface CatalogAssortmentCounts {
  sellable_count: number;
  total_count: number;
}

export interface ActiveSourceError {
  error: "unknown_erp_source" | "invalid_active_source" | "invalid_body" | "internal_error";
  detail?: string;
}
