import request from '@/utils/http/axios'

// 获取知识库权限
export const getLibDocPartner = (data) => {
  return request.get({
    url: '/manage/getLibDocPartner',
    params: data
  })
}

// 获取目录
export const getDocList = (data) => {
  return request.get({
    url: '/manage/getCatalog',
    params: data
  })
}

export const saveLibDoc = (data) => {
  return request.post({
    url: '/manage/saveLibDoc',
    data: data
  })
}

export const saveDraftLibDoc = (data) => {
  return request.post({
    url: '/manage/draftLibDoc',
    data: data
  })
}

export const deleteLibDoc = (data) => {
  return request.post({
    url: '/manage/deleteLibDoc',
    data: data
  })
}

export const getLibDocInfo = (data) => {
  return request.get({
    url: '/manage/getLibDocInfo',
    params: data
  })
}

export const updateDocSort = (data) => {
  return request.post({
    url: '/manage/changeLibDoc',
    data: data
  })
}
// 保存问题引导
export const saveQuestionGuide = (data) => {
  return request.post({
    url: '/manage/saveQuestionGuide',
    data: data
  })
}
// 删除问题引导
export const deleteQuestionGuide = (data) => {
  return request.post({
    url: '/manage/deleteQuestionGuide',
    data: data
  })
}
// 保存协作者
export const saveLibDocPartner = (data) => {
  return request.post({
    url: '/manage/saveLibDocPartner',
    data: data
  })
}

// 获取协作者
export const getLibDocPartnerList = (data) => {
  return request.get({
    url: '/manage/libDocPartnerList',
    params: data
  })
}

// 删除协作者
export const deleteLibDocPartner = (data) => {
  return request.post({
    url: '/manage/deleteLibDocPartner',
    data: data
  })
}

// 保存seo
export const saveLibDocSeo = (data) => {
  return request.post({
    url: '/manage/saveLibDocSeo',
    data: data
  })
}
