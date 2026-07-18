import { fireEvent, render, screen } from "@testing-library/react";
import { Button } from "./Button";
import { SurfaceCard } from "./SurfaceCard";
import { StatCard } from "./StatCard";

describe("Button", () => {
  it("renders children and uses semantic variant tokens", () => {
    render(<Button>Secondary</Button>);
    const secondary = screen.getByRole("button", { name: "Secondary" });
    expect(secondary).toHaveClass("bg-surface", "text-ink", "border-border", "rounded-control");

    render(<Button variant="primary">Primary</Button>);
    expect(screen.getByRole("button", { name: "Primary" })).toHaveClass("bg-accent", "text-white");

    render(<Button variant="danger">Danger</Button>);
    expect(screen.getByRole("button", { name: "Danger" })).toHaveClass("bg-warn", "text-white");
  });

  it("disables while loading and renders the spinner", () => {
    render(<Button loading>Loading</Button>);
    const button = screen.getByRole("button", { name: "Loading" });
    expect(button).toBeDisabled();
    expect(button.querySelector("svg")).toBeInTheDocument();
  });

  it("disables when disabled is passed", () => {
    render(<Button disabled>Disabled</Button>);
    expect(screen.getByRole("button", { name: "Disabled" })).toBeDisabled();
  });

  it("only fires onClick when enabled", () => {
    const onClick = vi.fn();
    const { rerender } = render(<Button onClick={onClick}>Action</Button>);
    const button = screen.getByRole("button", { name: "Action" });
    fireEvent.click(button);
    expect(onClick).toHaveBeenCalledTimes(1);

    rerender(
      <Button disabled onClick={onClick}>
        Action
      </Button>,
    );
    fireEvent.click(screen.getByRole("button", { name: "Action" }));
    expect(onClick).toHaveBeenCalledTimes(1);

    rerender(
      <Button loading onClick={onClick}>
        Action
      </Button>,
    );
    fireEvent.click(screen.getByRole("button", { name: "Action" }));
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});

describe("SurfaceCard", () => {
  it("renders children with semantic surface tokens and merges className", () => {
    render(<SurfaceCard className="custom-class">Content</SurfaceCard>);
    const section = screen.getByText("Content").closest("section");
    expect(section).toHaveClass("bg-surface", "border-border", "rounded-card", "custom-class");
  });
});

describe("StatCard", () => {
  it("renders label, value, and optional sub with honest value passthrough", () => {
    const { rerender } = render(<StatCard label="Status" value="—" sub="Unavailable" />);
    const value = screen.getByText("—");
    expect(value).toHaveClass("text-ink");
    expect(value).toHaveStyle({ fontFamily: "var(--font-mono)" });
    expect(screen.getByText("Status")).toHaveClass("text-muted");
    expect(screen.getByText("Unavailable")).toBeInTheDocument();

    rerender(<StatCard label="Count" value="1.234" />);
    expect(screen.getByText("1.234")).toBeInTheDocument();
    expect(screen.queryByText("Unavailable")).not.toBeInTheDocument();
  });
});
