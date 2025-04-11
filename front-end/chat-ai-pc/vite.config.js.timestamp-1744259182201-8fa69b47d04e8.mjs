// vite.config.js
import { fileURLToPath, URL } from "node:url";
import path from "path";
import { defineConfig, loadEnv as loadEnv2 } from "file:///F:/zhima_2022/zhima_chatwiki/front-end/chat-ai-pc/node_modules/vite/dist/node/index.js";
import vue from "file:///F:/zhima_2022/zhima_chatwiki/front-end/chat-ai-pc/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///F:/zhima_2022/zhima_chatwiki/front-end/chat-ai-pc/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import VueI18nPlugin from "file:///F:/zhima_2022/zhima_chatwiki/front-end/chat-ai-pc/node_modules/@intlify/unplugin-vue-i18n/lib/vite.mjs";
import { visualizer } from "file:///F:/zhima_2022/zhima_chatwiki/front-end/chat-ai-pc/node_modules/rollup-plugin-visualizer/dist/plugin/index.js";
import externalGlobals from "file:///F:/zhima_2022/zhima_chatwiki/front-end/chat-ai-pc/node_modules/rollup-plugin-external-globals/index.js";
import { createSvgIconsPlugin } from "file:///F:/zhima_2022/zhima_chatwiki/front-end/chat-ai-pc/node_modules/vite-plugin-svg-icons/dist/index.mjs";

// proxy_config.ts
import { loadEnv } from "file:///F:/zhima_2022/zhima_chatwiki/front-end/chat-ai-pc/node_modules/vite/dist/node/index.js";
var getProxyConfig = (opt) => {
  const { mode } = opt;
  const env = loadEnv(mode, process.cwd(), "");
  const proxyApis = ["/static", "/common", "/manage", "/app", "/chat", "/upload", "/public"];
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
var __vite_injected_original_dirname = "F:\\zhima_2022\\zhima_chatwiki\\front-end\\chat-ai-pc";
var __vite_injected_original_import_meta_url = "file:///F:/zhima_2022/zhima_chatwiki/front-end/chat-ai-pc/vite.config.js";
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiLCAicHJveHlfY29uZmlnLnRzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiRjpcXFxcemhpbWFfMjAyMlxcXFx6aGltYV9jaGF0d2lraVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1wY1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiRjpcXFxcemhpbWFfMjAyMlxcXFx6aGltYV9jaGF0d2lraVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1wY1xcXFx2aXRlLmNvbmZpZy5qc1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vRjovemhpbWFfMjAyMi96aGltYV9jaGF0d2lraS9mcm9udC1lbmQvY2hhdC1haS1wYy92aXRlLmNvbmZpZy5qc1wiO2ltcG9ydCB7IGZpbGVVUkxUb1BhdGgsIFVSTCB9IGZyb20gJ25vZGU6dXJsJztcbmltcG9ydCBwYXRoIGZyb20gJ3BhdGgnO1xuaW1wb3J0IHsgZGVmaW5lQ29uZmlnLCBsb2FkRW52IH0gZnJvbSAndml0ZSc7XG5pbXBvcnQgdnVlIGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZSc7XG5pbXBvcnQgdnVlSnN4IGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZS1qc3gnO1xuaW1wb3J0IFZ1ZUkxOG5QbHVnaW4gZnJvbSAnQGludGxpZnkvdW5wbHVnaW4tdnVlLWkxOG4vdml0ZSc7XG5pbXBvcnQgeyB2aXN1YWxpemVyIH0gZnJvbSAncm9sbHVwLXBsdWdpbi12aXN1YWxpemVyJztcbmltcG9ydCBleHRlcm5hbEdsb2JhbHMgZnJvbSAncm9sbHVwLXBsdWdpbi1leHRlcm5hbC1nbG9iYWxzJztcbmltcG9ydCB7IGNyZWF0ZVN2Z0ljb25zUGx1Z2luIH0gZnJvbSAndml0ZS1wbHVnaW4tc3ZnLWljb25zJztcbmltcG9ydCB7IGdldFByb3h5Q29uZmlnIH0gZnJvbSAnLi9wcm94eV9jb25maWcnO1xuZXhwb3J0IGRlZmF1bHQgZGVmaW5lQ29uZmlnKGZ1bmN0aW9uIChvcHQpIHtcbiAgICB2YXIgY29tbWFuZCA9IG9wdC5jb21tYW5kLCBtb2RlID0gb3B0Lm1vZGU7XG4gICAgdmFyIGVudiA9IGxvYWRFbnYobW9kZSwgcHJvY2Vzcy5jd2QoKSwgJycpO1xuICAgIHZhciBiYXNlID0gY29tbWFuZCA9PT0gJ3NlcnZlJyA/ICcvJyA6ICcuLyc7XG4gICAgdmFyIGdsb2JhbHMgPSBleHRlcm5hbEdsb2JhbHMoe30pO1xuICAgIHJldHVybiB7XG4gICAgICAgIHBsdWdpbnM6IFtcbiAgICAgICAgICAgIHZ1ZSgpLFxuICAgICAgICAgICAgdnVlSnN4KCksXG4gICAgICAgICAgICBjcmVhdGVTdmdJY29uc1BsdWdpbih7XG4gICAgICAgICAgICAgICAgLy8gXHU2MzA3XHU1QjlBXHU5NzAwXHU4OTgxXHU3RjEzXHU1QjU4XHU3Njg0XHU1NkZFXHU2ODA3XHU2NTg3XHU0RUY2XHU1OTM5XG4gICAgICAgICAgICAgICAgaWNvbkRpcnM6IFtwYXRoLnJlc29sdmUocHJvY2Vzcy5jd2QoKSwgJ3NyYy9hc3NldHMvaWNvbnMnKV0sXG4gICAgICAgICAgICAgICAgLy8gXHU2MzA3XHU1QjlBc3ltYm9sSWRcdTY4M0NcdTVGMEZcbiAgICAgICAgICAgICAgICBzeW1ib2xJZDogJ1tuYW1lXSdcbiAgICAgICAgICAgIH0pLFxuICAgICAgICAgICAgVnVlSTE4blBsdWdpbih7XG4gICAgICAgICAgICAgICAgcnVudGltZU9ubHk6IHRydWUsXG4gICAgICAgICAgICAgICAgY29tcG9zaXRpb25Pbmx5OiB0cnVlLFxuICAgICAgICAgICAgICAgIGluY2x1ZGU6IFtwYXRoLnJlc29sdmUoX19kaXJuYW1lLCAnLi9zcmMvbG9jYWxlcy9sYW5nKionKV1cbiAgICAgICAgICAgIH0pXG4gICAgICAgIF0sXG4gICAgICAgIHJlc29sdmU6IHtcbiAgICAgICAgICAgIGFsaWFzOiB7XG4gICAgICAgICAgICAgICAgJ0AnOiBmaWxlVVJMVG9QYXRoKG5ldyBVUkwoJy4vc3JjJywgaW1wb3J0Lm1ldGEudXJsKSlcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSxcbiAgICAgICAgYmFzZTogYmFzZSxcbiAgICAgICAgc2VydmVyOiB7XG4gICAgICAgICAgICBwcm94eTogZ2V0UHJveHlDb25maWcob3B0KVxuICAgICAgICB9LFxuICAgICAgICBidWlsZDoge1xuICAgICAgICAgICAgb3V0RGlyOiBmaWxlVVJMVG9QYXRoKG5ldyBVUkwoJy4uLy4uL3N0YXRpYy9jaGF0LWFpLXBjL3dlYicsIGltcG9ydC5tZXRhLnVybCkpLFxuICAgICAgICAgICAgZW1wdHlPdXREaXI6IHRydWUsXG4gICAgICAgICAgICBhc3NldHNEaXI6ICdhc3NldHMnLFxuICAgICAgICAgICAgc291cmNlbWFwOiBlbnYuVklURV9TT1VSQ0VNQVAgPT09ICd0cnVlJyxcbiAgICAgICAgICAgIHJvbGx1cE9wdGlvbnM6IHtcbiAgICAgICAgICAgICAgICAvLyBleHRlcm5hbDogWydtb21lbnQnLCAndmlkZW8uanMnLCAnanNwZGYnLCAneGxzeCddLFxuICAgICAgICAgICAgICAgIGV4dGVybmFsOiBbXSxcbiAgICAgICAgICAgICAgICBwbHVnaW5zOiBbZ2xvYmFscywgZW52LlZJVEVfVVNFX0JVTkRMRV9BTkFMWVpFUiA9PT0gJ3RydWUnID8gdmlzdWFsaXplcigpIDogdW5kZWZpbmVkXSxcbiAgICAgICAgICAgICAgICBvdXRwdXQ6IHtcbiAgICAgICAgICAgICAgICAgICAgLy8gXHU4MUVBXHU1QjlBXHU0RTQ5Y2h1bmtGaWxlTmFtZVx1NzUxRlx1NjIxMFx1ODlDNFx1NTIxOVxuICAgICAgICAgICAgICAgICAgICBjaHVua0ZpbGVOYW1lczogJ2Fzc2V0cy9qcy9bbmFtZV0tW2hhc2hdLmpzJyxcbiAgICAgICAgICAgICAgICAgICAgZW50cnlGaWxlTmFtZXM6ICdhc3NldHMvanMvW25hbWVdLVtoYXNoXS5qcycsXG4gICAgICAgICAgICAgICAgICAgIGFzc2V0RmlsZU5hbWVzOiBmdW5jdGlvbiAoYXNzZXRJbmZvKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICB2YXIgZmllbF9uYW1lID0gYXNzZXRJbmZvLm5hbWUudG9Mb3dlckNhc2UoKTtcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChmaWVsX25hbWUuZW5kc1dpdGgoJy5jc3MnKSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiAnYXNzZXRzL2Nzcy9bbmFtZV0tW2hhc2hdLltleHRdJztcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChbJ3BuZycsICdqcGcnLCAnanBlZycsICdzdmcnXS5zb21lKGZ1bmN0aW9uIChleHQpIHsgcmV0dXJuIGZpZWxfbmFtZS5lbmRzV2l0aChleHQpOyB9KSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiAnYXNzZXRzL2ltZy9bbmFtZV0tW2hhc2hdLltleHRdJztcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChbJ3R0ZicsICd3b2ZmJywgJ3dvZmYyJ10uc29tZShmdW5jdGlvbiAoZXh0KSB7IHJldHVybiBmaWVsX25hbWUuZW5kc1dpdGgoZXh0KTsgfSkpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9mb250cy9bbmFtZV0tW2hhc2hdLltleHRdJztcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiAnYXNzZXRzL1tleHRdL1tuYW1lXS1baGFzaF0uW2V4dF0nO1xuICAgICAgICAgICAgICAgICAgICB9LFxuICAgICAgICAgICAgICAgICAgICAvLyBcdThCRTVcdTkwMDlcdTk4NzlcdTUxNDFcdThCQjhcdTRGNjBcdTUyMUJcdTVFRkFcdTgxRUFcdTVCOUFcdTRFNDlcdTc2ODRcdTUxNkNcdTUxNzEgY2h1bmtcbiAgICAgICAgICAgICAgICAgICAgbWFudWFsQ2h1bmtzOiB7XG4gICAgICAgICAgICAgICAgICAgICAgICAndnVlLWNodW5rcyc6IFsndnVlJywgJ3Z1ZS1yb3V0ZXInLCAncGluaWEnLCAndnVlLWkxOG4nXSxcbiAgICAgICAgICAgICAgICAgICAgICAgIGRheWpzOiBbJ2RheWpzJ10sXG4gICAgICAgICAgICAgICAgICAgICAgICBheGlvczogWydheGlvcyddLFxuICAgICAgICAgICAgICAgICAgICAgICAgJ2NyeXB0by1qcyc6IFsnY3J5cHRvLWpzJ10sXG4gICAgICAgICAgICAgICAgICAgICAgICBxczogWydxcyddLFxuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSxcbiAgICAgICAgICAgIGNzc0NvZGVTcGxpdDogIShlbnYuVklURV9VU0VfQ1NTX1NQTElUID09PSAnZmFsc2UnKVxuICAgICAgICB9LFxuICAgICAgICBvcHRpbWl6ZURlcHM6IHtcbiAgICAgICAgICAgIGluY2x1ZGU6IFsndnVlJywgJ3Z1ZS1yb3V0ZXInLCAncGluaWEnLCAndnVlLWkxOG4nLCAnZGF5anMnLCAnYXhpb3MnLCAnY3J5cHRvLWpzJywgJ3FzJ11cbiAgICAgICAgfVxuICAgIH07XG59KTtcbiIsICJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiRjpcXFxcemhpbWFfMjAyMlxcXFx6aGltYV9jaGF0d2lraVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1wY1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiRjpcXFxcemhpbWFfMjAyMlxcXFx6aGltYV9jaGF0d2lraVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1wY1xcXFxwcm94eV9jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0Y6L3poaW1hXzIwMjIvemhpbWFfY2hhdHdpa2kvZnJvbnQtZW5kL2NoYXQtYWktcGMvcHJveHlfY29uZmlnLnRzXCI7aW1wb3J0IHsgbG9hZEVudiB9IGZyb20gJ3ZpdGUnO1xyXG5cclxuLy8gXHU1MDQ3XHU4QkJFXHU3M0FGXHU1ODgzXHU1M0Q4XHU5MUNGXHU1NDhDXHU5MTREXHU3RjZFXHU2NTg3XHU0RUY2XHU2ODNDXHU1RjBGXHU1REYyXHU3RUNGXHU1QjlBXHU0RTQ5XHU1OTdEXHJcbmludGVyZmFjZSBQcm94eUNvbmZpZyB7XHJcbiAgW2tleTogc3RyaW5nXToge1xyXG4gICAgdGFyZ2V0OiBzdHJpbmc7XHJcbiAgICBjaGFuZ2VPcmlnaW46IGJvb2xlYW47XHJcbiAgfTtcclxufVxyXG5cclxuZXhwb3J0IGNvbnN0IGdldFByb3h5Q29uZmlnID0gKG9wdDogeyBtb2RlOiBzdHJpbmcgfSk6IFByb3h5Q29uZmlnID0+IHtcclxuICBjb25zdCB7IG1vZGUgfSA9IG9wdDtcclxuICAgIGNvbnN0IGVudiA9IGxvYWRFbnYobW9kZSwgcHJvY2Vzcy5jd2QoKSwgJycpO1xyXG5cclxuICAgIGNvbnN0IHByb3h5QXBpcyA9IFsnL3N0YXRpYycsICcvY29tbW9uJywgJy9tYW5hZ2UnLCAnL2FwcCcsICcvY2hhdCcsICcvdXBsb2FkJywgJy9wdWJsaWMnXTtcclxuICAgIGNvbnN0IHByb3h5OiBQcm94eUNvbmZpZyA9IHt9O1xyXG5cclxuICAgIGNvbnNvbGUubG9nKCdQUk9YWV9CQVNFX0FQSV9VUkw6ICcgKyBlbnYuUFJPWFlfQkFTRV9BUElfVVJMKTtcclxuXHJcbiAgICBwcm94eUFwaXMuZm9yRWFjaCgoa2V5KSA9PiB7XHJcbiAgICAgIHByb3h5W2tleV0gPSB7XHJcbiAgICAgICAgdGFyZ2V0OiBlbnYuUFJPWFlfQkFTRV9BUElfVVJMLFxyXG4gICAgICAgIGNoYW5nZU9yaWdpbjogdHJ1ZSxcclxuICAgICAgfTtcclxuICAgIH0pO1xyXG5c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
