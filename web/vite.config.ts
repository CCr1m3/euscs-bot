import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/api": "http://localhost:9000",
      "/auth": "http://localhost:9000",
      "/twitch": "http://localhost:9000",
    },
  },
});
