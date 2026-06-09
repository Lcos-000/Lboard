import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  // 配置Vite服务器端口为5173，同时配置反向代理
  // 代理/api、/healthz、/readyz这三个路径到本地8080端口
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8080',
      '/healthz': 'http://localhost:8080',
      '/readyz': 'http://localhost:8080'
    }
  }
});
