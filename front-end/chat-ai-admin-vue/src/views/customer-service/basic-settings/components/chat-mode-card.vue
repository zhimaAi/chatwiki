<template>
  <EditBox class="chat-mode-card" :title="commonT('chat_mode')" icon-name="chat-mode" @edit="open">
    <template #tip>
      <a-tooltip :title="summaryT('tip_knowledge_base_only') + '\n' + summaryT('tip_direct_mode')"><QuestionCircleOutlined /></a-tooltip>
    </template>
    <div class="summary workspace-summary" style="--summary-columns: 3">
      <div class="summary-item"><span class="summary-label">{{ summaryLabel('label_chat_mode') }}</span><span class="summary-value">{{ summaryT(savedMode.summary) }}</span></div>
      <template v-if="savedMode.switchKey">
        <div class="summary-item"><span class="summary-label">{{ summaryLabel('label_direct_reply_qa') }}</span><span class="summary-value summary-status" :class="{ 'is-enabled': saved[savedMode.switchKey] }">{{ commonT(saved[savedMode.switchKey] ? 'setting_on' : 'setting_off') }}</span></div>
        <div v-if="saved[savedMode.switchKey]" class="summary-item"><span class="summary-label">{{ summaryLabel('label_similarity_threshold') }}</span><span class="summary-value">{{ saved[savedMode.scoreKey] }}</span></div>
      </template>
    </div>
    <SettingsDialog v-if="editing" :open="true" :width="1120" :title="t('title_chat_mode_settings')" :ok-text="commonT('confirm')" :cancel-text="commonT('cancel')" :confirm-loading="saving" :closable="!saving" :mask-closable="!saving" :keyboard="!saving" :cancel-button-props="{ disabled: saving }" @ok="persist" @cancel="close">
      <div class="chat-mode-form" :inert="saving ? '' : null">
        <div class="mode-options" :class="{ cropped: form.chat_type !== '2' }" role="radiogroup" :aria-label="t('title_chat_mode_settings')">
          <button v-for="(mode, index) in modes" :key="mode.value" :ref="element => modeButtons[index] = element" type="button" class="mode-option" :class="{ active: form.chat_type === mode.value }" role="radio" :aria-checked="form.chat_type === mode.value" :tabindex="form.chat_type === mode.value ? 0 : -1" :disabled="saving" @click="selectMode(mode.value)" @keydown="event => navigateMode(event, index)">
            <img class="mode-illustration" :src="mode.image" alt="" />
            <strong class="mode-title">{{ t(mode.title) }}</strong>
            <p class="description">{{ t(mode.description) }}</p>
            <img v-if="form.chat_type === mode.value" class="selected-mark" src="@/assets/img/customer-service/basic-settings/dialogs/selected.png" alt="" />
          </button>
        </div>
        <div class="qa-settings" :class="{ 'is-disabled': !selectedMode.switchKey }">
          <a-checkbox :disabled="qaDisabled" :checked="form[qaMode.switchKey]" @change="event => setDirectReply(qaMode, event.target.checked)">{{ t('label_reply_directly_when_qa_matched') }}</a-checkbox>
          <div class="threshold">
            <span>{{ t('label_similarity_exceeds') }}</span>
            <a-input-number :value="form[qaMode.scoreKey]" :min="0" :max="1" :step="0.01" :precision="2" :controls="false" :disabled="qaDisabled || !form[qaMode.switchKey]" :status="isValidScore(form[qaMode.scoreKey]) ? undefined : 'error'" :aria-label="t('label_similarity_exceeds')" @change="value => setScore(qaMode, value)" @blur="restoreScore(qaMode)" />
            <span>{{ t('label_reply_directly_when_exceeds') }}</span>
          </div>
        </div>
      </div>
      <a-alert v-if="saveFailed" type="error" :message="commonT('save_failed')" show-icon>
        <template #action><a-button size="small" :loading="saving" @click="persist">{{ commonT('retry') }}</a-button></template>
      </a-alert>
    </SettingsDialog>
  </EditBox>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import EditBox from '@/views/robot/robot-config/basic-config/components/edit-box.vue'
import { saveConfigurationPatch } from '../configuration'
import SettingsDialog from './settings-dialog.vue'
import libraryImage from '@/assets/img/customer-service/basic-settings/dialogs/library.png'
import mixedImage from '@/assets/img/customer-service/basic-settings/dialogs/mixed.png'
import directImage from '@/assets/img/customer-service/basic-settings/dialogs/direct.png'

const { t } = useI18n('views.robot.robot-config.basic-config.components.chat-mode.chat-edit-alert')
const { t: commonT } = useI18n('views.customer-service.index')
const { t: summaryT } = useI18n('views.robot.robot-config.basic-config.components.chat-mode.index')
const summaryLabel = key => summaryT(key).replace(/[：:]\s*$/, '')
const context = inject('customerRobotContext')
const modes = [
  { value: '1', title: 'title_library_mode', summary: 'chat_mode_knowledge_base_only', description: 'desc_library_mode', switchKey: 'library_qa_direct_reply_switch', scoreKey: 'library_qa_direct_reply_score', image: libraryImage },
  { value: '3', title: 'title_mixture_mode', summary: 'chat_mode_mixed', description: 'desc_mixture_mode', switchKey: 'mixture_qa_direct_reply_switch', scoreKey: 'mixture_qa_direct_reply_score', image: mixedImage },
  { value: '2', title: 'title_direct_mode', summary: 'chat_mode_direct', description: 'desc_direct_mode', image: directImage }
]
const isValidScore = value => value !== null && value !== '' && Number.isFinite(Number(value)) && Number(value) >= 0 && Number(value) <= 1
const enabled = value => value === true || value === 'true' || value === 1 || value === '1'
function getConfig() {
  const info = context.robotInfo
  return {
    chat_type: String(info.chat_type || '1'),
    library_qa_direct_reply_switch: enabled(info.library_qa_direct_reply_switch),
    library_qa_direct_reply_score: isValidScore(info.library_qa_direct_reply_score) ? Number(info.library_qa_direct_reply_score) : 0.9,
    mixture_qa_direct_reply_switch: enabled(info.mixture_qa_direct_reply_switch),
    mixture_qa_direct_reply_score: isValidScore(info.mixture_qa_direct_reply_score) ? Number(info.mixture_qa_direct_reply_score) : 0.9
  }
}
const form = reactive(getConfig())
const saved = ref(getConfig())
const savedMode = computed(() => modes.find(mode => mode.value === saved.value.chat_type) || modes[0])
const selectedMode = computed(() => modes.find(mode => mode.value === form.chat_type) || modes[0])
const lastQaMode = ref('3')
const qaMode = computed(() => selectedMode.value.switchKey ? selectedMode.value : modes.find(mode => mode.value === lastQaMode.value))
const qaDisabled = computed(() => !selectedMode.value.switchKey || saving.value)
const modeButtons = ref([])
const editing = ref(false)
const saving = ref(false)
const saveFailed = ref(false)
let alive = true

function open() {
  Object.assign(form, saved.value)
  lastQaMode.value = savedMode.value.switchKey ? savedMode.value.value : '3'
  saveFailed.value = false
  editing.value = true
}

function selectMode(value) {
  if (saving.value || value === form.chat_type) return
  if (selectedMode.value.scoreKey) restoreScore(selectedMode.value)
  form.chat_type = value
  if (selectedMode.value.switchKey) lastQaMode.value = value
}
function navigateMode(event, index) {
  const nextIndex = { ArrowRight: (index + 1) % modes.length, ArrowDown: (index + 1) % modes.length, ArrowLeft: (index + modes.length - 1) % modes.length, ArrowUp: (index + modes.length - 1) % modes.length, Home: 0, End: modes.length - 1 }[event.key]
  if (nextIndex === undefined || saving.value) return
  event.preventDefault()
  selectMode(modes[nextIndex].value)
  modeButtons.value[nextIndex]?.focus()
}
function close() {
  if (!saving.value) editing.value = false
}
function setDirectReply(mode, checked) {
  if (qaDisabled.value || form[mode.switchKey] === checked) return
  form[mode.switchKey] = checked
}
function setScore(mode, value) {
  if (qaDisabled.value || !form[mode.switchKey] || form[mode.scoreKey] === value) return
  form[mode.scoreKey] = value
}
function restoreScore(mode) {
  if (!isValidScore(form[mode.scoreKey])) form[mode.scoreKey] = saved.value[mode.scoreKey]
}
async function persist() {
  if (saving.value || !modes.filter(mode => mode.scoreKey).every(mode => isValidScore(form[mode.scoreKey]))) return
  const patch = {
    chat_type: form.chat_type,
    library_qa_direct_reply_switch: String(form.library_qa_direct_reply_switch),
    library_qa_direct_reply_score: String(form.library_qa_direct_reply_score),
    mixture_qa_direct_reply_switch: String(form.mixture_qa_direct_reply_switch),
    mixture_qa_direct_reply_score: String(form.mixture_qa_direct_reply_score)
  }
  saving.value = true
  try {
    const success = await saveConfigurationPatch(context, patch)
    if (!alive) return
    saveFailed.value = !success
    if (success) {
      saved.value = getConfig()
      editing.value = false
    }
  } catch {
    if (alive) saveFailed.value = true
  } finally {
    if (alive) saving.value = false
  }
}
watch(context.savedRevision, () => { saved.value = getConfig() })
onBeforeUnmount(() => { alive = false })
</script>

<style lang="less" scoped>
.chat-mode-form { display: flex; flex-direction: column; gap: 24px; padding: 8px 19px 8px 8px; }
.mode-options { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 16px; }
.mode-option { position: relative; display: flex; flex-direction: column; align-items: stretch; min-width: 0; min-height: 170px; padding: 24px; border: 1px solid #d8d8d8; border-radius: 8px; background: #fff; overflow: hidden; cursor: pointer; font-family: inherit; text-align: left; }
.mode-option.active { border-color: #146ef5; background: #f1f6fe; }
.mode-option:focus-visible { outline: 2px solid #146ef5; outline-offset: 2px; }
.mode-illustration { position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: cover; pointer-events: none; }
.cropped .mode-illustration { height: 149.55%; object-fit: fill; }
.mode-title { position: relative; display: flex; gap: 8px; color: #080808; font-size: 16px; line-height: 24px; font-weight: 500; &::before { content: ''; width: 4px; height: 16px; flex-shrink: 0; margin-top: 4px; border-radius: 4px; background: #146ef5; } }
.description { position: relative; margin: 8px 0 0; color: #5a5a5a; font-size: 14px; line-height: 22px; }
.selected-mark { position: absolute; top: 0; right: 0; width: 32px; height: 32px; }
.qa-settings { display: flex; flex-direction: column; gap: 16px; padding: 16px; border: 1px solid #e8e8e8; border-radius: 8px; background: #f5f5f5; color: #222; font-size: 14px; line-height: 22px; }
.qa-settings :deep(.ant-checkbox-wrapper) { display: inline-flex; align-items: center; gap: 8px; margin: 0; color: inherit; line-height: 24px; }
.qa-settings :deep(.ant-checkbox) { top: 0; padding: 4px; }
.qa-settings :deep(.ant-checkbox-wrapper > span:last-child) { padding: 0; }
.qa-settings :deep(.ant-checkbox-inner) { width: 16px; height: 16px; border-radius: 4px; }
.qa-settings :deep(.ant-checkbox-checked .ant-checkbox-inner) { background: #146ef5; border-color: #146ef5; }
.qa-settings :deep(.ant-checkbox-checked .ant-checkbox-inner::after) { top: 1px; left: 1px; width: 12px; height: 12px; border: 0; transform: none; background: url('@/assets/img/customer-service/basic-settings/dialogs/checkbox-check.svg') center no-repeat; }
.qa-settings :deep(.ant-checkbox-disabled.ant-checkbox-checked .ant-checkbox-inner) { background: #f5f5f5; border-color: #e8e8e8; }
.qa-settings :deep(.ant-checkbox-disabled.ant-checkbox-checked .ant-checkbox-inner::after) { background-image: url('@/assets/img/customer-service/basic-settings/dialogs/checkbox-check-disabled.svg'); }
.qa-settings.is-disabled { color: #ababab; }
.qa-settings.is-disabled :deep(.ant-checkbox-wrapper-disabled) { color: #ababab; }
.threshold { display: flex; flex-wrap: wrap; align-items: center; gap: 12px; }
.threshold :deep(.ant-input-number) { width: 88px; height: 32px; border-color: #d8d8d8; border-radius: 6px; }
.threshold :deep(.ant-input-number-input) { height: 30px; padding: 0 12px; }
.threshold :deep(.ant-input-number-disabled) { background: #f5f5f5; border-color: #e8e8e8; color: #ababab; }
.summary { display: flex; flex-wrap: wrap; gap: 8px 16px; font-size: 14px; line-height: 22px; color: #595959; }
.chat-mode-form + :deep(.ant-alert) { margin-top: 12px; }
@media (max-width: 760px) { .chat-mode-form { padding: 8px; } .mode-options { grid-template-columns: 1fr; } .mode-option { min-height: 0; } }
</style>
