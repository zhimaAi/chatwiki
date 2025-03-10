/* eslint-disable no-undef */
import { fileURLToPath, URL } from 'node:url'
import path from 'path'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
// import VueDevTools from 'vite-plugin-vue-devtools'
import { visualizer } from 'rollup-plugin-visualizer'
import externalGlobals from 'rollup-plugin-external-globals'
import Components from 'unplugin-vue-components/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import { getProxyConfig } from './proxy_config'

const { resolve } = path
const root = process.cwd()
function pathResolve(dir) {
  return resolve(root, '.', dir)
}

// const globals = externalGlobals({
//   moment: 'moment',
//   'video.js': 'videojs',
//   jspdf: 'jspdf',
//   xlsx: 'XLSX',
// });

const globals = externalGlobals({})

// https://vitejs.dev/config/
export default defineConfig((opt) => {
  const { command, mode } = opt
  // eslint-disable-next-line no-unused-vars
  const env = loadEnv(mode, process.cwd(), '')
  const base = command === 'serve' ? '/' : '/'

  return {
    plugins: [
      vue(),
      vueJsx(),
      // VueDevTools(),
      Components({
        resolvers: [
          AntDesignVueResolver({
            importStyle: false // css in js
          })
        ]
      }),
      createSvgIconsPlugin({
        // 指定需要缓存的图标文件夹
        iconDirs: [resolve('./src/assets/svg')],
        // 指定symbolId格式
        symbolId: '[name]'
      }),
      VueI18nPlugin({
        runtimeOnly: true,
        compositionOnly: true,
        include: [path.resolve(__dirname, './src/locales/lang**')]
      })
      // copyIndex(),
    ],
    resolve: {
      alias: [
        {
          find: 'vue-i18n',
          replacement: 'vue-i18n/dist/vue-i18n.cjs.js'
        },
        {
          find: /@\//,
          replacement: `${pathResolve('src')}/`
        }
      ]
    },
    base: base,
    experimental: {
      // 进阶基础路径选项
    },
    server: {
      host: '0.0.0.0',
      proxy: getProxyConfig(opt),
      port: 5520,
      open: true,
    },
    esbuild: {
      pure: env.VITE_DROP_CONSOLE === 'true' ? ['console.log'] : undefined,
      drop: env.VITE_DROP_DEBUGGER === 'true' ? ['debugger'] : undefined
    },
    build: {
      outDir: fileURLToPath(new URL('../../static/chat-ai-admin', import.meta.url)),
      emptyOutDir: true,
      assetsDir: 'assets',
      sourcemap: env.VITE_SOURCEMAP === 'true',
      rollupOptions: {
        // external: ['moment', 'video.js', 'jspdf', 'xlsx'],
        external: [],
        plugins: [globals, env.VITE_USE_BUNDLE_ANALYZER === 'true' ? visualizer() : undefined],
        output: {
          // 自定义chunkFileName生成规则
          chunkFileNames: 'assets/js/[name]-[hash].js',
          entryFileNames: '[name]-[hash].js',
          assetFileNames(assetInfo) {
            let fiel_name = assetInfo.name.toLowerCase()

            if (fiel_name.endsWith('.css')) {
              return 'assets/css/[name]-[hash].[ext]'
            }
            if (['png', 'jpg', 'jpeg', 'svg'].some((ext) => fiel_name.endsWith(ext))) {
              return 'assets/img/[name]-[hash].[ext]'
            }
            if (['ttf', 'woff', 'woff2'].some((ext) => fiel_name.endsWith(ext))) {
              return 'assets/fonts/[name]-[hash].[ext]'
            }
            return 'assets/[name]-[hash].[ext]'
          },
          // 该选项允许你创建自定义的公共 chunk
          manualChunks: {
            'vue-chunks': ['vue', 'vue-router', 'pinia', 'vue-i18n'],
            dayjs: ['dayjs'],
            axios: ['axios'],
            'crypto-js': ['crypto-js'],
            qs: ['qs'],
            'vue-pdf-embed': ['vue-pdf-embed']
          }
        }
      },
      cssCodeSplit: !(env.VITE_USE_CSS_SPLIT === 'false')
    },
    optimizeDeps: {
      include: ['vue', 'vue-router', 'pinia', 'vue-i18n', 'dayjs', 'axios', 'crypto-js', 'qs']
    }
  }
})
