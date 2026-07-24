# PC Chat 消息流说明

本文档说明 PC/PC 插件如何将旧版和 Agent/Clawbot SSE 事件转换为统一的过程状态。所有机器人文本消息均以 `process_steps` 为过程数据源，并通过 `ProcessTimeline` 展示。

## 架构与职责

```text
旧版/新版 SSE
      ↓ normalizeSseEvent
NormalizedChatEvent
      ↓ applyNormalizedChatEvent
process_steps + content
      ↓
ProcessTimeline + 最终正文
```

- `src/stores/modules/chat.ts`：解析 SSE、转换内部事件、更新消息和历史兼容数据。
- `src/views/chat/index.vue`：页面编排、发送、停止和滚动，不解析 SSE。
- `components/messages/message-item.vue`：组合引用、过程时间线、正文、菜单和操作区。
- `components/messages/process-timeline.vue`：展示思考、工具、技能步骤，展开状态只保存在组件本地。

后端按照 `application_type` 选择旧链路或新链路，正常情况下同一回复不会同时发送 `sending` 和 `stream_message.content`。前端不增加协议模式字段，也不根据 `application_type` 切换渲染。

## 原始事件与内部事件

| 原始事件 | 内部事件 | 语义 |
| --- | --- | --- |
| `reasoning_content` | `thinking_delta` | 旧链路思考增量 |
| `stream_message.reasoning_content` | `thinking_delta` | 新链路思考增量 |
| `sending` | `answer_delta` | 旧链路正文增量，保持原自动滚动行为 |
| `stream_message.content` | `answer_delta` | 新链路正文增量，保持过程查看位置 |
| `llm_rounds:begin` | `round_begin` | 开始下一轮模型推理 |
| `llm_rounds:finish` | `round_finish` | 兜底完成活动思考 |
| `tool_call_full` | `tool_start` | 创建工具或技能步骤 |
| `tool_result` | `tool_finish` | 按 `tool_call_id` 完成步骤 |
| `FileOperation`、`ExecuteCommand` | `operation` | 保存隐藏操作记录 |
| `finish` | `process_finalize/finish` | 正常结束过程 |
| 用户停止 | `process_finalize/stop` | 标记停止并完成运行步骤 |
| SSE `onClose` | `process_finalize/close` | 断流清理和释放发送锁 |
| `ai_message` | `final_snapshot` | 同步服务端最终消息快照 |

`stream_message` 会被转换成事件数组。如果思考和正文异常地同时非空，必须先消费 `thinking_delta`，再消费 `answer_delta`，确保当前思考正确完成。

## 完整 SSE 事件表

| 事件 | 处理方式 |
| --- | --- |
| `dialogue_id` | 更新当前对话 ID |
| `c_message` | 插入用户消息 |
| `robot` | 插入预创建的机器人占位消息 |
| `reply_content_list` | 更新多内容回复列表 |
| `reasoning_content` | 转换为思考增量 |
| `sending` | 转换为正文增量 |
| `llm_rounds` | 转换为轮次开始或结束 |
| `stream_message` | 按字段转换为思考、正文增量 |
| `tool_call_full` | 创建普通工具或技能步骤 |
| `tool_result` | 完成已关联的工具或技能步骤 |
| `FileOperation` | 保存隐藏操作步骤 |
| `ExecuteCommand` | 保存隐藏操作步骤 |
| `start_quote_file` | 开始知识库检索 |
| `quote_file` | 保存引用文件并结束检索状态 |
| `ai_message` | 同步最终正文、类型、菜单、引用和回复列表 |
| `debug` | 保存调试数据 |
| `chat_prompt_variables` | 更新会话变量 |
| `finish` | 正常完成过程 |
| SSE `onClose` | 完成运行步骤、关闭 loading、释放锁 |

心跳、保活、空帧或未消费事件不会改变过程状态。

## 消息和步骤状态

| 字段 | 说明 |
| --- | --- |
| `content` | 正文增量；`ai_message` 到达后由最终快照覆盖 |
| `reasoning_content` | 接口和历史兼容字段，不直接驱动 UI |
| `process_steps` | 唯一的过程渲染数据源 |
| `current_round_index` | 当前推理轮次 |
| `active_thinking_step_id` | 当前运行中的思考步骤 ID |
| `startLoading` | 是否仍在等待首段正文 |
| `loading` | 当前回复是否仍在传输 |
| `is_stopped` | 是否由用户主动停止 |

步骤类型：

- `thinking`：模型思考，正文出现后从 `running` 变为 `done`。
- `tool`：普通工具，保存原始输入和输出。
- `skill`：技能调用，标题取 `function.arguments.skill`。
- `operation`：文件或命令操作，只记录、不展示、不在前端执行。

## 状态转换规则

1. 非空思考增量创建或追加当前 `thinking/running` 步骤。
2. 旧链路没有 `llm_rounds:begin` 时，首个思考至少归入第 1 轮。
3. 第一段非空正文立即完成活动思考步骤。
4. 空正文、空思考和元数据帧不结束思考，也不创建空步骤。
5. `llm_rounds:finish` 是正文未到达时的完成兜底。
6. 没有思考、直接收到 `sending` 时只追加正文，不生成虚假的完成步骤。
7. `finish`、停止和断流都会完成剩余的 `running` 步骤并关闭 loading。
8. `ai_message` 是最终权威快照，同时保留菜单、引用、推荐回复和语音解析逻辑。

## 工具和技能关联

`tool_call_full.id` 保存为 `tool_call_id`。普通工具使用 `function.name` 作为标题并保存原始 `function.arguments`；当函数名为 `skill` 时创建 `skill` 步骤，标题读取参数中的 `skill`。

`tool_result.tool_call_id` 优先精确匹配运行中的工具或技能。旧数据没有 ID 时，才按最近运行步骤和工具名兼容匹配。仍无法关联时安全忽略，避免错误完成并行调用。

`grep`、`read_file`、`bash`、`glob`、`execute`、`ls` 属于内部普通工具，工具名去除首尾空格并转为小写后精确匹配。命中后步骤保留在 `process_steps` 并标记为 `hidden`，继续参与结果关联和过程收尾，但不进入时间线；同名技能不隐藏。`FileOperation`、`ExecuteCommand` 仍按隐藏的 `operation` 步骤处理。

## 历史消息兼容

历史 `process_steps` 同时兼容数组和 JSON 字符串。已有步骤时先补齐内部工具的隐藏标记；没有步骤且机器人文本消息包含非空 `reasoning_content` 时，生成一个确定性 ID 的 `thinking/done` 步骤：

```text
legacy-thinking-{message_id | id | uid}
```

转换范围仅限 `is_customer != 1 && msg_type == 1`。用户消息、菜单、图片、文件和空思考不转换。确定性 ID 和“已有步骤优先”规则共同避免分页或重新进入时重复生成。

## UI 和样式边界

- PC/PC 插件的机器人文本消息统一使用过程时间线，包括只发送 `sending` 的旧链路。
- 等待时显示时间线标题；有运行步骤后 loading 移到步骤前。
- 配置的答案提示语开启且非空时作为运行标题，否则显示“思考中”。
- 完成后显示“思考完成”，用户停止显示“已停止”。
- 思考默认收起并显示单行摘要；用户展开后由组件本地状态控制。
- 工具和技能详情最大高度 `200px`，超出后局部滚动。
- 知识库 loading 可见时不重复显示时间线标题 loading。
- 时间线样式仅存在于消息组件，不修改页面最大宽度、桌面布局、头像和输入区。
- 菜单、图片、文件、引用、推荐问题、语音和用户消息继续使用原渲染分支。

## 配置兼容

答案生成提示语优先读取当前语言配置，缺失时回退机器人顶层配置；语言配置显式关闭时不回退顶层开关。多语言配置兼容数组和 JSON 字符串，开关兼容布尔值 `true` 和字符串 `"true"`，非法配置安全回退。
