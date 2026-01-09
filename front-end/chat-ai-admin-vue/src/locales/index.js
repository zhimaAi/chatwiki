import { createI18n } from 'vue-i18n'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { setHtmlPageLang } from './helper'

export let i18n = ''

const createI18nOptions = async () => {
  const localeStore = useLocaleStoreWithOut()
  const locale = localeStore.getCurrentLocale
  const localeMap = localeStore.getLocaleMap
  const defaultLocal = await import(`./lang/${locale.lang}.js`)

  const message = defaultLocal.default ?? {}

  setHtmlPageLang(locale.lang)

  localeStore.setCurrentLocale({
    lang: locale.lang || 'zh-CN',
  })

  return {
    legacy: false,
    locale: locale.lang,
    fallbackLocale: locale.lang,
    messages: {
      [locale.lang]: message
    },
    availableLocales: localeMap.map((v) => v.lang),
    sync: true,
    silentTranslationWarn: import.meta.env.PROD,
    missingWarn: !import.meta.env.PROD,
    silentFallbackWarn: import.meta.env.PROD
  }
}

export const setupI18n = async (app) => {
  if (i18n) return i18n

  const options = await createI18nOptions()
  
  i18n = createI18n(options)
  app.use(i18n)
}
