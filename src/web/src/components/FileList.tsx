import type { QueuedFile } from "../types/manifest";

// FileList renders selected files with their status and results.
//
// STORY 14: show filename + size + a remove button (when pending).
// STORY 17: show variants for completed files, errors for failed ones, and a
// retry button per failed file. Make "partial" visually distinct.
export function FileList(_props: {
  files: QueuedFile[];
  onRemove: (id: string) => void;
  onRetry: (id: string) => void;
}) {
  // TODO: map over files and render each row.
  return <ul />;
}
