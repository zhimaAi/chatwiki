// 清除 markdown 标记(...)
function clearMarkdownTag(str: string): string {
  const input = String(str || '')
  // 过滤点开头的```markdown 和结束的```并保留中间的内容
  let out = input.replace(/```markdown([\s\S]*?)```/g, (_, content) => {
    return String(content || '').trim()
  })
  // 过滤掉没有```结束符的```markdown
  out = out.replace(/```markdown/g, '')
  return out
}

function escapeAttr(value: string): string {
  return String(value || '')
    .replace(/&/g, '&amp;')
    .replace(/"/g, '&quot;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}

function isVideoUrl(url: string): boolean {
  const val = String(url || '').trim()
  if (!val) return false

  const isSafeSrc = /^https?:\/\//i.test(val) || /^\/(?!\/)/.test(val)
  if (!isSafeSrc) return false

  return (
    /\.(mp4|webm|ogg)(\?|#|$)/i.test(val) ||
    /library_video/i.test(val) ||
    /\/video\//i.test(val)
  )
}

function buildVideoTag(url: string, poster = ''): string {
  const src = escapeAttr(url)
  const posterAttr = poster ? ` poster="${escapeAttr(poster)}"` : ''
  return `<video class="markdown-video" controls preload="metadata" playsinline webkit-playsinline src="${src}"${posterAttr}></video>`
}

// Support:
// 1) ![video](https://xx/video.mp4)
// 2) ![video](https://xx/video.mp4 "poster-url")
function replaceVideoSyntax(str: string): string {
  return str.replace(/!\[([^\]]*)\]\(([^)\s]+)(?:\s+"([^"]*)")?\)/g, (match, alt, url, poster = '') => {
    const isVideoToken = String(alt || '').trim().toLowerCase() === 'video'
    return isVideoToken && isVideoUrl(url) ? buildVideoTag(url, poster) : match
  })
}

function textParseProcessing(str: string): string {
  let out = clearMarkdownTag(str)
  out = replaceVideoSyntax(out)
  return out
}

export default textParseProcessing
