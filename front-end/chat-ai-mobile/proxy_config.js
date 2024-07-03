import { loadEnv } from 'vite';
export var getProxyConfig = function (opt) {
    var mode = opt.mode;
    var env = loadEnv(mode, process.cwd(), '');
    var proxyApis = ['/static', '/common', '/manage', '/app', '/chat', '/upload'];
    var proxy = {};
    console.log('PROXY_BASE_API_URL: ' + env.PROXY_BASE_API_URL);
    proxyApis.forEach(function (key) {
        proxy[key] = {
            target: env.PROXY_BASE_API_URL,
            changeOrigin: true,
        };
    });
    return proxy;
};
