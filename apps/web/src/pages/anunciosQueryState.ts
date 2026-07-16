import type {
  ListingException,
  ListingLinkState,
  ListingListOptions,
  ListingSyncState,
} from "@marketplace-central/sdk-runtime";

export type AnunciosTab = "todos" | "ativos" | "pausados" | "pendencia";

export interface AnunciosFilterState {
  exception?: ListingException;
  sync_state?: ListingSyncState;
  link_state?: ListingLinkState;
  listing_type_code?: string;
}

export interface AnunciosQueryState {
  tab: AnunciosTab;
  q: string;
  filters: AnunciosFilterState;
}

const exceptionValues = ["sync_error", "stale", "unlinked", "below_margin"] as const;
const syncStateValues = ["synced", "error", "stale", "queued", "syncing", "paused_sync"] as const;
const linkStateValues = ["unresolved", "conflict", "resolved", "rejected"] as const;

function isListingException(value: string): value is ListingException {
  return exceptionValues.includes(value as ListingException);
}

function isListingSyncState(value: string): value is ListingSyncState {
  return syncStateValues.includes(value as ListingSyncState);
}

function isListingLinkState(value: string): value is ListingLinkState {
  return linkStateValues.includes(value as ListingLinkState);
}

function setFilterParam(searchParams: URLSearchParams, key: string, value: string | undefined): void {
  if (value) searchParams.set(`filter.${key}`, value);
}

export function parseAnunciosQueryState(searchParams: URLSearchParams): AnunciosQueryState {
  const tabValue = searchParams.get("tab");
  const tab: AnunciosTab =
    tabValue === "ativos" || tabValue === "pausados" || tabValue === "pendencia"
      ? tabValue
      : "todos";
  const exception = searchParams.get("filter.exception");
  const syncState = searchParams.get("filter.sync_state");
  const linkState = searchParams.get("filter.link_state");
  const listingTypeCode = searchParams.get("filter.listing_type_code");
  const filters: AnunciosFilterState = {};

  if (exception && isListingException(exception)) filters.exception = exception;
  if (syncState && isListingSyncState(syncState)) filters.sync_state = syncState;
  if (linkState && isListingLinkState(linkState)) filters.link_state = linkState;
  if (listingTypeCode) filters.listing_type_code = listingTypeCode;

  return {
    tab,
    q: searchParams.get("q") ?? "",
    filters,
  };
}

export function applyAnunciosQueryState(
  searchParams: URLSearchParams,
  state: AnunciosQueryState,
): URLSearchParams {
  const next = new URLSearchParams(searchParams);
  next.delete("tab");
  next.delete("q");
  for (const key of Array.from(next.keys())) {
    if (key.startsWith("filter.")) next.delete(key);
  }

  if (state.tab !== "todos") next.set("tab", state.tab);
  if (state.q) next.set("q", state.q);

  const filters = state.filters;
  if (filters.exception && isListingException(filters.exception)) {
    setFilterParam(next, "exception", filters.exception);
  }
  if (filters.sync_state && isListingSyncState(filters.sync_state)) {
    setFilterParam(next, "sync_state", filters.sync_state);
  }
  if (filters.link_state && isListingLinkState(filters.link_state)) {
    setFilterParam(next, "link_state", filters.link_state);
  }
  setFilterParam(next, "listing_type_code", filters.listing_type_code);
  return next;
}

export function clearFilters(searchParams: URLSearchParams): URLSearchParams {
  const next = new URLSearchParams(searchParams);
  for (const key of Array.from(next.keys())) {
    if (key.startsWith("filter.")) next.delete(key);
  }
  return next;
}

export function toListingListOptions(
  state: AnunciosQueryState,
  installationId: string,
): ListingListOptions {
  const options: ListingListOptions = { installation_id: installationId };
  if (state.q) options.q = state.q;

  if (state.tab === "ativos") options.status = "active";
  if (state.tab === "pausados") options.status = "paused";
  if (state.tab === "pendencia") options.has_exception = true;

  if (state.filters.exception && isListingException(state.filters.exception)) {
    options.exception = state.filters.exception;
  }
  if (state.filters.sync_state && isListingSyncState(state.filters.sync_state)) {
    options.sync_state = state.filters.sync_state;
  }
  if (state.filters.link_state && isListingLinkState(state.filters.link_state)) {
    options.link_state = state.filters.link_state;
  }
  if (state.filters.listing_type_code) options.listing_type_code = state.filters.listing_type_code;

  return options;
}
