import { useUserStore } from '@/stores/modules/user'

/**
 * 自定义hook，用于打开带token的链接
 * @returns {Function} openUrlWithToken函数，接收一个URL参数
 */
export const useOpenUrlWithToken = () => {
  // 获取用户store实例
  const userStore = useUserStore()
  
  /**
   * 打开带token的链接
   * @param {string} url - 需要打开的链接
   */
  const openUrlWithToken = (url) => {
    // 从store中获取token
    const token = userStore.getToken
    
    // 检查token和URL是否存在
    if (!token) {
      console.warn('Token is missing')
      // 可以选择是否在没有token的情况下打开链接
      // window.open(url, '_blank')
      return
    }
    
    if (!url) {
      console.warn('URL is missing')
      return
    }
    
    // 检查URL是否已经包含查询参数
    const separator = url.includes('?') ? '&' : '?'
    
    // 构造带token的完整URL
    const urlWithToken = `${url}${separator}token=${token}`
    
    // 在新窗口中打开链接
    window.open(urlWithToken, '_blank')
  }
  
  return {
    openUrlWithToken
  }
}