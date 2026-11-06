<template>
  <a-popover :open="open" :trigger="['click']" placement="bottomLeft" :arrow="false"
    :align="{ offset: [0, 8] }" :overlay-inner-style="popoverInnerStyle" @open-change="emit('update:open', $event)">
    <span ref="triggerHost" class="picker-trigger"><slot /></span>
    <template #content>
      <a-popover :open="open && !!previewOption && !interactionDisabled" :trigger="[]" placement="rightTop" :arrow="false"
        :align="{ offset: [8, 0] }" :overlay-inner-style="popoverInnerStyle" :get-popup-container="getPickerContainer">
        <div id="customer-service-app-picker" ref="pickerPanel" class="robot-selector" role="dialog" :aria-label="t('select_app')"
          :aria-busy="status === 'loading'" @mouseenter="cancelPreviewClose" @mouseleave="schedulePreviewClose" @focusout="handleFocusOut">
          <div v-if="status === 'idle' || status === 'loading'" class="picker-feedback" role="status">
            <a-spin /><span>{{ t('loading_robots') }}</span>
          </div>
          <div v-else-if="status === 'error'" class="picker-feedback">
            <a-alert type="error" show-icon :message="t('error_robots')" />
            <a-button @click="emit('retry')">{{ t('retry') }}</a-button>
          </div>
          <div v-else-if="options.length === 0" class="picker-feedback">
            <a-empty :description="t('empty_robot')" /><a-button @click="emit('retry')">{{ t('retry') }}</a-button>
          </div>
          <div v-else class="robot-list-scroll" @scroll="clearPreview">
            <draggable v-model="robotDraft" item-key="id" handle=".drag-btn" :disabled="busy || creating || status !== 'ready'"
              class="robot-list" @start="startDrag" @end="saveSort">
              <template #item="{ element: option }">
                <article class="robot-row" :class="{ selected: isSelected(option), previewed: previewId === String(option.id), unavailable: !canSelect(option) }"
                  @mouseenter="showPreview(option)" @focusin="showPreview(option)">
                  <button type="button" class="drag-btn" :disabled="interactionDisabled" :aria-label="t('drag_sort') + '：' + option.name"
                    :title="t('drag_sort')" @click.stop><img :src="icons.drag" alt="" draggable="false" /></button>
                  <button type="button" class="robot-row-select" :disabled="!canSelect(option) || interactionDisabled"
                    :aria-pressed="isSelected(option)" :aria-label="option.name" :title="canSelect(option) ? option.name : t('invalid_robot')"
                    @click="selectRobot(option)">
                    <img class="robot-avatar" :src="option.avatar || defaultAvatar(option)" alt="" draggable="false" @error="useDefaultAvatar($event, option)" />
                    <span class="robot-name">{{ option.name }}</span>
                  </button>
                </article>
              </template>
            </draggable>
          </div>
          <template v-if="canCreateRobot || canCreateAgent">
            <div class="picker-divider" role="separator"></div>
            <a-popover :open="open && createMenuOpen && !interactionDisabled" :trigger="['hover', 'focus', 'click']" placement="rightTop" :arrow="false"
              :align="{ offset: [8, 0] }" :overlay-inner-style="popoverInnerStyle" :get-popup-container="getPickerContainer" @open-change="setCreateMenuOpen">
              <button type="button" class="picker-add" :class="{ active: createMenuOpen }" :disabled="interactionDisabled"
                :aria-expanded="createMenuOpen" aria-haspopup="menu" @mouseenter="clearPreview" @focus="clearPreview">
                <img :src="icons.add" alt="" /><span>{{ t('add_app') }}</span><img :src="icons.arrowRight" alt="" />
              </button>
              <template #content>
                <div class="picker-create-menu" role="menu">
                  <button v-if="canCreateRobot" type="button" role="menuitem" :disabled="interactionDisabled" @click="createRobot">
                    <img :src="icons.robotSmall" alt="" /><span>{{ t('add_robot') }}</span>
                  </button>
                  <button v-if="canCreateAgent" type="button" role="menuitem" :disabled="interactionDisabled" @click="createAgent">
                    <img :src="icons.agentSmall" alt="" /><span>{{ t('add_agent') }}</span>
                  </button>
                </div>
              </template>
            </a-popover>
          </template>
        </div>
        <template #content>
          <section v-if="previewOption" class="robot-preview" :aria-label="previewOption.name" @mouseenter="cancelPreviewClose" @mouseleave="schedulePreviewClose">
            <div class="preview-header">
              <img class="preview-avatar" :src="previewOption.avatar || defaultAvatar(previewOption, true)" alt="" @error="useDefaultAvatar($event, previewOption, true)" />
              <span class="robot-type" :class="{ agent: Number(previewOption.applicationType) === 2 }">{{ t(Number(previewOption.applicationType) === 2 ? 'type_agent' : 'type_chat_robot') }}</span>
            </div>
            <h3 :title="previewOption.name">{{ previewOption.name }}</h3>
            <p class="robot-description" :title="previewOption.intro">{{ previewOption.intro || '—' }}</p>
            <div class="robot-actions">
              <button type="button" :disabled="interactionDisabled" @click="handleAction('copy', previewOption)"><img :src="icons.copy" alt="" />{{ t('copy') }}</button>
              <template v-if="getRobotPermission(previewOption.id) === 4">
                <span class="action-divider" role="separator"></span>
                <button type="button" :disabled="interactionDisabled" @click="handleAction('delete', previewOption)"><img :src="icons.delete" alt="" />{{ t('delete') }}</button>
              </template>
            </div>
          </section>
        </template>
      </a-popover>
    </template>
  </a-popover>
  <AddRobotAlert ref="addRobotRef" hide-group create-redirect="none" @ok="emit('created', $event)" />
  <AddAgentAlert ref="addAgentRef" @ok="emit('created', $event)" />
</template>

<script setup>
import { computed, createVNode, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import draggable from 'vuedraggable'
import { Modal, message } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { deleteRobot, moveRobotSort, robotCopy } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'
import { checkRole, getRobotPermission } from '@/utils/permission'
import AddRobotAlert from '@/views/robot/robot-list/components/add-robot-alert.vue'
import AddAgentAlert from '@/views/clawbot/components/add-agent-alert.vue'
import dragIcon from '@/assets/img/customer-service/app-picker/drag.svg'
import robotSmallIcon from '@/assets/img/customer-service/app-picker/robot-small.svg'
import agentSmallIcon from '@/assets/img/customer-service/app-picker/agent-small.svg'
import addIcon from '@/assets/img/customer-service/app-picker/add.svg'
import arrowRightIcon from '@/assets/img/customer-service/app-picker/arrow-right.svg'
import robotIcon from '@/assets/img/customer-service/app-picker/robot.svg'
import agentIcon from '@/assets/img/customer-service/app-picker/agent.svg'
import copyIcon from '@/assets/img/customer-service/app-picker/copy.svg'
import deleteIcon from '@/assets/img/customer-service/app-picker/delete.svg'

const icons = { drag: dragIcon, robotSmall: robotSmallIcon, agentSmall: agentSmallIcon, add: addIcon, arrowRight: arrowRightIcon, robot: robotIcon, agent: agentIcon, copy: copyIcon, delete: deleteIcon }
const props = defineProps({
  open: { type: Boolean, default: false },
  value: { type: String, default: '' },
  options: { type: Array, default: () => [] },
  status: { type: String, default: 'idle' }
})
const emit = defineEmits(['update:open', 'change', 'retry', 'removed', 'created'])
const { t } = useI18n('views.customer-service.index')
const { t: listT } = useI18n('views.robot.robot-list.robot-list')
const addRobotRef = ref(null)
const addAgentRef = ref(null)
const pickerPanel = ref(null)
const triggerHost = ref(null)
const creating = ref(false)
const busy = ref(false)
const dragging = ref(false)
const robotDraft = ref([])
const previewId = ref('')
const createMenuOpen = ref(false)
const interactionDisabled = computed(() => busy.value || creating.value || dragging.value)
const previewOption = computed(() => previewId.value ? robotDraft.value.find(option => String(option.id) === previewId.value) : null)
const canCreateRobot = computed(() => checkRole(['RobotCreate']))
const canCreateAgent = computed(() => checkRole(['ClawBotCreate']))
const popoverInnerStyle = { padding: 0, background: 'transparent', boxShadow: 'none' }
let dragSnapshot = []
let confirmation = null
let confirmationTimer = null
let previewCloseTimer = null
let alive = true

// Store 是列表真值，本地数组只承接拖拽；刷新后重新同步接口顺序。
watch(() => props.options, options => { robotDraft.value = options.slice() }, { immediate: true })
const cancelPreviewClose = () => {
  if (previewCloseTimer !== null) clearTimeout(previewCloseTimer)
  previewCloseTimer = null
}
const clearPreview = () => { cancelPreviewClose(); previewId.value = '' }
const schedulePreviewClose = () => {
  cancelPreviewClose()
  previewCloseTimer = setTimeout(clearPreview, 180)
}
const showPreview = option => {
  if (!props.open || interactionDisabled.value) return
  cancelPreviewClose()
  createMenuOpen.value = false
  previewId.value = String(option.id)
}
const handleFocusOut = event => {
  if (!pickerPanel.value?.contains(event.relatedTarget)) schedulePreviewClose()
}
const getPickerContainer = () => pickerPanel.value || document.body
const setCreateMenuOpen = open => {
  createMenuOpen.value = props.open && !interactionDisabled.value && open
  if (createMenuOpen.value) clearPreview()
}
const close = async (restoreFocus = false) => {
  emit('update:open', false)
  if (restoreFocus) {
    await nextTick()
    if (alive) triggerHost.value?.querySelector('button')?.focus({ preventScroll: true })
  }
}
watch(() => props.open, (open, _, onCleanup) => {
  clearPreview()
  createMenuOpen.value = false
  if (!open) {
    if (dragging.value) robotDraft.value = dragSnapshot.slice()
    dragging.value = false
    return
  }
  const escape = event => {
    if (event.key === 'Escape' && !confirmation) { event.preventDefault(); close(true) }
  }
  document.addEventListener('keydown', escape)
  onCleanup(() => document.removeEventListener('keydown', escape))
}, { immediate: true })
watch(interactionDisabled, disabled => {
  if (disabled) { clearPreview(); createMenuOpen.value = false }
})
watch(() => props.status, status => { if (status !== 'ready') clearPreview() })
const isSelected = option => Boolean(props.value) && props.value === String(option.id)
const canSelect = option => props.status === 'ready' && Boolean(option.id && option.robotKey)
const selectRobot = option => {
  if (!canSelect(option) || interactionDisabled.value) return
  emit('change', option.id)
  close()
}
const startDrag = () => {
  dragSnapshot = robotDraft.value.slice()
  dragging.value = true
  clearPreview()
}
const saveSort = async ({ oldIndex, newIndex }) => {
  if (!dragging.value) return
  dragging.value = false
  if (oldIndex === newIndex) return
  const previousOrder = dragSnapshot
  const movedRobot = robotDraft.value[newIndex]
  if (!movedRobot) { robotDraft.value = previousOrder; return }
  busy.value = true
  try {
    const response = await moveRobotSort({ move_robot_id: movedRobot.id, to_robot_id: newIndex > 0 ? robotDraft.value[newIndex - 1].id : 0 })
    if (!alive) return
    if (!response) { robotDraft.value = previousOrder; return }
    message.success(t('sort_saved'))
    emit('retry')
  } catch {
    // 请求封装负责提示；失败只恢复顺序，不改变当前应用。
    if (alive) robotDraft.value = previousOrder
  } finally { if (alive) busy.value = false }
}
const defaultAvatar = (option, large = false) => Number(option.applicationType) === 2 ? (large ? icons.agent : icons.agentSmall) : (large ? icons.robot : icons.robotSmall)
const useDefaultAvatar = (event, option, large = false) => {
  const fallback = defaultAvatar(option, large)
  if (event.target.getAttribute('src') !== fallback) event.target.src = fallback
}
const createRobot = async () => {
  if (!canCreateRobot.value || interactionDisabled.value) return
  creating.value = true
  close()
  try {
    addRobotRef.value.setGroupId('0')
    await addRobotRef.value.open(0)
  } catch {
    // 创建表单加载失败仍由统一请求封装提示。
  } finally { if (alive) creating.value = false }
}
const createAgent = async () => {
  if (!canCreateAgent.value || interactionDisabled.value) return
  creating.value = true
  close()
  try { await addAgentRef.value.open() }
  catch {
    // 创建表单加载失败仍由统一请求封装提示。
  }
  finally { if (alive) creating.value = false }
}
const clearConfirmationTimer = () => {
  if (confirmationTimer !== null) clearInterval(confirmationTimer)
  confirmationTimer = null
}
const confirmAction = (action, option) => {
  if (confirmation || interactionDisabled.value || !option.id) return
  if (action === 'delete' && getRobotPermission(option.id) !== 4) {
    message.error(listT('msg_no_delete_permission'))
    return
  }
  let secondsRemaining = 3
  confirmation = Modal.confirm({
    title: listT(action === 'delete' ? 'title_delete_robot' : 'title_copy_robot') + ' ' + option.name,
    content: listT(action === 'delete' ? 'msg_confirm_delete_robot' : 'msg_confirm_copy_robot'),
    icon: createVNode(ExclamationCircleOutlined),
    okText: secondsRemaining + ' ' + listT('btn_confirm'),
    okType: 'danger',
    cancelText: listT('btn_cancel'),
    okButtonProps: { disabled: true },
    async onOk() {
      if (!alive) return
      clearConfirmationTimer()
      busy.value = true
      try {
        const response = await (action === 'delete' ? deleteRobot : robotCopy)({ id: option.id })
        if (!alive || !response) return
        message.success(listT(action === 'delete' ? 'msg_delete_success' : 'msg_copy_success'))
        if (action === 'delete') emit('removed', option.id)
        else emit('retry')
      } finally {
        // 请求拒绝交给确认框处理，失败时不刷新或切换应用。
        if (alive) busy.value = false
      }
    },
    onCancel: clearConfirmationTimer,
    afterClose() { clearConfirmationTimer(); confirmation = null }
  })
  confirmationTimer = setInterval(() => {
    secondsRemaining -= 1
    confirmation?.update({
      okText: secondsRemaining > 0 ? secondsRemaining + ' ' + listT('btn_confirm') : listT('btn_confirm'),
      okButtonProps: { disabled: secondsRemaining > 0 }
    })
    if (secondsRemaining === 0) clearConfirmationTimer()
  }, 1000)
}
const handleAction = (action, option) => {
  if (!props.open || interactionDisabled.value || !['copy', 'delete'].includes(action)) return
  clearPreview()
  createMenuOpen.value = false
  confirmAction(action, option)
}
onBeforeUnmount(() => {
  alive = false
  cancelPreviewClose()
  clearConfirmationTimer()
  confirmation?.destroy()
})
</script>

<style lang="less" scoped>
.picker-trigger { display: block; width: 100%; }
.robot-selector, .robot-preview, .picker-create-menu {
  box-sizing: border-box; border: 1px solid #e8e8e8; border-radius: 8px; color: #222; background: #fff;
  box-shadow: 0 8px 24px 8px rgba(8, 8, 8, .04), 0 6px 12px rgba(8, 8, 8, .04), 0 4px 8px -8px rgba(8, 8, 8, .06);
}
.robot-selector { position: relative; width: 320px; max-width: calc(100vw - 24px); padding: 4px; }
.robot-list-scroll { max-height: max(40px, min(440px, calc(100vh - 160px))); overflow-y: auto; scrollbar-width: thin; }
.robot-list { display: flex; flex-direction: column; gap: 4px; }
.robot-row {
  display: flex; align-items: center; height: 40px; flex-shrink: 0; border-radius: 4px;
  &:hover, &:focus-within, &.previewed { background: #f0f0f0; }
  &.selected { background: #e8f0fe; .robot-row-select { color: #146ef5; } }
}
.robot-row-select {
  display: flex; align-items: center; gap: 8px; flex: 1; min-width: 0; height: 40px; padding: 8px;
  border: 0; border-radius: 0 4px 4px 0; color: #222; background: transparent; font: inherit; text-align: left; cursor: pointer;
}
.robot-avatar { flex-shrink: 0; width: 24px; height: 24px; border-radius: 6px; object-fit: cover; }
.robot-name { min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 14px; line-height: 22px; }
.drag-btn {
  display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0; width: 32px; height: 32px;
  padding: 0; border: 0; border-radius: 4px; background: transparent; cursor: grab; touch-action: none;
  &:active { cursor: grabbing; }
}
.picker-divider { height: 1px; margin: 4px 0; background: #e8e8e8; }
.picker-add {
  display: flex; align-items: center; gap: 8px; width: 100%; min-height: 40px; padding: 8px 8px 8px 40px;
  border: 0; border-radius: 4px; background: transparent; color: #222; font: inherit; font-size: 14px; line-height: 22px; cursor: pointer;
  span { flex: 1; min-width: 0; text-align: left; }
  img { flex-shrink: 0; }
  &:hover, &.active { background: #f0f0f0; }
}
.picker-create-menu {
  display: flex; flex-direction: column; gap: 2px; width: 180px; padding: 4px;
  button {
    display: flex; align-items: center; gap: 8px; width: 100%; min-height: 32px; padding: 4px 8px;
    border: 0; border-radius: 4px; color: #222; background: transparent; font: inherit; font-size: 14px; line-height: 22px; text-align: left; cursor: pointer;
    img { flex-shrink: 0; }
    &:hover { background: #f0f0f0; }
  }
}
.robot-preview {
  display: flex; flex-direction: column; gap: 16px; width: 280px; max-width: calc(100vw - 24px); padding: 24px;
  h3 { margin: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 16px; line-height: 24px; font-weight: 500; color: #080808; }
}
.preview-header { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.preview-avatar { flex-shrink: 0; width: 40px; height: 40px; border-radius: 10px; object-fit: cover; }
.robot-type {
  padding: 2px 6px; border-radius: 4px; background: #f1f6fe; color: #146ef5; font-size: 12px; line-height: 20px; white-space: nowrap;
  &.agent { color: #7a3dff; background: #f5f5f5; }
}
.robot-description { display: -webkit-box; min-height: 44px; max-height: 44px; margin: 0; overflow: hidden; color: #5a5a5a; font-size: 14px; line-height: 22px; overflow-wrap: anywhere; -webkit-line-clamp: 2; -webkit-box-orient: vertical; }
.robot-actions {
  display: flex; align-items: center; gap: 8px;
  button {
    display: inline-flex; align-items: center; justify-content: center; gap: 8px; flex: 1; min-width: 0; min-height: 40px;
    padding: 8px; border: 0; border-radius: 6px; color: #5a5a5a; background: transparent; font: inherit; font-size: 14px; line-height: 22px; cursor: pointer;
    img { flex-shrink: 0; }
    &:hover { background: #f0f0f0; }
  }
}
.action-divider { flex-shrink: 0; width: 1px; height: 24px; background: #e8e8e8; }
.picker-feedback { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 12px; min-height: 128px; padding: 16px; color: #5a5a5a; font-size: 14px; }
.robot-row-select:disabled, .drag-btn:disabled, .picker-add:disabled, .picker-create-menu button:disabled, .robot-actions button:disabled { cursor: not-allowed; opacity: .5; }
.robot-row-select:focus-visible, .drag-btn:focus-visible, .picker-add:focus-visible, .picker-create-menu button:focus-visible, .robot-actions button:focus-visible { outline: 2px solid #146ef5; outline-offset: -2px; }
</style>
