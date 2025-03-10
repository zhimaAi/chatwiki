import { useStorage } from '../hooks/web/useStorage'
import { TOKEN_KEY } from '@/constants/index'

const { setStorage, getStorage, removeStorage } = useStorage('localStorage')

export const setTokenCache = (token) => {
  setStorage(TOKEN_KEY, token)
}

export const getTokenCache = () => {
  getStorage(TOKEN_KEY)
}

export const removeTokenCache = () => {
  removeStorage(TOKEN_KEY)
}

export const setRememberMeCache = (rememberMe) => {
  setStorage('rememberMe', rememberMe)
}

export const getRememberMeCache = () => {
  getStorage('rememberMe')
}

export const removeRememberMeCache = () => {
  removeStorage('rememberMe')
}
