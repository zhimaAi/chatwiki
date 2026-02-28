<template>
  <div class="icon-smile-wrapper" :style="wrapperStyle">
    <span class="icon-smile" :style="iconStyle"></span>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
  wrapperSize?: number
  iconSize?: number
  bgColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  wrapperSize: 80,
  iconSize: 48,
  bgColor: '#ff8c00'
})

const wrapperStyle = computed(() => ({
  width: `${props.wrapperSize}px`,
  height: `${props.wrapperSize}px`,
  background: props.bgColor
}))

const iconStyle = computed(() => ({
  width: `${props.iconSize}px`,
  height: `${props.iconSize}px`,
  '--eye-size': `${props.iconSize * 0.167}px`,
  '--eye-top': `${props.iconSize * 0.292}px`,
  '--eye-left': `${props.iconSize * 0.25}px`,
  '--eye-right': `${props.iconSize * 0.25}px`,
  '--mouth-width': `${props.iconSize * 0.5}px`,
  '--mouth-height': `${props.iconSize * 0.25}px`,
  '--mouth-border': `${props.iconSize * 0.0625}px`,
  '--mouth-top': `${props.iconSize * 0.167}px`
}))
</script>

<style lang="less" scoped>
.icon-smile-wrapper {
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;

  // 嘴巴
  &::after {
    content: '';
    position: absolute;
    width: var(--mouth-width);
    height: var(--mouth-height);
    border: var(--mouth-border) solid #fff;
    border-top: none;
    border-radius: 0 0 var(--mouth-width) var(--mouth-width);
    margin-top: var(--mouth-top);
  }
}

.icon-smile {
  display: inline-block;
  background: #fff;
  border-radius: 50%;
  position: relative;

  // 左眼
  &::before {
    content: '';
    position: absolute;
    width: var(--eye-size);
    height: var(--eye-size);
    background: #ff8c00;
    border-radius: 50%;
    top: var(--eye-top);
    left: var(--eye-left);
  }

  // 右眼
  &::after {
    content: '';
    position: absolute;
    width: var(--eye-size);
    height: var(--eye-size);
    background: #ff8c00;
    border-radius: 50%;
    top: var(--eye-top);
    right: var(--eye-right);
  }
}
</style>
