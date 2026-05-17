// vite.config.ts
import { defineConfig } from "file:///mnt/d/ai-projects/datamanager/DataRegistry/frontend/node_modules/vite/dist/node/index.js";
import vue from "file:///mnt/d/ai-projects/datamanager/DataRegistry/frontend/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import { resolve } from "path";
var __vite_injected_original_dirname = "/mnt/d/ai-projects/datamanager/DataRegistry/frontend";
var vite_config_default = defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": resolve(__vite_injected_original_dirname, "src")
    }
  },
  server: {
    port: 4004,
    host: "0.0.0.0",
    proxy: {
      "/api": {
        target: "http://localhost:18421",
        changeOrigin: true
      },
      "/public": {
        target: "http://localhost:18421",
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: "../backend/cmd/server/web",
    emptyOutDir: true,
    // 代码分割优化
    rollupOptions: {
      output: {
        manualChunks: {
          // vue 核心框架独立 chunk
          "vue-vendor": ["vue", "vue-router", "pinia"],
          // echarts 独立 chunk（懒加载 + 大小优化）
          "echarts": ["echarts"],
          // element-plus 独立 chunk
          "element-plus": ["element-plus", "@element-plus/icons-vue"]
        }
      }
    },
    // 启用 CSS 代码分割
    cssCodeSplit: true,
    // 生成 sourcemap 便于调试
    sourcemap: false,
    // 启用 gzip 压缩（服务器也建议配置）
    chunkSizeWarningLimit: 1e3
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvbW50L2QvYWktcHJvamVjdHMvZGF0YW1hbmFnZXIvRGF0YVJlZ2lzdHJ5L2Zyb250ZW5kXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCIvbW50L2QvYWktcHJvamVjdHMvZGF0YW1hbmFnZXIvRGF0YVJlZ2lzdHJ5L2Zyb250ZW5kL3ZpdGUuY29uZmlnLnRzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9tbnQvZC9haS1wcm9qZWN0cy9kYXRhbWFuYWdlci9EYXRhUmVnaXN0cnkvZnJvbnRlbmQvdml0ZS5jb25maWcudHNcIjtpbXBvcnQgeyBkZWZpbmVDb25maWcgfSBmcm9tICd2aXRlJ1xuaW1wb3J0IHZ1ZSBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUnXG5pbXBvcnQgeyByZXNvbHZlIH0gZnJvbSAncGF0aCdcblxuZXhwb3J0IGRlZmF1bHQgZGVmaW5lQ29uZmlnKHtcbiAgcGx1Z2luczogW3Z1ZSgpXSxcbiAgcmVzb2x2ZToge1xuICAgIGFsaWFzOiB7XG4gICAgICAnQCc6IHJlc29sdmUoX19kaXJuYW1lLCAnc3JjJylcbiAgICB9XG4gIH0sXG4gIHNlcnZlcjoge1xuICAgIHBvcnQ6IDQwMDQsXG4gICAgaG9zdDogJzAuMC4wLjAnLFxuICAgIHByb3h5OiB7XG4gICAgICAnL2FwaSc6IHtcbiAgICAgICAgdGFyZ2V0OiAnaHR0cDovL2xvY2FsaG9zdDoxODQyMScsXG4gICAgICAgIGNoYW5nZU9yaWdpbjogdHJ1ZVxuICAgICAgfSxcbiAgICAgICcvcHVibGljJzoge1xuICAgICAgICB0YXJnZXQ6ICdodHRwOi8vbG9jYWxob3N0OjE4NDIxJyxcbiAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlXG4gICAgICB9XG4gICAgfVxuICB9LFxuICBidWlsZDoge1xuICAgIG91dERpcjogJy4uL2JhY2tlbmQvY21kL3NlcnZlci93ZWInLFxuICAgIGVtcHR5T3V0RGlyOiB0cnVlLFxuICAgIC8vIFx1NEVFM1x1NzgwMVx1NTIwNlx1NTI3Mlx1NEYxOFx1NTMxNlxuICAgIHJvbGx1cE9wdGlvbnM6IHtcbiAgICAgIG91dHB1dDoge1xuICAgICAgICBtYW51YWxDaHVua3M6IHtcbiAgICAgICAgICAvLyB2dWUgXHU2ODM4XHU1RkMzXHU2ODQ2XHU2N0I2XHU3MkVDXHU3QUNCIGNodW5rXG4gICAgICAgICAgJ3Z1ZS12ZW5kb3InOiBbJ3Z1ZScsICd2dWUtcm91dGVyJywgJ3BpbmlhJ10sXG4gICAgICAgICAgLy8gZWNoYXJ0cyBcdTcyRUNcdTdBQ0IgY2h1bmtcdUZGMDhcdTYxRDJcdTUyQTBcdThGN0QgKyBcdTU5MjdcdTVDMEZcdTRGMThcdTUzMTZcdUZGMDlcbiAgICAgICAgICAnZWNoYXJ0cyc6IFsnZWNoYXJ0cyddLFxuICAgICAgICAgIC8vIGVsZW1lbnQtcGx1cyBcdTcyRUNcdTdBQ0IgY2h1bmtcbiAgICAgICAgICAnZWxlbWVudC1wbHVzJzogWydlbGVtZW50LXBsdXMnLCAnQGVsZW1lbnQtcGx1cy9pY29ucy12dWUnXSxcbiAgICAgICAgfVxuICAgICAgfVxuICAgIH0sXG4gICAgLy8gXHU1NDJGXHU3NTI4IENTUyBcdTRFRTNcdTc4MDFcdTUyMDZcdTUyNzJcbiAgICBjc3NDb2RlU3BsaXQ6IHRydWUsXG4gICAgLy8gXHU3NTFGXHU2MjEwIHNvdXJjZW1hcCBcdTRGQkZcdTRFOEVcdThDMDNcdThCRDVcbiAgICBzb3VyY2VtYXA6IGZhbHNlLFxuICAgIC8vIFx1NTQyRlx1NzUyOCBnemlwIFx1NTM4Qlx1N0YyOVx1RkYwOFx1NjcwRFx1NTJBMVx1NTY2OFx1NEU1Rlx1NUVGQVx1OEJBRVx1OTE0RFx1N0Y2RVx1RkYwOVxuICAgIGNodW5rU2l6ZVdhcm5pbmdMaW1pdDogMTAwMFxuICB9XG59KVxuIl0sCiAgIm1hcHBpbmdzIjogIjtBQUE4VSxTQUFTLG9CQUFvQjtBQUMzVyxPQUFPLFNBQVM7QUFDaEIsU0FBUyxlQUFlO0FBRnhCLElBQU0sbUNBQW1DO0FBSXpDLElBQU8sc0JBQVEsYUFBYTtBQUFBLEVBQzFCLFNBQVMsQ0FBQyxJQUFJLENBQUM7QUFBQSxFQUNmLFNBQVM7QUFBQSxJQUNQLE9BQU87QUFBQSxNQUNMLEtBQUssUUFBUSxrQ0FBVyxLQUFLO0FBQUEsSUFDL0I7QUFBQSxFQUNGO0FBQUEsRUFDQSxRQUFRO0FBQUEsSUFDTixNQUFNO0FBQUEsSUFDTixNQUFNO0FBQUEsSUFDTixPQUFPO0FBQUEsTUFDTCxRQUFRO0FBQUEsUUFDTixRQUFRO0FBQUEsUUFDUixjQUFjO0FBQUEsTUFDaEI7QUFBQSxNQUNBLFdBQVc7QUFBQSxRQUNULFFBQVE7QUFBQSxRQUNSLGNBQWM7QUFBQSxNQUNoQjtBQUFBLElBQ0Y7QUFBQSxFQUNGO0FBQUEsRUFDQSxPQUFPO0FBQUEsSUFDTCxRQUFRO0FBQUEsSUFDUixhQUFhO0FBQUE7QUFBQSxJQUViLGVBQWU7QUFBQSxNQUNiLFFBQVE7QUFBQSxRQUNOLGNBQWM7QUFBQTtBQUFBLFVBRVosY0FBYyxDQUFDLE9BQU8sY0FBYyxPQUFPO0FBQUE7QUFBQSxVQUUzQyxXQUFXLENBQUMsU0FBUztBQUFBO0FBQUEsVUFFckIsZ0JBQWdCLENBQUMsZ0JBQWdCLHlCQUF5QjtBQUFBLFFBQzVEO0FBQUEsTUFDRjtBQUFBLElBQ0Y7QUFBQTtBQUFBLElBRUEsY0FBYztBQUFBO0FBQUEsSUFFZCxXQUFXO0FBQUE7QUFBQSxJQUVYLHVCQUF1QjtBQUFBLEVBQ3pCO0FBQ0YsQ0FBQzsiLAogICJuYW1lcyI6IFtdCn0K
