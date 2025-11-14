import request from "@/utils/http/axios/index.js";

export function getTMcpProviders(params = {}) {
  return request.get({url: '/manage/getMcpProviderList', params})
}

export function getTMcpProviderInfo(params = {}) {
  return request.get({url: '/manage/getMcpProviderDetail', params})
}

export function saveTMcpProvider(data = {}) {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/saveMcpProvider',
    data
  })
}

export function authTMcpProvider(data = {}) {
  return request.post({url: '/manage/authMcpProvider', data})
}

export function cancelAuthTMcpProvider(data = {}) {
  return request.post({url: '/manage/cancelAuthMcpProvider', data})
}

export function delTMcpProvider(data = {}) {
  return request.post({url: '/manage/deleteMcpProvider', data})
}

