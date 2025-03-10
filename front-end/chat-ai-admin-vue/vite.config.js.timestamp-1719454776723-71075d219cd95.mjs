// vite.config.js
import { fileURLToPath, URL } from "node:url";
import path from "path";
import { defineConfig, loadEnv as loadEnv2 } from "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/vite/dist/node/index.js";
import vue from "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import { visualizer } from "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/rollup-plugin-visualizer/dist/plugin/index.js";
import externalGlobals from "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/rollup-plugin-external-globals/index.js";
import Components from "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/unplugin-vue-components/dist/vite.js";
import { AntDesignVueResolver } from "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/unplugin-vue-components/dist/resolvers.js";
import VueI18nPlugin from "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/@intlify/unplugin-vue-i18n/lib/vite.mjs";
import { createSvgIconsPlugin } from "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/vite-plugin-svg-icons/dist/index.mjs";

// proxy_config.js
import { loadEnv } from "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/vite/dist/node/index.js";
var getProxyConfig = (opt) => {
  const { mode } = opt;
  const env = loadEnv(mode, process.cwd(), "");
  let proxyApis = ["/static", "/common", "/manage", "/app", "/chat", "/upload"];
  let proxy = {};
  console.log(env.PROXY_BASE_API_URL);
  proxyApis.forEach((key) => {
    proxy[key] = {
      target: env.PROXY_BASE_API_URL,
      changeOrigin: true
    };
  });
  return {
    ...proxy
  };
};

// vite.config.js
var __vite_injected_original_dirname = "/Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue";
var __vite_injected_original_import_meta_url = "file:///Users/aixiangfei/code/zhima_chat_ai/front-end/chat-ai-admin-vue/vite.config.js";
var { resolve } = path;
var root = process.cwd();
function pathResolve(dir) {
  return resolve(root, ".", dir);
}
var globals = externalGlobals({});
var vite_config_default = defineConfig((opt) => {
  const { command, mode } = opt;
  const env = loadEnv2(mode, process.cwd(), "");
  const base = command === "serve" ? "/" : "/";
  return {
    plugins: [
      vue(),
      vueJsx(),
      // VueDevTools(),
      Components({
        resolvers: [
          AntDesignVueResolver({
            importStyle: false
            // css in js
          })
        ]
      }),
      createSvgIconsPlugin({
        // 指定需要缓存的图标文件夹
        iconDirs: [resolve("./src/assets/svg")],
        // 指定symbolId格式
        symbolId: "[name]"
      }),
      VueI18nPlugin({
        runtimeOnly: true,
        compositionOnly: true,
        include: [path.resolve(__vite_injected_original_dirname, "./src/locales/lang**")]
      })
      // copyIndex(),
    ],
    resolve: {
      alias: [
        {
          find: "vue-i18n",
          replacement: "vue-i18n/dist/vue-i18n.cjs.js"
        },
        {
          find: /@\//,
          replacement: `${pathResolve("src")}/`
        }
      ]
    },
    base,
    experimental: {
      // 进阶基础路径选项
    },
    server: {
      proxy: getProxyConfig(opt)
    },
    esbuild: {
      pure: env.VITE_DROP_CONSOLE === "true" ? ["console.log"] : void 0,
      drop: env.VITE_DROP_DEBUGGER === "true" ? ["debugger"] : void 0
    },
    build: {
      outDir: fileURLToPath(new URL("../../static/chat-ai-admin", __vite_injected_original_import_meta_url)),
      assetsDir: "assets",
      sourcemap: env.VITE_SOURCEMAP === "true",
      rollupOptions: {
        // external: ['moment', 'video.js', 'jspdf', 'xlsx'],
        external: [],
        plugins: [globals, env.VITE_USE_BUNDLE_ANALYZER === "true" ? visualizer() : void 0],
        output: {
          // 自定义chunkFileName生成规则
          chunkFileNames: "assets/js/[name]-[hash].js",
          entryFileNames: "[name].js",
          assetFileNames(assetInfo) {
            let fiel_name = assetInfo.name.toLowerCase();
            if (fiel_name.endsWith(".css")) {
              return "assets/css/[name]-[hash].[ext]";
            }
            if (["png", "jpg", "jpeg", "svg"].some((ext) => fiel_name.endsWith(ext))) {
              return "assets/img/[name]-[hash].[ext]";
            }
            if (["ttf", "woff", "woff2"].some((ext) => fiel_name.endsWith(ext))) {
              return "assets/fonts/[name]-[hash].[ext]";
            }
            return "assets/[name]-[hash].[ext]";
          },
          // 该选项允许你创建自定义的公共 chunk
          manualChunks: {
            "vue-chunks": ["vue", "vue-router", "pinia", "vue-i18n"],
            dayjs: ["dayjs"],
            axios: ["axios"],
            "crypto-js": ["crypto-js"],
            qs: ["qs"],
            "vue-pdf-embed": ["vue-pdf-embed"]
          }
        }
      },
      cssCodeSplit: !(env.VITE_USE_CSS_SPLIT === "false")
    },
    optimizeDeps: {
      include: ["vue", "vue-router", "pinia", "vue-i18n", "dayjs", "axios", "crypto-js", "qs"]
    }
  };
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiLCAicHJveHlfY29uZmlnLmpzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiL1VzZXJzL2FpeGlhbmdmZWkvY29kZS96aGltYV9jaGF0X2FpL2Zyb250LWVuZC9jaGF0LWFpLWFkbWluLXZ1ZVwiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiL1VzZXJzL2FpeGlhbmdmZWkvY29kZS96aGltYV9jaGF0X2FpL2Zyb250LWVuZC9jaGF0LWFpLWFkbWluLXZ1ZS92aXRlLmNvbmZpZy5qc1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vVXNlcnMvYWl4aWFuZ2ZlaS9jb2RlL3poaW1hX2NoYXRfYWkvZnJvbnQtZW5kL2NoYXQtYWktYWRtaW4tdnVlL3ZpdGUuY29uZmlnLmpzXCI7LyogZXNsaW50LWRpc2FibGUgbm8tdW5kZWYgKi9cbmltcG9ydCB7IGZpbGVVUkxUb1BhdGgsIFVSTCB9IGZyb20gJ25vZGU6dXJsJ1xuaW1wb3J0IHBhdGggZnJvbSAncGF0aCdcblxuaW1wb3J0IHsgZGVmaW5lQ29uZmlnLCBsb2FkRW52IH0gZnJvbSAndml0ZSdcbmltcG9ydCB2dWUgZnJvbSAnQHZpdGVqcy9wbHVnaW4tdnVlJ1xuaW1wb3J0IHZ1ZUpzeCBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUtanN4J1xuLy8gaW1wb3J0IFZ1ZURldlRvb2xzIGZyb20gJ3ZpdGUtcGx1Z2luLXZ1ZS1kZXZ0b29scydcbmltcG9ydCB7IHZpc3VhbGl6ZXIgfSBmcm9tICdyb2xsdXAtcGx1Z2luLXZpc3VhbGl6ZXInXG5pbXBvcnQgZXh0ZXJuYWxHbG9iYWxzIGZyb20gJ3JvbGx1cC1wbHVnaW4tZXh0ZXJuYWwtZ2xvYmFscydcbmltcG9ydCBDb21wb25lbnRzIGZyb20gJ3VucGx1Z2luLXZ1ZS1jb21wb25lbnRzL3ZpdGUnXG5pbXBvcnQgeyBBbnREZXNpZ25WdWVSZXNvbHZlciB9IGZyb20gJ3VucGx1Z2luLXZ1ZS1jb21wb25lbnRzL3Jlc29sdmVycydcbmltcG9ydCBWdWVJMThuUGx1Z2luIGZyb20gJ0BpbnRsaWZ5L3VucGx1Z2luLXZ1ZS1pMThuL3ZpdGUnXG5pbXBvcnQgeyBjcmVhdGVTdmdJY29uc1BsdWdpbiB9IGZyb20gJ3ZpdGUtcGx1Z2luLXN2Zy1pY29ucydcbmltcG9ydCB7IGdldFByb3h5Q29uZmlnIH0gZnJvbSAnLi9wcm94eV9jb25maWcnXG5cbmNvbnN0IHsgcmVzb2x2ZSB9ID0gcGF0aFxuY29uc3Qgcm9vdCA9IHByb2Nlc3MuY3dkKClcbmZ1bmN0aW9uIHBhdGhSZXNvbHZlKGRpcikge1xuICByZXR1cm4gcmVzb2x2ZShyb290LCAnLicsIGRpcilcbn1cblxuLy8gY29uc3QgZ2xvYmFscyA9IGV4dGVybmFsR2xvYmFscyh7XG4vLyAgIG1vbWVudDogJ21vbWVudCcsXG4vLyAgICd2aWRlby5qcyc6ICd2aWRlb2pzJyxcbi8vICAganNwZGY6ICdqc3BkZicsXG4vLyAgIHhsc3g6ICdYTFNYJyxcbi8vIH0pO1xuXG5jb25zdCBnbG9iYWxzID0gZXh0ZXJuYWxHbG9iYWxzKHt9KVxuXG4vLyBodHRwczovL3ZpdGVqcy5kZXYvY29uZmlnL1xuZXhwb3J0IGRlZmF1bHQgZGVmaW5lQ29uZmlnKChvcHQpID0+IHtcbiAgY29uc3QgeyBjb21tYW5kLCBtb2RlIH0gPSBvcHRcbiAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lIG5vLXVudXNlZC12YXJzXG4gIGNvbnN0IGVudiA9IGxvYWRFbnYobW9kZSwgcHJvY2Vzcy5jd2QoKSwgJycpXG4gIGNvbnN0IGJhc2UgPSBjb21tYW5kID09PSAnc2VydmUnID8gJy8nIDogJy8nXG5cbiAgcmV0dXJuIHtcbiAgICBwbHVnaW5zOiBbXG4gICAgICB2dWUoKSxcbiAgICAgIHZ1ZUpzeCgpLFxuICAgICAgLy8gVnVlRGV2VG9vbHMoKSxcbiAgICAgIENvbXBvbmVudHMoe1xuICAgICAgICByZXNvbHZlcnM6IFtcbiAgICAgICAgICBBbnREZXNpZ25WdWVSZXNvbHZlcih7XG4gICAgICAgICAgICBpbXBvcnRTdHlsZTogZmFsc2UgLy8gY3NzIGluIGpzXG4gICAgICAgICAgfSlcbiAgICAgICAgXVxuICAgICAgfSksXG4gICAgICBjcmVhdGVTdmdJY29uc1BsdWdpbih7XG4gICAgICAgIC8vIFx1NjMwN1x1NUI5QVx1OTcwMFx1ODk4MVx1N0YxM1x1NUI1OFx1NzY4NFx1NTZGRVx1NjgwN1x1NjU4N1x1NEVGNlx1NTkzOVxuICAgICAgICBpY29uRGlyczogW3Jlc29sdmUoJy4vc3JjL2Fzc2V0cy9zdmcnKV0sXG4gICAgICAgIC8vIFx1NjMwN1x1NUI5QXN5bWJvbElkXHU2ODNDXHU1RjBGXG4gICAgICAgIHN5bWJvbElkOiAnW25hbWVdJ1xuICAgICAgfSksXG4gICAgICBWdWVJMThuUGx1Z2luKHtcbiAgICAgICAgcnVudGltZU9ubHk6IHRydWUsXG4gICAgICAgIGNvbXBvc2l0aW9uT25seTogdHJ1ZSxcbiAgICAgICAgaW5jbHVkZTogW3BhdGgucmVzb2x2ZShfX2Rpcm5hbWUsICcuL3NyYy9sb2NhbGVzL2xhbmcqKicpXVxuICAgICAgfSlcbiAgICAgIC8vIGNvcHlJbmRleCgpLFxuICAgIF0sXG4gICAgcmVzb2x2ZToge1xuICAgICAgYWxpYXM6IFtcbiAgICAgICAge1xuICAgICAgICAgIGZpbmQ6ICd2dWUtaTE4bicsXG4gICAgICAgICAgcmVwbGFjZW1lbnQ6ICd2dWUtaTE4bi9kaXN0L3Z1ZS1pMThuLmNqcy5qcydcbiAgICAgICAgfSxcbiAgICAgICAge1xuICAgICAgICAgIGZpbmQ6IC9AXFwvLyxcbiAgICAgICAgICByZXBsYWNlbWVudDogYCR7cGF0aFJlc29sdmUoJ3NyYycpfS9gXG4gICAgICAgIH1cbiAgICAgIF1cbiAgICB9LFxuICAgIGJhc2U6IGJhc2UsXG4gICAgZXhwZXJpbWVudGFsOiB7XG4gICAgICAvLyBcdThGREJcdTk2MzZcdTU3RkFcdTc4NDBcdThERUZcdTVGODRcdTkwMDlcdTk4NzlcbiAgICB9LFxuICAgIHNlcnZlcjoge1xuICAgICAgcHJveHk6IGdldFByb3h5Q29uZmlnKG9wdClcbiAgICB9LFxuICAgIGVzYnVpbGQ6IHtcbiAgICAgIHB1cmU6IGVudi5WSVRFX0RST1BfQ09OU09MRSA9PT0gJ3RydWUnID8gWydjb25zb2xlLmxvZyddIDogdW5kZWZpbmVkLFxuICAgICAgZHJvcDogZW52LlZJVEVfRFJPUF9ERUJVR0dFUiA9PT0gJ3RydWUnID8gWydkZWJ1Z2dlciddIDogdW5kZWZpbmVkXG4gICAgfSxcbiAgICBidWlsZDoge1xuICAgICAgb3V0RGlyOiBmaWxlVVJMVG9QYXRoKG5ldyBVUkwoJy4uLy4uL3N0YXRpYy9jaGF0LWFpLWFkbWluJywgaW1wb3J0Lm1ldGEudXJsKSksXG4gICAgICBhc3NldHNEaXI6ICdhc3NldHMnLFxuICAgICAgc291cmNlbWFwOiBlbnYuVklURV9TT1VSQ0VNQVAgPT09ICd0cnVlJyxcbiAgICAgIHJvbGx1cE9wdGlvbnM6IHtcbiAgICAgICAgLy8gZXh0ZXJuYWw6IFsnbW9tZW50JywgJ3ZpZGVvLmpzJywgJ2pzcGRmJywgJ3hsc3gnXSxcbiAgICAgICAgZXh0ZXJuYWw6IFtdLFxuICAgICAgICBwbHVnaW5zOiBbZ2xvYmFscywgZW52LlZJVEVfVVNFX0JVTkRMRV9BTkFMWVpFUiA9PT0gJ3RydWUnID8gdmlzdWFsaXplcigpIDogdW5kZWZpbmVkXSxcbiAgICAgICAgb3V0cHV0OiB7XG4gICAgICAgICAgLy8gXHU4MUVBXHU1QjlBXHU0RTQ5Y2h1bmtGaWxlTmFtZVx1NzUxRlx1NjIxMFx1ODlDNFx1NTIxOVxuICAgICAgICAgIGNodW5rRmlsZU5hbWVzOiAnYXNzZXRzL2pzL1tuYW1lXS1baGFzaF0uanMnLFxuICAgICAgICAgIGVudHJ5RmlsZU5hbWVzOiAnW25hbWVdLmpzJyxcbiAgICAgICAgICBhc3NldEZpbGVOYW1lcyhhc3NldEluZm8pIHtcbiAgICAgICAgICAgIGxldCBmaWVsX25hbWUgPSBhc3NldEluZm8ubmFtZS50b0xvd2VyQ2FzZSgpXG5cbiAgICAgICAgICAgIGlmIChmaWVsX25hbWUuZW5kc1dpdGgoJy5jc3MnKSkge1xuICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9jc3MvW25hbWVdLVtoYXNoXS5bZXh0XSdcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGlmIChbJ3BuZycsICdqcGcnLCAnanBlZycsICdzdmcnXS5zb21lKChleHQpID0+IGZpZWxfbmFtZS5lbmRzV2l0aChleHQpKSkge1xuICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9pbWcvW25hbWVdLVtoYXNoXS5bZXh0XSdcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGlmIChbJ3R0ZicsICd3b2ZmJywgJ3dvZmYyJ10uc29tZSgoZXh0KSA9PiBmaWVsX25hbWUuZW5kc1dpdGgoZXh0KSkpIHtcbiAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvZm9udHMvW25hbWVdLVtoYXNoXS5bZXh0XSdcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHJldHVybiAnYXNzZXRzL1tuYW1lXS1baGFzaF0uW2V4dF0nXG4gICAgICAgICAgfSxcbiAgICAgICAgICAvLyBcdThCRTVcdTkwMDlcdTk4NzlcdTUxNDFcdThCQjhcdTRGNjBcdTUyMUJcdTVFRkFcdTgxRUFcdTVCOUFcdTRFNDlcdTc2ODRcdTUxNkNcdTUxNzEgY2h1bmtcbiAgICAgICAgICBtYW51YWxDaHVua3M6IHtcbiAgICAgICAgICAgICd2dWUtY2h1bmtzJzogWyd2dWUnLCAndnVlLXJvdXRlcicsICdwaW5pYScsICd2dWUtaTE4biddLFxuICAgICAgICAgICAgZGF5anM6IFsnZGF5anMnXSxcbiAgICAgICAgICAgIGF4aW9zOiBbJ2F4aW9zJ10sXG4gICAgICAgICAgICAnY3J5cHRvLWpzJzogWydjcnlwdG8tanMnXSxcbiAgICAgICAgICAgIHFzOiBbJ3FzJ10sXG4gICAgICAgICAgICAndnVlLXBkZi1lbWJlZCc6IFsndnVlLXBkZi1lbWJlZCddXG4gICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICB9LFxuICAgICAgY3NzQ29kZVNwbGl0OiAhKGVudi5WSVRFX1VTRV9DU1NfU1BMSVQgPT09ICdmYWxzZScpXG4gICAgfSxcbiAgICBvcHRpbWl6ZURlcHM6IHtcbiAgICAgIGluY2x1ZGU6IFsndnVlJywgJ3Z1ZS1yb3V0ZXInLCAncGluaWEnLCAndnVlLWkxOG4nLCAnZGF5anMnLCAnYXhpb3MnLCAnY3J5cHRvLWpzJywgJ3FzJ11cbiAgICB9XG4gIH1cbn0pXG4iLCAiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIi9Vc2Vycy9haXhpYW5nZmVpL2NvZGUvemhpbWFfY2hhdF9haS9mcm9udC1lbmQvY2hhdC1haS1hZG1pbi12dWVcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIi9Vc2Vycy9haXhpYW5nZmVpL2NvZGUvemhpbWFfY2hhdF9haS9mcm9udC1lbmQvY2hhdC1haS1hZG1pbi12dWUvcHJveHlfY29uZmlnLmpzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9Vc2Vycy9haXhpYW5nZmVpL2NvZGUvemhpbWFfY2hhdF9haS9mcm9udC1lbmQvY2hhdC1haS1hZG1pbi12dWUvcHJveHlfY29uZmlnLmpzXCI7LyogZXNsaW50LWRpc2FibGUgbm8tdW5kZWYgKi9cbmltcG9ydCB7IGxvYWRFbnYgfSBmcm9tICd2aXRlJ1xuXG5leHBvcnQgY29uc3QgZ2V0UHJveHlDb25maWcgPSAob3B0KSA9PiB7XG4gIGNvbnN0IHsgbW9kZSB9ID0gb3B0XG4gIGNvbnN0IGVudiA9IGxvYWRFbnYobW9kZSwgcHJvY2Vzcy5jd2QoKSwgJycpXG5cbiAgbGV0IHByb3h5QXBpcyA9IFsnL3N0YXRpYycsICcvY29tbW9uJywgJy9tYW5hZ2UnLCAnL2FwcCcsICcvY2hhdCcsICcvdXBsb2FkJ11cbiAgbGV0IHByb3h5ID0ge31cblxuICBjb25zb2xlLmxvZyhlbnYuUFJPWFlfQkFTRV9BUElfVVJMKVxuXG4gIHByb3h5QXBpcy5mb3JFYWNoKChrZXkpID0+IHtcbiAgICBwcm94eVtrZXldID0ge1xuICAgICAgdGFyZ2V0OiBlbnYuUFJPWFlfQkFTRV9BUElfVVJMLFxuICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlXG4gICAgfVxuICB9KVxuXG4gIHJldHVybiB7XG4gICAgLi4ucHJveHlcbiAgfVxufVxuIl0sCiAgIm1hcHBpbmdzIjogIjtBQUNBLFNBQVMsZUFBZSxXQUFXO0FBQ25DLE9BQU8sVUFBVTtBQUVqQixTQUFTLGNBQWMsV0FBQUEsZ0JBQWU7QUFDdEMsT0FBTyxTQUFTO0FBQ2hCLE9BQU8sWUFBWTtBQUVuQixTQUFTLGtCQUFrQjtBQUMzQixPQUFPLHFCQUFxQjtBQUM1QixPQUFPLGdCQUFnQjtBQUN2QixTQUFTLDRCQUE0QjtBQUNyQyxPQUFPLG1CQUFtQjtBQUMxQixTQUFTLDRCQUE0Qjs7O0FDWnJDLFNBQVMsZUFBZTtBQUVqQixJQUFNLGlCQUFpQixDQUFDLFFBQVE7QUFDckMsUUFBTSxFQUFFLEtBQUssSUFBSTtBQUNqQixRQUFNLE1BQU0sUUFBUSxNQUFNLFFBQVEsSUFBSSxHQUFHLEVBQUU7QUFFM0MsTUFBSSxZQUFZLENBQUMsV0FBVyxXQUFXLFdBQVcsUUFBUSxTQUFTLFNBQVM7QUFDNUUsTUFBSSxRQUFRLENBQUM7QUFFYixVQUFRLElBQUksSUFBSSxrQkFBa0I7QUFFbEMsWUFBVSxRQUFRLENBQUMsUUFBUTtBQUN6QixVQUFNLEdBQUcsSUFBSTtBQUFBLE1BQ1gsUUFBUSxJQUFJO0FBQUEsTUFDWixjQUFjO0FBQUEsSUFDaEI7QUFBQSxFQUNGLENBQUM7QUFFRCxTQUFPO0FBQUEsSUFDTCxHQUFHO0FBQUEsRUFDTDtBQUNGOzs7QUR0QkEsSUFBTSxtQ0FBbUM7QUFBK0wsSUFBTSwyQ0FBMkM7QUFnQnpSLElBQU0sRUFBRSxRQUFRLElBQUk7QUFDcEIsSUFBTSxPQUFPLFFBQVEsSUFBSTtBQUN6QixTQUFTLFlBQVksS0FBSztBQUN4QixTQUFPLFFBQVEsTUFBTSxLQUFLLEdBQUc7QUFDL0I7QUFTQSxJQUFNLFVBQVUsZ0JBQWdCLENBQUMsQ0FBQztBQUdsQyxJQUFPLHNCQUFRLGFBQWEsQ0FBQyxRQUFRO0FBQ25DLFFBQU0sRUFBRSxTQUFTLEtBQUssSUFBSTtBQUUxQixRQUFNLE1BQU1DLFNBQVEsTUFBTSxRQUFRLElBQUksR0FBRyxFQUFFO0FBQzNDLFFBQU0sT0FBTyxZQUFZLFVBQVUsTUFBTTtBQUV6QyxTQUFPO0FBQUEsSUFDTCxTQUFTO0FBQUEsTUFDUCxJQUFJO0FBQUEsTUFDSixPQUFPO0FBQUE7QUFBQSxNQUVQLFdBQVc7QUFBQSxRQUNULFdBQVc7QUFBQSxVQUNULHFCQUFxQjtBQUFBLFlBQ25CLGFBQWE7QUFBQTtBQUFBLFVBQ2YsQ0FBQztBQUFBLFFBQ0g7QUFBQSxNQUNGLENBQUM7QUFBQSxNQUNELHFCQUFxQjtBQUFBO0FBQUEsUUFFbkIsVUFBVSxDQUFDLFFBQVEsa0JBQWtCLENBQUM7QUFBQTtBQUFBLFFBRXRDLFVBQVU7QUFBQSxNQUNaLENBQUM7QUFBQSxNQUNELGNBQWM7QUFBQSxRQUNaLGFBQWE7QUFBQSxRQUNiLGlCQUFpQjtBQUFBLFFBQ2pCLFNBQVMsQ0FBQyxLQUFLLFFBQVEsa0NBQVcsc0JBQXNCLENBQUM7QUFBQSxNQUMzRCxDQUFDO0FBQUE7QUFBQSxJQUVIO0FBQUEsSUFDQSxTQUFTO0FBQUEsTUFDUCxPQUFPO0FBQUEsUUFDTDtBQUFBLFVBQ0UsTUFBTTtBQUFBLFVBQ04sYUFBYTtBQUFBLFFBQ2Y7QUFBQSxRQUNBO0FBQUEsVUFDRSxNQUFNO0FBQUEsVUFDTixhQUFhLEdBQUcsWUFBWSxLQUFLLENBQUM7QUFBQSxRQUNwQztBQUFBLE1BQ0Y7QUFBQSxJQUNGO0FBQUEsSUFDQTtBQUFBLElBQ0EsY0FBYztBQUFBO0FBQUEsSUFFZDtBQUFBLElBQ0EsUUFBUTtBQUFBLE1BQ04sT0FBTyxlQUFlLEdBQUc7QUFBQSxJQUMzQjtBQUFBLElBQ0EsU0FBUztBQUFBLE1BQ1AsTUFBTSxJQUFJLHNCQUFzQixTQUFTLENBQUMsYUFBYSxJQUFJO0FBQUEsTUFDM0QsTUFBTSxJQUFJLHVCQUF1QixTQUFTLENBQUMsVUFBVSxJQUFJO0FBQUEsSUFDM0Q7QUFBQSxJQUNBLE9BQU87QUFBQSxNQUNMLFFBQVEsY0FBYyxJQUFJLElBQUksOEJBQThCLHdDQUFlLENBQUM7QUFBQSxNQUM1RSxXQUFXO0FBQUEsTUFDWCxXQUFXLElBQUksbUJBQW1CO0FBQUEsTUFDbEMsZUFBZTtBQUFBO0FBQUEsUUFFYixVQUFVLENBQUM7QUFBQSxRQUNYLFNBQVMsQ0FBQyxTQUFTLElBQUksNkJBQTZCLFNBQVMsV0FBVyxJQUFJLE1BQVM7QUFBQSxRQUNyRixRQUFRO0FBQUE7QUFBQSxVQUVOLGdCQUFnQjtBQUFBLFVBQ2hCLGdCQUFnQjtBQUFBLFVBQ2hCLGVBQWUsV0FBVztBQUN4QixnQkFBSSxZQUFZLFVBQVUsS0FBSyxZQUFZO0FBRTNDLGdCQUFJLFVBQVUsU0FBUyxNQUFNLEdBQUc7QUFDOUIscUJBQU87QUFBQSxZQUNUO0FBQ0EsZ0JBQUksQ0FBQyxPQUFPLE9BQU8sUUFBUSxLQUFLLEVBQUUsS0FBSyxDQUFDLFFBQVEsVUFBVSxTQUFTLEdBQUcsQ0FBQyxHQUFHO0FBQ3hFLHFCQUFPO0FBQUEsWUFDVDtBQUNBLGdCQUFJLENBQUMsT0FBTyxRQUFRLE9BQU8sRUFBRSxLQUFLLENBQUMsUUFBUSxVQUFVLFNBQVMsR0FBRyxDQUFDLEdBQUc7QUFDbkUscUJBQU87QUFBQSxZQUNUO0FBQ0EsbUJBQU87QUFBQSxVQUNUO0FBQUE7QUFBQSxVQUVBLGNBQWM7QUFBQSxZQUNaLGNBQWMsQ0FBQyxPQUFPLGNBQWMsU0FBUyxVQUFVO0FBQUEsWUFDdkQsT0FBTyxDQUFDLE9BQU87QUFBQSxZQUNmLE9BQU8sQ0FBQyxPQUFPO0FBQUEsWUFDZixhQUFhLENBQUMsV0FBVztBQUFBLFlBQ3pCLElBQUksQ0FBQyxJQUFJO0FBQUEsWUFDVCxpQkFBaUIsQ0FBQyxlQUFlO0FBQUEsVUFDbkM7QUFBQSxRQUNGO0FBQUEsTUFDRjtBQUFBLE1BQ0EsY0FBYyxFQUFFLElBQUksdUJBQXVCO0FBQUEsSUFDN0M7QUFBQSxJQUNBLGNBQWM7QUFBQSxNQUNaLFNBQVMsQ0FBQyxPQUFPLGNBQWMsU0FBUyxZQUFZLFNBQVMsU0FBUyxhQUFhLElBQUk7QUFBQSxJQUN6RjtBQUFBLEVBQ0Y7QUFDRixDQUFDOyIsCiAgIm5hbWVzIjogWyJsb2FkRW52IiwgImxvYWRFbnYiXQp9Cg==
