"use client";

import { createContext, useContext, useState, useEffect, type ReactNode } from "react";
import { apiFetch } from "@/lib/api";

type FooterConfig = {
  footerText: string;
  filingNumber: string;
};

type FooterConfigContextType = FooterConfig & {
  loading: boolean;
};

const FooterConfigContext = createContext<FooterConfigContextType>({
  footerText: "",
  filingNumber: "",
  loading: true,
});

export function FooterConfigProvider({ children }: { children: ReactNode }) {
  const [config, setConfig] = useState<FooterConfig>({
    footerText: "",
    filingNumber: "",
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiFetch("/api/v1/config")
      .then((r) => r.json())
      .then((data: FooterConfig) => setConfig(data))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  return (
    <FooterConfigContext.Provider value={{ ...config, loading }}>
      {children}
    </FooterConfigContext.Provider>
  );
}

export function useFooterConfig() {
  return useContext(FooterConfigContext);
}
