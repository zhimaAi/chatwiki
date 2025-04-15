<style lang="less" scoped>
.cu-tabs {
  display: flex;
  align-items: center;

  .tab-item {
    position: relative;
    line-height: 24px;
    padding: 16px 0;
    margin-right: 32px;
    border-radius: 6px;
    font-size: 16px;
    font-weight: 400;
    color: #595959;
    cursor: pointer;

    &.active {
      font-weight: 600;
      color: #2475fc;
      &::after {
        content: '';
        position: absolute;
        left: 0;
        right: 0;
        bottom: 0;
        height: 2px;
        background-color: #2475fc;
      }
    }

    &:hover {
      color: #2475fc;
    }
  }
}
</style>

<template>
  <div class="cu-tabs">
    <div
      class="tab-item"
      :class="{ active: props.active === tab.path }"
      @click="handleClickTab(tab)"
      v-for="tab in tabs"
      :key="tab.path"
    >
      {{ tab.title }}
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
const emit = defineEmits(['change'])

const props = defineProps({
  active: {
    type: [Number, String],
    required: true
  },
  tabs: {
    type: Array,
    required: true,
    default: () => []
  }
})

const handleClickTab = (tab) => {
  emit('change', tab)
  if(tab.path){
    router.push(tab.path)
  }
}
</script>