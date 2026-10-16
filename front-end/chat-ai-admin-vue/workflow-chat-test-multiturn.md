# 工作流聊天测试多轮对话 - 前端改动文档

> 关联需求：【chatwiki】工作流内调试支持进行多轮对话
> 后端改动已完成（见本文「接口契约」），本文档描述前端需要配合的改动。

## 一、背景与目标

- 工作流聊天测试（`/chat/callWorkFlowDialog`）当前每轮独立：正常结束后前端清空会话标识并清空消息列表，页面不留历史。
- LLM 类节点（AI 对话 / 问题分类 / 问题优化 / 参数提取）的【上下文数量】（`context_pair`）在调试模式下取不到历史（临时会话 ID 不落库），多轮上下文恒为空。
- 目标：
  1. 聊天测试支持连续多轮对话，页面记录全部多轮消息；
  2. 每轮请求携带已完成轮次的问答历史，后端按各节点【上下文数量】裁剪注入。

## 二、涉及文件

仅一个文件：

```
src/views/workflow/components/run-test/ChatTestPanel.vue
```

`src/api/robot/index.js` 的 `callWorkFlowDialog` 透传 `data`，无需修改。

## 三、接口契约（后端已就绪）

`POST /chat/callWorkFlowDialog` 新增表单字段 `context_list`：

- 类型：JSON 数组的字符串形式，按时间正序，仅含**已完成轮次**（当前这轮的问题仍走 `question` 字段）：
  ```json
  [{"question": "第一轮问题", "answer": "第一轮回答"}, {"question": "...", "answer": "..."}]
  ```
- 后端容错规则（前端无需重复处理，但不要发送明显非法数据）：
  - 非法 JSON / 非数组 → 忽略，降级为单轮；
  - 最多 50 对（与节点 `context_pair` 上限一致），超出取最近 50 对；
  - 单条 question/answer 超过 10000 字符截断；空 question 或空 answer 的项丢弃。

## 四、改动一：会话标识常驻

`handleSend` 中，原逻辑仅在问答节点暂停时延续会话；改为只要会话标识存在就延续：

```js
// 原：
const shouldContinue = isAwaitingFollowUp.value && dialogue_id.value && session_id.value && testUserId.value
if (!shouldContinue) {
  resultList.value = []
  currentNodeKey.value = ''
  dialogue_id.value = ''
  session_id.value = ''
  testUserId.value = generateRandomId(16)
  messageList.value = []
}
isAwaitingFollowUp.value = false

// 改为：
const shouldContinue = !!(dialogue_id.value && session_id.value && testUserId.value)
if (!shouldContinue) {
  // 首次或出错后：全新会话
  resultList.value = []
  currentNodeKey.value = ''
  dialogue_id.value = ''
  session_id.value = ''
  testUserId.value = generateRandomId(16)
  messageList.value = []
} else if (!isAwaitingFollowUp.value) {
  // 新一轮运行：清空上一轮节点日志，保留聊天记录与会话标识
  resultList.value = []
  currentNodeKey.value = ''
}
isAwaitingFollowUp.value = false
```

`finalizeCurrentRun` 中删除正常结束时的清空分支（会话标识常驻）：

```js
// 原：
isAwaitingFollowUp.value = lastNode?.node_type === 43
if (!isAwaitingFollowUp.value) {
  dialogue_id.value = ''
  session_id.value = ''
  testUserId.value = ''
}

// 改为：
isAwaitingFollowUp.value = lastNode?.node_type === 43
```

语义说明：

- `isAwaitingFollowUp` 不再表示「能否继续对话」，仅用于区分**同轮恢复**（问答节点暂停后续跑，保留节点日志）与**新一轮**（清空节点日志）。
- 会话标识（`dialogue_id` / `session_id` / `testUserId`）在整个面板会话期间常驻。正常结束后复用同一组 ID 再发消息时，后端查无暂停快照，会从开始节点全新运行，行为正确。
- `onError` 回调仍清空会话标识：出错后下次发送即全新会话。

## 五、改动二：发送时携带多轮上下文

新增三个辅助函数：

```js
// 与后端 GetFirstQuestionByInput 对齐：多模态 JSON 取第一个文本项
const extractQuestionText = (content) => {
  if (typeof content !== 'string' || !content) return ''
  const text = content.trim()
  if (!text.startsWith('[') && !text.startsWith('{')) return text
  try {
    const items = JSON.parse(text)
    if (Array.isArray(items)) {
      const textItem = items.find((it) => it?.type === 'text' && it.text)
      if (textItem) return textItem.text
      const media = items.find((it) => ['image_url', 'input_audio', 'video_url'].includes(it?.type))
      if (media) return `[${media.type === 'image_url' ? '图片' : media.type === 'input_audio' ? '语音' : '视频'}]`
    }
  } catch {
    /* 非多模态 JSON，按原文返回 */
  }
  return text
}

// 与后端 takeReplyText 对齐（Description 优先），智能菜单补充 menu_description
const extractAnswerText = (msg) => {
  if (msg?.content) return String(msg.content)
  const list = Array.isArray(msg?.reply_content_list) ? msg.reply_content_list : []
  for (const item of list) {
    if (item?.description) return item.description
    if (item?.title) return item.title
    if (item?.smart_menu?.menu_description) return item.smart_menu.menu_description
  }
  return ''
}

const buildContextList = () => {
  const pairs = []
  let question = ''
  messageList.value.forEach((msg) => {
    if (msg.is_customer === 1) {
      question = extractQuestionText(msg.content)
      return
    }
    if (!question) return
    const answer = extractAnswerText(msg)
    if (answer) {
      pairs.push({ question, answer })
    }
    question = ''
  })
  return pairs.slice(-50) // 与后端 context_pair 上限一致
}
```

在 `handleSend` 中，于 `messageList.value.push(userMessage, aiMessage)` **之前**构建历史（保证不含当前轮），并在请求参数中加入 `context_list`：

```js
const contextList = buildContextList()
// ... push 当前 userMessage / aiMessage 之后：
const sse = callWorkFlowDialog({
  robot_key: query.robot_key,
  openid: testUserId.value,
  question,
  dialogue_id: dialogue_id.value,
  session_id: session_id.value,
  context_list: JSON.stringify(contextList)
})
```

说明：

- 历史只取文本（多模态消息取文本部分，媒体项以 `[图片]` 等占位），与正式对话的 `context_qa` 行为一致。
- 智能菜单轮次的回答文本取 `smart_menu.menu_description`（问答节点的问题文案），保证暂停轮也能进入上下文。

## 六、改动三：页面记录多轮

`messageList` 全程累积（改动一的分支已保证新一轮不清空）；仅两处重置：

1. `open()`：打开测试面板时整体重置（现有逻辑，保留）；
2. SSE `onError`：出错重置（现有逻辑，保留）。

无需新增 UI 或清空按钮。

## 七、验收标准

1. 连续多轮提问，页面保留全部聊天记录；节点日志面板每轮独立（新一轮清空，问答节点恢复时保留）。
2. 第二轮起，请求 form data 中包含 `context_list`，question/answer 成对且不含当前轮问题。
3. 节点【上下文数量】= 0：调试日志 `llm_debug_log` 无 `context_qa`；= 1：仅最近一轮；= N：最多 N 对。
4. 问答节点暂停 → 点选智能菜单选项：恢复执行，且上下文包含暂停那一轮。
5. 出错后重发、或重新打开测试面板：均为全新会话（页面清空、无历史上下文）。
