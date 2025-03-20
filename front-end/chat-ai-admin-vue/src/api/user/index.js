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

export const checkLogin = (data) => {
  return request.get({ url: '/manage/checkLogin', data })
}

// 刷新用户token
export const refreshUserToken = (data) => {
  return request.post({ url: '/manage/refreshUserToken', data })
}

export const getKefuNav = (data) => {
  return request.post({ url: '/manage/getKefuNav', data })
}

export const getUnReadMessageTotal = (data) => {
  return request.post({ url: '/message/getUnReadMessageTotalCount', data })
}

export const loginOutApi = () => {
  return request.get({ url: '/mock/user/loginOut' })
}

export const getUserListApi = ({ params }) => {
  return request.get({ url: '/mock/user/list', params })
}

export const getAdminRoleApi = (params) => {
  return request.get({ url: '/mock/role/list', params })
}

export const getTestRoleApi = (params) => {
  return request.get({ url: '/mock/role/list2', params })
}

export const getCompany = () => {
  return request.get({ url: '/manage/getCompany' })
}

export const saveCompany = (data) => {
  return request.post({ url: '/manage/saveCompany', data })
}

export const getClientSideLoginSwitch = () => {
  return request.get({ url: '/manage/getClientSideLoginSwitch' })
}

export const setClientSideLoginSwitch = ({ client_side_login_switch }) => {
  return request.post({
    url: '/manage/setClientSideLoginSwitch',
    data: { client_side_login_switch }
  })
}

export const clientSideDownload = ({ domain }) => {
  return request.post({
    url: '/manage/clientSideDownload',
    data: { domain }
  })
}

// 获取自定义域名
export const getDomainList = (data) => {
    return request.get({ url: '/manage/diyDomainList', data })
}

// 保存自定义域名
export const saveDomain = (data) => {
    return request.post({ url: '/manage/saveDiyDomain', data })
}

// 删除自定义域名
export const deleteDomain = (data) => {
    return request.post({ url: '/manage/deleteDiyDomain', data })
}

// 上传证书
export const uploadCertificate = (data) => {
    return request.post({ url: '/manage/uploadCertificate', data })
}

// 上传验证文件
export const uploadCheckFile = (data) => {
    return request.post({ url: '/manage/uploadCheckFile', data })
}