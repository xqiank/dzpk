import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'node:path'

export default defineConfig({
  base: '/admin/',
  plugins: [vue()],
  resolve: { alias: { '@': path.resolve(__dirname, 'src') } },
  server: {
    host: '127.0.0.1', port: 8722,
    proxy: { '/api': { target: 'http://127.0.0.1:8720', changeOrigin: true, ws: true } }
  },
  build: { outDir: 'dist', sourcemap: false }
})
