import request from '@/utils/http/axios'

export const getDomainList = (data) => {
  return request.post({ url: '/yun/yunAdminApi/GetDomainList', data })
}

export const getSubjectInfoList = (data) => {
  return request.post({ url: '/yun/YunOneSubMsgApi/GetSubjectInfoList', data })
}

export const addChannel = (data) => {
  return request.post({ url: '/yun/yunAdmin/ChannelCreate', data })
}
