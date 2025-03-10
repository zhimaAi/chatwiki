// 获取传入的值的类型
const getValueType = (value) => {
  const type = Object.prototype.toString.call(value)
  return type.slice(8, -1)
}
// type = 'sessionStorage' | 'localStorage'
export const useStorage = (type = 'localStorage') => {
  const setStorage = (key, value) => {
    const valueType = getValueType(value)
    window[type].setItem(key, JSON.stringify({ type: valueType, value }))
  }

  const getStorage = (key) => {
    const value = window[type].getItem(key)
    if (value) {
      const { value: val } = JSON.parse(value)
      return val
    } else {
      return value
    }
  }

  const removeStorage = (key) => {
    window[type].removeItem(key)
  }

  const clear = (excludes) => {
    // 获取排除项
    const keys = Object.keys(window[type])
    const defaultExcludes = ['dynamicRouter', 'serverDynamicRouter']
    const excludesArr = excludes ? [...excludes, ...defaultExcludes] : defaultExcludes
    const excludesKeys = excludesArr ? keys.filter((key) => !excludesArr.includes(key)) : keys
    // 排除项不清除
    excludesKeys.forEach((key) => {
      window[type].removeItem(key)
    })
    // window[type].clear()
  }

  return {
    setStorage,
    getStorage,
    removeStorage,
    clear
  }
}
