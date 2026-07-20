import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { Header } from "./Header";
import { InstallationProvider } from "./InstallationContext";

const { listIntegrationInstallations } = vi.hoisted(() => ({
  listIntegrationInstallations: vi.fn(),
}));

vi.mock("./ClientContext", () => ({
  useClient: () => ({ listIntegrationInstallations }),
}));

function renderHeader(initialEntry = "/") {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[initialEntry]}>
        <InstallationProvider>
          <Header />
        </InstallationProvider>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("Header", () => {
  beforeEach(() => {
    listIntegrationInstallations.mockReset();
    listIntegrationInstallations.mockResolvedValue({ items: [] });
  });

  it("renders the six navigation pills in canonical order", () => {
    renderHeader();

    const labels = ["Visão geral", "Anúncios", "Mercado", "Simulador", "Pedidos", "Repasses"];
    const navigation = screen.getByRole("navigation");
    // Each pill is a direct child of <nav> (enabled = <a>, disabled = <span>); the pill's first
    // child node is the label text node (the "em breve" badge is a separate nested span, so
    // filtering by full textContent would miss the disabled pills).
    const renderedLabels = Array.from(navigation.children).map(
      (pill) => pill.childNodes[0]?.textContent?.trim() ?? "",
    );

    expect(renderedLabels).toEqual(labels);
  });

  it("keeps Repasses visible but non-navigable", () => {
    renderHeader();

    const navigation = screen.getByRole("navigation");
    const pill = Array.from(navigation.children).find((el) =>
      el.textContent?.includes("Repasses"),
    ) as HTMLElement | undefined;
    expect(pill).toBeDefined();
    expect(pill!.tagName).toBe("SPAN"); // a disabled stub, never an <a> link
    expect(pill!.querySelector("a")).toBeNull();
    expect(within(pill!).getByText("em breve")).toBeInTheDocument();
  });

  it("renders enabled pills as links", () => {
    renderHeader();

    for (const label of ["Visão geral", "Anúncios", "Mercado", "Simulador", "Pedidos"]) {
      expect(screen.getByRole("link", { name: label })).toBeInTheDocument();
    }
  });

  it.each([
    ["/precos", "Simulador", "Anúncios"],
    ["/", "Visão geral", "Anúncios"],
  ])("derives the active pill from the router at %s", (entry, active, inactive) => {
    renderHeader(entry);

    expect(screen.getByRole("link", { name: active })).toHaveClass(
      "bg-accent-soft",
      "text-accent-ink",
    );
    expect(screen.getByRole("link", { name: inactive })).not.toHaveClass("bg-accent-soft");
  });

  it("exposes the exact gear menu entries and DIFAL target", () => {
    renderHeader();

    // The gear is a native <summary aria-label="Abrir menu"> — jsdom does not expose <summary>
    // as role "button", so target it by its glyph to toggle the <details> open.
    fireEvent.click(screen.getByText("⚙"));

    const menu = screen.getByRole("group", { name: "Menu de configurações" });
    expect(within(menu).getByText("Configurações")).toBeInTheDocument();
    expect(within(menu).getByRole("link", { name: "DIFAL" })).toHaveAttribute(
      "href",
      "/precos?params=1",
    );
    expect(within(menu).getByRole("link", { name: "Integrações" })).toBeInTheDocument();
    expect(within(menu).getByRole("link", { name: "Catálogo" })).toBeInTheDocument();
    expect(within(menu).getByRole("link", { name: "Estoque" })).toBeInTheDocument();
    expect(screen.queryByText("Vínculos")).not.toBeInTheDocument();
  });

  it("wires the theme control to useTheme", () => {
    renderHeader();

    expect(() => fireEvent.click(screen.getByRole("button", { name: "Alternar tema" }))).not.toThrow();
  });
});
