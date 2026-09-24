# Background — Story 15: Build the Frontend Upload Queue

This story is about async JavaScript, immutable state updates over time, and modeling a
per-item lifecycle. These patterns are the heart of any "process a list with progress" UI.

## `fetch` + `FormData` = one upload

The browser's `fetch` sends HTTP requests. To send a file as multipart (what the Go API
expects from story 11), use `FormData`:

```ts
export async function uploadImage(file: File): Promise<Manifest> {
  const body = new FormData();
  body.append("image", file); // field name "image" — must match r.FormFile("image")
  const res = await fetch("/v1/assets", { method: "POST", body });
  if (!res.ok) {
    throw new Error(`upload failed: ${res.status}`);
  }
  return (await res.json()) as Manifest;
}
```

Key points:
- **Don't set `Content-Type` yourself.** When you pass a `FormData` body, the browser sets
  `multipart/form-data` *with the boundary* automatically. Setting it manually breaks the
  boundary — a classic bug.
- `/v1/assets` is relative; Vite's dev proxy forwards it to `:8080` (story 14).
- `res.ok` is true for 2xx. On anything else, throw so the caller marks the file failed.
- `res.json()` parses the body; the `as Manifest` cast trusts it matches your type (the
  field names must line up with the Go JSON tags).

## async/await and Promises

An `async` function returns a **Promise** — a value that arrives later. `await` pauses until
it resolves:

```ts
const manifest = await uploadImage(file); // waits for the server
```

`await` only works inside `async` functions. Under the hood it's non-blocking — the browser
stays responsive while waiting. Errors from an awaited call are thrown, so you catch them
with ordinary `try/catch`:

```ts
try {
  const m = await uploadImage(file);
  markCompleted(id, m);
} catch (err) {
  markFailed(id, err instanceof Error ? err.message : String(err));
}
```

(`err` is typed `unknown` in TS — narrow it before using `.message`.)

## Sequential processing

"One at a time" is just a `for...of` loop with `await` inside:

```ts
async function start() {
  for (const qf of files) {
    setStatus(qf.id, "uploading");
    try {
      const m = await uploadImage(qf.file);
      setFiles(prev => prev.map(f => f.id === qf.id
        ? { ...f, status: "completed", manifest: m } : f));
    } catch (err) {
      setFiles(prev => prev.map(f => f.id === qf.id
        ? { ...f, status: "failed", error: String(err) } : f));
    }
  }
}
```

Because each iteration `await`s, the next upload doesn't start until the current one
finishes — sequential by construction. The `try/catch` per item is what makes "one failure
doesn't stop the queue" true: a thrown error is handled locally and the loop continues.
Story 16 will relax this to run two at once.

## Immutable updates by id

You'll update one file inside a list many times (pending → uploading → completed). The
idiomatic move is `map` returning a new array, replacing just the matching item with a spread
copy:

```ts
setFiles(prev => prev.map(f => f.id === id ? { ...f, status } : f));
```

`{ ...f, status }` copies `f` and overrides `status`. Every other item is returned unchanged
(same reference). Never mutate `f` in place — React won't notice, and you'll get stale UI.

## Guarding against double-start

If the user clicks "Process" twice, you'd launch the queue twice. Track a running flag:

```ts
const running = useRef(false);
async function start() {
  if (running.current) return;
  running.current = true;
  try { /* ...loop... */ } finally { running.current = false; }
}
```

A `useRef` holds a mutable value that survives re-renders *without* causing one — perfect for
a "am I already running?" latch. (A `useState<boolean>` works too and also lets you disable
the button in the UI; using both is common.)

## The lifecycle as a union type

Recall `FileStatus = "pending" | "uploading" | "processing" | "completed" | "failed"`. Since
the API is synchronous (one request does upload *and* processing), you can treat
`uploading` and `processing` as one visible state — the requirements allow collapsing them.
The union type means when you render status-specific UI, TypeScript can warn you if you
forget a case.

## Try this

Implement `uploadImage` and test it by hand: run the Go API, select one image, and log the
returned manifest. Then implement the sequential `start()` loop with per-item `try/catch`.
Verify that if one file fails (e.g. select a `.gif`, which the API rejects), the others still
complete. Add the double-start guard and a `clear()`.
