<style lang="less" scoped>
.page-mini-title {
  height: 38px;
  line-height: 38px;
  padding: 0 16px;
  border-radius: 2px;
  background-color: #f2f4f7;
  margin-bottom: 24px;
}
</style>

<template>
  <div class="page-mini-title">
    <slot>{{ title }}</slot>
  </div>
</template>
<script>
export default {
  props: {
    title: {
      type: String,
      default: ''
    }
  }
}
</script>
