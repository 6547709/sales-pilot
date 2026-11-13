import { describe, expect, it } from "vitest";
import { toLocalUser, toPageResult, toProduct, toProductDTO } from "./adapters";

describe("toProduct", () => {
  it("maps snake_case DTO to camelCase model and parses json fields", () => {
    const p = toProduct({
      id: 1,
      name: "A",
      category: "c",
      description: "",
      highlights: '["h1","h2"]',
      target_personas: { x: "y" },
      discovery_questions: "",
      trigger_events: "",
      competitor_analysis: "",
      roi_metrics: "",
      manufacturer_name: "M",
      vendor_market: "domestic",
      is_draft: true,
      updated_at: "2026-01-01T00:00:00Z",
      created_at: "2026-01-01T00:00:00Z",
    });
    expect(p.manufacturerName).toBe("M");
    expect(p.vendorMarket).toBe("domestic");
    expect(p.isDraft).toBe(true);
    expect(p.highlights).toEqual(["h1", "h2"]);
    expect(p.targetPersonas).toEqual({ x: "y" });
  });

  it("round-trips back to the DTO shape the API expects", () => {
    const dto = toProductDTO(
      toProduct({
        id: 2,
        name: "B",
        category: "",
        description: "",
        highlights: ["h"],
        target_personas: {},
        discovery_questions: ["q"],
        trigger_events: "",
        competitor_analysis: "",
        roi_metrics: "",
        manufacturer_name: "M2",
        updated_at: "2026-01-01T00:00:00Z",
      }),
    );
    expect(dto.manufacturer_name).toBe("M2");
    expect(dto.highlights).toEqual(["h"]);
    expect(dto.discovery_questions).toEqual(["q"]);
  });
});

describe("toLocalUser", () => {
  it("keeps active flag and iso times", () => {
    const u = toLocalUser({
      id: 2,
      username: "u",
      email: "e",
      role: "admin",
      is_active: false,
      created_at: "2026-01-01T00:00:00Z",
      updated_at: "2026-01-01T00:00:00Z",
    });
    expect(u.isActive).toBe(false);
    expect(u.createdAt).toBe("2026-01-01T00:00:00Z");
  });
});

describe("toPageResult", () => {
  it("wraps a bare array as a single page", () => {
    const r = toPageResult<{ id: number }, number>(
      [{ id: 1 }],
      2,
      10,
      (d) => d.id,
    );
    expect(r.items).toEqual([1]);
    expect(r.total).toBe(1);
    expect(r.page).toBe(1);
  });

  it("reads the paginated envelope", () => {
    const r = toPageResult<{ id: number }, number>(
      { items: [{ id: 1 }, { id: 2 }], total: 42, page: 3, page_size: 2 },
      1,
      10,
      (d) => d.id,
    );
    expect(r.items).toEqual([1, 2]);
    expect(r.total).toBe(42);
    expect(r.page).toBe(3);
    expect(r.pageSize).toBe(2);
  });
});
