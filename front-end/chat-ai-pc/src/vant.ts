import 'vant/lib/index.css'
import type { App } from 'vue'
import { Button, ConfigProvider, Icon } from 'vant'

export const setupVant = (app: App<Element>) => {
  app.use(ConfigProvider)
  app.use(Button)
  app.use(Icon)
}
