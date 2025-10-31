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
    font-family: 'PingFang SC';
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
    padding-top: 24px;
  }
}
</style>

<template>
  <div class="user-model-page">
    <!-- <div class="page-title">统计分析</div> -->
    <a-tabs class="tab-wrapper" @change="changeMenu" v-model:activeKey="activeKey">
      <a-tab-pane :key="1" tab="基础统计"></a-tab-pane>
      <a-tab-pane :key="3" tab="命中统计"></a-tab-pane>
    </a-tabs>
    <div class="list-wrapper">
      <StatisticalAnalysis v-if="activeKey === 1" />
      <HitStatics v-if="activeKey === 3" />
    </div>
  </div>
</template>

<script setup>
import StatisticalAnalysis from './statistical_analysis.vue'
import HitStatics from './hit-statics.vue'
import { ref } from 'vue'
const activeLocalKey = '/robot/config/statistical_analysis/activeKey'

const activeKey = ref(+localStorage.getItem(activeLocalKey) || 1)

const changeMenu = () => {
  localStorage.setItem(activeLocalKey, activeKey.value)
}
</script>

<style lang="less" scoped>
::v-deep(.ant-tabs-nav-wrap) {
  padding-left: 24px;
  background-color: #fff;
}
.tab-wrapper ::v-deep(.ant-tabs-nav) {
  margin-bottom: 0;
}
</style>
