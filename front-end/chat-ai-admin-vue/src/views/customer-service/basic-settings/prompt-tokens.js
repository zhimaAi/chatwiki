export function getPromptOptions(variables = [], cards = []) {
  return {
    variables: [
      { id: 'system_current_time', label: '当前时间-系统变量', name: '当前时间', description: '系统变量', value: '【system_var:current_time】' },
      { id: 'system_current_date', label: '当前日期-系统变量', name: '当前日期', description: '系统变量', value: '【system_var:current_date】' },
      ...(variables || []).map(item => ({ ...item, label: item.variable_name || item.variable_key, name: item.variable_name || item.variable_key, description: item.variable_key, value: `【chat_variable:${item.variable_key}】` }))
    ],
    cards: (cards || []).map(item => ({ ...item, label: `小程序卡片：${item.title}`, value: `【mini_card:${item.id}】` }))
  }
}

export function tokenizePrompt(value, options = []) {
  const text = String(value ?? '')
  const labels = new Map(options.map(option => [option.value, option]))
  const parts = []
  let position = 0
  for (const match of text.matchAll(/【(system_var|chat_variable|mini_card):[^】\n]+】/g)) {
    if (match.index > position) parts.push({ text: text.slice(position, match.index) })
    const option = labels.get(match[0])
    const type = match[1] === 'mini_card' ? 'card' : 'variable'
    parts.push({ token: { ...option, value: match[0], type, label: option?.label || (type === 'card' ? '小程序卡片：已删除' : '未找到的变量') } })
    position = match.index + match[0].length
  }
  if (position < text.length) parts.push({ text: text.slice(position) })
  return parts
}

export function readPromptText(root) {
  let text = ''
  for (const node of root.childNodes) {
    if (node.nodeType === 3) text += node.textContent
    else if (node.classList?.contains('j-mention-at')) text += node.dataset.value || ''
    else if (node.nodeName === 'BR') text += '\n'
    else {
      const block = ['DIV', 'P'].includes(node.nodeName)
      if (block && text && !text.endsWith('\n')) text += '\n'
      const content = readPromptText(node)
      text += node.dataset?.promptCaret ? content.replace(/\u200b/g, '') : content
      if (block && node.nextSibling && !text.endsWith('\n')) text += '\n'
    }
  }
  return text
}

export function getPromptSelection(root) {
  const selection = root.ownerDocument.getSelection()
  if (!selection.rangeCount) return null
  const range = selection.getRangeAt(0)
  if (!root.contains(range.startContainer) || !root.contains(range.endContainer)) return null
  const offset = (node, position) => {
    const prefix = root.ownerDocument.createRange()
    prefix.selectNodeContents(root)
    prefix.setEnd(node, position)
    return readPromptText(prefix.cloneContents()).length
  }
  return { start: offset(range.startContainer, range.startOffset), end: offset(range.endContainer, range.endOffset) }
}

export function restorePromptSelection(root, offsets) {
  if (!offsets) return
  const point = offset => {
    for (const node of root.childNodes) {
      const token = node.classList?.contains('j-mention-at')
      const caret = node.dataset?.promptCaret
      const length = token ? node.dataset.value.length : caret ? 0 : node.textContent.length
      if (token && offset > 0 && offset <= length && node.nextSibling?.dataset?.promptCaret) return [node.nextSibling.firstChild, 1]
      if (caret && offset === 0) return [node.firstChild, 1]
      if (offset <= length) return token ? [root, [...root.childNodes].indexOf(node) + Number(offset > 0)] : [node, offset]
      offset -= length
    }
    return [root, root.childNodes.length]
  }
  const range = root.ownerDocument.createRange()
  range.setStart(...point(offsets.start))
  range.setEnd(...point(offsets.end))
  const selection = root.ownerDocument.getSelection()
  selection.removeAllRanges()
  selection.addRange(range)
}

export function createPromptToken(token, document) {
  const span = document.createElement('span')
  span.className = `j-mention-at${token.type === 'card' ? ' mini-card-tag' : ''}`
  span.setAttribute('contenteditable', 'false')
  span.dataset.value = token.value
  span.dataset.type = token.type || 'variable'
  if (token.type === 'card') {
    const icon = document.createElement('span')
    icon.className = 'svg-action'
    icon.innerHTML = '<svg aria-hidden="true" class="svg-icon" fill="currentColor" width="12" height="12"><use href="#applet"></use></svg>'
    span.append(icon)
  }
  span.append(document.createTextNode(token.label))
  return span
}

export function ensurePromptCaretSlots(root) {
  const createSlot = () => {
    const span = root.ownerDocument.createElement('span')
    span.dataset.promptCaret = 'true'
    span.append(root.ownerDocument.createTextNode('\u200b'))
    return span
  }
  for (const token of root.querySelectorAll('.j-mention-at')) {
    if (!token.previousSibling || (token.previousSibling.nodeType === 3 && !token.previousSibling.textContent)) token.before(createSlot())
    if (!token.nextSibling?.dataset?.promptCaret) token.after(createSlot())
  }
}

export function renderPromptText(root, value, options) {
  const offsets = root.ownerDocument.activeElement === root ? getPromptSelection(root) : null
  const fragment = root.ownerDocument.createDocumentFragment()
  for (const part of tokenizePrompt(value, options)) fragment.append(part.token ? createPromptToken(part.token, root.ownerDocument) : root.ownerDocument.createTextNode(part.text))
  root.replaceChildren(fragment)
  // 非编辑标签两侧保留文字落点；占位字符只存在于 DOM，不提交到提示词。
  ensurePromptCaretSlots(root)
  restorePromptSelection(root, offsets)
}
