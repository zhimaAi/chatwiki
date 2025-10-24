<template>
  <div class="error-page">
    <img src="@/assets/img/404.svg" alt="" />
    <!-- <h1 class="error-title">404</h1> -->
    <p class="error-desc">抱歉，您访问的页面不存在</p>
    <!-- <a href="/home" class="back-home">返回首页</a> -->
  </div>
</template>

<script setup>
// 404页面组件
</script>

<style lang="less" scoped>
.error-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  text-align: center;
  padding: 0 20px;
}
.error-page img {
  max-width: 360px;
  width: 100%;
  height: auto;
  margin-bottom: 24px;
}
.error-title {
  font-size: 120px;
  font-weight: bold;
  color: #262626;
  margin: 0;
  line-height: 1;
}
.error-desc {
  font-size: 18px;
  color: #595959;
  margin: 24px 0 32px;
}
.back-home {
  display: inline-block;
  padding: 0 24px;
  height: 40px;
  line-height: 40px;
  background: #1677ff;
  border-radius: 6px;
  color: #fff;
  text-decoration: none;
  font-size: 16px;
  transition: background-color 0.3s;
}
.back-home:hover {
  background: #4096ff;
}
</style>
