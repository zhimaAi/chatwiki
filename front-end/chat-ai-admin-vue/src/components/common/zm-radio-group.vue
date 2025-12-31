<template>
  <div class="custom-radio-box">
    <div
      v-for="item in options"
      :key="item.value"
      class="custom-radio-item"
      :class="{ active: value === item.value }"
      @click="onSelect(item.value)"
    >
      {{ item.label }}
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  value: {
    type: [String, Number, Boolean],
    default: ''
  },
  options: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['update:value', 'change'])

function onSelect(val) {
  if (val !== props.value) {
    emit('update:value', val)
    emit('change', val)
  }
}
</script>
<style lang="less" scoped>
.custom-radio-box {
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  padding: 2px;
  gap: 4px;
  border-radius: 24px;
  background: #E4E6EB;
  height: 34px;

  .custom-radio-item {
    color: #3a4559;
    font-size: 12px;
    font-weight: 400;
    padding: 4px 12px;
    cursor: pointer;
    &.active,
    &:hover {
      color: #2475fc;
      border-radius: 24px;
      background: #FFF;
    }
  }
}
</style>
