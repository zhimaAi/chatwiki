import 'vant/lib/index.css'
import type { App } from 'vue'
import { Button, ConfigProvider, Icon, Popover, Popup, Field, Form, Loading, Dialog } from 'vant'

export const setupVant = (app: App<Element>) => {
  app.use(ConfigProvider)
  app.use(Button)
  app.use(Icon)
  app.use(Popover)
  app.use(Popup)
  app.use(Field)
  app.use(Form)
  app.use(Loading)
  app.use(Dialog)
}
