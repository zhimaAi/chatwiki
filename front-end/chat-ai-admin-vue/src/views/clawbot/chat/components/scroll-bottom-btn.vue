<template>
  <button
    class="scroll-bottom-btn"
    :class="{ 'is-loading': loading }"
    type="button"
    v-if="visible"
    @click="handleClick"
  >
    <ArrowDownOutlined />
    <svg v-if="loading" class="loading-ring" viewBox="0 0 44 44">
      <circle class="loading-ring-track" cx="22" cy="22" r="18" />
      <circle class="loading-ring-indicator" cx="22" cy="22" r="18" />
    </svg>
  </button>
</template>

<script setup>
import { ArrowDownOutlined } from '@ant-design/icons-vue'

defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])

const handleClick = () => {
  emit('click')
}
</script>

<style lang="less" scoped>
.scroll-bottom-btn {
  position: absolute;
  left: 50%;
  bottom: 16px;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  padding: 0;
  border: 0;
  border-radius: 50%;
  background: #fff;
  color: #000;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.12);
  transform: translateX(-50%);
  cursor: pointer;

  :deep(.anticon) {
    font-size: 18px;
  }

  .loading-ring {
    position: absolute;
    top: -4px;
    left: -4px;
    width: 44px;
    height: 44px;
    animation: loading-spin 1s linear infinite;
    pointer-events: none;

    circle {
      fill: none;
      stroke-width: 3;
    }
    .loading-ring-track {
      stroke: #e0e0e0;
    }
    .loading-ring-indicator {
      stroke: #659efc;
      stroke-linecap: round;
      stroke-dasharray: 100;
      stroke-dashoffset: 70;
    }
  }
}

@keyframes loading-spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
