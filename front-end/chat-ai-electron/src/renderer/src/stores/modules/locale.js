import { defineStore } from 'pinia'
import { store } from '../index'
import zhCn from 'ant-design-vue/es/locale/zh_CN'
import en from 'ant-design-vue/es/locale/en_US'
import Storage from '@/utils/storage'

const antdvLocaleMap = {
  'zh-CN': zhCn,
  en: en
}

export const useLocaleStore = defineStore('locales', {
  state: () => {
    return {
      currentLocale: {
        lang: Storage.get('lang') || 'zh-CN',
        elLocale: antdvLocaleMap[Storage.get('lang') || 'zh-CN']
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
      this.currentLocale.elLocale = antdvLocaleMap[localeMap?.lang]
      Storage.set('lang', localeMap?.lang)
    }
  }
})

export const useLocaleStoreWithOut = () => {
  return useLocaleStore(store)
}
