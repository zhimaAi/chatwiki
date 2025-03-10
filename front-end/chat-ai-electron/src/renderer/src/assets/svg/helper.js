const icons = import.meta.glob('./**/*.svg', { eager: true })

export function getAllIcons(prefix = '') {
  let list = []

  Object.keys(icons).forEach((key) => {
    //const iconFileModule = icons[key].default

    let fileName = key.replace(`./${prefix}/`, '').replace(/^\.\//, '')
    const lastIndex = fileName.lastIndexOf('.')

    fileName = fileName.substring(0, lastIndex)

    list.push({
      name: fileName
    })
  })

  return list
}
