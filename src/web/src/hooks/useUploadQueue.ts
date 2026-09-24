import { useState } from "react";
import type { QueuedFile } from "../types/manifest";

// useUploadQueue owns the list of selected files and their per-file state.
// Keep state LOCAL and simple — no Redux/Zustand.
//
// STORY 15: sequential processing (one at a time).
// STORY 16: bounded concurrency (default 2 at once).
// STORY 17: expose retry for a single failed file.
export function useUploadQueue() {
  const [files, setFiles] = useState<QueuedFile[]>([]);

  // TODO (story 14): addFiles(FileList) -> append as "pending"
  // TODO (story 14): removeFile(id) -> drop before processing
  // TODO (story 15): start() -> process the queue; guard against double-start
  // TODO (story 16): add a concurrency limit (default 2)
  // TODO (story 17): retry(id) -> re-run only that file
  // TODO: clear() -> reset after completion

  return { files, setFiles };
}
