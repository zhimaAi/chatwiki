<template>
  <span class="icon-attachment" :style="iconStyle"></span>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
  size?: number
  bgColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 20,
  bgColor: '#bfbfbf'
})

const iconStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  background: props.bgColor
}))
</script>

<style lang="less" scoped>
.icon-attachment {
  display: inline-block;
  border-radius: 4px;
  position: relative;
  cursor: pointer;
  transition: background 0.2s;

  &::before {
    content: '';
    position: absolute;
    width: 50%;
    height: 50%;
    border: 2px solid #fff;
    border-radius: 50% 50% 0 50%;
    transform: rotate(-45deg);
    top: 50%;
    left: 50%;
    margin-top: -30%;
    margin-left: -30%;
  }

  &:hover {
    background: #999 !important;
  }
}
</style>
