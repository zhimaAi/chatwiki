import { getUuid, extractVoiceInfo, removeVoiceFormat } from '@/utils/index'

/**
 * @typedef {'finish' | 'stop' | 'close' | 'error'} ProcessFinalizeReason
 * @typedef {Object} NormalizedChatEvent
 * @property {'thinking_delta' | 'answer_delta' | 'answer_begin' | 'round_begin' | 'round_finish' | 'tool_start' | 'tool_finish' | 'operation' | 'process_finalize' | 'final_snapshot'} type
 * @property {string} [content]
 * @property {string} [source]
 * @property {'tool' | 'skill'} [stepType]
 * @property {string} [title]
 * @property {string} [paramsText]
 * @property {string} [toolCallId]
 * @property {string} [toolName]
 * @property {string} [result]
 * @property {string} [eventName]
 * @property {string} [rawData]
 * @property {ProcessFinalizeReason} [reason]
 * @property {Object} [message]
 */
export const safeParseJson = (value, fallback) => {
  if (typeof value !== 'string') {
    return value ?? fallback
  }

  try {
    return JSON.parse(value)
  } catch {
    return fallback
  }
}

export const ensureProcessStepState = (msg) => {
  if (!Array.isArray(msg.process_steps)) {
    msg.process_steps = []
  }
  if (typeof msg.current_round_index !== 'number') {
    msg.current_round_index = 0
  }
  if (typeof msg.active_thinking_step_id !== 'string') {
    msg.active_thinking_step_id = ''
  }
}

const createProcessStep = (step = {}) => ({
  id: step.id || getUuid(32),
  type: step.type || 'thinking',
  title: step.title || '',
  status: step.status || '',
  expanded: step.expanded === true,
  hidden: step.hidden === true,
  roundIndex: step.roundIndex || 0,
  contentText: step.contentText || '',
  paramsText: step.paramsText || '',
  tool_name: step.tool_name || '',
  resultText: step.resultText || '',
  eventName: step.eventName || '',
  tool_call_id: step.tool_call_id || '',
})

const appendProcessStep = (msg, step = {}) => {
  ensureProcessStepState(msg)
  const nextStep = createProcessStep(step)
  msg.process_steps.push(nextStep)
  return nextStep
}

const getProcessStepById = (msg, stepId) => {
  ensureProcessStepState(msg)
  return msg.process_steps.find((item) => item.id === stepId)
}

const finishActiveThinkingStep = (msg) => {
  const activeStep = getProcessStepById(msg, msg.active_thinking_step_id)
  if (activeStep) {
    if (!activeStep.contentText) {
      msg.process_steps = msg.process_steps.filter((item) => item.id !== activeStep.id)
    } else {
      activeStep.status = 'done'
      activeStep.resultText = ''
    }
  }
  msg.active_thinking_step_id = ''
}

const finalizeRunningProcessSteps = (msg, reason) => {
  const hasRunning = msg.process_steps?.some((step) => step.status === 'running')
  if (msg.process_finish_reason && !hasRunning) return
  reason = reason === 'close' ? 'stop' : reason
  msg.process_finish_reason = reason
  msg.is_stopped = reason === 'stop'
  ensureProcessStepState(msg)
  msg.process_steps = msg.process_steps
    .filter((step) => !(step.type === 'thinking' && step.status === 'running' && !step.contentText))
    .map((step) => {
      if (step.status !== 'running') {
        return step
      }
      return {
        ...step,
        status: reason === 'stop' ? 'stopped' : reason === 'error' ? 'error' : 'done',
        resultText: step.type === 'thinking' && !step.contentText ? '' : step.resultText,
      }
    })
  msg.active_thinking_step_id = ''
  msg.loading = false
  msg.startLoading = false
  msg.quote_loading = false
}

/**
 * 将后端原始 SSE 事件转换为页面无关的过程事件。
 *
 * 后端根据 application_type 选择旧链路或新链路，正常情况下两套正文事件互斥：
 * - reasoning_content -> thinking_delta
 * - stream_message.reasoning_content -> thinking_delta
 * - sending -> answer_delta
 * - stream_message.content -> answer_delta
 * - llm_rounds:begin / finish -> round_begin / round_finish
 * - tool_call_full / tool_result -> tool_start / tool_finish
 * - FileOperation、ExecuteCommand -> operation
 * - finish、用户停止、SSE 关闭 -> process_finalize
 * - ai_message -> final_snapshot
 *
 * stream_message 可能同时包含思考和正文，因此返回事件数组并保证先思考、后正文。
 * 首段非空正文表示当前思考阶段结束；空正文、空思考和元数据帧不改变过程状态。
 * ai_message 转为 final_snapshot，继续作为服务端持久化后的最终权威快照。
 * debug、耗时、引用和会话变量属于管理端附属状态，不进入过程状态机。
 *
 * @param {string} type 原始事件名
 * @param {*} content 原始事件数据
 * @returns {NormalizedChatEvent[]}
 */
export const normalizeSseEvent = (type, content) => {
  if (type === 'reasoning_content') {
    return [{ type: 'thinking_delta', content: String(content || ''), source: 'reasoning_content' }]
  }

  if (type === 'sending') {
    return [{ type: 'answer_delta', content: String(content || ''), source: 'sending' }]
  }

  if (type === 'stream_message') {
    const data = safeParseJson(content, {})
    if (data?.role && data.role !== 'assistant') return []
    const events = []
    if (data?.reasoning_content) {
      events.push({
        type: 'thinking_delta',
        content: String(data.reasoning_content),
        source: 'stream_message',
      })
    }
    if (data?.content) {
      events.push({
        type: 'answer_delta',
        content: String(data.content),
        source: 'stream_message',
      })
    }
    return events
  }

  if (type === 'llm_rounds') {
    if (content === 'begin') {
      return [{ type: 'round_begin' }]
    }
    if (content === 'finish') {
      return [{ type: 'round_finish' }]
    }
    return []
  }

  if (type === 'tool_call_full') {
    const functionInfo = content?.function || {}
    const isSkill = functionInfo.name === 'skill'
    const functionArgs = isSkill ? safeParseJson(functionInfo.arguments, {}) : {}
    return [{
      type: 'tool_start',
      stepType: isSkill ? 'skill' : 'tool',
      title: isSkill ? functionArgs?.skill || 'skill' : functionInfo.name || 'tool',
      paramsText: functionInfo.arguments || '',
      toolName: functionInfo.name || '',
      toolCallId: content?.id || '',
    }]
  }

  if (type === 'tool_result') {
    const data = safeParseJson(content, {})
    return [{
      type: 'tool_finish',
      toolCallId: data?.tool_call_id || '',
      toolName: data?.tool_name || '',
      result: data?.content || '',
    }]
  }

  if (type === 'process_event') {
    return [{
      type: 'operation',
      eventName: content?.eventName || '',
      rawData: content?.rawData || '',
    }]
  }

  if (type === 'finalize_process_steps') {
    return [{ type: 'process_finalize', reason: content?.reason || 'finish' }]
  }

  if (type === 'ai_message') {
    return [{ type: 'final_snapshot', message: content }]
  }

  return []
}

/**
 * @param {Object} msg 当前机器人消息
 * @param {NormalizedChatEvent} event 统一过程事件
 */
export const applyNormalizedChatEvent = (msg, event) => {
  if (event.type === 'thinking_delta') {
    if (!event.content) {
      return
    }
    ensureProcessStepState(msg)
    let activeStep = getProcessStepById(msg, msg.active_thinking_step_id)
    if (!activeStep) {
      // 旧链路没有 llm_rounds:begin，首个思考增量至少归入第 1 轮。
      msg.current_round_index = Math.max(1, msg.current_round_index || 0)
      activeStep = appendProcessStep(msg, {
        type: 'thinking',
        status: 'running',
        expanded: false,
        roundIndex: msg.current_round_index,
        eventName: event.source,
      })
      msg.active_thinking_step_id = activeStep.id
    }
    activeStep.contentText = `${activeStep.contentText || ''}${event.content}`
    activeStep.resultText = ''
    // 原字段只用于接口和历史兼容，页面统一从 process_steps 渲染。
    msg.reasoning_content = `${msg.reasoning_content || ''}${event.content}`
    return
  }

  if (event.type === 'answer_begin') {
    msg.startLoading = false
    finishActiveThinkingStep(msg)
    return
  }

  if (event.type === 'answer_delta') {
    // 空正文及元数据帧不能结束思考，首个非空正文才代表当前思考阶段完成。
    if (!event.content) {
      return
    }
    msg.startLoading = false
    finishActiveThinkingStep(msg)
    msg.content = `${msg.content || ''}${event.content}`
    msg.voice_content = extractVoiceInfo(msg.content)
    msg.content = removeVoiceFormat(msg.content)
    return
  }

  if (event.type === 'round_begin') {
    ensureProcessStepState(msg)
    msg.current_round_index += 1
    return
  }

  if (event.type === 'round_finish') {
    finishActiveThinkingStep(msg)
    return
  }

  if (event.type === 'tool_start') {
    appendProcessStep(msg, {
      type: event.stepType,
      title: event.title,
      status: 'running',
      expanded: false,
      roundIndex: msg.current_round_index || 0,
      paramsText: event.paramsText,
      tool_call_id: event.toolCallId,
      tool_name: event.toolName,
      eventName: 'tool_call_full',
    })
    return
  }

  if (event.type === 'tool_finish') {
    ensureProcessStepState(msg)
    const candidates = msg.process_steps.filter((step) => {
      return ['tool', 'skill'].includes(step.type) && step.status === 'running' &&
        (!event.toolName || (step.tool_name || step.title) === event.toolName)
    })
    const matchedStep = event.toolCallId
      ? msg.process_steps.find((step) => {
          return ['tool', 'skill'].includes(step.type) &&
            step.status === 'running' &&
            step.tool_call_id === event.toolCallId
        })
      : candidates.length === 1 ? candidates[0] : null

    // 结果无法关联时不补造步骤，避免错误完成其他并行工具。
    if (matchedStep) {
      matchedStep.status = 'done'
      matchedStep.resultText = event.result
    }
    return
  }

  if (event.type === 'operation') {
    // 文件和命令事件只保留审计数据，前端不展示也不会执行其中内容。
    appendProcessStep(msg, {
      type: 'operation',
      title: event.eventName,
      status: '',
      hidden: true,
      expanded: false,
      roundIndex: msg.current_round_index || 0,
      paramsText: event.rawData,
      eventName: event.eventName,
    })
    return
  }

  if (event.type === 'process_finalize') {
    // 正常完成、用户停止和断流统一清理 running，避免页面残留 spinner。
    finalizeRunningProcessSteps(msg, event.reason)
    return
  }

  if (event.type === 'final_snapshot') {
    const snapshot = event.message || {}
    if (!restoreAgentProcess(msg, { debugLog: snapshot.debug_log })) finalizeRunningProcessSteps(msg, 'finish')
    msg.loading = false
    msg.startLoading = false
    msg.quote_loading = false
    if (snapshot.debug_log !== undefined) msg.debug_log = snapshot.debug_log
    if (snapshot.menu_json && snapshot.msg_type == 2) {
      const menuJson = safeParseJson(snapshot.menu_json, {})
      msg.question = menuJson.question || []
    }
    msg.id = snapshot.id
    msg.message_id = snapshot.message_id || snapshot.id || msg.message_id
    if (snapshot.msg_type !== undefined) {
      msg.msg_type = snapshot.msg_type
    }
    msg.content = snapshot.content || ''
    msg.voice_content = extractVoiceInfo(msg.content)
    msg.content = removeVoiceFormat(msg.content)
    if (snapshot.reply_content_list !== undefined) {
      msg.reply_content_list = safeParseJson(snapshot.reply_content_list, snapshot.reply_content_list || [])
    }
    if (snapshot.quote_file !== undefined) {
      msg.quote_file = safeParseJson(snapshot.quote_file, snapshot.quote_file || [])
    }
  }
}

/**
 * 旧历史机器人文本消息没有 process_steps 时，将 reasoning_content 转成确定性的完成步骤。
 * 已有步骤、用户消息、菜单、图片和空思考均不转换，避免分页或重进时重复生成。
 */
export const normalizeHistoricalProcessSteps = (msg) => {
  if (restoreAgentProcess(msg, { historical: true })) return msg.process_steps
  const parsedSteps = Array.isArray(msg.process_steps)
    ? msg.process_steps
    : safeParseJson(msg.process_steps, [])
  if (Array.isArray(parsedSteps) && parsedSteps.length > 0) {
    return parsedSteps
  }

  const reasoningContent = typeof msg.reasoning_content === 'string' ? msg.reasoning_content : ''
  if (msg.is_customer == 1 || msg.msg_type != 1 || !reasoningContent.trim()) {
    return []
  }

  return [createProcessStep({
    id: `legacy-thinking-${msg.message_id || msg.id || msg.uid}`,
    type: 'thinking',
    status: 'done',
    expanded: false,
    hidden: false,
    roundIndex: 1,
    contentText: reasoningContent,
    resultText: '',
    eventName: 'reasoning_content',
  })]
}

// New process logs are intentionally independent of legacy prompt/tool debug entries.
export const restoreAgentProcess = (msg, { historical = false, debugLog = msg.debug_log } = {}) => {
  if (msg.is_customer == 1 || msg.msg_type != 1) return false
  const logs = safeParseJson(debugLog, [])
  if (!Array.isArray(logs)) return false
  const processLogs = logs.filter((item) => item?.type === 'agent_process' &&
    typeof item.content === 'string' && ['thinking', 'tool_call', 'tool_result', 'finish'].includes(item.event))
  if (!processLogs.length) return false

  const replay = { process_steps: [], content: '', reasoning_content: '' }
  processLogs.forEach((item) => {
    const roundIndex = Number(item.round_index)
    if (Number.isInteger(roundIndex) && roundIndex > 0) replay.current_round_index = roundIndex
    if (item.event === 'thinking') {
      applyNormalizedChatEvent(replay, { type: 'thinking_delta', content: item.content, source: 'stream_message' })
      applyNormalizedChatEvent(replay, { type: 'round_finish' })
    } else if (item.event === 'tool_call') {
      normalizeSseEvent('tool_call_full', {
        id: item.tool_call_id,
        function: { name: item.tool_name, arguments: item.content }
      }).forEach((event) => applyNormalizedChatEvent(replay, event))
    } else if (item.event === 'tool_result') {
      applyNormalizedChatEvent(replay, {
        type: 'tool_finish', toolCallId: item.tool_call_id, toolName: item.tool_name, result: item.content
      })
    } else if (['finish', 'stop', 'error'].includes(item.content)) {
      applyNormalizedChatEvent(replay, { type: 'process_finalize', reason: item.content })
    }
  })
  const finishReason = replay.process_finish_reason || 'finish'
  if (!replay.process_finish_reason) {
    applyNormalizedChatEvent(replay, { type: 'process_finalize', reason: finishReason })
  }

  const previous = (Array.isArray(msg.process_steps) ? msg.process_steps : [])
    .filter((step) => step && typeof step === 'object')
  replay.process_steps.forEach((step, index) => {
    const matchedIndex = previous.findIndex((old) => step.tool_call_id
      ? old.tool_call_id === step.tool_call_id
      : old.type === step.type && old.roundIndex === step.roundIndex)
    const old = matchedIndex >= 0 ? previous.splice(matchedIndex, 1)[0] : null
    step.id = old?.id || `agent-process-${msg.message_id || msg.id || msg.uid}-${index}`
    step.expanded = historical ? false : old?.expanded === true
  })
  msg.process_steps = replay.process_steps
  msg.current_round_index = replay.current_round_index || 0
  msg.active_thinking_step_id = ''
  msg.process_finish_reason = finishReason
  msg.is_stopped = finishReason === 'stop'
  if (historical) msg.process_expanded = false
  return true
}

export const hasVisibleProcessSteps = (msg) => Array.isArray(msg?.process_steps) &&
  msg.process_steps.some((step) => step && step.hidden !== true &&
    (step.title || step.contentText || step.paramsText || step.resultText || step.status === 'running'))
