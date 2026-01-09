import request from '@/utils/http/axios'

export const getPaymentSetting = (params = {}) => {
  return request.get({url: '/manage/robotPayment/getPaymentSetting', params})
}

export const savePaymentSetting = (data = {}) => {
  return request.post({url: '/manage/robotPayment/savePaymentSetting', data})
}

export const copyPaymentSetting = (data = {}) => {
  return request.post({url: '/manage/robotPayment/copyPaymentSetting', data})
}

export const addAuthCode = (data = {}) => {
  return request.post({url: '/manage/robotPayment/addAuthCode', data})
}

export const delAuthCode = (data = {}) => {
  return request.post({url: '/manage/robotPayment/deleteAuthCode', data})
}

export const getAuthCodeList = (params = {}) => {
  return request.get({url: '/manage/robotPayment/getAuthCodeList', params})
}

export const getAuthCodeStats = (params = {}) => {
  return request.get({url: '/manage/robotPayment/getAuthCodeStats', params})
}

export const addAuthCodeManager = (data = {}) => {
  return request.post({url: '/manage/robotPayment/addAuthCodeManager', data})
}

export const getAuthCodeManager = (params = {}) => {
  return request.get({url: '/manage/robotPayment/getAuthCodeManager', params})
}

export const delAuthCodeManager = (data = {}) => {
  return request.post({url: '/manage/robotPayment/deleteAuthCodeManager', data})
}
