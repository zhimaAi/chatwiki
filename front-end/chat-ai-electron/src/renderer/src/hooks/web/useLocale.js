import { i18n } from '@/locales'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { setHtmlPageLang } from '@/locales/helper'

const setI18nLanguage = (locale) => {
  const localeStore = useLocaleStoreWithOut()

  if (i18n.mode === 'legacy') {
    i18n.global.locale = locale
  } else {
    i18n.global.locale.value = locale
  }
  localeStore.setCurrentLocale({
    lang: locale
  })
  setHtmlPageLang(locale)
}

export const useLocale = () => {
  // Switching the language will change the locale of useI18n
  // And submit to configuration modification
  const changeLocale = async (locale) => {
    const globalI18n = i18n.global

    const langModule = await import(`../../locales/lang/${locale}.js`)

    globalI18n.setLocaleMessage(locale, langModule.default)

    setI18nLanguage(locale)
  }

  return {
    changeLocale
  }
}
