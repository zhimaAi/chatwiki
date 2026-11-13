import request from '@/utils/http/axios'

export const getOptimizeContext = (params) =>
  request.get({ url: '/manage/getPromptOptimizeContext', params })
export const getOptimizeRecords = (params) =>
  request.get({ url: '/manage/getPromptOptimizeRecords', params })
export const getOptimizeVersions = (params) =>
  request.get({ url: '/manage/getPromptOptimizeVersions', params })
export const getOptimizeVersion = (params) =>
  request.get({ url: '/manage/getPromptOptimizeVersion', params })
