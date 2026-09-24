import type { Manifest } from "../types/manifest";

// uploadImage sends ONE file to the Go API and returns its manifest.
// The frontend sends one request per image — there is no batch endpoint.
//
// STORY 15: implement this with fetch + FormData. Throw on non-2xx so callers
// can mark the file failed.
export async function uploadImage(file: File): Promise<Manifest> {
  // TODO (story 15):
  //   const body = new FormData();
  //   body.append("image", file);
  //   const res = await fetch("/v1/assets", { method: "POST", body });
  //   if (!res.ok) throw new Error(`upload failed: ${res.status}`);
  //   return (await res.json()) as Manifest;
  throw new Error("uploadImage not implemented (story 15)");
}
