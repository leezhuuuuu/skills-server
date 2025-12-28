import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    // 输出到 Go 服务的嵌入目录
    outDir: '../cmd/server/web_dist',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // 代理 .md 文件，模拟 LLM 访问
      '/skill.md': 'http://localhost:8080',
      '^/skill/.*\\.md': 'http://localhost:8080'
    }
  }
})
