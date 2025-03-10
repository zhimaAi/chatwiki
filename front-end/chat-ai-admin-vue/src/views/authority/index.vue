<template>
  <div class="authority-page">
    <div class="center-content">
      <img src="../../assets/no-authority.png" alt="">
      <div class="title">无管理权限</div>
      <div class="desc">请联系管理员添加响应权限</div>
    </div>
  </div>
</template>

<script setup>

</script>

<style lang="less" scoped>
.authority-page{
  background: #fff;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  .center-content{
    margin-top: 140px;
    text-align: center;
  }
  img{
    width: 200px;
    height: 200px;
  }
  .title{
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }
  .desc{
    margin-top: 8px;
    font-weight: 400;
    font-size: 14px;
    line-height: 22px;
    color: #8c8c8c;
  }
}
</style>