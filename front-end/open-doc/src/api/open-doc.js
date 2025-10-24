import request from '@/utils/http/axios'
import SSE from '@/utils/http/sse'

const baseURL = import.meta.env.VITE_BASE_API_URL

// 获取默认库的接口
export const getBindLibList = ({ library_key }) => {
  return request.get({ url: '/open/doc/bindLibList/api?library_key=' + library_key })
}

export const getAiMessage = (data) => {
  return new SSE({
    method: 'GET',
    url: baseURL + '/open/summary/' + data.id + '?v=' + data.v,
  })
}

export const getOpenDoc = (params) => {
  return request.get({ url: '/open/doc/api/' + params.id })
}

export const getOpenHome = (params) => {
  return request.get({ url: '/open/home/api/' + params.id })
}

export const getOpenSearch = (params) => {
  return request.get({
    url: '/open/search/api/' + params.id,
    params: { v: params.v },
  })
}

export const getSearchResult = (params) => {
  return request.get({
    url: '/open/search/query/' + params.id,
    params: { v: params.v },
  })
}
