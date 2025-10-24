import { fileURLToPath, URL } from 'node:url'
import path from 'path'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import { getProxyConfig } from './proxy_config'
import { visualizer } from 'rollup-plugin-visualizer'

// https://vite.dev/config/
export default defineConfig((opt) => {
  const { command, mode } = opt

  // eslint-disable-next-line no-undef
  const env = loadEnv(mode, process.cwd(), '')
  const base = command === 'serve' ? '/' : '/open-doc'

  return {
    plugins: [vue(), vueJsx()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    base: base,
    server: {
      // host: '0.0.0.0',
      port: 5521,
      proxy: getProxyConfig(opt),
      open: true,
    },
    esbuild: {
      pure: env.VITE_DROP_CONSOLE === 'true' ? ['console.log'] : undefined,
      drop: env.VITE_DROP_DEBUGGER === 'true' ? ['debugger'] : undefined,
    },
    css: {
      postcss: {
        plugins: [],
      },
      extract:
        env.VITE_USE_CSS_SPLIT === 'true'
          ? {
              filename: 'assets/css/[name]-[hash].css',
            }
          : false,
    },
    optimizeDeps: {
      include: [
        'vue',
        'vue-router',
        'pinia',
        'axios',
        'dayjs',
        'crypto-js',
        'qs',
        'cherry-markdown',
      ],
      exclude: ['vue-devtools'], // 排除开发工具
    },
    build: {
      outDir: fileURLToPath(new URL('../../static/open-doc', import.meta.url)),
      emptyOutDir: true,
      assetsDir: 'assets',
      sourcemap: env.VITE_SOURCEMAP === 'true',
      minify: 'esbuild',
      target: 'esnext',
      cssCodeSplit: !(env.VITE_USE_CSS_SPLIT === 'false'),
      rollupOptions: {
        plugins: [env.VITE_USE_BUNDLE_ANALYZER === 'true' ? visualizer() : undefined],
        output: {
          advancedChunks: {
            groups: [
              {
                test: /node_modules\/(vue|vue-router|vue-i18n)/,
                name: 'vue-chunks',
              },
              {
                test: /node_modules\/pinia/,
                name: 'pinia',
              },
              {
                test: /node_modules\/dayjs/,
                name: 'dayjs',
              },
              {
                test: /node_modules\/axios/,
                name: 'axios',
              },
              {
                test: /node_modules\/crypto-js/,
                name: 'crypto-js',
              },
              {
                test: /node_modules\/qs/,
                name: 'qs',
              },
              {
                test: /node_modules\/cherry-markdown/,
                name: 'cherry-markdown',
              },
            ],
          },
          chunkFileNames: 'assets/js/[name]-[hash].js',
          entryFileNames: 'assets/js/[name]-[hash].js',
          assetFileNames: 'assets/[ext]/[name]-[hash].[ext]',
        },
      },
    },
  }
})
