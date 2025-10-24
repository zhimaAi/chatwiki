import request from '@/utils/http/axios'

export const getLibraryInfo = ({ id }) => {
  return request.get({
    url: '/manage/getLibraryInfo',
    params: {
      id
    }
  })
}

export const editLibrary = (data) => {
  return request.post({
    url: '/manage/editLibrary',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data: data
  })
}

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
// 获取文档信息
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

// 获取文档首页配置
export const getLibDocHomeConfig = (data) => {
  return request.get({
    url: '/manage/getLibDocInfo',
    params: {
      library_key: data.library_key,
      is_index: 1
    }
  })
}

// 保存seo
export const saveLibDocSeo = (data) => {
  return request.post({
    url: '/manage/saveLibDocSeo',
    data: data
  })
}

// 获取预览地址
export const getPreviewUrl = (data) => {
  return request.post({
    url: '/manage/previewLibDoc',
    data: {
      library_key: data.library_key,
      doc_id: data.doc_id
    }
  })
}

// 保存快捷方式
export const saveLibDocIndexQuickDoc = (data) => {
  return request.post({
    url: '/manage/saveLibDocIndexQuickDoc',
    data: {
      library_key: data.library_key,
      doc_id: data.doc_id,
      quick_doc_content: data.quick_doc_content
    }
  })
}

// 保存首页banner
export const saveLibDocBannerImg = (data) => {
  return request.post({
    url: '/manage/saveLibDocBannerImg',
    data: {
      library_key: data.library_key,
      doc_id: data.doc_id,
      banner_img_url: data.banner_img_url
    }
  })
}