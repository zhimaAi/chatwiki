import qs from 'qs'
import { SUCCESS_CODE, TRANSFORM_REQUEST_DATA } from '@/constants'
import { useUserStore } from '@/stores/modules/user'
import { objToFormData } from '@/utils'
import { showToast } from 'vant'
import { useI18n } from '@/hooks/web/useI18n'
import { getErrorMsg } from './errorMsg'

// 全局错误提示
let isGlobalError = false
function showGlobalErrorMsg(msg) {
  if (isGlobalError) {
    return
  }

  isGlobalError = true

  showToast(msg)

  setTimeout(() => {
    isGlobalError = false
  }, 200)
}

const defaultRequestInterceptors = (config) => {
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
    return response.data
  } else {
    showToast(response?.data?.message)

    if (response?.data?.code === 401) {
      const userStore = useUserStore()

      userStore.logout()
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
  const { t } = useI18n()
  const { response, code, message } = error || {}
  const msg = response?.data?.error?.message ?? ''
  const err = error?.toString?.() ?? ''

  let errMessage = ''

  if (code === 'ECONNABORTED' && message.indexOf('timeout') !== -1) {
    errMessage = t('common.apiTimeoutMessage')
  } else if (err?.includes('Network Error')) {
    errMessage = t('common.networkExceptionMsg')
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
    const userStore = useUserStore()
    userStore.logout()
  }

  return Promise.reject(error)
}

export { defaultResponseInterceptors, defaultRequestInterceptors, responseInterceptorsCatch }
