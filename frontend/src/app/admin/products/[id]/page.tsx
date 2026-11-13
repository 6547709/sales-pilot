"use client";

import { useParams } from "next/navigation";
import { ProductEditor } from "@/features/products/product-editor";

export default function AdminProductEditPage() {
  const params = useParams();
  return <ProductEditor productId={String(params.id)} />;
}
