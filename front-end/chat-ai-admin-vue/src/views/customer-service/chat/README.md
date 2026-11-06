# 客服对话页

入口为 `/customer-service/chat`。页面使用客服公共布局当前选中的应用，支持聊天机器人（类型0）和智能体（类型2），两者均留在客服模块。交互参考 `src/views/clawbot/chat/index.vue`，输入区不提供 Skill 入口。新建会话显示 Hero，不插入欢迎消息；服务端实际返回的工具或 Skill 过程仍可在消息中展示。

当前实现包含会话列表与历史分页、文本/图片发送、流式消息与停止、模型和人设正式保存、会话变量及知识库入口。未新增会话搜索、重命名、删除、消息赞踩或重新生成。旧 `/robot/test` 不作为本页实现依据。

## 文件与状态归属

以下路径均相对项目根目录。

| 文件 | 职责与维护边界 |
| --- | --- |
| `src/stores/modules/customer-service.js` | `useCustomerServiceStore` 持有已校验的当前应用和正式配置；列表与详情只接受类型0/2。URL身份、默认选择和应用切换由客服公共布局负责。 |
| `src/stores/modules/customer-service-chat.js` | `useCustomerServiceChatStore` 持有会话身份、消息、列表、分页、变量、发送状态及自己的 SSE 连接；只读基础 Store 身份，不选择或创建机器人。 |
| `src/views/customer-service/chat/index.vue` | 编排页面、草稿/附件、弹层、滚动、跳转、配置保存反馈与异步结果接纳。 |
| `src/views/customer-service/chat/components/message-input.vue` | 受控输入、图片上传、模型选择与快捷入口；通过事件请求保存或跳转。 |
| `src/views/customer-service/chat/components/prompt-drawer.vue` | 人设查看/编辑、变量与卡片辅助操作；保存通过事件交给页面。 |
| `src/views/customer-service/chat/components/variable-modal.vue` | 会话变量本地副本及必填校验；提交通过事件交给页面。 |
| `src/locales/lang/{zh-CN,en-US}/views/customer-service/chat.json` | 客服局部提示与错误文案；共用组件继续使用其既有文案。 |

当前身份以基础 Store 的 `robotStatus === 'ready'` 及已校验 `robot.id/robot_key` 为准，来源为客服 URL 的 `id/robot_key`。聊天页不从 Agent 当前项或浏览器缓存另选机器人。业务状态不使用 `src/stores/modules/clawbot.js`、`clawbot-chat.js`、`robot.js`、`chat.js` 的实例，也不嵌入 Agent 整页。

两个客服 Store 均无 Pinia 持久化配置。首轮变量沿既有协议使用 `chat_prompt_variables_${robot_key}` 本地缓存，首次接纳发送后移除；草稿和附件不持久化。`openid` 优先使用页面传入值，否则沿既有 `getOpenid()` 获取。

## 只读共享依赖

| 准确路径 | 使用方式 |
| --- | --- |
| `src/views/clawbot/chat/components/chat-list.vue` | 会话日期分组、选中和滚动事件。 |
| `src/views/clawbot/chat/components/message-list.vue` 及其消息展示叶子组件 | 消息、过程步骤、引用、菜单及滚动定位。 |
| `src/views/clawbot/chat/components/conversation-bar.vue` | 当前会话与知识库展示。 |
| `src/views/clawbot/chat/components/scroll-bottom-btn.vue` | 滚动到底按钮。 |
| `src/views/clawbot/chat/components/chat-model-select.vue` | 由局部输入区传正式配置，保留既有模型选项交互。 |
| `src/views/clawbot/chat/components/file-toolbar.vue` | 由局部输入区传附件展示。 |
| `src/utils/agent-process.js` | 纯事件归一化、过程步骤更新与历史步骤恢复。 |

上述组件在当前使用范围内通过 props/events 或展示接口工作，没有直接读取 Agent 业务 Store；引用它们仍形成展示层共享依赖。后续修改其接口或深层依赖时应同时审阅客服和 Agent 消费者，不能据目录归属认定可随意修改。

## 聊天 Store 契约

主要状态：`chatStatus` 为 `idle/loading/ready/error`，`chatLoading` 为派生值；`sendStatus` 为 `idle/checking/streaming`，`sendLock` 覆盖检查至连接结束。会话身份使用 `dialogue_id/openid/rel_user_id`，变量使用完整 `chat_variables` 对象；`session_id` 不替代 `dialogue_id`。

`messageList/lastPushedUserMessageUid` 驱动消息及发送定位。`myChatList` 配套 `myChatListLoading/myChatListLoadCompleted/myChatListError`；历史配套 `chatMessageLoading/chatMessageLoadCompleted/chatMessageError`。`chatError/sendError/variablesError` 是提示 key，敏感词另存 `sensitiveWords`；请求失败与空列表、无更多分别处理。

| 动作 | 参数与结果 |
| --- | --- |
| `initialize(data)` | welcome 与列表独立启动；返回 welcome 响应或 `null`。 |
| `createChat(data)` | 新建传 `dialogue_id: 0`；历史传对应 ID。可传 `robot_key/openid/avatar/name/nickname/rel_user_id`；返回 welcome 响应或 `null`，保留同机器人列表。 |
| `openChat(data)` | welcome 后加载首批历史；返回历史响应或 `null`，页面不再重复请求首批。 |
| `getMyChatList({ refresh = false } = {})` | 列表每页25条；成功返回响应，失败/过期/不满足条件返回 `null`。刷新成功后更新列表，失败保留已有数据。 |
| `onGetChatMessage()` | 历史每页20条，按 ID 升序前插、恢复过程步骤；返回响应或 `null`。 |
| `sendMessage({ message, checkQuestion, global })` | 返回 `Promise<boolean>`；`message` 是协议 question，`checkQuestion` 可为原文本，`global` 默认 JSON 空对象。接纳与草稿消费规则见下节。 |
| `stopMessage()` | 只停止自己的当前检查/连接，不自动重发。 |
| `handleEditVariables({ chat_prompt_variables: values })` | 新会话缓存成功返回 `true`；历史提交成功返回接口响应；失败/失效返回 `null`，保留原变量。 |
| `$reset()` | 停止自己的流并清理会话、列表、分页、变量及发送状态，使在途结果失效。 |

### API、发送与恢复

请求均沿现有封装：`src/api/chat/index.js` 的 welcome、会话列表、历史消息、SSE发送和变量提交；`src/api/robot/index.js` 的发前权限/词检查。列表和历史分别使用 `min_id` 分页，请求失败不推进游标。

1. 页面先判断空内容、附件失败/上传中及图片开关；上传沿既有 `useUpload`，最多10张、单张10MB。
2. Store 检查变量和 `checkChatRequestPermission` 后建立 SSE。文本 question 沿原字符串；多模态为 JSON 数组字符串，包含 `{type:'text',uid,text}` 和 `{type:'image_url',uid,image_url:{url}}`，不新增 `msg_type` 请求字段。
3. `sendMessage()` 只有收到当前有效 `c_message/robot` 业务接纳才返回 `true`，不把构造连接或打开传输当成功。此前词检查失败、传输错误、停止、切会话/机器人或关闭均返回 `false`，调用方保留输入与附件。接纳只结算一次，`c_message` 不重复插入用户消息。
4. 页面收到 `true` 后仍核对会话代次、发送序号和提交时草稿版本，只有草稿未变才清理。消息菜单发送不覆盖或消费正在编辑的草稿；接纳后的错误保留已经展示的用户内容和回复输出。
5. 新会话首次接纳后，捕获实际 `dialogue_id`，按现有 `size: 1` 请求补入真实列表项。补入属于机器人上下文，切同机器人会话不会丢弃；不能覆盖更晚刷新或冒充其他记录。补入失败或最新项不匹配会通过toast提示 `error_chat_list`。

事件保留 `dialogue_id/c_message/robot`、`reasoning_content/llm_rounds/stream_message/sending`、`tool_call_full/tool_result/FileOperation/ExecuteCommand`、`ai_message/finish`、`reply_content_list/start_quote_file/quote_file`、`debug/error/recall_time/request_time/chat_prompt_variables`，过程与历史使用 `agent-process` 工具维护 `process_steps`。只展示实际服务端事件，不伪造工具执行。

业务 `error` 事件沿 `clawbot-chat.js` 记录消息错误并结束当前过程状态，不设置发送失败提示、不终止SSE连接，后续 `sending/ai_message` 等事件继续接收；传输错误由连接的 `onError` 收尾。

共用 `robot-message-item.vue` 在当前回复底部展示非空 `item.error`：浅红背景与边框、错误图标及“异常信息”标题，正文按纯文本保留换行、自动换行并可选中复制。无正文时也展示，后续回复不会覆盖错误提示；空错误不显示。客服Store只将业务 `error` 写入该字段，连接异常通过 `sendError` 交页面toast提示。

页面及人设抽屉的操作错误、输入检查、列表/历史/辅助选项加载错误统一使用toast，不显示常驻错误条、不提供错误重试入口。会话建立失败仅保留“会话暂不可用”状态占位，具体错误通过toast提示；会话变量保存失败保留弹窗中的编辑内容。

**`finish` 只结束过程步骤，不能直接 abort 整个连接。** 其后合法的最终内容、引用、耗时和变量仍需接收；`onClose/onError/stop` 才幂等收尾本次连接和锁。旧回调不能解锁新流。分页失败不推进游标；页面前插历史前后校验会话代次，并按可见消息 uid 恢复滚动锚点。

## 正式模型、人设与变量

模型和人设修改会保存**正式机器人配置**，不是本地临时聊天参数。基础 Store 的三个窄动作如下：

| 动作 | 约束 |
| --- | --- |
| `mergeCurrentRobotInfo(expectedIdentity, partialRobot)` | 返回是否接纳 welcome 配置；核对原对象、身份和类型，缺类型可沿用当前值，显式类型必须与当前应用一致。 |
| `refreshCurrentRobot(expectedIdentity)` | 返回详情或 `null`；校验ID、key和应用类型不变，保持现有 robot/ready，不走清空详情的初始化路径。 |
| `saveRobotSettings(patch)` | 只接 `{ prompt }` 或 `{ model_config_id, use_model }`，单保存；失败抛错，成功返回 `{ saved: true, refreshed: boolean }`。 |

`expectedIdentity` 必须是**本次操作前基础 Store 的原 `robot` 对象引用**，每次 welcome/刷新重新捕获，不能克隆、重造 `{id, robot_key}` 或长期复用初始快照。配置写入统一替换对象，结合私有序号拒绝旧配置回填；快照过期仅拒绝配置，不丢弃同身份/会话代次仍有效的 customer、会话或变量结果。

保存沿既有 `saveRobot`，从完整配置构建请求，处理欢迎语、未知问题、结构化人设、缓存及头像序列化，显式传当前应用真实的 `application_type`（0或2），不以逻辑或将有效的0默认成类型2。页面不乐观改正式值；输入区在保存结束后按正式 props 重建模型选择器。人设取消不保存，保存失败保留编辑。

若结果为 `{ saved: true, refreshed: false }`，表示服务端保存已成功、配置刷新失败或身份已变，**不能再次保存来重试**。当前页面仅在身份仍有效时通过warning toast提示“配置已保存，但页面配置未能同步”，不提供刷新重试入口；不清消息、不重建当前会话。离页或切机器人无法撤销已到服务端的正式保存，旧结果不能回填新身份。

局部表单接口：

- `MessageInput`：props 为 `value/fileList/loading/modelLoading/modeId/modeName/showUpload/emptyMode`；事件为 `update:value/update:fileList/send/stop/changeModel/openPrompt/openKnowledge`。输入法组合期间不发送，Enter发送、Shift+Enter换行。
- `PromptDrawer`：props 为 `open/prompt/robotKey/loading`；事件为 `close/save(prompt)/variablesLoaded(list)`。变量及小程序卡片辅助查询保留，身份/查询序号保护旧结果；不直接保存配置或改业务 Store。
- `VariableModal`：props 为 `open/variables/loading`，`variables` 是完整 `chat_variables` 对象；事件为 `submit(values)/close`。表单编辑本地副本，页面交 Store 提交/缓存；失败保留表单。发送中先停止再修改会话变量。

## 身份切换与知识库边界

- 首次 ready 后初始化；无效或残缺客服身份不请求聊天，由公共布局处理机器人提示。
- 新建/打开历史先停止旧发送，递增会话代次，清草稿、附件、弹层和滚动。局部输入/变量表单按会话代次重建，防旧上传或校验回填；同机器人会话列表保留。
- 切机器人及离页递增上下文代次并 reset，拒绝旧 welcome、检查、列表、历史、变量和流结果，包括 A→B→A。只清客服聊天状态和自己的连接，不取消其他模块请求。
- 同身份模型/人设配置刷新只更新配置相关 UI，不重新 initialize、不清消息。离开再进入则按当前身份重新初始化。
- 库元数据沿 `getLibraryList/getSpecifyAbilityConfig`，默认库优先，微信库遵守现有能力开关；失败使用toast提示。知识条详情新窗打开既有 `/library/details/knowledge-document?id=...`。
- 输入区“知识库”跳转客服 `/customer-service/knowledge-base`，仅携带当前 `id/robot_key`。**目标仍是占位页，本轮没有实现知识库管理。**

## 交付与运行边界

已对照当前源码进行人工范围、契约和可接管性阅读；**未执行验证**。没有运行开发服务、浏览器、业务接口试请求、lint、静态/类型检查、测试或构建，不能将人工阅读视为真实链路通过。

建议用户手动检查类型0/2的切换、刷新恢复、普通消息及已有工具任务的实际过程展示，确认生成中切应用不串消息、配置保存保持原类型，并检查正式保存成功但刷新失败的恢复。这些是后续检查建议，尚未执行。发现需改 API、后端或共享组件契约时先协调范围，不能在本页静默扩展。
