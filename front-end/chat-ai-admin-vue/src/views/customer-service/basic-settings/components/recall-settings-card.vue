<template>
  <EditBox class="recall-card" :title="t('knowledge_recall_settings')" icon-name="guanlianzhishiku" @edit="open">
    <div class="summary workspace-summary">
      <div class="summary-item"><span class="summary-label">{{ summaryLabel('label_search_mode') }}</span><span class="summary-value">{{ searchMode }}</span></div>
      <div class="summary-item"><span class="summary-label">{{ summaryLabel('label_top_k') }}</span><span class="summary-value">{{ saved.top_k }}</span></div>
      <div v-if="Number(saved.search_type) <= 2" class="summary-item"><span class="summary-label">{{ summaryLabel('label_similarity_threshold') }}</span><span class="summary-value">{{ saved.similarity }}</span></div>
      <div class="summary-item"><span class="summary-label">{{ summaryLabel('label_rerank_model') }}</span><span class="summary-value">{{ Number(saved.rerank_status) === 1 ? modelName : recallT('label_off') }}</span></div>
    </div>
    <RecallSettingsAlert ref="dialog" variant="workspace" :submit="save" />
  </EditBox>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getModelNameText } from '@/components/model-select'
import RecallSettingsAlert from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/recall-settings-alert.vue'
import EditBox from '@/views/robot/robot-config/basic-config/components/edit-box.vue'
import { getRecallConfig, saveConfigurationPatch } from '../configuration'
import { cloneConfig } from '../../config-session'

const { t } = useI18n('views.customer-service.index')
const { t: recallT } = useI18n('views.robot.robot-config.basic-config.components.associated-knowledge-base.index')
const summaryLabel = key => recallT(key).replace(/[：:]\s*$/, '')
const context = inject('customerRobotContext')
const saved = ref(getRecallConfig(context.robotInfo))
const dialog = ref(null)
let alive = true
const searchMode = computed(() => ({ 1: recallT('label_hybrid_search'), 2: recallT('label_vector_search'), 3: recallT('label_fulltext_search') })[saved.value.search_type] || '')
const modelName = computed(() => getModelNameText(saved.value.rerank_model_config_id, saved.value.rerank_use_model, 'RERANK'))

watch(context.savedRevision, () => { saved.value = getRecallConfig(context.robotInfo) })
function open() { dialog.value.open(cloneConfig(saved.value), cloneConfig(context.robotInfo)) }
async function save(data) {
  const patch = getRecallConfig(data)
  patch.rrf_weight = JSON.stringify(patch.rrf_weight || {})
  const success = await saveConfigurationPatch(context, patch)
  if (!success || !alive) return false
  saved.value = getRecallConfig(context.robotInfo)
  return true
}
onBeforeUnmount(() => { alive = false })
</script>

<style lang="less" scoped>
.summary { display: flex; flex-wrap: wrap; gap: 8px 16px; font-size: 14px; line-height: 22px; color: #595959; }
</style>
