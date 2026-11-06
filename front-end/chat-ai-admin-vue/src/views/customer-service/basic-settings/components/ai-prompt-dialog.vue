<template>
  <SettingsDialog v-if="open" :open="true" :width="880" :title="t('title_ai_generate_prompt')" @cancel="close">
    <template #closeIcon><img src="@/assets/img/customer-service/basic-settings/prompt/close.svg" alt="" /></template>
    <div class="ai-prompt-form">
      <div class="demand-row">
        <a-input v-model:value="demand" :placeholder="t('ph_input_prompt')" :aria-label="t('ph_input_prompt')" :disabled="loading" @press-enter="generate" />
        <a-button type="primary" :loading="loading" :disabled="!demand.trim()" @click="generate"><template #icon><img v-if="!loading" src="@/assets/img/customer-service/basic-settings/prompt/ai-generate.svg" alt="" /></template>{{ loading ? t('txt_generating') : t('btn_generate') }}</a-button>
      </div>
      <div class="quick-prompts"><span>{{ t('label_quick_generate').replace(/[：:]\s*$/, '') }}</span><button v-for="key in quickKeys" :key="key" type="button" :disabled="loading" @click="demand = t(key)">{{ t(key) }}</button></div>
      <div class="generated-prompt" :class="{ ready: result }" :aria-busy="loading" aria-live="polite">
        <a-spin v-if="loading" :tip="t('tip_generating')" />
        <template v-else-if="result"><section v-for="(section, index) in sections" :key="index"><div>{{ section.subject }}</div><PromptTokenText :value="section.describe" :options="tokenOptions" /></section></template>
        <a-alert v-else-if="error" type="error" :message="error" show-icon />
        <span v-else class="empty-prompt">{{ commonT('prompt_ai_empty') }}</span>
      </div>
    </div>
    <template #footer><a-button @click="close">{{ commonT('cancel') }}</a-button><a-button type="primary" :disabled="!result || loading" @click="apply">{{ t('btn_use_prompt') }}</a-button></template>
  </SettingsDialog>
</template>
<script setup>
import { computed, inject, onBeforeUnmount, ref } from 'vue'
import { createPromptByAi } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'
import { useRobotStore } from '@/stores/modules/robot'
import { useMiniCardStore } from '@/stores/modules/mini-card'
import { cloneConfig } from '../../config-session'
import { getPromptOptions } from '../prompt-tokens'
import PromptTokenText from './prompt-token-text.vue'
import SettingsDialog from './settings-dialog.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.ai-create-prompt')
const { t: commonT } = useI18n('views.customer-service.index')
const context = inject('customerRobotContext')
const robotStore = useRobotStore()
const miniCardStore = useMiniCardStore()
const emit = defineEmits(['handleAiSave'])
const open = ref(false)
const demand = ref('')
const loading = ref(false)
const error = ref('')
const result = ref(null)
const quickKeys = ['quick_ecommerce_after_sales', 'quick_english_education_pre_sales']
const fields = ['role', 'task', 'constraints', 'skill', 'output', 'tone']
const sections = computed(() => result.value ? [...fields.map(key => result.value[key]), ...result.value.custom].filter(section => section?.describe) : [])
const tokenOptions = computed(() => { const options = getPromptOptions(robotStore.chatVariables, miniCardStore.cardList); return [...options.variables, ...options.cards] })
let requestVersion = 0
let alive = true

function show() {
  requestVersion++
  demand.value = ''
  result.value = null
  loading.value = false
  error.value = ''
  open.value = true
}
function close() { requestVersion++; open.value = false; loading.value = false }
async function generate() {
  if (loading.value || !demand.value.trim()) return
  const version = ++requestVersion
  loading.value = true
  error.value = ''
  result.value = null
  try {
    const response = await createPromptByAi({ id: context.robotInfo.id, demand: demand.value })
    if (!alive || !open.value || version !== requestVersion) return
    if (response.res !== undefined && response.res !== 0) throw new Error(response.msg || commonT('prompt_ai_failed'))
    const raw = response.data?.promptStruct
    const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
    if (!parsed || !fields.some(key => parsed[key]?.describe)) throw new Error(commonT('prompt_ai_failed'))
    result.value = { ...Object.fromEntries(fields.map(key => [key, { subject: '', describe: '', ...parsed[key] }])), custom: Array.isArray(parsed.custom) ? parsed.custom : [] }
  } catch (cause) {
    if (alive && open.value && version === requestVersion) error.value = cause instanceof SyntaxError ? commonT('prompt_ai_failed') : cause.msg || cause.message || commonT('prompt_ai_failed')
  } finally {
    if (alive && version === requestVersion) loading.value = false
  }
}
function apply() {
  if (!result.value || loading.value) return
  emit('handleAiSave', cloneConfig(result.value))
  close()
}
onBeforeUnmount(() => { alive = false; requestVersion++ })
defineExpose({ show })
</script>
<style lang="less" scoped>
.ai-prompt-form { display: flex; flex-direction: column; gap: 24px; padding: 8px 19px 8px 8px; color: #222; font-size: 14px; line-height: 22px; }
.demand-row { display: flex; gap: 8px; :deep(.ant-input) { min-width: 0; height: 32px; border-radius: 6px; border-color: #d8d8d8; } :deep(.ant-btn) { display: inline-flex; align-items: center; gap: 6px; height: 32px; padding: 0 12px; border-radius: 6px; background: #146ef5; border-color: #146ef5; } }
.quick-prompts { display: flex; align-items: center; flex-wrap: wrap; gap: 8px; min-height: 32px; button { min-height: 32px; padding: 0 12px; border: 0; border-radius: 6px; background: #f5f5f5; color: #363636; font-size: 14px; cursor: pointer; } button:disabled { cursor: default; opacity: .6; } }
.generated-prompt { display: flex; align-items: center; justify-content: center; min-height: 240px; padding: 8px 12px; border-radius: 6px; background: #f5f5f5; overflow-wrap: anywhere; }
.generated-prompt.ready { display: block; min-height: 360px; max-height: 440px; overflow-y: auto; border: 1px solid #d8d8d8; background: #f5f9ff; section { margin-bottom: 22px; white-space: pre-wrap; } section:last-child { margin-bottom: 0; } }
.empty-prompt { color: #7a7a7a; text-align: center; }
</style>
