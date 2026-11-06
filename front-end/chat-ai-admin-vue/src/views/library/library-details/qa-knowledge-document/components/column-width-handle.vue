<template>
  <span class="column-width-handle" role="separator" tabindex="0" aria-orientation="vertical" :aria-label="label" :aria-valuenow="width" :aria-valuemin="minWidth"
    @pointerdown.stop.prevent="startResize" @click.stop @dblclick.stop.prevent="$emit('reset')" @keydown.stop="handleKeydown">
    <img :src="icon" alt="" draggable="false" />
  </span>
</template>

<script setup>
import { onBeforeUnmount } from 'vue'
import icons from '../knowledge-icons'

const props = defineProps({ width: { type: Number, required: true }, minWidth: { type: Number, default: 64 }, label: { type: String, required: true }, icon: { type: String, default: icons.resize } })
const emit = defineEmits(['resize', 'reset'])
let cleanup
function startResize(event) {
  if (event.button !== 0) return
  cleanup?.()
  const startX = event.clientX
  const startWidth = props.width
  const cursor = document.body.style.cursor
  const userSelect = document.body.style.userSelect
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
  const move = event => emit('resize', Math.max(props.minWidth, startWidth + event.clientX - startX))
  cleanup = () => {
    document.removeEventListener('pointermove', move)
    document.removeEventListener('pointerup', stop)
    document.removeEventListener('pointercancel', stop)
    document.body.style.cursor = cursor
    document.body.style.userSelect = userSelect
    cleanup = undefined
  }
  const stop = () => cleanup?.()
  document.addEventListener('pointermove', move)
  document.addEventListener('pointerup', stop)
  document.addEventListener('pointercancel', stop)
}
function handleKeydown(event) {
  if (event.key === 'ArrowLeft' || event.key === 'ArrowRight') {
    event.preventDefault()
    emit('resize', Math.max(props.minWidth, props.width + (event.key === 'ArrowLeft' ? -8 : 8)))
  }
}
onBeforeUnmount(() => cleanup?.())
</script>

<style scoped lang="less">
.column-width-handle {
  position: absolute; top: 0; right: -6px; z-index: 2; display: flex; align-items: center; justify-content: center; width: 12px; height: 100%; cursor: col-resize; touch-action: none;
  &:focus-visible { outline: 2px solid #146ef5; outline-offset: -2px; }
}
</style>
