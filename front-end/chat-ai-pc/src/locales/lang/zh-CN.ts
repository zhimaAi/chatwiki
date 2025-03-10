import { genMessage } from '../helper'

const modulesFiles = import.meta.glob<Recordable>('./zh-CN/**/*.json', { eager: true })

export default {
  ...genMessage(modulesFiles, 'zh-CN')
}
