import request from '@/utils/http/axios'

export const getRobotList = (params = {}) => {
  return request.get({
    url: '/manage/getRobotList',
    params: params
  })
}

export const saveRobot = (data = {}) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/saveRobot',
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
