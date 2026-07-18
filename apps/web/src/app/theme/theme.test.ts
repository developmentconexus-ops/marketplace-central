import { act, renderHook } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import {
  applyTheme,
  getInitialTheme,
  persistTheme,
  readStoredTheme,
  STORAGE_KEY,
} from "./theme";
import { useTheme } from "./useTheme";

describe("theme storage", () => {
  afterEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    document.documentElement.removeAttribute("data-theme");
  });

  it.each([
    ["dark", "dark"],
    ["light", "light"],
    [null, "light"],
    ["blue", "light"],
    ["", "light"],
  ] as const)("reads %s as %s", (stored, expected) => {
    if (stored !== null) localStorage.setItem(STORAGE_KEY, stored);

    expect(readStoredTheme()).toBe(expected);
  });

  it("defaults to light when storage cannot be read", () => {
    vi.spyOn(Storage.prototype, "getItem").mockImplementation(() => {
      throw new DOMException("denied", "SecurityError");
    });

    expect(readStoredTheme()).toBe("light");
  });

  it("does not throw when storage cannot be written", () => {
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {
      throw new DOMException("full", "QuotaExceededError");
    });

    expect(() => persistTheme("dark")).not.toThrow();
  });

  it("prefers the bootstrapped root theme", () => {
    localStorage.setItem(STORAGE_KEY, "light");
    applyTheme("dark");

    expect(getInitialTheme()).toBe("dark");
  });
});

describe("useTheme", () => {
  afterEach(() => {
    localStorage.clear();
    document.documentElement.removeAttribute("data-theme");
  });

  it("applies a change immediately and retains it on re-initialization", () => {
    applyTheme("light");
    const first = renderHook(() => useTheme());

    act(() => first.result.current.setTheme("dark"));

    expect(document.documentElement.getAttribute("data-theme")).toBe("dark");
    expect(getInitialTheme()).toBe("dark");
    first.unmount();

    const second = renderHook(() => useTheme());
    expect(second.result.current.theme).toBe("dark");
  });
});
