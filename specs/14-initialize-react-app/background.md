# Background — Story 14: Initialize the React Application

You plan to AI-generate most of the frontend, which is fine — but understanding these
concepts lets you *prompt precisely* and *review* what you get instead of pasting it
blindly. This doc is lighter than the Go ones on purpose; it gives you the model.

## Vite, in one breath

Vite is the build tool + dev server. `npm run dev` starts a hot-reloading server;
`npm run build` type-checks (`tsc`) and produces static files. The scaffold's
`vite.config.ts` also **proxies `/v1` to `http://localhost:8080`**, so in dev the browser
can call your Go API at a relative path without CORS configuration. Run the Go API and the
Vite server side by side.

## Components and JSX

A React component is a function returning JSX (HTML-like markup):

```tsx
export function App() {
  return <h1>Media Pipeline</h1>;
}
```

JSX quirks to know: attributes are camelCase (`className`, not `class`), `{}` embeds
JavaScript expressions, and you render lists by mapping to elements with a stable `key`:

```tsx
<ul>{files.map(f => <li key={f.id}>{f.file.name}</li>)}</ul>
```

The `key` must be stable and unique per item (that's why each `QueuedFile` gets a
`crypto.randomUUID()` id, not its array index) — React uses it to track items across
re-renders.

## State with `useState`

State is data that, when it changes, re-renders the component:

```tsx
const [files, setFiles] = useState<QueuedFile[]>([]);
```

`files` is the current value; `setFiles` updates it and triggers a re-render. Two rules:

1. **Never mutate state directly.** `files.push(x)` won't re-render and corrupts React's
   model. Always create a *new* array/object:
   ```tsx
   setFiles(prev => [...prev, newFile]);          // add
   setFiles(prev => prev.filter(f => f.id !== id)); // remove
   ```
2. **Use the updater form** (`prev => ...`) when the new value depends on the old — it
   avoids bugs when multiple updates batch together.

This immutable-update discipline is the single most important React habit. It's also why
`useUploadQueue` centralizes the state: one place owns `files`, everything else calls its
functions.

## TypeScript interfaces

TypeScript adds types to JavaScript. The `types/manifest.ts` file mirrors the Go manifest:

```ts
export interface Variant {
  name: string;
  width: number;
  size_bytes: number; // snake_case to match the Go JSON tags
}
```

Note the field names are `snake_case` to match what the Go API sends — the wire format and
your TS types must agree, or `response.json()` gives you `undefined`s. A **union type**
models the file lifecycle:

```ts
export type FileStatus = "pending" | "uploading" | "processing" | "completed" | "failed";
```

A value of this type can only be one of those five strings — the compiler rejects typos and
forces you to handle each case. You'll lean on this in story 15.

## The File API

The browser gives you `File` objects from the input:

```tsx
<input type="file" multiple accept="image/jpeg,image/png"
       onChange={e => {
         const picked = Array.from(e.target.files ?? []);
         setFiles(prev => [...prev, ...picked.map(f => ({
           id: crypto.randomUUID(), file: f, status: "pending" as const,
         }))]);
       }} />
```

- `multiple` allows selecting several; `accept` hints the OS picker toward images (it's a
  hint, not enforcement — the Go API still validates).
- `e.target.files` is a `FileList` (array-like, not an array) — `Array.from(...)` converts
  it. `?? []` guards the null case.
- A `File` has `.name`, `.size` (bytes), and `.type` (MIME). You don't read pixels; you just
  hand the whole `File` to `fetch` later.
- `as const` narrows `"pending"` to the literal type so it fits `FileStatus`.

## Where AI helps and where you should read

Let AI generate the JSX and CSS scaffolding. But *you* should understand: the state shape
(`QueuedFile[]`), the immutable update pattern, and the fact that the frontend never
transforms images — it selects, lists, and (next story) uploads. Review generated code
against the acceptance criteria, especially "files can be removed before processing" and "no
image transformation in the browser."

## Try this

Get `npm run dev` running and rendering "Media Pipeline". Add the file input and the
add/remove logic in `useUploadQueue` + `FileList`. Confirm `npm run build` passes (strict TS
will catch unused vars and type mismatches — treat those as useful feedback).
