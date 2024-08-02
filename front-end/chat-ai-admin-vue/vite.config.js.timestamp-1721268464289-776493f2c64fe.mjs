// vite.config.js
import { fileURLToPath, URL } from "node:url";
import path from "path";
import { defineConfig, loadEnv as loadEnv2 } from "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/vite/dist/node/index.js";
import vue from "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import { visualizer } from "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/rollup-plugin-visualizer/dist/plugin/index.js";
import externalGlobals from "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/rollup-plugin-external-globals/index.js";
import Components from "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/unplugin-vue-components/dist/vite.js";
import { AntDesignVueResolver } from "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/unplugin-vue-components/dist/resolvers.js";
import VueI18nPlugin from "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/@intlify/unplugin-vue-i18n/lib/vite.mjs";
import { createSvgIconsPlugin } from "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/vite-plugin-svg-icons/dist/index.mjs";

// proxy_config.js
import { loadEnv } from "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/vite/dist/node/index.js";
var getProxyConfig = (opt) => {
  const { mode } = opt;
  const env = loadEnv(mode, process.cwd(), "");
  let proxyApis = ["/static", "/common", "/manage", "/app", "/chat", "/upload", "/public"];
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
var __vite_injected_original_dirname = "/Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue";
var __vite_injected_original_import_meta_url = "file:///Users/aixiangfei/code/zhima/zhima_chat_ai/front-end/chat-ai-admin-vue/vite.config.js";
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
      emptyOutDir: true,
      assetsDir: "assets",
      sourcemap: env.VITE_SOURCEMAP === "true",
      rollupOptions: {
        // external: ['moment', 'video.js', 'jspdf', 'xlsx'],
        external: [],
        plugins: [globals, env.VITE_USE_BUNDLE_ANALYZER === "true" ? visualizer() : void 0],
        output: {
          // 自定义chunkFileName生成规则
          chunkFileNames: "assets/js/[name]-[hash].js",
          entryFileNames: "[name]-[hash].js",
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiLCAicHJveHlfY29uZmlnLmpzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiL1VzZXJzL2FpeGlhbmdmZWkvY29kZS96aGltYS96aGltYV9jaGF0X2FpL2Zyb250LWVuZC9jaGF0LWFpLWFkbWluLXZ1ZVwiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiL1VzZXJzL2FpeGlhbmdmZWkvY29kZS96aGltYS96aGltYV9jaGF0X2FpL2Zyb250LWVuZC9jaGF0LWFpLWFkbWluLXZ1ZS92aXRlLmNvbmZpZy5qc1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vVXNlcnMvYWl4aWFuZ2ZlaS9jb2RlL3poaW1hL3poaW1hX2NoYXRfYWkvZnJvbnQtZW5kL2NoYXQtYWktYWRtaW4tdnVlL3ZpdGUuY29uZmlnLmpzXCI7LyogZXNsaW50LWRpc2FibGUgbm8tdW5kZWYgKi9cbmltcG9ydCB7IGZpbGVVUkxUb1BhdGgsIFVSTCB9IGZyb20gJ25vZGU6dXJsJ1xuaW1wb3J0IHBhdGggZnJvbSAncGF0aCdcblxuaW1wb3J0IHsgZGVmaW5lQ29uZmlnLCBsb2FkRW52IH0gZnJvbSAndml0ZSdcbmltcG9ydCB2dWUgZnJvbSAnQHZpdGVqcy9wbHVnaW4tdnVlJ1xuaW1wb3J0IHZ1ZUpzeCBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUtanN4J1xuLy8gaW1wb3J0IFZ1ZURldlRvb2xzIGZyb20gJ3ZpdGUtcGx1Z2luLXZ1ZS1kZXZ0b29scydcbmltcG9ydCB7IHZpc3VhbGl6ZXIgfSBmcm9tICdyb2xsdXAtcGx1Z2luLXZpc3VhbGl6ZXInXG5pbXBvcnQgZXh0ZXJuYWxHbG9iYWxzIGZyb20gJ3JvbGx1cC1wbHVnaW4tZXh0ZXJuYWwtZ2xvYmFscydcbmltcG9ydCBDb21wb25lbnRzIGZyb20gJ3VucGx1Z2luLXZ1ZS1jb21wb25lbnRzL3ZpdGUnXG5pbXBvcnQgeyBBbnREZXNpZ25WdWVSZXNvbHZlciB9IGZyb20gJ3VucGx1Z2luLXZ1ZS1jb21wb25lbnRzL3Jlc29sdmVycydcbmltcG9ydCBWdWVJMThuUGx1Z2luIGZyb20gJ0BpbnRsaWZ5L3VucGx1Z2luLXZ1ZS1pMThuL3ZpdGUnXG5pbXBvcnQgeyBjcmVhdGVTdmdJY29uc1BsdWdpbiB9IGZyb20gJ3ZpdGUtcGx1Z2luLXN2Zy1pY29ucydcbmltcG9ydCB7IGdldFByb3h5Q29uZmlnIH0gZnJvbSAnLi9wcm94eV9jb25maWcnXG5cbmNvbnN0IHsgcmVzb2x2ZSB9ID0gcGF0aFxuY29uc3Qgcm9vdCA9IHByb2Nlc3MuY3dkKClcbmZ1bmN0aW9uIHBhdGhSZXNvbHZlKGRpcikge1xuICByZXR1cm4gcmVzb2x2ZShyb290LCAnLicsIGRpcilcbn1cblxuLy8gY29uc3QgZ2xvYmFscyA9IGV4dGVybmFsR2xvYmFscyh7XG4vLyAgIG1vbWVudDogJ21vbWVudCcsXG4vLyAgICd2aWRlby5qcyc6ICd2aWRlb2pzJyxcbi8vICAganNwZGY6ICdqc3BkZicsXG4vLyAgIHhsc3g6ICdYTFNYJyxcbi8vIH0pO1xuXG5jb25zdCBnbG9iYWxzID0gZXh0ZXJuYWxHbG9iYWxzKHt9KVxuXG4vLyBodHRwczovL3ZpdGVqcy5kZXYvY29uZmlnL1xuZXhwb3J0IGRlZmF1bHQgZGVmaW5lQ29uZmlnKChvcHQpID0+IHtcbiAgY29uc3QgeyBjb21tYW5kLCBtb2RlIH0gPSBvcHRcbiAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lIG5vLXVudXNlZC12YXJzXG4gIGNvbnN0IGVudiA9IGxvYWRFbnYobW9kZSwgcHJvY2Vzcy5jd2QoKSwgJycpXG4gIGNvbnN0IGJhc2UgPSBjb21tYW5kID09PSAnc2VydmUnID8gJy8nIDogJy8nXG5cbiAgcmV0dXJuIHtcbiAgICBwbHVnaW5zOiBbXG4gICAgICB2dWUoKSxcbiAgICAgIHZ1ZUpzeCgpLFxuICAgICAgLy8gVnVlRGV2VG9vbHMoKSxcbiAgICAgIENvbXBvbmVudHMoe1xuICAgICAgICByZXNvbHZlcnM6IFtcbiAgICAgICAgICBBbnREZXNpZ25WdWVSZXNvbHZlcih7XG4gICAgICAgICAgICBpbXBvcnRTdHlsZTogZmFsc2UgLy8gY3NzIGluIGpzXG4gICAgICAgICAgfSlcbiAgICAgICAgXVxuICAgICAgfSksXG4gICAgICBjcmVhdGVTdmdJY29uc1BsdWdpbih7XG4gICAgICAgIC8vIFx1NjMwN1x1NUI5QVx1OTcwMFx1ODk4MVx1N0YxM1x1NUI1OFx1NzY4NFx1NTZGRVx1NjgwN1x1NjU4N1x1NEVGNlx1NTkzOVxuICAgICAgICBpY29uRGlyczogW3Jlc29sdmUoJy4vc3JjL2Fzc2V0cy9zdmcnKV0sXG4gICAgICAgIC8vIFx1NjMwN1x1NUI5QXN5bWJvbElkXHU2ODNDXHU1RjBGXG4gICAgICAgIHN5bWJvbElkOiAnW25hbWVdJ1xuICAgICAgfSksXG4gICAgICBWdWVJMThuUGx1Z2luKHtcbiAgICAgICAgcnVudGltZU9ubHk6IHRydWUsXG4gICAgICAgIGNvbXBvc2l0aW9uT25seTogdHJ1ZSxcbiAgICAgICAgaW5jbHVkZTogW3BhdGgucmVzb2x2ZShfX2Rpcm5hbWUsICcuL3NyYy9sb2NhbGVzL2xhbmcqKicpXVxuICAgICAgfSlcbiAgICAgIC8vIGNvcHlJbmRleCgpLFxuICAgIF0sXG4gICAgcmVzb2x2ZToge1xuICAgICAgYWxpYXM6IFtcbiAgICAgICAge1xuICAgICAgICAgIGZpbmQ6ICd2dWUtaTE4bicsXG4gICAgICAgICAgcmVwbGFjZW1lbnQ6ICd2dWUtaTE4bi9kaXN0L3Z1ZS1pMThuLmNqcy5qcydcbiAgICAgICAgfSxcbiAgICAgICAge1xuICAgICAgICAgIGZpbmQ6IC9AXFwvLyxcbiAgICAgICAgICByZXBsYWNlbWVudDogYCR7cGF0aFJlc29sdmUoJ3NyYycpfS9gXG4gICAgICAgIH1cbiAgICAgIF1cbiAgICB9LFxuICAgIGJhc2U6IGJhc2UsXG4gICAgZXhwZXJpbWVudGFsOiB7XG4gICAgICAvLyBcdThGREJcdTk2MzZcdTU3RkFcdTc4NDBcdThERUZcdTVGODRcdTkwMDlcdTk4NzlcbiAgICB9LFxuICAgIHNlcnZlcjoge1xuICAgICAgcHJveHk6IGdldFByb3h5Q29uZmlnKG9wdClcbiAgICB9LFxuICAgIGVzYnVpbGQ6IHtcbiAgICAgIHB1cmU6IGVudi5WSVRFX0RST1BfQ09OU09MRSA9PT0gJ3RydWUnID8gWydjb25zb2xlLmxvZyddIDogdW5kZWZpbmVkLFxuICAgICAgZHJvcDogZW52LlZJVEVfRFJPUF9ERUJVR0dFUiA9PT0gJ3RydWUnID8gWydkZWJ1Z2dlciddIDogdW5kZWZpbmVkXG4gICAgfSxcbiAgICBidWlsZDoge1xuICAgICAgb3V0RGlyOiBmaWxlVVJMVG9QYXRoKG5ldyBVUkwoJy4uLy4uL3N0YXRpYy9jaGF0LWFpLWFkbWluJywgaW1wb3J0Lm1ldGEudXJsKSksXG4gICAgICBlbXB0eU91dERpcjogdHJ1ZSxcbiAgICAgIGFzc2V0c0RpcjogJ2Fzc2V0cycsXG4gICAgICBzb3VyY2VtYXA6IGVudi5WSVRFX1NPVVJDRU1BUCA9PT0gJ3RydWUnLFxuICAgICAgcm9sbHVwT3B0aW9uczoge1xuICAgICAgICAvLyBleHRlcm5hbDogWydtb21lbnQnLCAndmlkZW8uanMnLCAnanNwZGYnLCAneGxzeCddLFxuICAgICAgICBleHRlcm5hbDogW10sXG4gICAgICAgIHBsdWdpbnM6IFtnbG9iYWxzLCBlbnYuVklURV9VU0VfQlVORExFX0FOQUxZWkVSID09PSAndHJ1ZScgPyB2aXN1YWxpemVyKCkgOiB1bmRlZmluZWRdLFxuICAgICAgICBvdXRwdXQ6IHtcbiAgICAgICAgICAvLyBcdTgxRUFcdTVCOUFcdTRFNDljaHVua0ZpbGVOYW1lXHU3NTFGXHU2MjEwXHU4OUM0XHU1MjE5XG4gICAgICAgICAgY2h1bmtGaWxlTmFtZXM6ICdhc3NldHMvanMvW25hbWVdLVtoYXNoXS5qcycsXG4gICAgICAgICAgZW50cnlGaWxlTmFtZXM6ICdbbmFtZV0tW2hhc2hdLmpzJyxcbiAgICAgICAgICBhc3NldEZpbGVOYW1lcyhhc3NldEluZm8pIHtcbiAgICAgICAgICAgIGxldCBmaWVsX25hbWUgPSBhc3NldEluZm8ubmFtZS50b0xvd2VyQ2FzZSgpXG5cbiAgICAgICAgICAgIGlmIChmaWVsX25hbWUuZW5kc1dpdGgoJy5jc3MnKSkge1xuICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9jc3MvW25hbWVdLVtoYXNoXS5bZXh0XSdcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGlmIChbJ3BuZycsICdqcGcnLCAnanBlZycsICdzdmcnXS5zb21lKChleHQpID0+IGZpZWxfbmFtZS5lbmRzV2l0aChleHQpKSkge1xuICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9pbWcvW25hbWVdLVtoYXNoXS5bZXh0XSdcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGlmIChbJ3R0ZicsICd3b2ZmJywgJ3dvZmYyJ10uc29tZSgoZXh0KSA9PiBmaWVsX25hbWUuZW5kc1dpdGgoZXh0KSkpIHtcbiAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvZm9udHMvW25hbWVdLVtoYXNoXS5bZXh0XSdcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHJldHVybiAnYXNzZXRzL1tuYW1lXS1baGFzaF0uW2V4dF0nXG4gICAgICAgICAgfSxcbiAgICAgICAgICAvLyBcdThCRTVcdTkwMDlcdTk4NzlcdTUxNDFcdThCQjhcdTRGNjBcdTUyMUJcdTVFRkFcdTgxRUFcdTVCOUFcdTRFNDlcdTc2ODRcdTUxNkNcdTUxNzEgY2h1bmtcbiAgICAgICAgICBtYW51YWxDaHVua3M6IHtcbiAgICAgICAgICAgICd2dWUtY2h1bmtzJzogWyd2dWUnLCAndnVlLXJvdXRlcicsICdwaW5pYScsICd2dWUtaTE4biddLFxuICAgICAgICAgICAgZGF5anM6IFsnZGF5anMnXSxcbiAgICAgICAgICAgIGF4aW9zOiBbJ2F4aW9zJ10sXG4gICAgICAgICAgICAnY3J5cHRvLWpzJzogWydjcnlwdG8tanMnXSxcbiAgICAgICAgICAgIHFzOiBbJ3FzJ10sXG4gICAgICAgICAgICAndnVlLXBkZi1lbWJlZCc6IFsndnVlLXBkZi1lbWJlZCddXG4gICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICB9LFxuICAgICAgY3NzQ29kZVNwbGl0OiAhKGVudi5WSVRFX1VTRV9DU1NfU1BMSVQgPT09ICdmYWxzZScpXG4gICAgfSxcbiAgICBvcHRpbWl6ZURlcHM6IHtcbiAgICAgIGluY2x1ZGU6IFsndnVlJywgJ3Z1ZS1yb3V0ZXInLCAncGluaWEnLCAndnVlLWkxOG4nLCAnZGF5anMnLCAnYXhpb3MnLCAnY3J5cHRvLWpzJywgJ3FzJ11cbiAgICB9XG4gIH1cbn0pXG4iLCAiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIi9Vc2Vycy9haXhpYW5nZmVpL2NvZGUvemhpbWEvemhpbWFfY2hhdF9haS9mcm9udC1lbmQvY2hhdC1haS1hZG1pbi12dWVcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIi9Vc2Vycy9haXhpYW5nZmVpL2NvZGUvemhpbWEvemhpbWFfY2hhdF9haS9mcm9udC1lbmQvY2hhdC1haS1hZG1pbi12dWUvcHJveHlfY29uZmlnLmpzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9Vc2Vycy9haXhpYW5nZmVpL2NvZGUvemhpbWEvemhpbWFfY2hhdF9haS9mcm9udC1lbmQvY2hhdC1haS1hZG1pbi12dWUvcHJveHlfY29uZmlnLmpzXCI7LyogZXNsaW50LWRpc2FibGUgbm8tdW5kZWYgKi9cbmltcG9ydCB7IGxvYWRFbnYgfSBmcm9tICd2aXRlJ1xuXG5leHBvcnQgY29uc3QgZ2V0UHJveHlDb25maWcgPSAob3B0KSA9PiB7XG4gIGNvbnN0IHsgbW9kZSB9ID0gb3B0XG4gIGNvbnN0IGVudiA9IGxvYWRFbnYobW9kZSwgcHJvY2Vzcy5jd2QoKSwgJycpXG5cbiAgbGV0IHByb3h5QXBpcyA9IFsnL3N0YXRpYycsICcvY29tbW9uJywgJy9tYW5hZ2UnLCAnL2FwcCcsICcvY2hhdCcsICcvdXBsb2FkJywgJy9wdWJsaWMnXVxuICBsZXQgcHJveHkgPSB7fVxuXG4gIGNvbnNvbGUubG9nKGVudi5QUk9YWV9CQVNFX0FQSV9VUkwpXG5cbiAgcHJveHlBcGlzLmZvckVhY2goKGtleSkgPT4ge1xuICAgIHByb3h5W2tleV0gPSB7XG4gICAgICB0YXJnZXQ6IGVudi5QUk9YWV9CQVNFX0FQSV9VUkwsXG4gICAgICBjaGFuZ2VPcmlnaW46IHRydWVcbiAgICB9XG4gIH0pXG5cbiAgcmV0dXJuIHtcbiAgICAuLi5wcm94eVxuICB9XG59XG4iXSwKICAibWFwcGluZ3MiOiAiO0FBQ0EsU0FBUyxlQUFlLFdBQVc7QUFDbkMsT0FBTyxVQUFVO0FBRWpCLFNBQVMsY0FBYyxXQUFBQSxnQkFBZTtBQUN0QyxPQUFPLFNBQVM7QUFDaEIsT0FBTyxZQUFZO0FBRW5CLFNBQVMsa0JBQWtCO0FBQzNCLE9BQU8scUJBQXFCO0FBQzVCLE9BQU8sZ0JBQWdCO0FBQ3ZCLFNBQVMsNEJBQTRCO0FBQ3JDLE9BQU8sbUJBQW1CO0FBQzFCLFNBQVMsNEJBQTRCOzs7QUNackMsU0FBUyxlQUFlO0FBRWpCLElBQU0saUJBQWlCLENBQUMsUUFBUTtBQUNyQyxRQUFNLEVBQUUsS0FBSyxJQUFJO0FBQ2pCLFFBQU0sTUFBTSxRQUFRLE1BQU0sUUFBUSxJQUFJLEdBQUcsRUFBRTtBQUUzQyxNQUFJLFlBQVksQ0FBQyxXQUFXLFdBQVcsV0FBVyxRQUFRLFNBQVMsV0FBVyxTQUFTO0FBQ3ZGLE1BQUksUUFBUSxDQUFDO0FBRWIsVUFBUSxJQUFJLElBQUksa0JBQWtCO0FBRWxDLFlBQVUsUUFBUSxDQUFDLFFBQVE7QUFDekIsVUFBTSxHQUFHLElBQUk7QUFBQSxNQUNYLFFBQVEsSUFBSTtBQUFBLE1BQ1osY0FBYztBQUFBLElBQ2hCO0FBQUEsRUFDRixDQUFDO0FBRUQsU0FBTztBQUFBLElBQ0wsR0FBRztBQUFBLEVBQ0w7QUFDRjs7O0FEdEJBLElBQU0sbUNBQW1DO0FBQTJNLElBQU0sMkNBQTJDO0FBZ0JyUyxJQUFNLEVBQUUsUUFBUSxJQUFJO0FBQ3BCLElBQU0sT0FBTyxRQUFRLElBQUk7QUFDekIsU0FBUyxZQUFZLEtBQUs7QUFDeEIsU0FBTyxRQUFRLE1BQU0sS0FBSyxHQUFHO0FBQy9CO0FBU0EsSUFBTSxVQUFVLGdCQUFnQixDQUFDLENBQUM7QUFHbEMsSUFBTyxzQkFBUSxhQUFhLENBQUMsUUFBUTtBQUNuQyxRQUFNLEVBQUUsU0FBUyxLQUFLLElBQUk7QUFFMUIsUUFBTSxNQUFNQyxTQUFRLE1BQU0sUUFBUSxJQUFJLEdBQUcsRUFBRTtBQUMzQyxRQUFNLE9BQU8sWUFBWSxVQUFVLE1BQU07QUFFekMsU0FBTztBQUFBLElBQ0wsU0FBUztBQUFBLE1BQ1AsSUFBSTtBQUFBLE1BQ0osT0FBTztBQUFBO0FBQUEsTUFFUCxXQUFXO0FBQUEsUUFDVCxXQUFXO0FBQUEsVUFDVCxxQkFBcUI7QUFBQSxZQUNuQixhQUFhO0FBQUE7QUFBQSxVQUNmLENBQUM7QUFBQSxRQUNIO0FBQUEsTUFDRixDQUFDO0FBQUEsTUFDRCxxQkFBcUI7QUFBQTtBQUFBLFFBRW5CLFVBQVUsQ0FBQyxRQUFRLGtCQUFrQixDQUFDO0FBQUE7QUFBQSxRQUV0QyxVQUFVO0FBQUEsTUFDWixDQUFDO0FBQUEsTUFDRCxjQUFjO0FBQUEsUUFDWixhQUFhO0FBQUEsUUFDYixpQkFBaUI7QUFBQSxRQUNqQixTQUFTLENBQUMsS0FBSyxRQUFRLGtDQUFXLHNCQUFzQixDQUFDO0FBQUEsTUFDM0QsQ0FBQztBQUFBO0FBQUEsSUFFSDtBQUFBLElBQ0EsU0FBUztBQUFBLE1BQ1AsT0FBTztBQUFBLFFBQ0w7QUFBQSxVQUNFLE1BQU07QUFBQSxVQUNOLGFBQWE7QUFBQSxRQUNmO0FBQUEsUUFDQTtBQUFBLFVBQ0UsTUFBTTtBQUFBLFVBQ04sYUFBYSxHQUFHLFlBQVksS0FBSyxDQUFDO0FBQUEsUUFDcEM7QUFBQSxNQUNGO0FBQUEsSUFDRjtBQUFBLElBQ0E7QUFBQSxJQUNBLGNBQWM7QUFBQTtBQUFBLElBRWQ7QUFBQSxJQUNBLFFBQVE7QUFBQSxNQUNOLE9BQU8sZUFBZSxHQUFHO0FBQUEsSUFDM0I7QUFBQSxJQUNBLFNBQVM7QUFBQSxNQUNQLE1BQU0sSUFBSSxzQkFBc0IsU0FBUyxDQUFDLGFBQWEsSUFBSTtBQUFBLE1BQzNELE1BQU0sSUFBSSx1QkFBdUIsU0FBUyxDQUFDLFVBQVUsSUFBSTtBQUFBLElBQzNEO0FBQUEsSUFDQSxPQUFPO0FBQUEsTUFDTCxRQUFRLGNBQWMsSUFBSSxJQUFJLDhCQUE4Qix3Q0FBZSxDQUFDO0FBQUEsTUFDNUUsYUFBYTtBQUFBLE1BQ2IsV0FBVztBQUFBLE1BQ1gsV0FBVyxJQUFJLG1CQUFtQjtBQUFBLE1BQ2xDLGVBQWU7QUFBQTtBQUFBLFFBRWIsVUFBVSxDQUFDO0FBQUEsUUFDWCxTQUFTLENBQUMsU0FBUyxJQUFJLDZCQUE2QixTQUFTLFdBQVcsSUFBSSxNQUFTO0FBQUEsUUFDckYsUUFBUTtBQUFBO0FBQUEsVUFFTixnQkFBZ0I7QUFBQSxVQUNoQixnQkFBZ0I7QUFBQSxVQUNoQixlQUFlLFdBQVc7QUFDeEIsZ0JBQUksWUFBWSxVQUFVLEtBQUssWUFBWTtBQUUzQyxnQkFBSSxVQUFVLFNBQVMsTUFBTSxHQUFHO0FBQzlCLHFCQUFPO0FBQUEsWUFDVDtBQUNBLGdCQUFJLENBQUMsT0FBTyxPQUFPLFFBQVEsS0FBSyxFQUFFLEtBQUssQ0FBQyxRQUFRLFVBQVUsU0FBUyxHQUFHLENBQUMsR0FBRztBQUN4RSxxQkFBTztBQUFBLFlBQ1Q7QUFDQSxnQkFBSSxDQUFDLE9BQU8sUUFBUSxPQUFPLEVBQUUsS0FBSyxDQUFDLFFBQVEsVUFBVSxTQUFTLEdBQUcsQ0FBQyxHQUFHO0FBQ25FLHFCQUFPO0FBQUEsWUFDVDtBQUNBLG1CQUFPO0FBQUEsVUFDVDtBQUFBO0FBQUEsVUFFQSxjQUFjO0FBQUEsWUFDWixjQUFjLENBQUMsT0FBTyxjQUFjLFNBQVMsVUFBVTtBQUFBLFlBQ3ZELE9BQU8sQ0FBQyxPQUFPO0FBQUEsWUFDZixPQUFPLENBQUMsT0FBTztBQUFBLFlBQ2YsYUFBYSxDQUFDLFdBQVc7QUFBQSxZQUN6QixJQUFJLENBQUMsSUFBSTtBQUFBLFlBQ1QsaUJBQWlCLENBQUMsZUFBZTtBQUFBLFVBQ25DO0FBQUEsUUFDRjtBQUFBLE1BQ0Y7QUFBQSxNQUNBLGNBQWMsRUFBRSxJQUFJLHVCQUF1QjtBQUFBLElBQzdDO0FBQUEsSUFDQSxjQUFjO0FBQUEsTUFDWixTQUFTLENBQUMsT0FBTyxjQUFjLFNBQVMsWUFBWSxTQUFTLFNBQVMsYUFBYSxJQUFJO0FBQUEsSUFDekY7QUFBQSxFQUNGO0FBQ0YsQ0FBQzsiLAogICJuYW1lcyI6IFsibG9hZEVudiIsICJsb2FkRW52Il0KfQo=
