import request from '@/utils/http/axios'

export const saveLibrarySearch = (data) => {
  return request.post({
    url: '/manage/saveLibrarySearch',
    data: data
  })
}

export const getAllDepartment = (params = {}) => {
  return request.get({
    url: '/manage/getAllDepartment',
    params: params
  })
}

export const getDepartmentList = (params = {}) => {
  return request.get({
    url: '/manage/getDepartmentList',
    params: params
  })
}

export const getPermissionManageList = (params = {}) => {
  return request.get({
    url: '/manage/getPermissionManageList',
    params: params
  })
}

export const getPartnerManageList = (params = {}) => {
  return request.get({
    url: '/manage/getPartnerManageList',
    params: params
  })
}

export const batchSavePermissionManage = (data) => {
  return request.post({
    url: '/manage/batchSavePermissionManage',
    data: data
  })
}

export const savePermissionManage = (data) => {
  return request.post({
    url: '/manage/savePermissionManage',
    data: data
  })
}

export const saveDepartment = (data) => {
  return request.post({
    url: '/manage/saveDepartment',
    data: data
  })
}

export const deletePermissionManage = (data) => {
  return request.post({
    url: '/manage/deletePermissionManage',
    data: data
  })
}


export const deleteDepartment = (data) => {
  return request.post({
    url: '/manage/deleteDepartment',
    data: data
  })
}


export const batchUpdateUserDepartment = (data) => {
  return request.post({
    url: '/manage/batchUpdateUserDepartment',
    data: data
  })
}
