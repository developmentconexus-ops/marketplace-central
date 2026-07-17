import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { ConflictTag } from "./ConflictTag";
import { FreshnessIndicator } from "./FreshnessIndicator";
import { UnknownValue } from "./UnknownValue";

describe("fact state components", () => {
  it("renders an unknown value and only adds a non-empty hint tooltip", () => {
    const { rerender } = render(<UnknownValue hint="x" />);

    expect(screen.getByText("—")).toHaveAttribute("title", "x");

    rerender(<UnknownValue />);

    expect(screen.getByText("—")).not.toHaveAttribute("title");

    rerender(<UnknownValue hint="" />);

    expect(screen.getByText("—")).not.toHaveAttribute("title");
  });

  it("renders the divergent conflict tag with its optional tooltip", () => {
    const { rerender } = render(<ConflictTag detail="ERP=35" />);
    const tag = screen.getByText("divergente");

    expect(tag.className).toContain("amber");
    expect(tag).toHaveAttribute("title", "ERP=35");

    rerender(<ConflictTag />);

    expect(screen.getByText("divergente")).not.toHaveAttribute("title");

    rerender(<ConflictTag detail="" />);

    expect(screen.getByText("divergente")).not.toHaveAttribute("title");
  });

  it("formats freshness through the pt-BR time representation", () => {
    const asOf = "2024-01-02T03:04:05-03:00";
    const expectedTime = new Date(asOf).toLocaleTimeString("pt-BR", {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    });

    render(<FreshnessIndicator asOf={asOf} />);

    expect(screen.getByLabelText("Atualização dos dados")).toHaveTextContent(
      `dados de ${expectedTime}`,
    );
  });

  it.each([null, undefined])("renders unknown freshness for %s", (asOf) => {
    render(<FreshnessIndicator asOf={asOf} />);

    expect(screen.getByLabelText("Atualização dos dados")).toHaveTextContent(
      "dados de desconhecido",
    );
  });
});
