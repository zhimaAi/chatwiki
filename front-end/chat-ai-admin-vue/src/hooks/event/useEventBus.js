import mitt from 'mitt'
import { onBeforeUnmount } from 'vue'

const emitter = mitt()

export const useEventBus = (option) => {
  if (option) {
    emitter.on(option.name, option.callback)

    onBeforeUnmount(() => {
      emitter.off(option.name)
    })
  }

  return {
    on: emitter.on,
    off: emitter.off,
    emit: emitter.emit,
    all: emitter.all
  }
}
