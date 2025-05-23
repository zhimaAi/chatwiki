import request from '@/utils/http/axios'

export const checkPermission = (params = {}) => {
  // 获取团队成员列表
  return request.get({
    url: '/manage/checkPermission',
    params: params
  })
}

export const getUserList = (params = {}) => {
  // 获取团队成员列表
  return request.get({
    url: '/manage/getUserList',
    params: params
  })
}

export const getTokenStats = (params = {}) => {
  // 获取模型的Token使用
  return request.get({
    url: '/manage/stats/token',
    params: params
  })
}

export const getFeedbackList = (params = {}) => {
  // 后台反馈列表
  return request.get({
    url: '/manage/feedback/list',
    params: params
  })
}

export const getFeedbackDetail = (params = {}) => {
  // 后台反馈详情
  return request.get({
    url: '/manage/feedback/detail',
    params: params
  })
}

export const getFeedbackStats = (params = {}) => {
  // 后台反馈统计
  return request.get({
    url: '/manage/feedback/stats',
    params: params
  })
}

export const getAnalyse = (params = {}) => {
  // 获取统计分析
  return request.get({
    url: '/manage/stats/analyse',
    params: params
  })
}

export const saveUser = (data = {}) => {
  // 添加成员
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/saveUser',
    data: data
  })
}

export const resetPass = (data = {}) => {
  // 重置密码
  return request.post({
    url: '/manage/resetPass',
    data: data
  })
}

export const delUser = (data = {}) => {
  // 删除用户
  return request.post({
    url: '/manage/delUser',
    data: data
  })
}


export const getRoleList = (params = {}) => {
  // 获取团队成员列表
  return request.get({
    url: '/manage/getRoleList',
    params: params
  })
}

export const getUser = (params = {}) => {
  // 获取团队成员列表
  return request.get({
    url: '/manage/getUser',
    params: params
  })
}

export const getRole = (params = {}) => {
  // 获取角色信息
  return request.get({
    url: '/manage/getRole',
    params: params
  })
}

export const saveUserManagedDataList = (data = {}) => {
  // 保存用户管理的机器人、知识库、数据库
  return request.post({
    url: '/manage/saveUserManagedDataList',
    data: data
  })
}

export const saveRole = (data = {}) => {
  // 保存角色
  return request.post({
    url: '/manage/saveRole',
    data: data
  })
}
export const delRole = (data = {}) => {
  // 删除角色
  return request.post({
    url: '/manage/delRole',
    data: data
  })
}
export const getMenu = (params = {}) => {
  // 获取菜单
  return request.get({
    url: '/manage/getMenu',
    params: params
  })
}
export const saveMenu = (data = {}) => {
  // 保存菜单
  return request.post({
    url: '/manage/saveMenu',
    data: data
  })
}

export const delMenu = (data = {}) => {
  // 删除菜单
  return request.post({
    url: '/manage/delMenu',
    data: data
  })
}

export const saveProfile = (data = {}) => {
  // 修改账户信息
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/saveProfile',
    data: data
  })
}

export const loginSwitch = (data = {}) => {
  // 登录开关
  return request.post({
    url: '/manage/loginSwitch',
    data: data
  })
}