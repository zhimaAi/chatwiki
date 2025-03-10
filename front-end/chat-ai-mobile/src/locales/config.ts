export const localeMap = {
  'zh-CN': 'zh-CN',
  'en-US': 'en-US'
} as const


export type LocaleType = keyof typeof localeMap

export interface Language {
  el: Recordable
  name: string
}

export interface LocaleDropdownType {
  lang: LocaleType
  name?: string
  vantLocale?: Language
}

export interface LocaleState {
  currentLocale: LocaleDropdownType
  localeMap: LocaleDropdownType[]
}

export const localeList = [
  {
    lang: localeMap['en-US'],
    label: 'English',
    icon: '🇺🇸',
    title: 'Language'
  },
  {
    lang: localeMap['zh-CN'],
    label: '简体中文',
    icon: '🇨🇳',
    title: '语言'
  }
] as const
