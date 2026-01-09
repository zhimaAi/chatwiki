import { defineStore } from 'pinia'
import { store } from '../index'
import zhCn from 'ant-design-vue/es/locale/zh_CN'
import en from 'ant-design-vue/es/locale/en_US'
import { useStorage } from '@/hooks/web/useStorage'

const { getStorage, setStorage } = useStorage('localStorage')

const antvLocaleMap = {
  'zh-CN': zhCn,
  'en-US': en,
  'en': en  // 兼容旧版本
}

export const useLocaleStore = defineStore('locales', {
  state: () => {
    return {
      currentLocale: {
        lang: getStorage('lang') || 'zh-CN',
        antvLocale: antvLocaleMap[getStorage('lang') || 'zh-CN']
      },
      // 多语言
      localeMap: [
        {
          lang: 'zh-CN',
          name: '简体中文'
        },
        {
          lang: 'en-US',
          name: 'English'
        }
      ]
    }
  },
  getters: {
    getCurrentLocale() {
      return this.currentLocale
    },
    getLocaleMap() {
      return this.localeMap
    },
    getSelectedLocale() {
      return this.localeMap.filter((item) => item.lang == this.currentLocale.lang)[0]
    }
  },
  actions: {
    setCurrentLocale(localeMap) {
      // this.locale = Object.assign(this.locale, localeMap)
      this.currentLocale.lang = localeMap?.lang
      this.currentLocale.antvLocale = antvLocaleMap[localeMap?.lang]
      setStorage('lang', localeMap?.lang)
    }
  }
})

export const useLocaleStoreWithOut = () => {
  return useLocaleStore(store)
}
