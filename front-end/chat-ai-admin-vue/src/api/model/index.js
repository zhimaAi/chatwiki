import request from '@/utils/http/axios'

export const getModelConfigList = () => {
  return request.get({ url: '/manage/getModelConfigList' })
}

export const addModelConfig = (data) => {
  return request.post({ url: '/manage/addModelConfig', data })
}

export const editModelConfig = (data) => {
  return request.post({ url: '/manage/editModelConfig', data })
}

export const delModelConfig = (data) => {
  return request.post({ url: '/manage/delModelConfig', data })
}

export const getModelConfigOption = (params) => {
  return request.get({ url: '/manage/getModelConfigOption', params })
}
