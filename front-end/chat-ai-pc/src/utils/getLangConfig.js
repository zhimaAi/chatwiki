import { useChatStore } from '@/stores/modules/chat'
export function getLang() {
  const chatStore = useChatStore()
  let externalConfigPC = chatStore.externalConfigPC
  let lang = externalConfigPC?.lang || 'zh-CN'
  if (lang === 'auto') {
    lang = navigator.language || navigator.userLanguage
    if(lang == 'zh'){
      lang = 'zh-CN'
    }
  }
  return lang
}

export function getCurrentConfig(multi_lang_configs) {
  let lang = getLang()
  let list = multi_lang_configs ? JSON.parse(multi_lang_configs) : []
  return list.find((item) => item.lang_key === lang) || list[1] || list[0]
}
