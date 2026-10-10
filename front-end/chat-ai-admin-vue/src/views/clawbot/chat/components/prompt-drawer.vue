<template>
  <a-drawer
    :open="props.open"
    placement="right"
    :width="500"
    :closable="false"
    :maskClosable="!props.loading"
    :bodyStyle="{ padding: '24px', background: '#fff' }"
    :headerStyle="{ padding: '10px 16px', borderBottom: '1px solid #f0f0f0' }"
    @close="handleClose"
  >
    <template #title>
      <div class="drawer-header">
        <span class="drawer-title">{{ t('title_persona') }}</span>
        <button class="drawer-close" type="button" :disabled="props.loading" @click="handleClose">
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
          <button v-else class="action-btn" type="button" @click="handleStartEdit">
            <EditOutlined />
            <span>{{ t('btn_edit') }}</span>
          </button>
        </div>
      </div>

      <div v-if="editing" class="prompt-insert-tip">{{ promptWordsT('msg_var_insert_tip_custom') }}</div>
      <div class="prompt-content-box" :class="{ editing, 'is-disabled': !editing || props.loading }">
        <AtInput
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
    <AddVariableModal ref="addVariableModalRef" :robotKey="props.robotKey" @ok="loadChatVariables" />
    <MiniCardFormModal ref="miniCardFormModalRef" @ok="() => miniCardStore.fetchCardList(true)" />
  </a-drawer>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { CloseOutlined, EditOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getChatVariables } from '@/api/robot/index'
import { useClawbotChatStore } from '@/stores/modules/clawbot-chat'
import { useMiniCardStore } from '@/stores/modules/mini-card'
import AtInput from '@/views/robot/robot-config/basic-config/components/at-input/index.vue'
import SlashPopupPanel from '@/views/robot/robot-config/basic-config/components/slash-popup-panel.vue'
import AddVariableModal from '@/views/robot/robot-config/basic-config/components/variable-setting/add-variable-modal.vue'
import MiniCardFormModal from '@/components/mini-card/mini-card-form-modal.vue'

const { t } = useI18n('views.clawbot.chat.components.prompt-drawer')
const { t: promptWordsT } = useI18n('views.robot.robot-config.basic-config.components.system-prompt-words')
const emit = defineEmits(['close', 'save'])

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
const chatStore = useClawbotChatStore()
const miniCardStore = useMiniCardStore()
const chatVariables = ref([])
let variableRequestVersion = 0

const systemVariables = [
  { id: 'system_current_time', label: '当前时间-系统变量', value: '【system_var:current_time】', typ: '系统变量' },
  { id: 'system_current_date', label: '当前日期-系统变量', value: '【system_var:current_date】', typ: '系统变量' }
]

const variableOptions = computed(() => [
  ...systemVariables,
  ...chatVariables.value.map((item) => ({
    ...item,
    label: `${item.variable_key}（${item.variable_name}）`,
    value: `【chat_variable:${item.variable_key}】`
  }))
])

const miniCardOptions = computed(() => miniCardStore.cardList.map((item) => ({
  ...item,
  label: `小程序卡片：${item.title}`,
  value: `【mini_card:${item.id}】`
})))

const loadChatVariables = async () => {
  const robotKey = props.robotKey
  const requestVersion = ++variableRequestVersion
  chatVariables.value = []
  if (!robotKey) return
  try {
    const res = await getChatVariables({ robot_key: robotKey })
    if (requestVersion === variableRequestVersion && robotKey === props.robotKey) {
      const list = res.data || []
      chatVariables.value = list
      if (chatStore.robot.robot_key === robotKey) {
        const fillVariables = chatStore.chat_variables.fill_variables || []
        chatStore.chat_variables.wait_variables = list
        chatStore.chat_variables.need_fill_variable = list.length > 0 && fillVariables.length === 0
      }
    }
  } catch {
    if (requestVersion === variableRequestVersion) chatVariables.value = []
  }
}

watch(
  () => [props.open, props.robotKey],
  ([open]) => {
    variableRequestVersion++
    chatVariables.value = []
    if (open) {
      editing.value = false
      draftPrompt.value = props.prompt || ''
      loadChatVariables()
      miniCardStore.fetchCardList()
    }
  }
)

watch([variableOptions, miniCardOptions], () => {
  nextTick(() => {
    if (editing.value) atInputRef.value?.updateValue()
    atInputRef.value?.refresh()
  })
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
  draftPrompt.value = props.prompt || ''
  editing.value = true
  nextTick(() => atInputRef.value?.refresh())
}

const handlePromptChange = (text) => {
  if (editing.value && !props.loading) draftPrompt.value = text
}

const handleAddVariable = () => {
  if (!props.robotKey) return
  if (chatVariables.value.length >= 10) {
    return message.error(promptWordsT('msg_max_variables'))
  }
  addVariableModalRef.value?.show()
}

const handleAddMiniCard = () => miniCardFormModalRef.value?.show()
const handleEditMiniCard = (card) => miniCardFormModalRef.value?.edit(card)

const handleCancelEdit = () => {
  if (props.loading) {
    return
  }
  editing.value = false
  draftPrompt.value = props.prompt || ''
}

const handleSave = () => {
  atInputRef.value?.updateValue()
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
  font-weight: 600;
  color: #262626;
}

.drawer-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  padding: 0;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: #595959;
  cursor: pointer;

  &:hover {
    background: #f5f5f5;
  }

  &:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }
}

.prompt-card {
  display: flex;
  flex-direction: column;
  height: 100%;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  background: #fff;
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
  color: #262626;
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
  color: #2475fc;
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;

  &:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }
}

.action-btn-cancel {
  color: #595959;
}

.prompt-content-box {
  flex: 1;
  margin: 0 16px 16px;
  position: relative;
  border-radius: 8px;
  background: #f2f4f7;
  overflow: hidden;
  &.editing {
    padding: 0;
  }
  &.is-disabled::before {
    content: '';
    position: absolute;
    inset: 0;
    z-index: 1;
    background: rgba(255, 255, 255, 0.25);
  }
  .prompt-at-input {
    height: 100%;
  }
  :deep(.mention-input-warpper) {
    height: 100%;
    border: none;
    padding: 12px 16px;
    background: #f2f4f7;
  }
  :deep(.j-mention.type-textarea) {
    min-height: 216px;
    height: 100%;
    overflow-y: auto;
    background: transparent;
    font-size: 16px;
    line-height: 24px;
    color: #262626;
  }
}

.prompt-insert-tip {
  margin: 0 16px 8px;
  color: #8c8c8c;
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
  color: #262626;
}

.prompt-textarea {
  width: 100%;
  height: 100% !important;
  border: none !important;
  box-shadow: none !important;
  resize: none !important;
  background: #f2f4f7 !important;

  :deep(textarea) {
    min-height: 216px !important;
    padding: 12px 16px !important;
    border: none !important;
    box-shadow: none !important;
    resize: none !important;
    background: #f2f4f7 !important;
    font-size: 16px;
    line-height: 24px;
    color: #262626;
  }
}
</style>
