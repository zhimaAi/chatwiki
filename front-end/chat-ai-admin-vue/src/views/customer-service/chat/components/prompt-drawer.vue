<template>
  <a-drawer
    :open="props.open"
    placement="right"
    :width="500"
    :closable="false"
    :maskClosable="!props.loading"
    :bodyStyle="{ padding: '24px' }"
    :headerStyle="{ padding: '10px 16px' }"
    @close="handleClose"
  >
    <template #title>
      <div class="drawer-header">
        <span class="drawer-title">{{ t('title_persona') }}</span>
        <button class="drawer-close" type="button" :aria-label="chatT('close')" :disabled="props.loading" @click="handleClose">
          <CloseOutlined />
        </button>
      </div>
    </template>

    <div class="prompt-card">
      <div class="prompt-card-header">
        <span class="prompt-label">{{ t('label_prompt') }}</span>
        <div class="prompt-actions">
          <template v-if="editing">
            <button class="action-btn action-btn-cancel" type="button" :disabled="props.loading" @click="handleCancelEdit">
              {{ t('btn_cancel') }}
            </button>
            <button class="action-btn" type="button" :disabled="props.loading" @click="handleSave">
              <a-spin v-if="props.loading" size="small" />
              <template v-else>{{ t('btn_save') }}</template>
            </button>
          </template>
          <button v-else class="action-btn" type="button" :disabled="props.loading" @click="handleStartEdit">
            <EditOutlined />
            <span>{{ t('btn_edit') }}</span>
          </button>
        </div>
      </div>

      <div v-if="editing" class="prompt-insert-tip">{{ promptWordsT('msg_var_insert_tip_custom') }}</div>
      <div class="prompt-content-box" :class="{ editing, 'is-disabled': !editing || props.loading }" :inert="!editing || props.loading ? true : undefined">
        <AtInput
          :key="editorVersion"
          ref="atInputRef"
          class="prompt-at-input"
          type="textarea"
          :defaultValue="draftPrompt"
          :placeholder="editing ? '' : t('msg_no_prompt')"
          :options="variableOptions"
          :cardOptions="miniCardOptions"
          :enableCardInsert="true"
          :checkAnyLevel="true"
          @change="handlePromptChange"
        >
          <template #dropdown="{ options, select, cardOptions, selectCard, close }">
            <SlashPopupPanel
              :options="options"
              :cardOptions="cardOptions"
              :select="select"
              :selectCard="selectCard"
              :close="close"
              @add-tag="handleAddVariable"
              @add-card="handleAddMiniCard"
              @edit-card="handleEditMiniCard"
            />
          </template>
        </AtInput>
      </div>
    </div>
    <template v-if="props.open && props.robotKey">
      <AddVariableModal :key="props.robotKey" ref="addVariableModalRef" :robotKey="props.robotKey" @ok="loadChatVariables" />
      <MiniCardFormModal :key="props.robotKey" ref="miniCardFormModalRef" @ok="loadMiniCards" />
    </template>
  </a-drawer>
</template>

<script setup>
import { computed, nextTick, onScopeDispose, ref, watch } from 'vue'
import { CloseOutlined, EditOutlined } from '@ant-design/icons-vue'
import { App } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getChatVariables } from '@/api/robot/index'
import { getMiniCardList } from '@/api/mini-card/index'
import AtInput from '@/views/robot/robot-config/basic-config/components/at-input/index.vue'
import SlashPopupPanel from '@/views/robot/robot-config/basic-config/components/slash-popup-panel.vue'
import AddVariableModal from '@/views/robot/robot-config/basic-config/components/variable-setting/add-variable-modal.vue'
import MiniCardFormModal from '@/components/mini-card/mini-card-form-modal.vue'

const { t } = useI18n('views.clawbot.chat.components.prompt-drawer')
const { t: promptWordsT } = useI18n('views.robot.robot-config.basic-config.components.system-prompt-words')
const { t: chatT } = useI18n('views.customer-service.chat')
const { message } = App.useApp()
// 只维护编辑副本和辅助选项；save交父页正式保存，变量定义用事件通知，不写聊天Store。
const emit = defineEmits(['close', 'save', 'variablesLoaded'])

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  prompt: {
    type: String,
    default: ''
  },
  robotKey: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const editing = ref(false)
const draftPrompt = ref('')
const atInputRef = ref(null)
const addVariableModalRef = ref(null)
const miniCardFormModalRef = ref(null)
const chatVariables = ref([])
const miniCards = ref([])
const editorVersion = ref(0)
let variableRequestVersion = 0
let cardRequestVersion = 0

const systemVariables = computed(() => [
  { id: 'system_current_time', label: chatT('system_current_time'), value: '【system_var:current_time】', typ: chatT('system_variable') },
  { id: 'system_current_date', label: chatT('system_current_date'), value: '【system_var:current_date】', typ: chatT('system_variable') }
])

const variableOptions = computed(() => [
  ...systemVariables.value,
  ...chatVariables.value.map((item) => ({
    ...item,
    label: `${item.variable_key}（${item.variable_name}）`,
    value: `【chat_variable:${item.variable_key}】`
  }))
])

const miniCardOptions = computed(() => miniCards.value.map((item) => ({
  ...item,
  label: chatT('mini_card_label', { title: item.title }),
  value: `【mini_card:${item.id}】`
})))

const loadChatVariables = async () => {
  const robotKey = props.robotKey
  const requestVersion = ++variableRequestVersion
  if (!props.open || !robotKey) return
  try {
    const res = await getChatVariables({ robot_key: robotKey })
    if (props.open && requestVersion === variableRequestVersion && robotKey === props.robotKey) {
      if (!Array.isArray(res?.data)) throw new Error('Invalid chat variable response')
      const list = res.data.map((item) => ({ ...item }))
      chatVariables.value = list
      emit('variablesLoaded', list.map((item) => ({ ...item })))
    }
  } catch {
    if (props.open && requestVersion === variableRequestVersion && robotKey === props.robotKey) {
      message.warning({ key: 'customer-service-prompt-options', content: chatT('error_prompt_options') })
    }
  }
}

const loadMiniCards = async () => {
  const robotKey = props.robotKey
  const requestVersion = ++cardRequestVersion
  if (!props.open || !robotKey) return
  try {
    const response = await getMiniCardList()
    if (!props.open || requestVersion !== cardRequestVersion || robotKey !== props.robotKey) return
    if (!Array.isArray(response?.data?.list)) throw new Error('Invalid mini card response')
    miniCards.value = response.data.list
  } catch {
    if (props.open && requestVersion === cardRequestVersion && robotKey === props.robotKey) {
      message.warning({ key: 'customer-service-prompt-options', content: chatT('error_prompt_options') })
    }
  }
}

watch(
  [() => props.open, () => props.robotKey],
  ([open]) => {
    // 关闭/身份变化立即失效旧查询；A→B→A也不能接纳第一次A的响应。
    variableRequestVersion++
    cardRequestVersion++
    editorVersion.value++
    chatVariables.value = []
    miniCards.value = []
    editing.value = false
    draftPrompt.value = props.prompt || ''
    if (open) {
      void loadChatVariables()
      void loadMiniCards()
    }
  },
  { immediate: true, flush: 'sync' }
)

onScopeDispose(() => {
  variableRequestVersion++
  cardRequestVersion++
})

watch([variableOptions, miniCardOptions], async () => {
  const version = editorVersion.value
  await nextTick()
  if (!props.open || version !== editorVersion.value) return
  if (editing.value) {
    // 选项回来时先同步尚未触发防抖的输入，再等defaultValue更新，避免刷新标签吞掉刚输入的字。
    atInputRef.value?.updateValue()
    await nextTick()
  }
  if (props.open && version === editorVersion.value) atInputRef.value?.refresh()
})

watch(
  () => props.prompt,
  (prompt) => {
    if (!editing.value) {
      draftPrompt.value = prompt || ''
    }
  }
)

const handleClose = () => {
  if (props.loading) {
    return
  }
  editing.value = false
  draftPrompt.value = props.prompt || ''
  emit('close')
}

const handleStartEdit = () => {
  if (props.loading) return
  draftPrompt.value = props.prompt || ''
  editing.value = true
  editorVersion.value++
  nextTick(() => atInputRef.value?.refresh())
}

const handlePromptChange = (text) => {
  if (editing.value && !props.loading) draftPrompt.value = text
}

const handleAddVariable = () => {
  if (!editing.value || props.loading || !props.robotKey) return
  if (chatVariables.value.length >= 10) {
    return message.error(promptWordsT('msg_max_variables'))
  }
  addVariableModalRef.value?.show()
}

const handleAddMiniCard = () => {
  if (editing.value && !props.loading) miniCardFormModalRef.value?.show()
}
const handleEditMiniCard = (card) => {
  if (editing.value && !props.loading) miniCardFormModalRef.value?.edit(card)
}

const handleCancelEdit = () => {
  if (props.loading) {
    return
  }
  editing.value = false
  draftPrompt.value = props.prompt || ''
  editorVersion.value++
}

const handleSave = () => {
  if (!editing.value || props.loading || !props.robotKey) return
  atInputRef.value?.updateValue()
  // 不提前退出编辑，父页保存失败时仍保留原草稿；成功后由父页关闭。
  emit('save', draftPrompt.value)
}
</script>

<style lang="less" scoped>
.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.drawer-title {
  font-size: 16px;
  line-height: 24px;
  font-weight: 500;
  color: var(--app-text-primary);
}

.drawer-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  padding: 0;
  border: 0;
  border-radius: var(--app-radius-control);
  background: transparent;
  color: var(--app-text-secondary);
  cursor: pointer;

  &:hover:not(:disabled) {
    background: var(--app-bg-hover);
  }

  &:active:not(:disabled) {
    background: var(--app-fill-secondary);
  }

  &:disabled {
    cursor: not-allowed;
    color: var(--app-text-disabled);
  }
}

.prompt-card {
  display: flex;
  flex-direction: column;
  height: 100%;
  border: 1px solid var(--app-border);
  border-radius: var(--app-radius-lg);
  box-shadow: var(--app-shadow-card);
  overflow: hidden;
  background: var(--app-bg-container);
}

.prompt-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
}

.prompt-label {
  font-size: 16px;
  line-height: 24px;
  color: var(--app-text-primary);
}

.prompt-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 0;
  border: 0;
  background: transparent;
  color: var(--app-primary);
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;

  &:hover:not(:disabled) {
    color: var(--app-primary-hover);
  }

  &:active:not(:disabled) {
    color: var(--app-primary-active);
  }

  &:disabled {
    cursor: not-allowed;
    color: var(--app-text-disabled);
  }
}

.action-btn-cancel {
  color: var(--app-text-secondary);
}

.drawer-close,
.action-btn {
  &:focus-visible {
    outline: 2px solid var(--app-primary);
    outline-offset: 2px;
  }
}

.prompt-content-box {
  flex: 1;
  margin: 0 16px 16px;
  position: relative;
  border-radius: var(--app-radius-lg);
  background: var(--app-bg-layout);
  overflow: hidden;
  &.editing {
    padding: 0;
    &:focus-within {
      box-shadow: var(--app-shadow-focus);
    }
  }
  &.is-disabled::before {
    content: '';
    position: absolute;
    inset: 0;
    z-index: 1;
    background: var(--app-bg-container);
    opacity: 0.25;
  }
  .prompt-at-input {
    height: 100%;
  }
  :deep(.mention-input-warpper) {
    height: 100%;
    border: none;
    padding: 12px 16px;
    background: var(--app-bg-layout);
  }
  :deep(.j-mention.type-textarea) {
    min-height: 216px;
    height: 100%;
    overflow-y: auto;
    background: transparent;
    font-size: 16px;
    line-height: 24px;
    color: var(--app-text-primary);
  }
}

.prompt-insert-tip {
  margin: 0 16px 8px;
  color: var(--app-text-tertiary);
  font-size: 12px;
}

.prompt-content {
  height: 100% !important;
  padding: 12px 16px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 16px;
  line-height: 24px;
  color: var(--app-text-primary);
}

.prompt-textarea {
  width: 100%;
  height: 100% !important;
  border: none !important;
  box-shadow: none !important;
  resize: none !important;
  background: var(--app-bg-layout) !important;

  :deep(textarea) {
    min-height: 216px !important;
    padding: 12px 16px !important;
    border: none !important;
    box-shadow: none !important;
    resize: none !important;
    background: var(--app-bg-layout) !important;
    font-size: 16px;
    line-height: 24px;
    color: var(--app-text-primary);
  }
}
</style>
