<template>
  <div class="empty-box">
    <img src="../../../../assets/img/library/preview/empty.png" alt="" />
    <div class="title">暂无分段</div>
    <div class="info">左侧进行分段设置后，点击“生成分段预览</div>
  </div>
</template>
<script setup>
</script>
<style lang="less" scoped>
.empty-box {
  padding-top: 70px;
  width: 100%;
  height: calc(100% - 20px);
  background: #f2f4f7;
  text-align: center;
  img {
    width: 200px;
    height: 200px;
  }
  .title {
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    color: #595959;
    margin-bottom: 8px;
  }

  .info {
    color: rgba(0, 0, 0, 0.45);
  }
}
</style>
