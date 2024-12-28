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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiLCAicHJveHlfY29uZmlnLnRzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiRDpcXFxcXHU4MjlEXHU5RUJCQUlcXFxcemhpbWFfY2hhdF9haVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1tb2JpbGVcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIkQ6XFxcXFx1ODI5RFx1OUVCQkFJXFxcXHpoaW1hX2NoYXRfYWlcXFxcZnJvbnQtZW5kXFxcXGNoYXQtYWktbW9iaWxlXFxcXHZpdGUuY29uZmlnLmpzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9EOi8lRTglOEElOUQlRTklQkElQkJBSS96aGltYV9jaGF0X2FpL2Zyb250LWVuZC9jaGF0LWFpLW1vYmlsZS92aXRlLmNvbmZpZy5qc1wiO2ltcG9ydCB7IGZpbGVVUkxUb1BhdGgsIFVSTCB9IGZyb20gJ25vZGU6dXJsJztcbmltcG9ydCBwYXRoIGZyb20gJ3BhdGgnO1xuaW1wb3J0IHsgZGVmaW5lQ29uZmlnLCBsb2FkRW52IH0gZnJvbSAndml0ZSc7XG5pbXBvcnQgdnVlIGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZSc7XG5pbXBvcnQgdnVlSnN4IGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZS1qc3gnO1xuaW1wb3J0IFZ1ZUkxOG5QbHVnaW4gZnJvbSAnQGludGxpZnkvdW5wbHVnaW4tdnVlLWkxOG4vdml0ZSc7XG5pbXBvcnQgeyB2aXN1YWxpemVyIH0gZnJvbSAncm9sbHVwLXBsdWdpbi12aXN1YWxpemVyJztcbmltcG9ydCBleHRlcm5hbEdsb2JhbHMgZnJvbSAncm9sbHVwLXBsdWdpbi1leHRlcm5hbC1nbG9iYWxzJztcbmltcG9ydCB7IGNyZWF0ZVN2Z0ljb25zUGx1Z2luIH0gZnJvbSAndml0ZS1wbHVnaW4tc3ZnLWljb25zJztcbmltcG9ydCB7IGdldFByb3h5Q29uZmlnIH0gZnJvbSAnLi9wcm94eV9jb25maWcnO1xuZXhwb3J0IGRlZmF1bHQgZGVmaW5lQ29uZmlnKGZ1bmN0aW9uIChvcHQpIHtcbiAgICB2YXIgY29tbWFuZCA9IG9wdC5jb21tYW5kLCBtb2RlID0gb3B0Lm1vZGU7XG4gICAgdmFyIGVudiA9IGxvYWRFbnYobW9kZSwgcHJvY2Vzcy5jd2QoKSwgJycpO1xuICAgIHZhciBiYXNlID0gY29tbWFuZCA9PT0gJ3NlcnZlJyA/ICcvJyA6ICcvJztcbiAgICB2YXIgZ2xvYmFscyA9IGV4dGVybmFsR2xvYmFscyh7fSk7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgcGx1Z2luczogW1xuICAgICAgICAgICAgdnVlKCksXG4gICAgICAgICAgICB2dWVKc3goKSxcbiAgICAgICAgICAgIGNyZWF0ZVN2Z0ljb25zUGx1Z2luKHtcbiAgICAgICAgICAgICAgICAvLyBcdTYzMDdcdTVCOUFcdTk3MDBcdTg5ODFcdTdGMTNcdTVCNThcdTc2ODRcdTU2RkVcdTY4MDdcdTY1ODdcdTRFRjZcdTU5MzlcbiAgICAgICAgICAgICAgICBpY29uRGlyczogW3BhdGgucmVzb2x2ZShwcm9jZXNzLmN3ZCgpLCAnc3JjL2Fzc2V0cy9pY29ucycpXSxcbiAgICAgICAgICAgICAgICAvLyBcdTYzMDdcdTVCOUFzeW1ib2xJZFx1NjgzQ1x1NUYwRlxuICAgICAgICAgICAgICAgIHN5bWJvbElkOiAnW25hbWVdJ1xuICAgICAgICAgICAgfSksXG4gICAgICAgICAgICBWdWVJMThuUGx1Z2luKHtcbiAgICAgICAgICAgICAgICBydW50aW1lT25seTogdHJ1ZSxcbiAgICAgICAgICAgICAgICBjb21wb3NpdGlvbk9ubHk6IHRydWUsXG4gICAgICAgICAgICAgICAgaW5jbHVkZTogW3BhdGgucmVzb2x2ZShfX2Rpcm5hbWUsICcuL3NyYy9sb2NhbGVzL2xhbmcqKicpXVxuICAgICAgICAgICAgfSlcbiAgICAgICAgXSxcbiAgICAgICAgcmVzb2x2ZToge1xuICAgICAgICAgICAgYWxpYXM6IHtcbiAgICAgICAgICAgICAgICAnQCc6IGZpbGVVUkxUb1BhdGgobmV3IFVSTCgnLi9zcmMnLCBpbXBvcnQubWV0YS51cmwpKVxuICAgICAgICAgICAgfVxuICAgICAgICB9LFxuICAgICAgICBiYXNlOiBiYXNlLFxuICAgICAgICBzZXJ2ZXI6IHtcbiAgICAgICAgICAgIHByb3h5OiBnZXRQcm94eUNvbmZpZyhvcHQpXG4gICAgICAgIH0sXG4gICAgICAgIGJ1aWxkOiB7XG4gICAgICAgICAgICBvdXREaXI6IGZpbGVVUkxUb1BhdGgobmV3IFVSTCgnLi4vLi4vc3RhdGljL2NoYXQtYWktbW9iaWxlJywgaW1wb3J0Lm1ldGEudXJsKSksXG4gICAgICAgICAgICBlbXB0eU91dERpcjogdHJ1ZSxcbiAgICAgICAgICAgIGFzc2V0c0RpcjogJ2Fzc2V0cycsXG4gICAgICAgICAgICBzb3VyY2VtYXA6IGVudi5WSVRFX1NPVVJDRU1BUCA9PT0gJ3RydWUnLFxuICAgICAgICAgICAgcm9sbHVwT3B0aW9uczoge1xuICAgICAgICAgICAgICAgIC8vIGV4dGVybmFsOiBbJ21vbWVudCcsICd2aWRlby5qcycsICdqc3BkZicsICd4bHN4J10sXG4gICAgICAgICAgICAgICAgZXh0ZXJuYWw6IFtdLFxuICAgICAgICAgICAgICAgIHBsdWdpbnM6IFtnbG9iYWxzLCBlbnYuVklURV9VU0VfQlVORExFX0FOQUxZWkVSID09PSAndHJ1ZScgPyB2aXN1YWxpemVyKCkgOiB1bmRlZmluZWRdLFxuICAgICAgICAgICAgICAgIG91dHB1dDoge1xuICAgICAgICAgICAgICAgICAgICAvLyBcdTgxRUFcdTVCOUFcdTRFNDljaHVua0ZpbGVOYW1lXHU3NTFGXHU2MjEwXHU4OUM0XHU1MjE5XG4gICAgICAgICAgICAgICAgICAgIGNodW5rRmlsZU5hbWVzOiAnYXNzZXRzL2pzL1tuYW1lXS1baGFzaF0uanMnLFxuICAgICAgICAgICAgICAgICAgICBlbnRyeUZpbGVOYW1lczogJ1tuYW1lXS1baGFzaF0uanMnLFxuICAgICAgICAgICAgICAgICAgICBhc3NldEZpbGVOYW1lczogZnVuY3Rpb24gKGFzc2V0SW5mbykge1xuICAgICAgICAgICAgICAgICAgICAgICAgdmFyIGZpZWxfbmFtZSA9IGFzc2V0SW5mby5uYW1lLnRvTG93ZXJDYXNlKCk7XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoZmllbF9uYW1lLmVuZHNXaXRoKCcuY3NzJykpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9jc3MvW25hbWVdLVtoYXNoXS5bZXh0XSc7XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoWydwbmcnLCAnanBnJywgJ2pwZWcnLCAnc3ZnJ10uc29tZShmdW5jdGlvbiAoZXh0KSB7IHJldHVybiBmaWVsX25hbWUuZW5kc1dpdGgoZXh0KTsgfSkpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9pbWcvW25hbWVdLVtoYXNoXS5bZXh0XSc7XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoWyd0dGYnLCAnd29mZicsICd3b2ZmMiddLnNvbWUoZnVuY3Rpb24gKGV4dCkgeyByZXR1cm4gZmllbF9uYW1lLmVuZHNXaXRoKGV4dCk7IH0pKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvZm9udHMvW25hbWVdLVtoYXNoXS5bZXh0XSc7XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9bbmFtZV0tW2hhc2hdLltleHRdJztcbiAgICAgICAgICAgICAgICAgICAgfSxcbiAgICAgICAgICAgICAgICAgICAgLy8gXHU4QkU1XHU5MDA5XHU5ODc5XHU1MTQxXHU4QkI4XHU0RjYwXHU1MjFCXHU1RUZBXHU4MUVBXHU1QjlBXHU0RTQ5XHU3Njg0XHU1MTZDXHU1MTcxIGNodW5rXG4gICAgICAgICAgICAgICAgICAgIG1hbnVhbENodW5rczoge1xuICAgICAgICAgICAgICAgICAgICAgICAgJ3Z1ZS1jaHVua3MnOiBbJ3Z1ZScsICd2dWUtcm91dGVyJywgJ3BpbmlhJywgJ3Z1ZS1pMThuJ10sXG4gICAgICAgICAgICAgICAgICAgICAgICBkYXlqczogWydkYXlqcyddLFxuICAgICAgICAgICAgICAgICAgICAgICAgYXhpb3M6IFsnYXhpb3MnXSxcbiAgICAgICAgICAgICAgICAgICAgICAgICdjcnlwdG8tanMnOiBbJ2NyeXB0by1qcyddLFxuICAgICAgICAgICAgICAgICAgICAgICAgcXM6IFsncXMnXSxcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgICBjc3NDb2RlU3BsaXQ6ICEoZW52LlZJVEVfVVNFX0NTU19TUExJVCA9PT0gJ2ZhbHNlJylcbiAgICAgICAgfSxcbiAgICAgICAgb3B0aW1pemVEZXBzOiB7XG4gICAgICAgICAgICBpbmNsdWRlOiBbJ3Z1ZScsICd2dWUtcm91dGVyJywgJ3BpbmlhJywgJ3Z1ZS1pMThuJywgJ2RheWpzJywgJ2F4aW9zJywgJ2NyeXB0by1qcycsICdxcyddXG4gICAgICAgIH1cbiAgICB9O1xufSk7XG4iLCAiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIkQ6XFxcXFx1ODI5RFx1OUVCQkFJXFxcXHpoaW1hX2NoYXRfYWlcXFxcZnJvbnQtZW5kXFxcXGNoYXQtYWktbW9iaWxlXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCJEOlxcXFxcdTgyOURcdTlFQkJBSVxcXFx6aGltYV9jaGF0X2FpXFxcXGZyb250LWVuZFxcXFxjaGF0LWFpLW1vYmlsZVxcXFxwcm94eV9jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0Q6LyVFOCU4QSU5RCVFOSVCQSVCQkFJL3poaW1hX2NoYXRfYWkvZnJvbnQtZW5kL2NoYXQtYWktbW9iaWxlL3Byb3h5X2NvbmZpZy50c1wiO2ltcG9ydCB7IGxvYWRFbnYgfSBmcm9tICd2aXRlJztcblxuLy8gXHU1MDQ3XHU4QkJFXHU3M0FGXHU1ODgzXHU1M0Q4XHU5MUNGXHU1NDhDXHU5MTREXHU3RjZFXHU2NTg3XHU0RUY2XHU2ODNDXHU1RjBGXHU1REYyXHU3RUNGXHU1QjlBXHU0RTQ5XHU1OTdEXG5pbnRlcmZhY2UgUHJveHlDb25maWcge1xuICBba2V5OiBzdHJpbmddOiB7XG4gICAgdGFyZ2V0OiBzdHJpbmc7XG4gICAgY2hhbmdlT3JpZ2luOiBib29sZWFuO1xuICB9O1xufVxuXG5leHBvcnQgY29uc3QgZ2V0UHJveHlDb25maWcgPSAob3B0OiB7IG1vZGU6IHN0cmluZyB9KTogUHJveHlDb25maWcgPT4ge1xuICBjb25zdCB7IG1vZGUgfSA9IG9wdDtcbiAgICBjb25zdCBlbnYgPSBsb2FkRW52KG1vZGUsIHByb2Nlc3MuY3dkKCksICcnKTtcblxuICAgIGNvbnN0IHByb3h5QXBpcyA9IFsnL3N0YXRpYycsICcvY29tbW9uJywgJy9tYW5hZ2UnLCAnL2FwcCcsICcvY2hhdCcsICcvdXBsb2FkJywgJy9wdWJsaWMnXTtcbiAgICBjb25zdCBwcm94eTogUHJveHlDb25maWcgPSB7fTtcblxuICAgIGNvbnNvbGUubG9nKCdQUk9YWV9CQVNFX0FQSV9VUkw6ICcgKyBlbnYuUFJPWFlfQkFTRV9BUElfVVJMKTtcblxuICAgIHByb3h5QXBpcy5mb3JFYWNoKChrZXkpID0+IHtcbiAgICAgIHByb3h5W2tleV0gPSB7XG4gICAgICAgIHRhcmdldDogZW52LlBST1hZX0JBU0VfQVBJX1VSTCxcbiAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlLFxuICAgICAgfTtcbiAgICB9KTtcblxuICAgIHJldHVybiBw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QVUsU0FBUyxHQUFHO0FBQUEsWUFBRyxDQUFDLEdBQUc7QUFDbkYscUJBQU87QUFBQSxZQUNYO0FBQ0EsbUJBQU87QUFBQSxVQUNYO0FBQUE7QUFBQSxVQUVBLGNBQWM7QUFBQSxZQUNWLGNBQWMsQ0FBQyxPQUFPLGNBQWMsU0FBUyxVQUFVO0FBQUEsWUFDdkQsT0FBTyxDQUFDLE9BQU87QUFBQSxZQUNmLE9BQU8sQ0FBQyxPQUFPO0FBQUEsWUFDZixhQUFhLENBQUMsV0FBVztBQUFBLFlBQ3pCLElBQUksQ0FBQyxJQUFJO0FBQUEsVUFDYjtBQUFBLFFBQ0o7QUFBQSxNQUNKO0FBQUEsTUFDQSxjQUFjLEVBQUUsSUFBSSx1QkFBdUI7QUFBQSxJQUMvQztBQUFBLElBQ0EsY0FBYztBQUFBLE1BQ1YsU0FBUyxDQUFDLE9BQU8sY0FBYyxTQUFTLFlBQVksU0FBUyxTQUFTLGFBQWEsSUFBSTtBQUFBLElBQzNGO0FBQUEsRUFDSjtBQUNKLENBQUM7IiwKICAibmFtZXMiOiBbImxvYWRFbnYiLCAibG9hZEVudiJdCn0K
