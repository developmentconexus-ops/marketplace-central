import { describe, expect, it } from "vitest";
import { resolveApiBaseUrl } from "./ClientContext";

describe("resolveApiBaseUrl", () => {
  it("prefers an explicit environment base URL", () => {
    expect(resolveApiBaseUrl("http://api.internal", true)).toBe("http://api.internal");
  });

  it("uses localhost backend by default in dev", () => {
    expect(resolveApiBaseUrl(undefined, true)).toBe("http://localhost:8080");
  });

  it("keeps same-origin requests in non-dev builds when no override exists", () => {
    expect(resolveApiBaseUrl(undefined, false)).toBe("");
  });
});
