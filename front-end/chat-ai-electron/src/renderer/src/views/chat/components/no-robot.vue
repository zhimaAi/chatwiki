<style lang="less" scoped>
.no-robot {
  width: 100%;
  height: 100%;
  display: flex;
  padding-top: 12.1%;
  justify-content: center;
  background-color: #f0f2f5;

  .no-data {
    text-align: center;
  }

  .tip-text-1 {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }
  .tip-text-2 {
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    margin-top: 4px;
    color: #7a8699;
  }
}
</style>

<template>
  <div class="no-robot">
    <div class="no-data">
      <img class="no-data-img" src="@/assets/img/no_robot.svg" alt="" />
      <div class="tip-text-1">暂无可用机器人</div>
      <div class="tip-text-2">请先去网页端创建智能机器人</div>
    </div>
  </div>
</template>

<script setup></script>
