<style lang="less" scoped>
.phone-wrapper {
  height: 100%;
  width: 399px;
  border-radius: 32px;
  border: 12px solid #fff;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
  .phone-box {
    display: flex;
    flex-direction: column;
    position: relative;
    width: 100%;
    height: 100%;

    &::before {
      position: absolute;
      left: 0;
      right: 0;
      top: 0;
      bottom: 0;
      content: ' ';
      z-index: 101;
      width: 100%;
      height: 100%;
      border-radius: 28px;
      box-shadow: inset 0 0 4px rgba(0, 0, 0, 0.1);
      overflow: hidden;
      pointer-events: none;
    }

    .phone-header {
      .navbar {
        display: block;
      }
    }
    .phone-body {
      flex: 1;
      overflow: hidden;
    }
    .phone-footer {
    }
  }
}
</style>
<template>
  <div class="phone-wrapper">
    <div class="phone-box">
      <div class="phone-header">
        <img class="navbar" src="../../../../assets/img/library/preview/phone-header.png" alt="" />
      </div>
      <div class="phone-body">
        <slot name="default"></slot>
      </div>
      <div class="phone-footer"></div>
    </div>
  </div>
</template>

<script setup></script>
