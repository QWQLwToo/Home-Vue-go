import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true,
      },
    },
  },
  server: {
    port: 1552,
    proxy: {
      '/api': {
        target: 'http://localhost:1551',
        changeOrigin: true,
      },
      '/uploads': {
        target: 'http://localhost:1551',
        changeOrigin: true,
      },
    },
  },
});