// vite.config.js
import { fileURLToPath, URL } from "node:url";
import path from "path";
import { defineConfig, loadEnv as loadEnv2 } from "file:///D:/%E8%8A%9D%E9%BA%BBAI/zhima_chat_ai/front-end/chat-ai-mobile/node_modules/vite/dist/node/index.js";
import vue from "file:///D:/%E8%8A%9D%E9%BA%BBAI/zhima_chat_ai/front-end/chat-ai-mobile/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///D:/%E8%8A%9D%E9%BA%BBAI/zhima_chat_ai/front-end/chat-ai-mobile/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import VueI18nPlugin from "file:///D:/%E8%8A%9D%E9%BA%BBAI/zhima_chat_ai/front-end/chat-ai-mobile/node_modules/@intlify/unplugin-vue-i18n/lib/vite.mjs";
import { visualizer } from "file:///D:/%E8%8A%9D%E9%BA%BBAI/zhima_chat_ai/front-end/chat-ai-mobile/node_modules/rollup-plugin-visualizer/dist/plugin/index.js";
import externalGlobals from "file:///D:/%E8%8A%9D%E9%BA%BBAI/zhima_chat_ai/front-end/chat-ai-mobile/node_modules/rollup-plugin-external-globals/index.js";
import { createSvgIconsPlugin } from "file:///D:/%E8%8A%9D%E9%BA%BBAI/zhima_chat_ai/front-end/chat-ai-mobile/node_modules/vite-plugin-svg-icons/dist/index.mjs";

// proxy_config.ts
import { loadEnv } from "file:///D:/%E8%8A%9D%E9%BA%BBAI/zhima_chat_ai/front-end/chat-ai-mobile/node_modules/vite/dist/node/index.js";
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
var __vite_injected_original_dirname = "D:\\\u829D\u9EBBAI\\zhima_chat_ai\\front-end\\chat-ai-mobile";
var __vite_injected_original_import_meta_url = "file:///D:/%E8%8A%9D%E9%BA%BBAI/zhima_chat_ai/front-end/chat-ai-mobile/vite.config.js";
var vite_config_default = defineConfig(function(opt) {
  var command = opt.command, mode = opt.mode;
  var env = loadEnv2(mode, process.cwd(), "");
  var base = command === "serve" ? "/" : "/";
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
      outDir: fileURLToPath(new URL("../../static/chat-ai-mobile", __vite_injected_original_import_meta_url)),
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
            return "assets/[name]-[hash].[ext]";
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiLCAicHJveHlfY29uZmlnLnRzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiRDpcXFxcXHU4MjlEXHU5RUJCQUlcXFxcemhpbWFfY2hhdF9haVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1tb2JpbGVcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIkQ6XFxcXFx1ODI5RFx1OUVCQkFJXFxcXHpoaW1hX2NoYXRfYWlcXFxcZnJvbnQtZW5kXFxcXGNoYXQtYWktbW9iaWxlXFxcXHZpdGUuY29uZmlnLmpzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9EOi8lRTglOEElOUQlRTklQkElQkJBSS96aGltYV9jaGF0X2FpL2Zyb250LWVuZC9jaGF0LWFpLW1vYmlsZS92aXRlLmNvbmZpZy5qc1wiO2ltcG9ydCB7IGZpbGVVUkxUb1BhdGgsIFVSTCB9IGZyb20gJ25vZGU6dXJsJztcbmltcG9ydCBwYXRoIGZyb20gJ3BhdGgnO1xuaW1wb3J0IHsgZGVmaW5lQ29uZmlnLCBsb2FkRW52IH0gZnJvbSAndml0ZSc7XG5pbXBvcnQgdnVlIGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZSc7XG5pbXBvcnQgdnVlSnN4IGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZS1qc3gnO1xuaW1wb3J0IFZ1ZUkxOG5QbHVnaW4gZnJvbSAnQGludGxpZnkvdW5wbHVnaW4tdnVlLWkxOG4vdml0ZSc7XG5pbXBvcnQgeyB2aXN1YWxpemVyIH0gZnJvbSAncm9sbHVwLXBsdWdpbi12aXN1YWxpemVyJztcbmltcG9ydCBleHRlcm5hbEdsb2JhbHMgZnJvbSAncm9sbHVwLXBsdWdpbi1leHRlcm5hbC1nbG9iYWxzJztcbmltcG9ydCB7IGNyZWF0ZVN2Z0ljb25zUGx1Z2luIH0gZnJvbSAndml0ZS1wbHVnaW4tc3ZnLWljb25zJztcbmltcG9ydCB7IGdldFByb3h5Q29uZmlnIH0gZnJvbSAnLi9wcm94eV9jb25maWcnO1xuZXhwb3J0IGRlZmF1bHQgZGVmaW5lQ29uZmlnKGZ1bmN0aW9uIChvcHQpIHtcbiAgICB2YXIgY29tbWFuZCA9IG9wdC5jb21tYW5kLCBtb2RlID0gb3B0Lm1vZGU7XG4gICAgdmFyIGVudiA9IGxvYWRFbnYobW9kZSwgcHJvY2Vzcy5jd2QoKSwgJycpO1xuICAgIHZhciBiYXNlID0gY29tbWFuZCA9PT0gJ3NlcnZlJyA/ICcvJyA6ICcvJztcbiAgICB2YXIgZ2xvYmFscyA9IGV4dGVybmFsR2xvYmFscyh7fSk7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgcGx1Z2luczogW1xuICAgICAgICAgICAgdnVlKCksXG4gICAgICAgICAgICB2dWVKc3goKSxcbiAgICAgICAgICAgIGNyZWF0ZVN2Z0ljb25zUGx1Z2luKHtcbiAgICAgICAgICAgICAgICAvLyBcdTYzMDdcdTVCOUFcdTk3MDBcdTg5ODFcdTdGMTNcdTVCNThcdTc2ODRcdTU2RkVcdTY4MDdcdTY1ODdcdTRFRjZcdTU5MzlcbiAgICAgICAgICAgICAgICBpY29uRGlyczogW3BhdGgucmVzb2x2ZShwcm9jZXNzLmN3ZCgpLCAnc3JjL2Fzc2V0cy9pY29ucycpXSxcbiAgICAgICAgICAgICAgICAvLyBcdTYzMDdcdTVCOUFzeW1ib2xJZFx1NjgzQ1x1NUYwRlxuICAgICAgICAgICAgICAgIHN5bWJvbElkOiAnW25hbWVdJ1xuICAgICAgICAgICAgfSksXG4gICAgICAgICAgICBWdWVJMThuUGx1Z2luKHtcbiAgICAgICAgICAgICAgICBydW50aW1lT25seTogdHJ1ZSxcbiAgICAgICAgICAgICAgICBjb21wb3NpdGlvbk9ubHk6IHRydWUsXG4gICAgICAgICAgICAgICAgaW5jbHVkZTogW3BhdGgucmVzb2x2ZShfX2Rpcm5hbWUsICcuL3NyYy9sb2NhbGVzL2xhbmcqKicpXVxuICAgICAgICAgICAgfSlcbiAgICAgICAgXSxcbiAgICAgICAgcmVzb2x2ZToge1xuICAgICAgICAgICAgYWxpYXM6IHtcbiAgICAgICAgICAgICAgICAnQCc6IGZpbGVVUkxUb1BhdGgobmV3IFVSTCgnLi9zcmMnLCBpbXBvcnQubWV0YS51cmwpKVxuICAgICAgICAgICAgfVxuICAgICAgICB9LFxuICAgICAgICBiYXNlOiBiYXNlLFxuICAgICAgICBzZXJ2ZXI6IHtcbiAgICAgICAgICAgIHByb3h5OiBnZXRQcm94eUNvbmZpZyhvcHQpXG4gICAgICAgIH0sXG4gICAgICAgIGJ1aWxkOiB7XG4gICAgICAgICAgICBvdXREaXI6IGZpbGVVUkxUb1BhdGgobmV3IFVSTCgnLi4vLi4vc3RhdGljL2NoYXQtYWktbW9iaWxlJywgaW1wb3J0Lm1ldGEudXJsKSksXG4gICAgICAgICAgICBhc3NldHNEaXI6ICdhc3NldHMnLFxuICAgICAgICAgICAgc291cmNlbWFwOiBlbnYuVklURV9TT1VSQ0VNQVAgPT09ICd0cnVlJyxcbiAgICAgICAgICAgIHJvbGx1cE9wdGlvbnM6IHtcbiAgICAgICAgICAgICAgICAvLyBleHRlcm5hbDogWydtb21lbnQnLCAndmlkZW8uanMnLCAnanNwZGYnLCAneGxzeCddLFxuICAgICAgICAgICAgICAgIGV4dGVybmFsOiBbXSxcbiAgICAgICAgICAgICAgICBwbHVnaW5zOiBbZ2xvYmFscywgZW52LlZJVEVfVVNFX0JVTkRMRV9BTkFMWVpFUiA9PT0gJ3RydWUnID8gdmlzdWFsaXplcigpIDogdW5kZWZpbmVkXSxcbiAgICAgICAgICAgICAgICBvdXRwdXQ6IHtcbiAgICAgICAgICAgICAgICAgICAgLy8gXHU4MUVBXHU1QjlBXHU0RTQ5Y2h1bmtGaWxlTmFtZVx1NzUxRlx1NjIxMFx1ODlDNFx1NTIxOVxuICAgICAgICAgICAgICAgICAgICBjaHVua0ZpbGVOYW1lczogJ2Fzc2V0cy9qcy9bbmFtZV0tW2hhc2hdLmpzJyxcbiAgICAgICAgICAgICAgICAgICAgZW50cnlGaWxlTmFtZXM6ICdbbmFtZV0uanMnLFxuICAgICAgICAgICAgICAgICAgICBhc3NldEZpbGVOYW1lczogZnVuY3Rpb24gKGFzc2V0SW5mbykge1xuICAgICAgICAgICAgICAgICAgICAgICAgdmFyIGZpZWxfbmFtZSA9IGFzc2V0SW5mby5uYW1lLnRvTG93ZXJDYXNlKCk7XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoZmllbF9uYW1lLmVuZHNXaXRoKCcuY3NzJykpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9jc3MvW25hbWVdLVtoYXNoXS5bZXh0XSc7XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoWydwbmcnLCAnanBnJywgJ2pwZWcnLCAnc3ZnJ10uc29tZShmdW5jdGlvbiAoZXh0KSB7IHJldHVybiBmaWVsX25hbWUuZW5kc1dpdGgoZXh0KTsgfSkpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9pbWcvW25hbWVdLVtoYXNoXS5bZXh0XSc7XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoWyd0dGYnLCAnd29mZicsICd3b2ZmMiddLnNvbWUoZnVuY3Rpb24gKGV4dCkgeyByZXR1cm4gZmllbF9uYW1lLmVuZHNXaXRoKGV4dCk7IH0pKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvZm9udHMvW25hbWVdLVtoYXNoXS5bZXh0XSc7XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9bbmFtZV0tW2hhc2hdLltleHRdJztcbiAgICAgICAgICAgICAgICAgICAgfSxcbiAgICAgICAgICAgICAgICAgICAgLy8gXHU4QkU1XHU5MDA5XHU5ODc5XHU1MTQxXHU4QkI4XHU0RjYwXHU1MjFCXHU1RUZBXHU4MUVBXHU1QjlBXHU0RTQ5XHU3Njg0XHU1MTZDXHU1MTcxIGNodW5rXG4gICAgICAgICAgICAgICAgICAgIG1hbnVhbENodW5rczoge1xuICAgICAgICAgICAgICAgICAgICAgICAgJ3Z1ZS1jaHVua3MnOiBbJ3Z1ZScsICd2dWUtcm91dGVyJywgJ3BpbmlhJywgJ3Z1ZS1pMThuJ10sXG4gICAgICAgICAgICAgICAgICAgICAgICBkYXlqczogWydkYXlqcyddLFxuICAgICAgICAgICAgICAgICAgICAgICAgYXhpb3M6IFsnYXhpb3MnXSxcbiAgICAgICAgICAgICAgICAgICAgICAgICdjcnlwdG8tanMnOiBbJ2NyeXB0by1qcyddLFxuICAgICAgICAgICAgICAgICAgICAgICAgcXM6IFsncXMnXSxcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgICBjc3NDb2RlU3BsaXQ6ICEoZW52LlZJVEVfVVNFX0NTU19TUExJVCA9PT0gJ2ZhbHNlJylcbiAgICAgICAgfSxcbiAgICAgICAgb3B0aW1pemVEZXBzOiB7XG4gICAgICAgICAgICBpbmNsdWRlOiBbJ3Z1ZScsICd2dWUtcm91dGVyJywgJ3BpbmlhJywgJ3Z1ZS1pMThuJywgJ2RheWpzJywgJ2F4aW9zJywgJ2NyeXB0by1qcycsICdxcyddXG4gICAgICAgIH1cbiAgICB9O1xufSk7XG4iLCAiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIkQ6XFxcXFx1ODI5RFx1OUVCQkFJXFxcXHpoaW1hX2NoYXRfYWlcXFxcZnJvbnQtZW5kXFxcXGNoYXQtYWktbW9iaWxlXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCJEOlxcXFxcdTgyOURcdTlFQkJBSVxcXFx6aGltYV9jaGF0X2FpXFxcXGZyb250LWVuZFxcXFxjaGF0LWFpLW1vYmlsZVxcXFxwcm94eV9jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0Q6LyVFOCU4QSU5RCVFOSVCQSVCQkFJL3poaW1hX2NoYXRfYWkvZnJvbnQtZW5kL2NoYXQtYWktbW9iaWxlL3Byb3h5X2NvbmZpZy50c1wiO2ltcG9ydCB7IGxvYWRFbnYgfSBmcm9tICd2aXRlJztcblxuLy8gXHU1MDQ3XHU4QkJFXHU3M0FGXHU1ODgzXHU1M0Q4XHU5MUNGXHU1NDhDXHU5MTREXHU3RjZFXHU2NTg3XHU0RUY2XHU2ODNDXHU1RjBGXHU1REYyXHU3RUNGXHU1QjlBXHU0RTQ5XHU1OTdEXG5pbnRlcmZhY2UgUHJveHlDb25maWcge1xuICBba2V5OiBzdHJpbmddOiB7XG4gICAgdGFyZ2V0OiBzdHJpbmc7XG4gICAgY2hhbmdlT3JpZ2luOiBib29sZWFuO1xuICB9O1xufVxuXG5leHBvcnQgY29uc3QgZ2V0UHJveHlDb25maWcgPSAob3B0OiB7IG1vZGU6IHN0cmluZyB9KTogUHJveHlDb25maWcgPT4ge1xuICBjb25zdCB7IG1vZGUgfSA9IG9wdDtcbiAgICBjb25zdCBlbnYgPSBsb2FkRW52KG1vZGUsIHByb2Nlc3MuY3dkKCksICcnKTtcblxuICAgIGNvbnN0IHByb3h5QXBpcyA9IFsnL3N0YXRpYycsICcvY29tbW9uJywgJy9tYW5hZ2UnLCAnL2FwcCcsICcvY2hhdCcsICcvdXBsb2FkJ107XG4gICAgY29uc3QgcHJveHk6IFByb3h5Q29uZmlnID0ge307XG5cbiAgICBjb25zb2xlLmxvZygnUFJPWFlfQkFTRV9BUElfVVJMOiAnICsgZW52LlBST1hZX0JBU0VfQVBJX1VSTCk7XG5cbiAgICBwcm94eUFwaXMuZm9yRWFjaCgoa2V5KSA9PiB7XG4gICAgICBwcm94eVtrZXldID0ge1xuICAgICAgICB0YXJnZXQ6IGVudi5QUk9YWV9CQVNFX0FQSV9VUkwsXG4gICAgICAgIGNoYW5nZU9yaWdpbjogdHJ1ZSxcbiAgICAgIH07XG4gICAgfSk7XG5cbiAgICByZXR1cm4gcHJveHk7XG59O1xuIl0sCiAgIm1hcHBpbmdzIjogIjtBQUFzVixTQUFTLGVBQWUsV0FBVztBQUN6WCxPQUFPLFVBQVU7QUFDakIsU0FBUyxjQUFjLFdBQUFBLGdCQUFlO0FBQ3RDLE9BQU8sU0FBUztBQUNoQixPQUFPLFlBQVk7QUFDbkIsT0FBTyxtQkFBbUI7QUFDMUIsU0FBUyxrQkFBa0I7QUFDM0IsT0FBTyxxQkFBcUI7QUFDNUIsU0FBUyw0QkFBNEI7OztBQ1JtVCxTQUFTLGVBQWU7QUFVelcsSUFBTSxpQkFBaUIsQ0FBQyxRQUF1QztBQUNwRSxRQUFNLEVBQUUsS0FBSyxJQUFJO0FBQ2YsUUFBTSxNQUFNLFFBQVEsTUFBTSxRQUFRLElBQUksR0FBRyxFQUFFO0FBRTNDLFFBQU0sWUFBWSxDQUFDLFdBQVcsV0FBVyxXQUFXLFFBQVEsU0FBUyxTQUFTO0FBQzlFLFFBQU0sUUFBcUIsQ0FBQztBQUU1QixVQUFRLElBQUkseUJBQXlCLElBQUksa0JBQWtCO0FBRTNELFlBQVUsUUFBUSxDQUFDLFFBQVE7QUFDekIsVUFBTSxHQUFHLElBQUk7QUFBQSxNQUNYLFFBQVEsSUFBSTtBQUFBLE1BQ1osY0FBYztBQUFBLElBQ2hCO0FBQUEsRUFDRixDQUFDO0FBRUQsU0FBTztBQUNYOzs7QUQzQkEsSUFBTSxtQ0FBbUM7QUFBb0ssSUFBTSwyQ0FBMkM7QUFVOVAsSUFBTyxzQkFBUSxhQUFhLFNBQVUsS0FBSztBQUN2QyxNQUFJLFVBQVUsSUFBSSxTQUFTLE9BQU8sSUFBSTtBQUN0QyxNQUFJLE1BQU1DLFNBQVEsTUFBTSxRQUFRLElBQUksR0FBRyxFQUFFO0FBQ3pDLE1BQUksT0FBTyxZQUFZLFVBQVUsTUFBTTtBQUN2QyxNQUFJLFVBQVUsZ0JBQWdCLENBQUMsQ0FBQztBQUNoQyxTQUFPO0FBQUEsSUFDSCxTQUFTO0FBQUEsTUFDTCxJQUFJO0FBQUEsTUFDSixPQUFPO0FBQUEsTUFDUCxxQkFBcUI7QUFBQTtBQUFBLFFBRWpCLFVBQVUsQ0FBQyxLQUFLLFFBQVEsUUFBUSxJQUFJLEdBQUcsa0JBQWtCLENBQUM7QUFBQTtBQUFBLFFBRTFELFVBQVU7QUFBQSxNQUNkLENBQUM7QUFBQSxNQUNELGNBQWM7QUFBQSxRQUNWLGFBQWE7QUFBQSxRQUNiLGlCQUFpQjtBQUFBLFFBQ2pCLFNBQVMsQ0FBQyxLQUFLLFFBQVEsa0NBQVcsc0JBQXNCLENBQUM7QUFBQSxNQUM3RCxDQUFDO0FBQUEsSUFDTDtBQUFBLElBQ0EsU0FBUztBQUFBLE1BQ0wsT0FBTztBQUFBLFFBQ0gsS0FBSyxjQUFjLElBQUksSUFBSSxTQUFTLHdDQUFlLENBQUM7QUFBQSxNQUN4RDtBQUFBLElBQ0o7QUFBQSxJQUNBO0FBQUEsSUFDQSxRQUFRO0FBQUEsTUFDSixPQUFPLGVBQWUsR0FBRztBQUFBLElBQzdCO0FBQUEsSUFDQSxPQUFPO0FBQUEsTUFDSCxRQUFRLGNBQWMsSUFBSSxJQUFJLCtCQUErQix3Q0FBZSxDQUFDO0FBQUEsTUFDN0UsV0FBVztBQUFBLE1BQ1gsV0FBVyxJQUFJLG1CQUFtQjtBQUFBLE1BQ2xDLGVBQWU7QUFBQTtBQUFBLFFBRVgsVUFBVSxDQUFDO0FBQUEsUUFDWCxTQUFTLENBQUMsU0FBUyxJQUFJLDZCQUE2QixTQUFTLFdBQVcsSUFBSSxNQUFTO0FBQUEsUUFDckYsUUFBUTtBQUFBO0FBQUEsVUFFSixnQkFBZ0I7QUFBQSxVQUNoQixnQkFBZ0I7QUFBQSxVQUNoQixnQkFBZ0IsU0FBVSxXQUFXO0FBQ2pDLGdCQUFJLFlBQVksVUFBVSxLQUFLLFlBQVk7QUFDM0MsZ0JBQUksVUFBVSxTQUFTLE1BQU0sR0FBRztBQUM1QixxQkFBTztBQUFBLFlBQ1g7QUFDQSxnQkFBSSxDQUFDLE9BQU8sT0FBTyxRQUFRLEtBQUssRUFBRSxLQUFLLFNBQVUsS0FBSztBQUFFLHFCQUFPLFVBQVUsU0FBUyxHQUFHO0FBQUEsWUFBRyxDQUFDLEdBQUc7QUFDeEYscUJBQU87QUFBQSxZQUNYO0FBQ0EsZ0JBQUksQ0FBQyxPQUFPLFFBQVEsT0FBTyxFQUFFLEtBQUssU0FBVSxLQUFLO0FBQUUscUJBQU8sVUFBVSxTQUFTLEdBQUc7QUFBQSxZQUFHLENBQUMsR0FBRztBQUNuRixxQkFBTztBQUFBLFlBQ1g7QUFDQSxtQkFBTztBQUFBLFVBQ1g7QUFBQTtBQUFBLFVBRUEsY0FBYztBQUFBLFlBQ1YsY0FBYyxDQUFDLE9BQU8sY0FBYyxTQUFTLFVBQVU7QUFBQSxZQUN2RCxPQUFPLENBQUMsT0FBTztBQUFBLFlBQ2YsT0FBTyxDQUFDLE9BQU87QUFBQSxZQUNmLGFBQWEsQ0FBQyxXQUFXO0FBQUEsWUFDekIsSUFBSSxDQUFDLElBQUk7QUFBQSxVQUNiO0FBQUEsUUFDSjtBQUFBLE1BQ0o7QUFBQSxNQUNBLGNBQWMsRUFBRSxJQUFJLHVCQUF1QjtBQUFBLElBQy9DO0FBQUEsSUFDQSxjQUFjO0FBQUEsTUFDVixTQUFTLENBQUMsT0FBTyxjQUFjLFNBQVMsWUFBWSxTQUFTLFNBQVMsYUFBYSxJQUFJO0FBQUEsSUFDM0Y7QUFBQSxFQUNKO0FBQ0osQ0FBQzsiLAogICJuYW1lcyI6IFsibG9hZEVudiIsICJsb2FkRW52Il0KfQo=
