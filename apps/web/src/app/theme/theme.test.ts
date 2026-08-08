import { act, renderHook } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import { applyTheme, getInitialTheme, persistTheme, readStoredTheme, STORAGE_KEY } from "./theme";
import { useTheme } from "./useTheme";

// jsdom under this repo's vitest env (jsdom 26 / vitest 3.2.4) does not instantiate Web Storage:
// window.localStorage is undefined (env-verified 2026-07-17). Install a spec-shaped, Map-backed
// localStorage whose methods live on Storage.prototype so vi.spyOn(Storage.prototype, ...) below
// still intercepts them. No-op in a real browser env where localStorage already exists.
if (typeof globalThis.localStorage === "undefined") {
  const store = new Map<string, string>();
  const def = (name: string, value: unknown) =>
    Object.defineProperty(Storage.prototype, name, { configurable: true, writable: true, value });
  def("getItem", (key: string) => (store.has(key) ? store.get(key)! : null));
  def("setItem", (key: string, value: string) => void store.set(key, String(value)));
  def("removeItem", (key: string) => void store.delete(key));
  def("clear", () => store.clear());
  def("key", (index: number) => Array.from(store.keys())[index] ?? null);
  const localStorage = Object.create(Storage.prototype) as Storage;
  Object.defineProperty(globalThis, "localStorage", { configurable: true, value: localStorage });
  Object.defineProperty(window, "localStorage", { configurable: true, value: localStorage });
}

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
