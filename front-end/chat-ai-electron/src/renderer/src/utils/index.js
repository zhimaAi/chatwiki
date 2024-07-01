import { message } from 'ant-design-vue'
import { Storage } from './storage'

export function showErrorMsg(msg) {
  message.destroy()
  message.error(msg)
}

export function showSuccessMsg(msg) {
  message.destroy()
  message.success(msg)
}

export function getAdminUserId() {
  let admin_user_id = window.__CONFIG__.ADMIN_USER_ID || import.meta.env.VITE_ADMIN_USER_ID

  return admin_user_id
}

export function getBaseApiUrl() {
  let base_api_url = window.__CONFIG__.BASE_API_URL || import.meta.env.VITE_BASE_API_URL

  return base_api_url
}

/**
 * 把对象转为formData
 */
export function objToFormData(obj) {
  const formData = new FormData()
  Object.keys(obj).forEach((key) => {
    if (obj[key] !== void 0 && obj[key] !== null) {
      formData.append(key, obj[key])
    }
  })
  return formData
}

// 字符串加密替换成*
export function strEncryption(str, start = 0, end = str.length - 1) {
  if (start < 0 || end >= str.length || start > end) {
    throw new Error('Invalid start or end position')
  }

  let maskLength = end - start + 1 // 计算需要被替换的字符数量
  let maskedPart = '*'.repeat(maskLength) // 生成替换的字符串

  // 截取字符串的各部分并拼接
  return str.substring(0, start) + maskedPart + str.substring(end + 1)
}

// 把基于字节的文件大小抓换成KB,MB,GB,TB
export function formatFileSize(size) {
  size = size ? Number(size) : 0

  if (size === 0) {
    return '0 B'
  }
  let k = 1024
  let sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = Math.floor(Math.log(size) / Math.log(k))
  return (size / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i]
}

export function getUuid(len, radix) {
  var chars = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'.split('')
  var uuid = [],
    i
  radix = radix || chars.length

  if (len) {
    // Compact form
    for (i = 0; i < len; i++) uuid[i] = chars[0 | (Math.random() * radix)]
  } else {
    // rfc4122, version 4 form
    var r

    // rfc4122 requires these characters
    uuid[8] = uuid[13] = uuid[18] = uuid[23] = '-'
    uuid[14] = '4'

    // Fill in random data.  At i==19 set the high bits of clock sequence as
    // per rfc4122, sec. 4.1.5
    for (i = 0; i < 36; i++) {
      if (!uuid[i]) {
        r = 0 | (Math.random() * 16)
        uuid[i] = chars[i == 19 ? (r & 0x3) | 0x8 : r]
      }
    }
  }

  return uuid.join('')
}

export function getOpenid() {
  const { get, set } = Storage
  let openid = get('openid')

  if (!openid) {
    openid = getUuid(16)

    set('openid', openid)
  }

  return openid
}

export function getBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => resolve(reader.result)
    reader.onerror = (error) => reject(error)
  })
}

/**
 * 复制文本
 * @param text
 */
export const copyText = (text) => {
  const copyInput = document.createElement('textarea')
  copyInput.setAttribute('readonly', 'readonly')
  copyInput.value = text
  document.body.appendChild(copyInput)
  copyInput.select()
  document.execCommand('copy')
  copyInput.remove()
}
