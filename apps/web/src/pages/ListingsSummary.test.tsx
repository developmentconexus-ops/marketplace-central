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
    render(
      <ListingsSummary
        isPending={false}
        isError={false}
        data={summary}
        onRetry={vi.fn()}
        onExceptionChipClick={vi.fn()}
      />,
    );

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
    render(
      <ListingsSummary
        isPending={false}
        isError
        data={undefined}
        onRetry={onRetry}
        onExceptionChipClick={vi.fn()}
      />,
    );

    expect(screen.getByRole("alert")).toHaveTextContent("Erro ao carregar.");
    fireEvent.click(screen.getByRole("button", { name: "Tentar novamente" }));
    expect(onRetry).toHaveBeenCalledOnce();
  });

  it("renders exception chips only for non-zero counters and calls onExceptionChipClick", () => {
    const onExceptionChipClick = vi.fn();
    const withNewExceptions: ListingSummary = {
      ...summary,
      exceptions: {
        ...summary.exceptions,
        sem_vinculo: 4,
        abaixo_custo: 0,
        sem_evidencia: null,
      },
    };
    render(
      <ListingsSummary
        isPending={false}
        isError={false}
        data={withNewExceptions}
        onRetry={vi.fn()}
        onExceptionChipClick={onExceptionChipClick}
      />,
    );

    const chip = screen.getByRole("button", { name: "Sem vínculo: 4" });
    expect(chip).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /Abaixo do custo/ })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /Sem evidência/ })).not.toBeInTheDocument();

    fireEvent.click(chip);
    expect(onExceptionChipClick).toHaveBeenCalledWith("sem_vinculo");
  });

  it("marks the active exception chip as pressed", () => {
    const withNewExceptions: ListingSummary = {
      ...summary,
      exceptions: { ...summary.exceptions, sem_vinculo: 4, abaixo_custo: 2 },
    };
    render(
      <ListingsSummary
        isPending={false}
        isError={false}
        data={withNewExceptions}
        onRetry={vi.fn()}
        onExceptionChipClick={vi.fn()}
        activeException="abaixo_custo"
      />,
    );

    expect(screen.getByRole("button", { name: "Sem vínculo: 4" })).toHaveAttribute("aria-pressed", "false");
    expect(screen.getByRole("button", { name: "Abaixo do custo: 2" })).toHaveAttribute("aria-pressed", "true");
  });

  it("hides exception chips when the summary query errors (isolated from the table)", () => {
    render(
      <ListingsSummary
        isPending={false}
        isError
        data={undefined}
        onRetry={vi.fn()}
        onExceptionChipClick={vi.fn()}
      />,
    );

    expect(screen.queryByRole("button", { name: /Sem vínculo/ })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /Abaixo do custo/ })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /Sem evidência/ })).not.toBeInTheDocument();
  });
});
