<template>
  <span class="icon-send" :class="{ active: isActive }" :style="iconStyle">
    <svg fill="none" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" class="design-iconfont">
      <rect width="32" height="32" rx="16" :fill="currentBgColor"/>
      <path d="M16 7.75L16 24.25" :stroke="currentIconColor" stroke-width="1.83333" stroke-linecap="round" stroke-linejoin="round"/>
      <path d="M10.5 13.25L16 7.75L21.5 13.25" :stroke="currentIconColor" stroke-width="1.83333" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
  </span>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
  size?: number
  bgColor?: string
  activeBgColor?: string
  iconColor?: string
  activeIconColor?: string
  isActive?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 32,
  bgColor: '#D8DDE6',
  activeBgColor: '#2475FC',
  iconColor: '#fff',
  activeIconColor: '#fff',
  isActive: false
})

const currentBgColor = computed(() =>
  props.isActive ? props.activeBgColor : props.bgColor
)

const currentIconColor = computed(() =>
  props.isActive ? props.activeIconColor : props.iconColor
)

const iconStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`
}))
</script>

<style lang="less" scoped>
.icon-send {
  display: inline-block;
}
</style>
