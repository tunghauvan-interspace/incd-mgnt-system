import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// Determine API target for dev proxy. When running inside Docker, use 'backend'.
const apiTarget = process.env.VITE_API_TARGET || process.env.VITE_API_HOST || 'http://backend:8080'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  build: {
    outDir: '../static',
    emptyOutDir: true,
    // Enable minification
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true
      }
    },
    // Enable gzip compression reporting
    reportCompressedSize: true,
    // Chunk size warning limit
    chunkSizeWarningLimit: 500,
    rollupOptions: {
      output: {
        // Organize assets by type with content hashing
        entryFileNames: 'js/[name]-[hash].js',
        chunkFileNames: 'js/[name]-[hash].js',
        assetFileNames: (assetInfo: any) => {
          const info = assetInfo.name?.split('.') || []
          const ext = info[info.length - 1]

          if (/\.(png|jpe?g|svg|gif|tiff|bmp|ico)$/i.test(assetInfo.name || '')) {
            return 'images/[name]-[hash][extname]'
          }

          if (/\.css$/i.test(assetInfo.name || '')) {
            return 'css/[name]-[hash][extname]'
          }

          return `${ext}/[name]-[hash][extname]`
        },
        // Manual chunk splitting for better caching
        manualChunks: {
          // Vendor libraries
          vendor: ['vue', 'vue-router', 'pinia'],
          charts: ['chart.js', 'vue-chartjs'],
          utils: ['axios']
        }
      }
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: apiTarget,
        changeOrigin: true
      }
    }
  },
  // Bundle analyzer (can be enabled when needed)
  define: {
    __VUE_PROD_DEVTOOLS__: false,
  }
})
