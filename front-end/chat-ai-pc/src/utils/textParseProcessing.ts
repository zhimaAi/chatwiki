// 清除 markdown 标记(...)
function clearMarkdownTag(str: string): string {
  // 过滤点开头的```markdown 和结束的```并保留中间的内容
  str = str.replace(/```markdown([\s\S]*)```/g, (_, content) => {
    console.log(content)
    return content.trim()
  })
  // 过滤掉没有```结束符的```markdown
  str = str.replace(/```markdown/g, '')
  return str
}

function textParseProcessing(str: string):string {
  str = clearMarkdownTag(str)
  return str
}

export default textParseProcessing
