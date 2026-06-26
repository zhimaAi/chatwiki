<template>
  <transition name="slide-down">
    <div class="scroll-to-bottom-btn" :class="{ 'is-loading': loading }" @click="handleClick" v-if="visible">
      <div class="bottom-btn">
        <svg-icon name="down-arrow" class="down-arrow-icon" />
      </div>
      
      <svg v-if="loading" class="loading-ring" viewBox="0 0 44 44">
        <circle class="loading-ring-track" cx="22" cy="22" r="18" />
        <circle class="loading-ring-indicator" cx="22" cy="22" r="18" />
      </svg>
    </div>
  </transition>
</template>

<script setup>
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
.scroll-to-bottom-btn {
  width: 40px;
  height: 40px;
  cursor: pointer;
  position: absolute;
  bottom: 20px;
  left: 50%;
  margin-left: -20px;

  .bottom-btn {
    width: 40px;
    height: 40px;
    border-radius: 40px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    border: 1px solid #fff;
    background: #fff;
    font-size: 16px;
    color: #659efc;
    box-shadow: 0 4px 16px 0 #0000001f;

    &:hover {
      border: 1px solid #659dfc;
    }
  }

  &.is-loading {
    .bottom-btn:hover{
      border: 1px solid #fff !important;
    }
  }

  .loading-ring {
    position: absolute;
    top: -6px;
    left: -6px;
    width: 52px;
    height: 52px;
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

/* 定义进入动画 */
.slide-down-enter-active {
  animation: slide-down-in 0.3s ease-in;
  position: absolute;
  z-index: 1;
}

/* 定义进入完成后的状态 */
.slide-down-enter-from {
  transform: translateY(150%);
}

/* 定义退出动画 */
.slide-down-leave-active {
  animation: slide-down-out 0.3s ease-out;
  position: absolute;
  z-index: 1;
}

/* 定义退出完成后的状态 */
.slide-down-leave-to {
  transform: translateY(150%);
}

@keyframes slide-down-in {
  from {
    transform: translateY(150%);
  }
  to {
    transform: translateY(0);
  }
}

@keyframes slide-down-out {
  from {
    transform: translateY(0);
  }
  to {
    transform: translateY(150%);
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
