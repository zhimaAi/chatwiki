<style lang="less" scoped>
.cu-tabs {
  display: flex;
  align-items: center;
  padding: 2px;
  font-size: 14px;
  color: #262626;
  border-radius: 6px;
  background-color: #f0f2f6;

  .tab-item {
    position: relative;
    line-height: 22px;
    padding: 5px 16px;
    border-radius: 6px;
    cursor: pointer;

    &.active {
      color: #2475fc;
      background-color: #FFFFFF;
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
      :class="{ active: props.value === tab.value }"
      @click="handleClickTab(tab.value)"
      v-for="tab in tabs"
      :key="tab.value"
    >
      {{ tab.title }}
      <slot name="extra" :tab="tab" :key="tab.value"></slot>
    </div>
  </div>
</template>

<script setup>
const emit = defineEmits(['update:value', 'change'])

const props = defineProps({
  value: {
    type: [Number, String],
    required: true
  },
  tabs: {
    type: Array,
    required: true,
    default: () => []
  }
})

const handleClickTab = (val) => {
  emit('update:value', val)
  emit('change', val)
}
</script>
