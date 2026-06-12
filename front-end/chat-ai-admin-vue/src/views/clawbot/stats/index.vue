<template>
  <div class="user-model-page">
    <a-tabs class="tab-wrapper" @change="changeMenu" v-model:activeKey="activeKey">
      <a-tab-pane :key="1" :tab="t('tab_basic_statistics')"></a-tab-pane>
      <a-tab-pane :key="3" :tab="t('tab_hit_statistics')"></a-tab-pane>
    </a-tabs>
    <div class="list-wrapper" v-if="robotId">
      <StatisticalAnalysis
        v-if="activeKey === 1"
        :robot-id="robotId"
        :robot-key="robotKey"
      />
      <HitStatics
        v-if="activeKey === 3"
        :robot-id="robotId"
        :robot-key="robotKey"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useClawbotStore } from '@/stores/modules/clawbot'
import StatisticalAnalysis from '@/views/robot/robot-config/statistical_analysis/statistical_analysis.vue'
import HitStatics from '@/views/robot/robot-config/statistical_analysis/hit-statics.vue'

// 复用 robot 模块的 i18n 命名空间，文案完全一致
const { t } = useI18n('views.robot.robot-config.statistical-analysis.index')

const clawbotStore = useClawbotStore()
const robotId = computed(() => clawbotStore.currentAssistant?.id || '')
const robotKey = computed(() => clawbotStore.currentAssistant?.robot_key || '')
const applicationType = computed(
  () => clawbotStore.currentAssistant?.application_type ?? 2
)

// activeKey 单独 localStorage key，避免与 robot 模块互相干扰
const activeLocalKey = '/clawbot/stats/activeKey'
const activeKey = ref(+localStorage.getItem(activeLocalKey) || 1)

const changeMenu = () => {
  localStorage.setItem(activeLocalKey, activeKey.value)
}
</script>

<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  border-bottom: 1px solid #fff;
  border-right: 1px solid #fff;
  background-color: #f2f4f7;

  .list-wrapper {
    background: #fff;
    height: calc(100% - 47px);
    overflow-x: hidden;
    overflow-y: auto;
    padding-top: 24px;
  }
}

::v-deep(.ant-tabs-nav-wrap) {
  padding-left: 24px;
  background-color: #fff;
}
.tab-wrapper ::v-deep(.ant-tabs-nav) {
  margin-bottom: 0;
}
</style>
