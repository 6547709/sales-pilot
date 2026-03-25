"use client";

import { Suspense, useCallback } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { ProductSpotlight } from "@/components/marketing/product-spotlight";
import { SolutionAtlas } from "@/components/marketing/solution-atlas";
import type { TopologyCategoryNode } from "@/lib/api";

function scrollToSpotlight() {
  requestAnimationFrame(() => {
    document.getElementById("solution-spotlight")?.scrollIntoView({
      behavior: "smooth",
      block: "start",
    });
  });
}

function HomeContent() {
  const sp = useSearchParams();
  const router = useRouter();
  const urlSol = sp.get("solution") ?? "";
  const urlKw = sp.get("kw") ?? "";
  const urlMarketRaw = sp.get("market") ?? "";
  const urlMarket =
    urlMarketRaw === "domestic" || urlMarketRaw === "foreign"
      ? urlMarketRaw
      : "";

  const handleSelectSolution = useCallback(
    (cat: TopologyCategoryNode | null) => {
      if (!cat) {
        router.push("/", { scroll: false });
        scrollToSpotlight();
        return;
      }
      const idStr = String(cat.id);
      if (urlSol === idStr && !urlMarket) {
        router.push("/", { scroll: false });
      } else if (urlSol === idStr && urlMarket) {
        router.push(`/?solution=${idStr}`, { scroll: false });
      } else {
        router.push(`/?solution=${idStr}`, { scroll: false });
      }
      scrollToSpotlight();
    },
    [router, urlSol, urlMarket],
  );

  return (
    <>
      <SolutionAtlas
        initialSolutionId={urlSol}
        initialKeyword={urlKw}
        initialMarket={urlMarket}
        onSelectSolution={handleSelectSolution}
      />
      <ProductSpotlight
        key={`${urlSol}-${urlKw}-${urlMarket}`}
        solutionId={urlSol}
        keyword={urlKw}
        vendorMarket={urlMarket || undefined}
      />
    </>
  );
}

export default function HomePage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-[40vh] items-center justify-center text-muted-foreground">
          加载中…
        </div>
      }
    >
      <HomeContent />
    </Suspense>
  );
}
