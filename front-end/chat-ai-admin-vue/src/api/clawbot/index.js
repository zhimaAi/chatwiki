import request from '@/utils/http/axios'

export const getE2bConf = (params = {}) => {
  return request.get({
    url: '/manage/getE2bConf',
    params
  })
}

export const saveE2bConf = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'application/json' },
    url: '/manage/saveE2bConf',
    data
  })
}

export const saveClawbotConf = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/saveClawbotConf',
    data
  })
}

export const uploadClawbotLocalDoc = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/uploadClawbotLocalDoc',
    data
  })
}

export const uploadClawbotSkillZip = (data = {}, onUploadProgress) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/uploadClawbotSkillZip',
    data,
    onUploadProgress
  })
}

export const saveClawbotSkill = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/saveClawbotSkill',
    data
  })
}

export const getClawbotSkillList = (params = {}) => {
  return request.get({
    url: '/manage/getClawbotSkillList',
    params
  })
}

export const getDocToSkillTaskList = (params = {}) => {
  return request.get({
    url: '/manage/getDocToSkillTaskList',
    params
  })
}

export const createDocToSkillTask = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/createDocToSkillTask',
    data
  })
}

export const stopDocToSkillTask = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    url: '/manage/stopDocToSkillTask',
    data
  })
}

export const installDocToSkill = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    url: '/manage/installDocToSkill',
    data
  })
}

export const getDocToSkillTaskInfo = (params = {}) => {
  return request.get({
    url: '/manage/getDocToSkillTaskInfo',
    params
  })
}

export const regenerateDocToSkillTask = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    url: '/manage/regenerateDocToSkillTask',
    data
  })
}

export const getWebToSkillTaskList = (params = {}) => {
  return request.get({
    url: '/manage/getWebToSkillTaskList',
    params
  })
}

export const createWebToSkillTask = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'application/json' },
    url: '/manage/createWebToSkillTask',
    data
  })
}

export const stopWebToSkillTask = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'application/json' },
    url: '/manage/stopWebToSkillTask',
    data
  })
}

export const regenerateWebToSkillTask = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'application/json' },
    url: '/manage/regenerateWebToSkillTask',
    data
  })
}

export const getWebToSkillTaskInfo = (params = {}) => {
  return request.get({
    url: '/manage/getWebToSkillTaskInfo',
    params
  })
}

export const installWebToSkill = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'application/json' },
    url: '/manage/installWebToSkill',
    data
  })
}

export const saveClawbotRobotSkills = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/saveClawbotRobotSkills',
    data
  })
}

export const getClawbotSkillInfo = (params = {}) => {
  return request.get({
    url: '/manage/getClawbotSkillInfo',
    params
  })
}

export const deleteClawbotSkill = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/deleteClawbotSkill',
    data
  })
}

export const getClawbotLocalDocList = (params = {}) => {
  return request.get({
    url: '/manage/getClawbotLocalDocList',
    params
  })
}

export const deleteClawbotLocalDoc = (data = {}) => {
  return request.post({
    headers: { 'Content-Type': 'multipart/form-data' },
    url: '/manage/deleteClawbotLocalDoc',
    data
  })
}
