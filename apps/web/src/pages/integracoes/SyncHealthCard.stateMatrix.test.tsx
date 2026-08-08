import type { QueryObserverResult, UseQueryResult } from "@tanstack/react-query";
import { MarketplaceCentralClientError, type SyncHealth } from "@marketplace-central/sdk-runtime";
import { render, screen, within } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { SyncHealthCard } from "./SyncHealthCard";

// Systematic state-space guard for the blank-card defect (live-drive,
// hub-confirmed TWICE — networkMode:"always" fixed one react-query internal
// path but the design flaw was untouched, so the same blank <h2>-only render
// happened again). The hub's ruling: a jsdom/Node-fetch-against-a-refused-port
// test (SyncHealthCard.realNetwork.test.tsx) cannot reproduce the browser's
// real onlineManager/connectivity-detection timing, so a green result there
// does not prove the live browser behaves the same way. This file proves the
// fix a different, stronger way: it mocks useSyncHealthQuery itself and
// directly injects every UseQueryResult<SyncHealth> SHAPE that
// @tanstack/query-core's own types (query-core/src/types.ts,
// QueryObserverPendingResult / QueryObserverLoadingResult /
// QueryObserverLoadingErrorResult / QueryObserverSuccessResult) say are
// reachable — including the exact one the hub observed live: status:
// "pending", fetchStatus: "paused" (or "idle"), isLoading: false, isError:
// false, data: undefined. No network, no timers, no race — pure state-space
// enumeration, so every combination is exercised deterministically on every
// run.
const mockUseSyncHealthQuery = vi.fn();

vi.mock("@marketplace-central/web-query", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@marketplace-central/web-query")>();
  return {
    ...actual,
    useSyncHealthQuery: (...args: unknown[]) => mockUseSyncHealthQuery(...args),
  };
});

vi.mock("../../app/ClientContext", () => ({
  // useSyncHealthQuery is mocked wholesale below, so the client itself is
  // never invoked — any non-null object satisfies useClient()'s null check.
  useClient: () => ({}),
}));

const FIXTURE_HEALTH: SyncHealth = {
  entities: [
    {
      entity: "products",
      last_success_at: "2026-08-01T11:55:00Z",
      last_incremental_at: "2026-08-01T11:55:00Z",
      consecutive_failures: 0,
      phase: "incremental",
      last_error: null,
    },
  ],
  webhook: { last_notification_at: null, pending: 0, dropped_24h: 0 },
};

function refetchStub() {
  return vi.fn(async () => ({}) as QueryObserverResult<SyncHealth>);
}

// Common QueryObserverBaseResult fields not relevant to the branch under
// test, held at plausible neutral defaults per case below.
const COMMON = {
  dataUpdatedAt: 0,
  errorUpdatedAt: 0,
  failureCount: 0,
  failureReason: null,
  errorUpdateCount: 0,
  isFetched: false,
  isFetchedAfterMount: false,
  isInitialLoading: false,
  isPlaceholderData: false,
  isRefetchError: false,
  isRefetching: false,
  isStale: false,
  isEnabled: true,
  promise: Promise.resolve(undefined as unknown as SyncHealth),
} as const;

type Case = {
  name: string;
  result: UseQueryResult<SyncHealth>;
};

const CASES: Case[] = [
  {
    name: "loading (first fetch in flight): status=pending, fetchStatus=fetching",
    result: {
      ...COMMON,
      status: "pending",
      fetchStatus: "fetching",
      isPending: true,
      isLoading: true,
      isFetching: true,
      isPaused: false,
      isError: false,
      isLoadingError: false,
      isSuccess: false,
      data: undefined,
      error: null,
      refetch: refetchStub(),
    } as UseQueryResult<SyncHealth>,
  },
  {
    name: "error: status=error",
    result: {
      ...COMMON,
      status: "error",
      fetchStatus: "idle",
      isPending: false,
      isLoading: false,
      isFetching: false,
      isPaused: false,
      isError: true,
      isLoadingError: true,
      isSuccess: false,
      data: undefined,
      error: new MarketplaceCentralClientError(
        500,
        "internal_error",
        "falha ao ler sync_state",
        {},
      ),
      refetch: refetchStub(),
    } as UseQueryResult<SyncHealth>,
  },
  {
    name: "success: status=success",
    result: {
      ...COMMON,
      status: "success",
      fetchStatus: "idle",
      isPending: false,
      isLoading: false,
      isFetching: false,
      isPaused: false,
      isError: false,
      isLoadingError: false,
      isSuccess: true,
      data: FIXTURE_HEALTH,
      error: null,
      refetch: refetchStub(),
    } as UseQueryResult<SyncHealth>,
  },
  {
    // THE live-observed defect shape: networkMode:"online" (or any gate that
    // pauses the fetch) leaves the query permanently in this exact state —
    // no cached data, no error surfaced, nothing in flight.
    name: "paused/indeterminate (THE hub-observed defect): status=pending, fetchStatus=paused, isLoading=false, isError=false, data=undefined",
    result: {
      ...COMMON,
      status: "pending",
      fetchStatus: "paused",
      isPending: true,
      isLoading: false,
      isFetching: false,
      isPaused: true,
      isError: false,
      isLoadingError: false,
      isSuccess: false,
      data: undefined,
      error: null,
      refetch: refetchStub(),
    } as UseQueryResult<SyncHealth>,
  },
  {
    // A second in-between shape distinct from "paused": status still
    // "pending" but fetchStatus "idle" (e.g. a disabled/not-yet-triggered
    // query) — also isLoading:false, isError:false, data:undefined, also not
    // any of the three originally-named branches.
    name: "idle/not-yet-fetching: status=pending, fetchStatus=idle, isLoading=false, isError=false, data=undefined",
    result: {
      ...COMMON,
      status: "pending",
      fetchStatus: "idle",
      isPending: true,
      isLoading: false,
      isFetching: false,
      isPaused: false,
      isError: false,
      isLoadingError: false,
      isSuccess: false,
      data: undefined,
      error: null,
      refetch: refetchStub(),
    } as UseQueryResult<SyncHealth>,
  },
];

describe("SyncHealthCard renders real content for every reachable UseQueryResult shape", () => {
  it.each(CASES)("$name", ({ result }) => {
    mockUseSyncHealthQuery.mockReturnValue(result);
    render(<SyncHealthCard />);

    const section = screen.getByRole("region", { name: "Saúde do sync" });

    // The exact regression this file exists to catch: the section holding
    // nothing but its own heading, forever.
    expect(section.textContent).not.toBe("Saúde do sync");

    // Strong universal assertion: SOME real, observable state is present —
    // a loading indicator, an alert (error/unknown), or actual
    // entities/webhook content. Never silence.
    const hasLoading = within(section).queryByRole("status") !== null;
    const hasAlert = within(section).queryByRole("alert") !== null;
    const hasEntitiesOrWebhook =
      within(section).queryByTestId("sync-health-entities") !== null ||
      within(section).queryByTestId("sync-health-webhook-initial") !== null ||
      within(section).queryByTestId("sync-health-webhook-active") !== null;

    expect(hasLoading || hasAlert || hasEntitiesOrWebhook).toBe(true);
  });
});
