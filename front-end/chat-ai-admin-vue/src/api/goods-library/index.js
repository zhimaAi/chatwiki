import request from '@/utils/http/axios'
import { useUserStore } from '@/stores/modules/user'

export const getGoodsGroupList = (params = {}) => {
  return request.get({
    url: '/manage/getGoodsGroupList',
    params
  })
}

export const saveGoodsGroup = (data) => {
  return request.post({
    url: '/manage/saveGoodsGroup',
    data
  })
}

export const deleteGoodsGroup = (data) => {
  return request.post({
    url: '/manage/deleteGoodsGroup',
    data
  })
}

export const sortGoodsGroup = (data) => {
  return request.post({
    url: '/manage/sortGoodsGroup',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

export const getGoodsList = (params = {}) => {
  return request.get({
    url: '/manage/getGoodsLibraryList',
    params
  })
}

export const getGoodsInfo = (params = {}) => {
  return request.get({
    url: '/manage/getGoodsLibraryInfo',
    params
  })
}

export const saveGoods = (data) => {
  return request.post({
    url: '/manage/saveGoodsLibrary',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

export const deleteGoods = (data) => {
  return request.post({
    url: '/manage/deleteGoodsLibrary',
    data
  })
}

export const updateGoodsEnabled = (data) => {
  return request.post({
    url: '/manage/updateGoodsLibrarySwitch',
    data
  })
}

export const uploadGoodsImage = (data) => {
  return request.post({
    url: '/manage/uploadGoodsLibraryImage',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export const importGoodsLibrary = (data) => {
  return request.post({
    url: '/manage/importGoodsLibrary',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export const downloadImportTemplate = () => {
  const userStore = useUserStore()
  window.open(
    `/manage/downloadGoodsLibraryImportTemplate?token=${userStore.getToken}`,
    '_blank',
    'noopener,noreferrer'
  )
}

export const exportGoodsList = (params = {}) => {
  const userStore = useUserStore()
  const searchParams = new URLSearchParams()

  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      searchParams.append(key, value)
    }
  })

  const query = searchParams.toString()
  const target = query
    ? `/manage/exportGoodsLibrary?${query}`
    : '/manage/exportGoodsLibrary'

  window.open(`${target}${query ? '&' : '?'}token=${userStore.getToken}`, '_blank', 'noopener,noreferrer')
}
