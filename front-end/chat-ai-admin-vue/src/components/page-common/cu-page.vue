<template>
  <div class="cu-page" :style="{ backgroundColor: props.bgColor, ...props.pageSyle }">
    <div class="page-body" :style="{ ...props.pageBodyStyle }">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  bgColor: {
    type: String,
    default: '#fff'
  },
  pageSyle: {
    type: Object,
    default: () => {
      return {}
    }
  },
  pageBodyStyle: {
    type: Object,
    default: () => {
      return {}
    }
  },
  loading: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  }
})
</script>

<style lang="less" scoped>
.cu-page {
  border-radius: 2px;
}
</style>
