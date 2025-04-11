import request from '@/utils/http/axios'

export const getReceiverList = (params = {}) => {
  return request.get({
    url: '/manage/getReceiverList',
    params: params
  })
}

export const getChatMessage = (data = {}) => {
  return request.post({
    url: '/chat/message',
    data: data
  })
}

export const setReceiverRead = (data = {}) => {
  return request.post({
    url: '/manage/setReceiverRead',
    data: data
  })
}
