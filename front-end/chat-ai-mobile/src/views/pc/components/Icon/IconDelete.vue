<template>
  <span class="icon-delete" :style="iconStyle"></span>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
  size?: number
  bgColor?: string
  lineColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 36,
  bgColor: '#f5f5f5',
  lineColor: '#999'
})

const iconStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  background: props.bgColor,
  '--line-color': props.lineColor
}))
</script>

<style lang="less" scoped>
.icon-delete {
  display: inline-block;
  border-radius: 6px;
  position: relative;
  cursor: pointer;
  transition: background 0.2s;

  &::before {
    content: '';
    position: absolute;
    width: 40%;
    height: 2px;
    background: var(--line-color);
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    border-radius: 1px;
  }

  &:hover {
    opacity: 0.8;
  }
}
</style>
