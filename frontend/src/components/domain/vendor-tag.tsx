export function VendorTag({ name, logo }: { name?: string; logo?: string }) {
  if (!name?.trim()) return null;
  return (
    <span className="inline-flex max-w-full items-center gap-1.5 rounded-md border border-border bg-background px-2 py-1 text-xs">
      {logo ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={logo} alt="" className="h-4 w-4 shrink-0 object-contain" />
      ) : null}
      <span className="truncate">{name}</span>
    </span>
  );
}
