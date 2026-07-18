import { render, screen } from "@testing-library/react";
import { MarginChip } from "./MarginChip";

function chip() {
  return screen.getByText(/%|—/);
}

describe("MarginChip", () => {
  it("renders unknown margins as a neutral placeholder", () => {
    render(<MarginChip marginPct={null} />);

    expect(screen.getByText("—")).toBeInTheDocument();
    expect(screen.queryByText("0%")).not.toBeInTheDocument();
    expect(chip()).not.toHaveTextContent(/\d/);
    expect(chip()).toHaveClass("bg-surface-2", "text-muted");
    expect(chip()).not.toHaveClass("bg-accent-soft");
  });

  it.each([
    [25, "25%", "bg-accent-soft", "text-accent-ink"],
    [18, "18%", "bg-accent-soft", "text-accent-ink"],
    [12, "12%", "bg-amber-soft", "text-amber"],
    [3, "3%", "bg-warn-soft", "text-warn"],
  ] as const)("classifies %s as %s", (marginPct, text, background, foreground) => {
    render(<MarginChip marginPct={marginPct} />);

    expect(screen.getByText(text)).toHaveClass(background, foreground);
  });

  it.each([Number.NaN, Number.POSITIVE_INFINITY, Number.NEGATIVE_INFINITY])(
    "renders non-finite margin %s as unknown",
    (marginPct) => {
      render(<MarginChip marginPct={marginPct} />);

      expect(screen.getByText("—")).toBeInTheDocument();
      expect(screen.queryByText("0%")).not.toBeInTheDocument();
      expect(chip()).toHaveClass("bg-surface-2", "text-muted");
    },
  );

  it("uses valid custom thresholds", () => {
    render(<MarginChip marginPct={25} thresholds={{ healthy: 30, tight: 20 }} />);

    expect(screen.getByText("25%")).toHaveClass("bg-amber-soft", "text-amber");
    expect(screen.getByText("25%")).not.toHaveClass("bg-accent-soft");
  });

  it.each([
    { healthy: 10, tight: 20 },
    { healthy: Number.POSITIVE_INFINITY, tight: 10 },
  ])("falls back to defaults and warns for invalid thresholds", (thresholds) => {
    const warning = vi.spyOn(console, "warn").mockImplementation(() => {});

    render(<MarginChip marginPct={25} thresholds={thresholds} />);

    expect(screen.getByText("25%")).toHaveClass("bg-accent-soft", "text-accent-ink");
    expect(warning).toHaveBeenCalledWith(expect.stringContaining("MarginChip"));
  });
});
