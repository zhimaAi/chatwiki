import copy from './v-copy/index'
import hasPermi from './permission/hasPermi'
import { timeFormatter } from "./ftime"
export default function registGlobaDirective(app) {
  app.directive('copy', copy)

  app.directive('hasPermi', hasPermi)
  
  app.directive('ftime', timeFormatter);
}
