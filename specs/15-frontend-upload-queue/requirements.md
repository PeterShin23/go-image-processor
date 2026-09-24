# Story 15 — Build the Frontend Upload Queue

**Time budget:** ~1 hour (Day 9)

## Goal

As a user, I want selected files processed one at a time so I can see the result of each
upload. One HTTP request per image — no batch endpoint.

## Files to work in (`/src/web`)

- `src/api/client.ts` — implement `uploadImage(file)` with `fetch` + `FormData`
- `src/hooks/useUploadQueue.ts` — implement sequential `start()`
- `src/App.tsx` — add a "Process" button and a "Clear" button

## Per-file lifecycle

For each file: `pending` → `uploading` (→ `processing`) → `completed` | `failed`.
You may collapse `uploading`/`processing` into one since the API is synchronous.

## Tasks

- [ ] Implement `uploadImage`: `FormData` with field `image`, `POST /v1/assets`, throw on
      non-2xx, return the parsed `Manifest`.
- [ ] Implement `start()`: iterate files sequentially; mark each `uploading`, call
      `uploadImage`, then set `completed` (store the manifest) or `failed` (store the error);
      continue to the next regardless.
- [ ] A failed upload must NOT stop the remaining queue.
- [ ] Guard against starting the same queue twice (disable/ignore while running).
- [ ] Add a `clear()` to reset after completion.

## Acceptance criteria

- [ ] Every selected file has independent state.
- [ ] Files are processed sequentially.
- [ ] A failed upload does not stop the rest of the queue.
- [ ] Completed files store their manifest.
- [ ] Failed files display a useful error.
- [ ] The same queue cannot be started twice.
- [ ] The queue can be cleared after completion.
- [ ] `npm run build` succeeds.

## Definition of done

Selecting 3 images and clicking Process uploads them one by one; each row ends
`completed` or `failed`. Commit: `story 15: sequential upload queue`.

## Hints (not solutions)

- Update one file immutably by id:
  `setFiles(prev => prev.map(f => f.id === id ? { ...f, status } : f))`.
- Sequential = `for (const f of files) { await processOne(f); }` inside an async function.
- Track running with a `useState<boolean>` or a `useRef` flag to block double-start.
- `try/catch` around `uploadImage`; the `catch` sets `failed` with `err.message`.
