import type { App } from 'vue'
import SvgIcon from './svg-icon/index.vue'

export const registGlobalComponent = (app: App<Element>): void => {
  // 全局注册组件
  app.component('SvgIcon', SvgIcon)
}
