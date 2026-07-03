import request from '@/utils/http/axios'

// 获取企业小程序卡片列表
export const getMiniCardList = (params = {}) => {
  return request.get({
    url: '/manage/AdminMiniCard/getList',
    params: params
  })
}

// 新增小程序卡片
export const addMiniCard = (data) => {
  return request.post({
    url: '/manage/AdminMiniCard/addOne',
    data: data
  })
}

// 编辑小程序卡片
export const updateMiniCard = (data) => {
  return request.post({
    url: '/manage/AdminMiniCard/updateOne',
    data: data
  })
}

// 删除小程序卡片
export const deleteMiniCard = (data) => {
  return request.post({
    url: '/manage/AdminMiniCard/deleteOne',
    data: data
  })
}

