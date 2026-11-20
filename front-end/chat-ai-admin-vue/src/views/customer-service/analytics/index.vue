<template>
  <section class="analytics-page">
    <CustomerPageHeader :active-key="tab" :title="t('analytics')" :tabs="headerTabs" @update:active-key="selectTab" />
    <div :key="robot.id" class="analytics-content">
      <template v-if="tab === 'data'">
        <a-tabs class="statistics-tabs" :active-key="statisticsTab" @change="selectStatistics">
          <a-tab-pane key="basic" :tab="tStatistics('tab_basic_statistics')" />
          <a-tab-pane key="hit" :tab="tStatistics('tab_hit_statistics')" />
        </a-tabs>
        <StatisticalAnalysis v-if="statisticsTab === 'basic'" :robot-id="robot.id" :robot-key="robot.robot_key" />
        <HitStatistics v-else :robot-id="robot.id" :robot-key="robot.robot_key" :filters="filters" embedded @unknown="showUnknown" />
      </template>
      <Feedback v-else-if="tab === 'feedback'" :robot-id="robot.id" embedded />
      <template v-else>
        <a-alert v-if="filters.app_type || filters.app_id" type="info" :message="t('unknown_all_channels')" show-icon />
        <UnknownStats v-if="tab === 'unknown'" :robot-id="robot.id" :filters="filters" embedded />
        <UnknownSummary v-else :robot-id="robot.id" :filters="filters" embedded />
      </template>
    </div>
  </section>
</template>

<script setup>
import { computed, defineAsyncComponent, inject, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import HitStatistics from '@/views/robot/robot-config/statistical_analysis/hit-statics.vue'
import Feedback from '@/views/robot/robot-config/qa-feedback/index.vue'
import UnknownStats from '@/views/robot/robot-config/unknown_issue/index.vue'
import UnknownSummary from '@/views/robot/robot-config/unknown_issue/summarize/index.vue'
import CustomerPageHeader from '../components/page-header.vue'
const StatisticalAnalysis = defineAsyncComponent(() => import('@/views/robot/robot-config/statistical_analysis/statistical_analysis.vue').then(module => module.default))
const { t } = useI18n('views.customer-service.index')
const { t: tStatistics } = useI18n('views.robot.robot-config.statistical-analysis.index')
const route = useRoute()
const router = useRouter()
const robot = inject('customerRobotContext').robotInfo
const pageTabs = ['data', 'feedback', 'unknown', 'summary']
const selectedTab = () => pageTabs.includes(route.query.tab) ? route.query.tab : 'data'
// 旧 Agent 链接 tab=basic 仍落到数据统计的基础统计。
const selectedStatistics = () => route.query.tab === 'basic' ? 'basic' : route.query.stats === 'hit' ? 'hit' : 'basic'
const tab = ref(selectedTab())
const statisticsTab = ref(selectedStatistics())
watch(() => [route.query.tab, route.query.stats], () => {
  tab.value = selectedTab()
  statisticsTab.value = selectedStatistics()
})
function selectTab(value) {
  tab.value = value
  // 即使目标等于当前 URL 也发起导航，以取消尚未完成的旧标签跳转。
  void router.replace({ query: { ...route.query, tab: value, stats: statisticsTab.value } })
}
function selectStatistics(value) {
  statisticsTab.value = value
  void router.replace({ query: { ...route.query, tab: 'data', stats: value } })
}
const headerTabs = computed(() => [
  { key: 'data', label: t('data_statistics') },
  { key: 'feedback', label: t('qa_feedback') },
  { key: 'unknown', label: t('unknown_issues') },
  { key: 'summary', label: t('unknown_summary') }
])
const filters = ref({})
function showUnknown(value) { filters.value = { ...value }; selectTab('unknown') }
</script>
<style lang="less" scoped>
.analytics-page { flex: 1; display: flex; flex-direction: column; min-height: 0; margin: 8px 8px 8px 0; border-radius: 12px; overflow: hidden; background: #fff; }
.analytics-content { flex: 1; min-height: 0; overflow: auto; padding-top: 20px; }
.statistics-tabs { padding: 0 24px; }
.analytics-content > .ant-alert { margin: 0 24px 16px; }
</style>
