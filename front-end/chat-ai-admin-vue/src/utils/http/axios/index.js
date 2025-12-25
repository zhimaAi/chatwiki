import service from './service'
import { CONTENT_TYPE } from '@/constants'
import { useUserStore } from '@/stores/modules/user'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'

const request = (option) => {
  const userStore = useUserStore()
  const localeStore = useLocaleStoreWithOut()
  const {
    url,
    method,
    params,
    data,
    headers,
    responseType,
    withToken,
    onUploadProgress,
    cancelToken
  } = option

  const currentLocale = localeStore.getCurrentLocale

  const defaultHeaders = {
    'Content-Type': CONTENT_TYPE,
    'X-Requested-With': 'XMLHttpRequest',
    'App-Type': '',
    lang: currentLocale.lang,
    ...headers
  }

  if (!withToken) {
    defaultHeaders.token = userStore.getToken ?? ''
  }

  return service.request({
    url: url,
    method,
    params,
    data: data,
    responseType: responseType,
    headers: {
      ...defaultHeaders
    },
    onUploadProgress,
    cancelToken
  })
}

export default {
  get: (option) => {
    return request({ method: 'get', ...option })
  },
  post: (option) => {
    return request({ method: 'post', ...option })
  },
  delete: (option) => {
    return request({ method: 'delete', ...option })
  },
  put: (option) => {
    return request({ method: 'put', ...option })
  },
  cancelRequest: (url) => {
    return service.cancelRequest(url)
  },
  cancelAllRequest: () => {
    return service.cancelAllRequest()
  }
}
