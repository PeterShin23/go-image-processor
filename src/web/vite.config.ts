import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// The Go API runs on :8080. Proxy /v1 during dev so the browser can call the
// backend without CORS setup.
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/v1": "http://localhost:8080",
    },
  },
});
