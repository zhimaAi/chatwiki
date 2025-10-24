import qs from 'qs'
import { SUCCESS_CODE, TRANSFORM_REQUEST_DATA } from '@/constants'
import { message } from 'ant-design-vue'
import { objToFormData } from '@/utils'
import { getErrorMsg } from './errorMsg'
import { useOpenDocStore } from '@/stores/open-doc'
import router from '@/router'

function showErrorMsg(msg) {
  message.destroy()
  message.error(msg)
}

// 全局错误提示
let isGlobalError = false
function showGlobalErrorMsg(msg) {
  if (isGlobalError) {
    return
  }

  isGlobalError = true

  showErrorMsg(msg)

  setTimeout(() => {
    isGlobalError = false
  }, 200)
}

const defaultRequestInterceptors = (config) => {
  const store = useOpenDocStore()
  if (config.method === 'get') {
    if (!config.params) {
      config.params = {}
    }
    if (store.token.length) {
      config.params.token = store.token
    }
    if (store.previewKey) {
      config.params.preview = store.previewKey
    }
  } else {
    if (!config.data) {
      config.data = {}
    }
    if (store.token.length) {
      config.data.token = store.token
    }
    if (store.previewKey) {
      config.data.preview = store.previewKey
    }
  }

  if (
    config.method === 'post' &&
    config.headers['Content-Type'] === 'application/x-www-form-urlencoded'
  ) {
    config.data = qs.stringify(config.data)
  } else if (
    TRANSFORM_REQUEST_DATA &&
    config.method === 'post' &&
    config.headers['Content-Type'] === 'multipart/form-data'
  ) {
    config.data = objToFormData(config.data)
  }

  if (config.method === 'get' && config.params) {
    let url = config.url
    url += '?'

    const keys = Object.keys(config.params)

    for (const key of keys) {
      if (config.params[key] !== void 0 && config.params[key] !== null) {
        url += `${key}=${encodeURIComponent(config.params[key])}&`
      }
    }
    url = url.substring(0, url.length - 1)
    config.params = {}
    config.url = url
  }

  return config
}

const defaultResponseInterceptors = (response) => {
  if (response.data && typeof response.data.res != 'undefined') {
    response.data.code = response.data.res
    response.data.message = response.data.msg || ''
  }

  if (response?.config?.responseType === 'blob') {
    // 如果是文件流，直接过
    return response
  } else if (response.data.code === SUCCESS_CODE) {
    let res = response.data.data
    if (res && res.is_404 === 1) {
      router.replace('/404')
      return Promise.reject({
        code: 1,
        msg: '页面不存在',
      })
    }

    return response.data
  } else {
    showErrorMsg(response?.data?.message)

    if (response?.data?.code === 401) {
      console.log('is 401')
    } else {
      return Promise.reject(response.data)
    }
  }
}

const responseInterceptorsCatch = (error) => {
  if (error.response) {
    console.log('response err： ' + error)
    return responseError(error)
  } else if (error.request) {
    console.log('request err： ' + error)
    return requestError(error)
  }

  return Promise.reject(error)
}

const requestError = (error) => {
  const { response, code, message } = error || {}
  const msg = response?.data?.error?.message ?? ''
  const err = error?.toString?.() ?? ''

  let errMessage = ''

  if (code === 'ECONNABORTED' && message.indexOf('timeout') !== -1) {
    errMessage = '接口请求超时,请刷新页面重试!'
  } else if (err?.includes('Network Error')) {
    errMessage = '网络异常，请检查您的网络连接是否正常!'
  } else {
    errMessage = msg
  }

  showGlobalErrorMsg(errMessage)

  return Promise.reject(error)
}

const responseError = (error) => {
  let status = error?.response?.status
  let msg = getErrorMsg(error)

  showGlobalErrorMsg(msg)

  if (status === 401) {
    console.log('is 401')
  }

  return Promise.reject(error)
}

export { defaultResponseInterceptors, defaultRequestInterceptors, responseInterceptorsCatch }
