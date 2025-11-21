// vite.config.js
import { fileURLToPath, URL } from "node:url";
import path from "path";
import { defineConfig, loadEnv as loadEnv2 } from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-mobile/node_modules/vite/dist/node/index.js";
import vue from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-mobile/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-mobile/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import VueI18nPlugin from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-mobile/node_modules/@intlify/unplugin-vue-i18n/lib/vite.mjs";
import { visualizer } from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-mobile/node_modules/rollup-plugin-visualizer/dist/plugin/index.js";
import externalGlobals from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-mobile/node_modules/rollup-plugin-external-globals/index.js";
import { createSvgIconsPlugin } from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-mobile/node_modules/vite-plugin-svg-icons/dist/index.mjs";

// proxy_config.ts
import { loadEnv } from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-mobile/node_modules/vite/dist/node/index.js";
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
var __vite_injected_original_dirname = "D:\\project\\zhima_chatwiki\\front-end\\chat-ai-mobile";
var __vite_injected_original_import_meta_url = "file:///D:/project/zhima_chatwiki/front-end/chat-ai-mobile/vite.config.js";
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiLCAicHJveHlfY29uZmlnLnRzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiRDpcXFxccHJvamVjdFxcXFx6aGltYV9jaGF0d2lraVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1tb2JpbGVcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIkQ6XFxcXHByb2plY3RcXFxcemhpbWFfY2hhdHdpa2lcXFxcZnJvbnQtZW5kXFxcXGNoYXQtYWktbW9iaWxlXFxcXHZpdGUuY29uZmlnLmpzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9EOi9wcm9qZWN0L3poaW1hX2NoYXR3aWtpL2Zyb250LWVuZC9jaGF0LWFpLW1vYmlsZS92aXRlLmNvbmZpZy5qc1wiO2ltcG9ydCB7IGZpbGVVUkxUb1BhdGgsIFVSTCB9IGZyb20gJ25vZGU6dXJsJztcclxuaW1wb3J0IHBhdGggZnJvbSAncGF0aCc7XHJcbmltcG9ydCB7IGRlZmluZUNvbmZpZywgbG9hZEVudiB9IGZyb20gJ3ZpdGUnO1xyXG5pbXBvcnQgdnVlIGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZSc7XHJcbmltcG9ydCB2dWVKc3ggZnJvbSAnQHZpdGVqcy9wbHVnaW4tdnVlLWpzeCc7XHJcbmltcG9ydCBWdWVJMThuUGx1Z2luIGZyb20gJ0BpbnRsaWZ5L3VucGx1Z2luLXZ1ZS1pMThuL3ZpdGUnO1xyXG5pbXBvcnQgeyB2aXN1YWxpemVyIH0gZnJvbSAncm9sbHVwLXBsdWdpbi12aXN1YWxpemVyJztcclxuaW1wb3J0IGV4dGVybmFsR2xvYmFscyBmcm9tICdyb2xsdXAtcGx1Z2luLWV4dGVybmFsLWdsb2JhbHMnO1xyXG5pbXBvcnQgeyBjcmVhdGVTdmdJY29uc1BsdWdpbiB9IGZyb20gJ3ZpdGUtcGx1Z2luLXN2Zy1pY29ucyc7XHJcbmltcG9ydCB7IGdldFByb3h5Q29uZmlnIH0gZnJvbSAnLi9wcm94eV9jb25maWcnO1xyXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoZnVuY3Rpb24gKG9wdCkge1xyXG4gICAgdmFyIGNvbW1hbmQgPSBvcHQuY29tbWFuZCwgbW9kZSA9IG9wdC5tb2RlO1xyXG4gICAgdmFyIGVudiA9IGxvYWRFbnYobW9kZSwgcHJvY2Vzcy5jd2QoKSwgJycpO1xyXG4gICAgdmFyIGJhc2UgPSBjb21tYW5kID09PSAnc2VydmUnID8gJy8nIDogJy8nO1xyXG4gICAgdmFyIGdsb2JhbHMgPSBleHRlcm5hbEdsb2JhbHMoe30pO1xyXG4gICAgcmV0dXJuIHtcclxuICAgICAgICBwbHVnaW5zOiBbXHJcbiAgICAgICAgICAgIHZ1ZSgpLFxyXG4gICAgICAgICAgICB2dWVKc3goKSxcclxuICAgICAgICAgICAgY3JlYXRlU3ZnSWNvbnNQbHVnaW4oe1xyXG4gICAgICAgICAgICAgICAgLy8gXHU2MzA3XHU1QjlBXHU5NzAwXHU4OTgxXHU3RjEzXHU1QjU4XHU3Njg0XHU1NkZFXHU2ODA3XHU2NTg3XHU0RUY2XHU1OTM5XHJcbiAgICAgICAgICAgICAgICBpY29uRGlyczogW3BhdGgucmVzb2x2ZShwcm9jZXNzLmN3ZCgpLCAnc3JjL2Fzc2V0cy9pY29ucycpXSxcclxuICAgICAgICAgICAgICAgIC8vIFx1NjMwN1x1NUI5QXN5bWJvbElkXHU2ODNDXHU1RjBGXHJcbiAgICAgICAgICAgICAgICBzeW1ib2xJZDogJ1tuYW1lXSdcclxuICAgICAgICAgICAgfSksXHJcbiAgICAgICAgICAgIFZ1ZUkxOG5QbHVnaW4oe1xyXG4gICAgICAgICAgICAgICAgcnVudGltZU9ubHk6IHRydWUsXHJcbiAgICAgICAgICAgICAgICBjb21wb3NpdGlvbk9ubHk6IHRydWUsXHJcbiAgICAgICAgICAgICAgICBpbmNsdWRlOiBbcGF0aC5yZXNvbHZlKF9fZGlybmFtZSwgJy4vc3JjL2xvY2FsZXMvbGFuZyoqJyldXHJcbiAgICAgICAgICAgIH0pXHJcbiAgICAgICAgXSxcclxuICAgICAgICByZXNvbHZlOiB7XHJcbiAgICAgICAgICAgIGFsaWFzOiB7XHJcbiAgICAgICAgICAgICAgICAnQCc6IGZpbGVVUkxUb1BhdGgobmV3IFVSTCgnLi9zcmMnLCBpbXBvcnQubWV0YS51cmwpKVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSxcclxuICAgICAgICBiYXNlOiBiYXNlLFxyXG4gICAgICAgIHNlcnZlcjoge1xyXG4gICAgICAgICAgICBwcm94eTogZ2V0UHJveHlDb25maWcob3B0KVxyXG4gICAgICAgIH0sXHJcbiAgICAgICAgYnVpbGQ6IHtcclxuICAgICAgICAgICAgb3V0RGlyOiBmaWxlVVJMVG9QYXRoKG5ldyBVUkwoJy4uLy4uL3N0YXRpYy9jaGF0LWFpLW1vYmlsZScsIGltcG9ydC5tZXRhLnVybCkpLFxyXG4gICAgICAgICAgICBlbXB0eU91dERpcjogdHJ1ZSxcclxuICAgICAgICAgICAgYXNzZXRzRGlyOiAnYXNzZXRzJyxcclxuICAgICAgICAgICAgc291cmNlbWFwOiBlbnYuVklURV9TT1VSQ0VNQVAgPT09ICd0cnVlJyxcclxuICAgICAgICAgICAgcm9sbHVwT3B0aW9uczoge1xyXG4gICAgICAgICAgICAgICAgLy8gZXh0ZXJuYWw6IFsnbW9tZW50JywgJ3ZpZGVvLmpzJywgJ2pzcGRmJywgJ3hsc3gnXSxcclxuICAgICAgICAgICAgICAgIGV4dGVybmFsOiBbXSxcclxuICAgICAgICAgICAgICAgIHBsdWdpbnM6IFtnbG9iYWxzLCBlbnYuVklURV9VU0VfQlVORExFX0FOQUxZWkVSID09PSAndHJ1ZScgPyB2aXN1YWxpemVyKCkgOiB1bmRlZmluZWRdLFxyXG4gICAgICAgICAgICAgICAgb3V0cHV0OiB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gXHU4MUVBXHU1QjlBXHU0RTQ5Y2h1bmtGaWxlTmFtZVx1NzUxRlx1NjIxMFx1ODlDNFx1NTIxOVxyXG4gICAgICAgICAgICAgICAgICAgIGNodW5rRmlsZU5hbWVzOiAnYXNzZXRzL2pzL1tuYW1lXS1baGFzaF0uanMnLFxyXG4gICAgICAgICAgICAgICAgICAgIGVudHJ5RmlsZU5hbWVzOiAnW25hbWVdLVtoYXNoXS5qcycsXHJcbiAgICAgICAgICAgICAgICAgICAgYXNzZXRGaWxlTmFtZXM6IGZ1bmN0aW9uIChhc3NldEluZm8pIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgdmFyIGZpZWxfbmFtZSA9IGFzc2V0SW5mby5uYW1lLnRvTG93ZXJDYXNlKCk7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChmaWVsX25hbWUuZW5kc1dpdGgoJy5jc3MnKSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvY3NzL1tuYW1lXS1baGFzaF0uW2V4dF0nO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChbJ3BuZycsICdqcGcnLCAnanBlZycsICdzdmcnXS5zb21lKGZ1bmN0aW9uIChleHQpIHsgcmV0dXJuIGZpZWxfbmFtZS5lbmRzV2l0aChleHQpOyB9KSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvaW1nL1tuYW1lXS1baGFzaF0uW2V4dF0nO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChbJ3R0ZicsICd3b2ZmJywgJ3dvZmYyJ10uc29tZShmdW5jdGlvbiAoZXh0KSB7IHJldHVybiBmaWVsX25hbWUuZW5kc1dpdGgoZXh0KTsgfSkpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiAnYXNzZXRzL2ZvbnRzL1tuYW1lXS1baGFzaF0uW2V4dF0nO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiAnYXNzZXRzL1tuYW1lXS1baGFzaF0uW2V4dF0nO1xyXG4gICAgICAgICAgICAgICAgICAgIH0sXHJcbiAgICAgICAgICAgICAgICAgICAgLy8gXHU4QkU1XHU5MDA5XHU5ODc5XHU1MTQxXHU4QkI4XHU0RjYwXHU1MjFCXHU1RUZBXHU4MUVBXHU1QjlBXHU0RTQ5XHU3Njg0XHU1MTZDXHU1MTcxIGNodW5rXHJcbiAgICAgICAgICAgICAgICAgICAgbWFudWFsQ2h1bmtzOiB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICd2dWUtY2h1bmtzJzogWyd2dWUnLCAndnVlLXJvdXRlcicsICdwaW5pYScsICd2dWUtaTE4biddLFxyXG4gICAgICAgICAgICAgICAgICAgICAgICBkYXlqczogWydkYXlqcyddLFxyXG4gICAgICAgICAgICAgICAgICAgICAgICBheGlvczogWydheGlvcyddLFxyXG4gICAgICAgICAgICAgICAgICAgICAgICAnY3J5cHRvLWpzJzogWydjcnlwdG8tanMnXSxcclxuICAgICAgICAgICAgICAgICAgICAgICAgcXM6IFsncXMnXSxcclxuICAgICAgICAgICAgICAgICAgICAgICAgJ3Z1ZS1wZGYtZW1iZWQnOiBbJ3Z1ZS1wZGYtZW1iZWQnXVxyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfSxcclxuICAgICAgICAgICAgY3NzQ29kZVNwbGl0OiAhKGVudi5WSVRFX1VTRV9DU1NfU1BMSVQgPT09ICdmYWxzZScpXHJcbiAgICAgICAgfSxcclxuICAgICAgICBvcHRpbWl6ZURlcHM6IHtcclxuICAgICAgICAgICAgaW5jbHVkZTogWyd2dWUnLCAndnVlLXJvdXRlcicsICdwaW5pYScsICd2dWUtaTE4bicsICdkYXlqcycsICdheGlvcycsICdjcnlwdG8tanMnLCAncXMnXVxyXG4gICAgICAgIH1cclxuICAgIH07XHJcbn0pO1xyXG4iLCAiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIkQ6XFxcXHByb2plY3RcXFxcemhpbWFfY2hhdHdpa2lcXFxcZnJvbnQtZW5kXFxcXGNoYXQtYWktbW9iaWxlXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCJEOlxcXFxwcm9qZWN0XFxcXHpoaW1hX2NoYXR3aWtpXFxcXGZyb250LWVuZFxcXFxjaGF0LWFpLW1vYmlsZVxcXFxwcm94eV9jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0Q6L3Byb2plY3QvemhpbWFfY2hhdHdpa2kvZnJvbnQtZW5kL2NoYXQtYWktbW9iaWxlL3Byb3h5X2NvbmZpZy50c1wiO2ltcG9ydCB7IGxvYWRFbnYgfSBmcm9tICd2aXRlJztcclxuXHJcbi8vIFx1NTA0N1x1OEJCRVx1NzNBRlx1NTg4M1x1NTNEOFx1OTFDRlx1NTQ4Q1x1OTE0RFx1N0Y2RVx1NjU4N1x1NEVGNlx1NjgzQ1x1NUYwRlx1NURGMlx1N0VDRlx1NUI5QVx1NEU0OVx1NTk3RFxyXG5pbnRlcmZhY2UgUHJveHlDb25maWcge1xyXG4gIFtrZXk6IHN0cmluZ106IHtcclxuICAgIHRhcmdldDogc3RyaW5nO1xyXG4gICAgY2hhbmdlT3JpZ2luOiBib29sZWFuO1xyXG4gIH07XHJcbn1cclxuXHJcbmV4cG9ydCBjb25zdCBnZXRQcm94eUNvbmZpZyA9IChvcHQ6IHsgbW9kZTogc3RyaW5nIH0pOiBQcm94eUNvbmZpZyA9PiB7XHJcbiAgY29uc3QgeyBtb2RlIH0gPSBvcHQ7XHJcbiAgICBjb25zdCBlbnYgPSBsb2FkRW52KG1vZGUsIHByb2Nlc3MuY3dkKCksICcnKTtcclxuXHJcbiAgICBjb25zdCBwcm94eUFwaXMgPSBbJy9zdGF0aWMnLCAnL2NvbW1vbicsICcvbWFuYWdlJywgJy9hcHAnLCAnL2NoYXQnLCAnL3VwbG9hZCcsICcvcHVibGljJ107XHJcbiAgICBjb25zdCBwcm94eTogUHJveHlDb25maWcgPSB7fTtcclxuXHJcbiAgICBjb25zb2xlLmxvZygnUFJPWFlfQkFTRV9BUElfVVJMOiAnICsgZW52LlBST1hZX0JBU0VfQVBJX1VSTCk7XHJcblxyXG4gICAgcHJveHlBcGlzLmZvckVhY2goKGtleSkgPT4ge1xyXG4gICAgICBwcm94eVtrZXldID0ge1xyXG4gICAgICAgIHRhcmdldDogZW52LlBST1hZX0JBU0VfQVBJX1VSTCxcclxuICAgICAgICBjaGFuZ2VPcmlnaW46IHRydWUsXHJcbiAgICAgIH07XHJcbiAgICB9KTtcclxuXHJcbiAgICByZXR1cm4gcHJveHk7XHJcbn07XHJcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFBa1YsU0FBUyxlQUFlLFdBQVc7QUFDclgsT0FBTyxVQUFVO0FBQ2pCLFNBQVMsY0FBYyxXQUFBQSxnQkFBZTtBQUN0QyxPQUFPLFNBQVM7QUFDaEIsT0FBTyxZQUFZO0FBQ25CLE9BQU8sbUJBQW1CO0FBQzFCLFNBQVMsa0JBQWtCO0FBQzNCLE9BQU8scUJBQXFCO0FBQzVCLFNBQVMsNEJBQTRCOzs7QUNSK1MsU0FBUyxlQUFlO0FBVXJXLElBQU0saUJBQWlCLENBQUMsUUFBdUM7QUFDcEUsUUFBTSxFQUFFLEtBQUssSUFBSTtBQUNmLFFBQU0sTUFBTSxRQUFRLE1BQU0sUUFBUSxJQUFJLEdBQUcsRUFBRTtBQUUzQyxRQUFNLFlBQVksQ0FBQyxXQUFXLFdBQVcsV0FBVyxRQUFRLFNBQVMsV0FBVyxTQUFTO0FBQ3pGLFFBQU0sUUFBcUIsQ0FBQztBQUU1QixVQUFRLElBQUkseUJBQXlCLElBQUksa0JBQWtCO0FBRTNELFlBQVUsUUFBUSxDQUFDLFFBQVE7QUFDekIsVUFBTSxHQUFHLElBQUk7QUFBQSxNQUNYLFFBQVEsSUFBSTtBQUFBLE1BQ1osY0FBYztBQUFBLElBQ2hCO0FBQUEsRUFDRixDQUFDO0FBRUQsU0FBTztBQUNYOzs7QUQzQkEsSUFBTSxtQ0FBbUM7QUFBNEssSUFBTSwyQ0FBMkM7QUFVdFEsSUFBTyxzQkFBUSxhQUFhLFNBQVUsS0FBSztBQUN2QyxNQUFJLFVBQVUsSUFBSSxTQUFTLE9BQU8sSUFBSTtBQUN0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
