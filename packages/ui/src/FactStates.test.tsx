import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { formatRelativeAge } from "@marketplace-central/web-query";
import { ConflictTag } from "./ConflictTag";
import { FreshnessIndicator } from "./FreshnessIndicator";
import { UnknownValue } from "./UnknownValue";

describe("fact state components", () => {
  it("renders an unknown value and only adds a non-empty hint tooltip", () => {
    const { rerender } = render(<UnknownValue hint="x" />);

    expect(screen.getByText("—")).toHaveAttribute("title", "x");
    expect(screen.getByText("—")).toHaveClass("text-faint");

    rerender(<UnknownValue />);

    expect(screen.getByText("—")).not.toHaveAttribute("title");

    rerender(<UnknownValue hint="" />);

    expect(screen.getByText("—")).not.toHaveAttribute("title");
  });

  it("renders the divergent conflict tag with its optional tooltip", () => {
    const { rerender } = render(<ConflictTag detail="ERP=35" />);
    const tag = screen.getByText("divergente");

    expect(tag.className).toContain("amber");
    expect(tag).toHaveClass("bg-amber-soft", "text-amber");
    expect(tag).toHaveAttribute("title", "ERP=35");

    rerender(<ConflictTag />);

    expect(screen.getByText("divergente")).not.toHaveAttribute("title");

    rerender(<ConflictTag detail="" />);

    expect(screen.getByText("divergente")).not.toHaveAttribute("title");
  });

  it("formats freshness through the pt-BR time representation", () => {
    const asOf = "2024-01-02T03:04:05-03:00";
    const now = new Date().getTime();
    const expectedAge = formatRelativeAge(asOf, now);

    render(<FreshnessIndicator asOf={asOf} />);

    expect(screen.getByLabelText("Atualização dos dados")).toHaveTextContent(expectedAge);
    expect(screen.getByLabelText("Atualização dos dados")).toHaveClass("text-muted");
  });

  it.each([null, undefined])("renders unknown freshness for %s", (asOf) => {
    render(<FreshnessIndicator asOf={asOf} />);

    expect(screen.getByLabelText("Atualização dos dados")).toHaveTextContent("idade desconhecida");
  });
});
