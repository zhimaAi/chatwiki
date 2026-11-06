import currentAntTheme from '../theme'

const palette = {
  primary: '#2475fc',
  page: '#ffffff',
  workspace: '#f5f6f8',
  statisticsCard: '#f2f4f7',
  text: '#262626',
  textSecondary: '#595959',
  textMuted: '#8c8c8c',
  headerBorder: '#f2f4f7',
  deleteText: '#fb363f'
}

const defaultTheme = {
  id: 'default',
  name: '默认主题',
  sources: [
    'src/config/theme.js',
    'src/assets/base.css',
    'src/assets/main.css',
    'src/layouts/AdminLayout/index.vue',
    'src/views/customer-service/index.vue',
    'src/views/customer-service/components/customer-sidebar.vue',
    'src/views/customer-service/components/customer-user-menu.vue',
    'src/views/robot/robot-config/statistical_analysis/hit-statics.vue'
  ],
  palette,
  // 默认主题保留历史业务配色；独立于 Ant Token，避免改变组件库原有外观。
  legacyBusiness: {
    layoutBackground: palette.workspace,
    containerBackground: palette.page,
    textPrimary: palette.text,
    textSecondary: palette.textSecondary,
    textTertiary: palette.textMuted,
    link: palette.primary,
    linkHover: palette.primary,
    linkActive: palette.primary
  },
  typography: {
    ui: currentAntTheme.token.fontFamily,
    body: { fontSize: 15, lineHeight: 1.6 },
    sidebar: { fontSize: 14, fontWeightSelected: 600 }
  },
  radius: {
    statisticsCard: 6,
    sidebarChildItem: 8,
    sidebarItem: 10,
    popover: 12
  },
  shadow: {
    popover: '0 8px 24px rgba(0, 0, 0, 0.12)'
  },
  // 系列颜色来自统计页，不能与主色一起批量替换。
  chart: {
    totalMessages: '#0079fe',
    hitMessages: '#00c292',
    missMessages: '#ff6b6b',
    hitRate: '#333333',
    pointBorder: '#ffffff'
  },
  // 保留当前显式配置；业务提取值未自动写入组件库，避免改变现有默认外观。
  antTheme: currentAntTheme
}

export default defaultTheme
