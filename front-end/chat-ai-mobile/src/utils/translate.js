import { useChatStore } from '@/stores/modules/chat'


let wordDictionary = {
  '深度思考中...': 'Deep thinking...',
  '已完成深度思考': 'Deep thinking completed',
  '正在检索知识库...': 'Querying knowledge base...',
  '复制': 'Copy',
  '点赞': 'Like',
  '点踩': 'Dislike',
  '结果反馈': 'Result feedback',
  '请反馈你觉得回答不满意的地方': 'Please feedback what you think is unsatisfactory',
  '取消': 'Cancel',
  '提交': 'Submit',
  '发送': 'Send',
  '复制成功': 'Copy successful',
  '感谢反馈': 'Thank you for your feedback',
  '在此输入您想了解的内容': 'Type your message here',
  '在此输入您想了解的内容，Shift+Enter换行': 'Type your message here(Shift+Enter for new line)',
  '由 ChatWiki 提供软件支持': 'Powered by ChatWiki',
  '猜你想问': 'Suggestions',
  '常见问题': 'FAQ',
}

export function translate(text) {
  const chatStore = useChatStore()
  let lang = chatStore.externalConfigH5.lang
  if (lang == 'en-US') {
    return wordDictionary[text] || text
  }
  return text
}


