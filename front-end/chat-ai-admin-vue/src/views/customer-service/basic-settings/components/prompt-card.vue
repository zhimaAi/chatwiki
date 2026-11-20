<template>
  <EditBox class="prompt-card" :title="t('prompt')" icon-name="prompt-icon" @edit="open">
    <div ref="preview" class="prompt-preview" :class="{ expanded }">
      <template v-if="Number(saved.prompt_type) === 1">
        <div v-for="(section, index) in sections" :key="index" class="prompt-section"><strong v-if="section.subject">#<PromptTokenText :value="section.subject" :options="tokenOptions" /></strong><div><PromptTokenText :value="section.describe" :options="tokenOptions" /></div></div>
        <span v-if="!sections.length" class="empty">{{ t('empty_prompt') }}</span>
      </template>
      <div v-else class="custom-prompt"><PromptTokenText :value="saved.prompt || t('empty_prompt')" :options="tokenOptions" /></div>
    </div>
    <div v-if="canExpand" class="expand"><a-button type="link" size="small" @click="expanded = !expanded">{{ t(expanded ? 'collapse' : 'expand') }}<UpOutlined v-if="expanded" /><DownOutlined v-else /></a-button></div>
    <PromptDialog v-if="editing" :robot-info="draft" :robot-list="robotList" :saving="saving" @save="save" @cancel="close" />
  </EditBox>
</template>

<script setup>
import { computed, inject, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { DownOutlined, UpOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import PromptDialog from './prompt-dialog.vue'
import EditBox from '@/views/robot/robot-config/basic-config/components/edit-box.vue'
import { getPromptConfig, saveConfigurationPatch } from '../configuration'
import { cloneConfig } from '../../config-session'
import { useRobotStore } from '@/stores/modules/robot'
import { useMiniCardStore } from '@/stores/modules/mini-card'
import { getPromptOptions } from '../prompt-tokens'
import PromptTokenText from './prompt-token-text.vue'

defineProps({ robotList: { type: Array, default: () => [] } })
const { t } = useI18n('views.customer-service.index')
const context = inject('customerRobotContext')
const robotStore = useRobotStore()
const miniCardStore = useMiniCardStore()
const tokenOptions = computed(() => { const options = getPromptOptions(robotStore.chatVariables, miniCardStore.cardList); return [...options.variables, ...options.cards] })
const saved = ref(getPromptConfig(context.robotInfo))
const draft = ref(null)
const editing = ref(false)
const saving = ref(false)
const preview = ref(null)
const expanded = ref(false)
const canExpand = ref(false)
let alive = true
let observer
const sections = computed(() => {
  const value = saved.value.prompt_struct || {}
  return [value.role, value.task, value.constraints, value.skill, value.output, value.tone, ...(value.custom || [])].filter(section => section?.describe?.trim())
})

function measure() { if (alive) canExpand.value = (preview.value?.scrollHeight || 0) > 144 }
watch(context.savedRevision, () => { saved.value = getPromptConfig(context.robotInfo) })
watch(saved, async () => { await nextTick(); measure() }, { deep: true })
function open() {
  draft.value = { ...cloneConfig(context.robotInfo), ...cloneConfig(saved.value) }
  editing.value = true
}
function close() { if (!saving.value) editing.value = false }
async function save(patch) {
  if (saving.value) return
  saving.value = true
  try {
    if (await saveConfigurationPatch(context, patch) && alive) {
      saved.value = getPromptConfig(context.robotInfo)
      editing.value = false
    }
  } finally { if (alive) saving.value = false }
}
onMounted(() => {
  Promise.resolve(miniCardStore.fetchCardList()).catch(() => {})
  measure()
  if (typeof ResizeObserver !== 'undefined') { observer = new ResizeObserver(measure); observer.observe(preview.value) }
})
onBeforeUnmount(() => { alive = false; observer?.disconnect() })
</script>

<style lang="less" scoped>
.prompt-preview { max-height: 144px; overflow: hidden; font-size: 14px; line-height: 24px; color: #595959; overflow-wrap: anywhere; }
.prompt-preview.expanded { max-height: none; }
.prompt-section { margin-bottom: 12px; white-space: pre-wrap; }
.prompt-section:last-child { margin-bottom: 0; }
.custom-prompt { white-space: pre-wrap; }
.empty { color: #999; }
.expand { text-align: center; margin-top: 12px; }
</style>
