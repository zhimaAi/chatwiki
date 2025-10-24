import SSE from '@/utils/http/sse'
import request from '@/utils/http/axios'

const baseURL = import.meta.env.VITE_BASE_API_URL

export const sendAiMessage = (data) => {
  return new SSE({
    url: baseURL + '/chat/request',
    data: data
  })
}

export const chatWelcome = (data) => {
  return request.post({ url: '/chat/welcome', data })
}

export const getDialogueList = (data) => {
  return request.post({ url: '/manage/getDialogueList', data })
}

export const getChatMessage = (data) => {
  return request.post({ url: '/chat/message', data })
}

// 获取会话记录列表
export const getSessionRecordList = (params) => {
  return request.get({ url: '/manage/getSessionRecordList', params })
}

// 创建会话记录导出任务
export const createSessionExport = (params) => {
  return request.get({ url: '/manage/createSessionExport', params })
}

// 获取导出任务列表
export const getExportTaskList = (params) => {
  return request.get({ url: '/manage/getExportTaskList', params })
}

// 获取未知问题统计列表
export const unknownIssueStats = (params) => {
  return request.get({ url: '/manage/unknownIssueStats', params })
}

// 导出任务下载文件
export const downloadExportFile = (params) => {
  return request.get({ url: '/manage/downloadExportFile', params })
}

// 获取会话来源列表
export const getSessionChannelList = (params) => {
  return request.get({ url: '/manage/getSessionChannelList', params })
}

// 获取答案来源
export const getAnswerSource = (params) => {
  return request.get({ url: '/manage/getAnswerSource', params })
}

// 猜你想问
export const questionGuide = (data) => {
  return request.post({ url: '/chat/questionGuide', data })
}