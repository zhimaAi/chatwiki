<style lang="less" scoped>
.node-form-layout {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.layout-header {
  padding: 16px 16px 0 16px;
}
.layout-body {
  flex: 1;
  padding: 0 16px 16px 16px;
  overflow-y: auto;
}
</style>

<template>
  <div class="node-form-layout">
    <div class="layout-header">
      <slot name="header"></slot>
    </div>
    <div class="layout-body">
      <slot name="default"></slot>
    </div>
  </div>
</template>

<script setup>

</script>
