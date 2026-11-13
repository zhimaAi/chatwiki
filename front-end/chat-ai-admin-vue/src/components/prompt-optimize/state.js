export function mergeOlderRecords(current, incoming) {
  const seen = new Set(current.map((row) => String(row.id)))
  return [...incoming]
    .reverse()
    .filter((row) => !seen.has(String(row.id)))
    .concat(current)
}

export function composeTimeline(history, current) {
  return [...history, { id: 'guide', guide: true }, ...current]
}

export function findSlash(text, caret) {
  const beforeCaret = text.slice(0, caret)
  const token = beforeCaret.match(/\S*$/)?.[0] || ''
  if (token.includes('://')) return null
  const match = beforeCaret.match(/\/([^\s/]*)$/)
  return match ? { query: match[1], start: caret - match[1].length - 1 } : null
}

export function insertPlainName(text, start, end, name) {
  return { text: text.slice(0, start) + name + text.slice(end), caret: start + name.length }
}
