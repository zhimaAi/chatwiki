import { theme } from 'ant-design-vue'

const palette = {
  nearBlack: '#080808',
  blue: '#146ef5',
  blue400: '#3b89ff',
  blue300: '#006acc',
  buttonHoverBlue: '#0055d4',
  purple: '#7a3dff',
  pink: '#ed52cb',
  green: '#00d722',
  orange: '#ff6b00',
  yellow: '#ffae13',
  red: '#ee1d36',
  gray800: '#222222',
  gray700: '#363636',
  gray300: '#ababab',
  midGray: '#5a5a5a',
  borderGray: '#d8d8d8',
  borderHover: '#898989',
  white: '#ffffff'
}

// 背景层级及文字层级是 ChatWiki 的语义映射，色值均取自原文。
const semantic = {
  colorPrimary: palette.blue,
  colorPrimaryHover: palette.buttonHoverBlue,
  colorBgLayout: palette.white,
  colorBgContainer: palette.white,
  colorBgElevated: palette.white,
  colorText: palette.nearBlack,
  colorTextSecondary: palette.gray800,
  colorTextTertiary: palette.gray300,
  colorTextPlaceholder: palette.gray300,
  colorBorder: palette.borderGray,
  colorBorderSecondary: palette.borderGray
}

// 链接跟随品牌主色，缺省按下色使用官方算法；不改动信息提示色 colorInfo。
const primaryTokens = theme.defaultAlgorithm({ ...theme.defaultSeed, colorPrimary: semantic.colorPrimary })
const linkTokens = {
  colorLink: semantic.colorPrimary,
  colorLinkHover: semantic.colorPrimaryHover,
  colorLinkActive: primaryTokens.colorPrimaryActive
}

// 品牌字体仅记录为参考；应用回退沿用项目中文字体，不假定已有品牌字体资源。
const uiFontFamily = '"PingFang SC", "Noto Sans SC", sans-serif'

const webflowTheme = {
  id: 'webflow',
  name: '天青蓝',
  sources: ['https://fchangjun.github.io/awesome-design-md-cn/design-md/webflow/DESIGN.md'],
  palette,
  semantic,
  typography: {
    ui: uiFontFamily,
    reference: {
      fontFamily: 'WF Visual Sans Variable',
      fallback: 'Arial',
      codeFontFamily: 'Inconsolata',
      body: { fontSize: 16, fontWeight: [400, 500], lineHeight: 1.6, letterSpacing: -0.16 },
      button: { fontSize: 16, fontWeight: 500, lineHeight: 1.6, letterSpacing: -0.16 },
      caption: { fontSize: 14, fontWeight: [400, 500], lineHeight: [1.4, 1.6] },
      heading: { fontSize: 32, fontWeight: 500, lineHeight: 1.3 }
    }
  },
  radius: {
    scale: [2, 4, 8],
    circular: '50%',
    cardRange: [4, 8],
    badge: 4,
    // 原文没有统一控件圆角，按其功能元素范围选择 4px 作为应用基准。
    control: 4
  },
  shadow: {
    elevated: '0 84px 24px rgba(0, 0, 0, 0), 0 54px 22px rgba(0, 0, 0, 0.01), 0 30px 18px rgba(0, 0, 0, 0.04), 0 13px 13px rgba(0, 0, 0, 0.08), 0 3px 7px rgba(0, 0, 0, 0.09)'
  },
  reference: {
    spacing: [1, 2.4, 3.2, 4, 5.6, 6, 7.2, 8, 9.6, 12, 16, 24],
    breakpoints: [479, 768, 992],
    badgeBackgroundOpacity: 0.1,
    buttonHoverTranslate: 6
  },
  antTheme: {
    token: {
      ...semantic,
      ...linkTokens,
      fontFamily: uiFontFamily,
      borderRadius: 4
    }
  }
}

export default webflowTheme
