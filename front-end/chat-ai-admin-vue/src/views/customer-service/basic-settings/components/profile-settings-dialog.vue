<template>
  <SettingsDialog v-if="editing" :open="true" :width="480" :title="commonT('profile_edit_title')" :ok-text="commonT('save')" :cancel-text="commonT('cancel')" :confirm-loading="saving" :ok-button-props="{ disabled: reading }" :closable="!saving" :mask-closable="!saving" :keyboard="!saving" :cancel-button-props="{ disabled: saving }" @ok="persist" @cancel="close">
    <div class="profile-form" :inert="saving ? '' : null">
      <div class="avatar-upload">
        <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/webp" hidden @change="changeAvatar" />
        <button type="button" class="avatar-button" :disabled="saving" :aria-label="commonT('profile_upload_avatar')" @click="chooseFile"><img :src="draft.robot_avatar_url" alt="" /></button>
        <span class="upload-hint">{{ commonT('profile_avatar_hint') }}</span>
      </div>
      <div class="default-avatar-field">
        <span>{{ t('label_default_avatar') }}</span>
        <div class="default-avatar-list">
          <button v-for="(avatar, index) in defaultAvatars" :key="avatar" type="button" class="default-avatar" :class="{ active: draft.robot_avatar_url === avatar }" :aria-label="`${t('label_default_avatar')} ${index + 1}`" :aria-pressed="draft.robot_avatar_url === avatar" @click="useDefaultAvatar(avatar)"><img :src="avatar" alt="" /></button>
        </div>
      </div>
      <div class="profile-field">
        <label :for="nameId">{{ commonT('profile_robot_name') }}</label>
        <a-input :id="nameId" v-model:value="draft.robot_name" :placeholder="t('ph_input_robot_name')" :status="nameInvalid ? 'error' : undefined" :aria-invalid="nameInvalid" @change="nameInvalid = false" @press-enter="persist" />
        <span v-if="nameInvalid" class="field-error" role="alert">{{ commonT('profile_name_required') }}</span>
      </div>
      <div class="profile-field intro-field">
        <label :for="introId">{{ t('label_intro') }}</label>
        <a-textarea :id="introId" v-model:value="draft.robot_intro" :placeholder="t('ph_input_robot_intro')" :rows="3" />
      </div>
      <a-alert v-if="saveFailed" type="error" :message="commonT('save_failed')" show-icon>
        <template #action><a-button size="small" :loading="saving" @click="persist">{{ commonT('retry') }}</a-button></template>
      </a-alert>
    </div>
  </SettingsDialog>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getBase64 } from '@/utils'
import { DEFAULT_ROBOT_AVATAR, DEFAULT_ROBOT_AVATAR_LIST, DEFAULT_WORKFLOW_AVATAR } from '@/constants'
import { saveConfigurationPatch } from '../configuration'
import SettingsDialog from './settings-dialog.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.basic-config')
const { t: commonT } = useI18n('views.customer-service.index')
const context = inject('customerRobotContext')
const defaultAvatar = computed(() => Number(context.robotInfo.application_type) === 1 ? DEFAULT_WORKFLOW_AVATAR : DEFAULT_ROBOT_AVATAR)
const defaultAvatars = computed(() => [defaultAvatar.value, ...DEFAULT_ROBOT_AVATAR_LIST])
const nameId = `profile-name-${context.robotInfo.id}`
const introId = `profile-intro-${context.robotInfo.id}`
const draft = reactive({ robot_name: '', robot_intro: '', robot_avatar: '', robot_avatar_url: '' })
const editing = ref(false)
const saving = ref(false)
const reading = ref(false)
const nameInvalid = ref(false)
const saveFailed = ref(false)
const fileInput = ref(null)
let previewVersion = 0
let alive = true

function open() {
  if (saving.value) return
  previewVersion++
  const info = context.robotInfo
  Object.assign(draft, { robot_name: info.robot_name || '', robot_intro: info.robot_intro || '', robot_avatar: info.robot_avatar_url || defaultAvatar.value, robot_avatar_url: info.robot_avatar_url || defaultAvatar.value })
  reading.value = false
  nameInvalid.value = false
  saveFailed.value = false
  editing.value = true
}
function close() {
  if (saving.value) return
  previewVersion++
  reading.value = false
  editing.value = false
}
function chooseFile() {
  fileInput.value.value = ''
  fileInput.value.click()
}
async function changeAvatar(event) {
  const file = event.target.files?.[0]
  if (!file || saving.value) return
  if (!(file.type ? /^image\/(jpeg|png|webp)$/.test(file.type) : /\.(jpe?g|png|webp)$/i.test(file.name))) {
    message.error(commonT('profile_avatar_type_error'))
    return
  }
  if (file.size > 100 * 1024) {
    message.error(t('msg_avatar_size_limit'))
    return
  }
  const version = ++previewVersion
  reading.value = true
  try {
    const url = await getBase64(file)
    if (!alive || !editing.value || version !== previewVersion) return
    draft.robot_avatar = file
    draft.robot_avatar_url = url
  } catch {
    if (alive && editing.value && version === previewVersion) message.error(commonT('profile_avatar_read_error'))
  } finally {
    if (alive && version === previewVersion) reading.value = false
  }
}
function useDefaultAvatar(avatar) {
  if (saving.value) return
  // 切回默认头像后，较慢的文件预览结果不能覆盖当前选择。
  previewVersion++
  reading.value = false
  draft.robot_avatar = avatar
  draft.robot_avatar_url = avatar
}
async function persist() {
  if (saving.value || reading.value) return
  nameInvalid.value = !draft.robot_name.trim()
  if (nameInvalid.value) return
  saving.value = true
  try {
    const success = await saveConfigurationPatch(context, { ...draft })
    if (!alive) return
    saveFailed.value = !success
    if (success) editing.value = false
  } catch {
    if (alive) saveFailed.value = true
  } finally {
    if (alive) saving.value = false
  }
}
onBeforeUnmount(() => { alive = false; previewVersion++ })
defineExpose({ open })
</script>

<style lang="less" scoped>
.profile-form { display: flex; flex-direction: column; gap: 24px; font-size: 14px; line-height: 22px; color: #222; }
.avatar-upload { display: flex; flex-direction: column; align-items: center; gap: 8px; }
.avatar-button, .default-avatar { display: flex; align-items: center; justify-content: center; width: 48px; height: 48px; flex-shrink: 0; padding: 0; border: 0; border-radius: 6px; background: #e8f0fe; overflow: hidden; cursor: pointer; }
.avatar-button img { width: 48px; height: 48px; object-fit: cover; }
.upload-hint { color: #5a5a5a; font-size: 12px; line-height: 20px; text-align: center; }
.default-avatar-field { display: flex; flex-direction: column; gap: 12px; }
.default-avatar-list { display: flex; flex-wrap: wrap; gap: 8px; }
.default-avatar { padding: 2px; border: 2px solid transparent; background: #fff; img { width: 40px; height: 40px; object-fit: contain; } &.active { border-color: #146ef5; } }
.avatar-button:focus-visible, .default-avatar:focus-visible { outline: 2px solid #146ef5; outline-offset: 2px; }
.profile-field { display: flex; flex-direction: column; gap: 8px; }
.profile-field :deep(.ant-input) { min-height: 32px; padding: 4px 12px; border-radius: 6px; border-color: #d8d8d8; font-size: 14px; line-height: 22px; }
.intro-field { padding-bottom: 6px; }
.intro-field :deep(textarea.ant-input) { height: 88px; min-height: 88px; padding: 8px 12px; resize: vertical; }
.field-error { color: #ff4d4f; font-size: 12px; line-height: 20px; }
</style>
