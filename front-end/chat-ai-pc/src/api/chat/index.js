import SSE from '@/utils/http/sse'
import request from '@/utils/http/axios'

const baseURL = import.meta.env.VITE_BASE_API_URL

// 获取ws连接地址
export const getWsUrl = ({ openid }) => {
  return request.get({
    url: '/chat/getWsUrl',
    params: {
      openid,
      debug: import.meta.env.DEV ? 1 : 0
    }
  })
}
export const sendAiMessage = (data) => {
  return new SSE({
    url: baseURL + '/chat/request',
    data: data
  })
}

export const chatWelcome = (data) => {
  return request.post({ url: '/chat/welcome', data })
}

export const getDialogueList = (data) => {
  return request.post({ url: '/manage/getDialogueList', data })
}

export const getChatMessage = (data) => {
  return request.post({ url: '/chat/message', data })
}

// 获取答案来源
export const getAnswerSource = (params) => {
  return request.get({ url: '/manage/getAnswerSource', params })
}

export const getFastCommandList = (params) => {
  return request.get({ url: '/chat/getFastCommandList', params })
}