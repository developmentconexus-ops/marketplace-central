import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { MarketplaceIcon } from "./MarketplaceIcon";

describe("MarketplaceIcon", () => {
  it("renders correct initial for known marketplace code", () => {
    render(<MarketplaceIcon code="mercado_livre" />);
    expect(screen.getByText("M")).toBeInTheDocument();
  });

  it("applies Mercado Livre brand color as background", () => {
    const { container } = render(<MarketplaceIcon code="mercado_livre" />);
    const el = container.firstChild as HTMLElement;
    expect(el.style.backgroundColor).toBe("rgb(255, 230, 0)"); // #FFE600
  });

  it("applies amazon brand color as background", () => {
    const { container } = render(<MarketplaceIcon code="amazon" />);
    const el = container.firstChild as HTMLElement;
    expect(el.style.backgroundColor).toBe("rgb(17, 24, 39)");
    expect(screen.getByText("A")).toBeInTheDocument();
  });

  it("renders unknown code with default color and correct initial", () => {
    const { container } = render(<MarketplaceIcon code="unknown_mp" />);
    const el = container.firstChild as HTMLElement;
    expect(el.style.backgroundColor).toBe("rgb(99, 102, 241)"); // #6366F1
    expect(screen.getByText("U")).toBeInTheDocument();
  });

  it("uses provided size", () => {
    const { container } = render(<MarketplaceIcon code="mercado_livre" size={48} />);
    const el = container.firstChild as HTMLElement;
    expect(el.style.width).toBe("48px");
    expect(el.style.height).toBe("48px");
  });
});

