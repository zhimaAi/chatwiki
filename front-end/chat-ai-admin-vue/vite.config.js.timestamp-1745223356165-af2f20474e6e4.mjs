// vite.config.js
import { fileURLToPath, URL } from "node:url";
import path from "path";
import { defineConfig, loadEnv as loadEnv2 } from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/node_modules/vite/dist/node/index.js";
import vue from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import { visualizer } from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/node_modules/rollup-plugin-visualizer/dist/plugin/index.js";
import externalGlobals from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/node_modules/rollup-plugin-external-globals/index.js";
import Components from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/node_modules/unplugin-vue-components/dist/vite.js";
import { AntDesignVueResolver } from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/node_modules/unplugin-vue-components/dist/resolvers.js";
import VueI18nPlugin from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/node_modules/@intlify/unplugin-vue-i18n/lib/vite.mjs";
import { createSvgIconsPlugin } from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/node_modules/vite-plugin-svg-icons/dist/index.mjs";

// proxy_config.js
import { loadEnv } from "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/node_modules/vite/dist/node/index.js";
var getProxyConfig = (opt) => {
  const { mode } = opt;
  const env = loadEnv(mode, process.cwd(), "");
  let proxyApis = ["/open/", "/static", "/common", "/manage", "/app", "/chat", "/upload", "/public"];
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
var __vite_injected_original_dirname = "D:\\project\\zhima_chatwiki\\front-end\\chat-ai-admin-vue";
var __vite_injected_original_import_meta_url = "file:///D:/project/zhima_chatwiki/front-end/chat-ai-admin-vue/vite.config.js";
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
      host: "0.0.0.0",
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiLCAicHJveHlfY29uZmlnLmpzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiRDpcXFxccHJvamVjdFxcXFx6aGltYV9jaGF0d2lraVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1hZG1pbi12dWVcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIkQ6XFxcXHByb2plY3RcXFxcemhpbWFfY2hhdHdpa2lcXFxcZnJvbnQtZW5kXFxcXGNoYXQtYWktYWRtaW4tdnVlXFxcXHZpdGUuY29uZmlnLmpzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9EOi9wcm9qZWN0L3poaW1hX2NoYXR3aWtpL2Zyb250LWVuZC9jaGF0LWFpLWFkbWluLXZ1ZS92aXRlLmNvbmZpZy5qc1wiOy8qIGVzbGludC1kaXNhYmxlIG5vLXVuZGVmICovXHJcbmltcG9ydCB7IGZpbGVVUkxUb1BhdGgsIFVSTCB9IGZyb20gJ25vZGU6dXJsJ1xyXG5pbXBvcnQgcGF0aCBmcm9tICdwYXRoJ1xyXG5cclxuaW1wb3J0IHsgZGVmaW5lQ29uZmlnLCBsb2FkRW52IH0gZnJvbSAndml0ZSdcclxuaW1wb3J0IHZ1ZSBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUnXHJcbmltcG9ydCB2dWVKc3ggZnJvbSAnQHZpdGVqcy9wbHVnaW4tdnVlLWpzeCdcclxuLy8gaW1wb3J0IFZ1ZURldlRvb2xzIGZyb20gJ3ZpdGUtcGx1Z2luLXZ1ZS1kZXZ0b29scydcclxuaW1wb3J0IHsgdmlzdWFsaXplciB9IGZyb20gJ3JvbGx1cC1wbHVnaW4tdmlzdWFsaXplcidcclxuaW1wb3J0IGV4dGVybmFsR2xvYmFscyBmcm9tICdyb2xsdXAtcGx1Z2luLWV4dGVybmFsLWdsb2JhbHMnXHJcbmltcG9ydCBDb21wb25lbnRzIGZyb20gJ3VucGx1Z2luLXZ1ZS1jb21wb25lbnRzL3ZpdGUnXHJcbmltcG9ydCB7IEFudERlc2lnblZ1ZVJlc29sdmVyIH0gZnJvbSAndW5wbHVnaW4tdnVlLWNvbXBvbmVudHMvcmVzb2x2ZXJzJ1xyXG5pbXBvcnQgVnVlSTE4blBsdWdpbiBmcm9tICdAaW50bGlmeS91bnBsdWdpbi12dWUtaTE4bi92aXRlJ1xyXG5pbXBvcnQgeyBjcmVhdGVTdmdJY29uc1BsdWdpbiB9IGZyb20gJ3ZpdGUtcGx1Z2luLXN2Zy1pY29ucydcclxuaW1wb3J0IHsgZ2V0UHJveHlDb25maWcgfSBmcm9tICcuL3Byb3h5X2NvbmZpZydcclxuXHJcbmNvbnN0IHsgcmVzb2x2ZSB9ID0gcGF0aFxyXG5jb25zdCByb290ID0gcHJvY2Vzcy5jd2QoKVxyXG5mdW5jdGlvbiBwYXRoUmVzb2x2ZShkaXIpIHtcclxuICByZXR1cm4gcmVzb2x2ZShyb290LCAnLicsIGRpcilcclxufVxyXG5cclxuLy8gY29uc3QgZ2xvYmFscyA9IGV4dGVybmFsR2xvYmFscyh7XHJcbi8vICAgbW9tZW50OiAnbW9tZW50JyxcclxuLy8gICAndmlkZW8uanMnOiAndmlkZW9qcycsXHJcbi8vICAganNwZGY6ICdqc3BkZicsXHJcbi8vICAgeGxzeDogJ1hMU1gnLFxyXG4vLyB9KTtcclxuXHJcbmNvbnN0IGdsb2JhbHMgPSBleHRlcm5hbEdsb2JhbHMoe30pXHJcblxyXG4vLyBodHRwczovL3ZpdGVqcy5kZXYvY29uZmlnL1xyXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoKG9wdCkgPT4ge1xyXG4gIGNvbnN0IHsgY29tbWFuZCwgbW9kZSB9ID0gb3B0XHJcbiAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lIG5vLXVudXNlZC12YXJzXHJcbiAgY29uc3QgZW52ID0gbG9hZEVudihtb2RlLCBwcm9jZXNzLmN3ZCgpLCAnJylcclxuICBjb25zdCBiYXNlID0gY29tbWFuZCA9PT0gJ3NlcnZlJyA/ICcvJyA6ICcvJ1xyXG5cclxuICByZXR1cm4ge1xyXG4gICAgcGx1Z2luczogW1xyXG4gICAgICB2dWUoKSxcclxuICAgICAgdnVlSnN4KCksXHJcbiAgICAgIC8vIFZ1ZURldlRvb2xzKCksXHJcbiAgICAgIENvbXBvbmVudHMoe1xyXG4gICAgICAgIHJlc29sdmVyczogW1xyXG4gICAgICAgICAgQW50RGVzaWduVnVlUmVzb2x2ZXIoe1xyXG4gICAgICAgICAgICBpbXBvcnRTdHlsZTogZmFsc2UgLy8gY3NzIGluIGpzXHJcbiAgICAgICAgICB9KVxyXG4gICAgICAgIF1cclxuICAgICAgfSksXHJcbiAgICAgIGNyZWF0ZVN2Z0ljb25zUGx1Z2luKHtcclxuICAgICAgICAvLyBcdTYzMDdcdTVCOUFcdTk3MDBcdTg5ODFcdTdGMTNcdTVCNThcdTc2ODRcdTU2RkVcdTY4MDdcdTY1ODdcdTRFRjZcdTU5MzlcclxuICAgICAgICBpY29uRGlyczogW3Jlc29sdmUoJy4vc3JjL2Fzc2V0cy9zdmcnKV0sXHJcbiAgICAgICAgLy8gXHU2MzA3XHU1QjlBc3ltYm9sSWRcdTY4M0NcdTVGMEZcclxuICAgICAgICBzeW1ib2xJZDogJ1tuYW1lXSdcclxuICAgICAgfSksXHJcbiAgICAgIFZ1ZUkxOG5QbHVnaW4oe1xyXG4gICAgICAgIHJ1bnRpbWVPbmx5OiB0cnVlLFxyXG4gICAgICAgIGNvbXBvc2l0aW9uT25seTogdHJ1ZSxcclxuICAgICAgICBpbmNsdWRlOiBbcGF0aC5yZXNvbHZlKF9fZGlybmFtZSwgJy4vc3JjL2xvY2FsZXMvbGFuZyoqJyldXHJcbiAgICAgIH0pXHJcbiAgICAgIC8vIGNvcHlJbmRleCgpLFxyXG4gICAgXSxcclxuICAgIHJlc29sdmU6IHtcclxuICAgICAgYWxpYXM6IFtcclxuICAgICAgICB7XHJcbiAgICAgICAgICBmaW5kOiAndnVlLWkxOG4nLFxyXG4gICAgICAgICAgcmVwbGFjZW1lbnQ6ICd2dWUtaTE4bi9kaXN0L3Z1ZS1pMThuLmNqcy5qcydcclxuICAgICAgICB9LFxyXG4gICAgICAgIHtcclxuICAgICAgICAgIGZpbmQ6IC9AXFwvLyxcclxuICAgICAgICAgIHJlcGxhY2VtZW50OiBgJHtwYXRoUmVzb2x2ZSgnc3JjJyl9L2BcclxuICAgICAgICB9XHJcbiAgICAgIF1cclxuICAgIH0sXHJcbiAgICBiYXNlOiBiYXNlLFxyXG4gICAgZXhwZXJpbWVudGFsOiB7XHJcbiAgICAgIC8vIFx1OEZEQlx1OTYzNlx1NTdGQVx1Nzg0MFx1OERFRlx1NUY4NFx1OTAwOVx1OTg3OVxyXG4gICAgfSxcclxuICAgIHNlcnZlcjoge1xyXG4gICAgICBob3N0OiAnMC4wLjAuMCcsXHJcbiAgICAgIHByb3h5OiBnZXRQcm94eUNvbmZpZyhvcHQpLFxyXG4gICAgICBwb3J0OiA1NTIwLFxyXG4gICAgICBvcGVuOiB0cnVlLFxyXG4gICAgfSxcclxuICAgIGVzYnVpbGQ6IHtcclxuICAgICAgcHVyZTogZW52LlZJVEVfRFJPUF9DT05TT0xFID09PSAndHJ1ZScgPyBbJ2NvbnNvbGUubG9nJ10gOiB1bmRlZmluZWQsXHJcbiAgICAgIGRyb3A6IGVudi5WSVRFX0RST1BfREVCVUdHRVIgPT09ICd0cnVlJyA/IFsnZGVidWdnZXInXSA6IHVuZGVmaW5lZFxyXG4gICAgfSxcclxuICAgIGJ1aWxkOiB7XHJcbiAgICAgIG91dERpcjogZmlsZVVSTFRvUGF0aChuZXcgVVJMKCcuLi8uLi9zdGF0aWMvY2hhdC1haS1hZG1pbicsIGltcG9ydC5tZXRhLnVybCkpLFxyXG4gICAgICBlbXB0eU91dERpcjogdHJ1ZSxcclxuICAgICAgYXNzZXRzRGlyOiAnYXNzZXRzJyxcclxuICAgICAgc291cmNlbWFwOiBlbnYuVklURV9TT1VSQ0VNQVAgPT09ICd0cnVlJyxcclxuICAgICAgcm9sbHVwT3B0aW9uczoge1xyXG4gICAgICAgIC8vIGV4dGVybmFsOiBbJ21vbWVudCcsICd2aWRlby5qcycsICdqc3BkZicsICd4bHN4J10sXHJcbiAgICAgICAgZXh0ZXJuYWw6IFtdLFxyXG4gICAgICAgIHBsdWdpbnM6IFtnbG9iYWxzLCBlbnYuVklURV9VU0VfQlVORExFX0FOQUxZWkVSID09PSAndHJ1ZScgPyB2aXN1YWxpemVyKCkgOiB1bmRlZmluZWRdLFxyXG4gICAgICAgIG91dHB1dDoge1xyXG4gICAgICAgICAgLy8gXHU4MUVBXHU1QjlBXHU0RTQ5Y2h1bmtGaWxlTmFtZVx1NzUxRlx1NjIxMFx1ODlDNFx1NTIxOVxyXG4gICAgICAgICAgY2h1bmtGaWxlTmFtZXM6ICdhc3NldHMvanMvW25hbWVdLVtoYXNoXS5qcycsXHJcbiAgICAgICAgICBlbnRyeUZpbGVOYW1lczogJ1tuYW1lXS1baGFzaF0uanMnLFxyXG4gICAgICAgICAgYXNzZXRGaWxlTmFtZXMoYXNzZXRJbmZvKSB7XHJcbiAgICAgICAgICAgIGxldCBmaWVsX25hbWUgPSBhc3NldEluZm8ubmFtZS50b0xvd2VyQ2FzZSgpXHJcblxyXG4gICAgICAgICAgICBpZiAoZmllbF9uYW1lLmVuZHNXaXRoKCcuY3NzJykpIHtcclxuICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9jc3MvW25hbWVdLVtoYXNoXS5bZXh0XSdcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBpZiAoWydwbmcnLCAnanBnJywgJ2pwZWcnLCAnc3ZnJ10uc29tZSgoZXh0KSA9PiBmaWVsX25hbWUuZW5kc1dpdGgoZXh0KSkpIHtcclxuICAgICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9pbWcvW25hbWVdLVtoYXNoXS5bZXh0XSdcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBpZiAoWyd0dGYnLCAnd29mZicsICd3b2ZmMiddLnNvbWUoKGV4dCkgPT4gZmllbF9uYW1lLmVuZHNXaXRoKGV4dCkpKSB7XHJcbiAgICAgICAgICAgICAgcmV0dXJuICdhc3NldHMvZm9udHMvW25hbWVdLVtoYXNoXS5bZXh0XSdcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICByZXR1cm4gJ2Fzc2V0cy9bbmFtZV0tW2hhc2hdLltleHRdJ1xyXG4gICAgICAgICAgfSxcclxuICAgICAgICAgIC8vIFx1OEJFNVx1OTAwOVx1OTg3OVx1NTE0MVx1OEJCOFx1NEY2MFx1NTIxQlx1NUVGQVx1ODFFQVx1NUI5QVx1NEU0OVx1NzY4NFx1NTE2Q1x1NTE3MSBjaHVua1xyXG4gICAgICAgICAgbWFudWFsQ2h1bmtzOiB7XHJcbiAgICAgICAgICAgICd2dWUtY2h1bmtzJzogWyd2dWUnLCAndnVlLXJvdXRlcicsICdwaW5pYScsICd2dWUtaTE4biddLFxyXG4gICAgICAgICAgICBkYXlqczogWydkYXlqcyddLFxyXG4gICAgICAgICAgICBheGlvczogWydheGlvcyddLFxyXG4gICAgICAgICAgICAnY3J5cHRvLWpzJzogWydjcnlwdG8tanMnXSxcclxuICAgICAgICAgICAgcXM6IFsncXMnXSxcclxuICAgICAgICAgICAgJ3Z1ZS1wZGYtZW1iZWQnOiBbJ3Z1ZS1wZGYtZW1iZWQnXVxyXG4gICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgICAgfSxcclxuICAgICAgY3NzQ29kZVNwbGl0OiAhKGVudi5WSVRFX1VTRV9DU1NfU1BMSVQgPT09ICdmYWxzZScpXHJcbiAgICB9LFxyXG4gICAgb3B0aW1pemVEZXBzOiB7XHJcbiAgICAgIGluY2x1ZGU6IFsndnVlJywgJ3Z1ZS1yb3V0ZXInLCAncGluaWEnLCAndnVlLWkxOG4nLCAnZGF5anMnLCAnYXhpb3MnLCAnY3J5cHRvLWpzJywgJ3FzJ11cclxuICAgIH1cclxuICB9XHJcbn0pXHJcbiIsICJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiRDpcXFxccHJvamVjdFxcXFx6aGltYV9jaGF0d2lraVxcXFxmcm9udC1lbmRcXFxcY2hhdC1haS1hZG1pbi12dWVcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIkQ6XFxcXHByb2plY3RcXFxcemhpbWFfY2hhdHdpa2lcXFxcZnJvbnQtZW5kXFxcXGNoYXQtYWktYWRtaW4tdnVlXFxcXHByb3h5X2NvbmZpZy5qc1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vRDovcHJvamVjdC96aGltYV9jaGF0d2lraS9mcm9udC1lbmQvY2hhdC1haS1hZG1pbi12dWUvcHJveHlfY29uZmlnLmpzXCI7LyogZXNsaW50LWRpc2FibGUgbm8tdW5kZWYgKi9cclxuaW1wb3J0IHsgbG9hZEVudiB9IGZyb20gJ3ZpdGUnXHJcblxyXG5leHBvcnQgY29uc3QgZ2V0UHJveHlDb25maWcgPSAob3B0KSA9PiB7XHJcbiAgY29uc3QgeyBtb2RlIH0gPSBvcHRcclxuICBjb25zdCBlbnYgPSBsb2FkRW52KG1vZGUsIHByb2Nlc3MuY3dkKCksICcnKVxyXG5cclxuICBsZXQgcHJveHlBcGlzID0gWycvb3Blbi8nLCAnL3N0YXRpYycsICcvY29tbW9uJywgJy9tYW5hZ2UnLCAnL2FwcCcsICcvY2hhdCcsICcvdXBsb2FkJywgJy9wdWJsaWMnXVxyXG4gIGxldCBwcm94eSA9IHt9XHJcblxyXG4gIGNvbnNvbGUubG9nKGVudi5QUk9YWV9CQVNFX0FQSV9VUkwpXHJcblxyXG4gIHByb3h5QXBpcy5mb3JFYWNoKChrZXkpID0+IHtcclxuICAgIHByb3h5W2tleV0gPSB7XHJcbiAgICAgIHRhcmdldDogZW52LlBST1hZX0JBU0VfQVBJX1VSTCxcclxuICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlXHJcbiAgICB9XHJcbiAgfSlcclxuXHJcbiAgcmV0dXJuIHtcclxuICAgIC4uLnByb3h5XHJcbiAgfVxyXG59XHJcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFDQSxTQUFTLGVBQWUsV0FBVztBQUNuQyxPQUFPLFVBQVU7QUFFakIsU0FBUyxjQUFjLFdBQUFBLGdCQUFlO0FBQ3RDLE9BQU8sU0FBUztBQUNoQixPQUFPLFlBQVk7QUFFbkIsU0FBUyxrQkFBa0I7QUFDM0IsT0FBTyxxQkFBcUI7QUFDNUIsT0FBTyxnQkFBZ0I7QUFDdkIsU0FBUyw0QkFBNEI7QUFDckMsT0FBTyxtQkFBbUI7QUFDMUIsU0FBUyw0QkFBNEI7OztBQ1py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
