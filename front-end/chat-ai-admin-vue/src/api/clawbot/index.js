import request from '@/utils/http/axios'

export const saveClawbotConf = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/saveClawbotConf',
    data
  })
}

export const uploadClawbotLocalDoc = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/uploadClawbotLocalDoc',
    data
  })
}

export const getClawbotLocalDocList = (params = {}) => {
  return request.get({
    url: '/manage/getClawbotLocalDocList',
    params
  })
}

export const deleteClawbotLocalDoc = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/deleteClawbotLocalDoc',
    data
  })
}
