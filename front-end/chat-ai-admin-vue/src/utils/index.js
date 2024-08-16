import { message } from 'ant-design-vue'
import { useStorage } from '../hooks/web/useStorage'

export function showErrorMsg(msg) {
  message.destroy()
  message.error(msg)
}

export function showSuccessMsg(msg) {
  message.destroy()
  message.success(msg)
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

/**
 * 数组对象传入关键字去重
 */
export function duplicateRemoval(arr, key) {
  const keyVals = new Set();
  return arr.filter(obj => {
      const val = obj[key];
      if (keyVals.has(val)) return false;
      keyVals.add(val);
      return true;
  });
}

/**
 * 数组对象合并去重
 */
export function removeRepeat() {
  let arr = [].concat.apply([], arguments)
  return Array.from(new Set(arr)).sort()
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
  const { setStorage, getStorage } = useStorage('localStorage')
  let openid = getStorage('openid')

  if (!openid) {
    openid = 'A' + getUuid(16)

    setStorage('openid', openid)
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

export const devLog = (...args) => {
  if (import.meta.env.MODE == 'production') {
    return
  }

  // 确保至少有一个参数
  if (args.length === 0) {
    console.log('No arguments provided.')
    return
  }

  // 添加一个时间戳前缀（可选）
  const timestamp = new Date().toISOString()
  const formattedArgs = [timestamp, ...args]

  // 输出日志
  console.log(...formattedArgs)
}

// 下载文件
export function downloadFile(filename, link) {
  const element = document.createElement('a')
  element.setAttribute('href', link)
  element.setAttribute('download', filename)

  element.style.display = 'none'
  document.body.appendChild(element)

  element.click()

  document.body.removeChild(element)
}

export function objectToQueryString(obj) {
  // 将对象转换为数组，然后map每个键值对到 'key=value' 字符串  
  // 使用 encodeURIComponent 来确保URL安全  
  return Object.keys(obj)
    .map(key => encodeURIComponent(key) + '=' + encodeURIComponent(obj[key]))
    .join('&'); // 用 '&' 连接所有的键值对  
}  