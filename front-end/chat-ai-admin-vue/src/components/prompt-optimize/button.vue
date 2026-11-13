<template>
  <div
    v-if="supported"
    class="prompt-optimize-entry"
    :class="{ 'is-conversation': scene === 'conversation' }"
  >
    <a-button type="link" size="small" :class="{ 'persona-entry': appearance === 'persona' }" @click="show"><SparkleIcon filled />{{ t(appearance === 'persona' ? 'shortTitle' : 'title') }}</a-button>
  </div>
  <PromptOptimize v-if="mounted" ref="modal" :savePrompt="save" :saveMode="saveMode" />
</template>
<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { getRobotInfo, saveRobot } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'
import PromptOptimize from './index.vue'
import SparkleIcon from './sparkle-icon.vue'
import { useRobotStore } from '@/stores/modules/robot'
import { useClawbotStore } from '@/stores/modules/clawbot'
import { useCustomerServiceStore } from '@/stores/modules/customer-service'
const props = defineProps({
  robot: { type: Object, required: true },
  messageId: { type: [String, Number], default: 0 },
  scene: { type: String, default: 'config' },
  appearance: { type: String, default: 'link' },
  saveMode: { type: String, default: 'persist', validator: (value) => ['persist', 'draft'].includes(value) },
  onSave: Function
})
const emit = defineEmits(['adopted'])
const { t } = useI18n('components.prompt-optimize')
const modal = ref(null)
defineOptions({ name: 'PromptOptimizeButton' })
const mounted = ref(false)
let targetId
watch(
  () => props.robot?.id,
  () => {
    mounted.value = false
  }
)
const supported = computed(
  () =>
    props.robot?.id &&
    [0, 2].includes(Number(props.robot.application_type)) &&
    (props.scene !== 'conversation' || Number(props.messageId) > 0)
)
const show = async () => {
  targetId = props.robot.id
  mounted.value = true
  await nextTick()
  modal.value.show({ robotId: targetId, scene: props.scene, messageId: props.messageId })
}
async function save(prompt) {
  const savedRobotId = String(targetId)
  if (props.saveMode === 'draft') {
    if (!props.onSave) return false
    const success = await props.onSave(prompt)
    if (success !== true || String(props.robot.id) !== savedRobotId) return false
    emit('adopted', prompt)
    return true
  }
  let success
  if (props.onSave) {
    success = await props.onSave(prompt)
  } else {
    // Fetch before saving: saveRobot expects the complete configuration.
    const response = await getRobotInfo({ id: savedRobotId })
    const data = { ...response.data, prompt, prompt_type: 0 }
    if (data.robot_avatar) {
      data.robot_avatar_url = data.robot_avatar
      delete data.robot_avatar
    }
    for (const key of ['welcomes', 'unknown_question_prompt']) {
      if (data[key] && typeof data[key] === 'object')
        data[key] = {
          ...data[key],
          question: (data[key].question || []).map((q) => (typeof q === 'string' ? q : q.content))
        }
    }
    for (const [key, value] of Object.entries(data))
      if (value !== null && typeof value === 'object') data[key] = JSON.stringify(value)
    const result = await saveRobot(data, Number(data.application_type))
    success = result.res === 0
  }
  if (success === false) return false

  // Saving survives closing the modal. Update only the robot this request saved.
  const robotStore = useRobotStore()
  const clawbotStore = useClawbotStore()
  if (String(robotStore.robotInfo.id) === savedRobotId) {
    robotStore.robotInfo.prompt = prompt
    robotStore.robotInfo.prompt_type = 0
  }
  if (String(clawbotStore.robotInfo?.id) === savedRobotId) {
    clawbotStore.updateRobotInfo({ prompt, prompt_type: 0 })
  }
  const customerServiceStore = useCustomerServiceStore()
  if (String(customerServiceStore.robot?.id) === savedRobotId) {
    customerServiceStore.mergeCurrentRobotInfo(customerServiceStore.robot, { prompt, prompt_type: 0 })
  }
  if (String(props.robot.id) === savedRobotId) emit('adopted', prompt)
  return true
}
</script>

<style scoped>
.prompt-optimize-entry {
  display: inline-block;
}

.prompt-optimize-entry.is-conversation {
  display: block;
}
.prompt-optimize-entry.is-conversation .ant-btn,
.prompt-optimize-entry.is-conversation .ant-btn:hover,
.prompt-optimize-entry.is-conversation .ant-btn:focus,
.prompt-optimize-entry.is-conversation .ant-btn:active {
  color: #2475fc;
  background: transparent;
  border-color: transparent;
  box-shadow: none;
}
.prompt-optimize-entry :deep(.ant-btn) { display: inline-flex; align-items: center; gap: 4px; }
.prompt-optimize-entry .ant-btn:focus-visible {
  outline: none;
  box-shadow: none;
}
.prompt-optimize-entry .persona-entry { height: 32px; padding: 0 14px; border: 0; border-radius: 8px; color: #2475fc; background: #e8f0ff; }
</style>
