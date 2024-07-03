import { Storage } from '@/utils/Storage'

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

  const maskLength = end - start + 1 // 计算需要被替换的字符数量
  const maskedPart = '*'.repeat(maskLength) // 生成替换的字符串

  // 截取字符串的各部分并拼接
  return str.substring(0, start) + maskedPart + str.substring(end + 1)
}

// 把基于字节的文件大小抓换成KB,MB,GB,TB
export function formatFileSize(size) {
  size = size ? Number(size) : 0

  if (size === 0) {
    return '0 B'
  }
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(size) / Math.log(k))
  return (size / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i]
}

export function getUuid(len: number = 16, radix: number = 0) {
  const chars = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'.split('')
  const uuid: Array<string> = []
  let i: number

  radix = radix || chars.length

  if (len) {
    // Compact form
    for (i = 0; i < len; i++) uuid[i] = chars[0 | (Math.random() * radix)]
  } else {
    // rfc4122, version 4 form
    let r: number

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
  let openid = Storage.get('openid')

  if (!openid) {
    openid = 'P' + getUuid(16)

    Storage.set('openid', openid)
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

// 用于转义HTML，防止XSS攻击
export function escapeHTML(unsafeText: string) {
  return unsafeText.replace(/[&<"']/g, function (m) {
    switch (m) {
      case '&':
        return '&amp;'
      case '<':
        return '&lt;'
      case '"':
        return '&quot;'
      default:
        return '&#039;'
    }
  })
}
