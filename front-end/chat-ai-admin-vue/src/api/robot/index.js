import request from '@/utils/http/axios'

export const getRobotList = (params = {}) => {
  return request.get({
    url: '/manage/getRobotList',
    params: params
  })
}

export const saveRobot = (data = {}, application_type = 0) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: application_type == 0 ? '/manage/saveRobot' : '/manage/addFlowRobot',
    data: data
  })
}

export const getRobotInfo = ({ id }) => {
  return request.get({
    url: '/manage/getRobotInfo',
    params: { id }
  })
}

export const updateFastCommandSwitch = (data = {}) => {
  return request.post({
    url: '/manage/updateFastCommandSwitch',
    data: data
  })
}

export const deleteRobot = ({ id }) => {
  return request.post({
    url: '/manage/deleteRobot',
    data: {
      id
    }
  })
}

export const robotCopy = ({ id }) => {
  return request.post({
    url: '/manage/robotCopy',
    data: {
      from_id: id
    }
  })
}

export const editPrompt = ({ id, prompt }) => {
  return request.post({
    url: '/manage/editPrompt',
    data: {
      id,
      prompt
    }
  })
}

export const editExternalConfig = (data = {}) => {
  return request.post({
    url: '/manage/editExternalConfig',
    data: data
  })
}

export const getWechatAppList = ({ robot_id, app_type, app_name }) => {
  return request.get({
    url: '/manage/getWechatAppList',
    params: {
      robot_id,
      app_type,
      app_name
    }
  })
}

export const saveWechatApp = (data) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/saveWechatApp',
    data: data
  })
}

export const deleteWechatApp = ({ id }) => {
  return request.post({
    url: '/manage/deleteWechatApp',
    data: { id }
  })
}

export const getFastCommandList = (params = {}) => {
  return request.get({
    url: '/manage/getFastCommandList',
    params: params
  })
}


export const getFastCommandInfo = (params = {}) => {
  return request.get({
    url: '/manage/GetFastCommandInfo',
    params: params
  })
}

export const saveFastCommand = (data = {}) => {
  return request.post({
    url: '/manage/saveFastCommand',
    data: data
  })
}

export const deleteFastCommand = (data = {}) => {
  return request.post({
    url: '/manage/deleteFastCommand',
    data: data
  })
}

export const sortFastCommand = (data = {}) => {
  return request.post({
    headers: {
      'Content-Type': 'application/json'
    },
    url: '/manage/sortFastCommand',
    data: data
  })
}

export const listRobotApikey = (data = {}) => {
  return request.post({
    url: '/manage/listRobotApikey',
    data: data
  })
}

export const updateRobotApikey = (data = {}) => {
  return request.post({
    url: '/manage/updateRobotApikey',
    data: data
  })
}

export const addRobotApikey = (data = {}) => {
  return request.post({
    url: '/manage/addRobotApikey',
    data: data
  })
}

export const deleteRobotApikey = (data = {}) => {
  return request.post({
    url: '/manage/deleteRobotApikey',
    data: data
  })
}


export const getNodeList = (params = {}) => {
  return request.get({
    url: '/manage/getNodeList',
    params: params
  })
}

export const saveNodes = (data = {}) => {
  return request.post({
    url: '/manage/saveNodes',
    data: data
  })
}

export const createPromptByAi = (params = {}) => {
  return request.get({
    url: '/manage/createPromptByAi',
    params: params
  })
}

export const getSensitiveWordsList = (params = {}) => {
  return request.get({
    url: '/manage/getSensitiveWordsList',
    params: params
  })
}

export const switchSensitiveWords = (data = {}) => {
  return request.post({
    url: '/manage/switchSensitiveWords',
    data: data
  })
}

export const deleteSensitiveWords = (data = {}) => {
  return request.post({
    url: '/manage/deleteSensitiveWords',
    data: data
  })
}


export const saveSensitiveWords = (data = {}) => {
  return request.post({
    url: '/manage/saveSensitiveWords',
    data: data
  })
}

export const checkSensitiveWords = (data = {}) => {
  return request.post({
    url: '/manage/checkSensitiveWords',
    data: data
  })
}

// 发起GPT提问前置判断
export const checkChatRequestPermission = (data = {}) => {
  return request.post({
    url: '/chat/checkChatRequestPermission',
    data: data
  })
}

// 编辑工作量基本信息
export const editBaseInfo = (data = {}) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/editBaseInfo',
    data: data
  })
}

export const relationWorkFlow = (data = {}) => {
  return request.post({
    url: '/manage/relationWorkFlow',
    data: data
  })
}

// 导入csl文件
export const robotImport = (data = {}) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/robotImport',
    data: data
  })
}

export const getUnknownIssueSummary = (params = {}) => {
  return request.get({
    url: '/manage/getUnknownIssueSummary',
    params: params
  })
}

export const setUnknownIssueSummary = (data = {}) => {
  return request.post({
    url: '/manage/setUnknownIssueSummary',
    data: data
  })
}


export const unknownIssueSummaryImport = (data = {}) => {
  return request.post({
    url: '/manage/unknownIssueSummaryImport',
    data: data
  })
}

export const unknownIssueSummaryAnswer = (data = {}) => {
  return request.post({
    url: '/manage/unknownIssueSummaryAnswer',
    data: data
  })
}

export const relationLibrary = (data = {}) => {
  return request.post({
    url: '/manage/relationLibrary',
    data: data
  })
}

export const robotAutoAdd = (data = {}) => {
  return request.post({
    url: '/manage/robotAutoAdd',
    data: data
  })
}

export const getRobotGroupList = (params = {}) => {
  return request.get({
    url: '/manage/getRobotGroupList',
    params: params
  })
}

export const saveRobotGroup = (data = {}) => {
  return request.post({
    url: '/manage/saveRobotGroup',
    data: data
  })
}

export const deleteRobotGroup = (data = {}) => {
  return request.post({
    url: '/manage/deleteRobotGroup',
    data: data
  })
}

export const relationRobotGroup = (data = {}) => {
  return request.post({
    url: '/manage/relationRobotGroup',
    data: data
  })
}

export const callWorkFlow = (data = {}) => {
  return request.post({
    url: '/chat/callWorkFlow',
    data: data
  })
}
export const workFlowNextVersion = (data = {}) => {
  return request.post({
    url: '/manage/workFlowNextVersion',
    data: data
  })
}

export const workFlowPublishVersion = (data = {}) => {
  return request.post({
    url: '/manage/workFlowPublishVersion',
    data: data
  })
}

export const workFlowVersions = (data = {}) => {
  return request.post({
    url: '/manage/workFlowVersions',
    data: data
  })
}

export const workFlowVersionDetail = (data = {}) => {
  return request.post({
    url: '/manage/workFlowVersionDetail',
    data: data
  })
}

export const getDraftKey = (params = {}) => {
  return request.get({
    url: '/manage/getDraftKey',
    params: params
  })
}

export const getAdminConfig = (params = {}) => {
  return request.get({
    url: '/manage/getAdminConfig',
    params: params
  })
}

export const saveDraftExTime = (data = {}) => {
  return request.post({
    url: '/manage/saveDraftExTime',
    data: data
  })
}

export const refreshAccountVerify   = (data = {}) => {
  return request.post({
    url: '/manage/refreshAccountVerify',
    data: data
  })
}

export const setWechatNotVerifyConfig = (data = {}) => {
  return request.post({
    url: '/manage/setWechatNotVerifyConfig',
    data: data
  })
}