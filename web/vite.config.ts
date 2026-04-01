import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    outDir: '../biz/handler/index/templates',
    emptyOutDir: true,
    assetsDir: 'assets',
    // Generate static files for embedding
    rollupOptions: {
      output: {
        manualChunks: undefined,
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: (assetInfo) => {
          const info = assetInfo.name?.split('.') || []
          const ext = info[info.length - 1]
          if (/\.(png|jpe?g|gif|svg|webp|ico)$/i.test(assetInfo.name || '')) {
            return 'assets/images/[name][extname]'
          }
          if (/\.(css)$/i.test(assetInfo.name || '')) {
            return 'assets/css/[name][extname]'
          }
          return 'assets/[name][extname]'
        }
      }
    }
  },
  // For development proxy
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:10000',
        changeOrigin: true
      },
      '/download': {
        target: 'http://localhost:10000',
        changeOrigin: true
      }
    }
  }
})
