export const REASONING_EFFORTS = ['none', 'minimal', 'low', 'medium', 'high', 'xhigh', 'max']
export const REASONING_EFFORT_MAX_LENGTH = 50

export function toReasoningEffort(enableThinking, reasoningEffort) {
  const enabled = enableThinking === true || Number(enableThinking) > 0
  return enabled ? reasoningEffort || 'medium' : 'none'
}

export function isValidReasoningEffort(value) {
  return typeof value === 'string' && Array.from(value).length <= REASONING_EFFORT_MAX_LENGTH
}
