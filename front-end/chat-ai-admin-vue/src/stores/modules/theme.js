import { defineStore } from 'pinia'
import { DEFAULT_THEME_ID, themes } from '@/config/themes'
import { useStorage } from '@/hooks/web/useStorage'

export const THEME_STORAGE_KEY = 'chatwiki_theme_v1'
const { getStorage, setStorage } = useStorage('localStorage')
const isThemeId = (value) => typeof value === 'string' && Object.prototype.hasOwnProperty.call(themes, value)

export const useThemeStore = defineStore('theme', {
  state: () => ({
    themeId: DEFAULT_THEME_ID,
    storageError: false
  }),
  getters: {
    currentTheme: (state) => themes[state.themeId],
    antTheme() {
      return this.currentTheme.antTheme
    }
  },
  actions: {
    initialize() {
      this.themeId = DEFAULT_THEME_ID
      this.storageError = false

      try {
        const saved = getStorage(THEME_STORAGE_KEY)
        this.themeId = isThemeId(saved) ? saved : DEFAULT_THEME_ID
      } catch {
        // 无效 JSON 或存储被禁用时，保持默认主题，仍允许当前页面切换。
        this.storageError = true
      }
    },
    savePreference() {
      try {
        setStorage(THEME_STORAGE_KEY, this.themeId)
        this.storageError = false
        return true
      } catch {
        this.storageError = true
        return false
      }
    },
    setTheme(themeId) {
      if (!isThemeId(themeId)) throw new Error('Invalid theme id')

      this.themeId = themeId
      return this.savePreference()
    },
    syncFromStorage(event) {
      if (event.key !== THEME_STORAGE_KEY && event.key !== null) return

      try {
        if (event.storageArea !== window.localStorage) return
      } catch {
        this.storageError = true
        return
      }

      // 删除偏好或清空存储时直接恢复默认。
      if (event.newValue === null) {
        this.themeId = DEFAULT_THEME_ID
        this.storageError = false
        return
      }
      this.initialize()
    }
  }
})
