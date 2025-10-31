import { message } from 'ant-design-vue'
import { useStorage } from '../hooks/web/useStorage'
import CryptoJS from 'crypto-js'
import dayjs from 'dayjs'

export function generateUniqueId(salt = '') {
  // 获取当前时间戳（毫秒级）
  const timestamp = Date.now()

  // 生成一个32位的随机数
  let chars = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
  let randomPart = ''

  // 循环32次，每次随机选择一个字符
  for (let i = 0; i < 36; i++) {
    // 随机选择一个字符并添加到salt中
    const randomChar = chars[Math.floor(Math.random() * chars.length)]
    randomPart += randomChar
  }

  // 将时间戳、随机数和盐值拼接起来
  const idParts = [timestamp, randomPart, salt]

  // 使用MD5或其他哈希算法对拼接后的字符串进行哈希，以生成固定长度的唯一ID

  // 对拼接后的字符串进行哈希，并截取前16位作为唯一ID的一部分
  const hashedPart = CryptoJS.MD5(idParts.join('')).toString()

  // 返回32位唯一ID
  return hashedPart
}

export function generateRandomId(length) {
  const possibleCharacters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789' // 所有可能的字符
  let randomString = ''

  for (let i = 0; i < length; i++) {
    const randomIndex = Math.floor(Math.random() * possibleCharacters.length) // 生成一个随机索引
    randomString += possibleCharacters.charAt(randomIndex) // 获取一个随机字符并追加到结果字符串
  }

  return randomString // 返回生成的随机字符串
}

// 防抖函数
const delayedClickMap = {}
export function delayedClick(key, time = 3000) {
  if (Object.prototype.hasOwnProperty.call(delayedClickMap, key)) {
    return false
  }

  delayedClickMap[key] = true

  setTimeout(() => {
    delete delayedClickMap[key]
  }, time)

  return true
}

// 分钟转秒
export function minutesToSeconds(minutes) {
  return minutes * 60
}

// 秒转分钟
export function secondsToMinutes(seconds) {
  return seconds / 60
}

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
  const keyVals = new Set()
  return arr.filter((obj) => {
    const val = obj[key]
    if (keyVals.has(val)) return false
    keyVals.add(val)
    return true
  })
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

  // 使用 toFixed 来确保固定的小数位数，这里选择2位小数
  // 你可以根据需要调整小数位数
  let formattedSize = (size / Math.pow(k, i)).toFixed(2)

  // 如果格式化后的结果末尾是.00，则去掉这两个0
  if (formattedSize.endsWith('.00')) {
    formattedSize = formattedSize.slice(0, -3)
  }

  return formattedSize + '' + sizes[i]
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

export function unicodeToBase64(str) {
  return btoa(
    encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, function (match, p1) {
      return String.fromCharCode(parseInt(p1, 16))
    })
  )
}

export function base64ToUnicode(base64) {
  return decodeURIComponent(
    Array.prototype.map
      .call(atob(base64), function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
      })
      .join('')
  )
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
    .map((key) => encodeURIComponent(key) + '=' + encodeURIComponent(obj[key]))
    .join('&') // 用 '&' 连接所有的键值对
}

export function tableToExcel(str, jsonData, fieds, name) {
  //jsonData要导出的json数据
  //str列标题，逗号隔开，每一个逗号就是隔开一个单元格
  for (let i = 0; i < jsonData.length; i++) {
    for (let item of fieds) {
      str += `"${jsonData[i][item] + '\t'}",`
    }
    str += '\n'
  }
  //encodeURIComponent解决中文乱码
  let uri = 'data:text/csv;charset=utf-8,\ufeff' + encodeURIComponent(str)
  //通过创建a标签实现
  let link = document.createElement('a')
  link.href = uri
  //对下载的文件命名
  link.download = name
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

export function exportToJsonWithSaver(data, filename = 'data.json') {
  const jsonString = JSON.stringify(data, null, 2) // 将数据转换为格式化的 JSON 字符串
  const blob = new Blob([jsonString], { type: 'application/json' })
  const url = URL.createObjectURL(blob)

  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()

  setTimeout(() => {
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
  }, 0)
}

export function formatPriceWithCommas(price) {
  // 将价格转换为字符串，以确保小数部分不会被截断
  const priceStr = price.toString();

  // 分离整数部分和小数部分
  const [integerPart, decimalPart] = priceStr.split('.');

  // 使用正则表达式将整数部分转换为千分位格式
  const formattedIntegerPart = integerPart.replace(/\B(?=(\d{3})+(?!\d))/g, ",");

  // 如果存在小数部分，则将其添加回结果中
  const formattedPrice = decimalPart ? `${formattedIntegerPart}.${decimalPart}` : formattedIntegerPart;

  return formattedPrice;
}

// 格式化显示时间
export function formatDisplayChatTime(time) {
  time = time * 1000

  const date = new Date(time)
  const now = new Date()
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)

  const isToday = date.toDateString() === now.toDateString()
  const isYesterday = date.toDateString() === yesterday.toDateString()

  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')
  const seconds = date.getSeconds().toString().padStart(2, '0')
  const timeStr = `${hours}:${minutes}:${seconds}`

  if (isToday) {
    return `今天 ${timeStr}`
  } else if (isYesterday) {
    return `昨天 ${timeStr}`
  } else {
    const year = date.getFullYear()
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const day = date.getDate().toString().padStart(2, '0')
    return `${year}-${month}-${day} ${timeStr}`
  }
}

/**
 * 在时间格式和分钟数之间转换
 * @param {string|number} input - 输入可以是时间格式字符串（如 "00:10"）或分钟数（如 10 或 "70"）
 * @returns {string|number} - 返回转换后的结果
 */
export function convertTime(input) {
  // 如果输入是数字，或者字符串形式的数字，转换为时间格式
  if (typeof input === 'number' || (typeof input === 'string' && !isNaN(input))) {
    const minutes = parseInt(input, 10)
    const hours = Math.floor(minutes / 60)
    const remainingMinutes = minutes % 60
    // 使用 padStart 确保两位数格式
    return `${String(hours).padStart(2, '0')}:${String(remainingMinutes).padStart(2, '0')}`
  }

  // 如果输入是时间格式字符串，转换为分钟数
  if (typeof input === 'string') {
    const parts = input.split(':')
    if (parts.length === 2) {
      const hours = parseInt(parts[0], 10) || 0
      const minutes = parseInt(parts[1], 10) || 0
      return hours * 60 + minutes
    }
  }

  // 如果输入无效，返回原始值或抛出错误
  return input
}

export function formateDepartmentData(data, level = 0) {
  if (!Array.isArray(data)) return data;
  
  return data.map(item => {
    // 创建新对象，避免修改原数据
    const newItem = { ...item };
    
    // 添加 level 字段
    newItem.level = level;
    
    // 添加 key 和 title
    newItem.key = newItem.id;
    newItem.title = newItem.department_name;
    
    // 处理 children 和 user_data
    newItem.children = [];
    
    // 如果有子部门，递归处理，层级+1
    if (item.children && Array.isArray(item.children)) {
      newItem.children = formateDepartmentData(item.children, level + 1);
    }
    
    // 添加用户数据到 children
    // if (item.user_data && Array.isArray(item.user_data)) {
    //   const users = item.user_data.map(user => ({
    //     ...user,
    //     key: user.id,
    //     title: user.user_name,
    //     // 用户数据不需要 department 相关字段
    //     id: user.id,
    //     avatar: user.avatar,
    //     user_name: user.user_name,
    //     is_user: true,
    //     // 用户数据的层级与当前部门相同
    //     level: level
    //   }));
    //   newItem.children = [...newItem.children, ...users];
    // }
    
    return newItem;
  });
}

export function formateDepartmentCascaderData(data, level = 0, path) {
  if (!Array.isArray(data)) return data;
  return data.map(item => {
    // 创建新对象，避免修改原数据
    const newItem = { ...item };
    
    // 添加 level 字段
    newItem.level = level;
    
    // 添加 key 和 title
    newItem.key = newItem.id;
    newItem.title = newItem.department_name;
    newItem.label = newItem.department_name
    newItem.value = newItem.key
    newItem.path = path ? [...path, item.id] : [item.id];
    
    // 处理 children 和 user_data
    newItem.children = [];
    
    // 如果有子部门，递归处理，层级+1
    if (item.children && Array.isArray(item.children)) {
      newItem.children = formateDepartmentCascaderData(item.children, level + 1, [...newItem.path]);
    }
    return newItem;
  });
}


export function getDateRangePresets() {
  return [
    {
      label: '今天',
      value: [dayjs().startOf('day'), dayjs().endOf('day')]
    },
    {
      label: '昨天',
      value: [dayjs().subtract(1, 'day').startOf('day'), dayjs().subtract(1, 'day').endOf('day')]
    },
    {
      label: '本周',
      value: [dayjs().startOf('week'), dayjs().endOf('week')]
    },
    {
      label: '本月',
      value: [dayjs().startOf('month'), dayjs().endOf('month')]
    },
    {
      label: '本年',
      value: [dayjs().startOf('year'), dayjs().endOf('year')]
    },
    {
      label: '近7天',
      value: [dayjs().subtract(6, 'day').startOf('day'), dayjs().endOf('day')]
    },
    {
      label: '近30天',
      value: [dayjs().subtract(29, 'day').startOf('day'), dayjs().endOf('day')]
    },
    {
      label: '近90天',
      value: [dayjs().subtract(89, 'day').startOf('day'), dayjs().endOf('day')]
    }
  ]
}