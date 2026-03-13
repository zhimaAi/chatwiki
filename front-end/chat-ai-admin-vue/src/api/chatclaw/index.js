import request from '@/utils/http/axios'

export const chatClawLoginApi = (user_name, password, clientInfo = {}) => {
  return request.post({
    withToken: true,
    url: '/manage/chatclaw/login',
    data: { user_name, password, ...clientInfo }
  })
}

export const getChatClawTokenLogListApi = (params = {}) => {
  return request.get({
    url: '/manage/chatclaw/tokenLogList',
    params: { page: params.page ?? 1, size: params.size ?? 10 }
  })
}

export const forceOfflineChatClawTokenApi = (data = {}) => {
  return request.post({
    url: '/manage/chatclaw/tokenForceOffline',
    data
  })
}
