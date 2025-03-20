/* eslint-disable no-undef */
import { loadEnv } from 'vite'

export const getProxyConfig = (opt) => {
  const { mode } = opt
  const env = loadEnv(mode, process.cwd(), '')

  let proxyApis = ['/open/', '/static', '/common', '/manage', '/app', '/chat', '/upload', '/public']
  let proxy = {}

  console.log(env.PROXY_BASE_API_URL)

  proxyApis.forEach((key) => {
    proxy[key] = {
      target: env.PROXY_BASE_API_URL,
      changeOrigin: true
    }
  })

  return {
    ...proxy
  }
}
