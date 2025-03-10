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

// 猜你想问
export const questionGuide = (data) => {
  return request.post({ url: '/chat/questionGuide', data })
}

export const getFastCommandList = (params) => {
  return request.get({ url: '/chat/getFastCommandList', params })
}

// 点赞/点踩
export const addFeedback = (data) => {
  return request.post({ url: '/chat/message/addFeedback', data })
}

// 取消点赞/点踩接口
export const delFeedback = (data) => {
  return request.post({ url: '/chat/message/delFeedback', data })
}
