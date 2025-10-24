const docIcons = [
  // 标题/结构类
  {
    content: '🥇',
    keywords: ['第一步', '第一名', '优先', '首选'],
  },
  {
    content: '🥈',
    keywords: ['第二步', '第二名', '次选'],
  },
  {
    content: '🥉',
    keywords: ['第三步', '第三名'],
  },
  {
    content: '🔹',
    keywords: ['小节标题', '列表前缀'],
  },
  {
    content: '🔸',
    keywords: ['小节标题', '重点提示'],
  },
  {
    content: '📌',
    keywords: ['固定', '关键点', '重要'],
  },
  {
    content: '📝',
    keywords: ['写作', '笔记', '记录'],
  },

  // 状态类
  {
    content: '✅',
    keywords: ['成功', '完成', '推荐', '正确'],
  },
  {
    content: '❌',
    keywords: ['失败', '禁止', '不要做', '错误'],
  },
  {
    content: '⚠️',
    keywords: ['警告', '注意事项'],
  },
  {
    content: '⛔',
    keywords: ['严重警告', '禁用'],
  },
  {
    content: '❗',
    keywords: ['强提醒', '高亮提示', '重要'],
  },
  {
    content: '🔔',
    keywords: ['通知', '提醒操作'],
  },
  {
    content: '🔕',
    keywords: ['关闭提醒', '静音状态'],
  },
  {
    content: '💡',
    keywords: ['建议', '灵感', '经验分享', '想法'],
  },
  {
    content: '🔄',
    keywords: ['更新', '刷新', '修改'],
  },
  {
    content: '🔧',
    keywords: ['配置', '调试', '设置'],
  },
  {
    content: '🛠️',
    keywords: ['工具', '维护', '手动操作'],
  },

  // 文件/操作类
  {
    content: '📁',
    keywords: ['文件夹', '路径', '目录'],
  },
  {
    content: '📄',
    keywords: ['文件', '文档', '内容'],
  },
  {
    content: '📂',
    keywords: ['展开文件夹', '打开'],
  },
  {
    content: '📋',
    keywords: ['复制内容', '剪贴板'],
  },
  {
    content: '🗑️',
    keywords: ['删除操作', '清理'],
  },
  {
    content: '💾',
    keywords: ['保存', '存储'],
  },
  {
    content: '🗂️',
    keywords: ['分类', '归档', '整理'],
  },

  // 系统/代码类
  {
    content: '💻',
    keywords: ['编程', '开发', '代码'],
  },
  {
    content: '🖥️',
    keywords: ['桌面', '系统操作'],
  },
  {
    content: '📟',
    keywords: ['命令行', '终端'],
  },
  {
    content: '📦',
    keywords: ['安装包', '依赖', '组件'],
  },
  {
    content: '🧰',
    keywords: ['工具集', '开发者工具箱'],
  },
  {
    content: '🧪',
    keywords: ['测试', '实验'],
  },
  {
    content: '🚀',
    keywords: ['发布', '启动', '上线'],
  },
  {
    content: '🧱',
    keywords: ['构建', '模块', '组件'],
  },

  // 强调/提示类
  {
    content: '🚨',
    keywords: ['紧急提示', '警报'],
  },
  {
    content: '🔥',
    keywords: ['高亮', '热度高', '重点', '热门'],
  },
  {
    content: '📣',
    keywords: ['公告', '提醒', '通知'],
  },
  {
    content: '📢',
    keywords: ['强调', '提示', '宣告'],
  },
  {
    content: '👀',
    keywords: ['注意', '查看', '易错点'],
  },

  // 导向/流程类
  {
    content: '🔜',
    keywords: ['接下来', '下一步'],
  },
  {
    content: '🔙',
    keywords: ['返回', '撤销'],
  },
  {
    content: '🔛',
    keywords: ['激活中', '当前'],
  },
  {
    content: '⏩',
    keywords: ['快进', '下一步'],
  },
  {
    content: '⏪',
    keywords: ['上一步', '回退'],
  },
  {
    content: '🔁',
    keywords: ['循环', '重复'],
  },
  {
    content: '🧭',
    keywords: ['指引', '导航', '引导文档'],
  },

  // 设置/配置类
  {
    content: '⚙️',
    keywords: ['高级设置', '参数调整'],
  },
  {
    content: '🎛️',
    keywords: ['面板设置', '可视化配置'],
  },

  // 关闭/禁用类
  {
    content: '🛑',
    keywords: ['强制停止', '立即终止'],
  },
  {
    content: '🚫',
    keywords: ['拒绝', '阻止操作'],
  },
  {
    content: '🧯',
    keywords: ['熄火', '冷处理', '取消计划'],
  },

  // 文档/写作类
  {
    content: '📚',
    keywords: ['文档集合', '教程'],
  },
  {
    content: '📰',
    keywords: ['公告', '资讯', '更新说明'],
  },
  {
    content: '✏️',
    keywords: ['编辑', '修改', '草稿'],
  },

  // 分类/标签类
  {
    content: '🏷️',
    keywords: ['标签分类'],
  },
  {
    content: '🧾',
    keywords: ['列表', '数据明细'],
  },
  {
    content: '🗃️',
    keywords: ['归档内容'],
  },

  // 时间/计划类
  {
    content: '⏱️',
    keywords: ['计时', '效率工具'],
  },
  {
    content: '⏳',
    keywords: ['等待中', '延迟'],
  },
  {
    content: '📅',
    keywords: ['日期', '计划安排'],
  },
  {
    content: '🕒',
    keywords: ['时间说明'],
  },
  {
    content: '📆',
    keywords: ['日历视图', '按日分类'],
  },

  // 成就/进度类
  {
    content: '🎯',
    keywords: ['目标', '任务完成'],
  },
  {
    content: '🎖️',
    keywords: ['奖励', '成就'],
  },
  {
    content: '📈',
    keywords: ['进度上升', '优化效果'],
  },
  {
    content: '📉',
    keywords: ['问题回退', '异常波动'],
  },

  // 数字类
  {
    content: '1️⃣',
    keywords: ['数字一', '第一步'],
  },
  {
    content: '2️⃣',
    keywords: ['数字二', '第二步'],
  },
  {
    content: '3️⃣',
    keywords: ['数字三', '第三步'],
  },
  {
    content: '4️⃣',
    keywords: ['数字四'],
  },
  {
    content: '5️⃣',
    keywords: ['数字五'],
  },
  {
    content: '6️⃣',
    keywords: ['数字六'],
  },
  {
    content: '7️⃣',
    keywords: ['数字七'],
  },
  {
    content: '8️⃣',
    keywords: ['数字八'],
  },
  {
    content: '9️⃣',
    keywords: ['数字九'],
  },
  {
    content: '🔟',
    keywords: ['数字十'],
  },
];

export default docIcons;