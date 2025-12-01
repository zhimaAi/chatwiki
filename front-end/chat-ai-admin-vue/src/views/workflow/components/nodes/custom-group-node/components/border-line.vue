<template>
  <div
    class="border-line line-left"
    :class="{ isSelected: isSelected, isHovered: isHovered }"
  ></div>
  <div
    class="border-line line-right"
    :class="{ isSelected: isSelected, isHovered: isHovered }"
  ></div>
  <div class="border-line line-top" :class="{ isSelected: isSelected, isHovered: isHovered }"></div>
  <div
    class="border-line line-botttom"
    :class="{ isSelected: isSelected, isHovered: isHovered }"
  ></div>
</template>

<script setup>
const props = defineProps({
  isSelected: {
    type: Boolean,
    default: false
  },
  isHovered: {
    type: Boolean,
    default: false
  }
})
</script>

<style lang="less" scoped>
.border-line {
  position: absolute;

  border-radius: 16px;
  &.line-left {
    left: 0;
    top: 0;
    bottom: 0;

    &.isHovered {
      width: 1px;
      background: #2475fc;
    }
    &.isSelected {
      width: 2px;
      background: #2475fc;
    }
  }
  &.line-right {
    right: 0;
    top: 0;
    bottom: 0;
    &.isHovered {
      width: 1px;
      background: #2475fc;
    }
    &.isSelected {
      width: 2px;
      background: #2475fc;
    }
  }
  &.line-top {
    left: 0;
    right: 0;
    top: 0;

    &.isHovered {
      height: 1px;
      background: #2475fc;
    }
    &.isSelected {
      height: 2px;
      background: #2475fc;
    }
  }
  &.line-botttom {
    left: 0;
    right: 0;
    bottom: 0;
    &.isHovered {
      height: 1px;
      background: #2475fc;
    }
    &.isSelected {
      height: 2px;
      background: #2475fc;
    }
  }
}
</style>
