import { describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen } from "@testing-library/react";
import { DataTable, DetailDrawer, MarginChip } from "./index";

describe("DetailDrawer", () => {
  it("renders the panel content, subtitle, and actions when open", () => {
    render(
      <DetailDrawer
        open
        onClose={vi.fn()}
        title="Detalhes do anúncio"
        subtitle="SKU-001"
        actions={<button>Salvar</button>}
      >
        <p>Conteúdo do detalhe</p>
      </DetailDrawer>,
    );

    expect(screen.getByRole("complementary", { name: "Detalhes do anúncio" })).toBeInTheDocument();
    expect(screen.getByText("SKU-001")).toBeInTheDocument();
    expect(screen.getByText("Conteúdo do detalhe")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Salvar" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Salvar" }).closest(".border-t")).not.toBeNull();
  });

  it("renders nothing when closed", () => {
    const { container } = render(
      <DetailDrawer open={false} onClose={vi.fn()} title="Detalhes do anúncio">
        <p>Conteúdo do detalhe</p>
      </DetailDrawer>,
    );

    expect(container).toBeEmptyDOMElement();
  });

  it("uses the default close label and delegates close interactions", () => {
    const onClose = vi.fn();
    render(
      <DetailDrawer open onClose={onClose} title="Detalhes do anúncio">
        <p>Conteúdo do detalhe</p>
      </DetailDrawer>,
    );

    fireEvent.click(screen.getByRole("button", { name: "Fechar" }));
    expect(onClose).toHaveBeenCalledOnce();

    fireEvent.keyDown(document, { key: "Escape" });
    expect(onClose).toHaveBeenCalledTimes(2);
  });

  it("uses the canonical width by default and accepts an override", () => {
    const { rerender } = render(
      <DetailDrawer open onClose={vi.fn()} title="Detalhes do anúncio">
        <p>Conteúdo do detalhe</p>
      </DetailDrawer>,
    );

    expect(screen.getByRole("complementary")).toHaveStyle({ width: "360px" });

    rerender(
      <DetailDrawer open onClose={vi.fn()} title="Detalhes do anúncio" width={320}>
        <p>Conteúdo do detalhe</p>
      </DetailDrawer>,
    );

    expect(screen.getByRole("complementary")).toHaveStyle({ width: "320px" });
  });

  it("inherits the semantic surface token from DetailPanel", () => {
    render(
      <DetailDrawer open onClose={vi.fn()} title="Detalhes do anúncio">
        <p>Conteúdo do detalhe</p>
      </DetailDrawer>,
    );

    expect(screen.getByRole("complementary")).toHaveClass("bg-surface", "border-border");
  });

  it("exports the F-03 primitives from the barrel", () => {
    expect(typeof MarginChip).toBe("function");
    expect(typeof DataTable).toBe("function");
    expect(typeof DetailDrawer).toBe("function");
  });
});
