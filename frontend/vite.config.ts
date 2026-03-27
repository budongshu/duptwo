import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 4004,
    host: '0.0.0.0',
    proxy: {
      '/api': {
        target: 'http://localhost:18421',
        changeOrigin: true
      },
      '/public': {
        target: 'http://localhost:18421',
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: '../backend/cmd/server/web',
    emptyOutDir: true,
    // 代码分割优化
    rollupOptions: {
      output: {
        manualChunks: {
          // vue 核心框架独立 chunk
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          // echarts 独立 chunk（懒加载 + 大小优化）
          'echarts': ['echarts'],
          // element-plus 独立 chunk
          'element-plus': ['element-plus', '@element-plus/icons-vue'],
        }
      }
    },
    // 启用 CSS 代码分割
    cssCodeSplit: true,
    // 生成 sourcemap 便于调试
    sourcemap: false,
    // 启用 gzip 压缩（服务器也建议配置）
    chunkSizeWarningLimit: 1000
  }
})
