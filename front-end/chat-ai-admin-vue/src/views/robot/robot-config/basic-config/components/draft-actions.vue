<template>
  <div class="draft-actions" v-bind="$attrs">
    <span v-if="robotInfo.has_draft" class="draft-tip">{{ t('msg_unpublished_changes') }}</span>
    <span class="published-time" v-else>{{ t('label_published_time', { time: publishedTime }) }}</span>
    <a-button :disabled="(!robotInfo.has_draft && !appliedAction) || saving || busy || !!action || appliedAction === 'publish'" @click="undoConfirmOpen = true">
      {{ t('btn_undo_changes') }}
    </a-button>
    <a-button class="publish-btn" type="primary" :loading="action === 'publish'" :disabled="(!robotInfo.has_draft && !appliedAction) || saving || busy || !!action || appliedAction === 'undo'" @click="handleAction('publish')">
      {{ t('btn_publish') }}
    </a-button>
  </div>
  <a-modal
    v-model:open="undoConfirmOpen"
    :footer="null"
    :closable="false"
    :mask-closable="!action"
    :keyboard="!action"
    :width="416"
    centered
  >
    <div class="undo-confirm-content">
      <ExclamationCircleFilled />
      <div>
        <div class="undo-confirm-title">{{ t('confirm_undo_changes') }}</div>
        <p>{{ t('confirm_undo_description') }}</p>
      </div>
    </div>
    <div class="undo-confirm-footer">
      <a-button :disabled="!!action" @click="undoConfirmOpen = false">{{ t('btn_cancel') }}</a-button>
      <a-button type="primary" :loading="action === 'undo'" :disabled="saving" @click="handleAction('undo')">{{ t('btn_confirm') }}</a-button>
    </div>
  </a-modal>
</template>

<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { CheckOutlined, ExclamationCircleFilled } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { useI18n } from '@/hooks/web/useI18n'
import { publishRobot, undoRobotDraft } from '@/api/robot'

defineOptions({ inheritAttrs: false })
const props = defineProps({
  robotInfo: { type: Object, required: true },
  saving: { type: Boolean, default: false },
  busy: { type: Boolean, default: false },
  refresh: { type: Function, required: true }
})
const emit = defineEmits(['update:busy'])
const { t } = useI18n('views.robot.robot-config.basic-config.index')
const action = ref('')
const appliedAction = ref('')
const undoConfirmOpen = ref(false)
const identity = computed(() => String(props.robotInfo.id) + ':' + props.robotInfo.robot_key)
const publishedTime = computed(() => props.robotInfo.published_time
  ? dayjs.unix(props.robotInfo.published_time).format('YY-MM-DD HH:mm')
  : '--')
let version = 0
let alive = true
watch(identity, () => {
  version++
  action.value = ''
  appliedAction.value = ''
  undoConfirmOpen.value = false
  emit('update:busy', false)
})
async function handleAction(kind) {
  if ((!props.robotInfo.has_draft && !appliedAction.value) || props.saving || props.busy || action.value) return
  if (appliedAction.value && appliedAction.value !== kind) return
  const current = ++version
  const source = identity.value
  const id = props.robotInfo.id
  const robotKey = props.robotInfo.robot_key
  const isCurrent = () => alive && current === version && source === identity.value
  action.value = kind
  emit('update:busy', true)
  try {
    const request = kind === 'publish' ? publishRobot : undoRobotDraft
    // 接口成功但回显刷新失败时，重试只刷新，避免重复发布或撤销。
    if (!appliedAction.value) {
      await request({ id })
      if (!isCurrent()) return
      appliedAction.value = kind
    }
    if (!isCurrent()) return
    const refreshed = await props.refresh(id, robotKey)
    if (!isCurrent() || refreshed === false) return
    appliedAction.value = ''
    undoConfirmOpen.value = false
    message.success(t(kind === 'publish' ? 'msg_publish_success' : 'msg_undo_success'))
  } catch {
    // 请求错误由 axios 或配置上下文提示，保留弹窗和配置以便重试。
  } finally {
    if (isCurrent()) {
      action.value = ''
      emit('update:busy', false)
    }
  }
}
onBeforeUnmount(() => { alive = false; version++; emit('update:busy', false) })
</script>

<style lang="less" scoped>
.draft-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 400;
  white-space: nowrap;

  :deep(.ant-btn) {
    min-width: 88px;
    height: 32px;
    border-radius: 6px;
  }

  :deep(.publish-btn) {
    min-width: 80px;
  }

  .published-time {
    color: #8c8c8c;
  }

  .draft-tip {
    color: #ff4d4f;
  }
}

.undo-confirm-content {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 8px 4px 0;
  line-height: 22px;

  > .anticon {
    flex-shrink: 0;
    margin-top: 1px;
    color: #ff9900;
    font-size: 20px;
  }

  .undo-confirm-title {
    font-size: 14px;
    font-weight: 600;
    color: #262626;
  }

  p {
    margin: 4px 0 0;
    font-size: 12px;
    line-height: 20px;
    color: #8c8c8c;
  }
}

.undo-confirm-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 20px 4px 0;

  :deep(.ant-btn) {
    min-width: 56px;
    height: 28px;
    padding: 0 12px;
    font-size: 12px;
    border-radius: 4px;
  }
}

</style>
