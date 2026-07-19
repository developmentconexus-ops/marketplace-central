import { describe, expect, it } from "vitest";
import { applyProdutoTab, parseProdutoTab } from "./productQueryState";

describe("parseProdutoTab", () => {
  it("defaults to veredicto when the param is missing", () => {
    expect(parseProdutoTab(new URLSearchParams())).toBe("veredicto");
  });

  it("defaults to veredicto when the param is empty", () => {
    expect(parseProdutoTab(new URLSearchParams("tab="))).toBe("veredicto");
  });

  it("defaults to veredicto when the param is invalid", () => {
    expect(parseProdutoTab(new URLSearchParams("tab=bogus"))).toBe("veredicto");
  });

  it("round-trips anuncios", () => {
    expect(parseProdutoTab(new URLSearchParams("tab=anuncios"))).toBe("anuncios");
  });

  it("round-trips estoque", () => {
    expect(parseProdutoTab(new URLSearchParams("tab=estoque"))).toBe("estoque");
  });
});

describe("applyProdutoTab", () => {
  it("sets the tab param for a non-default tab", () => {
    const next = applyProdutoTab(new URLSearchParams(), "estoque");
    expect(next.get("tab")).toBe("estoque");
  });

  it("clears the tab param back to the default", () => {
    const next = applyProdutoTab(new URLSearchParams("tab=estoque"), "veredicto");
    expect(next.has("tab")).toBe(false);
  });
});
