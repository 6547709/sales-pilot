import Link from "next/link";
import { Frame, FramePanel } from "@/components/reui/frame";
import { Badge } from "@/components/reui/badge";
import { MarketBadge } from "@/components/domain/market-badge";
import { VendorTag } from "@/components/domain/vendor-tag";
import type { Product } from "@/lib/api";

type ProductCardProps = {
  product: Product;
  href: string;
};

export function ProductCard({ product, href }: ProductCardProps) {
  return (
    <Link href={href} className="block h-full">
      <Frame className="h-full transition-colors hover:border-primary/40">
        <FramePanel>
          <div className="flex flex-col gap-2">
            <div className="flex items-start justify-between gap-2">
              <p className="min-w-0 truncate font-medium">{product.name}</p>
              <MarketBadge market={product.vendorMarket} />
            </div>
            <p className="line-clamp-2 min-h-8 text-xs text-muted-foreground">
              {product.description}
            </p>
            <div className="flex flex-wrap items-center gap-1.5">
              <VendorTag name={product.manufacturerName} logo={product.manufacturerLogo} />
              {product.highlights.slice(0, 2).map((h) => (
                <Badge key={h} variant="outline">
                  {h}
                </Badge>
              ))}
            </div>
          </div>
        </FramePanel>
      </Frame>
    </Link>
  );
}
