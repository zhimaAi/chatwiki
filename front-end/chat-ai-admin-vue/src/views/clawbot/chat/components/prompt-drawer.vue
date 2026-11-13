<template>
  <a-drawer
    :open="props.open"
    placement="right"
    :width="drawerWidth"
    :closable="false"
    :maskClosable="!props.loading && !resizing"
    :contentWrapperStyle="{ maxWidth: '100vw', transition: resizing ? 'none' : undefined }"
    :drawerStyle="{ position: 'relative', cursor: resizing ? 'ew-resize' : undefined, userSelect: resizing ? 'none' : undefined }"
    :bodyStyle="{ display: 'flex', flexDirection: 'column', minHeight: 0, overflow: 'hidden', padding: '24px', background: '#fff' }"
    :headerStyle="{ padding: '10px 16px', borderBottom: '1px solid #f0f0f0' }"
    @close="handleClose"
  >
    <div
      class="drawer-resize-handle"
      :class="{ resizing }"
      role="separator"
      tabindex="0"
      aria-orientation="vertical"
      :aria-label="t('title_persona')"
      :aria-valuemin="Math.min(minDrawerWidth, viewportWidth)"
      :aria-valuemax="viewportWidth"
      :aria-valuenow="drawerWidth"
      @pointerdown="startResize"
      @pointermove="resizeDrawer"
      @pointerup="stopResize"
      @pointercancel="stopResize"
      @lostpointercapture="stopResize"
      @keydown="handleResizeKeydown"
    />
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
        <div class="prompt-insert-tip">{{ promptWordsT('msg_var_insert_tip_custom') }}</div>
        <div class="prompt-actions">
          <PromptOptimizeButton
            v-if="props.open && !props.loading"
            appearance="persona"
            :robot="props.robot"
            save-mode="draft"
            :onSave="applyOptimizedPromptDraft"
          />
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
          :readonly="!editing || props.loading"
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
import PromptOptimizeButton from '@/components/prompt-optimize/button.vue'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
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
  robot: { type: Object, default: () => ({}) },
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

const minDrawerWidth = 360
const viewportWidth = ref(window.innerWidth)
const drawerWidth = ref(Math.min(500, viewportWidth.value))
const resizing = ref(false)
let resizeSession = null

const setDrawerWidth = (width) => {
  const minWidth = Math.min(minDrawerWidth, viewportWidth.value)
  drawerWidth.value = Math.min(viewportWidth.value, Math.max(minWidth, width))
}

const stopResize = () => {
  const session = resizeSession
  resizeSession = null
  resizing.value = false
  if (session?.target.hasPointerCapture(session.pointerId)) {
    session.target.releasePointerCapture(session.pointerId)
  }
}

const startResize = (event) => {
  if (event.button !== 0 || !event.isPrimary) return
  event.preventDefault()
  resizeSession = {
    target: event.currentTarget,
    pointerId: event.pointerId,
    startX: event.clientX,
    startWidth: drawerWidth.value
  }
  resizing.value = true
  // 捕获指针，鼠标移出窄拖拽条后仍能继续调宽并收到松开事件。
  event.currentTarget.setPointerCapture(event.pointerId)
}

const resizeDrawer = (event) => {
  if (!resizeSession || event.pointerId !== resizeSession.pointerId) return
  // 抽屉右侧固定，左边缘向左移动时增加宽度。
  setDrawerWidth(resizeSession.startWidth + resizeSession.startX - event.clientX)
}

const handleResizeKeydown = (event) => {
  if (!['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(event.key)) return
  event.preventDefault()
  if (event.key === 'Home') setDrawerWidth(minDrawerWidth)
  else if (event.key === 'End') setDrawerWidth(viewportWidth.value)
  else setDrawerWidth(drawerWidth.value + (event.key === 'ArrowLeft' ? 20 : -20))
}

const handleWindowResize = () => {
  stopResize()
  viewportWidth.value = window.innerWidth
  setDrawerWidth(drawerWidth.value)
}

watch(() => props.open, (open) => {
  if (open) handleWindowResize()
  else stopResize()
})

onMounted(() => window.addEventListener('resize', handleWindowResize))
onBeforeUnmount(() => {
  stopResize()
  window.removeEventListener('resize', handleWindowResize)
})

const editing = ref(false)
const draftPrompt = ref('')
const draftPromptType = ref(undefined)
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
    draftPromptType.value = undefined
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
  draftPromptType.value = undefined
  draftPrompt.value = props.prompt || ''
  emit('close')
}

const handleStartEdit = () => {
  if (props.loading) return
  draftPromptType.value = undefined
  draftPrompt.value = props.prompt || ''
  editing.value = true
  nextTick(() => atInputRef.value?.refresh())
}

const applyOptimizedPromptDraft = (prompt) => {
  if (!props.open || props.loading) return false
  editing.value = true
  draftPrompt.value = prompt
  draftPromptType.value = 0
  nextTick(() => atInputRef.value?.refresh())
  return true
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
  draftPromptType.value = undefined
  draftPrompt.value = props.prompt || ''
}

const handleSave = () => {
  if (!editing.value || props.loading) return
  atInputRef.value?.updateValue()
  emit('save', draftPrompt.value, draftPromptType.value)
}
</script>

<style lang="less" scoped>
.drawer-resize-handle {
  position: absolute;
  inset: 0 auto 0 0;
  z-index: 2;
  width: 6px;
  cursor: ew-resize;
  touch-action: none;
  user-select: none;

  &:hover,
  &:focus-visible,
  &.resizing {
    background: rgba(36, 117, 252, 0.15);
  }
}

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
  flex: 1;
  min-height: 0;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  background: #fff;
}

.prompt-card-header {
  display: flex;
  flex-shrink: 0;
  align-items: flex-start;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
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
  margin-top: 8px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 0;
  border: 0;
  background: transparent;
  color: #666;
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
  min-height: 0;
  margin: 0;
  position: relative;
  border-radius: 8px;
  background: #fff;
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
    // 只读由输入组件控制，遮罩仅保留视觉效果，不拦截滚轮和滚动条。
    pointer-events: none;
  }
  .prompt-at-input {
    height: 100%;
  }
  :deep(.mention-input-warpper) {
    height: 100%;
    box-sizing: border-box;
    border: none;
    padding: 12px 16px;
    background: #fff;
    overflow-y: auto;
    overflow-x: hidden;
    overscroll-behavior-y: contain;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-thumb {
      border-radius: 3px;
      background: #bfbfbf;
    }
  }
  :deep(.j-mention.type-textarea) {
    min-height: 100%;
    height: auto;
    overflow: visible;
    overflow-wrap: anywhere;
    cursor: text;
    background: transparent;
    font-size: 16px;
    line-height: 24px;
    color: #262626;

  }
}

.prompt-insert-tip {
  flex-shrink: 0;
  margin: 0;
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
