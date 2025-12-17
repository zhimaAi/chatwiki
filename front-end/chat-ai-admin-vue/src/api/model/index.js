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

export const getTokenModels = (params) => {
  return request.get({ url: '/manage/stats/getActiveModels', params })
}

export const getSelfModelBuylog = (params) => {
  return request.get({ url: '/manage/getSelfModelBuylog', params })
}

export const getSelfModelConfigs = (params) => {
  return request.get({ url: '/manage/getSelfModelConfigs', params })
}

export const showModelConfigList = () => {
  return request.get({ url: '/manage/showModelConfigList' })
}

export const saveUseModelConfig = (data) => {
  return request.post({ url: '/manage/saveUseModelConfig', data })
}

export const delUseModelConfig = (data) => {
  return request.post({ url: '/manage/delUseModelConfig', data })
}