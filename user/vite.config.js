import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

const baseUrl = {
  development: './',
  production: '/users/',
}

export default ({ mode }) => defineConfig({
  plugins: [vue()],
  base: baseUrl[mode],
  resolve: {
    alias: {
      '~': path.resolve(__dirname, './'),
      '@': path.resolve(__dirname, 'src')
    }
  },
  build: {
    outDir: '../server/Tel/dist/user', //打包路径
    chunkSizeWarningLimit: 3500, //打包大小限制
    emptyOutDir: true, // 构建时清空该目录
  },
  css: {
    preprocessorOptions: {
      less: {
        modifyVars: { // 在这里自定义主题色等样式 7f8fef
          'primary-color': '#7556fd'
        },
        javascriptEnabled: true,
      },
    },
  },
  server: {
    host: '0.0.0.0',
    cors: true, // 默认启用并允许任何源
    // proxy: {
    //     '/api': {
    //         target: 'http://192.168.1.13',
    //         changeOrigin: true,
    //         rewrite: path => path.replace(/^\/api/, '')
    //     }
    // }
  }
})
