import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { EmptyState } from "./EmptyState";
import { ErrorState } from "./ErrorState";
import { LoadingState } from "./LoadingState";

describe("load state components", () => {
  it("renders the loading copy", () => {
    render(<LoadingState />);

    expect(screen.getByText("Carregando…")).toBeInTheDocument();
  });

  it("renders the generic error and retries once", () => {
    const onRetry = vi.fn();
    render(<ErrorState onRetry={onRetry} />);

    expect(screen.getByRole("alert")).toHaveTextContent(/^Erro ao carregar\.$/);
    expect(screen.getByRole("button", { name: "Tentar novamente" })).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Tentar novamente" }));

    expect(onRetry).toHaveBeenCalledOnce();
  });

  it("appends error detail after the fixed prefix", () => {
    render(<ErrorState onRetry={vi.fn()} detail="Tente novamente mais tarde." />);

    expect(screen.getByRole("alert")).toHaveTextContent(
      "Erro ao carregar. Tente novamente mais tarde.",
    );
  });

  it("renders an empty hint only when provided", () => {
    const { container, rerender } = render(<EmptyState hint="Ajuste os filtros." />);

    expect(screen.getByText("Nenhum registro encontrado.")).toBeInTheDocument();
    expect(screen.getByText("Ajuste os filtros.")).toBeInTheDocument();
    expect(container.querySelectorAll("p")).toHaveLength(2);

    rerender(<EmptyState />);

    expect(screen.getByText("Nenhum registro encontrado.")).toBeInTheDocument();
    expect(container.querySelectorAll("p")).toHaveLength(1);
  });
});
