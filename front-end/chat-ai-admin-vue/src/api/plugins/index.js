import request from "@/utils/http/axios/index.js";

export function getPluginTypes(params = {}) {
  return request.get({url: '/manage/plugin/remote-plugins/typeList', params})
}

export function getRemotePlugins(params = {}) {
  return request.get({url: '/manage/plugin/remote-plugins/list', params})
}

export function downloadPlugin(data = {}) {
  return request.post({url: '/manage/plugin/remote-plugins/download', data})
}

export function getInstallPlugins(params = {}) {
  return request.get({url: '/manage/plugin/local-plugins/list', params})
}

export function triggerConfigList(params = {}) {
  return request.get({url: '/manage/triggerConfigList', params})
}

export function uninstallPlugin(data = {}) {
  return request.post({url: '/manage/plugin/local-plugins/destroy', data})
}

export function openPlugin(data = {}) {
  return request.post({url: '/manage/plugin/local-plugins/load', data})
}

export function closePlugin(data = {}) {
  return request.post({url: '/manage/plugin/local-plugins/unload', data})
}

export function getPluginInfo(params = {}) {
  return request.get({url: '/manage/plugin/local-plugins/detail', params})
}

export function getPluginConfig(params = {}) {
  return request.get({url: '/manage/plugin/local-plugins/config', params})
}

export function setPluginConfig(data = {}) {
  return request.post({url: '/manage/plugin/local-plugins/config', data})
}

export function runPlugin(data = {}) {
  return request.post({url: '/manage/plugin/local-plugins/run', data})
}


export function triggerSwitch(data = {}) {
  return request.post({url: '/manage/triggerSwitch', data})
}