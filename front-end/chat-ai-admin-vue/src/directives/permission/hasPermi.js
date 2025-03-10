import { checkPermi } from '@/utils/permission'

function checkPermission(el, binding) {
  const { value } = binding

  if (value && Array.isArray(value) && value.length > 0) {
    if (!checkPermi(value)) {
      el.parentNode && el.parentNode.removeChild(el)
    }
  } else {
    throw new Error('请设置操作权限标签值')
  }
}

export default {
  mounted(el, binding) {
    checkPermission(el, binding)
  }
  // updated(el, binding) {
  //   checkPermission(el, binding)
  // }
}
