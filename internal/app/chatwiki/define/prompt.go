// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

const (
	PromptRoleTypeSystem = 0
	PromptRoleTypeUser   = 1
	PromptRoleUser       = `user`
)

var PromptRoleTypeMap = map[int]string{
	PromptRoleTypeSystem: `system`,
	PromptRoleTypeUser:   `user`,
}

const PromptDefaultQuestionOptimize = `
# 请根据以下步骤分析用户当前问题及上下文和对话背景，补全用户的提问
1、上下文分析
- 提取用户问题中的核心关键词
- 识别当前问题的潜在信息缺口（如：问题主语/时间范围/专业领域/具体场景）
- 推测用户潜在需求层级（基础信息/对比分析/解决方案/原理推导）
2、优化策略选择
- 【多维度扩展】添加对比维度（如：时间/空间/类型对比）
- 【精准化升级】补充限定条件（行业/地域/技术参数）
- 【结构化重组】将开放式问题转化为封闭式+开放式复合结构
- 【语义聚焦】消除歧义词汇，强化领域术语
3、生成优化方案
- 输出你优化后的问题，只用输出1个问题
- 如果你认为用户问题不需要补全，请返回用户问题原文

# 优化原则：
- 保持自然对话衔接，避免机械式提问
- 优先保留用户原始意图关键词

{{dialogue_background}}
`

const PromptDefaultGraphConstruct = `你是一个知识图谱构建专家。请分析以下文本，提取实体、关系和属性，构建知识图谱。
输出格式为JSON数组，每个元素包含以下字段：
- subject: 主体实体
- predicate: 关系
- object: 客体实体或属性值
- confidence: 置信度(0-1)

文本内容:
"""
{{content}}
"""

请确保输出的JSON格式正确，只返回JSON数组，注意，不要有其他文字说明，关系和实体中都不要有括号、引号、逗号等任何特殊字符或标点符号。
`

const PromptDefaultQuestionGuide = `
你是一个AI智能助手，可以回答和解决我的问题。请结合历史对话记录，帮我生成 {{num}} 个问题，引导我继续提问。
历史记录:
"""
{{histories}}
"""
每个问题的长度应小于20个字符，按如下格式返回: ["问题1", "问题2", "问题3", "问题4", "问题5", "问题6", "问题7", "问题8", "问题9", "问题10"]`

const PromptDefaultEntityExtract = `请从下面的问题中提取关键实体，以便进行知识图谱检索。
尽量从不同角度提取关键概念和实体。注意一定要只输出JSON数组格式，每个元素是一个实体字符串，不要有其他文字说明。

问题:
"""
{{question}}
"""

输出示例：["实体1", "实体2", "实体3", "实体4", "实体5"]`

const PromptDefaultReplyMarkdown = `- 请使用markdown格式回答问题。`
const PromptDefaultAnswerImage = `- 当你选择的知识点中包含图片、链接数据时，你需要在你的答案对应位置输出这些数据，不要改写或忽略这些数据。`
const PromptDefaultCreatePrompt = `# 角色
你是一个智能提示词工程师，擅长将用户需求转化为结构化的大模型指令模板。
# 任务
根据用户输入的领域需求，自动生成包含必填模块和可扩展自定义模块的提示词框架。
# 输出
## 基础结构说明
{
    "role": {"subject":"", "describe":"..."},
    "task": {"subject":"", "describe":"..."},
    "constraints": {"subject":"", "describe":"..."},
    "output": {"subject":"", "describe":"..."},
    "tone": {"subject":"", "describe":"..."},
    "custom": [{"subject":"", "describe":"..."}]
}
## 字段说明
- role：角色定义，需要包含模拟角色的专业领域和核心能力，
- task：核心任务，描述需要解决的问题和预期结果
- constraints：约束条件，用换行符分隔的条款式描述
- output：输出规范
- tone：表达风格
- custom：自定义拓展字段，根据用户要求智能添加
## 输出示例
{
    "role": {
        "subject": "",
        "describe": "你是智能家居领域的AI顾问，精通设备联动方案设计"
    },
    "task": {
        "subject": "",
        "describe": "根据用户家庭户型图，推荐最优设备布局方案"
    },
    "constraints": {
        "subject": "",
        "describe": "- 方案需标注设备间最大距离\n- 优先考虑节能配置\n- 避开承重墙位置"
    },
    "output": {
        "subject": "",
        "describe": "- 输出带标注的平面图PDF\n- 附设备清单表格"
    },
    "tone": {
        "subject": "",
        "describe": "使用家居杂志编辑风格，搭配场景化描述"
    },
    "custom": [
        {
            "subject": "扩展兼容",
            "describe": "支持导入主流智能家居平台现有配置"
        }
    ]
}
# 要求
- 根据用户提供的领域关键词，自动补充行业特定的约束条件和输出要求
- 必须严格保持JSON格式，描述使用的语言与用户输入的语言保持一致。比如，用户输入的是中文，请使用中文输出。
- 必填模块必须包含role/task/constraints/output/tone五个字段；
- custom模块需智能判断用户需求添加，每个条目应有明确subject和具体描述；
- 所有describe字段需用自然口语化表达，避免技术术语。`

const PrumptAiChunk = `你是一位文章分段助手，根据文章内容的语义进行合理分段，确保每个分段表述一个完整的语义，每个分段字数控制在500字左右，最大不超过1000字。请严格按照文章内容进行分段，不要对文章内容进行加工，分段完成后输出分段后的内容。`

const PromptGenerateSimilarQuestions = `有如下问题对：
------
question: {{question}}
answer: {{answer}}
------
你需要生成{{num}}条相似问题。
输出示例：["问题1", "问题2", "问题3", "问题4", "问题5"]]
`

const PromptLibAiSummary = `将提交的内容进行智能总结,不要随意发挥`

const PromptAiGenerate = `将提交的内容生成用于大模型对话的提示词,不超过%v字`

const PromptWorkFlowQuestionOptimize = `# 请根据以下步骤分析用户当前问题及上下文和对话背景，补全用户的提问
1、上下文分析
- 提取用户问题中的核心关键词
- 识别当前问题的潜在信息缺口（如：问题主语/时间范围/专业领域/具体场景）
- 推测用户潜在需求层级（基础信息/对比分析/解决方案/原理推导）
2、优化策略选择
- 【多维度扩展】添加对比维度（如：时间/空间/类型对比）
- 【精准化升级】补充限定条件（行业/地域/技术参数）
- 【结构化重组】将开放式问题转化为封闭式+开放式复合结构
- 【语义聚焦】消除歧义词汇，强化领域术语
3、生成优化方案
- 输出你优化后的问题，只用输出1个问题
- **如果你认为用户问题语意完整，不需要补全，请返回用户问题原文**

# 优化原则：
- 保持自然对话衔接，避免机械式提问
- 优先保留用户原始意图关键词
- **注意，你只是优化和补全用户的问题，而不是要回答用户的问题！**`

const PromptFaqFileAiChunk = `根据user角色提供的文本，学习和分析它，并整理学习成果：
- 提出问题并给出每个问题的答案。
- 答案需详细完整，尽可能保留原文描述，可以适当扩展答案描述。
- 答案可以包含普通文字、链接、代码、表格、公示、媒体链接等 Markdown 元素。
- 最多提出 50 个问题。
- 生成的问题和答案和源文本语言相同。
`

const CompletionGenerateJsonPrompt = ` ### 角色
你是一个结构化输出助手，负责根据提供的特定标准提取结构化信息。请遵循以下指南以确保一致性和准确性。说明：
以下提供了一些额外信息，如果信息有用需尽量参考。
<instructions>
%s
</instructions>

### 提取参数
需要从输入文本中提取以下信息。<structure> 标签指定了要提取的信息的'typ'、'key'、'enum'和'vals'。然后将提取的内容写到vals中;
<structure>
%s
</structure>

### 流程
步骤1：仔细阅读输入并理解预期输出的结构。
步骤2：根据对象的名称和描述从提供的文本中提取相关参数。
步骤3：按照 <structure> 中指定的方式将提取的参数构造为 JSON 对象。
步骤4：确保 JSON 对象格式正确且有效。输出不应包含任何标签。只能输出可解析的 JSON 对象。

### 记忆
从user角色提供的文本中提取相关信息。
### 示例
以下是预期输出的结构，你要始终遵循下面结构输出。
[
    {
        "key": "properties1",
        "typ": "string",
        "enum": "提供的枚举值",
        "vals": [
            "提取的相关文本"
        ]
    },
    {
        "key": "properties2",
        "typ": "array<string>",
        "enum": "",
        "vals": [
            "提取的相关文本"
        ]
    }
]
### 回答
你始终输出一个有效的可用于解析的 JSON 对象,去掉markdown标签,除了 JSON 结构外不要输出其他内容。
`

const OfficialAccountCommentCheckPrompt = `### 角色
你是一个顶尖的自媒体行业内容分析师，负责面向用户的舆情处置工作，能够从用户的角度来理解他们发表的想法。
# 任务
根据提供给你的评论内容，判断这个评论是否需要删除、回复、置顶
# 输出
## 基础结构说明
{
    "need_delete": "这个字段表示是否需要删除",
    "need_reply": "这个字段表示是否需要回复",
    "need_top": "这个字段表示是否需要置顶",
    "reply_content": "这里是需要回复的内容"
}
# 要求
- 根据我提供的评论内容和需求给出对应的结果，如果是判断类的，以true或者false返回，不需要返回全部的字段，针对需求返回即可
- 必须严格保持JSON格式。
`
