import { describe, expect, it } from "vitest";
import {
  formatCollectedAt,
  formatMoney,
  formatPosition,
} from "./mercadoFormatters";

describe("mercadoFormatters", () => {
  it("formats a decimal-string amount as pt-BR currency", () => {
    expect(formatMoney({ amount: "229.20", currency: "BRL" })).toBe("R$ 229,20");
    expect(formatMoney({ amount: "1179.9", currency: "BRL" })).toBe("R$ 1.179,90");
  });

  it("returns null (honest dash upstream) for a missing/unparseable amount — never R$ NaN", () => {
    expect(formatMoney(null)).toBeNull();
    expect(formatMoney({ amount: "", currency: "BRL" })).toBeNull();
    expect(formatMoney({ amount: "abc", currency: "BRL" })).toBeNull();
  });

  it("formats position as rankº/total, null when absent", () => {
    expect(formatPosition({ rank: 9, total: 14 })).toBe("9º/14");
    expect(formatPosition(null)).toBeNull();
  });

  it("treats a zero-date / missing timestamp as an honest unknown", () => {
    expect(formatCollectedAt("0001-01-01T00:00:00Z")).toBeNull();
    expect(formatCollectedAt(null)).toBeNull();
    expect(formatCollectedAt("")).toBeNull();
    expect(formatCollectedAt("not-a-date")).toBeNull();
  });
});
