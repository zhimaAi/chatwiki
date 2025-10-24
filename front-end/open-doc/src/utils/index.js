/**
 * 把对象转为formData
 */
export function objToFormData(obj) {
  const formData = new FormData()
  Object.keys(obj).forEach((key) => {
    if (obj[key] !== void 0 && obj[key] !== null) {
      formData.append(key, obj[key])
    }
  })
  return formData
}
