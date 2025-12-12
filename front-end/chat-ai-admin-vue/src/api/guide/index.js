import request from '@/utils/http/axios'

export const useGuideProcess = (params = {}) => {
  return request.get({
    url: '/manage/useGuideProcess',
    params: params
  })
}

export const requestNotStream = (data) => {
  return request.post({
    url: '/chat/requestNotStream',
    data: data
  })
}
