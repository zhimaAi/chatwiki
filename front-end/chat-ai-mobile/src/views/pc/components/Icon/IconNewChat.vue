<template>
  <span class="icon-new-chat" :style="iconStyle"></span>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
  size?: number
  bgColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 14,
  bgColor: '#1890ff'
})

const iconStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  background: props.bgColor
}))
</script>

<style lang="less" scoped>
.icon-new-chat {
  display: inline-block;
  border-radius: 50%;
  position: relative;

  &::before,
  &::after {
    content: '';
    position: absolute;
    background: #fff;
    border-radius: 1px;
  }

  &::before {
    width: 57%;
    height: 2px;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
  }

  &::after {
    width: 2px;
    height: 57%;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
  }
}
</style>
