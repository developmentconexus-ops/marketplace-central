import { render, screen } from "@testing-library/react";
import { Badge } from "./Badge";

describe("Badge", () => {
  it("renders pending status", () => {
    render(<Badge status="pending" />);
    const badge = screen.getByText("Pending");
    expect(badge).toBeInTheDocument();
    expect(badge).toHaveClass("bg-surface-2", "text-muted");
  });

  it("renders succeeded status", () => {
    render(<Badge status="succeeded" />);
    const badge = screen.getByText("Succeeded");
    expect(badge).toBeInTheDocument();
    expect(badge).toHaveClass("bg-accent-soft", "text-accent-ink");
  });

  it("renders failed status", () => {
    render(<Badge status="failed" />);
    const badge = screen.getByText("Failed");
    expect(badge).toBeInTheDocument();
    expect(badge).toHaveClass("bg-warn-soft", "text-warn");
  });

  it("renders in_progress status", () => {
    render(<Badge status="in_progress" />);
    const badge = screen.getByText("In Progress");
    expect(badge).toBeInTheDocument();
    expect(badge).toHaveClass("bg-info-soft", "text-info");
  });

  it("renders completed status", () => {
    render(<Badge status="completed" />);
    const badge = screen.getByText("Completed");
    expect(badge).toBeInTheDocument();
    expect(badge).toHaveClass("bg-accent-soft", "text-accent-ink");
  });
});
