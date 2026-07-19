export interface DashboardLastImport {
  at: string;
  age_seconds: number;
}

export type DashboardDegradedSource = "listings" | "linkage" | "orders" | "sync" | "erp_import";

export interface DashboardOverview {
  sync_errors: number | null;
  pending_links: number | null;
  below_margin: number | null;
  missing_gtin: number | null;
  orders_today: number | null;
  orders_7d: number | null;
  anuncios_ativos: number | null;
  last_sync_at: Record<string, string | null> | null;
  degraded: DashboardDegradedSource[];
  last_import: DashboardLastImport | null;
  as_of: string;
}
