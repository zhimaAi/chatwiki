<style lang="less" scoped>
.collapse-btn{
  position: absolute;
  display: flex;
  align-items: center;
  height: 30px;
  padding: 0 12px;
  top: 52px;
  left: 16px;
  padding: 4px 12px;
  border-radius: 6px;
  cursor: pointer;
  background-color: #fff;

  .collapse-icon{
    margin-right: 8px;
    font-size: 16px;
    color: #595959;
  }

  .collapse-label{
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: #262626;
  }
}
</style>

<template>
  <div class="collapse-btn" @click="toggleCollapse">
    <span class="collapse-icon">
      <svg-icon name="collapse-close" v-if="value"></svg-icon>
      <svg-icon name="collapse-open" v-else></svg-icon>
    </span>
    <span class="collapse-label">{{value ? '收起' : '展开'}}</span>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['change'])

const value = ref(true)

const toggleCollapse = () => {
  value.value = !value.value

  emit('change', value.value)
}
</script>
