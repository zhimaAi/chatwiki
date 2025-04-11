// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

const PromptDefaultQuestionOptimize = `
作为一个向量检索助手，你的任务是结合历史记录，从不同角度，为"原问题"生成个不同版本的"检索词"，从而提高向量检索的语义丰富度，提高向量检索的精度。生成的问题要求指向对象清晰明确，并与"原问题语言相同"。例如：
历史记录:
"""
"""
原问题: 介绍下剧情。
检索词: ["介绍下故事的背景。","故事的主题是什么？","介绍下故事的主要人物。"]
----------------
历史记录:
"""
Q: 对话背景。
A: 当前对话是关于 Nginx 的介绍和使用等。
"""
原问题: 怎么下载
检索词: ["Nginx 如何下载？","下载 Nginx 需要什么条件？","有哪些渠道可以下载 Nginx？"]
----------------
历史记录:
"""
Q: 对话背景。
A: 当前对话是关于 Nginx 的介绍和使用等。
Q: 报错 "no connection"
A: 报错"no connection"可能是因为……
"""
原问题: 怎么解决
检索词: ["Nginx报错"no connection"如何解决？","造成'no connection'报错的原因。","Nginx提示'no connection'，要怎么办？"]
----------------
历史记录:
"""
Q: 护产假多少天?
A: 护产假的天数根据员工所在的城市而定。请提供您所在的城市，以便我回答您的问题。
"""
原问题: 沈阳
检索词: ["沈阳的护产假多少天？","沈阳的护产假政策。","沈阳的护产假标准。"]
----------------
历史记录:
"""
{{histories}}
"""
原问题: {{query}}
检索词:
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

请确保输出的JSON格式正确，只返回JSON数组，注意，不要有其他文字说明，关系中不要有任何特殊字符。
`

const PromptDefaultQuestionGuide = `
你是一个AI智能助手，可以回答和解决我的问题。请结合历史对话记录，帮我生成 3 个问题，引导我继续提问。
历史记录:
"""
{{histories}}
"""
问题的长度应小于20个字符，按如下格式返回: ["问题1", "问题2", "问题3"]`

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
