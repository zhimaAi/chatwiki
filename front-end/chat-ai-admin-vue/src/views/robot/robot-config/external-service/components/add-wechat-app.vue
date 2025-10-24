<style lang="less" scoped>
.add-wechat-app {
  width: 325px;
  height: 74px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  border-radius: 2px;
  font-size: 14px;
  border: 1px solid #d9d9d9;
  color: #595959;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
  &:hover {
    color: #1890ff;
    box-shadow:
      0 3px 14px 2px #0000000d,
      0 8px 10px 1px #0000000f,
      0 5px 5px -3px #0000001a;
  }
  .add-icon {
    font-size: 16px;
    margin-right: 4px;
  }
}
</style>

<template>
  <div class="add-wechat-app" @click="handleClick">
    <PlusOutlined class="add-icon" />
    <span class="add-label">{{ props.label }}</span>
  </div>
</template>

<script setup>
import { PlusOutlined } from '@ant-design/icons-vue'

const emit = defineEmits(['click'])
const props = defineProps({
  label: {
    type: String,
    default: ''
  }
})

const handleClick = () => {
  emit('click')
}
</script>
