"use client";

import { createContext, useCallback, useContext, useEffect, useState } from "react";
import type { ViewMode } from "./view-mode-shared";

export type { ViewMode } from "./view-mode-shared";
export { withView } from "./view-mode-shared";

type ViewModeContextValue = {
  mode: ViewMode;
  setMode: (mode: ViewMode) => void;
};

const ViewModeContext = createContext<ViewModeContextValue>({
  mode: "sales",
  setMode: () => {},
});

function readInitialMode(): ViewMode {
  if (typeof window === "undefined") return "sales";
  const url = new URLSearchParams(window.location.search).get("view");
  if (url === "customer" || url === "sales") return url;
  const saved = window.localStorage.getItem("sp_view");
  return saved === "customer" ? "customer" : "sales";
}

export function ViewModeProvider({ children }: { children: React.ReactNode }) {
  const [mode, setModeState] = useState<ViewMode>("sales");

  useEffect(() => {
    setModeState(readInitialMode());
  }, []);

  const setMode = useCallback((next: ViewMode) => {
    setModeState(next);
    window.localStorage.setItem("sp_view", next);
    const url = new URL(window.location.href);
    url.searchParams.set("view", next);
    window.history.replaceState(null, "", url.toString());
  }, []);

  return (
    <ViewModeContext.Provider value={{ mode, setMode }}>{children}</ViewModeContext.Provider>
  );
}

export function useViewMode() {
  return useContext(ViewModeContext);
}
