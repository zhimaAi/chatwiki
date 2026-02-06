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
  console.log(typeof env.VITE_USE_BUNDLE_ANALYZER)
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
        // {
        //   find: 'vue-i18n',
        //   replacement: 'vue-i18n/dist/vue-i18n.cjs.js'
        // },
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
      ssr: false,
      outDir: fileURLToPath(new URL('../../static/chat-ai-admin', import.meta.url)),
      emptyOutDir: true,
      assetsDir: 'assets',
      sourcemap: env.VITE_SOURCEMAP === 'true',
      reportCompressedSize: false,
      rollupOptions: {
        // external: ['moment', 'video.js', 'jspdf', 'xlsx'],
        external: [],
        plugins: [globals, env.VITE_USE_BUNDLE_ANALYZER === 'true' ? visualizer({open: true}) : undefined],
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
          manualChunks(id) {
            // 语言文件按语言自动分离
            if (id.includes('/src/locales/lang/')) {
              const langMatch = id.match(/\/src\/locales\/lang\/([a-zA-Z0-9-]+)\//)
              if (langMatch) {
                return `lang-${langMatch[1]}`
              }
            }

            // 其他第三方库分离
            if (id.includes('node_modules')) {
              // 截取 node_modules 后的路径，避免项目名干扰
              const moduleId = id.split('node_modules/')[1]
              
              // Cherry Markdown - 必须在所有包含"vue"的库之前
              if (moduleId.includes('cherry') || moduleId.includes('Cherry')) {
                return 'editor-cherry'
              }

              // 重要：先匹配包含"vue"的具体库，再匹配核心vue（避免误判）
              // UI组件库 - 包含 vue 关键字，必须在 vue 之前
              if (moduleId.includes('ant-design-vue') || moduleId.includes('@ant-design')) {
                return 'ui-antd'
              }
              // 编辑器库 - 包含 vue 关键字，必须在 vue 之前
              if (moduleId.includes('@wangeditor')) {
                return 'wange-ditor'
              }
              if (moduleId.includes('vue-pdf-embed')) {
                return 'vue-pdf-embed'
              }
              if (moduleId.includes('codemirror') || moduleId.includes('markdown-it')) {
                return 'editor-code'
              }
              // Medium Editor
              if (moduleId.includes('medium-editor')) {
                return 'medium-editor'
              }
              // Emoji Picker
              if (moduleId.includes('emoji-mart')) {
                return 'emoji-mart'
              }
              // Lodash
              if (moduleId.includes('lodash')) {
                return 'lodash'
              }
              // LogicFlow
              if (moduleId.includes('@logicflow')) {
                return 'logicflow'
              }
              // 图表库
              if (moduleId.includes('echarts')) {
                return 'charts'
              }
              // 工具库
              if (moduleId.includes('axios') || moduleId.includes('dayjs') || moduleId.includes('crypto-js') ||
                  moduleId.includes('qs') || moduleId.includes('js-cookie') || moduleId.includes('file-saver')) {
                return 'utils'
              }
              // 布局相关
              if (moduleId.includes('elkjs')) {
                return 'elkjs'
              }
              // Neo4j相关
              if (moduleId.includes('@neo4j')) {
                return 'neo4j'
              }
              // 其他
              if (moduleId.includes('html2canvas') || moduleId.includes('v-viewer')) {
                return 'other'
              }
              // 框架核心 - 放在最后
              if (moduleId.includes('vue') || moduleId.includes('pinia') || moduleId.includes('vue-router') || moduleId.includes('vue-i18n')) {
                return 'vue-chunks'
              }
            }
          }
        }
      },
      cssCodeSplit: !(env.VITE_USE_CSS_SPLIT === 'false')
    },
    optimizeDeps: {
      exclude: ['@neo4j-nvl/layout-workers'],
      include: [
        'vue', 
        'vue-router',
        'ant-design-vue',
        'pinia', 
        'vue-i18n', 
        'dayjs', 
        'axios', 
        'crypto-js', 
        'qs',
        'echarts',
        '@neo4j-nvl/layout-workers > cytoscape',
        '@neo4j-nvl/layout-workers > cytoscape-cose-bilkent',
        '@neo4j-nvl/layout-workers > @neo4j-bloom/dagre',
        '@neo4j-nvl/layout-workers > bin-pack',
        '@neo4j-nvl/layout-workers > graphlib',
      ]
    }
  }
})
