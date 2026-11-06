// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

import { ref } from 'vue'

// Preview links are fetched when needed and are never added to saved messages.
export function useAnswerSourcePreview(getAnswerSource) {
  const sources = ref([])
  const sourceUrl = ref('')
  const viewOpen = ref(false)
  let context = null
  let generation = 0
  let previewRequest = 0

  const reset = () => {
    generation++
    previewRequest++
    context = null
    sources.value = []
    sourceUrl.value = ''
    viewOpen.value = false
  }

  const refreshSources = async () => {
    if (!context) return null
    const messageId = Number(context.message_id)
    if (!Number.isSafeInteger(messageId) || messageId <= 0) return null
    const currentGeneration = generation
    try {
      const response = await getAnswerSource({ ...context })
      if (generation !== currentGeneration) return null
      const latest = Array.isArray(response?.data) ? response.data : []
      // Keep streamed text visible while its answer sources are being saved.
      if (latest.length) sources.value = latest
      return latest
    } catch (_error) {
      return null
    }
  }

  const setContext = (params, preloaded = []) => {
    reset()
    context = {
      message_id: params.message_id,
      file_id: params.file_id,
      robot_key: params.robot_key,
      openid: params.openid,
    }
    sources.value = Array.isArray(preloaded) ? preloaded : []
    if (sources.value.length) return Promise.resolve(sources.value)
    return refreshSources()
  }

  const previewSource = async (item) => {
    const currentGeneration = generation
    const currentRequest = ++previewRequest
    sourceUrl.value = ''
    viewOpen.value = false
    const latest = await refreshSources()
    if (!latest || currentGeneration !== generation || currentRequest !== previewRequest) return false
    const paragraphId = item.id || item.paragraph_id
    const source = latest.find(candidate =>
      Number(paragraphId) > 0 &&
      String(candidate.id) === String(paragraphId) &&
      Number(candidate.page_num) > 0 &&
      Number(candidate.page_num) === Number(item.page_num)
    )
    if (!source || typeof source.preview_url !== 'string' || !source.preview_url) return false
    sourceUrl.value = source.preview_url
    viewOpen.value = true
    return true
  }

  return { sources, sourceUrl, viewOpen, setContext, previewSource, reset }
}
