import SSE from '@/utils/http/sse'
import request from '@/utils/http/axios'

const baseURL = import.meta.env.VITE_BASE_API_URL

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

// 猜你想问
export const questionGuide = (data) => {
  return request.post({ url: '/chat/questionGuide', data })
}