/**
 * @type {import('vite').UserConfig}
 */
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: "cmd/web/dist",
    sourcemap: true,
    manifest: true,
    ssrManifest: true,
  },
});
