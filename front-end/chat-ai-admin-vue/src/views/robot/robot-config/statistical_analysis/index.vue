<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  border-bottom: 1px solid #fff;
  border-right: 1px solid #fff;
  background-color: #f2f4f7;

  .page-title {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 24px 24px 16px;
    background-color: #fff;
    color: #000000;
    font-family: "PingFang SC";
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }

  .list-wrapper {
    background: #fff;
    height: calc(100% - 47px);
    overflow-x: hidden;
    overflow-y: auto;
  }
}
</style>

<template>
  <div class="user-model-page">
    <div class="page-title">统计分析</div>
    <div class="list-wrapper">
      <StatisticalAnalysis></StatisticalAnalysis>
    </div>
  </div>
</template>

<script setup>
import StatisticalAnalysis from './statistical_analysis.vue'
</script>
