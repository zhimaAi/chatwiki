import axios from 'axios'
import {
  defaultRequestInterceptors,
  defaultResponseInterceptors,
  responseInterceptorsCatch
} from './config'

import { REQUEST_TIMEOUT } from '@/constants'

export const PATH_URL = import.meta.env.VITE_BASE_API_URL

const pending = new Map()
const CancelToken = axios.CancelToken

const removePending = (config) => {
  const key = `${config.url}&${config.method}`
  if (pending.has(key)) {
    const cancel = pending.get(key)
    cancel(key)
    pending.delete(key)
  }
}

const axiosInstance = axios.create({
  timeout: REQUEST_TIMEOUT,
  baseURL: PATH_URL
})

axiosInstance.interceptors.request.use((config) => {
  if (config.cancelToken) {
    removePending(config) // 取消之前的请求
    config.cancelToken = new CancelToken((c) => {
      const key = `${config.url}&${config.method}`
      pending.set(key, c)
    })
  }
  return config
})

axiosInstance.interceptors.response.use(
  (res) => {
    // 在响应成功后，从 pending 中移除请求
    removePending(res.config)
    // 这里不能做任何处理，否则后面的 interceptors 拿不到完整的上下文了
    return res
  },
  (error) => {
    if (axios.isCancel(error)) {
      console.log('Request canceled', error.message)
    } else {
      // 在响应失败后，也从 pending 中移除请求
      if (error.config) {
        removePending(error.config)
      }
    }
    return Promise.reject(error)
  }
)

axiosInstance.interceptors.request.use(defaultRequestInterceptors)
axiosInstance.interceptors.response.use(defaultResponseInterceptors)
axiosInstance.interceptors.response.use(undefined, responseInterceptorsCatch)

const service = {
  request: (config) => {
    return new Promise((resolve, reject) => {
      if (config.interceptors?.requestInterceptors) {
        config = config.interceptors.requestInterceptors(config)
      }

      axiosInstance
        .request(config)
        .then((res) => {
          resolve(res)
        })
        .catch((err) => {
          reject(err)
        })
    })
  },
  cancelRequest: (url) => {
    for (const [key, cancel] of pending) {
      const urlList = Array.isArray(url) ? url : [url]
      if (urlList.some((_url) => key.startsWith(_url))) {
        cancel(`Request canceled: ${key}`)
        pending.delete(key)
      }
    }
  },
  cancelAllRequest() {
    for (const [key, cancel] of pending) {
      cancel(`All requests canceled`)
      pending.delete(key)
    }
  }
}

export default service
