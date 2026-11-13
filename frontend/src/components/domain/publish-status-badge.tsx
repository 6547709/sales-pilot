import { Badge } from "@/components/reui/badge";

export function PublishStatusBadge({ isDraft }: { isDraft: boolean }) {
  return isDraft ? (
    <Badge variant="warning-light">草稿</Badge>
  ) : (
    <Badge variant="success-light">已发布</Badge>
  );
}
