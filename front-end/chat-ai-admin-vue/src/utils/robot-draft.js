const languageFields = [
  'welcomes',
  'unknown_question_prompt',
  'tips_before_answer_switch',
  'tips_before_answer_content',
  'enable_common_question',
  'common_question_list'
]

export const resolveRobotDraft = (published, useDraft = true, lang = 'zh-CN') => {
  const hasDraft = published.has_draft === true || published.has_draft === 'true'
  const draft = useDraft && hasDraft ? published.draft_data : null
  const info = { ...published, ...draft?.robot, has_draft: hasDraft }
  const publishedConfigs = typeof published.multi_lang_configs === 'string'
    ? JSON.parse(published.multi_lang_configs || '[]')
    : published.multi_lang_configs || []
  const configs = new Map(publishedConfigs.map((item) => [item.lang_key, { ...item }]))

  for (const item of draft?.lang_configs || []) {
    configs.set(item.lang_key, { ...configs.get(item.lang_key), ...item })
  }

  // 草稿多语言开关是 boolean，现有表单使用字符串 checkedValue。
  info.multi_lang_configs = [...configs.values()].map((item) => {
    for (const key of ['tips_before_answer_switch', 'enable_common_question']) {
      if (typeof item[key] === 'boolean') item[key] = String(item[key])
    }
    return item
  })

  const currentConfig = info.multi_lang_configs.find((item) => item.lang_key === lang)
  for (const key of languageFields) {
    if (currentConfig?.[key] !== undefined) info[key] = currentConfig[key]
  }

  return info
}
