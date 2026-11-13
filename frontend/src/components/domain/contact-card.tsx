import { Frame, FramePanel } from "@/components/reui/frame";

type ContactCardProps = {
  kind: "sales" | "presales";
  name?: string;
  phone?: string;
  email?: string;
};

export function ContactCard({ kind, name, phone, email }: ContactCardProps) {
  if (!name && !phone && !email) return null;
  return (
    <Frame>
      <FramePanel>
        <p className="text-xs text-muted-foreground">
          {kind === "sales" ? "销售联系人" : "售前联系人"}
        </p>
        <p className="mt-1 text-sm font-medium">{name || "—"}</p>
        <div className="mt-1 flex flex-col gap-0.5 text-xs text-muted-foreground">
          {phone ? (
            <a href={`tel:${phone}`} className="hover:text-foreground">
              {phone}
            </a>
          ) : null}
          {email ? (
            <a href={`mailto:${email}`} className="hover:text-foreground">
              {email}
            </a>
          ) : null}
        </div>
      </FramePanel>
    </Frame>
  );
}
