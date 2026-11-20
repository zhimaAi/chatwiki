<template>
  <PromptDialog v-if="props.open && draft" :robot-info="draft" :robot-list="workflows" :saving="props.loading" @save="patch => emit('save', patch)" @cancel="emit('close')" />
</template>

<script setup>
import { onScopeDispose, ref, watch } from 'vue'
import { App } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getChatVariables, getRobotList } from '@/api/robot'
import { useRobotStore } from '@/stores/modules/robot'
import { cloneConfig } from '../../config-session'
import PromptDialog from '../../basic-settings/components/prompt-dialog.vue'

const props = defineProps({
  robot: { type: Object, default: null },
  open: { type: Boolean, default: false },
  robotKey: { type: String, default: '' },
  loading: { type: Boolean, default: false }
})
const emit = defineEmits(['close', 'save', 'variablesLoaded'])
const { t } = useI18n('views.customer-service.chat')
const { t: commonT } = useI18n('views.customer-service.index')
const { message } = App.useApp()
const robotStore = useRobotStore()
const draft = ref(null)
const workflows = ref([])
let requestVersion = 0
let variablesReady = false

watch([() => props.open, () => props.robotKey], async ([open, robotKey]) => {
  const version = ++requestVersion
  draft.value = null
  workflows.value = []
  variablesReady = false
  if (!open || !robotKey || !props.robot) return
  // 复用设置页的数据规范化；编辑器只修改副本，不把未保存的提示词写入共享表单。
  try {
    robotStore.setRobotInfo(cloneConfig(props.robot))
    draft.value = cloneConfig(robotStore.robotInfo)
  } catch {
    message.error(commonT('error_robot'))
    emit('close')
    return
  }
  robotStore.chatVariables = []
  const results = await Promise.allSettled([
    getChatVariables({ robot_key: robotKey }),
    getRobotList({ application_type: 1 })
  ])
  if (!props.open || robotKey !== props.robotKey || version !== requestVersion) return
  variablesReady = true
  try {
    const [variables, robots] = results
    if (robots.status === 'fulfilled' && Array.isArray(robots.value?.data)) workflows.value = robots.value.data
    if (variables.status === 'fulfilled' && Array.isArray(variables.value?.data)) {
      const list = variables.value.data.map(item => ({ ...item, options: item.options ? JSON.parse(item.options) : [] }))
      robotStore.chatVariables = list
    }
    if (results.some(result => result.status === 'rejected' || !Array.isArray(result.value?.data))) {
      message.warning(t('error_prompt_options'))
    }
  } catch {
    message.warning(t('error_prompt_options'))
  }
}, { immediate: true, flush: 'sync' })

// 编辑器新增变量后也同步当前对话的待填项，保留会话中的已填值。
watch(() => robotStore.chatVariables, list => {
  if (props.open && draft.value && variablesReady) emit('variablesLoaded', cloneConfig(list))
}, { flush: 'sync' })
onScopeDispose(() => { requestVersion++ })
</script>
