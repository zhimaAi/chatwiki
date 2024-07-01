import request from '@/utils/http/axios'

export const getCompany = () => {
  return request.get({
    url: '/manage/clientSide/getCompany',
    params: {}
  })
}

export const login = (user_name, password) => {
  return request.post({
    withToken: true,
    notDisplayErrorTips: true,
    url: 'manage/clientSide/login',
    data: { user_name: user_name, password }
  })
}
