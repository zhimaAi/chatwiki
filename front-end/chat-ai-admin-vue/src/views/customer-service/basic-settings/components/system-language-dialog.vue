<template>
  <SettingsDialog v-if="editing" :open="true" :width="800" :ok-text="commonT('save')" :cancel-text="commonT('cancel')" :confirm-loading="saving" :closable="!saving" :mask-closable="!saving" :keyboard="!saving" :cancel-button-props="{ disabled: saving }" @ok="persist" @cancel="close">
    <template #title>
      <div class="system-language-header">
        <span>系统提示语</span>
        <a-select class="current-language" :value="languageKey" :options="languageOptions" :bordered="false" :disabled="saving" aria-label="当前编辑语言" @update:value="changeLanguage">
          <template #suffixIcon><img src="@/assets/img/customer-service/basic-settings/system-language/chevron-down.svg" alt="" /></template>
        </a-select>
        <a-popover v-model:open="moreLanguageOpen" trigger="click" placement="bottomRight" overlay-class-name="workspace-system-language-popover" :overlay-inner-style="{ padding: 16 }">
          <template #content>
            <div class="language-options-header"><strong>语言</strong><span>隐藏/显示</span></div>
            <div class="language-options-list">
              <div v-for="language in moreLanguages" :key="language.value" class="more-language-item">
                <span class="language-name">{{ language.label }}</span>
                <a-switch :checked="enabledLanguages.includes(language.value)" :disabled="saving" :aria-label="`${language.label}隐藏/显示`" @change="checked => emit('visibility-change', language.value, checked)" />
              </div>
            </div>
          </template>
          <a-button class="more-language-button" :disabled="saving"><img src="@/assets/img/customer-service/basic-settings/system-language/language.svg" alt="" />更多语言</a-button>
        </a-popover>
      </div>
    </template>
    <template #closeIcon><img src="@/assets/img/customer-service/basic-settings/prompt/close.svg" alt="" /></template>
    <div class="system-language-form" :inert="saving ? '' : null">
      <fieldset class="language-sections" :disabled="saving">
        <section v-for="section in sections" :key="section.key" class="language-section" :data-section="section.key">
          <div class="prompt-field">
            <label :for="`${section.key}-${languageKey}`">{{ section.title }}</label>
            <a-textarea :id="`${section.key}-${languageKey}`" v-model:value="form[section.key].content" :placeholder="section.placeholder" />
          </div>
          <draggable v-model="form[section.key].question" item-key="id" handle=".question-drag" class="guide-questions" :animation="150" :disabled="saving">
            <template #item="{ element, index }">
              <div class="guide-question">
                <button type="button" class="question-drag icon-button" :aria-label="`调整第${index + 1}个引导问题的顺序`" title="按上下方向键调整顺序" @keydown.up.prevent="moveQuestion(section.key, index, -1, $event)" @keydown.down.prevent="moveQuestion(section.key, index, 1, $event)"><img src="@/assets/img/customer-service/basic-settings/system-language/drag.svg" alt="" /></button>
                <a-input v-model:value="element.content" :placeholder="welcomeT('ph_input_question')" :aria-label="`${section.title}第${index + 1}个引导问题`" />
                <button type="button" class="icon-button" :aria-label="`删除第${index + 1}个引导问题`" @click="form[section.key].question.splice(index, 1)"><img src="@/assets/img/customer-service/basic-settings/dialogs/close.svg" alt="" /></button>
              </div>
            </template>
          </draggable>
          <a-button class="add-question" type="dashed" @click="addQuestion(section.key)"><img src="@/assets/img/customer-service/basic-settings/prompt/add-theme.svg" alt="" />{{ welcomeT('btn_add_guide_question') }}</a-button>
        </section>
        <section class="language-section answer-section">
          <div class="answer-label"><label :for="`answer-prompt-${languageKey}`">{{ answerT('title_answer_generation_prompt') }}</label><a-switch v-model:checked="form.tips_before_answer_switch" checked-value="true" un-checked-value="false" :aria-label="answerT('title_answer_generation_prompt')" /></div>
          <a-input :id="`answer-prompt-${languageKey}`" v-model:value="form.tips_before_answer_content" :max-length="10" :placeholder="answerT('ph_input')" @blur="restoreAnswerPrompt" />
        </section>
      </fieldset>
      <a-alert v-if="saveFailed" type="error" :message="commonT('save_failed')" show-icon />
    </div>
  </SettingsDialog>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import draggable from 'vuedraggable'
import { useI18n } from '@/hooks/web/useI18n'
import { getUuid } from '@/utils'
import { cloneConfig } from '../../config-session'
import SettingsDialog from './settings-dialog.vue'

const props = defineProps({
  configs: { type: Array, default: () => [] },
  languages: { type: Array, default: () => [] },
  moreLanguages: { type: Array, default: () => [] },
  enabledLanguages: { type: Array, default: () => [] },
  isWorkflow: Boolean,
  submit: { type: Function, required: true }
})
const emit = defineEmits(['visibility-change'])
const { t: commonT } = useI18n('views.customer-service.index')
const { t: welcomeT } = useI18n('views.robot.robot-config.basic-config.components.welcome-words')
const { t: unknownT } = useI18n('views.robot.robot-config.basic-config.components.unknown-problem-prompt')
const { t: answerT } = useI18n('views.robot.robot-config.basic-config.components.answer-prompt')
const editing = ref(false)
const saving = ref(false)
const saveFailed = ref(false)
const moreLanguageOpen = ref(false)
const languageKey = ref('zh-CN')
const draft = ref([])
const form = ref(null)
const languageOptions = computed(() => props.languages.map(item => ({ ...item, label: item.value === 'zh-CN' ? '简体中文' : item.label })))
const sections = computed(() => [
  { key: 'welcomes', title: welcomeT('title_welcome'), placeholder: welcomeT('ph_input_welcome') },
  ...(!props.isWorkflow ? [{ key: 'unknown_question_prompt', title: unknownT('title_unknown_problem_prompt'), placeholder: unknownT('ph_input_unknown_prompt') }] : [])
])

function readMenu(value) {
  const menu = value ? JSON.parse(value) : {}
  return { ...menu, content: menu.content || '', question: (menu.question || []).map(content => ({ id: getUuid(8), content })) }
}
function readForm() {
  const config = draft.value.find(item => item.lang_key === languageKey.value) || { lang_key: languageKey.value }
  form.value = {
    ...config,
    welcomes: readMenu(config.welcomes),
    unknown_question_prompt: readMenu(config.unknown_question_prompt),
    tips_before_answer_switch: String(config.tips_before_answer_switch) === 'true' ? 'true' : 'false',
    tips_before_answer_content: config.tips_before_answer_content ?? answerT('msg_thinking_please_wait')
  }
}
function storeForm() {
  const config = { ...form.value }
  for (const { key } of sections.value) {
    config[key] = JSON.stringify({ ...config[key], question: config[key].question.map(item => item.content) })
  }
  // 智能体不展示未知问题设置，保存时保留接口中的原值。
  if (props.isWorkflow) config.unknown_question_prompt = draft.value.find(item => item.lang_key === config.lang_key)?.unknown_question_prompt
  const index = draft.value.findIndex(item => item.lang_key === config.lang_key)
  if (index < 0) draft.value.push(config)
  else draft.value.splice(index, 1, config)
}
function open(key = 'zh-CN') {
  if (saving.value) return
  draft.value = cloneConfig(props.configs)
  languageKey.value = props.languages.some(item => item.value === key) ? key : 'zh-CN'
  readForm()
  saveFailed.value = false
  moreLanguageOpen.value = false
  editing.value = true
}
function close() {
  if (saving.value) return
  moreLanguageOpen.value = false
  editing.value = false
}
function changeLanguage(key) {
  if (saving.value || key === languageKey.value) return
  storeForm()
  languageKey.value = key
  readForm()
}
watch(() => props.languages, languages => {
  if (editing.value && !languages.some(item => item.value === languageKey.value)) changeLanguage('zh-CN')
})
async function addQuestion(key) {
  if (saving.value) return
  form.value[key].question.push({ id: getUuid(8), content: '' })
  await nextTick()
  document.getElementById(`${key}-${languageKey.value}`)?.closest('section')?.querySelector('.guide-question:last-child input')?.focus()
}
async function moveQuestion(key, index, step, event) {
  if (saving.value) return
  const items = form.value[key].question
  const next = index + step
  if (next < 0 || next >= items.length) return
  const button = event.currentTarget
  items.splice(next, 0, items.splice(index, 1)[0])
  await nextTick()
  button.focus()
}
function restoreAnswerPrompt() {
  if (!form.value.tips_before_answer_content.trim()) form.value.tips_before_answer_content = answerT('msg_thinking_please_wait')
}
async function persist() {
  if (saving.value) return
  restoreAnswerPrompt()
  storeForm()
  saving.value = true
  moreLanguageOpen.value = false
  saveFailed.value = false
  try {
    const success = await props.submit(cloneConfig(draft.value))
    saveFailed.value = !success
    if (success) editing.value = false
  } catch {
    saveFailed.value = true
  } finally {
    saving.value = false
  }
}
defineExpose({ open })
</script>

<style lang="less" scoped>
.system-language-header { display: flex; align-items: center; gap: 8px; min-width: 0; min-height: 32px; }
.current-language { min-width: 94px; color: #5a5a5a; font-size: 14px; font-weight: 400; }
.current-language :deep(.ant-select-selector) { padding: 0 8px !important; }
.current-language :deep(.ant-select-arrow) { right: 8px; width: 16px; height: 16px; margin-top: -8px; }
.more-language-button { display: inline-flex; align-items: center; gap: 6px; height: 32px; margin-left: auto; padding: 0 8px; border: 0; border-radius: 6px; background: #f0f0f0; color: #5a5a5a; box-shadow: none; }
.system-language-form { padding: 8px; color: #222; font-size: 14px; line-height: 22px; }
.language-sections { display: flex; flex-direction: column; gap: 16px; min-width: 0; margin: 0; padding: 8px; border: 0; border-radius: 8px; background: #f5f5f5; }
.language-section { display: flex; flex-direction: column; gap: 12px; min-width: 0; padding: 16px; border-radius: 6px; background: #fff; }
.prompt-field { display: flex; flex-direction: column; gap: 8px; }
.language-section :deep(.ant-input) { min-width: 0; height: 32px; padding: 4px 12px; border-radius: 6px; border-color: #d8d8d8; font-size: 14px; line-height: 22px; }
.language-section :deep(textarea.ant-input) { height: 84px; min-height: 84px; padding: 8px 12px; resize: vertical; }
.guide-questions { display: flex; flex-direction: column; gap: 8px; }
.guide-questions:empty { display: none; }
.guide-question { display: grid; grid-template-columns: 32px minmax(0, 1fr) 32px; align-items: center; gap: 8px; }
.icon-button { display: flex; align-items: center; justify-content: center; width: 24px; height: 24px; padding: 0; border: 0; border-radius: 6px; background: transparent; cursor: pointer; }
.question-drag { cursor: grab; touch-action: none; }
.icon-button:focus-visible { outline: 2px solid #146ef5; outline-offset: 2px; }
.add-question { display: flex; align-items: center; justify-content: center; gap: 8px; height: 32px; border-radius: 6px; border-color: #e8e8e8; color: #5a5a5a; }
.answer-label { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.answer-label :deep(.ant-switch) { min-width: 36px; height: 20px; background: #d8d8d8; }
.answer-label :deep(.ant-switch-checked) { background: #146ef5; }
.answer-label :deep(.ant-switch-handle) { width: 16px; height: 16px; top: 2px; }
.answer-label :deep(.ant-switch-checked .ant-switch-handle) { inset-inline-start: calc(100% - 18px); }
.system-language-form > :deep(.ant-alert) { margin-top: 12px; }
@media (max-width: 560px) {
  .system-language-header { flex-wrap: wrap; }
  .more-language-button { margin-left: 0; }
  .system-language-form { padding: 0; }
  .language-section { padding: 12px; }
  .guide-question { grid-template-columns: 24px minmax(0, 1fr) 24px; gap: 4px; }
}
</style>

<style lang="less">
.workspace-system-language-popover {
  max-width: calc(100vw - 32px);
  .ant-popover-inner { width: 300px; max-width: calc(100vw - 32px); border-radius: 12px; color: #222; font-size: 14px; line-height: 22px; }
  .language-options-header { display: flex; align-items: center; justify-content: space-between; span { color: #5a5a5a; } }
  .language-options-list { display: flex; flex-direction: column; gap: 12px; max-height: 226px; margin-top: 16px; overflow-y: auto; overscroll-behavior: contain; }
  .more-language-item { display: flex; align-items: center; justify-content: space-between; gap: 16px; min-height: 22px; }
  .ant-switch { min-width: 36px; height: 20px; flex-shrink: 0; background: #d8d8d8; }
  .ant-switch-checked { background: #146ef5; }
  .ant-switch-handle { width: 16px; height: 16px; top: 2px; }
  .ant-switch-checked .ant-switch-handle { inset-inline-start: calc(100% - 18px); }
}
</style>
