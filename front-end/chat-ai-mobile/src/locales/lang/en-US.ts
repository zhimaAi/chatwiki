import { genMessage } from '../helper'

const modulesFiles = import.meta.glob<Recordable>('/en-US/**/*.json', { eager: true })

export default {
  ...genMessage(modulesFiles, 'en-US')
}
