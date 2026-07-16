import type { ListingSummary } from "@marketplace-central/sdk-runtime";
import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { ListingsSummary } from "./ListingsSummary";

const summary: ListingSummary = {
  total: 10,
  active: 6,
  paused: 4,
  exceptions: {
    sync_error: 1,
    stale: 2,
    unlinked: 3,
    below_margin_worst_case: null,
    margin_unknown: null,
  },
  as_of: "2026-07-16T12:00:00Z",
};

describe("ListingsSummary", () => {
  it("renders counters from the summary and preserves nullable unknowns", () => {
    render(<ListingsSummary isPending={false} isError={false} data={summary} onRetry={vi.fn()} />);

    expect(screen.getByText("10")).toBeInTheDocument();
    expect(screen.getByText("6")).toBeInTheDocument();
    expect(screen.getByText("4")).toBeInTheDocument();
    expect(screen.getByText("1")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();
    expect(screen.getByText("3")).toBeInTheDocument();
    expect(screen.getAllByText("—")).toHaveLength(2);
    expect(screen.getByLabelText("Atualização dos dados")).toBeInTheDocument();
  });

  it("renders an error with a working retry action", () => {
    const onRetry = vi.fn();
    render(<ListingsSummary isPending={false} isError data={undefined} onRetry={onRetry} />);

    expect(screen.getByRole("alert")).toHaveTextContent("Erro ao carregar.");
    fireEvent.click(screen.getByRole("button", { name: "Tentar novamente" }));
    expect(onRetry).toHaveBeenCalledOnce();
  });
});
