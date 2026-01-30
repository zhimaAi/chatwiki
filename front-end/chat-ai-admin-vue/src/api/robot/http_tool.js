import request from '@/utils/http/axios'

export function getHttpTools(params = {}) {
  return request.get({url: '/manage/httpTool/getHttpToolList', params})
}

export const saveHttpTool = (data = {}) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/httpTool/saveHttpTool',
    data
  })
}

export const delHttpTool = (data = {}) => {
  return request.post({
    url: '/manage/httpTool/deleteHttpTool',
    data
  })
}

export const saveHttpToolItem = (data = {}) => {
  return request.post({
    url: '/manage/httpTool/saveHttpToolNode',
    data
  })
}

export const delHttpToolItem = (data = {}) => {
  return request.post({
    url: '/manage/httpTool/deleteHttpToolNode',
    data
  })
}
