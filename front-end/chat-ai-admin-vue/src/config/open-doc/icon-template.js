import { useI18n } from '@/hooks/web/useI18n'

const iconTemplatesData = [
  {
    id: 1,
    nameKey: 'template_1',
    preview: [
      { level: 0, textKey: 'level_1_doc', icon: '📢', color: '' },
      { level: 0, textKey: 'level_1_folder', icon: '📁', color: '' },
      { level: 1, textKey: 'level_2_doc', icon: '📋', color: '' },
      { level: 1, textKey: 'level_2_folder', icon: '📂', color: '' },
      { level: 2, textKey: 'level_3_doc', icon: '📋', color: '' }
    ],
    levels: {
      0: { doc_icon: '📢', folder_icon: '📁' },
      1: { doc_icon: '📋', folder_icon: '📂' },
      2: { doc_icon: '📋', folder_icon: '📂' },
    }
  },
  {
    id: 2,
    nameKey: 'template_2',
    preview: [
      { level: 0, textKey: 'level_1_doc', icon: '🥇', color: '' },
      { level: 0, textKey: 'level_1_folder', icon: '📁', color: '' },
      { level: 1, textKey: 'level_2_doc', icon: '🥈', color: '' },
      { level: 1, textKey: 'level_2_folder', icon: '📂', color: '' },
      { level: 2, textKey: 'level_3_doc', icon: '🥉', color: '' }
    ],
    levels: {
      0: { doc_icon: '🥇', folder_icon: '📂' },
      1: { doc_icon: '🥈', folder_icon: '📂' },
      2: { doc_icon: '🥉', folder_icon: '📂' },
    }
  },
  {
    id: 3,
    nameKey: 'template_3',
    preview: [
      { level: 0, textKey: 'level_1_doc', icon: '📝', color: '' },
      { level: 0, textKey: 'level_1_folder', icon: '📁', color: '' },
      { level: 1, textKey: 'level_2_doc', icon: '📝', color: '' },
      { level: 1, textKey: 'level_2_folder', icon: '📂', color: '' },
      { level: 2, textKey: 'level_3_doc', icon: '📝', color: '' }
    ],
    levels: {
      0: { doc_icon: '📝', folder_icon: '📁' },
      1: { doc_icon: '📝', folder_icon: '📂' },
      2: { doc_icon: '📝', folder_icon: '📂' },
    }
  }
]

export function getIconTemplateList() {
  const { t } = useI18n('config.open-doc.icon-template')
  
  return iconTemplatesData.map(template => ({
    ...template,
    name: t(template.nameKey),
    preview: template.preview.map(item => ({
      ...item,
      text: t(item.textKey)
    }))
  }))
}

export function getIconTemplateById(id) {
  const { t } = useI18n('config.open-doc.icon-template')
  
  let template = iconTemplatesData.find((item) => item.id == id)
  if (!template) {
    template = iconTemplatesData[0]
  }

  template = JSON.parse(JSON.stringify(template))
  
  // 执行翻译
  template.name = t(template.nameKey)
  template.preview = template.preview.map(item => ({
    ...item,
    text: t(item.textKey)
  }))

  let icons = template.levels[2]

  for (let i = 3; i < 100; i++) {
    // template.levels[i] = template.levels[i % 3];
    template.levels[i] = icons;
  }

  return template;
}