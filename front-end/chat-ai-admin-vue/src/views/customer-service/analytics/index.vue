<template>
  <section class="analytics-page">
    <CustomerPageHeader :active-key="tab" :title="t('analytics')" :tabs="headerTabs" @update:active-key="selectTab" />
    <div class="analytics-content">
      <StatisticalAnalysis v-if="tab === 'basic'" :robot-id="robot.id" :robot-key="robot.robot_key" />
      <HitStatistics v-else-if="tab === 'data'" :robot-id="robot.id" :robot-key="robot.robot_key" :filters="filters" embedded @unknown="showUnknown" />
      <Feedback v-else-if="tab === 'feedback'" :robot-id="robot.id" />
      <template v-else>
        <a-radio-group v-model:value="unknownTab" class="unknown-tabs"><a-radio-button value="stats">{{ t('unknown_issues') }}</a-radio-button><a-radio-button value="summary">{{ t('unknown_summary') }}</a-radio-button></a-radio-group>
        <a-alert v-if="filters.app_type || filters.app_id" type="info" :message="t('unknown_all_channels')" show-icon />
        <UnknownStats v-if="unknownTab === 'stats'" :robot-id="robot.id" :filters="filters" embedded />
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
const isAgent = computed(() => Number(robot.application_type) === 2)
const agentTabs = ['basic', 'data', 'feedback', 'unknown']
const selectedTab = () => isAgent.value && agentTabs.includes(route.query.tab) ? route.query.tab : 'data'
const tab = ref(selectedTab())
watch(() => route.query.tab, () => { if (isAgent.value) tab.value = selectedTab() })
function selectTab(value) {
  tab.value = value
  // 即使目标等于当前 URL 也发起导航，以取消尚未完成的旧标签跳转。
  if (isAgent.value) void router.replace({ query: { ...route.query, tab: value } })
}
const headerTabs = computed(() => [
  ...(isAgent.value ? [{ key: 'basic', label: tStatistics('tab_basic_statistics') }] : []),
  { key: 'data', label: t('data_statistics') },
  { key: 'feedback', label: t('qa_feedback') },
  { key: 'unknown', label: t('unknown_issues') }
])
const unknownTab = ref('stats')
const filters = ref({})
function showUnknown(value) { filters.value = { ...value }; unknownTab.value = 'stats'; selectTab('unknown') }
</script>
<style lang="less" scoped>
.analytics-page { flex: 1; display: flex; flex-direction: column; min-height: 0; margin: 8px 8px 8px 0; border-radius: 12px; overflow: hidden; background: #fff; }
.analytics-content { flex: 1; min-height: 0; overflow: auto; padding-top: 20px; }
.unknown-tabs { margin: 0 24px 16px; }
.analytics-content > .ant-alert { margin: 0 24px 16px; }
</style>
