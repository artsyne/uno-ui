import { defineConfig } from 'vite';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  root: '.', 
  build: {
    outDir: 'dist',
    lib: {
      entry: resolve(__dirname, 'src/index.js'),
      name: 'Uno',
      fileName: (format) => `uno.${format}.js`,
    },
    rollupOptions: {
      external: ['alpinejs', 'htmx.org'],
      output: {
        globals: {
          alpinejs: 'Alpine',
          'htmx.org': 'htmx',
        },
      },
    },
  },
  server: {
    open: '/examples/index.html',
    port: 5176,
    strictPort: true,
  },
});
