import { fileURLToPath, URL } from 'node:url'
import path from 'path'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import { visualizer } from 'rollup-plugin-visualizer'
import externalGlobals from 'rollup-plugin-external-globals'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import { getProxyConfig } from './proxy_config'
import type { UserConfig, ConfigEnv, Plugin } from 'vite'
import type {} from 'vite/types/metadata.d.ts'

function previewWarmup(): Plugin {
  const chatPage = fileURLToPath(new URL('./src/views/chat/index.vue', import.meta.url))
  return {
    name: 'preview-warmup',
    apply: 'build',
    generateBundle: {
      order: 'post',
      handler(_, bundle) {
        const chat = Object.values(bundle).find(
          (item) => item.type === 'chunk' && item.facadeModuleId === chatPage
        )
        if (chat?.type !== 'chunk') this.error('PC preview chat chunk not found')

        // Warm the lazy chat chunk from the preview origin without starting a session.
        this.emitFile({
          type: 'asset',
          fileName: 'preview-warmup.html',
          source: [
            '<!doctype html><meta charset="utf-8"><title>ChatWiki</title>',
            `<link rel="preload" as="script" href="./${chat.fileName}" crossorigin fetchpriority="low">`,
            ...Array.from(chat.viteMetadata?.importedCss ?? [], (css) =>
              `<link rel="preload" as="style" href="./${css}" fetchpriority="low">`
            )
          ].join('\n')
        })
      }
    }
  }
}

export default defineConfig((opt: ConfigEnv): UserConfig => {
  const { command, mode } = opt
  const env = loadEnv(mode, process.cwd(), '')
  const base = command === 'serve' ? '/' : './'

  const globals = externalGlobals({})

  return {
    plugins: [
      vue(),
      vueJsx(),
      previewWarmup(),
      createSvgIconsPlugin({
        // 指定需要缓存的图标文件夹
        iconDirs: [path.resolve(process.cwd(), 'src/assets/icons')],
        // 指定symbolId格式
        symbolId: '[name]'
      }),
      VueI18nPlugin({
        runtimeOnly: true,
        compositionOnly: true,
        include: [path.resolve(__dirname, './src/locales/lang**')]
      })
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    base: base,
    server: {
      proxy: getProxyConfig(opt)
    },
    build: {
      outDir: fileURLToPath(new URL('../../static/chat-ai-pc/web', import.meta.url)),
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
          entryFileNames: 'assets/js/[name]-[hash].js',
          assetFileNames(assetInfo: any) {
            const fiel_name = assetInfo.name.toLowerCase()

            if (fiel_name.endsWith('.css')) {
              return 'assets/css/[name]-[hash].[ext]'
            }
            if (['png', 'jpg', 'jpeg', 'svg'].some((ext) => fiel_name.endsWith(ext))) {
              return 'assets/img/[name]-[hash].[ext]'
            }
            if (['ttf', 'woff', 'woff2'].some((ext) => fiel_name.endsWith(ext))) {
              return 'assets/fonts/[name]-[hash].[ext]'
            }
            return 'assets/[ext]/[name]-[hash].[ext]'
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
