// These types mirror the JSON manifest the Go API returns from POST /v1/assets.
// Keep them in sync with the Go domain.Manifest JSON tags.

export type ProcessingStatus = "completed" | "partial" | "failed";

export interface Variant {
  name: string;
  width: number;
  height: number;
  format: string;
  path: string;
  size_bytes: number;
}

export interface ProcessingError {
  profile: string;
  stage: string;
  message: string;
}

export interface Manifest {
  asset_id: string;
  original_filename: string;
  status: ProcessingStatus;
  variants: Variant[];
  errors: ProcessingError[];
}

// The per-file state the browser tracks for each selected file.
export type FileStatus =
  | "pending"
  | "uploading"
  | "processing"
  | "completed"
  | "failed";

export interface QueuedFile {
  id: string; // stable client-side id
  file: File;
  status: FileStatus;
  manifest?: Manifest;
  error?: string;
}
