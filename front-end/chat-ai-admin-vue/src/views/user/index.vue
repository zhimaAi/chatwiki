<style lang="less" scoped>
.user-page-wrapper {
  display: flex;
  margin-top: 16px;
  border-radius: 2px;
  width: 100%;
  height: 100%;
  overflow: hidden;

  .page-left {
    height: 100%;
    width: 255px;
    background-color: #fff;
    border-right: 1px solid #f2f4f7;
  }
  .page-body {
    flex: 1;
    height: 100%;
    overflow: hidden;
  }
}
</style>

<template>
  <div class="user-page-wrapper">
    <div class="page-left">
      <div class="left-menus-box">
        <LeftMenus />
      </div>
    </div>
    <div class="page-body">
      <router-view></router-view>
    </div>
  </div>
</template>

<script setup>
import LeftMenus from './components/left-menus.vue'
</script>
