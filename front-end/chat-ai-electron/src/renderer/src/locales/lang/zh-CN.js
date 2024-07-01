import { genMessage } from '../helper'

const modulesFiles = import.meta.glob('./zh-CN/**/*.json', { eager: true })

export default {
  ...genMessage(modulesFiles, 'zh-CN')
}
