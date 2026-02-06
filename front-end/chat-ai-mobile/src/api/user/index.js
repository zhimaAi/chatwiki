import request from '@/utils/http/axios'

export const loginApi = (user_name, password) => {
  return request.post({
    withToken: true,
    url: '/manage/login',
    data: { user_name: user_name, password }
  })
}

export const getUserInfo = (data) => {
  return request.get({ url: '/manage/checkLogin', data })
}

export const getCookieTip = (params) => {
  return request.get({ url: '/manage/getCookieTip', params })
}