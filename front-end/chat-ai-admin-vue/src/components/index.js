import SvgIcon from './svg-icon/index.vue'
import CuPage from './page-common/cu-page.vue'
import PageMiniTitle from './page-common/page-mini-title.vue'

export const registGlobalComponent = (app) => {
  // 全局注册组件
  app.component('SvgIcon', SvgIcon)
  app.component('CuPage', CuPage)
  app.component('PageMiniTitle', PageMiniTitle)
}
