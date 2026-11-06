import { theme } from 'ant-design-vue'

const palette = {
  nearBlack: '#141413',
  terracotta: '#c96442',
  coral: '#d97757',
  errorCrimson: '#b53333',
  focusBlue: '#3898ec',
  parchment: '#f5f4ed',
  ivory: '#faf9f5',
  white: '#ffffff',
  warmSand: '#e8e6dc',
  darkSurface: '#30302e',
  charcoal: '#4d4c48',
  oliveGray: '#5e5d59',
  stoneGray: '#87867f',
  darkWarm: '#3d3d3a',
  warmSilver: '#b0aea5',
  borderCream: '#f0eee6',
  ringWarm: '#d1cfc5',
  ringDeep: '#c2c0b6'
}

// 浮层采用原文的 Ivory 容器色；这是应用映射，不是原文单独定义的浮层规则。
const semantic = {
  colorPrimary: palette.terracotta,
  colorBgLayout: palette.parchment,
  colorBgContainer: palette.ivory,
  colorBgElevated: palette.ivory,
  colorText: palette.nearBlack,
  colorTextSecondary: palette.oliveGray,
  colorTextTertiary: palette.stoneGray,
  colorBorder: palette.borderCream,
  colorBorderSecondary: palette.warmSand,
  colorError: palette.errorCrimson
}

// 原文没有完整链接状态，使用官方主色算法补齐；信息提示色保留组件库默认语义。
const primaryTokens = theme.defaultAlgorithm({ ...theme.defaultSeed, colorPrimary: semantic.colorPrimary })
const linkTokens = {
  colorLink: semantic.colorPrimary,
  colorLinkHover: primaryTokens.colorPrimaryHover,
  colorLinkActive: primaryTokens.colorPrimaryActive
}

const uiFontFamily = '"Anthropic Sans", Arial, "PingFang SC", "Noto Sans SC", sans-serif'

const claudeTheme = {
  id: 'claude',
  name: '扶光金',
  sources: ['https://fchangjun.github.io/awesome-design-md-cn/design-md/claude/DESIGN.md'],
  palette,
  semantic,
  typography: {
    ui: uiFontFamily,
    reference: {
      headingFontFamily: 'Anthropic Serif',
      headingFallback: 'Georgia',
      uiFontFamily: 'Anthropic Sans',
      uiFallback: 'Arial',
      codeFontFamily: 'Anthropic Mono',
      heading: { fontSize: 32, fontWeight: 500, lineHeight: 1.1 },
      body: { fontSize: 16, fontWeight: [400, 500], lineHeight: [1.25, 1.6] },
      caption: { fontSize: 14, fontWeight: 400, lineHeight: 1.43 },
      code: { fontSize: 15, fontWeight: 400, lineHeight: 1.6, letterSpacing: -0.32 }
    }
  },
  radius: {
    scale: [4, 6, 7.5, 8, 8.5, 12, 16, 24, 32],
    card: 8,
    primaryButtonRange: [8, 12],
    input: 12,
    featured: 16,
    // 通用控件取原文标准按钮的 8px；输入框 12px 留给组件级适配。
    control: 8
  },
  shadow: {
    ring: `0 0 0 1px ${palette.ringWarm}`,
    ringActive: `0 0 0 1px ${palette.ringDeep}`,
    elevated: '0 4px 24px rgba(0, 0, 0, 0.05)'
  },
  reference: {
    spacing: [3, 4, 6, 8, 10, 12, 16, 20, 24, 30],
    breakpoints: [479, 640, 768, 992],
    focusBorder: palette.focusBlue
  },
  antTheme: {
    token: {
      ...semantic,
      ...linkTokens,
      fontFamily: uiFontFamily,
      borderRadius: 8
    },
    components: {
      Input: { borderRadius: 12 }
    }
  }
}

export default claudeTheme
