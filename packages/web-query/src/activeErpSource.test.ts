import { act, renderHook } from "@testing-library/react";
import { beforeEach, describe, expect, it } from "vitest";
import {
  DEFAULT_ERP_SOURCE,
  getActiveErpSource,
  setActiveErpSource,
  useActiveErpSource,
} from "./activeErpSource";

describe("activeErpSource", () => {
  beforeEach(() => {
    window.localStorage.clear();
  });

  it("defaults to xlsx (Sankhya) when nothing is stored", () => {
    expect(DEFAULT_ERP_SOURCE).toBe("xlsx");
    expect(getActiveErpSource()).toBe("xlsx");
  });

  it("ignores an unrecognized stored value and falls back to the default", () => {
    window.localStorage.setItem("mc.active-erp-source", "sankhya");
    expect(getActiveErpSource()).toBe("xlsx");
  });

  it("persists a selection and reads it back", () => {
    setActiveErpSource("catalogo_cliente");
    expect(getActiveErpSource()).toBe("catalogo_cliente");
    expect(window.localStorage.getItem("mc.active-erp-source")).toBe("catalogo_cliente");
  });

  it("updates every subscribed hook when the source changes", () => {
    const a = renderHook(() => useActiveErpSource());
    const b = renderHook(() => useActiveErpSource());
    expect(a.result.current[0]).toBe("xlsx");

    act(() => a.result.current[1]("catalogo_cliente"));

    expect(a.result.current[0]).toBe("catalogo_cliente");
    expect(b.result.current[0]).toBe("catalogo_cliente");
  });
});
