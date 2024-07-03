<template>
  <div class="color-picker-wrapper">
    <a-input :value="props.value" :readonly="true" placeholder="Basic usage" />
    <span class="color-picker-action">
      <span ref="colorPicker"></span>
    </span>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import '@simonwep/pickr/dist/themes/nano.min.css'
import Pickr from '@simonwep/pickr'

const emit = defineEmits(['update:value', 'change'])

const props = defineProps({
  value: {
    type: String,
    default: '#000000'
  }
})
const colorPicker = ref(null)

onMounted(() => {
  const pickr = Pickr.create({
    el: colorPicker.value,
    container: 'body',
    theme: 'nano',
    default: props.value,
    defaultRepresentation: 'HEX',
    swatches: [
      'rgba(244, 67, 54, 1)',
      'rgba(233, 30, 99, 0.95)',
      'rgba(156, 39, 176, 0.9)',
      'rgba(103, 58, 183, 0.85)',
      'rgba(63, 81, 181, 0.8)',
      'rgba(33, 150, 243, 0.75)',
      'rgba(3, 169, 244, 0.7)',
      'rgba(0, 188, 212, 0.7)',
      'rgba(0, 150, 136, 0.75)',
      'rgba(76, 175, 80, 0.8)',
      'rgba(139, 195, 74, 0.85)',
      'rgba(205, 220, 57, 0.9)',
      'rgba(255, 235, 59, 0.95)',
      'rgba(255, 193, 7, 1)'
    ],

    components: {
      // Main components
      preview: true,
      opacity: true,
      hue: true,

      // Input / output Options
      interaction: {
        hex: false,
        rgba: false,
        hsla: false,
        hsva: false,
        cmyk: false,
        input: true,
        clear: false,
        save: true
      }
    }
  })

  pickr.on('save', (color) => {
    let val = color.toHEXA().toString()
    emit('update:value', val)
    emit('change', val)
  })
})
</script>

<style lang="less" scoped>
.color-picker-wrapper {
  display: inline-block;
  position: relative;
  font-size: 0;
  .color-picker-action {
    position: absolute;
    right: 6px;
    top: 6px;
    width: 20px;
    height: 20px;
    border-bottom: 2px;
    box-shadow: 0 0 4px 0 rgba(0, 0, 0, 0.16);
    overflow: hidden;
  }
}
:deep(.pickr) .pcr-button {
  padding: 0;
}
</style>
