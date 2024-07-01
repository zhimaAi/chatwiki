<template>
  <span class="svg-action">
    <svg aria-hidden="true" class="svg-icon" :fill="props.color" :width="width" :height="height">
      <use :xlink:href="symbolId" :fill="props.color" />
    </svg>
  </span>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({
  name: {
    type: String,
    required: true
  },
  size: {
    type: String,
    default: '1em'
  },
  width: {
    type: String
  },
  height: {
    type: String
  },
  color: {
    type: String,
    default: 'currentColor'
  }
})
const symbolId = computed(() => `#${props.name}`)
const width = computed(() => (props.width ? props.width : props.size))
const height = computed(() => (props.height ? props.height : props.size))
</script>

<style lang="less" scoped>
.svg-action {
  display: inline-flex;
  align-items: center;
  color: inherit;
  font-style: normal;
  line-height: 0;
  text-align: center;
  text-transform: none;
  vertical-align: -0.125em;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  .svg-icon {
    line-height: 1;
    display: inline-block;
    fill: currentColor;
    outline: none;

    &.default-font-size {
      font-size: 14px;
    }
  }
}
</style>
