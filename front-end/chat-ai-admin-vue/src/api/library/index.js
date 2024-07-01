import request from '@/utils/http/axios'

export const getLibraryList = (params = {}) => {
  return request.get({
    url: '/manage/getLibraryList',
    params: params
  })
}

export const deleteLibrary = ({ id }) => {
  return request.post({
    url: '/manage/deleteLibrary',
    data: {
      id
    }
  })
}

export const createLibrary = (data) => {
  return request.post({
    headers: {
      //'Content-Type': 'multipart/form-data'
    },
    url: '/manage/createLibrary',
    data: data
  })
}

export const getLibraryFileList = ({ library_id, file_name = undefined, page = 1, size = 20 }) => {
  return request.get({
    url: '/manage/getLibFileList',
    params: {
      library_id,
      file_name,
      page,
      size
    }
  })
}

export const getLibraryInfo = ({ id }) => {
  return request.get({
    url: '/manage/getLibraryInfo',
    params: {
      id
    }
  })
}

export const delLibraryFile = ({ id }) => {
  return request.post({
    url: '/manage/delLibraryFile',
    data: {
      id
    }
  })
}

export const addLibraryFile = (data) => {
  return request.post({
    url: '/manage/addLibraryFile',
    data: data
  })
}

export const editLibrary = (data) => {
  return request.post({
    url: '/manage/editLibrary',
    data: data
  })
}

// 获取分隔符列表
export const getSeparatorsList = () => {
  return request.get({
    url: '/manage/getSeparatorsList',
    params: {}
  })
}

// 获取文档拆分
export const getLibFileSplit = (params) => {
  return request.get({
    url: '/manage/getLibFileSplit',
    params: params
  })
}

export const getLibFileInfo = (params = {}) => {
  return request.get({
    url: '/manage/getLibFileInfo',
    params: params
  })
}

export const saveLibFileSplit = (data) => {
  return request.post({
    url: '/manage/saveLibFileSplit',
    data: data
  })
}

export const getParagraphList = (params = {}) => {
  return request.get({
    url: '/manage/getParagraphList',
    params: params
  })
}

export const editParagraph = (data) => {
  return request.post({
    url: '/manage/editParagraph',
    data: data
  })
}

export const deleteParagraph = (data) => {
  return request.post({
    url: '/manage/deleteParagraph',
    data: data
  })
}

export const libraryRecallTest = (data) => {
  return request.post({
    url: '/manage/libraryRecallTest',
    data: data
  })
}

export const getLibFileExcelTitle = (params) => {
  return request.get({
    url: '/manage/getLibFileExcelTitle',
    params: params
  })
}

export const editLibFile = (data) => {
  return request.post({
    url: '/manage/editLibFile',
    data: data
  })
}
