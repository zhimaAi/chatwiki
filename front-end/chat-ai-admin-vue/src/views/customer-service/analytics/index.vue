<template>
  <section class="analytics-page">
    <CustomerPageHeader v-model:active-key="tab" :title="t('analytics')" :tabs="headerTabs" />
    <div class="analytics-content">
      <HitStatistics v-if="tab === 'data'" :robot-id="robot.id" :robot-key="robot.robot_key" :filters="filters" embedded @unknown="showUnknown" />
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
import { computed, inject, ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import HitStatistics from '@/views/robot/robot-config/statistical_analysis/hit-statics.vue'
import Feedback from '@/views/robot/robot-config/qa-feedback/index.vue'
import UnknownStats from '@/views/robot/robot-config/unknown_issue/index.vue'
import UnknownSummary from '@/views/robot/robot-config/unknown_issue/summarize/index.vue'
import CustomerPageHeader from '../components/page-header.vue'
const { t } = useI18n('views.customer-service.index')
const robot = inject('customerRobotContext').robotInfo
const tab = ref('data')
const headerTabs = computed(() => [
  { key: 'data', label: t('data_statistics') },
  { key: 'feedback', label: t('qa_feedback') },
  { key: 'unknown', label: t('unknown_issues') }
])
const unknownTab = ref('stats')
const filters = ref({})
function showUnknown(value) { filters.value = { ...value }; unknownTab.value = 'stats'; tab.value = 'unknown' }
</script>
<style lang="less" scoped>
.analytics-page { flex: 1; display: flex; flex-direction: column; min-height: 0; margin: 8px 8px 8px 0; border-radius: 12px; overflow: hidden; background: #fff; }
.analytics-content { flex: 1; min-height: 0; overflow: auto; padding-top: 20px; }
.unknown-tabs { margin: 0 24px 16px; }
.analytics-content > .ant-alert { margin: 0 24px 16px; }
</style>
