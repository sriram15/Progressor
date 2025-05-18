import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      "@": "/src",
      "@bindings": "/bindings",
      "@bindings_service":
        "/bindings/github.com/sriram15/progressor-todo-app/internal/service",
      "@bindings_database":
        "/bindings/github.com/sriram15/progressor-todo-app/internal/database",
    },
  },
});
