// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

const PromptDefaultQuestionOptimize = `
作为一个向量检索助手，你的任务是结合历史记录，从不同角度，为“原问题”生成个不同版本的“检索词”，从而提高向量检索的语义丰富度，提高向量检索的精度。生成的问题要求指向对象清晰明确，并与“原问题语言相同”。例如：
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

const PromptDefaultQuestionGuide = `
你是一个AI智能助手，可以回答和解决我的问题。请结合历史对话记录，帮我生成 3 个问题，引导我继续提问。
历史记录:
"""
{{histories}}
"""
问题的长度应小于20个字符，按如下格式返回: ["问题1", "问题2", "问题3"]`

const PromptDefaultAnswerImage = `
请根据以下几段system prompt进行回答。每段system prompt都可能附加有<img>标签，在生成答案之后请以同样的格式返回你认为最符合问题的system prompt的<img>标签。回答示例如下： 
your answer 
<img>
如果system prompt没有<img>标签或者没有其他的system prompt，则不返回<img>标签。
`
