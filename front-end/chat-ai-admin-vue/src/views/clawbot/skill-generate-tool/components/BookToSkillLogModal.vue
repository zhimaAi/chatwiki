<template>
  <a-modal
    class="book-skill-log-modal"
    :open="visible"
    :title="t('modal_task_log')"
    :width="680"
    :footer="null"
    :destroyOnClose="false"
    @cancel="handleCancel"
  >
    <a-spin :spinning="loading">
      <div v-if="errorMsg" class="error-box">
        <div class="block-title">{{ t('label_error_message') }}</div>
        <pre class="log-content error-content">{{ errorMsg }}</pre>
      </div>

      <div class="log-box" v-if="false">
        <div class="block-title">{{ t('label_execution_log') }}</div>
        <pre v-if="logContent" class="log-content">{{ logContent }}</pre>
        <a-empty v-else-if="!loading" :description="t('empty_log')" />
      </div>
    </a-spin>
  </a-modal>
</template>

<script setup>
import { ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getBookToSkillTaskLog } from '@/api/clawbot'

const { t } = useI18n('views.clawbot.skill-generate-tool.index')

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  taskId: {
    type: [String, Number],
    default: ''
  }
})

const emit = defineEmits(['update:visible'])

const loading = ref(false)
const logContent = ref('')
const errorMsg = ref('')

watch(
  () => [props.visible, props.taskId],
  ([visible]) => {
    if (visible) {
      loadTaskLog()
    }
  }
)

const resetState = () => {
  logContent.value = ''
  errorMsg.value = ''
  loading.value = false
}

const loadTaskLog = async () => {
  resetState()
  if (!props.taskId) {
    return
  }

  loading.value = true
  try {
    const res = await getBookToSkillTaskLog({
      task_id: props.taskId
    })
    if (res && (res.res === 0 || res.code === 0)) {
      const data = res.data || {}
      logContent.value = data.log_content || ''
      errorMsg.value = data.error_msg || ''
    } else {
      message.error(res?.msg || t('msg_fetch_task_log_failed'))
    }
  } catch (err) {
    console.error('获取Book转Skill任务日志失败', err)
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  emit('update:visible', false)
}
</script>

<style lang="less" scoped>
.book-skill-log-modal {
  :deep(.ant-modal-content) {
    border-radius: 12px;
  }

  :deep(.ant-modal-title) {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }
}

.error-box {
  margin-bottom: 16px;
}

.block-title {
  margin-bottom: 8px;
  color: #262626;
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
}

.log-content {
  max-height: 360px;
  margin: 0;
  padding: 12px;
  border-radius: 6px;
  color: #3d3d3d;
  background: #f7f8fa;
  font-size: 13px;
  line-height: 20px;
  white-space: pre-wrap;
  word-break: break-word;
  overflow: auto;
}

.error-content {
  color: #d4380d;
  background: #fff2e8;
}
</style>
