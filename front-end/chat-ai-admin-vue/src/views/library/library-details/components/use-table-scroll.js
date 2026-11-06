import { nextTick, ref, watch } from 'vue'

export function useTableScroll(container, enabled) {
  const scrollY = ref(1)

  watch([container, enabled], async ([element, active], _, onCleanup) => {
    let disposed = false
    let observer
    const measure = () => {
      if (disposed || !element || !active) return
      const header = element.querySelector('.ant-table-header')
      scrollY.value = Math.max(1, element.clientHeight - (header?.getBoundingClientRect().height || 0))
    }
    onCleanup(() => {
      disposed = true
      observer?.disconnect()
      window.removeEventListener('resize', measure)
    })
    if (!element || !active) return
    await nextTick()
    if (disposed) return
    measure()
    if (typeof ResizeObserver !== 'undefined') {
      observer = new ResizeObserver(measure)
      observer.observe(element)
      const header = element.querySelector('.ant-table-header')
      if (header) observer.observe(header)
    } else {
      window.addEventListener('resize', measure)
    }
  }, { immediate: true, flush: 'post' })

  return { scrollY }
}
