import { describe, expect, it } from "vitest";
import { ptBrMoneyToDot, ptBrRateToDot } from "./ptbrDecimal";

describe("ptBrMoneyToDot", () => {
  it("promotes a pt-BR comma-decimal to dot-decimal (F-P7-2 core case)", () => {
    expect(ptBrMoneyToDot("150,00")).toBe("150.00");
  });

  it("strips thousands dots only when a comma is present", () => {
    expect(ptBrMoneyToDot("1.500,90")).toBe("1500.90");
    expect(ptBrMoneyToDot("1.234.567,89")).toBe("1234567.89");
  });

  it("passes a no-comma dot-decimal through untouched (backend / scenario values)", () => {
    // Guards the scenario-reload path: "77.00" must NOT become "7700".
    expect(ptBrMoneyToDot("77.00")).toBe("77.00");
    expect(ptBrMoneyToDot("89.00")).toBe("89.00");
  });

  it("passes a bare integer through", () => {
    expect(ptBrMoneyToDot("1500")).toBe("1500");
  });

  it("trims and treats empty as empty", () => {
    expect(ptBrMoneyToDot("")).toBe("");
    expect(ptBrMoneyToDot("  ")).toBe("");
    expect(ptBrMoneyToDot("  150,00  ")).toBe("150.00");
  });

  it("refuses to guess ambiguous / non-pt-BR shapes (ADR-17: honest 422 over silent mis-scale)", () => {
    // US comma-thousands + dot-decimal: a dot AFTER the comma ⇒ not pt-BR. Must
    // NOT become "1.23456" — pass raw so the backend rejects it honestly.
    expect(ptBrMoneyToDot("1,234.56")).toBe("1,234.56");
    // Dot after the last comma, general case.
    expect(ptBrMoneyToDot("12,5.0")).toBe("12,5.0");
    // Multiple commas = malformed ⇒ raw passthrough, never a fabricated number.
    expect(ptBrMoneyToDot("1,2,3")).toBe("1,2,3");
  });
});

describe("ptBrRateToDot", () => {
  it("promotes a pt-BR comma-decimal rate to dot-decimal", () => {
    expect(ptBrRateToDot("9,25")).toBe("9.25");
  });

  it("passes an integer or dot-decimal rate through untouched", () => {
    expect(ptBrRateToDot("18")).toBe("18");
    expect(ptBrRateToDot("20.5")).toBe("20.5");
  });

  it("trims and treats empty as empty", () => {
    expect(ptBrRateToDot("")).toBe("");
    expect(ptBrRateToDot("  40  ")).toBe("40");
  });
});
