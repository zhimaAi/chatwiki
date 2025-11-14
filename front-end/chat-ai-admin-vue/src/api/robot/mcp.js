import request from "@/utils/http/axios/index.js";

export function getMcpServers(params = {}) {
  return request.get({url: '/manage/getMcpServerList', params})
}

export function saveMcpServer(data = {}) {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/saveMcpServer',
    data
  })
}

export function delMcpServer(data = {}) {
  return request.post({url: '/manage/deleteMcpServer', data})
}

export function updateMcpSrvStatus(data = {}) {
  return request.post({url: '/manage/updateMcpServerPublishStatus', data})
}

export function saveMcpTool(data = {}) {
  return request.post({url: '/manage/saveMcpTool', data})
}

export function editMcpTool(data = {}) {
  return request.post({url: '/manage/editMcpTool', data})
}

export function delMcpTool(data = {}) {
  return request.post({url: '/manage/deleteMcpTool', data})
}
