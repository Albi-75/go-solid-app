import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  plugins: [tailwindcss(),solid()],
server: {
  proxy: {
    '/api': {
      target: 'http://backend:8080',
      changeOrigin: true,
      rewrite: path => path.replace(/^\/api/, ''),
    }
  }
}
})
