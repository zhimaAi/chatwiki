// vite.config.js
import { fileURLToPath, URL } from "node:url";
import path from "path";
import { defineConfig, loadEnv as loadEnv2 } from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-pc/node_modules/vite/dist/node/index.js";
import vue from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-pc/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-pc/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import VueI18nPlugin from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-pc/node_modules/@intlify/unplugin-vue-i18n/lib/vite.mjs";
import { visualizer } from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-pc/node_modules/rollup-plugin-visualizer/dist/plugin/index.js";
import externalGlobals from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-pc/node_modules/rollup-plugin-external-globals/index.js";
import { createSvgIconsPlugin } from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-pc/node_modules/vite-plugin-svg-icons/dist/index.mjs";

// proxy_config.ts
import { loadEnv } from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-pc/node_modules/vite/dist/node/index.js";
var getProxyConfig = (opt) => {
  const { mode } = opt;
  const env = loadEnv(mode, process.cwd(), "");
  const proxyApis = ["/static", "/common", "/manage", "/app", "/chat", "/upload"];
  const proxy = {};
  console.log("PROXY_BASE_API_URL: " + env.PROXY_BASE_API_URL);
  proxyApis.forEach((key) => {
    proxy[key] = {
      target: env.PROXY_BASE_API_URL,
      changeOrigin: true
    };
  });
  return proxy;
};

// vite.config.js
var __vite_injected_original_dirname = "D:\\project\\zhima_chat_ai\\front-end\\chat-ai-pc";
var __vite_injected_original_import_meta_url = "file:///D:/project/zhima_chat_ai/front-end/chat-ai-pc/vite.config.js";
var vite_config_default = defineConfig(function(opt) {
  var command = opt.command, mode = opt.mode;
  var env = loadEnv2(mode, process.cwd(), "");
  var base = command === "serve" ? "/" : "./";
  var globals = externalGlobals({});
  return {
    plugins: [
      vue(),
      vueJsx(),
      createSvgIconsPlugin({
        // 指定需要缓存的图标文件夹
        iconDirs: [path.resolve(process.cwd(), "src/assets/icons")],
        // 指定symbolId格式
        symbolId: "[name]"
      }),
      VueI18nPlugin({
        runtimeOnly: true,
        compositionOnly: true,
        include: [path.resolve(__vite_injected_original_dirname, "./src/locales/lang**")]
      })
    ],
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", __vite_injected_original_import_meta_url))
      }
    },
    base,
    server: {
      proxy: getProxyConfig(opt)
    },
    build: {
      outDir: fileURLToPath(new URL("../../static/chat-ai-pc/web", __vite_injected_original_import_meta_url)),
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
          entryFileNames: "assets/js/[name]-[hash].js",
          assetFileNames: function(assetInfo) {
            var fiel_name = assetInfo.name.toLowerCase();
            if (fiel_name.endsWith(".css")) {
              return "assets/css/[name]-[hash].[ext]";
            }
            if (["png", "jpg", "jpeg", "svg"].some(function(ext) {
              return fiel_name.endsWith(ext);
            })) {
              return "assets/img/[name]-[hash].[ext]";
            }
            if (["ttf", "woff", "woff2"].some(function(ext) {
              return fiel_name.endsWith(ext);
            })) {
              return "assets/fonts/[name]-[hash].[ext]";
            }
            return "assets/[ext]/[name]-[hash].[ext]";
          },
          // 该选项允许你创建自定义的公共 chunk
          manualChunks: {
            "vue-chunks": ["vue", "vue-router", "pinia", "vue-i18n"],
            dayjs: ["dayjs"],
            axios: ["axios"],
            "crypto-js": ["crypto-js"],
            qs: ["qs"]
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiLCAicHJveHlfY29uZmlnLnRzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiRDpcXFxccHJvamVjdFxcXFx6aGltYV9jaGF0X2FpXFxcXGZyb250LWVuZFxcXFxjaGF0LWFpLXBjXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCJEOlxcXFxwcm9qZWN0XFxcXHpoaW1hX2NoYXRfYWlcXFxcZnJvbnQtZW5kXFxcXGNoYXQtYWktcGNcXFxcdml0ZS5jb25maWcuanNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0Q6L3Byb2plY3QvemhpbWFfY2hhdF9haS9mcm9udC1lbmQvY2hhdC1haS1wYy92aXRlLmNvbmZpZy5qc1wiO2ltcG9ydCB7IGZpbGVVUkxUb1BhdGgsIFVSTCB9IGZyb20gJ25vZGU6dXJsJztcclxuaW1wb3J0IHBhdGggZnJvbSAncGF0aCc7XHJcbmltcG9ydCB7IGRlZmluZUNvbmZpZywgbG9hZEVudiB9IGZyb20gJ3ZpdGUnO1xyXG5pbXBvcnQgdnVlIGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZSc7XHJcbmltcG9ydCB2dWVKc3ggZnJvbSAnQHZpdGVqcy9wbHVnaW4tdnVlLWpzeCc7XHJcbmltcG9ydCBWdWVJMThuUGx1Z2luIGZyb20gJ0BpbnRsaWZ5L3VucGx1Z2luLXZ1ZS1pMThuL3ZpdGUnO1xyXG5pbXBvcnQgeyB2aXN1YWxpemVyIH0gZnJvbSAncm9sbHVwLXBsdWdpbi12aXN1YWxpemVyJztcclxuaW1wb3J0IGV4dGVybmFsR2xvYmFscyBmcm9tICdyb2xsdXAtcGx1Z2luLWV4dGVybmFsLWdsb2JhbHMnO1xyXG5pbXBvcnQgeyBjcmVhdGVTdmdJY29uc1BsdWdpbiB9IGZyb20gJ3ZpdGUtcGx1Z2luLXN2Zy1pY29ucyc7XHJcbmltcG9ydCB7IGdldFByb3h5Q29uZmlnIH0gZnJvbSAnLi9wcm94eV9jb25maWcnO1xyXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoZnVuY3Rpb24gKG9wdCkge1xyXG4gICAgdmFyIGNvbW1hbmQgPSBvcHQuY29tbWFuZCwgbW9kZSA9IG9wdC5tb2RlO1xyXG4gICAgdmFyIGVudiA9IGxvYWRFbnYobW9kZSwgcHJvY2Vzcy5jd2QoKSwgJycpO1xyXG4gICAgdmFyIGJhc2UgPSBjb21tYW5kID09PSAnc2VydmUnID8gJy8nIDogJy4vJztcclxuICAgIHZhciBnbG9iYWxzID0gZXh0ZXJuYWxHbG9iYWxzKHt9KTtcclxuICAgIHJldHVybiB7XHJcbiAgICAgICAgcGx1Z2luczogW1xyXG4gICAgICAgICAgICB2dWUoKSxcclxuICAgICAgICAgICAgdnVlSnN4KCksXHJcbiAgICAgICAgICAgIGNyZWF0ZVN2Z0ljb25zUGx1Z2luKHtcclxuICAgICAgICAgICAgICAgIC8vIFx1NjMwN1x1NUI5QVx1OTcwMFx1ODk4MVx1N0YxM1x1NUI1OFx1NzY4NFx1NTZGRVx1NjgwN1x1NjU4N1x1NEVGNlx1NTkzOVxyXG4gICAgICAgICAgICAgICAgaWNvbkRpcnM6IFtwYXRoLnJlc29sdmUocHJvY2Vzcy5jd2QoKSwgJ3NyYy9hc3NldHMvaWNvbnMnKV0sXHJcbiAgICAgICAgICAgICAgICAvLyBcdTYzMDdcdTVCOUFzeW1ib2xJZFx1NjgzQ1x1NUYwRlxyXG4gICAgICAgICAgICAgICAgc3ltYm9sSWQ6ICdbbmFtZV0nXHJcbiAgICAgICAgICAgIH0pLFxyXG4gICAgICAgICAgICBWdWVJMThuUGx1Z2luKHtcclxuICAgICAgICAgICAgICAgIHJ1bnRpbWVPbmx5OiB0cnVlLFxyXG4gICAgICAgICAgICAgICAgY29tcG9zaXRpb25Pbmx5OiB0cnVlLFxyXG4gICAgICAgICAgICAgICAgaW5jbHVkZTogW3BhdGgucmVzb2x2ZShfX2Rpcm5hbWUsICcuL3NyYy9sb2NhbGVzL2xhbmcqKicpXVxyXG4gICAgICAgICAgICB9KVxyXG4gICAgICAgIF0sXHJcbiAgICAgICAgcmVzb2x2ZToge1xyXG4gICAgICAgICAgICBhbGlhczoge1xyXG4gICAgICAgICAgICAgICAgJ0AnOiBmaWxlVVJMVG9QYXRoKG5ldyBVUkwoJy4vc3JjJywgaW1wb3J0Lm1ldGEudXJsKSlcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH0sXHJcbiAgICAgICAgYmFzZTogYmFzZSxcclxuICAgICAgICBzZXJ2ZXI6IHtcclxuICAgICAgICAgICAgcHJveHk6IGdldFByb3h5Q29uZmlnKG9wdClcclxuICAgICAgICB9LFxyXG4gICAgICAgIGJ1aWxkOiB7XHJcbiAgICAgICAgICAgIG91dERpcjogZmlsZVVSTFRvUGF0aChuZXcgVVJMKCcuLi8uLi9zdGF0aWMvY2hhdC1haS1wYy93ZWInLCBpbXBvcnQubWV0YS51cmwpKSxcclxuICAgICAgICAgICAgZW1wdHlPdXREaXI6IHRydWUsXHJcbiAgICAgICAgICAgIGFzc2V0c0RpcjogJ2Fzc2V0cycsXHJcbiAgICAgICAgICAgIHNvdXJjZW1hcDogZW52LlZJVEVfU09VUkNFTUFQID09PSAndHJ1ZScsXHJcbiAgICAgICAgICAgIHJvbGx1cE9wdGlvbnM6IHtcclxuICAgICAgICAgICAgICAgIC8vIGV4dGVybmFsOiBbJ21vbWVudCcsICd2aWRlby5qcycsICdqc3BkZicsICd4bHN4J10sXHJcbiAgICAgICAgICAgICAgICBleHRlcm5hbDogW10sXHJcbiAgICAgICAgICAgICAgICBwbHVnaW5zOiBbZ2xvYmFscywgZW52LlZJVEVfVVNFX0JVTkRMRV9BTkFMWVpFUiA9PT0gJ3RydWUnID8gdmlzdWFsaXplcigpIDogdW5kZWZpbmVkXSxcclxuICAgICAgICAgICAgICAgIG91dHB1dDoge1xyXG4gICAgICAgICAgICAgICAgICAgIC8vIFx1ODFFQVx1NUI5QVx1NEU0OWNodW5rRmlsZU5hbWVcdTc1MUZcdTYyMTBcdTg5QzRcdTUyMTlcclxuICAgICAgICAgICAgICAgICAgICBjaHVua0ZpbGVOYW1lczogJ2Fzc2V0cy9qcy9bbmFtZV0tW2hhc2hdLmpzJyxcclxuICAgICAgICAgICAgICAgICAgICBlbnRyeUZpbGVOYW1lczogJ2Fzc2V0cy9qcy9bbmFtZV0tW2hhc2hdLmpzJyxcclxuICAgICAgICAgICAgICAgICAgICBhc3NldEZpbGVOYW1lczogZnVuY3Rpb24gKGFzc2V0SW5mbykge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB2YXIgZmllbF9uYW1lID0gYXNzZXRJbmZvLm5hbWUudG9Mb3dlckNhc2UoKTtcclxuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKGZpZWxfbmFtZS5lbmRzV2l0aCgnLmNzcycpKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9jc3MvW25hbWVdLVtoYXNoXS5bZXh0XSc7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKFsncG5nJywgJ2pwZycsICdqcGVnJywgJ3N2ZyddLnNvbWUoZnVuY3Rpb24gKGV4dCkgeyByZXR1cm4gZmllbF9uYW1lLmVuZHNXaXRoKGV4dCk7IH0pKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9pbWcvW25hbWVdLVtoYXNoXS5bZXh0XSc7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKFsndHRmJywgJ3dvZmYnLCAnd29mZjInXS5zb21lKGZ1bmN0aW9uIChleHQpIHsgcmV0dXJuIGZpZWxfbmFtZS5lbmRzV2l0aChleHQpOyB9KSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvZm9udHMvW25hbWVdLVtoYXNoXS5bZXh0XSc7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvW2V4dF0vW25hbWVdLVtoYXNoXS5bZXh0XSc7XHJcbiAgICAgICAgICAgICAgICAgICAgfSxcclxuICAgICAgICAgICAgICAgICAgICAvLyBcdThCRTVcdTkwMDlcdTk4NzlcdTUxNDFcdThCQjhcdTRGNjBcdTUyMUJcdTVFRkFcdTgxRUFcdTVCOUFcdTRFNDlcdTc2ODRcdTUxNkNcdTUxNzEgY2h1bmtcclxuICAgICAgICAgICAgICAgICAgICBtYW51YWxDaHVua3M6IHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgJ3Z1ZS1jaHVua3MnOiBbJ3Z1ZScsICd2dWUtcm91dGVyJywgJ3BpbmlhJywgJ3Z1ZS1pMThuJ10sXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGRheWpzOiBbJ2RheWpzJ10sXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGF4aW9zOiBbJ2F4aW9zJ10sXHJcbiAgICAgICAgICAgICAgICAgICAgICAgICdjcnlwdG8tanMnOiBbJ2NyeXB0by1qcyddLFxyXG4gICAgICAgICAgICAgICAgICAgICAgICBxczogWydxcyddLFxyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfSxcclxuICAgICAgICAgICAgY3NzQ29kZVNwbGl0OiAhKGVudi5WSVRFX1VTRV9DU1NfU1BMSVQgPT09ICdmYWxzZScpXHJcbiAgICAgICAgfSxcclxuICAgICAgICBvcHRpbWl6ZURlcHM6IHtcclxuICAgICAgICAgICAgaW5jbHVkZTogWyd2dWUnLCAndnVlLXJvdXRlcicsICdwaW5pYScsICd2dWUtaTE4bicsICdkYXlqcycsICdheGlvcycsICdjcnlwdG8tanMnLCAncXMnXVxyXG4gICAgICAgIH1cclxuICAgIH07XHJcbn0pO1xyXG4iLCAiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIkQ6XFxcXHByb2plY3RcXFxcemhpbWFfY2hhdF9haVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1wY1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiRDpcXFxccHJvamVjdFxcXFx6aGltYV9jaGF0X2FpXFxcXGZyb250LWVuZFxcXFxjaGF0LWFpLXBjXFxcXHByb3h5X2NvbmZpZy50c1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vRDovcHJvamVjdC96aGltYV9jaGF0X2FpL2Zyb250LWVuZC9jaGF0LWFpLXBjL3Byb3h5X2NvbmZpZy50c1wiO2ltcG9ydCB7IGxvYWRFbnYgfSBmcm9tICd2aXRlJztcclxuXHJcbi8vIFx1NTA0N1x1OEJCRVx1NzNBRlx1NTg4M1x1NTNEOFx1OTFDRlx1NTQ4Q1x1OTE0RFx1N0Y2RVx1NjU4N1x1NEVGNlx1NjgzQ1x1NUYwRlx1NURGMlx1N0VDRlx1NUI5QVx1NEU0OVx1NTk3RFxyXG5pbnRlcmZhY2UgUHJveHlDb25maWcge1xyXG4gIFtrZXk6IHN0cmluZ106IHtcclxuICAgIHRhcmdldDogc3RyaW5nO1xyXG4gICAgY2hhbmdlT3JpZ2luOiBib29sZWFuO1xyXG4gIH07XHJcbn1cclxuXHJcbmV4cG9ydCBjb25zdCBnZXRQcm94eUNvbmZpZyA9IChvcHQ6IHsgbW9kZTogc3RyaW5nIH0pOiBQcm94eUNvbmZpZyA9PiB7XHJcbiAgY29uc3QgeyBtb2RlIH0gPSBvcHQ7XHJcbiAgICBjb25zdCBlbnYgPSBsb2FkRW52KG1vZGUsIHByb2Nlc3MuY3dkKCksICcnKTtcclxuXHJcbiAgICBjb25zdCBwcm94eUFwaXMgPSBbJy9zdGF0aWMnLCAnL2NvbW1vbicsICcvbWFuYWdlJywgJy9hcHAnLCAnL2NoYXQnLCAnL3VwbG9hZCddO1xyXG4gICAgY29uc3QgcHJveHk6IFByb3h5Q29uZmlnID0ge307XHJcblxyXG4gICAgY29uc29sZS5sb2coJ1BST1hZX0JBU0VfQVBJX1VSTDogJyArIGVudi5QUk9YWV9CQVNFX0FQSV9VUkwpO1xyXG5cclxuICAgIHByb3h5QXBpcy5mb3JFYWNoKChrZXkpID0+IHtcclxuICAgICAgcHJveHlba2V5XSA9IHtcclxuICAgICAgICB0YXJnZXQ6IGVudi5QUk9YWV9CQVNFX0FQSV9VUkwsXHJcbiAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlLFxyXG4gICAgICB9O1xyXG4gICAgfSk7XHJcblxyXG4gICAgcmV0dXJuIHByb3h5O1xyXG59O1xyXG4iXSwKICAibWFwcGluZ3MiOiAiO0FBQW1VLFNBQVMsZUFBZSxXQUFXO0FBQ3RXLE9BQU8sVUFBVTtBQUNqQixTQUFTLGNBQWMsV0FBQUEsZ0JBQWU7QUFDdEMsT0FBTyxTQUFTO0FBQ2hCLE9BQU8sWUFBWTtBQUNuQixPQUFPLG1CQUFtQjtBQUMxQixTQUFTLGtCQUFrQjtBQUMzQixPQUFPLHFCQUFxQjtBQUM1QixTQUFTLDRCQUE0Qjs7O0FDUmdTLFNBQVMsZUFBZTtBQVV0VixJQUFNLGlCQUFpQixDQUFDLFFBQXVDO0FBQ3BFLFFBQU0sRUFBRSxLQUFLLElBQUk7QUFDZixRQUFNLE1BQU0sUUFBUSxNQUFNLFFBQVEsSUFBSSxHQUFHLEVBQUU7QUFFM0MsUUFBTSxZQUFZLENBQUMsV0FBVyxXQUFXLFdBQVcsUUFBUSxTQUFTLFNBQVM7QUFDOUUsUUFBTSxRQUFxQixDQUFDO0FBRTVCLFVBQVEsSUFBSSx5QkFBeUIsSUFBSSxrQkFBa0I7QUFFM0QsWUFBVSxRQUFRLENBQUMsUUFBUTtBQUN6QixVQUFNLEdBQUcsSUFBSTtBQUFBLE1BQ1gsUUFBUSxJQUFJO0FBQUEsTUFDWixjQUFjO0FBQUEsSUFDaEI7QUFBQSxFQUNGLENBQUM7QUFFRCxTQUFPO0FBQ1g7OztBRDNCQSxJQUFNLG1DQUFtQztBQUFrSyxJQUFNLDJDQUEyQztBQVU1UCxJQUFPLHNCQUFRLGFBQWEsU0FBVSxLQUFLO0FBQ3ZDLE1BQUksVUFBVSxJQUFJLFNBQVMsT0FBTyxJQUFJO0FBQ3RDLE1BQUksTUFBTUMsU0FBUSxNQUFNLFFBQVEsSUFBSSxHQUFHLEVBQUU7QUFDekMsTUFBSSxPQUFPLFlBQVksVUFBVSxNQUFNO0FBQ3ZDLE1BQUksVUFBVSxn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
