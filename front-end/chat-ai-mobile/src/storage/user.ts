import { Storage } from '@/utils/Storage'
import { TOKEN_KEY } from '@/constants/index'

export const setTokenCache = (token) => {
  Storage.set(TOKEN_KEY, token)
}

export const getTokenCache = () => {
  Storage.get(TOKEN_KEY)
}

export const removeTokenCache = () => {
  Storage.remove(TOKEN_KEY)
}

export const setRememberMeCache = (rememberMe) => {
  Storage.set('rememberMe', rememberMe)
}

export const getRememberMeCache = () => {
  Storage.get('rememberMe')
}

export const removeRememberMeCache = () => {
  Storage.remove('rememberMe')
}
