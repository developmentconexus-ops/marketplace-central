import { onlineManager, QueryClientProvider } from "@tanstack/react-query";
import { createMarketplaceCentralClient } from "@marketplace-central/sdk-runtime";
import { createWebQueryClient } from "@marketplace-central/web-query";
import { render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { SyncHealthCard } from "./SyncHealthCard";

// Regression guard for the live-drive defect: the hub killed the backend
// container, loaded /integracoes in a real browser, and captured the
// rendered section's innerHTML as literally just the <h2> heading — no
// LoadingState, no ErrorState, indefinitely. DevTools Network showed exactly
// ONE request (GET /sync/health, net::ERR_CONNECTION_REFUSED) with zero
// retries after.
//
// Root cause (confirmed here with real instrumentation, not asserted): the
// production QueryClient (createWebQueryClient) leaves networkMode at its
// TanStack Query default of "online", which gates retries on
// onlineManager.isOnline(). The injected-state test in
// SyncHealthCard.test.tsx mocks getSyncHealth to synchronously reject,
// which bypasses react-query's own fetch/retry machinery entirely and can
// never observe this: it goes straight from "pending" to "error" because no
// real retry cycle - and therefore no onlineManager gate - is ever
// exercised.
//
// This test drives the REAL client (createMarketplaceCentralClient, the same
// one ClientContext.tsx wires up in production) against a real TCP
// connection refusal (127.0.0.1 on a port nothing listens on), through the
// REAL production QueryClient factory (createWebQueryClient, the same one
// main.tsx uses) — and flips onlineManager to "offline" right after the
// first real fetch attempt fails, mirroring a network blip landing at the
// same moment as the outage. Logged instrumentation (see task notes)
// confirmed the resulting query state is exactly
// { status: "pending", fetchStatus: "paused", failureCount: 1 } forever:
// query.isLoading is false (fetchStatus isn't "fetching"), query.isError is
// false (status never reaches "error"), query.data is undefined — so
// SyncHealthCard's `isLoading ? … : isError ? … : data ? … : null` chain
// falls through all three branches, exactly reproducing the blank card.
vi.mock("../../app/ClientContext", () => ({
  useClient: () =>
    createMarketplaceCentralClient({
      // Nothing listens here — a real ECONNREFUSED, not a mocked rejection.
      baseUrl: "http://127.0.0.1:59999",
    }),
}));

describe("SyncHealthCard against a real connection-refused network failure", () => {
  afterEach(() => {
    // onlineManager is a module-level singleton; restore it so this test's
    // manual offline flip never leaks into other test files in the run.
    onlineManager.setOnline(true);
  });

  it("reaches the named ErrorState instead of staying silently blank when the browser reports itself offline mid-outage", async () => {
    const queryClient = createWebQueryClient();
    render(
      <QueryClientProvider client={queryClient}>
        <SyncHealthCard />
      </QueryClientProvider>,
    );

    // Give the first real fetch attempt a moment to fire and fail, then flip
    // the browser's reported connectivity — the condition that turns an
    // ordinary connection refusal into a permanently paused, silently blank
    // query under the default networkMode.
    await waitFor(() => {
      const query = queryClient.getQueryCache().find({ queryKey: ["sync", "health"] });
      expect(query?.state.fetchFailureCount).toBeGreaterThanOrEqual(1);
    });
    onlineManager.setOnline(false);

    // The card must not sit blank forever: it has to land on an observable
    // state (ErrorState) that is not "nothing rendered". waitFor polls on
    // real timers, matching react-query's own real retry/backoff timers.
    const alert = await waitFor(() => screen.getByRole("alert"), { timeout: 10_000 });
    expect(alert).toBeInTheDocument();

    // The section must never be left with nothing but its heading — that is
    // the exact blank-card defect the hub captured live.
    const section = screen.getByRole("region", { name: "Saúde do sync" });
    expect(section.textContent).not.toBe("Saúde do sync");
  }, 15_000);
});
