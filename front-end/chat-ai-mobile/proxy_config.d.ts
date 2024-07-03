interface ProxyConfig {
    [key: string]: {
        target: string;
        changeOrigin: boolean;
    };
}
export declare const getProxyConfig: (opt: {
    mode: string;
}) => ProxyConfig;
export {};
