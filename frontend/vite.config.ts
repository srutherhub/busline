import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  resolve: {
    dedupe: ["react", "react-dom"],
  },
  optimizeDeps: {
    exclude: ["maplibre-gl"],
    include: ["react", "react-dom", "@deck.gl/react"],
  },
  server: {
    proxy: { "/api": "http://localhost" },
  },
});
