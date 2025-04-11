<template>
  <div class="no-data">
    <svg-icon name="empty" :style="{'font-size': `${props.size}px`}" />
    <div class="no-data-text">
      <slot name="default">{{ props.text }}</slot>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  text: {
    type: String,
    default: '暂无数据'
  },
  size: {
    type: [Number, String],
    default: 100 
  }
});
</script>

<style lang="less" scoped>
.no-data {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  .no-data-text {
    margin-top: 16px;
    font-size: 14px;
    color: #999;
    text-align: center;
  }
}
</style>
