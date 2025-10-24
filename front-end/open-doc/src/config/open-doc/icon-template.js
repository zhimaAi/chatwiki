const iconTemplates = [
  {
    id: 1,
    name: '模板1',
    preview: [
      { level: 0, text: '一级文档', icon: '📢', color: '' },
      { level: 0, text: '一级文件夹', icon: '📁', color: '' },
      { level: 1, text: '二级文档', icon: '📋', color: '' },
      { level: 1, text: '二级文件夹', icon: '📂', color: '' },
      { level: 2, text: '三级文档', icon: '📋', color: '' }
    ],
    levels: {
      0: { doc_icon: '📢', folder_icon: '📁' },
      1: { doc_icon: '📋', folder_icon: '📂' },
      2: { doc_icon: '📋', folder_icon: '📂' },
    }
  },
  {
    id: 2,
    name: '模板2',
    preview: [
      { level: 0, text: '一级文档', icon: '🥇', color: '' },
      { level: 0, text: '一级文件夹', icon: '📁', color: '' },
      { level: 1, text: '二级文档', icon: '🥈', color: '' },
      { level: 1, text: '二级文件夹', icon: '📂', color: '' },
      { level: 2, text: '三级文档', icon: '🥉', color: '' }
    ],
    levels: {
      0: { doc_icon: '🥇', folder_icon: '📁' },
      1: { doc_icon: '🥈', folder_icon: '📂' },
      2: { doc_icon: '🥉', folder_icon: '📂' },
    }
  },
  {
    id: 3,
    name: '模板3',
    preview: [
      { level: 0, text: '一级文档', icon: '📝', color: '' },
      { level: 0, text: '一级文件夹', icon: '📁', color: '' },
      { level: 1, text: '二级文档', icon: '📝', color: '' },
      { level: 1, text: '二级文件夹', icon: '📂', color: '' },
      { level: 2, text: '三级文档', icon: '📝', color: '' }
    ],
    levels: {
      0: { doc_icon: '📝', folder_icon: '📁' },
      1: { doc_icon: '📝', folder_icon: '📂' },
      2: { doc_icon: '📝', folder_icon: '📂' },
    }
  }
]

export function getIconTemplateList() {
  return JSON.parse(JSON.stringify(iconTemplates))
}

export function getIconTemplateById(id) {
  let template = iconTemplates.find((item) => item.id == id)
  if (!template) {
    template = iconTemplates[0]
  }

  template = JSON.parse(JSON.stringify(template))

  let icons = template.levels[2]

  for (let i = 3; i < 100; i++) {
    // template.levels[i] = template.levels[i % 3];
    template.levels[i] = icons;
  }

  return template;
}
