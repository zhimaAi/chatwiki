<template>
  <EditBox class="recall-card" :title="t('knowledge_recall_settings')" icon-name="guanlianzhishiku" @edit="open">
    <div class="summary workspace-summary">
      <div class="summary-item"><span class="summary-label">{{ summaryLabel('label_search_mode') }}</span><span class="summary-value">{{ searchMode }}</span></div>
      <div class="summary-item"><span class="summary-label">{{ summaryLabel('label_top_k') }}</span><span class="summary-value">{{ saved.top_k }}</span></div>
      <div v-if="Number(saved.search_type) <= 2" class="summary-item"><span class="summary-label">{{ summaryLabel('label_similarity_threshold') }}</span><span class="summary-value">{{ saved.similarity }}</span></div>
      <div class="summary-item"><span class="summary-label">{{ summaryLabel('label_rerank_model') }}</span><span class="summary-value">{{ Number(saved.rerank_status) === 1 ? modelName : recallT('label_off') }}</span></div>
    </div>
    <RecallSettingsAlert ref="dialog" variant="workspace" :submit="save" />
    <NoOpenGraphModal ref="graphWarning" :list="noGraphLibraries" @refresh-list="loadGraphLibraries" />
  </EditBox>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getLibraryList } from '@/api/library'
import { getModelNameText } from '@/components/model-select'
import RecallSettingsAlert from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/recall-settings-alert.vue'
import NoOpenGraphModal from '@/views/robot/robot-config/basic-config/components/associated-knowledge-base/no-open-graph-modal.vue'
import EditBox from '@/views/robot/robot-config/basic-config/components/edit-box.vue'
import { getRecallConfig, saveConfigurationPatch } from '../configuration'
import { cloneConfig } from '../../config-session'

const { t } = useI18n('views.customer-service.index')
const { t: recallT } = useI18n('views.robot.robot-config.basic-config.components.associated-knowledge-base.index')
const summaryLabel = key => recallT(key).replace(/[：:]\s*$/, '')
const context = inject('customerRobotContext')
const saved = ref(getRecallConfig(context.robotInfo))
const dialog = ref(null)
const graphWarning = ref(null)
const noGraphLibraries = ref([])
let alive = true
let requestVersion = 0
const searchMode = computed(() => recallT({ 1: 'label_hybrid_search', 2: 'label_vector_search', 3: 'label_fulltext_search', 4: 'label_graph_search' }[saved.value.search_type]))
const modelName = computed(() => getModelNameText(saved.value.rerank_model_config_id, saved.value.rerank_use_model, 'RERANK'))

watch(context.savedRevision, () => { saved.value = getRecallConfig(context.robotInfo) })
function open() { dialog.value.open(cloneConfig(saved.value), cloneConfig(context.robotInfo)) }
async function loadGraphLibraries() {
  const version = ++requestVersion
  const response = await getLibraryList({ type: '', show_open_docs: 1 })
  if (!alive || version !== requestVersion) return
  const ids = String(context.robotInfo.library_ids || '').split(',')
  noGraphLibraries.value = (response.data || []).filter(item => ids.includes(String(item.id)) && Number(item.graph_switch) === 0)
}
async function save(data) {
  const patch = getRecallConfig(data)
  patch.rrf_weight = JSON.stringify(patch.rrf_weight || {})
  const success = await saveConfigurationPatch(context, patch)
  if (!success || !alive) return false
  saved.value = getRecallConfig(context.robotInfo)
  if ([1, 4].includes(Number(saved.value.search_type))) {
    try {
      await loadGraphLibraries()
      if (alive && noGraphLibraries.value.length) graphWarning.value?.show()
    } catch { /* 图谱检查失败不改变已经保存的召回结果。 */ }
  }
  return true
}
onBeforeUnmount(() => { alive = false; requestVersion++ })
</script>

<style lang="less" scoped>
.summary { display: flex; flex-wrap: wrap; gap: 8px 16px; font-size: 14px; line-height: 22px; color: #595959; }
</style>
