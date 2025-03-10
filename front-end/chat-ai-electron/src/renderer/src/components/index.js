import SvgIcon from './svg-icon/index.vue'

export const registGlobalComponent = (app) => {
  // 全局注册组件
  app.component('SvgIcon', SvgIcon)
}
