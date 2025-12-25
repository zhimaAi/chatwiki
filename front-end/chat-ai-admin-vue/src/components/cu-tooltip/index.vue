<template>
  <ATooltip v-bind="filteredAttrs" v-if="!disabled">
    <slot></slot>
  </ATooltip>
  <div v-else v-bind="wrapperAttrs">
    <slot></slot>
  </div>
</template>

<script setup>
import { useAttrs, computed } from 'vue'
import { Tooltip as ATooltip } from 'ant-design-vue'

//  定义组件名称以便调试
const props = defineProps({
  disabled: {
    type: Boolean,
    default: false
  }
})

const attrs = useAttrs()

// 过滤掉disabled属性，避免传递给a-tooltip
const filteredAttrs = computed(() => {
  const { disabled, ...rest } = attrs
  return rest
})

// 为禁用状态下的包装元素准备属性（排除tooltip相关属性）
const wrapperAttrs = computed(() => {
  const { title, overlayClassName, overlayStyle, mouseEnterDelay, mouseLeaveDelay, ...rest } = attrs
  return rest
})
</script>
