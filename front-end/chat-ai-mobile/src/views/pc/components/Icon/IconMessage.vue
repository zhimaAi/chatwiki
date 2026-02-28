<template>
  <span class="icon-message" :class="{ active: isActive }" :style="iconStyle"></span>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
  size?: number
  isActive?: boolean
  activeColor?: string
  inactiveColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 16,
  isActive: false,
  activeColor: '#1890ff',
  inactiveColor: '#bfbfbf'
})

const iconStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  background: props.isActive ? props.activeColor : props.inactiveColor
}))
</script>

<style lang="less" scoped>
.icon-message {
  display: inline-block;
  border-radius: 3px;
  position: relative;
  flex-shrink: 0;

  &::before {
    content: '';
    position: absolute;
    width: 50%;
    height: 37.5%;
    background: #fff;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    border-radius: 1px;
  }

  &::after {
    content: '';
    position: absolute;
    width: 0;
    height: 0;
    border-left: 3px solid transparent;
    border-right: 3px solid transparent;
    border-top: 4px solid currentColor;
    bottom: -3px;
    left: 50%;
    transform: translateX(-50%);
  }
}
</style>
