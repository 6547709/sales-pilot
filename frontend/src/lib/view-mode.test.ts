import { describe, expect, it } from "vitest";
import { withView } from "./view-mode-shared";

describe("withView", () => {
  it("appends view to a path without query", () => {
    expect(withView("/api/v1/products", "customer")).toBe("/api/v1/products?view=customer");
  });

  it("keeps existing query params", () => {
    expect(withView("/api/v1/products?q=a", "customer")).toBe("/api/v1/products?q=a&view=customer");
  });
});
