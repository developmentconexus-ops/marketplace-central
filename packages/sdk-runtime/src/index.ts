import type {
  ActiveSourceConfig,
  CatalogAssortmentCounts,
  SellableAssortmentConfig,
  SetActiveSourceRequest,
  SetSellableAssortmentRequest,
} from "./activeSource";
import type { ErpImportChain, ErpImportCreated, ErpImportDetail, ErpImportList, ErpImportSourceInput } from "./erpImport";

export * from "./erpImport";
export * from "./market";
export * from "./dashboard";
export * from "./activeSource";

export interface CatalogProduct {
  product_id: string;
  sku: string;
  name: string;
  description: string;
  brand_name: string;
  status: string;
  cost_amount: number;
  price_amount: number;
  stock_quantity: number;
  ean: string;
  reference: string;
  taxonomy_node_id: string;
  taxonomy_name: string;
  suggested_price: number | null;
  height_cm: number | null;
  width_cm: number | null;
  length_cm: number | null;
  weight_g: number | null;
}

export interface TaxonomyNode {
  node_id: string;
  name: string;
  level: number;
  level_label: string;
  parent_node_id: string;
  is_active: boolean;
  product_count: number;
}

export interface ProductEnrichment {
  product_id: string;
  tenant_id?: string;
  height_cm: number | null;
  width_cm: number | null;
  length_cm: number | null;
  weight_g: number | null;
  suggested_price_amount: number | null;
}

export interface Classification {
  classification_id: string;
  tenant_id?: string;
  name: string;
  ai_context: string;
  product_ids: string[];
  product_count: number;
  created_at: string;
  updated_at: string;
}

export interface CreateClassificationRequest {
  name: string;
  ai_context: string;
  product_ids: string[];
}

export interface UpdateClassificationRequest {
  name: string;
  ai_context: string;
  product_ids: string[];
}

export interface CredentialField {
  key: string;
  label: string;
  secret: boolean;
}

export type CapabilityStatus = 'supported' | 'unsupported' | 'degraded' | 'blocked'

export interface CapabilityProfile {
  listing_read: CapabilityStatus
  stock_read: CapabilityStatus
  stock_write: CapabilityStatus
  order_read: CapabilityStatus
  messages: CapabilityStatus
  questions: CapabilityStatus
  freight_quotes: CapabilityStatus
  webhooks: CapabilityStatus
  sandbox: CapabilityStatus
}

export interface PluginMetadata {
  icon_url?: string
  color?: string
  docs_url?: string
  rollout_stage: 'v1' | 'wave_2' | 'blocked'
  execution_mode: 'live' | 'blocked'
}

export interface MarketplaceDefinition {
  code: string
  display_name: string
  auth_strategy: 'oauth2' | 'lwa' | 'shopee_partner' | 'api_key' | 'token' | 'unknown'
  is_active: boolean
  capability_profile: CapabilityProfile
  metadata: PluginMetadata
}

export interface IntegrationProviderMetadata {
  country?: string;
  rollout_stage?: "v1" | "wave_2" | "blocked";
  execution_mode?: "available" | "blocked" | "planned";
  unavailable_reason?: string;
  fee_source?: "api_sync" | "seed" | "manual";
  baseline_commission_percent?: number;
  baseline_fixed_fee_amount?: number;
  credential_schema?: CredentialField[];
  docs_url?: string;
}

export interface IntegrationProviderDefinition {
  provider_code: string;
  tenant_id: string;
  family: "marketplace";
  display_name: string;
  auth_strategy: "oauth2" | "lwa" | "shopee_partner" | "api_key" | "token" | "none" | "unknown";
  install_mode: "interactive" | "manual" | "hybrid";
  metadata?: IntegrationProviderMetadata & Record<string, unknown>;
  declared_capabilities: string[];
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface IntegrationInstallation {
  installation_id: string;
  tenant_id: string;
  provider_code: string;
  family: "marketplace";
  display_name: string;
  status:
    | "draft"
    | "pending_connection"
    | "connected"
    | "degraded"
    | "requires_reauth"
    | "disconnected"
    | "suspended"
    | "failed";
  health_status: "healthy" | "warning" | "critical";
  external_account_id: string;
  external_account_name: string;
  active_credential_id?: string;
  last_verified_at?: string;
  connection: IntegrationConnectionSnapshot;
  runtime_capabilities: IntegrationRuntimeCapability[];
  created_at: string;
  updated_at: string;
}

export type CanonicalSourceFactQuality = "current" | "stale" | "unknown" | "conflict";

export interface CanonicalNumericSourceFact {
  source: string;
  value: number | null;
  quality: CanonicalSourceFactQuality;
  observed_at: string | null;
  quality_reason: string | null;
}

/**
 * Closed set of source/identity quality flags a canonical product row may carry.
 * `complete` is the default (no anomaly); `invalid_ean`/`ean_collision` are identity
 * warnings; the rest surface product-existence and per-source gaps from the reader.
 */
export type QualityFlag =
  | "complete"
  | "missing_product"
  | "ambiguous_product"
  | "missing_stock"
  | "missing_price"
  | "missing_cost"
  | "ambiguous_price"
  | "missing_tax"
  | "stale_source"
  | "invalid_ean"
  | "ean_collision";

export interface CanonicalCatalogProduct {
  internal_product_id: number;
  name: string;
  ean: string | null;
  manufacturer_reference: string | null;
  seller_sku: string | null;
  brand_name: string | null;
  product_group_name: string | null;
  ncm: string | null;
  quality_flags: QualityFlag[];
  cost_amount: CanonicalNumericSourceFact;
  price_amount: CanonicalNumericSourceFact;
  stock_quantity: CanonicalNumericSourceFact;
}

export interface CatalogProductFact {
  internal_product_id: number;
  /** @deprecated Compatibility alias of manufacturer_reference. */
  reference: string | null;
  description: string | null;
  ean: string | null;
  manufacturer_reference: string | null;
  brand_name: string | null;
  ncm: string | null;
  quality_flags: QualityFlag[];
  active: boolean;
  sellable_stock: {
    quantity: number | null;
    quality: string[];
  };
  current_price: {
    amount: string | null;
    currency: string;
    quality: string[];
  };
  cost: {
    amount: string | null;
    currency: string;
    quality: string[];
  };
}

export interface CatalogProductFactPage {
  items: CatalogProductFact[];
  next_cursor: string | null;
  page_size: number;
  as_of: string;
}

export interface CatalogPageOptions {
  cursor?: string;
  limit?: number;
  erp_source?: ErpImportSourceInput;
  include_all?: boolean;
}

export interface CatalogProductFactsByIdsOptions {
  /** 1..100 internal_product_ids. Ids the active source does not carry come back absent. */
  ids: number[];
  erp_source?: ErpImportSourceInput;
}

export interface CatalogAssortmentCountsOptions {
  erp_source?: ErpImportSourceInput;
}

export interface CatalogSearchPageOptions {
  q: string;
  /** Opaque cursor from the previous page of matches. Search pages like the list read. */
  cursor?: string;
  limit?: number;
  erp_source?: ErpImportSourceInput;
  include_all?: boolean;
}

/**
 * Error code a catalog page read (list, by-ids, search, counts) answers with on
 * any non-2xx (400 for the invalid_* codes, 503 source_unavailable, 504
 * deadline_exceeded, 500 internal_error). The wire shape is the nested
 * ErrorResponse envelope (apierror.Write), like every other module — see
 * CatalogPageErrorResponse in the OpenAPI spec.
 */
export type CatalogErrorCode =
  | "invalid_cursor"
  | "invalid_limit"
  | "invalid_q"
  | "invalid_erp_source"
  | "invalid_ids"
  | "invalid_include_all"
  | "source_unavailable"
  | "deadline_exceeded"
  | "internal_error";

export type ListingStatus = "active" | "paused" | "closed" | "unknown" | "under_review" | "inactive" | "payment_required" | "not_yet_active";
export type ListingSyncState = "synced" | "error" | "stale" | "queued" | "syncing" | "paused_sync";
export type ListingLinkState = "unresolved" | "conflict" | "resolved" | "rejected";
export type ListingException = "sync_error" | "stale" | "unlinked" | "below_margin" | "sem_vinculo" | "abaixo_custo" | "sem_evidencia";

export interface ListingListOptions {
  installation_id: string;
  q?: string;
  status?: ListingStatus;
  sync_state?: ListingSyncState;
  link_state?: ListingLinkState;
  exception?: ListingException;
  has_exception?: boolean;
  listing_type_code?: string;
  product_id?: string;
  limit?: number;
  cursor?: string;
}

export interface RefreshListingsRequest {
  installation_id: string;
}

export interface RefreshListingsAccepted {
  operation_run_id: string;
}

export interface ListingType {
  code: string;
  label: string;
}

export interface ListingMoney {
  amount: string;
  currency: string;
}

// Listings-owned composite over IC-03 states (not a market enum) — see
// domain invariant in the M-05 slice card: SEM_VINCULO (no codprod link) !=
// NO_PRICE_EVIDENCE (linked, no market evidence) != STALE (evidence older
// than TTL), never collapsed, never fabricated for an unlinked listing.
export type SignalStatus = "OK" | "SEM_VINCULO" | "NO_PRICE_EVIDENCE" | "STALE";

export interface ListingSignalEvidence {
  source: string;
  fetched_at: string;
  freshness: string;
}

export interface ListingMarketSignal {
  status: SignalStatus;
  position: { rank: number; total: number } | null;
  price_to_win: ListingMoney | null;
  delta_pct: string | null;
  match_status: "ACCEPT" | "REVIEW" | "REJECT" | "NO_CANDIDATE" | null;
  n_offers: number | null;
  n_sellers: number | null;
  /** Competitor price range (faixa de mercado), our own offer excluded.
   * null when no valid stats — render "—", never 0 (ADR-17). */
  median: ListingMoney | null;
  min_valid: ListingMoney | null;
  max_valid: ListingMoney | null;
  evidence: ListingSignalEvidence | null;
}

export interface ListingLink {
  state: ListingLinkState;
  product_id: string | null;
  seller_sku: string | null;
}

export interface ListingSyncError {
  code: string;
  message_pt: string;
  message_provider: string | null;
}

export type ListingPendingIssueKind = "sync_error" | "stale" | "unlinked" | "below_margin" | "attribute_required";

export interface ListingPendingIssue {
  kind: ListingPendingIssueKind;
  message_pt: string;
}

export interface ListingICMSWorstCase {
  destination_uf: string;
  worst_case_icms_pct: string | null;
  price_net_basis: string | null;
  below_margin_at_uf: boolean | null;
}

export interface ListingReadModel {
  listing_id: string;
  installation_id: string;
  provider: string;
  provider_listing_id: string;
  /** Provider variation id; null when the listing has no variations. */
  variation_id: string | null;
  title: string;
  listing_type: ListingType;
  status: ListingStatus;
  link: ListingLink;
  price: ListingMoney;
  published_quantity: number;
  sync_state: ListingSyncState;
  sync_error: ListingSyncError | null;
  // quality_score e sales_30d removidos: selo PLANEJADO (ADR-C3). Voltam com
  // o produtor — quality_score no P4, sales_30d quando houver agregador de
  // pedidos por janela.
  pending_issue: ListingPendingIssue | null;
  cost: ListingMoney | null;
  below_margin_worst_case: boolean | null;
  icms_worst_case_by_uf: ListingICMSWorstCase[] | null;
  fetched_at: string | null;
  market_signal?: ListingMarketSignal | null;
  signal_status?: SignalStatus;
}

export interface ListingPage {
  items: ListingReadModel[];
  next_cursor: string | null;
  page_size: number;
  as_of: string;
}

export interface ListingGroup {
  product_id: string | null;
  product_title: string | null;
  listing_count: number;
  group_state: "ok" | "attention" | "error";
  listings: ListingReadModel[];
}

export interface ListingGroupPage {
  groups: ListingGroup[];
  next_cursor: string | null;
  page_size: number;
  as_of: string;
}

export interface ListingTimelineEvent {
  at: string;
  kind: string;
  message_pt: string;
}

export interface ListingDetail extends ListingReadModel {
  timeline: ListingTimelineEvent[];
}

export interface ListingSummaryExceptions {
  sync_error: number;
  stale: number;
  unlinked: number;
  below_margin_worst_case: number | null;
  margin_unknown: number | null;
  sem_vinculo?: number | null;
  abaixo_custo?: number | null;
  sem_evidencia?: number | null;
}

export interface ListingSummary {
  total: number;
  active: number;
  paused: number;
  exceptions: ListingSummaryExceptions;
  as_of: string;
}

export interface CategoryAttributeValue {
  id: string;
  name: string;
}

export interface CategoryAttribute {
  id: string;
  name: string;
  required: boolean;
  value_type: string;
  values: CategoryAttributeValue[] | null;
  constraints: Record<string, unknown> | null;
}

export interface CategoryAttributesResponse {
  category_id: string;
  attributes: CategoryAttribute[];
}

export interface IntegrationConnectionSnapshot {
  state: "draft" | "pending_connection" | "connected" | "degraded" | "needs_reauth" | "disconnected";
  health: "healthy" | "warning" | "critical";
  provider_code: string;
  external_account_id: string;
  external_account_name: string;
  auth_strategy: "oauth2" | "lwa" | "shopee_partner" | "api_key" | "token" | "none" | "unknown";
  last_verified_at?: string | null;
  expires_at?: string | null;
  next_action: "none" | "authorize" | "reauth" | "configure" | "retry";
  reauth_reason?: string;
}

export interface IntegrationRuntimeCapability {
  code: string;
  state: "available" | "unavailable" | "needs_auth" | "degraded" | "not_configured";
  executable: boolean;
  live_validated: boolean;
  local_validated: boolean;
  unavailable_reason?: string;
  last_validated_at?: string | null;
}

export interface IntegrationOperationRun {
  operation_run_id: string;
  installation_id: string;
  operation_type: string;
  status: "queued" | "running" | "succeeded" | "failed" | "cancelled";
  result_code: string;
  failure_code: string;
  translated_error_code?: string;
  attempt_count: number;
  actor_type: string;
  actor_id: string;
  provider_evidence?: Record<string, unknown>;
  duration_ms?: number;
  started_at?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface IntegrationAccountProbe {
  provider_code: string;
  provider_account_id: string;
  provider_account_name: string;
  site_id: string;
  status: string;
  fetched_at: string;
  raw_provider_ref?: unknown;
}

export interface IntegrationListingVariationProbe {
  provider_variation_id: string;
  seller_sku?: string;
  ean?: string;
  available_quantity?: number | null;
}

export interface IntegrationListingProbe {
  provider_code: string;
  provider_item_id: string;
  provider_variation_id?: string;
  provider_status?: string;
  seller_sku?: string;
  ean?: string;
  title?: string;
  available_quantity?: number | null;
  source_updated_at?: string | null;
  fetched_at: string;
  raw_provider_ref?: unknown;
  variations?: IntegrationListingVariationProbe[];
}

export interface IntegrationOrderItemProbe {
  provider_item_id: string;
  provider_variation_id?: string;
  seller_sku?: string;
  ean?: string;
  title?: string;
  quantity: number;
  unit_price?: number | null;
  sale_fee_amount?: number | null;
}

export interface IntegrationOrderPaymentProbe {
  payment_id: string;
  status?: string;
  amount?: number | null;
  transaction_amount?: number | null;
  total_paid_amount?: number | null;
}

export interface IntegrationOrderProbe {
  provider_code: string;
  provider_order_id: string;
  provider_status?: string;
  provider_status_detail?: string;
  provider_created_at?: string | null;
  provider_closed_at?: string | null;
  provider_updated_at?: string | null;
  fetched_at: string;
  raw_provider_ref?: unknown;
  items: IntegrationOrderItemProbe[];
  sale_fee_amount?: number | null;
  payments: IntegrationOrderPaymentProbe[];
  shipping_id?: string;
  cancellation_detail?: string;
  tags?: string[];
}

export type OrderLinkQuality = "resolved" | "rejected" | "conflict" | "unresolved" | "missing";

/** @deprecated per hub lock-exception D-02 */
export interface MarketplaceOrderItem {
  provider_item_id: string;
  provider_variation_id?: string;
  seller_sku?: string;
  title?: string;
  quantity: number;
  unit_price?: number | null;
  sale_fee_amount?: number | null;
  link_quality: OrderLinkQuality;
  internal_product_id?: number | null;
}

/** @deprecated per hub lock-exception D-02 */
export interface MarketplaceOrderPayment {
  provider_payment_id: string;
  provider_status?: string;
  transaction_amount?: number | null;
  total_paid_amount?: number | null;
}

/** @deprecated per hub lock-exception D-02 */
export interface MarketplaceOrder {
  installation_id: string;
  provider_code: string;
  provider_order_id: string;
  provider_status?: string;
  provider_status_detail?: string;
  provider_created_at?: string | null;
  provider_closed_at?: string | null;
  provider_updated_at?: string | null;
  fetched_at: string;
  shipping_id?: string;
  cancellation_detail?: string;
  tags?: string[];
  raw_provider_ref?: unknown;
  items: MarketplaceOrderItem[];
  payments: MarketplaceOrderPayment[];
  created_at?: string;
  updated_at?: string;
}

export interface OrderReadItem {
  provider_item_id: string;
  provider_variation_id?: string;
  seller_sku?: string;
  title?: string;
  quantity: number;
  unit_price?: number;
  sale_fee_amount?: number;
  link_quality: OrderLinkQuality;
  internal_product_id?: number;
  custo_unitario?: number;
  /** The ERP source observed this cost at another instant than the order date
   * (single snapshot, no cost history): the amount is the closest honest
   * observation, taken at custo_observado_em — not an exact as-of answer. */
  custo_aproximado?: boolean;
  custo_observado_em?: string;
}

export interface OrderReadPayment {
  provider_payment_id: string;
  provider_status?: string;
  transaction_amount?: number;
  total_paid_amount?: number;
}

export interface OrderBuyer {
  display: string;
  city?: string;
  uf?: string;
}

export interface OrderSla {
  due?: string;
  atrasado?: boolean;
}

export interface OrderRastreio {
  shipment_id: string;
  status: string;
  // ML shipment sub-status (e.g. out_for_delivery, receiver_absent). Additive/
  // optional: absent when the provider reports no sub-status, never fabricated (ADR-17).
  substatus?: string;
  // Carrier name + tracking URL from the shipment /carrier sub-resource. Additive/
  // optional: absent (carrier 404 / not yet dispatched) is honest absence, never
  // a fabricated carrier (ADR-17).
  transportadora?: string;
  url_rastreio?: string;
}

// OrderFreteReal is the real shipment freight actuals from the ML shipment
// /costs sub-resource — distinct from the modeled decomposicao.frete. Every
// amount is number|null|undefined: absent/null means the cost could not be
// honestly sourced, never a fabricated zero (ADR-17).
export interface OrderFreteReal {
  bruto?: number | null;
  receiver?: number | null;
  sender?: number | null;
}

// OrderBucket is the workflow bucket derived from DeriveOrderBucket
// (server-side single authority). Always present on OrderRead (ADR-17).
export type OrderBucket = "novo" | "faturar" | "enviar" | "enviado" | "cancelado";

// OrderDecomposicao is the order-level cost/fee decomposition (F01-C1).
// Every amount is number|null: null means the component could not be
// honestly sourced (ADR-17). componentes_desconhecidos names every unknown
// source component so a null amount is always explained.
// imposto is the retired D-38 legacy field (P2.b T5) — superseded end to end
// by icms_saida + difal (D-41 per-item path) and always null from this
// endpoint now; kept only for old consumers still reading the key.
// icms_saida/pis_cofins/restituicao_st are the D-41 per-item tax components
// (P2.b T5/T6), summed across the order's items. restituicao_st is the one
// component that ADDS to margem_valor instead of subtracting — an explicit 0
// (e.g. intra-MG destination) is a real fact, distinct from null (unresolved
// product link/fiscal snapshot for at least one item).
export interface OrderDecomposicao {
  comissao: number | null;
  taxa_fixa: number | null;
  frete: number | null;
  imposto: number | null;
  difal: number | null;
  icms_saida: number | null;
  pis_cofins: number | null;
  tarifa_full: number | null;
  custo: number | null;
  restituicao_st: number | null;
  margem_valor: number | null;
  margem_pct: number | null;
  componentes_desconhecidos: string[];
}

// OrderDifal is the order-level DIFAL fact (F01-C1). All fields nullable:
// null means unknown (ADR-17).
export interface OrderDifal {
  amount: number | null;
  uf_route: string | null;
  due_date: string | null;
  paid: boolean | null;
}

// OrderCompradorFiscalEndereco is the buyer billing address. Every field is
// optional: absent (masked/blank) is honest absence, never a fabricated value
// (ADR-17). uf_codigo/uf_nome are the provider's own state code + label,
// rendered verbatim.
export interface OrderCompradorFiscalEndereco {
  logradouro?: string;
  numero?: string;
  cidade?: string;
  uf_codigo?: string;
  uf_nome?: string;
  cep?: string;
  pais?: string;
}

// OrderCompradorFiscal is the buyer fiscal identity for ERP registration,
// resolved read-time via the provider's two-step billing-info flow. Additive/
// optional on the order: absent when the buyer has no billing data (404 = a
// normal "none" condition, ADR-17), never a fabricated identity. doc_tipo is
// opaque — render as-is, never map to a CPF/CNPJ enum. doc_numero is rendered
// by the client only and MUST never be logged (LGPD).
export interface OrderCompradorFiscal {
  nome?: string;
  doc_tipo?: string;
  doc_numero?: string;
  endereco?: OrderCompradorFiscalEndereco;
}

export interface OrderRead {
  provider_order_id: string;
  provider_code: string;
  status: string;
  provider_status_detail: string | null;
  buyer_nickname: string | null;
  total: number | null;
  currency: string;
  fulfillment: string | null;
  cancellation_detail: string | null;
  nf_state: string | null;
  created_at: string | null;
  provider_created_at: string | null;
  provider_closed_at: string | null;
  provider_updated_at: string | null;
  // Operator's "Foi faturado" mark (our-DB only, never an ML write). Optional:
  // absent means not yet invoiced — an honest unknown (ADR-17), never inferred
  // from shipment presence. Drives the Faturar->Enviar workflow bucket.
  faturado_at?: string;
  items: OrderReadItem[];
  payments: OrderReadPayment[];
  vinculo_status?: string;
  buyer?: OrderBuyer;
  componentes_desconhecidos?: string[];
  sla?: OrderSla;
  destino_uf?: string;
  // Shipment destination postal code + receiver name (ML destination.*). Additive/
  // optional: absent (masked/not yet readable — ML obfuscates until payment) is
  // honest absence, never fabricated (ADR-17).
  destino_cep?: string;
  destinatario?: string;
  frete_real?: OrderFreteReal;
  rastreio?: OrderRastreio;
  bucket: OrderBucket;
  retorno_liquido?: number | null;
  margem_pct?: number | null;
  decomposicao: OrderDecomposicao;
  difal: OrderDifal;
  comprador_fiscal?: OrderCompradorFiscal;
}

export interface OrderPage {
  items: OrderRead[];
  next_cursor: string | null;
}

export interface OrderSummary {
  today: number;
  seven_days: number;
  by_status?: {
    novo: number;
    faturar: number;
    enviar: number;
    enviado: number;
  };
}

export interface OrderListOptions {
  installation_id: string;
  limit?: number;
  cursor?: string;
  status?: string;
  date_from?: string;
  date_to?: string;
  q?: string;
}

export interface ImportMarketplaceOrdersResponse {
  installation_id: string;
  imported_count: number;
  skipped_count: number;
  items: MarketplaceOrder[];
}

export interface AssistedSankhyaCandidateLine {
  internal_document_id: number;
  internal_line_number: number;
  product_id: number | null;
  quantity: number | null;
  amount: number | null;
}

export interface AssistedSankhyaCandidate {
  internal_document_id: number;
  document_number: number | null;
  operation_code: 313;
  negotiated_at: string | null;
  amount: number | null;
  lines: AssistedSankhyaCandidateLine[];
}

export interface AssistedSankhyaCandidateListResponse {
  installation_id: string;
  provider_order_id: string;
  candidates: AssistedSankhyaCandidate[];
}

export interface AssistedSankhyaMappingLine {
  mpc_line_id: string;
  internal_document_id: number;
  internal_line_number: number;
}

export interface AssistedSankhyaAudit {
  readonly event_id: string;
  readonly event_type: "confirmed";
  readonly actor_type: "operator_supplied_unverified";
  actor_id: string;
  reason: string;
  source_at: string;
  readonly recorded_at: string;
  readonly configuration_revision: string;
  readonly evidence_state: "unknown" | "exact";
  readonly evidence_reference: string;
}

export interface AssistedSankhyaLinkageResponse {
  installation_id: string;
  provider_order_id: string;
  internal_document_id: number;
  lines: AssistedSankhyaMappingLine[];
  audit: AssistedSankhyaAudit;
}

export interface ConfirmAssistedSankhyaLineSelection {
  mpc_line_id: string;
  internal_document_id: number;
  internal_line_number: number;
}

export interface ConfirmAssistedSankhyaLinkageRequest {
  selected_document_id: number;
  selections: ConfirmAssistedSankhyaLineSelection[];
  actor_id: string;
  reason: string;
  source_at: string;
  idempotency_key: string;
}

export type AssistedSankhyaLineageState = "none" | "partial" | "complete" | "conflict" | "unavailable";

export interface AssistedSankhyaDescendant {
  internal_document_id: number;
  internal_line_number: number;
  attended_quantity: number | null;
}

export interface AssistedSankhyaLineage {
  mpc_line_id: string;
  internal_document_id: number;
  internal_line_number: number;
  state: AssistedSankhyaLineageState;
  descendants: AssistedSankhyaDescendant[];
}

export interface ConfirmAssistedSankhyaLinkageResponse extends AssistedSankhyaLinkageResponse {
  lineage: AssistedSankhyaLineage[];
}

export type ProfitabilityInputScope = "order" | "item";
export type ProfitabilityInputKind =
  | "revenue"
  | "sale_fee"
  | "cost"
  | "tax_icms"
  | "tax_ipi"
  | "tax_pis"
  | "tax_cofins"
  | "freight"
  | "commission";
export type ProfitabilityInputQuality =
  | "complete"
  | "missing"
  | "unresolved_link"
  | "rejected_link"
  | "conflict_link"
  | "stale"
  | "manual"
  | "partial";
export type ProfitabilityManualAdjustmentCategory = "freight" | "commission" | "cost" | "generic_adjustment";

export interface ProfitabilityActor {
  actor_type: string;
  actor_id: string;
  actor_name?: string;
}

export interface ProfitabilityMarginInput {
  installation_id: string;
  provider_order_id: string;
  provider_item_id?: string;
  provider_variation_id?: string;
  scope: ProfitabilityInputScope;
  kind: ProfitabilityInputKind;
  amount?: number | null;
  currency: string;
  source_system: string;
  source_reference?: string;
  observed_at?: string | null;
  quality: ProfitabilityInputQuality;
  quality_reason?: string;
  created_at: string;
  updated_at: string;
}

export interface ProfitabilityManualAdjustment {
  adjustment_id: string;
  idempotency_key: string;
  installation_id: string;
  provider_order_id: string;
  provider_item_id?: string;
  provider_variation_id?: string;
  scope: ProfitabilityInputScope;
  category: ProfitabilityManualAdjustmentCategory;
  amount: number;
  currency: string;
  reason: string;
  actor: ProfitabilityActor;
  created_at: string;
}

export interface ImportProfitabilityMarginInputsResponse {
  installation_id: string;
  imported_count: number;
  items: ProfitabilityMarginInput[];
}

export type OrderRealizationState = "realized" | "not_realized" | "unknown";
export type ProfitSnapshotQuality = "complete" | "incomplete" | "negative_margin" | "not_realized";
export type ProfitSnapshotFlag =
  | "missing_revenue"
  | "missing_sale_fee"
  | "missing_cost"
  | "missing_tax"
  | "missing_freight"
  | "missing_commission"
  | "missing_link"
  | "negative_margin"
  | "order_cancelled"
  | "order_state_unknown";

export interface ProfitabilityProfitSnapshot {
  installation_id: string;
  provider_order_id: string;
  provider_item_id?: string;
  provider_variation_id?: string;
  scope: ProfitabilityInputScope;
  currency: string;
  realization_state: OrderRealizationState;
  revenue_amount?: number | null;
  sale_fee_amount?: number | null;
  cost_amount?: number | null;
  tax_amount?: number | null;
  freight_amount?: number | null;
  commission_amount?: number | null;
  adjustment_amount?: number | null;
  contribution_amount?: number | null;
  margin_percent?: number | null;
  quality: ProfitSnapshotQuality;
  flags?: ProfitSnapshotFlag[];
  created_at: string;
  updated_at: string;
}

export interface CalculateProfitabilityProfitSnapshotsResponse {
  installation_id: string;
  calculated_count: number;
  items: ProfitabilityProfitSnapshot[];
}

export interface IntegrationFeeQuoteSnapshot {
  provider_code: string;
  site_id: string;
  category_id: string;
  listing_type_id: string;
  price_amount: number;
  currency_id: string;
  commission_percent?: number | null;
  fixed_fee_amount?: number | null;
  source_updated_at?: string | null;
  fetched_at: string;
  raw_provider_ref?: unknown;
}

export interface IntegrationStockProbe {
  provider_code: string;
  provider_item_id: string;
  provider_variation_id?: string;
  available_quantity: number;
  provider_status?: string;
  seller_sku?: string;
  ean?: string;
  title?: string;
  source_updated_at?: string | null;
  fetched_at: string;
  raw_provider_ref?: unknown;
  scope: "item" | "variation";
}

export interface ProductLinkListingSnapshot {
  installation_id: string;
  provider_code: string;
  provider_item_id: string;
  provider_variation_id?: string;
  provider_status?: string;
  seller_sku?: string;
  ean?: string;
  title?: string;
  available_quantity?: number | null;
  source_updated_at?: string | null;
  fetched_at: string;
  created_at?: string;
  updated_at?: string;
}

export interface ProductLinkListingSnapshotImportResult {
  installation_id: string;
  imported_count: number;
  items: ProductLinkListingSnapshot[];
}

export type ProductLinkCandidateState =
  | "manual"
  | "exact_sku"
  | "exact_ean"
  | "title_match"
  | "unresolved"
  | "conflict";

export type ProductLinkCandidateMatchInput = "manual" | "seller_sku" | "ean" | "title" | "none";

export type ProductLinkConfidenceBand = "ALTA" | "MEDIA" | "BAIXA";

/** CONFIRM is the confirmation queue (D-121-2): one anchor resolved a single
 * product with nothing against it and nothing corroborating it, so it waits for
 * a human yes. Distinct from REVIEW, which means the anchors disagree or are
 * ambiguous — the two are counted and filtered separately. */
export type ProductLinkMatchStatus =
  | "ACCEPT"
  | "REVIEW"
  | "REJECT"
  | "NO_CANDIDATE"
  | "CONFIRM";

export type ProductLinkReasonDirection = "FOR" | "AGAINST" | "UNAVAILABLE" | "INCOMPARABLE";
export type ProductLinkReasonSide = "provider" | "erp" | "both";

export interface ProductLinkReason {
  anchor: string;
  direction: ProductLinkReasonDirection;
  side?: ProductLinkReasonSide;
  detail?: string;
}

export interface ProductLinkCandidateItem {
  candidate_id: string;
  installation_id: string;
  provider_code: string;
  provider_item_id: string;
  provider_variation_id?: string;
  internal_product_id?: number;
  internal_product_name?: string;
  internal_reference_code?: string;
  state: ProductLinkCandidateState;
  match_input: ProductLinkCandidateMatchInput;
  match_value?: string;
  source_snapshot_fetched_at?: string | null;
  confidence: number;
  confidence_band: ProductLinkConfidenceBand;
  match_status: ProductLinkMatchStatus;
  reasons: ProductLinkReason[];
  created_at: string;
  updated_at: string;
}

export interface ProductLinkCandidateGenerationResult {
  installation_id: string;
  generated_count: number;
  items: ProductLinkCandidateItem[];
}

export type ProductLinkState = "none" | "unresolved" | "conflict" | "resolved" | "rejected";

export type ProductLinkAction = "approve_candidate" | "reject_listing" | "manual_resolve";

export interface ProductLinkActor {
  actor_type: string;
  actor_id: string;
  actor_name?: string;
}

export interface ProductLinkListingIdentity {
  installation_id: string;
  provider_item_id: string;
  provider_variation_id?: string;
}

export interface ProductLink {
  installation_id: string;
  provider_code: string;
  provider_item_id: string;
  provider_variation_id?: string;
  state: ProductLinkState;
  source_candidate_id?: string;
  internal_product_id?: number;
  internal_product_name?: string;
  internal_reference_code?: string;
  created_at: string;
  updated_at: string;
}

export interface ProductLinkAuditEntry {
  audit_id: string;
  installation_id: string;
  provider_code: string;
  provider_item_id: string;
  provider_variation_id?: string;
  action: ProductLinkAction;
  reason?: string;
  source_candidate_id?: string;
  actor: ProductLinkActor;
  previous_state: ProductLinkState;
  next_state: ProductLinkState;
  previous_internal_product_id?: number;
  next_internal_product_id?: number;
  created_at: string;
}

export interface ProductLinkWorkflowItem {
  identity: ProductLinkListingIdentity;
  current_link?: ProductLink;
  candidates: ProductLinkCandidateItem[];
  audit: ProductLinkAuditEntry[];
}

export interface ProductLinkResolutionResult {
  link: ProductLink;
  audit: ProductLinkAuditEntry;
}

export interface ProductLinkBatchApproval {
  candidate_id: string;
}

export interface ProductLinkBatchPreviewItem {
  candidate_id: string;
  status: "OK" | "FAILED";
  cause?: string;
}

export interface PreviewProductLinkBatchResponse {
  items: ProductLinkBatchPreviewItem[];
}

export interface ProductLinkAppliedBatchItem {
  candidate_id: string;
}

export interface ProductLinkFailedBatchItem {
  candidate_id: string;
  cause: string;
}

export interface ApplyProductLinkBatchResponse {
  batch_id: string;
  applied: ProductLinkAppliedBatchItem[];
  failed: ProductLinkFailedBatchItem[];
}

export interface ProductLinkUndoBatchItem {
  audit_id: string;
  candidate_id?: string;
  status: "OK" | "FAILED";
  cause?: string;
}

export interface UndoProductLinkBatchResponse {
  batch_id: string;
  reverted: ProductLinkUndoBatchItem[];
  failed: ProductLinkUndoBatchItem[];
}

export type StockRiskState =
  | "healthy"
  | "oversell"
  | "undersell"
  | "stale"
  | "unresolved"
  | "conflict"
  | "ineligible"
  | "unsupported";

export type StockRiskActionability = "actionable" | "blocked";

export interface InventoryListingIdentity {
  installation_id: string;
  provider_item_id: string;
  provider_variation_id?: string;
}

export interface InventoryBlockingReason {
  code: string;
  message: string;
}

export interface InventoryStockRiskItem {
  identity: InventoryListingIdentity;
  provider_code: string;
  provider_status?: string;
  seller_sku?: string;
  ean?: string;
  title?: string;
  link_state: ProductLinkState;
  internal_product_id?: number;
  internal_product_name?: string;
  internal_reference_code?: string;
  state: StockRiskState;
  actionability: StockRiskActionability;
  actionable: boolean;
  internal_quantity?: number | null;
  provider_quantity?: number | null;
  recommended_quantity?: number | null;
  policy_id: string;
  internal_observed_at?: string | null;
  /**
   * How the internal quantity was produced: a live ERP read, or an uploaded
   * snapshot whose data is only as new as the last import.
   */
  internal_source_kind?: "upload_snapshot" | "live_read_through";
  provider_observed_at?: string | null;
  blocking_reason?: InventoryBlockingReason;
}

export type InventoryStockActionState = "proposed" | "blocked" | "approved" | "applied" | "failed" | "skipped";

export interface InventoryStockActionOperator {
  actor_type: string;
  actor_id: string;
  actor_name?: string;
}

export interface InventoryManualApproval {
  approved: boolean;
  approved_at: string;
  operator: InventoryStockActionOperator;
}

export interface InventoryStockActionAuditEvent {
  state: InventoryStockActionState;
  reason: string;
  provider_status?: "applied" | "rejected" | "transient_failure" | "unsupported_shape";
  occurred_at: string;
}

export interface InventoryStockWriteResult {
  status?: "applied" | "rejected" | "transient_failure" | "unsupported_shape";
  idempotency_key?: string;
  provider_code?: string;
  provider_item_id?: string;
  provider_variation_id?: string;
  message?: string;
  response_summary?: string;
}

export interface InventoryStockActionRecord {
  action_id: string;
  state: InventoryStockActionState;
  trigger: "manual";
  provider_ref: {
    tenant_id: string;
    installation_id: string;
    provider_code: string;
    provider_account_id: string;
    provider_item_id: string;
    provider_variation_id?: string;
  };
  before_quantity?: number | null;
  requested_quantity: number;
  recommended_quantity?: number | null;
  policy_id: string;
  internal_observed_at?: string | null;
  provider_observed_at?: string | null;
  operator: InventoryStockActionOperator;
  reason?: string;
  idempotency_key: string;
  blocking_reason?: InventoryBlockingReason;
  failure_reason?: InventoryBlockingReason;
  provider_result?: InventoryStockWriteResult;
  audit_events: InventoryStockActionAuditEvent[];
  created_at: string;
  updated_at: string;
}

export interface ApplyInventoryStockActionResponse {
  action: InventoryStockActionRecord;
  risk: InventoryStockRiskItem;
}

export interface IntegrationFeeSyncAccepted {
  installation_id: string;
  operation_run_id: string;
  status: "queued";
}

export interface IntegrationAuthorizeResponse {
  installation_id: string;
  provider_code: string;
  state: string;
  auth_url: string;
  expires_in: number;
}

export interface IntegrationAuthStatusResponse {
  installation_id: string;
  status:
    | "draft"
    | "pending_connection"
    | "connected"
    | "degraded"
    | "requires_reauth"
    | "disconnected"
    | "suspended"
    | "failed";
  health_status: "healthy" | "warning" | "critical";
  provider_code?: string;
  external_account_id?: string;
  external_account_name?: string;
  connection: IntegrationConnectionSnapshot;
}

export type SubmitIntegrationCredentialsRequest =
  { metadata?: Record<string, string> } & (
    | { api_key: string; credentials?: Record<string, string> }
    | { api_key?: string; credentials: Record<string, string> }
  );

export interface CreateIntegrationInstallationRequest {
  installation_id: string;
  provider_code: string;
  family: "marketplace";
  display_name: string;
}

export interface MarketplaceFeeSchedule {
  id: string;
  marketplace_code: string;
  category_id: string;
  listing_type: string;
  commission_percent: number;
  fixed_fee_amount: number;
  notes: string;
  source: 'api_sync' | 'seeded' | 'manual';
  synced_at: string;
}

export interface MarketplaceAccount {
  account_id: string;
  tenant_id: string;
  channel_code: string;
  marketplace_code: string;
  display_name: string;
  status: string;
  connection_mode: string;
}

export interface CreateMarketplaceAccountRequest {
  account_id: string;
  channel_code: string;
  display_name: string;
  connection_mode: string;
  marketplace_code?: string;
  credentials_json?: Record<string, string>;
}

export type ShippingProvider = "fixed" | "melhor_envio" | "marketplace";

export interface MarketplacePolicy {
  policy_id: string;
  tenant_id: string;
  account_id: string;
  marketplace_code: string;
  commission_percent: number;
  commission_override: number | null;
  fixed_fee_amount: number;
  default_shipping: number;
  tax_percent: number;
  min_margin_percent: number;
  sla_question_minutes: number;
  sla_dispatch_hours: number;
  shipping_provider: string;
}

export interface CreateMarketplacePolicyRequest {
  policy_id: string;
  account_id: string;
  commission_percent: number;
  commission_override?: number | null;
  fixed_fee_amount: number;
  default_shipping: number;
  min_margin_percent: number;
  sla_question_minutes: number;
  sla_dispatch_hours: number;
  shipping_provider?: string;
}

export interface PricingSimulation {
  simulation_id: string;
  tenant_id: string;
  product_id: string;
  account_id: string;
  margin_amount: number;
  margin_percent: number;
  status: string;
}

export interface RunPricingSimulationRequest {
  simulation_id: string;
  product_id: string;
  account_id: string;
  base_price_amount: number;
  cost_amount: number;
  commission_percent: number;
  fixed_fee_amount: number;
  shipping_amount: number;
  min_margin_percent: number;
}

export interface BatchSimulationRequest {
  product_ids: string[];
  policy_ids: string[];
  origin_cep: string;
  destination_cep: string;
  price_source: "my_price" | "suggested_price";
  price_overrides?: Record<string, number>;
}

export interface BatchSimulationItem {
  product_id: string;
  policy_id: string;
  selling_price: number;
  cost_amount: number;
  commission_amount: number;
  freight_amount: number;
  fixed_fee_amount: number;
  margin_amount: number;
  margin_percent: number;
  status: string;
  commission_source: string;
  freight_source: string;
}

// --- IC-04 pricing calculator (M-07) ---
// Money and percentage fields are decimal strings (never float), and unknown
// components are null (ADR-17: unknown ≠ zero), never a fabricated 0.
export interface PricingCalcProfile {
  regime: string;
  aliquota_pct: string;
  limiar_verde_pct: string;
  limiar_amarelo_pct: string;
  tarifa_full: string | null;
  difal_enabled: boolean;
  difal_destino_uf: string | null;
  origem: string;
}

export interface PricingDifalOverride {
  interna_pct: string;
  updated_at: string;
  actor: string;
}

export interface PricingDifalRate {
  uf: string;
  interna_pct: string;
  interestadual_pct: string;
  efetivo_pct: string;
  origem_versao: string;
  override: PricingDifalOverride | null;
}

export interface PricingDifalListResponse {
  items: PricingDifalRate[];
  disclaimer: string;
}

export interface PricingDifalOverrideRequest {
  interna_pct: string;
  actor: string;
}

export interface PricingDifalOverrideResponse {
  persisted: boolean;
  rate: PricingDifalRate;
}

export interface PricingScenario {
  id: string;
  name: string;
  payload: Record<string, unknown>;
  created_at: string;
}

export interface PricingScenarioListResponse {
  items: PricingScenario[];
}

export interface PricingScenarioRequest {
  id: string;
  name: string;
  payload: Record<string, unknown>;
}

export interface PricingCalcInput {
  preco?: string;
  margem_alvo_pct?: string;
  comissao_pct?: string;
  modalidade: string;
  tarifa_full?: string | null;
  frete_produto?: string | null;
  custo?: string | null;
  product_id?: number | null;
}

export interface PricingDecomposition {
  preco: string;
  comissao: string;
  taxa_fixa: string;
  frete: string | null;
  /**
   * Legacy D-38 regime-aliquota field, in retirement (Task A7). Always
   * present, never optional. null means this item's tax was apportioned per
   * fiscal cell, so no aggregate legacy figure exists for it — a STRUCTURAL
   * absence, not an unknown fact (do not render it with the affordance used
   * for componentes_desconhecidos). Never a fabricated "0.00"/"". The
   * per-component successor fields are not on this interface yet.
   */
  imposto: string | null;
  difal: string | null;
  tarifa_full: string | null;
  custo: string | null;
  margem_valor: string | null;
  margem_pct: string | null;
  componentes_desconhecidos: string[];
}

export interface PricingTarifaComponent {
  valor: string | null;
  fonte: string;
  degrau: number;
  data: string | null;
  estimativa: boolean;
}

export interface PricingTarifaFrete extends PricingTarifaComponent {
  sem_dados: boolean;
}

export interface PricingTarifa {
  comissao: PricingTarifaComponent;
  frete: PricingTarifaFrete;
}

export interface PricingTariffDefaults {
  comissao_classico_pct: string;
  comissao_premium_pct: string;
  frete_estimativa_amount: string | null;
  frete_policy: "estimativa" | "sem_dados";
}

export interface PricingDecomposeResponse {
  decomposition: PricingDecomposition;
  blocking_state: string | null;
  tarifa?: PricingTarifa | null;
}

export interface PricingSolveResponse {
  reached: boolean;
  preco: string | null;
  ceiling_pct: string;
  desconhecidos: string[];
  frete_desconhecido: boolean;
  blocking_state: string | null;
  tarifa?: PricingTarifa | null;
  code?: string;
}

export type MutationType =
  | "price_update"
  | "stock_correct"
  | "link_apply"
  | "listing_pause"
  | "listing_resync"
  | "listing_edit"
  | "listing_create";

export type MutationState =
  | "draft"
  | "previewed"
  | "approved"
  | "applying"
  | "applied"
  | "partially_failed"
  | "failed_preserved"
  | "cancelled";

export type MutationItemState = "previewed" | "approved" | "applying" | "applied" | "failed" | "skipped";

export interface CreateMutationRequest {
  installation_id: string;
  type: MutationType;
  actor: string;
  intent: Record<string, unknown>;
  selection: Record<string, unknown>;
}

export interface ApproveMutationRequest {
  execute: true;
}

export interface MutationProtocol {
  protocol_id: string;
  installation_id: string;
  type: MutationType;
  state: MutationState;
  actor: string;
  intent: Record<string, unknown>;
  selection: Record<string, unknown>;
  totals: Record<string, unknown>;
  source_as_of: string | null;
  retried_from: string | null;
  created_at: string;
  previewed_at: string | null;
  approved_at: string | null;
  finished_at: string | null;
}

export interface MutationItem {
  seq: number;
  item_id: string;
  listing_id: string;
  idempotency_key: string;
  before: Record<string, unknown> | null;
  after: Record<string, unknown>;
  state: MutationItemState;
  failure: Record<string, unknown> | null;
  applied_at: string | null;
}

export interface MutationPreview extends MutationProtocol {
  items: MutationItem[];
}

export interface MutationPage<T> {
  items: T[];
  next_cursor: string | null;
  page_size: number;
}

export interface MutationListOptions {
  installation_id: string;
  state?: MutationState;
  type?: MutationType;
  limit?: number;
  cursor?: string;
}

export interface MutationItemListOptions {
  state?: MutationItemState;
  limit?: number;
  cursor?: string;
}

export interface ListResponse<T> {
  items: T[];
}

export interface ErrorResponse {
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}

// Per-domain error code unions, sourced verbatim from the `code` enum on each
// per-domain error schema in contracts/api/marketplace-central.openapi.yaml.
// ApiErrorCode is the union of all of them and is what hasCode() checks
// against — a bogus code (typo, code from a domain not wired here) is a
// compile error at the call site (see errorContract.test.ts).
export type ConnectorErrorCode =
  | "CONNECTORS_PUBLISH_MISSING_NAME"
  | "CONNECTORS_PUBLISH_MISSING_SKU"
  | "CONNECTORS_PUBLISH_MISSING_CATEGORY"
  | "CONNECTORS_PUBLISH_MISSING_BRAND"
  | "CONNECTORS_PUBLISH_MISSING_PRICE"
  | "CONNECTORS_PUBLISH_MISSING_STOCK"
  | "CONNECTORS_PUBLISH_ALREADY_IN_PROGRESS"
  | "CONNECTORS_PUBLISH_DEPENDENCY_FAILED"
  | "CONNECTORS_METHOD_NOT_ALLOWED"
  | "CONNECTORS_ME_METHOD_NOT_ALLOWED"
  | "CONNECTORS_ME_STATE_GENERATION_FAILED"
  | "CONNECTORS_ME_STATE_MISSING"
  | "CONNECTORS_ME_STATE_MISMATCH"
  | "CONNECTORS_ME_CODE_MISSING"
  | "CONNECTORS_ME_TOKEN_EXCHANGE_FAILED"
  | "CONNECTORS_ME_TOKEN_SAVE_FAILED"
  | "CONNECTORS_ME_STATUS_STORE_FAILED"
  | "CONNECTORS_EXECUTOR_INTERNAL"
  | "CONNECTORS_RETRY_MISSING_PRODUCT"
  | "CONNECTORS_BATCH_NOT_FOUND";

export type MutationErrorCode =
  | "invalid_body"
  | "invalid_intent"
  | "installation_required"
  | "invalid_filter"
  | "invalid_cursor"
  | "actor_required"
  | "type_not_enabled"
  | "empty_selection"
  | "selection_too_large"
  | "execute_required"
  | "source_time_unavailable"
  | "preview_stale"
  | "invalid_state"
  | "nothing_to_retry"
  | "protocol_not_found"
  | "method_not_allowed";

// Merges ErpImportError (400/422/404/500) and ErpImportConflict (409) — both
// answer the same /erp/imports* surface and are one domain from the client's
// point of view.
export type ErpImportErrorCode =
  | "invalid_file"
  | "missing_required_column"
  | "invalid_import_id"
  | "import_not_found"
  | "internal_error"
  | "duplicate_file"
  | "import_in_progress";

export type ActiveSourceErrorCode =
  | "unknown_erp_source"
  | "invalid_active_source"
  | "invalid_body"
  | "internal_error";

export type SellableAssortmentErrorCode =
  | "unknown_erp_source"
  | "invalid_body"
  | "internal_error";

export type ApiErrorCode =
  | CatalogErrorCode
  | ConnectorErrorCode
  | MutationErrorCode
  | ErpImportErrorCode
  | ActiveSourceErrorCode
  | SellableAssortmentErrorCode;

/**
 * The one typed throw path for every SDK client (VC-5). A real `Error`
 * subclass — `instanceof` works, a bare `console.error(e)` prints a human
 * message, and `.error` stays available for call sites that already read
 * `e.error.code`. `details` is always an object, never `undefined`: an
 * unknown operational fact (a malformed body, a proxy's HTML error page)
 * lands as code `internal_error` with the original payload preserved at
 * `details.raw` (ADR-17) rather than being silently dropped.
 */
export class MarketplaceCentralClientError extends Error {
  readonly status: number;
  readonly code: string;
  readonly details: Record<string, unknown>;
  readonly error: { code: string; message: string; details: Record<string, unknown> };

  constructor(status: number, code: string, message: string, details: Record<string, unknown>) {
    super(message);
    this.name = "MarketplaceCentralClientError";
    this.status = status;
    this.code = code;
    this.details = details;
    this.error = { code, message, details };
  }
}

export function isApiError(value: unknown): value is MarketplaceCentralClientError {
  return value instanceof MarketplaceCentralClientError;
}

/**
 * `details` fields the spec actually guarantees, keyed by the code that
 * guarantees them — sourced from CatalogPageErrorResponse (allowed_range,
 * populated only for invalid_erp_source/invalid_limit/invalid_ids — see
 * catalog/transport/http_handler.go:419-426, "gating this on one hardcoded
 * code made allowed_range dead data on the two other sites"), ErpImportError
 * (column, populated only for missing_required_column — erp_import's
 * writeImportError, CodeMissingRequiredColumn branch) and ErpImportConflict
 * (import_id + protocol, populated only for duplicate_file — the
 * DuplicateFileError branch of the same function).
 *
 * import_in_progress is deliberately ABSENT: its branch in writeImportError
 * passes `nil` for details, so there is nothing to guarantee there — adding
 * it here would be the SDK asserting a field the server never sends.
 *
 * A flat schema (`details: { additionalProperties: true }`) cannot express
 * "this field exists only for these 3 of 9 codes", so this map plus the
 * runtime check in hasCode() below is what makes the per-code narrowing
 * honest instead of asserted.
 */
export interface ApiErrorDetailsByCode {
  invalid_erp_source: { allowed_range: string };
  invalid_limit: { allowed_range: string };
  invalid_ids: { allowed_range: string };
  missing_required_column: { column: string };
  duplicate_file: { import_id: string; protocol: string };
}

// Keyed off ApiErrorDetailsByCode rather than off `string`, so tsc enforces
// that the runtime check and the type-level claim cannot drift apart. Both
// directions are caught: a code added to the interface with no entry here is
// a missing-property error, an entry here for a code not in the interface is
// an excess-property error, and a misspelled field name is rejected because
// the element type is `keyof` the mapped shape rather than plain `string`.
// The drift this forbids is the whole defect class this chip exists to kill:
// the type promising `string` while nothing actually validates it.
type RequiredDetailFields = {
  [C in keyof ApiErrorDetailsByCode]: readonly (keyof ApiErrorDetailsByCode[C] & string)[];
};

const REQUIRED_DETAIL_FIELDS: RequiredDetailFields = {
  invalid_erp_source: ["allowed_range"],
  invalid_limit: ["allowed_range"],
  invalid_ids: ["allowed_range"],
  missing_required_column: ["column"],
  duplicate_file: ["import_id", "protocol"],
};

// details is narrowed to the mapped shape (intersected with the general
// Record so callers can still read an unmapped field) when C is a code the
// spec guarantees fields for; otherwise it stays the general shape.
export type ApiErrorDetails<C extends ApiErrorCode> = C extends keyof ApiErrorDetailsByCode
  ? ApiErrorDetailsByCode[C] & Record<string, unknown>
  : Record<string, unknown>;

// Shared "is this field present and a string" idiom — also used by
// isRefreshInProgressError below, rather than growing a second copy of it.
function hasStringField(details: Record<string, unknown>, field: string): boolean {
  return typeof details[field] === "string";
}

/**
 * The code check alone is not the contract: for codes the spec guarantees
 * fields for (see ApiErrorDetailsByCode), hasCode also VALIDATES those
 * fields are present and are strings before narrowing. A code match with a
 * missing or wrong-typed field returns false, so a consumer falls through to
 * its generic branch instead of reading `undefined` while the type claims
 * `string` — the type-level narrowing above is only honest because this
 * runtime check backs it.
 */
export function hasCode<C extends ApiErrorCode>(
  value: unknown,
  code: C,
): value is MarketplaceCentralClientError & { code: C; details: ApiErrorDetails<C> } {
  if (!isApiError(value) || value.code !== code) return false;
  // Widening cast, not a narrowing one: most ApiErrorCode members are not keys
  // of the map, so the lookup is legitimately a miss for them and `undefined`
  // is the honest result. The cast only loosens the key type for the read; it
  // cannot weaken the completeness the RequiredDetailFields literal enforces.
  const requiredFields = (REQUIRED_DETAIL_FIELDS as Record<string, readonly string[] | undefined>)[code];
  if (requiredFields === undefined) return true;
  return requiredFields.every((field) => hasStringField(value.details, field));
}

/**
 * POST /listings/refresh's 409 is not backed by a per-domain enum in the
 * spec (it answers the generic ErrorResponse — see
 * apps/server_core/internal/modules/listings/transport/http_handler.go:66),
 * so "refresh_in_progress" is not a member of ApiErrorCode and hasCode()
 * cannot narrow it. A dedicated type-level guard gets the same narrowing
 * power hasCode gives ApiErrorCode members, without inventing a code that
 * isn't in the spec's closed unions.
 */
export function isRefreshInProgressError(
  value: unknown,
): value is MarketplaceCentralClientError & {
  status: 409;
  code: "refresh_in_progress";
  details: { operation_run_id: string } & Record<string, unknown>;
} {
  return (
    isApiError(value) &&
    value.status === 409 &&
    value.code === "refresh_in_progress" &&
    hasStringField(value.details, "operation_run_id")
  );
}

function isErrorEnvelope(
  data: unknown,
): data is { error: { code: string; message?: unknown; details?: unknown } } {
  if (typeof data !== "object" || data === null) return false;
  const error = (data as Record<string, unknown>).error;
  if (typeof error !== "object" || error === null) return false;
  return typeof (error as Record<string, unknown>).code === "string";
}

/**
 * The one parse path every throw site funnels into (VC-5's core). Validates
 * rather than casts: a well-formed envelope (`error` is an object with a
 * string `code`) is used as-is, with `details` defaulting to `{}` when
 * absent. Anything else — a bare string, an HTML proxy error page, `null`, a
 * legacy flat `{"error":"some_code"}`, a body-less 502 — becomes
 * `internal_error` with the ORIGINAL payload preserved at `details.raw`. It
 * never returns `undefined` and never drops the body.
 *
 * Exported (not just used internally) so market.ts's client — which stays
 * independent of the rest of this barrel per D-F4-o — can still funnel its
 * two throw sites through the SAME parse path instead of growing a second
 * one, per this slice's VC-5.
 */
export function parseApiError(status: number, data: unknown): MarketplaceCentralClientError {
  if (isErrorEnvelope(data)) {
    const { code, message, details } = data.error;
    return new MarketplaceCentralClientError(
      status,
      code,
      typeof message === "string" ? message : code,
      details !== null && typeof details === "object" ? (details as Record<string, unknown>) : {},
    );
  }
  return new MarketplaceCentralClientError(status, "internal_error", "erro inesperado do servidor", { raw: data });
}

export function createMarketplaceCentralClient(options: {
  baseUrl: string;
  fetchImpl?: typeof fetch;
}) {
  const fetchImpl = options.fetchImpl ?? fetch;

  async function getJson<T>(path: string): Promise<T> {
    const response = await fetchImpl(`${options.baseUrl}${path}`, { method: "GET" });
    const data = await response.json();
    if (!response.ok) {
      throw parseApiError(response.status, data);
    }
    return data as T;
  }

  function catalogQuery(params: Record<string, string | number | boolean | undefined>): string {
    const query = new URLSearchParams();
    for (const [key, value] of Object.entries(params)) {
      if (value !== undefined) {
        query.set(key, String(value));
      }
    }
    const encoded = query.toString();
    return encoded ? `?${encoded}` : "";
  }

  function listingQuery(params: ListingListOptions): string {
    const query = new URLSearchParams();
    query.set("installation_id", params.installation_id);
    if (params.q !== undefined) query.set("q", params.q);
    const filters: Array<[string, string | number | boolean | undefined]> = [
      ["status", params.status],
      ["sync_state", params.sync_state],
      ["link_state", params.link_state],
      ["exception", params.exception],
      ["has_exception", params.has_exception],
      ["listing_type_code", params.listing_type_code],
      ["product_id", params.product_id],
    ];
    for (const [key, value] of filters) {
      if (value !== undefined) query.set(`filter.${key}`, String(value));
    }
    if (params.limit !== undefined) query.set("limit", String(params.limit));
    if (params.cursor !== undefined) query.set("cursor", params.cursor);
    return `?${query.toString()}`;
  }

  function syncRunQuery(params: SyncRunListOptions): string {
    const query = new URLSearchParams();
    query.set("installation_id", params.installation_id);
    if (params.cursor !== undefined) query.set("cursor", params.cursor);
    if (params.limit !== undefined) query.set("limit", String(params.limit));
    if (params.module !== undefined) query.set("filter.module", params.module);
    if (params.status !== undefined) query.set("filter.status", params.status);
    return `?${query.toString()}`;
  }

  function orderQuery(params: OrderListOptions): string {
    const query = new URLSearchParams();
    query.set("installation_id", params.installation_id);
    if (params.cursor !== undefined) query.set("cursor", params.cursor);
    if (params.limit !== undefined) query.set("limit", String(params.limit));
    if (params.status !== undefined) query.set("status", params.status);
    if (params.date_from !== undefined) query.set("date_from", params.date_from);
    if (params.date_to !== undefined) query.set("date_to", params.date_to);
    if (params.q !== undefined) query.set("q", params.q);
    return `?${query.toString()}`;
  }

  function mutationQuery(params: Record<string, string | number | undefined>): string {
    const query = new URLSearchParams();
    for (const [key, value] of Object.entries(params)) {
      if (value !== undefined) query.set(key, String(value));
    }
    return `?${query.toString()}`;
  }

  async function putJson<T>(path: string, body: unknown): Promise<T> {
    const response = await fetchImpl(`${options.baseUrl}${path}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    const data = await response.json();
    if (!response.ok) {
      throw parseApiError(response.status, data);
    }
    return data as T;
  }

  async function deleteJson(path: string): Promise<void> {
    const response = await fetchImpl(`${options.baseUrl}${path}`, { method: "DELETE" });
    if (!response.ok) {
      const data = await response.json();
      throw parseApiError(response.status, data);
    }
  }

  async function postJson<T>(path: string, body: unknown): Promise<T> {
    const response = await fetchImpl(`${options.baseUrl}${path}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    const data = await response.json();
    if (!response.ok) {
      throw parseApiError(response.status, data);
    }
    return data as T;
  }

  async function postVoid(path: string, body: unknown): Promise<void> {
    const response = await fetchImpl(`${options.baseUrl}${path}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    if (!response.ok) {
      const data = await response.json();
      throw parseApiError(response.status, data);
    }
  }

  // postMultipart posts a FormData body. The Content-Type header is intentionally
  // NOT set so the runtime supplies the multipart boundary. Non-2xx bodies are
  // the nested ErpImportError/ErpImportConflict envelope; callers branch on
  // status (409 duplicate/in-progress, 422 missing column) and read the extra
  // fields off `e.details` (`column`, `import_id`, `protocol`). The raw parsed
  // body is available at `e.details.raw` only when the body did NOT parse as a
  // well-formed envelope — there is no second `body` property to read.
  async function postMultipart<T>(path: string, body: FormData): Promise<T> {
    const response = await fetchImpl(`${options.baseUrl}${path}`, {
      method: "POST",
      body,
    });
    const data = await response.json();
    if (!response.ok) {
      throw parseApiError(response.status, data);
    }
    return data as T;
  }

  return {
    listCatalogProductFacts: (options: CatalogPageOptions = {}) =>
      getJson<CatalogProductFactPage>(
        `/catalog/products${catalogQuery({ cursor: options.cursor, limit: options.limit, erp_source: options.erp_source, include_all: options.include_all })}`,
      ),
    catalogProductFactsByIds: (options: CatalogProductFactsByIdsOptions) =>
      getJson<CatalogProductFactPage>(
        `/catalog/products${catalogQuery({ ids: options.ids.join(","), erp_source: options.erp_source })}`,
      ),
    searchCatalogProductFacts: (options: CatalogSearchPageOptions) =>
      getJson<CatalogProductFactPage>(
        `/catalog/products/search${catalogQuery({ q: options.q, cursor: options.cursor, limit: options.limit, erp_source: options.erp_source, include_all: options.include_all })}`,
      ),
    getCatalogAssortmentCounts: (options: CatalogAssortmentCountsOptions = {}) =>
      getJson<CatalogAssortmentCounts>(
        `/catalog/products/counts${catalogQuery({ erp_source: options.erp_source })}`,
      ),
    listListings: (options: ListingListOptions) =>
      getJson<ListingPage>(`/listings${listingQuery(options)}`),
    listListingsByProduct: (options: ListingListOptions) =>
      getJson<ListingGroupPage>(`/listings/by-product${listingQuery(options)}`),
    getListing: (id: string) =>
      getJson<ListingDetail>(`/listings/${encodeURIComponent(id)}`),
    getCategoryAttributes: (categoryId: string) =>
      getJson<CategoryAttributesResponse>(
        `/listings/categories/${encodeURIComponent(categoryId)}/attributes`,
      ),
    getListingsSummary: (installationId: string) =>
      getJson<ListingSummary>(`/listings/summary?installation_id=${encodeURIComponent(installationId)}`),
    getDashboardSummary: (installationId: string) =>
      getJson<DashboardSummary>(`/dashboard/summary?installation_id=${encodeURIComponent(installationId)}`),
    listSyncRuns: (options: SyncRunListOptions) =>
      getJson<SyncRunPage>(`/sync/runs${syncRunQuery(options)}`),
    getSyncHealth: () => getJson<SyncHealth>("/sync/health"),
    refreshListings: (req: RefreshListingsRequest) =>
      postJson<RefreshListingsAccepted>("/listings/refresh", req),
    /** @deprecated Use listCatalogProductFacts; the endpoint is now paginated. */
    listCatalogProducts: (options: CatalogPageOptions = {}) =>
      getJson<CatalogProductFactPage>(
        `/catalog/products${catalogQuery({ cursor: options.cursor, limit: options.limit })}`,
      ),
    createErpImport: (file: File | Blob, source: ErpImportSourceInput, fileName?: string) => {
      const form = new FormData();
      form.append("file", file, fileName ?? (file instanceof File ? file.name : "import.xlsx"));
      form.append("source", source);
      return postMultipart<ErpImportCreated>("/erp/imports", form);
    },
    // The tenant's active source lives in the database and drives every
    // mirror-backed read server-side. The client reads and writes it here so the
    // "Fonte ativa" toggle changes what the platform actually reads, instead of
    // a browser-local preference the backend never sees.
    getActiveSource: () => getJson<ActiveSourceConfig>("/config/active-source"),
    setActiveSource: (req: SetActiveSourceRequest) =>
      putJson<ActiveSourceConfig>("/config/active-source", req),
    getSellableAssortment: () =>
      getJson<SellableAssortmentConfig>("/config/sellable-assortment"),
    setSellableAssortment: (req: SetSellableAssortmentRequest) =>
      putJson<SellableAssortmentConfig>("/config/sellable-assortment", req),
    listErpImports: () => getJson<ErpImportList>("/erp/imports"),
    getErpImport: (id: string) => getJson<ErpImportDetail>(`/erp/imports/${encodeURIComponent(id)}`),
    getErpImportChain: (id: string) => getJson<ErpImportChain>(`/erp/imports/${encodeURIComponent(id)}/chain`),
    listMarketplaceAccounts: () => getJson<ListResponse<MarketplaceAccount>>("/marketplaces/accounts"),
    listMarketplacePolicies: () => getJson<ListResponse<MarketplacePolicy>>("/marketplaces/policies"),
    listMarketplaceDefinitions: () => getJson<ListResponse<MarketplaceDefinition>>("/marketplaces/definitions"),
    listIntegrationProviders: () => getJson<ListResponse<IntegrationProviderDefinition>>("/integrations/providers"),
    listIntegrationInstallations: () => getJson<ListResponse<IntegrationInstallation>>("/integrations/installations"),
    listIntegrationOperationRuns: (installationId: string) =>
      getJson<ListResponse<IntegrationOperationRun>>(`/integrations/installations/${installationId}/operations`),
    startIntegrationAuthorization: (installationId: string) =>
      postJson<IntegrationAuthorizeResponse>(`/integrations/installations/${installationId}/auth/authorize`, {}),
    startIntegrationReauthorization: (installationId: string) =>
      postJson<IntegrationAuthorizeResponse>(`/integrations/installations/${installationId}/reauth/authorize`, {}),
    submitIntegrationCredentials: (installationId: string, req: SubmitIntegrationCredentialsRequest) =>
      postJson<IntegrationAuthStatusResponse>(`/integrations/installations/${installationId}/auth/credentials`, req),
    disconnectIntegrationInstallation: (installationId: string) =>
      postJson<IntegrationAuthStatusResponse>(`/integrations/installations/${installationId}/disconnect`, {}),
    getIntegrationAuthStatus: (installationId: string) =>
      getJson<IntegrationAuthStatusResponse>(`/integrations/installations/${installationId}/auth/status`),
    startIntegrationFeeSync: (installationId: string) =>
      postJson<IntegrationFeeSyncAccepted>(`/integrations/installations/${installationId}/fee-sync`, {}),
    probeIntegrationAccount: (installationId: string) =>
      postJson<IntegrationAccountProbe>(`/integrations/installations/${installationId}/probes/account`, {}),
    probeIntegrationListings: (installationId: string, limit = 20) =>
      getJson<ListResponse<IntegrationListingProbe>>(
        `/integrations/installations/${installationId}/probes/listings?limit=${encodeURIComponent(String(limit))}`,
      ),
    probeIntegrationOrders: (installationId: string, limit = 20) =>
      getJson<ListResponse<IntegrationOrderProbe>>(
        `/integrations/installations/${installationId}/probes/orders?limit=${encodeURIComponent(String(limit))}`,
      ),
    probeIntegrationFeeQuote: (
      installationId: string,
      input: { listing_type_id: string; price: number; site_id?: string; category_id?: string; currency_id?: string },
    ) =>
      getJson<IntegrationFeeQuoteSnapshot>(
        `/integrations/installations/${installationId}/probes/fee-quote?listing_type_id=${encodeURIComponent(input.listing_type_id)}&price=${encodeURIComponent(String(input.price))}${input.site_id ? `&site_id=${encodeURIComponent(input.site_id)}` : ""}${input.category_id ? `&category_id=${encodeURIComponent(input.category_id)}` : ""}${input.currency_id ? `&currency_id=${encodeURIComponent(input.currency_id)}` : ""}`,
      ),
    probeIntegrationStock: (
      installationId: string,
      input: { provider_item_id: string; provider_variation_id?: string },
    ) =>
      getJson<IntegrationStockProbe>(
        `/integrations/installations/${installationId}/probes/stock?provider_item_id=${encodeURIComponent(input.provider_item_id)}${input.provider_variation_id ? `&provider_variation_id=${encodeURIComponent(input.provider_variation_id)}` : ""}`,
      ),
    importProductLinkListingSnapshots: (req: { installation_id: string; limit?: number }) =>
      postJson<ProductLinkListingSnapshotImportResult>("/product-links/listing-snapshots/imports", req),
    generateProductLinkCandidates: (req: { installation_id: string; limit?: number }) =>
      postJson<ProductLinkCandidateGenerationResult>("/product-links/link-candidates/generations", req),
    listProductLinkCandidates: (installationId: string, limit = 20) =>
      getJson<ListResponse<ProductLinkCandidateItem>>(
        `/product-links/link-candidates?installation_id=${encodeURIComponent(installationId)}&limit=${encodeURIComponent(String(limit))}`,
      ),
    listProductLinkWorkflows: (installationId: string, limit = 20) =>
      getJson<ListResponse<ProductLinkWorkflowItem>>(
        `/product-links/link-workflows?installation_id=${encodeURIComponent(installationId)}&limit=${encodeURIComponent(String(limit))}`,
      ),
    approveProductLinkCandidate: (req: { candidate_id: string; reason?: string; actor: ProductLinkActor }) =>
      postJson<ProductLinkResolutionResult>("/product-links/link-resolutions/approve-candidate", req),
    rejectProductLinkListing: (req: {
      installation_id: string;
      provider_code: string;
      provider_item_id: string;
      provider_variation_id?: string;
      reason?: string;
      actor: ProductLinkActor;
    }) => postJson<ProductLinkResolutionResult>("/product-links/link-resolutions/reject-listing", req),
    manualResolveProductLink: (req: {
      installation_id: string;
      provider_code: string;
      provider_item_id: string;
      provider_variation_id?: string;
      internal_product_id: number;
      internal_product_name?: string;
      internal_reference_code?: string;
      reason?: string;
      actor: ProductLinkActor;
    }) => {
      if (!Number.isInteger(req.internal_product_id) || req.internal_product_id <= 0) {
        throw new Error("invalid_identity: internal_product_id must be a positive integer");
      }
      return postJson<ProductLinkResolutionResult>("/product-links/link-resolutions/manual-resolve", req);
    },
    previewProductLinkBatch: (req: { approvals: ProductLinkBatchApproval[] }) =>
      postJson<PreviewProductLinkBatchResponse>("/product-links/link-resolutions/batch-preview", req),
    applyProductLinkBatch: (req: { approvals: ProductLinkBatchApproval[]; actor: ProductLinkActor }) =>
      postJson<ApplyProductLinkBatchResponse>("/product-links/link-resolutions/batch", req),
    undoProductLinkResolution: (auditId: string, req: { reason?: string; actor?: ProductLinkActor } = {}) =>
      postJson<ProductLinkResolutionResult>(`/product-links/link-resolutions/${encodeURIComponent(auditId)}/undo`, req),
    undoProductLinkBatch: (batchId: string) =>
      postJson<UndoProductLinkBatchResponse>(`/product-links/link-resolutions/batch/${encodeURIComponent(batchId)}/undo`, {}),
    listInventoryStockRisks: (input: {
      installation_id: string;
      state?: StockRiskState;
      link_state?: ProductLinkState;
      actionability?: StockRiskActionability;
      limit?: number;
    }) =>
      getJson<ListResponse<InventoryStockRiskItem>>(
        `/inventory/stock-risks?installation_id=${encodeURIComponent(input.installation_id)}${input.state ? `&state=${encodeURIComponent(input.state)}` : ""}${input.link_state ? `&link_state=${encodeURIComponent(input.link_state)}` : ""}${input.actionability ? `&actionability=${encodeURIComponent(input.actionability)}` : ""}${input.limit ? `&limit=${encodeURIComponent(String(input.limit))}` : ""}`,
      ),
    applyInventoryManualStockAction: (req: {
      stock_action_id: string;
      installation_id: string;
      provider_item_id: string;
      provider_variation_id?: string;
      requested_quantity: number;
      reason?: string;
      approval: InventoryManualApproval;
    }) => postJson<ApplyInventoryStockActionResponse>("/inventory/stock-actions/manual-apply", req),
    importMarketplaceOrders: (req: { installation_id: string; limit?: number }) =>
      postJson<ImportMarketplaceOrdersResponse>("/orders/import", req),
    /** @deprecated per hub lock-exception D-02 */
    listMarketplaceOrders: (installationId: string, limit = 20) =>
      getJson<ListResponse<MarketplaceOrder>>(
        `/orders?installation_id=${encodeURIComponent(installationId)}&limit=${encodeURIComponent(String(limit))}`,
      ),
    listOrders: (options: OrderListOptions) =>
      getJson<OrderPage>(`/orders${orderQuery(options)}`),
    getOrder: (installationId: string, providerOrderId: string) =>
      getJson<OrderRead>(
        `/orders/${encodeURIComponent(providerOrderId)}?installation_id=${encodeURIComponent(installationId)}`,
      ),
    // Records our own faturado_at fact so the order moves from the "A FATURAR"
    // bucket to "A ENVIAR" (goal B). OUR-DB only — never a Mercado Livre write.
    markOrderFaturado: (installationId: string, providerOrderId: string) =>
      postVoid(`/orders/${encodeURIComponent(providerOrderId)}/faturado`, {
        installation_id: installationId,
      }),
    getOrderSummary: (installationId: string, options?: { by?: "status" }) =>
      getJson<OrderSummary>(
        `/orders/summary?installation_id=${encodeURIComponent(installationId)}${options?.by ? `&by=${encodeURIComponent(options.by)}` : ""}`,
      ),
    getAssistedSankhyaLinkage: (installationId: string, providerOrderId: string) =>
      getJson<AssistedSankhyaLinkageResponse>(
        `/orders/${encodeURIComponent(providerOrderId)}/sankhya-linkage?installation_id=${encodeURIComponent(installationId)}`,
      ),
    listAssistedSankhyaLinkageCandidates: (installationId: string, providerOrderId: string) =>
      getJson<AssistedSankhyaCandidateListResponse>(
        `/orders/${encodeURIComponent(providerOrderId)}/sankhya-linkage/candidates?installation_id=${encodeURIComponent(installationId)}`,
      ),
    confirmAssistedSankhyaLinkage: (
      installationId: string,
      providerOrderId: string,
      req: ConfirmAssistedSankhyaLinkageRequest,
    ) =>
      postJson<ConfirmAssistedSankhyaLinkageResponse>(
        `/orders/${encodeURIComponent(providerOrderId)}/sankhya-linkage/confirm?installation_id=${encodeURIComponent(installationId)}`,
        req,
      ),
    importProfitabilityMarginInputs: (req: { installation_id: string; limit?: number }) =>
      postJson<ImportProfitabilityMarginInputsResponse>("/profitability/margin-inputs/import", req),
    listProfitabilityMarginInputs: (installationId: string, limit = 50) =>
      getJson<ListResponse<ProfitabilityMarginInput>>(
        `/profitability/margin-inputs?installation_id=${encodeURIComponent(installationId)}&limit=${encodeURIComponent(String(limit))}`,
      ),
    createProfitabilityManualAdjustment: (req: {
      installation_id: string;
      idempotency_key: string;
      provider_order_id: string;
      provider_item_id?: string;
      provider_variation_id?: string;
      scope?: ProfitabilityInputScope;
      category?: ProfitabilityManualAdjustmentCategory;
      amount: number;
      currency?: string;
      reason: string;
      actor: ProfitabilityActor;
    }) => postJson<ProfitabilityManualAdjustment>("/profitability/manual-adjustments", req),
    listProfitabilityManualAdjustments: (installationId: string, limit = 50) =>
      getJson<ListResponse<ProfitabilityManualAdjustment>>(
        `/profitability/manual-adjustments?installation_id=${encodeURIComponent(installationId)}&limit=${encodeURIComponent(String(limit))}`,
      ),
    calculateProfitabilityProfitSnapshots: (req: { installation_id: string; limit?: number }) =>
      postJson<CalculateProfitabilityProfitSnapshotsResponse>("/profitability/profit-snapshots/calculate", req),
    listProfitabilityProfitSnapshots: (installationId: string, limit = 50) =>
      getJson<ListResponse<ProfitabilityProfitSnapshot>>(
        `/profitability/profit-snapshots?installation_id=${encodeURIComponent(installationId)}&limit=${encodeURIComponent(String(limit))}`,
      ),
    listMarketplaceFeeSchedules: (marketplaceCode: string) =>
      getJson<ListResponse<MarketplaceFeeSchedule>>(`/marketplaces/fee-schedules?marketplace_code=${encodeURIComponent(marketplaceCode)}`),
    listPricingSimulations: () => getJson<ListResponse<PricingSimulation>>("/pricing/simulations"),
    createMarketplaceAccount: (req: CreateMarketplaceAccountRequest) =>
      postJson<MarketplaceAccount>("/marketplaces/accounts", req),
    createMarketplacePolicy: (req: CreateMarketplacePolicyRequest) =>
      postJson<MarketplacePolicy>("/marketplaces/policies", req),
    createIntegrationInstallation: (req: CreateIntegrationInstallationRequest) =>
      postJson<IntegrationInstallation>("/integrations/installations", req),
    runPricingSimulation: (req: RunPricingSimulationRequest) =>
      postJson<PricingSimulation>("/pricing/simulations", req),
    runBatchSimulation: (req: BatchSimulationRequest) =>
      postJson<{ items: BatchSimulationItem[] }>("/pricing/simulations/batch", req),
    getPricingProfile: () => getJson<PricingCalcProfile>("/pricing/profile"),
    putPricingProfile: (req: PricingCalcProfile) =>
      putJson<PricingCalcProfile>("/pricing/profile", req),
    listPricingDifal: () => getJson<PricingDifalListResponse>("/pricing/difal"),
    putPricingDifalOverride: (uf: string, req: PricingDifalOverrideRequest) =>
      putJson<PricingDifalOverrideResponse>(`/pricing/difal/${encodeURIComponent(uf)}`, req),
    listPricingScenarios: () => getJson<PricingScenarioListResponse>("/pricing/scenarios"),
    createPricingScenario: (req: PricingScenarioRequest) =>
      postJson<PricingScenario>("/pricing/scenarios", req),
    deletePricingScenario: (id: string) =>
      deleteJson(`/pricing/scenarios/${encodeURIComponent(id)}`),
    pricingDecompose: (req: PricingCalcInput) =>
      postJson<PricingDecomposeResponse>("/pricing/decompose", req),
    pricingSolveTarget: (req: PricingCalcInput) =>
      postJson<PricingSolveResponse>("/pricing/solve", req),
    getMelhorEnvioStatus: () =>
      getJson<{ connected: boolean }>("/connectors/melhor-envio/status"),
    createMutation: (req: CreateMutationRequest) =>
      postJson<MutationProtocol>("/mutations", req),
    previewMutation: (protocolId: string) =>
      postJson<MutationPreview>(`/mutations/${encodeURIComponent(protocolId)}/preview`, {}),
    approveMutation: (protocolId: string, req: ApproveMutationRequest) =>
      postJson<MutationProtocol>(`/mutations/${encodeURIComponent(protocolId)}/approve`, req),
    cancelMutation: (protocolId: string) =>
      postJson<MutationProtocol>(`/mutations/${encodeURIComponent(protocolId)}/cancel`, {}),
    retryMutationFailures: (protocolId: string) =>
      postJson<MutationProtocol>(`/mutations/${encodeURIComponent(protocolId)}/retry`, {}),
    listMutations: (options: MutationListOptions) =>
      getJson<MutationPage<MutationProtocol>>(
        `/mutations${mutationQuery({
          installation_id: options.installation_id,
          state: options.state,
          type: options.type,
          limit: options.limit,
          cursor: options.cursor,
        })}`,
      ),
    getMutation: (protocolId: string) =>
      getJson<MutationProtocol>(`/mutations/${encodeURIComponent(protocolId)}`),
    listMutationItems: (protocolId: string, options: MutationItemListOptions = {}) =>
      getJson<MutationPage<MutationItem>>(
        `/mutations/${encodeURIComponent(protocolId)}/items${mutationQuery({
          state: options.state,
          limit: options.limit,
          cursor: options.cursor,
        })}`,
      ),

    // Catalog
    /** @deprecated Use searchCatalogProductFacts; the endpoint now returns the IC-01 page envelope. */
    searchCatalogProducts: (query: string, limit?: number) =>
      getJson<CatalogProductFactPage>(`/catalog/products/search${catalogQuery({ q: query, limit })}`),
    getCatalogProduct: (productId: number) =>
      getJson<CanonicalCatalogProduct>(`/catalog/products/${productId}`),
    listTaxonomyNodes: () =>
      getJson<ListResponse<TaxonomyNode>>("/catalog/taxonomy"),
    getProductEnrichment: (productId: string) =>
      getJson<ProductEnrichment>(`/catalog/products/${productId}/enrichment`),
    updateProductEnrichment: (productId: string, data: Partial<ProductEnrichment>) =>
      putJson<ProductEnrichment>(`/catalog/products/${productId}/enrichment`, data),

    // Classifications
    listClassifications: () =>
      getJson<ListResponse<Classification>>("/classifications"),
    createClassification: (req: CreateClassificationRequest) =>
      postJson<Classification>("/classifications", req),
    getClassification: (id: string) =>
      getJson<Classification>(`/classifications/${id}`),
    updateClassification: (id: string, req: UpdateClassificationRequest) =>
      putJson<Classification>(`/classifications/${id}`, req),
    deleteClassification: (id: string) =>
      deleteJson(`/classifications/${id}`),
    listMarketObservations: (installationId: string, listingIds: string[]) => {
      const query = new URLSearchParams();
      query.set("installation_id", installationId);
      query.set("listing_ids", listingIds.join(","));
      return getJson<MarketObservationPage>(`/market/observations?${query.toString()}`);
    },
    listMarketReferences: (productIds: string[]) => {
      const query = new URLSearchParams();
      query.set("product_ids", productIds.join(","));
      return getJson<MarketReferencePage>(`/market/references?${query.toString()}`);
    },
  };
}

export interface DashboardSummary {
  sync_errors: number | null;
  pending_links: number | null;
  below_margin: number | null;
  missing_gtin: number | null;
  orders_today: number | null;
  orders_7d: number | null;
  last_sync_at: Record<string, string | null> | null;
  degraded: Array<"listings" | "linkage" | "orders" | "sync">;
  as_of: string;
}

export type SyncRunStatus = "queued" | "running" | "succeeded" | "failed" | "cancelled";

export interface SyncRun {
  operation_run_id: string;
  installation_id: string;
  module: string;
  status: SyncRunStatus;
  result_code: string;
  failure_code: string;
  translated_error_code: string;
  attempt_count: number;
  duration_ms: number;
  started_at: string;
  finished_at: string | null;
}

export interface SyncRunPage {
  items: SyncRun[];
  next_cursor: string | null;
  page_size: number;
  as_of: string;
}

export interface SyncRunListOptions {
  installation_id: string;
  cursor?: string;
  limit?: number;
  module?: string;
  status?: SyncRunStatus;
}

export interface SyncHealthError {
  message: string;
  at: string;
}

export interface SyncHealthEntity {
  entity: string;
  last_success_at: string | null;
  last_incremental_at: string | null;
  consecutive_failures: number;
  phase: string | null;
  last_error: SyncHealthError | null;
}

export interface SyncHealthWebhook {
  last_notification_at: string | null;
  pending: number;
  dropped_24h: number;
}

export interface SyncHealth {
  entities: SyncHealthEntity[];
  webhook: SyncHealthWebhook;
}

export interface MarketMoney {
  amount: string;
  currency: "BRL";
}

export interface MarketCatalogStats {
  min: string;
  p25: string;
  median: string;
  trimmed_mean: string;
  p75: string;
  max: string;
  n_offers: number;
  n_sellers: number;
}

export type MarketCompetitiveStatus = "winning" | "competing" | "not_listed";
export type MarketEvidenceState = "observed" | "insufficient_market" | "no_price_evidence";
export type MarketSource = "official_api" | "vendor" | "manual";
export type MarketMatchState = "accept" | "review" | "reject" | "no_candidate";

export interface MarketObservation {
  listing_id: string;
  our_sale_price: MarketMoney | null;
  competitive_status: MarketCompetitiveStatus | null;
  winner_price: MarketMoney | null;
  competitive_target: MarketMoney | null;
  catalog_offer_price: MarketMoney | null;
  catalog_stats: MarketCatalogStats | null;
  evidence_state: MarketEvidenceState;
  captured_at: string | null;
  source: MarketSource | null;
}

export interface MarketReference {
  product_id: string;
  catalog_product_id: string | null;
  match_state: MarketMatchState;
  match_method: string | null;
  catalog_stats: MarketCatalogStats | null;
  evidence_state: MarketEvidenceState;
  captured_at: string | null;
  source: MarketSource | null;
}

export interface MarketObservationPage {
  items: MarketObservation[];
  as_of: string;
}

export interface MarketReferencePage {
  items: MarketReference[];
  as_of: string;
}
