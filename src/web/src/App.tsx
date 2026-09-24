import { useUploadQueue } from "./hooks/useUploadQueue";
import { FileList } from "./components/FileList";

// App ties together file selection, the queue, and the results display.
//
// STORY 14: a multi-file <input accept="image/jpeg,image/png"> that adds files.
// STORY 15: a "Process" button that starts the queue.
// STORY 16: concurrency of 2 under the hood.
// STORY 17: per-file results + retry.
export function App() {
  const { files } = useUploadQueue();

  return (
    <main style={{ maxWidth: 720, margin: "2rem auto", fontFamily: "system-ui" }}>
      <h1>Media Pipeline</h1>
      {/* TODO (story 14): file input + selected count */}
      {/* TODO (story 15): Process / Clear buttons */}
      <FileList files={files} onRemove={() => {}} onRetry={() => {}} />
    </main>
  );
}
