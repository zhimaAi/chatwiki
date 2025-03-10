import { resolve } from 'path'
import { defineConfig, externalizeDepsPlugin, defineViteConfig, loadEnv } from 'electron-vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'

export default defineConfig({
  main: {
    plugins: [externalizeDepsPlugin()]
  },
  preload: {
    plugins: [externalizeDepsPlugin()]
  },
  renderer: defineViteConfig((mode) => {
    const env = loadEnv(mode)
    console.log(env)
    return {
      plugins: [
        vue(),
        Components({
          resolvers: [AntDesignVueResolver({ importStyle: false })]
        }),
        createSvgIconsPlugin({
          // 指定需要缓存的图标文件夹
          iconDirs: [resolve('src/renderer/src/assets/svg')],
          // 指定symbolId格式
          symbolId: '[name]'
        }),
        VueI18nPlugin({
          runtimeOnly: true,
          compositionOnly: true,
          include: [resolve(__dirname, 'src/renderer/src/locales/lang**')]
        })
      ],
      resolve: {
        alias: {
          '@renderer': resolve('src/renderer/src'),
          '@': resolve('src/renderer/src'),
          'vue-i18n': 'vue-i18n/dist/vue-i18n.cjs.js'
        }
      }
    }
  })
})
