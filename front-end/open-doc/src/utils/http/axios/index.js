import service from './service'
import { CONTENT_TYPE } from '@/constants'

const request = (option) => {
  const { url, method, params, data, headers, responseType, withToken } = option

  const defaultHeaders = {
    'Content-Type': CONTENT_TYPE,
    'X-Requested-With': 'XMLHttpRequest',
    'App-Type': '',
    ...headers,
  }

  if (!withToken) {
    // console.log('without token')
  }

  return service.request({
    url: url,
    method,
    params,
    data: data,
    responseType: responseType,
    headers: {
      ...defaultHeaders,
    },
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
  },
}
