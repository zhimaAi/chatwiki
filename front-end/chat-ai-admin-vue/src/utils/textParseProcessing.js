// Strip ```markdown ... ``` wrapper while keeping the content.
function clearMarkdownTag(str) {
  str = String(str || '')
  str = str.replace(/```markdown([\s\S]*?)```/g, (_, content) => content.trim())
  str = str.replace(/```markdown/g, '')
  return str
}

function escapeAttr(value) {
  return String(value || '')
    .replace(/&/g, '&amp;')
    .replace(/"/g, '&quot;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}

function isVideoUrl(url) {
  const val = String(url || '').trim()
  if (!val) return false
  const isSafeSrc =
    /^https?:\/\//i.test(val) ||
    /^\/(?!\/)/.test(val)
  if (!isSafeSrc) return false
  return (
    /\.(mp4|webm|ogg)(\?|#|$)/i.test(val) ||
    /library_video/i.test(val) ||
    /\/video\//i.test(val)
  )
}

function buildVideoTag(url, poster) {
  const src = escapeAttr(url)
  const posterAttr = poster ? ` poster="${escapeAttr(poster)}"` : ''
  return `<video class="markdown-video" controls preload="metadata" playsinline webkit-playsinline src="${src}"${posterAttr}></video>`
}

// Support:
// 1) ![video](https://xx/video.mp4)
// 2) ![video](https://xx/video.mp4 "poster-url")
function replaceVideoSyntax(str) {
  str = str.replace(/!\[([^\]]*)\]\(([^)\s]+)(?:\s+"([^"]*)")?\)/g, (match, alt, url, poster = '') => {
    const isVideoToken = String(alt || '').trim().toLowerCase() === 'video'
    return isVideoToken && isVideoUrl(url) ? buildVideoTag(url, poster) : match
  })

  return str
}

function textParseProcessing(str) {
  str = clearMarkdownTag(str)
  str = replaceVideoSyntax(str)
  return str
}

export default textParseProcessing
