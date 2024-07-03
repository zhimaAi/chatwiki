import copy from './v-copy/index'
import hasPermi from './permission/hasPermi'
export default function registGlobaDirective(app) {
  app.directive('copy', copy)

  app.directive('hasPermi', hasPermi)
}
