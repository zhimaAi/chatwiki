import { loadEnv } from 'vite';

// 假设环境变量和配置文件格式已经定义好
interface ProxyConfig {
  [key: string]: {
    target: string;
    changeOrigin: boolean;
  };
}

export const getProxyConfig = (opt: { mode: string }): ProxyConfig => {
  const { mode } = opt;
    const env = loadEnv(mode, process.cwd(), '');

    const proxyApis = ['/static', '/common', '/manage', '/app', '/chat', '/upload'];
    const proxy: ProxyConfig = {};

    console.log('PROXY_BASE_API_URL: ' + env.PROXY_BASE_API_URL);

    proxyApis.forEach((key) => {
      proxy[key] = {
        target: env.PROXY_BASE_API_URL,
        changeOrigin: true,
      };
    });

    return proxy;
};
