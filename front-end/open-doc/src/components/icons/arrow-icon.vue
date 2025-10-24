<style lang="less" scoped>
.arrow-icon {
  display: inline-block;
  transition: transform 0.2s ease;
  
  &.rotate-up {
    transform: rotate(180deg);
  }
  
  &.rotate-left {
    transform: rotate(90deg);
  }
  
  &.rotate-right {
    transform: rotate(-90deg);
  }
}
</style>

<template>
  <span class="arrow-icon" :class="rotationClass">
    <svg
      :width="size"
      :height="size"
      viewBox="0 0 16 16"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M12 6L8 10L4 6"
        :stroke="color"
        stroke-width="1.3"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
  </span>
</template>

<script>
export default {
  props: {
    direction: {
      type: String,
      default: 'bottom',
      validator: (value) => ['top', 'right', 'bottom', 'left'].includes(value)
    },
    color: {
      type: String,
      default: 'currentColor'
    },
    size: {
      type: [String, Number],
      default: 16
    }
  },
  computed: {
    rotationClass() {
      const directionMap = {
        top: 'rotate-up',
        left: 'rotate-left',
        right: 'rotate-right',
        bottom: ''
      }
      return directionMap[this.direction]
    }
  }
}
</script>