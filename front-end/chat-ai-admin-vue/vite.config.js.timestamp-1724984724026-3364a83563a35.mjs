// vite.config.js
import { fileURLToPath, URL } from "node:url";
import path from "path";
import { defineConfig, loadEnv as loadEnv2 } from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/vite/dist/node/index.js";
import vue from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import { visualizer } from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/rollup-plugin-visualizer/dist/plugin/index.js";
import externalGlobals from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/rollup-plugin-external-globals/index.js";
import Components from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/unplugin-vue-components/dist/vite.js";
import { AntDesignVueResolver } from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/unplugin-vue-components/dist/resolvers.js";
import VueI18nPlugin from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/@intlify/unplugin-vue-i18n/lib/vite.mjs";
import { createSvgIconsPlugin } from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/vite-plugin-svg-icons/dist/index.mjs";

// proxy_config.js
import { loadEnv } from "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/node_modules/vite/dist/node/index.js";
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
var __vite_injected_original_dirname = "D:\\project\\zhima_chat_ai\\front-end\\chat-ai-admin-vue";
var __vite_injected_original_import_meta_url = "file:///D:/project/zhima_chat_ai/front-end/chat-ai-admin-vue/vite.config.js";
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
      proxy: getProxyConfig(opt),
      port: 5520,
      open: true
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiLCAicHJveHlfY29uZmlnLmpzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiRDpcXFxccHJvamVjdFxcXFx6aGltYV9jaGF0X2FpXFxcXGZyb250LWVuZFxcXFxjaGF0LWFpLWFkbWluLXZ1ZVwiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiRDpcXFxccHJvamVjdFxcXFx6aGltYV9jaGF0X2FpXFxcXGZyb250LWVuZFxcXFxjaGF0LWFpLWFkbWluLXZ1ZVxcXFx2aXRlLmNvbmZpZy5qc1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vRDovcHJvamVjdC96aGltYV9jaGF0X2FpL2Zyb250LWVuZC9jaGF0LWFpLWFkbWluLXZ1ZS92aXRlLmNvbmZpZy5qc1wiOy8qIGVzbGludC1kaXNhYmxlIG5vLXVuZGVmICovXHJcbmltcG9ydCB7IGZpbGVVUkxUb1BhdGgsIFVSTCB9IGZyb20gJ25vZGU6dXJsJ1xyXG5pbXBvcnQgcGF0aCBmcm9tICdwYXRoJ1xyXG5cclxuaW1wb3J0IHsgZGVmaW5lQ29uZmlnLCBsb2FkRW52IH0gZnJvbSAndml0ZSdcclxuaW1wb3J0IHZ1ZSBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUnXHJcbmltcG9ydCB2dWVKc3ggZnJvbSAnQHZpdGVqcy9wbHVnaW4tdnVlLWpzeCdcclxuLy8gaW1wb3J0IFZ1ZURldlRvb2xzIGZyb20gJ3ZpdGUtcGx1Z2luLXZ1ZS1kZXZ0b29scydcclxuaW1wb3J0IHsgdmlzdWFsaXplciB9IGZyb20gJ3JvbGx1cC1wbHVnaW4tdmlzdWFsaXplcidcclxuaW1wb3J0IGV4dGVybmFsR2xvYmFscyBmcm9tICdyb2xsdXAtcGx1Z2luLWV4dGVybmFsLWdsb2JhbHMnXHJcbmltcG9ydCBDb21wb25lbnRzIGZyb20gJ3VucGx1Z2luLXZ1ZS1jb21wb25lbnRzL3ZpdGUnXHJcbmltcG9ydCB7IEFudERlc2lnblZ1ZVJlc29sdmVyIH0gZnJvbSAndW5wbHVnaW4tdnVlLWNvbXBvbmVudHMvcmVzb2x2ZXJzJ1xyXG5pbXBvcnQgVnVlSTE4blBsdWdpbiBmcm9tICdAaW50bGlmeS91bnBsdWdpbi12dWUtaTE4bi92aXRlJ1xyXG5pbXBvcnQgeyBjcmVhdGVTdmdJY29uc1BsdWdpbiB9IGZyb20gJ3ZpdGUtcGx1Z2luLXN2Zy1pY29ucydcclxuaW1wb3J0IHsgZ2V0UHJveHlDb25maWcgfSBmcm9tICcuL3Byb3h5X2NvbmZpZydcclxuXHJcbmNvbnN0IHsgcmVzb2x2ZSB9ID0gcGF0aFxyXG5jb25zdCByb290ID0gcHJvY2Vzcy5jd2QoKVxyXG5mdW5jdGlvbiBwYXRoUmVzb2x2ZShkaXIpIHtcclxuICByZXR1cm4gcmVzb2x2ZShyb290LCAnLicsIGRpcilcclxufVxyXG5cclxuLy8gY29uc3QgZ2xvYmFscyA9IGV4dGVybmFsR2xvYmFscyh7XHJcbi8vICAgbW9tZW50OiAnbW9tZW50JyxcclxuLy8gICAndmlkZW8uanMnOiAndmlkZW9qcycsXHJcbi8vICAganNwZGY6ICdqc3BkZicsXHJcbi8vICAgeGxzeDogJ1hMU1gnLFxyXG4vLyB9KTtcclxuXHJcbmNvbnN0IGdsb2JhbHMgPSBleHRlcm5hbEdsb2JhbHMoe30pXHJcblxyXG4vLyBodHRwczovL3ZpdGVqcy5kZXYvY29uZmlnL1xyXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoKG9wdCkgPT4ge1xyXG4gIGNvbnN0IHsgY29tbWFuZCwgbW9kZSB9ID0gb3B0XHJcbiAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lIG5vLXVudXNlZC12YXJzXHJcbiAgY29uc3QgZW52ID0gbG9hZEVudihtb2RlLCBwcm9jZXNzLmN3ZCgpLCAnJylcclxuICBjb25zdCBiYXNlID0gY29tbWFuZCA9PT0gJ3NlcnZlJyA/ICcvJyA6ICcvJ1xyXG5cclxuICByZXR1cm4ge1xyXG4gICAgcGx1Z2luczogW1xyXG4gICAgICB2dWUoKSxcclxuICAgICAgdnVlSnN4KCksXHJcbiAgICAgIC8vIFZ1ZURldlRvb2xzKCksXHJcbiAgICAgIENvbXBvbmVudHMoe1xyXG4gICAgICAgIHJlc29sdmVyczogW1xyXG4gICAgICAgICAgQW50RGVzaWduVnVlUmVzb2x2ZXIoe1xyXG4gICAgICAgICAgICBpbXBvcnRTdHlsZTogZmFsc2UgLy8gY3NzIGluIGpzXHJcbiAgICAgICAgICB9KVxyXG4gICAgICAgIF1cclxuICAgICAgfSksXHJcbiAgICAgIGNyZWF0ZVN2Z0ljb25zUGx1Z2luKHtcclxuICAgICAgICAvLyBcdTYzMDdcdTVCOUFcdTk3MDBcdTg5ODFcdTdGMTNcdTVCNThcdTc2ODRcdTU2RkVcdTY4MDdcdTY1ODdcdTRFRjZcdTU5MzlcclxuICAgICAgICBpY29uRGlyczogW3Jlc29sdmUoJy4vc3JjL2Fzc2V0cy9zdmcnKV0sXHJcbiAgICAgICAgLy8gXHU2MzA3XHU1QjlBc3ltYm9sSWRcdTY4M0NcdTVGMEZcclxuICAgICAgICBzeW1ib2xJZDogJ1tuYW1lXSdcclxuICAgICAgfSksXHJcbiAgICAgIFZ1ZUkxOG5QbHVnaW4oe1xyXG4gICAgICAgIHJ1bnRpbWVPbmx5OiB0cnVlLFxyXG4gICAgICAgIGNvbXBvc2l0aW9uT25seTogdHJ1ZSxcclxuICAgICAgICBpbmNsdWRlOiBbcGF0aC5yZXNvbHZlKF9fZGlybmFtZSwgJy4vc3JjL2xvY2FsZXMvbGFuZyoqJyldXHJcbiAgICAgIH0pXHJcbiAgICAgIC8vIGNvcHlJbmRleCgpLFxyXG4gICAgXSxcclxuICAgIHJlc29sdmU6IHtcclxuICAgICAgYWxpYXM6IFtcclxuICAgICAgICB7XHJcbiAgICAgICAgICBmaW5kOiAndnVlLWkxOG4nLFxyXG4gICAgICAgICAgcmVwbGFjZW1lbnQ6ICd2dWUtaTE4bi9kaXN0L3Z1ZS1pMThuLmNqcy5qcydcclxuICAgICAgICB9LFxyXG4gICAgICAgIHtcclxuICAgICAgICAgIGZpbmQ6IC9AXFwvLyxcclxuICAgICAgICAgIHJlcGxhY2VtZW50OiBgJHtwYXRoUmVzb2x2ZSgnc3JjJyl9L2BcclxuICAgICAgICB9XHJcbiAgICAgIF1cclxuICAgIH0sXHJcbiAgICBiYXNlOiBiYXNlLFxyXG4gICAgZXhwZXJpbWVudGFsOiB7XHJcbiAgICAgIC8vIFx1OEZEQlx1OTYzNlx1NTdGQVx1Nzg0MFx1OERFRlx1NUY4NFx1OTAwOVx1OTg3OVxyXG4gICAgfSxcclxuICAgIHNlcnZlcjoge1xyXG4gICAgICBwcm94eTogZ2V0UHJveHlDb25maWcob3B0KSxcclxuICAgICAgcG9ydDogNTUyMCxcclxuICAgICAgb3BlbjogdHJ1ZSxcclxuICAgIH0sXHJcbiAgICBlc2J1aWxkOiB7XHJcbiAgICAgIHB1cmU6IGVudi5WSVRFX0RST1BfQ09OU09MRSA9PT0gJ3RydWUnID8gWydjb25zb2xlLmxvZyddIDogdW5kZWZpbmVkLFxyXG4gICAgICBkcm9wOiBlbnYuVklURV9EUk9QX0RFQlVHR0VSID09PSAndHJ1ZScgPyBbJ2RlYnVnZ2VyJ10gOiB1bmRlZmluZWRcclxuICAgIH0sXHJcbiAgICBidWlsZDoge1xyXG4gICAgICBvdXREaXI6IGZpbGVVUkxUb1BhdGgobmV3IFVSTCgnLi4vLi4vc3RhdGljL2NoYXQtYWktYWRtaW4nLCBpbXBvcnQubWV0YS51cmwpKSxcclxuICAgICAgZW1wdHlPdXREaXI6IHRydWUsXHJcbiAgICAgIGFzc2V0c0RpcjogJ2Fzc2V0cycsXHJcbiAgICAgIHNvdXJjZW1hcDogZW52LlZJVEVfU09VUkNFTUFQID09PSAndHJ1ZScsXHJcbiAgICAgIHJvbGx1cE9wdGlvbnM6IHtcclxuICAgICAgICAvLyBleHRlcm5hbDogWydtb21lbnQnLCAndmlkZW8uanMnLCAnanNwZGYnLCAneGxzeCddLFxyXG4gICAgICAgIGV4dGVybmFsOiBbXSxcclxuICAgICAgICBwbHVnaW5zOiBbZ2xvYmFscywgZW52LlZJVEVfVVNFX0JVTkRMRV9BTkFMWVpFUiA9PT0gJ3RydWUnID8gdmlzdWFsaXplcigpIDogdW5kZWZpbmVkXSxcclxuICAgICAgICBvdXRwdXQ6IHtcclxuICAgICAgICAgIC8vIFx1ODFFQVx1NUI5QVx1NEU0OWNodW5rRmlsZU5hbWVcdTc1MUZcdTYyMTBcdTg5QzRcdTUyMTlcclxuICAgICAgICAgIGNodW5rRmlsZU5hbWVzOiAnYXNzZXRzL2pzL1tuYW1lXS1baGFzaF0uanMnLFxyXG4gICAgICAgICAgZW50cnlGaWxlTmFtZXM6ICdbbmFtZV0tW2hhc2hdLmpzJyxcclxuICAgICAgICAgIGFzc2V0RmlsZU5hbWVzKGFzc2V0SW5mbykge1xyXG4gICAgICAgICAgICBsZXQgZmllbF9uYW1lID0gYXNzZXRJbmZvLm5hbWUudG9Mb3dlckNhc2UoKVxyXG5cclxuICAgICAgICAgICAgaWYgKGZpZWxfbmFtZS5lbmRzV2l0aCgnLmNzcycpKSB7XHJcbiAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvY3NzL1tuYW1lXS1baGFzaF0uW2V4dF0nXHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgaWYgKFsncG5nJywgJ2pwZycsICdqcGVnJywgJ3N2ZyddLnNvbWUoKGV4dCkgPT4gZmllbF9uYW1lLmVuZHNXaXRoKGV4dCkpKSB7XHJcbiAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvaW1nL1tuYW1lXS1baGFzaF0uW2V4dF0nXHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgaWYgKFsndHRmJywgJ3dvZmYnLCAnd29mZjInXS5zb21lKChleHQpID0+IGZpZWxfbmFtZS5lbmRzV2l0aChleHQpKSkge1xyXG4gICAgICAgICAgICAgIHJldHVybiAnYXNzZXRzL2ZvbnRzL1tuYW1lXS1baGFzaF0uW2V4dF0nXHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvW25hbWVdLVtoYXNoXS5bZXh0XSdcclxuICAgICAgICAgIH0sXHJcbiAgICAgICAgICAvLyBcdThCRTVcdTkwMDlcdTk4NzlcdTUxNDFcdThCQjhcdTRGNjBcdTUyMUJcdTVFRkFcdTgxRUFcdTVCOUFcdTRFNDlcdTc2ODRcdTUxNkNcdTUxNzEgY2h1bmtcclxuICAgICAgICAgIG1hbnVhbENodW5rczoge1xyXG4gICAgICAgICAgICAndnVlLWNodW5rcyc6IFsndnVlJywgJ3Z1ZS1yb3V0ZXInLCAncGluaWEnLCAndnVlLWkxOG4nXSxcclxuICAgICAgICAgICAgZGF5anM6IFsnZGF5anMnXSxcclxuICAgICAgICAgICAgYXhpb3M6IFsnYXhpb3MnXSxcclxuICAgICAgICAgICAgJ2NyeXB0by1qcyc6IFsnY3J5cHRvLWpzJ10sXHJcbiAgICAgICAgICAgIHFzOiBbJ3FzJ10sXHJcbiAgICAgICAgICAgICd2dWUtcGRmLWVtYmVkJzogWyd2dWUtcGRmLWVtYmVkJ11cclxuICAgICAgICAgIH1cclxuICAgICAgICB9XHJcbiAgICAgIH0sXHJcbiAgICAgIGNzc0NvZGVTcGxpdDogIShlbnYuVklURV9VU0VfQ1NTX1NQTElUID09PSAnZmFsc2UnKVxyXG4gICAgfSxcclxuICAgIG9wdGltaXplRGVwczoge1xyXG4gICAgICBpbmNsdWRlOiBbJ3Z1ZScsICd2dWUtcm91dGVyJywgJ3BpbmlhJywgJ3Z1ZS1pMThuJywgJ2RheWpzJywgJ2F4aW9zJywgJ2NyeXB0by1qcycsICdxcyddXHJcbiAgICB9XHJcbiAgfVxyXG59KVxyXG4iLCAiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIkQ6XFxcXHByb2plY3RcXFxcemhpbWFfY2hhdF9haVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1hZG1pbi12dWVcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIkQ6XFxcXHByb2plY3RcXFxcemhpbWFfY2hhdF9haVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1hZG1pbi12dWVcXFxccHJveHlfY29uZmlnLmpzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9EOi9wcm9qZWN0L3poaW1hX2NoYXRfYWkvZnJvbnQtZW5kL2NoYXQtYWktYWRtaW4tdnVlL3Byb3h5X2NvbmZpZy5qc1wiOy8qIGVzbGludC1kaXNhYmxlIG5vLXVuZGVmICovXHJcbmltcG9ydCB7IGxvYWRFbnYgfSBmcm9tICd2aXRlJ1xyXG5cclxuZXhwb3J0IGNvbnN0IGdldFByb3h5Q29uZmlnID0gKG9wdCkgPT4ge1xyXG4gIGNvbnN0IHsgbW9kZSB9ID0gb3B0XHJcbiAgY29uc3QgZW52ID0gbG9hZEVudihtb2RlLCBwcm9jZXNzLmN3ZCgpLCAnJylcclxuXHJcbiAgbGV0IHByb3h5QXBpcyA9IFsnL3N0YXRpYycsICcvY29tbW9uJywgJy9tYW5hZ2UnLCAnL2FwcCcsICcvY2hhdCcsICcvdXBsb2FkJywgJy9wdWJsaWMnXVxyXG4gIGxldCBwcm94eSA9IHt9XHJcblxyXG4gIGNvbnNvbGUubG9nKGVudi5QUk9YWV9CQVNFX0FQSV9VUkwpXHJcblxyXG4gIHByb3h5QXBpcy5mb3JFYWNoKChrZXkpID0+IHtcclxuICAgIHByb3h5W2tleV0gPSB7XHJcbiAgICAgIHRhcmdldDogZW52LlBST1hZX0JBU0VfQVBJX1VSTCxcclxuICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlXHJcbiAgICB9XHJcbiAgfSlcclxuXHJcbiAgcmV0dXJuIHtcclxuICAgIC4uLnByb3h5XHJcbiAgfVxyXG59XHJcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFDQSxTQUFTLGVBQWUsV0FBVztBQUNuQyxPQUFPLFVBQVU7QUFFakIsU0FBUyxjQUFjLFdBQUFBLGdCQUFlO0FBQ3RDLE9BQU8sU0FBUztBQUNoQixPQUFPLFlBQVk7QUFFbkIsU0FBUyxrQkFBa0I7QUFDM0IsT0FBTyxxQkFBcUI7QUFDNUIsT0FBTyxnQkFBZ0I7QUFDdkIsU0FBUyw0QkFBNEI7QUFDckMsT0FBTyxtQkFBbUI7QUFDMUIsU0FBUyw0QkFBNEI7OztBQ1pyQyxTQUFTLGVBQWU7QUFFakIsSUFBTSxpQkFBaUIsQ0FBQyxRQUFRO0FB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
