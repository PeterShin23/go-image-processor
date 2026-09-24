# Story 14 — Initialize the React Application

**Time budget:** ~1.5 hours (Day 9)

## Goal

As a user, I want a simple browser interface where I can select multiple images, see
them listed, and remove any before processing.

## Tech (and only this)

React · TypeScript · Vite · native `fetch` · plain CSS. **Do not** add Redux, Zustand,
React Router, a component framework, or SSR.

## Files to work in (`/src/web`)

Scaffold is provided:
- `package.json`, `tsconfig.json`, `vite.config.ts`, `index.html`
- `src/main.tsx`, `src/App.tsx`
- `src/types/manifest.ts` — manifest + `QueuedFile` types (matches the Go JSON)
- `src/hooks/useUploadQueue.ts` — the queue hook (mostly TODOs)
- `src/components/FileList.tsx` — the list component (mostly TODOs)
- `src/api/client.ts` — `uploadImage` stub (implemented in story 15)

## Tasks

- [ ] `cd src/web && npm install` (or `make web-install`), then `npm run dev`; confirm the
      page loads at the Vite URL.
- [ ] Add a multi-file `<input type="file" multiple accept="image/jpeg,image/png">`.
- [ ] On change, add the selected files to `QueuedFile[]` state (status `"pending"`).
- [ ] Render each file's name and human-readable size in `FileList`.
- [ ] Add a remove button that drops a file (while it's still pending).
- [ ] Give each queued file a stable client id (e.g. `crypto.randomUUID()`).

## Acceptance criteria

- [ ] The frontend runs locally (`npm run dev`).
- [ ] Multiple files can be selected.
- [ ] Selected files appear in a list with name + size.
- [ ] Files can be removed before processing.
- [ ] The frontend does not transform images (it only selects/lists them).
- [ ] `npm run build` succeeds (types check).

## Definition of done

You can select several images, see them listed, and remove one. Commit:
`story 14: react file selection UI`.

## Hints (not solutions)

- `e.target.files` is a `FileList`; spread it: `Array.from(e.target.files)`.
- Format sizes with a tiny helper (KB/MB) — a `File` has `.name` and `.size`.
- Immutable updates: `setFiles(prev => [...prev, ...newOnes])` and
  `setFiles(prev => prev.filter(f => f.id !== id))`.
