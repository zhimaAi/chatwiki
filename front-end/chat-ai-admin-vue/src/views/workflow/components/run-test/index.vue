<template>
  <div>
    <cu-modal
      v-model:open="show"
      :width="modalWidth"
      :destroyOnClose="false"
      wrapClassName="run-test-panel-wrapper"
      @cancel="handleCancel"
    >
      <div class="run-test-modal-shell">
        <div class="modal-title-block">
          <a-dropdown :trigger="['click']">
            <div class="mode-title">
              <span class="title-text">{{ currentMode === 'chat' ? t('btn_chat_test') : t('btn_run_test') }}</span>
              <svg-icon class="arrow-icon" name="arrow-down"></svg-icon>
            </div>
            <template #overlay>
              <a-menu @click="handleModeMenuClick">
                <a-menu-item key="chat" v-if="hasSessionTrigger">
                  {{ t('btn_chat_test') }}
                </a-menu-item>
                <a-menu-item key="run">
                  {{ t('btn_run_test') }}
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
          <div class="run-detail" v-if="panelState.hasRunResult">
            <svg-icon name="status-success" class="success-icon"></svg-icon>
            <span>{{ formatTime(panelState.use_mills) }}</span>
            <i class="line"></i>
            <span>{{ panelState.use_token }} {{ t('label_tokens') }}</span>
          </div>
        </div>

        <RunTestPanel
          v-if="currentMode === 'run'"
          ref="runTestRef"
          :style="runTestLayoutStyle"
          :start_node_params="startNodeParams"
          @stateChange="handlePanelStateChange"
        />
        <ChatTestPanel
          v-if="currentMode === 'chat'"
          ref="chatTestRef"
          :style="chatTestLayoutStyle"
          :start_node_params="startNodeParams"
          @stateChange="handlePanelStateChange"
        />
      </div>
    </cu-modal>
  </div>
</template>

<script setup>
import { computed, nextTick, reactive, ref } from 'vue'
import { formatTime } from '../util'
import { useI18n } from '@/hooks/web/useI18n'
import RunTestPanel from './RunTestPanel.vue'
import ChatTestPanel from './ChatTestPanel.vue'
import CuModal from '@/components/common/cu-modal.vue'

const { t } = useI18n('views.workflow.components.run-test.index')

const CHAT_TEST_CHAT_WIDTH = 382
const CHAT_TEST_FLOW_WIDTH = 326
const CHAT_TEST_LOG_WIDTH = 382
const CHAT_TEST_TOTAL_WIDTH = CHAT_TEST_CHAT_WIDTH + CHAT_TEST_FLOW_WIDTH + CHAT_TEST_LOG_WIDTH

const props = defineProps({
  lf: {
    type: Object
  },
  start_node_params: {
    default: () => {},
    type: Object
  },
  isLockedByOther: { type: Boolean, default: false }
})

defineEmits(['save'])

const show = ref(false)
const currentMode = ref('')
const runTestRef = ref(null)
const chatTestRef = ref(null)
const panelState = reactive({
  hasRunTested: false,
  hasRunResult: false,
  use_token: 0,
  use_mills: 0
})

const startNodeParams = ref({
  diy_global: [],
  sys_global: [],
  trigger_list: []
})
const hasSessionTrigger = ref(false)

const refreshStartNodeParams = () => {
  const graphData = props.lf?.getGraphRawData?.()
  const startNode = graphData?.nodes?.find((node) => node.type === 'start-node')
  const nodeParams = startNode?.properties?.node_params

  if (nodeParams) {
    try {
      const params = typeof nodeParams === 'string' ? JSON.parse(nodeParams) : nodeParams
      const startParams = params?.start || params || {}
      startNodeParams.value = {
        ...startParams,
        diy_global: startParams.diy_global || [],
        sys_global: startParams.sys_global || [],
        trigger_list: startParams.trigger_list || (startParams.trigger ? [startParams.trigger] : [])
      }
      return
    } catch {
      // 保持默认参数，避免画布中的异常数据阻断测试弹窗打开。
    }
  }
  startNodeParams.value = {
    diy_global: [],
    sys_global: [],
    trigger_list: []
  }
}

const updateSessionTrigger = () => {
  const triggerList = startNodeParams.value.trigger_list
  hasSessionTrigger.value = triggerList.some((item) => item.trigger_type == 1)
  return hasSessionTrigger.value
}

const modalWidth = computed(() => {
  if (currentMode.value === 'chat') {
    return CHAT_TEST_TOTAL_WIDTH
  }
  return panelState.hasRunTested ? CHAT_TEST_TOTAL_WIDTH : 768
})

const chatTestLayoutStyle = computed(() => ({
  '--chat-test-chat-width': `${CHAT_TEST_CHAT_WIDTH}px`,
  '--chat-test-flow-width': `${CHAT_TEST_FLOW_WIDTH}px`,
  '--chat-test-log-width': `${CHAT_TEST_LOG_WIDTH}px`,
  '--node-run-logs-panel-box-sizing': 'border-box'
}))

const runTestLayoutStyle = computed(() => {
  if (!panelState.hasRunTested) {
    return {}
  }
  return chatTestLayoutStyle.value
})

const handlePanelStateChange = (state) => {
  panelState.hasRunTested = !!state.hasRunTested
  panelState.hasRunResult = !!state.hasRunResult
  panelState.use_token = state.use_token || 0
  panelState.use_mills = state.use_mills || 0
}

const resetPanelState = () => {
  handlePanelStateChange({
    hasRunTested: false,
    hasRunResult: false,
    use_token: 0,
    use_mills: 0
  })
}

const openCurrentPanel = () => {
  nextTick(() => {
    if (currentMode.value === 'chat') {
      chatTestRef.value?.open()
      return
    }
    runTestRef.value?.open()
  })
}

const handleModeMenuClick = ({ key }) => {
  if (key === currentMode.value) {
    return
  }
  chatTestRef.value?.abort()
  currentMode.value = key
  resetPanelState()
  openCurrentPanel()
}

const handleCancel = () => {
  chatTestRef.value?.abort()
  currentMode.value = ''
}

const open = () => {
  refreshStartNodeParams()
  currentMode.value = updateSessionTrigger() ? 'chat' : 'run'
  resetPanelState()
  show.value = true
  openCurrentPanel()
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.run-test-modal-shell {
  overflow: hidden;
  border-radius: 8px;
  background: #fff;
}

.modal-title-block {
  display: flex;
  align-items: center;
  min-height: 56px;
  padding: 0 48px 0 24px;
  border-bottom: 1px solid #f0f0f0;
  gap: 16px;
  .run-detail {
    display: flex;
    align-items: center;
    height: 20px;
    padding: 0 6px;
    font-size: 12px;
    color: #21A665;
    font-weight: 400;
    border-radius: 4px;
    background: #CAFCE4;

    .success-icon{
      margin-right: 4px;
      font-size: 14px;
      color: #fff;
    }

    .line {
      display: inline-block;
      width: 1px;
      height: 12px;
      margin: 0 8px;
      background: rgba(0, 0, 0, 0.25);
    }
  }
}

.mode-title {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  user-select: none;
  color: #262626;
  .title-text{
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }
  .arrow-icon{
    font-size: 16px;
  }
}

</style>
