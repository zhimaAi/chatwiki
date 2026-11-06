<template>
  <section class="chat-preview">
    <header><span>{{ t('chat_test') }}</span><a-button size="small" @click="reset">{{ t('new_chat') }}</a-button></header>
    <div v-if="error" class="error"><a-alert type="error" :message="error" show-icon /><a-button @click="reset">{{ t('retry') }}</a-button></div>
    <div class="messages" ref="scrollBox">
      <a-spin v-if="initializing" />
      <MessageList v-else show-sender-name :messages="session.messageList" :robot-info="context.robotInfo" @click-msg-meun="value = $event" @open-prompt-log="item => promptLog?.open(item)" @open-library="openLibrary" />
    </div>
    <div class="composer">
      <a-button v-if="session.chat_variables.need_fill_variable" type="link" @click="variables?.handleEdit()">{{ t('fill_variables') }}</a-button>
      <MessageInput v-model:value="value" v-model:file-list="files" :loading="busy || session.sendLock || initializing" :show-stop="!initializing" :show-upload="Number(context.robotInfo.question_multiple_switch) === 1" @send="send" @stop="stop" />
    </div>
    <VariableModal ref="variables" />
    <PromptLogAlert ref="promptLog" />
    <LibraryInfoAlert ref="libraryInfo" />
  </section>
</template>
<script setup>
import { inject, nextTick, onBeforeUnmount, provide, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useUserStore } from '@/stores/modules/user'
import { createRobotChatSession } from '@/stores/modules/chat'
import { checkChatRequestPermission } from '@/api/robot'
import MessageList from '@/views/robot/robot-test/components/message-list.vue'
import MessageInput from '@/views/robot/robot-test/components/message-input.vue'
import VariableModal from '@/views/robot/robot-test/components/variable-modal/index.vue'
import PromptLogAlert from '@/views/robot/robot-test/components/prompt-log-alert.vue'
import LibraryInfoAlert from '@/views/robot/robot-test/components/library-info-alert.vue'
import { generateRandomId } from '@/utils'
import { serializeRobotConfig } from '../config-session'
const props = defineProps({ resetKey: { type: Number, default: 0 } })
const { t } = useI18n('views.customer-service.index')
const { t: chatTestT } = useI18n('views.robot.robot-test.index')
const user = useUserStore()
const context = inject('customerRobotContext')
const session = reactive(createRobotChatSession({ isolatedIM: true }))
provide('robotChatSession', session)
const value = ref('')
const files = ref([])
const error = ref('')
const busy = ref(false)
const initializing = ref(false)
const scrollBox = ref(null)
const variables = ref(null)
const promptLog = ref(null)
const libraryInfo = ref(null)
let version = 0
let alive = true
let configurationRevision = 0
const welcomeParams = { robot_key: context.robotInfo.robot_key, openid: user.user_id, name: user.user_name, nickname: user.user_name, avatar: user.avatar, is_background: 1 }

function stop() { version++; busy.value = false; initializing.value = false; session.stopMessage() }
async function reset() {
  const current = ++version
  session.$reset()
  value.value = ''
  files.value = []
  busy.value = false
  error.value = ''
  initializing.value = true
  try {
    await context.waitForSave()
    if (!alive || current !== version) return
    const response = await session.createChat({ ...welcomeParams, dialogue_id: 0 })
    if (!alive || current !== version) return
    const welcomeMessage = session.messageList.find(item => item.isWelcome)
    if (welcomeMessage && !response?.data?.robot?.welcomes) {
      const language = response?.data?.robot?.lang_key || 'zh-CN'
      const languageConfig = context.robotInfo.multi_lang_configs?.find(item => item.lang_key === language)
      const welcomes = languageConfig ? languageConfig.welcomes : context.robotInfo.welcomes
      if (welcomes) {
        const menu = JSON.parse(serializeRobotConfig({ welcomes }).welcomes)
        welcomeMessage.menu_json = menu
        session.robot.welcomes = menu
      }
    }
    configurationRevision = context.savedRevision?.value || 0
  } catch (err) { if (alive && current === version) error.value = err.message || t('preview_failed') }
  finally { if (alive && current === version) initializing.value = false }
}
async function send() {
  if (busy.value || initializing.value || session.sendLock || (!value.value.trim() && !files.value.length)) return
  if (files.value.some(file => file.status !== 'done')) { message.error(t('wait_for_attachments')); return }
  const current = version
  busy.value = true
  try {
    await context.waitForSave()
    if (!alive || current !== version) return
    const latestRevision = context.savedRevision?.value || 0
    if (latestRevision !== configurationRevision) {
      await session.createChat({ ...welcomeParams, dialogue_id: session.dialogue_id, refreshOnly: true })
      if (!alive || current !== version) return
      configurationRevision = latestRevision
      if (session.chat_variables.need_fill_variable) return
    }
    session.robot.robot_name = context.robotInfo.robot_name || session.robot.robot_name
    session.robot.robot_avatar = context.robotInfo.robot_avatar_url || session.robot.robot_avatar
    session.robot.enable_question_guide = String(context.robotInfo.enable_question_guide) === 'true'
    session.robot.enable_common_question = String(context.robotInfo.enable_common_question) === 'true'
    const checked = await checkChatRequestPermission({ robot_key: session.robot.robot_key, openid: session.robot.openid, question: value.value, form_ids: session.robot.form_ids, dialogue_id: session.dialogue_id })
    if (!alive || current !== version) return
    if (checked.data?.words?.length) { message.error(chatTestT('msg_sensitive_words', { words: checked.data.words.join(';') })); return }
    let content = value.value.trim()
    if (Number(context.robotInfo.question_multiple_switch) === 1) {
      content = JSON.stringify([
        ...(content ? [{ type: 'text', uid: generateRandomId(16), text: content }] : []),
        ...files.value.map(file => ({ type: 'image_url', uid: file.uid, image_url: { url: file.url } }))
      ])
    }
    session.sendMessage({ message: content, global: JSON.stringify({ id: context.robotInfo.id, robot_key: context.robotInfo.robot_key }), msg_type: Number(context.robotInfo.question_multiple_switch) === 1 ? 99 : 1 })
    value.value = ''
    files.value = []
  } catch (err) { if (alive && current === version) message.error(err.message || t('preview_failed')) }
  finally { if (alive && current === version) busy.value = false }
}
function openLibrary(files, file) { libraryInfo.value?.open(files, file) }
watch(() => props.resetKey, reset, { immediate: true })
watch(() => session.messageList, async () => { await nextTick(); if (scrollBox.value) scrollBox.value.scrollTop = scrollBox.value.scrollHeight }, { deep: true })
onBeforeUnmount(() => { alive = false; version++; session.$reset() })
</script>
<style lang="less" scoped>
.chat-preview { height: 100%; display: flex; flex-direction: column; background: #fff; overflow: hidden; }
header { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 24px; font-size: 16px; line-height: 24px; font-weight: 500; color: #222; }
header :deep(.ant-btn) { height: 32px; padding: 0 12px; border-color: #d8d8d8; border-radius: 6px; font-size: 14px; font-weight: 400; box-shadow: none; }
.messages { flex: 1; min-height: 0; overflow: auto; padding: 0 24px; }
.composer { padding: 16px 24px 24px; }
.composer :deep(.message-input-wrapper) { border-radius: 12px; border-color: #d8d8d8; }
.error { padding: 0 24px 16px; }
.messages :deep(.robot-message) {
  display: block;
  &:first-child { margin-top: 0; }
  > .itme-left { display: flex; align-items: center; gap: 8px; margin: 0 0 12px; }
  .robot-avatar { width: 28px; height: 28px; border-radius: 8px; object-fit: cover; }
  .sender-name { min-width: 0; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; font-size: 12px; line-height: 20px; color: #7a7a7a; }
  .item-body { padding: 0; border-radius: 0; background: transparent; --robot-bubble-color: #5a5a5a; }
  .item-body > .message-content, .item-body > .markdown-content { padding: 16px; border-radius: 12px; background: #f5f5f5; }
  .message-menus .menu-item { padding: 12px; margin-top: 8px; border: 1px solid #d8d8d8; border-radius: 6px; background: #fff; color: #2475fc; overflow-wrap: anywhere; }
}
</style>
