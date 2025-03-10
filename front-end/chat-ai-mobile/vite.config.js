import { fileURLToPath, URL } from 'node:url';
import path from 'path';
import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite';
import { visualizer } from 'rollup-plugin-visualizer';
import externalGlobals from 'rollup-plugin-external-globals';
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons';
import { getProxyConfig } from './proxy_config';
export default defineConfig(function (opt) {
    var command = opt.command, mode = opt.mode;
    var env = loadEnv(mode, process.cwd(), '');
    var base = command === 'serve' ? '/' : '/';
    var globals = externalGlobals({});
    return {
        plugins: [
            vue(),
            vueJsx(),
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
            outDir: fileURLToPath(new URL('../../static/chat-ai-mobile', import.meta.url)),
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
                    assetFileNames: function (assetInfo) {
                        var fiel_name = assetInfo.name.toLowerCase();
                        if (fiel_name.endsWith('.css')) {
                            return 'assets/css/[name]-[hash].[ext]';
                        }
                        if (['png', 'jpg', 'jpeg', 'svg'].some(function (ext) { return fiel_name.endsWith(ext); })) {
                            return 'assets/img/[name]-[hash].[ext]';
                        }
                        if (['ttf', 'woff', 'woff2'].some(function (ext) { return fiel_name.endsWith(ext); })) {
                            return 'assets/fonts/[name]-[hash].[ext]';
                        }
                        return 'assets/[name]-[hash].[ext]';
                    },
                    // 该选项允许你创建自定义的公共 chunk
                    manualChunks: {
                        'vue-chunks': ['vue', 'vue-router', 'pinia', 'vue-i18n'],
                        dayjs: ['dayjs'],
                        axios: ['axios'],
                        'crypto-js': ['crypto-js'],
                        qs: ['qs'],
                    }
                }
            },
            cssCodeSplit: !(env.VITE_USE_CSS_SPLIT === 'false')
        },
        optimizeDeps: {
            include: ['vue', 'vue-router', 'pinia', 'vue-i18n', 'dayjs', 'axios', 'crypto-js', 'qs']
        }
    };
});
