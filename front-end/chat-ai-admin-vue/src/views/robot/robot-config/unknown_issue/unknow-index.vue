<template>
  <a-tabs class="tab-wrapper" v-model:activeKey="activeKey" @change="handleChangeTab">
    <template v-if="unknown_summary_status">
      <a-tab-pane :key="2" tab="未知问题总结"></a-tab-pane>
      <a-tab-pane :key="1" tab="未知问题统计"></a-tab-pane>
    </template>
    <template v-else>
      <a-tab-pane :key="1" tab="未知问题统计"></a-tab-pane>
      <a-tab-pane :key="2" tab="未知问题总结"></a-tab-pane>
    </template>
  </a-tabs>
  <div class="user-model-page">
    <cu-scroll>
      <UnknownIssue v-if="activeKey == 1" />
      <UnknownIssueSummarize v-else />
    </cu-scroll>
  </div>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { useRobotStore } from '@/stores/modules/robot'
import UnknownIssue from './index.vue'
import UnknownIssueSummarize from './summarize/index.vue'
const robotStore = useRobotStore()
const activeLocalKey = '/robot/config/unknown_issue/activeKey'
const activeKey = ref(+localStorage.getItem(activeLocalKey) || 1)

const robotInfo = computed(() => {
  return robotStore.robotInfo
})
const unknown_summary_status = computed(() => {
  return robotInfo.value.unknown_summary_status == 1
})

const handleChangeTab = () => {
  localStorage.setItem(activeLocalKey, activeKey.value)
}
</script>

<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: calc(100% - 46px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
::v-deep(.ant-tabs-nav-wrap) {
  padding-left: 24px;
}
.tab-wrapper ::v-deep(.ant-tabs-nav) {
  margin-bottom: 0;
}
</style>
