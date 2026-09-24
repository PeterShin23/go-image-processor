# Background — Story 17: Display Results and Retry Failures

The last frontend story is about rendering derived state and doing a targeted re-run. Light
on new concepts, heavy on "get the details right." This is where the app becomes genuinely
usable.

## Rendering conditionally

React renders different UI based on state. Two common tools:

```tsx
{f.status === "failed" && <button onClick={() => onRetry(f.id)}>Retry</button>}
{f.manifest && <VariantList variants={f.manifest.variants} />}
```

- `cond && <JSX/>` renders the element only when `cond` is truthy (a falsy value renders
  nothing). Watch the classic trap: `array.length && <JSX/>` renders a literal `0` when the
  array is empty — use `array.length > 0 && ...` instead.
- For multiple branches, a small helper or a `switch` returning JSX reads cleaner than nested
  ternaries.

## Three success states, not two

An upload can end in three meaningfully different places — model all three:

1. **completed** — the request succeeded and `manifest.status === "completed"` (all variants
   built).
2. **partial** — the request succeeded but `manifest.status === "partial"` (some variants
   failed; `manifest.errors` is non-empty). This is a *server* outcome, delivered in a 200
   response.
3. **failed** — the request itself failed (network error, 4xx/5xx). There's no manifest; you
   have an error string.

Note partial vs failed come from different layers: `partial` is inside a successful
response's manifest; `failed` is the fetch/HTTP failing. Your `QueuedFile` already
distinguishes them (`status` + optional `manifest` + optional `error`). Derive the badge:

```tsx
function badge(f: QueuedFile): { label: string; color: string } {
  if (f.status === "failed") return { label: "Failed", color: "red" };
  if (f.manifest?.status === "partial") return { label: "Partial", color: "goldenrod" };
  if (f.status === "completed") return { label: "Completed", color: "green" };
  return { label: f.status, color: "gray" };
}
```

The acceptance criterion "partial success is visually distinguishable" means partial must
*not* look like full success — a different color/label is enough.

## Listing variants

The manifest's `variants` array maps straight to rows:

```tsx
<ul>
  {variants.map(v => (
    <li key={v.name}>
      {v.name}: {v.width}×{v.height} {v.format} ({Math.round(v.size_bytes / 1024)} KB)
    </li>
  ))}
</ul>
```

`v.name` is unique per manifest, so it's a fine `key`. Show the fields the requirements list
(name, width, height, format, size). For partial results, also render `manifest.errors` (the
`profile` + `message`) so the user sees which variant failed and why.

## Retry: re-run one file

Retry reuses your existing single-file processing, scoped to one id:

```ts
async function retry(id: string) {
  const qf = files.find(f => f.id === id);
  if (!qf) return;
  setFiles(prev => prev.map(f => f.id === id
    ? { ...f, status: "uploading", error: undefined, manifest: undefined } : f));
  await processOne(qf); // the same function the queue uses
}
```

The important guarantee: retry touches **only** that file. Because you update by id and call
`processOne` for a single `QueuedFile`, successful files are never re-uploaded — satisfying
"retrying does not resend successful files." Clear the previous `error`/`manifest` so the row
doesn't show stale results while retrying.

## Keeping it usable with many files

"Usable with at least 10 files" mostly means: stable list order (story 16), no layout that
collapses, and results that don't shift around as items complete. You get this for free if you
render from the unordered `files` array and update by id. Don't add virtualization or
pagination — that's over-engineering for this scope.

## Stretch (only if time remains)

Image previews (`URL.createObjectURL(file)`), drag-and-drop, and per-file progress percentages
are listed as optional stretch goals. Skip them until every acceptance criterion passes and
the README (story 18) is written.

## Try this

Render the badge + variant list for completed files first, then add the partial styling using
a file whose manifest has errors (you can force a server partial by including a profile that
fails, or just trust the shape). Finally wire `retry(id)` and confirm clicking it re-runs only
the failed row. Run `npm run build` to catch type issues.
