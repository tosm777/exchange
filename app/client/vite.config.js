import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: [{ find: '@', replacement: '/src' }]
  },
  server: {
    // https: true,
    port: 443,
  }
    // hmr: {
    //   overlay: false,
    // }
})
