import request from '@/utils/http/axios'

export const saveUserAbility = (data) => {
  return request.post({
    url: '/manage/ability/saveUserAbility',
    data: data
  })
}

export const getAbilityList = (params = {}) => {
  return request.get({
    url: '/manage/ability/getAbilityList',
    params: params
  })
}

export const saveRobotAbility = (data) => {
  return request.post({
    url: '/manage/ability/saveRobotAbility',
    data: data
  })
}

export const saveRobotAbilitySwitchStatus = (data) => {
  return request.post({
    url: '/manage/ability/saveRobotAbilitySwitchStatus',
    data: data
  })
}

export const saveRobotAbilityFixedMenu = (data) => {
  return request.post({
    url: '/manage/ability/saveRobotAbilityFixedMenu',
    data: data
  })
}

export const saveRobotAbilityAiReplyStatus = (data) => {
  return request.post({
    url: '/manage/ability/saveRobotAbilityAiReplyStatus',
    data: data
  })
}

export const saveRobotReceivedMessageReply = (data) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/ability/saveRobotReceivedMessageReply',
    data: data
  })
}

export const getRobotAbilityList = (params = {}) => {
  return request.get({
    url: '/manage/ability/getRobotAbilityList',
    params: params
  })
}

export const checkKeyWordRepeat = (data) => {
  return request.post({
    url: '/manage/ability/checkKeyWordRepeat',
    data: data
  })
}

export const saveRobotKeywordReply = (data) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/ability/saveRobotKeywordReply',
    data: data
  })
}

export const deleteRobotKeywordReply = (data) => {
  return request.post({
    url: '/manage/ability/deleteRobotKeywordReply',
    data: data
  })
}

export const deleteRobotReceivedMessageReply = (data) => {
  return request.post({
    url: '/manage/ability/deleteRobotReceivedMessageReply',
    data: data
  })
}

export const getRobotKeywordReply = (params = {}) => {
  return request.get({
    url: '/manage/ability/getRobotKeywordReply',
    params: params
  })
}

export const getRobotReceivedMessageReply = (params = {}) => {
  return request.get({
    url: '/manage/ability/getRobotReceivedMessageReply',
    params: params
  })
}

export const getRobotKeywordReplyList = (params = {}) => {
  return request.get({
    url: '/manage/ability/getRobotKeywordReplyList',
    params: params
  })
}

export const getRobotReceivedMessageReplyList = (params = {}) => {
  return request.get({
    url: '/manage/ability/getRobotReceivedMessageReplyList',
    params: params
  })
}

export const updateRobotKeywordReplySwitchStatus = (data) => {
  return request.post({
    url: '/manage/ability/updateRobotKeywordReplySwitchStatus',
    data: data
  })
}

export const updateRobotReceivedMessageReplySwitchStatus = (data) => {
  return request.post({
    url: '/manage/ability/updateRobotReceivedMessageReplySwitchStatus',
    data: data
  })
}

export const updateRobotReceivedMessageReplyPriorityNum = (data) => {
  return request.post({
    url: '/manage/ability/updateRobotReceivedMessageReplyPriorityNum',
    data: data
  })
}

// Subscribe reply APIs
export const saveRobotSubscribeReply = (data) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/ability/saveRobotSubscribeReply',
    data: data
  })
}

export const deleteRobotSubscribeReply = (data) => {
  return request.post({
    url: '/manage/ability/deleteRobotSubscribeReply',
    data: data
  })
}

export const getRobotSubscribeReply = (params = {}) => {
  return request.get({
    url: '/manage/ability/getRobotSubscribeReply',
    params: params
  })
}

export const getRobotSubscribeReplyList = (params = {}) => {
  return request.get({
    url: '/manage/ability/getRobotSubscribeReplyList',
    params: params
  })
}

export const updateRobotSubscribeReplySwitchStatus = (data) => {
  return request.post({
    url: '/manage/ability/updateRobotSubscribeReplySwitchStatus',
    data: data
  })
}

export const updateRobotSubscribeReplyPriorityNum = (data) => {
  return request.post({
    url: '/manage/ability/updateRobotSubscribeReplyPriorityNum',
    data: data
  })
}