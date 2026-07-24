<template>
  <a-modal
    wrapClassName="remove-model-padding"
    :open="visible"
    :title="t('modal_task_log')"
    :width="680"
    :footer="null"
    :destroyOnClose="false"
    :bodyStyle="{'max-height': '70vh', 'overflow-y': 'auto', 'padding': '0 24px 24px'}"
    @cancel="handleCancel"
  >
    <a-spin :spinning="loading">
      <div v-if="errorMsg" class="error-box">
        <div class="block-title">{{ t('label_error_message') }}</div>
        <pre class="log-content error-content">{{ errorMsg }}</pre>
      </div>

      <div class="log-box">
        <div class="block-title">{{ t('label_model_result') }}</div>
        <template v-if="debugLogs.length">
          <div v-for="(item, index) in debugLogs" :key="index" class="log-item">
            <div class="log-type">{{ item?.type || 'log' }}</div>
            <pre class="log-content">{{ formatContent(item?.content ?? item) }}</pre>
          </div>
        </template>
        <a-empty v-else-if="!loading" :description="t('empty_log')" />
      </div>
    </a-spin>
  </a-modal>
</template>

<script setup>
import { ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getDocToSkillTaskInfo } from '@/api/clawbot'

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
const debugLogs = ref([])
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
  debugLogs.value = []
  errorMsg.value = ''
  loading.value = false
}

const formatContent = (content) => {
  if (typeof content === 'string') {
    return content
  }
  try {
    return JSON.stringify(content, null, 2)
  } catch {
    return String(content || '')
  }
}

const loadTaskLog = async () => {
  resetState()
  if (!props.taskId) {
    return
  }

  loading.value = true
  try {
    const res = await getDocToSkillTaskInfo({ id: props.taskId })
    if (res && (res.res === 0 || res.code === 0)) {
      const data = res.data || {}
      debugLogs.value = Array.isArray(data.debug_log) ? data.debug_log : []
      errorMsg.value = data.err_msg || ''
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

.error-box {
  margin-bottom: 16px;
}

.block-title,
.log-type {
  color: #262626;
  font-size: 14px;
  line-height: 22px;
}

.block-title {
  margin-bottom: 8px;
  font-weight: 600;
}

.log-item + .log-item {
  margin-top: 12px;
}

.log-type {
  margin-bottom: 4px;
  color: #595959;
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
