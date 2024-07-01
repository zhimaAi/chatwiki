import { contextBridge } from 'electron'
import { electronAPI } from '@electron-toolkit/preload'
const fs = require('fs')
import { join } from 'path'

const configPath = join(__dirname, '../../resources/config.json')
const __CONFIG__ = {}

try {
  // 使用 fs.readFileSync 同步读取文件内容
  const configData = fs.readFileSync(configPath, 'utf8')
  // 将 JSON 字符串解析为 JavaScript 对象
  const config = JSON.parse(configData)

  Object.assign(__CONFIG__, config)
} catch (error) {
  console.error('Error reading config file:', error)
}

// Custom APIs for renderer
const api = {}

// Use `contextBridge` APIs to expose Electron APIs to
// renderer only if context isolation is enabled, otherwise
// just add to the DOM global.
if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld('__CONFIG__', __CONFIG__)
    contextBridge.exposeInMainWorld('electron', electronAPI)
    contextBridge.exposeInMainWorld('api', api)
  } catch (error) {
    console.error(error)
  }
} else {
  window.__CONFIG__ = __CONFIG__
  window.electron = electronAPI
  window.api = api
}
