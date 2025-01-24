import request from '@/utils/http/axios'
import { objectToQueryString } from '@/utils/index.js'
import { useUserStore } from '@/stores/modules/user'
const userStore = useUserStore()
export const getFormList = (params = {}) => {
  return request.get({
    url: '/manage/getFormList',
    params: params
  })
}

export const getFormInfo = (params = {}) => {
  return request.get({
    url: '/manage/getFormInfo',
    params: params
  })
}

export const addForm = (data) => {
  return request.post({
    url: '/manage/addForm',
    data: data
  })
}

export const editForm = (data) => {
  return request.post({
    url: '/manage/editForm',
    data: data
  })
}


export const delForm = (data) => {
  return request.post({
    url: '/manage/delForm',
    data: data
  })
}


export const getFormFieldList = (params = {}) => {
  return request.get({
    url: '/manage/getFormFieldList',
    params: params
  })
}

export const addFormField = (data) => {
  let url = data.id ? '/manage/editFormField' : '/manage/addFormField'
  return request.post({
    url,
    data,
  })
}

export const delFormField = (data) => {
  return request.post({
    url: '/manage/delFormField',
    data: data
  })
}

export const updateFormRequired = (data) => {
  return request.post({
    url: '/manage/updateFormRequired',
    data: data
  })
}

export const getFormEntryList = (params = {}) => {
  return request.get({
    url: '/manage/getFormEntryList',
    params: params
  })
}

export const delFormEntry = (data) => {
  return request.post({
    url: '/manage/delFormEntry',
    data: data
  })
}

export const emptyFormEntry = (data) => {
  return request.post({
    url: '/manage/emptyFormEntry',
    data: data
  })
}

export const addFormEntry = (data) => {
  let url = data.id ? '/manage/editFormEntry' : '/manage/addFormEntry'
  return request.post({
    url,
    data,
  })
}

export const exportFormEntry = (params = {}) => {
  window.open(`/manage/exportFormEntry?${objectToQueryString(params)}` + '&token=' + userStore.getToken)
}

export const getFormFilterList = (params = {}) => {
  return request.get({
    url: '/manage/getFormFilterList',
    params: params
  })
}

export const getFormFilterInfo = (params = {}) => {
  return request.get({
    url: '/manage/getFormFilterInfo',
    params: params
  })
}

export const addFormFilter = (data) => {
  return request.post({
    url: '/manage/addFormFilter',
    data: data
  })
}
export const editFormFilter = (data) => {
  return request.post({
    url: '/manage/editFormFilter',
    data: data
  })
}
export const delFormFilter = (data) => {
  return request.post({
    url: '/manage/delFormFilter',
    data: data
  })
}

export const updateFormFilterEnabled = (data) => {
  return request.post({
    url: '/manage/updateFormFilterEnabled',
    data: data
  })
}

export const updateFormFilterSort = (data) => {
  return request.post({
    url: '/manage/updateFormFilterSort',
    data: data
  })
}

export const uploadFormFile = (data) => {
  return request.post({
    url: '/manage/uploadFormFile',
    data: data
  })
}


export const getUploadFormFileProc = (data) => {
  return request.post({
    url: '/manage/getUploadFormFileProc',
    data: data
  })
}