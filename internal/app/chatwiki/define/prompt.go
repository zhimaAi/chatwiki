// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

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
