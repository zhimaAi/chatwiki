export function cloneConfig(value) {
  if (value === null || typeof value !== 'object' || value instanceof File) return value
  if (Array.isArray(value)) return value.map(cloneConfig)
  return Object.fromEntries(Object.entries(value).map(([key, item]) => [key, cloneConfig(item)]))
}

export function serializeRobotConfig(state) {
  const data = cloneConfig(state)
  if (data.robot_avatar instanceof File) delete data.robot_avatar_url
  else if (data.robot_avatar) {
    data.robot_avatar_url = data.robot_avatar
    delete data.robot_avatar
  }
  for (const key of ['welcomes', 'unknown_question_prompt']) {
    if (typeof data[key] !== 'string') {
      const config = data[key] || { content: '', question: [] }
      data[key] = JSON.stringify({ ...config, question: (config.question || []).map(item => typeof item === 'string' ? item : item.content) })
    }
  }
  for (const key of ['prompt_struct', 'cache_config']) {
    if (typeof data[key] !== 'string') data[key] = JSON.stringify(data[key] || {})
  }
  return data
}

// 保存失败不能中断后续重试；发送测试消息必须等待所有已提交的保存。
export function createSaveQueue() {
  let tail = Promise.resolve()
  let lastError = null
  return {
    enqueue(task) {
      const result = tail.then(task).then(value => { lastError = null; return value }, error => { lastError = error; throw error })
      tail = result.catch(() => {})
      return result
    },
    async flush() {
      let current
      do { current = tail; await current } while (current !== tail)
      if (lastError) throw lastError
    }
  }
}
