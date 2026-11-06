<style lang="less" scoped>
.setting-box {
  position: relative;
  .robot-info-box {
    .robot-prompt {
      display: block;
      line-height: 22px;
      font-size: 14px;
      white-space: normal;
      word-break: break-word;
      color: #595959;
    }

    .guide-num-input {
      width: 64px;
      margin: 0 6px;
      vertical-align: middle;
    }

    .guide-action-row {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-top: 16px;
      line-height: 22px;
      font-size: 14px;
      color: #595959;
    }
  }
  .modal-item-box {
    padding: 24px 0;
    display: flex;
    flex-direction: column;
    gap: 24px;

    .modal-item {
      .label {
        margin-bottom: 8px;
        color: #262626;
        font-size: 14px;
        line-height: 22px;
      }
    }
  }
  .current-mode {
    max-width: 180px;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    color: #595959;
    font-size: 14px;
  }
}
</style>

<template>
  <edit-box class="setting-box" :compact="compact" :title="t('title_suggested_questions')" icon-name="suggested-issues">
    <template #extra>
      <template v-if="compact">
        <span class="current-mode" :title="t('label_current_mode') + currentModeText">{{ t('label_current_mode') }}{{ currentModeText }}</span>
        <a-button size="small" @click="openSettingModal">{{ t('btn_generation_setting') }}</a-button>
      </template>
      <a-switch v-if="!compact"
        @change="handleEdit"
        checkedValue="true"
        unCheckedValue="false"
        v-model:checked="robotInfo.enable_question_guide"
        :checked-children="t('btn_on')"
        :un-checked-children="t('btn_off')"
      />
    </template>
    <template #switch>
      <a-switch v-if="compact"
        @change="handleEdit"
        checkedValue="true"
        unCheckedValue="false"
        v-model:checked="robotInfo.enable_question_guide"
        :checked-children="t('btn_on')"
        :un-checked-children="t('btn_off')"
      />
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt" :title="t('msg_guide_description', { num: robotInfo.question_guide_num })">
        <span>{{ descriptionParts[0] }}</span>
        <a-input-number
          v-model:value="robotInfo.question_guide_num"
          class="guide-num-input"
          size="small"
          :precision="0"
          :min="1"
          :max="10"
          @blur="handleEdit"
        />
        <span>{{ descriptionParts[1] }}</span>
      </div>
      <div class="guide-action-row" v-if="!compact">
        <a-button @click="openSettingModal">{{ t('btn_generation_setting') }}</a-button>
        <span>{{ t('label_current_mode') }}{{ currentModeText }}</span>
      </div>
    </div>

    <div class="modal-box" ref="modalBoxRef">
      <component
        :is="compact ? SettingsDialog : AntModal"
        :getContainer="compact ? undefined : () => $refs.modalBoxRef"
        :open="settingVisible"
        :width="compact ? 640 : 520"
        :title="t('title_generation_mode')"
        @ok="handleSaveSetting"
        :ok-text="compact ? commonT('confirm') : undefined"
        :cancel-text="compact ? commonT('cancel') : undefined"
        :confirm-loading="compact && saving"
        :closable="!compact || !saving"
        :mask-closable="!compact || !saving"
        :keyboard="!compact || !saving"
        :cancel-button-props="{ disabled: compact && saving }"
        @cancel="closeSettingModal"
      >
        <template v-if="compact" #closeIcon><img src="@/assets/img/customer-service/basic-settings/prompt/close.svg" alt="" /></template>
        <div class="modal-item-box" :class="{ 'workspace-generation-form': compact }" :inert="compact && saving ? '' : null">
          <div class="modal-item">
            <div class="label">{{ t('label_generation_mode') }}</div>
            <a-radio-group v-model:value="formState.question_guide_mode" :disabled="compact && saving" :aria-label="t('label_generation_mode')">
              <a-radio :value="0">{{ t('mode_default') }}</a-radio>
              <a-radio :value="1">{{ t('mode_prompt') }}</a-radio>
              <a-radio :value="2">{{ t('mode_workflow') }}</a-radio>
            </a-radio-group>
          </div>
          <div class="modal-item" v-if="formState.question_guide_mode == 1">
            <div class="label">{{ t('label_prompt') }}</div>
            <a-textarea
              v-model:value="formState.question_guide_prompt"
              :maxlength="3000"
              :show-count="!compact"
              :style="{ height: compact ? '128px' : '160px' }"
              :placeholder="t('ph_prompt')"
            />
          </div>
          <div class="modal-item" v-if="formState.question_guide_mode == 2">
            <div class="label">{{ t('label_workflow') }}</div>
            <a-select
              v-model:value="formState.question_guide_workflow_key"
              style="width: 100%"
              :placeholder="t('ph_workflow')"
              :options="workflowOptions"
            />
          </div>
          <a-alert v-if="compact && saveFailed" type="error" :message="commonT('save_failed')" show-icon />
        </div>
      </component>
    </div>
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, computed } from 'vue'
import { message, Modal as AntModal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import EditBox from './edit-box.vue'
import SettingsDialog from '@/views/customer-service/basic-settings/components/settings-dialog.vue'
import { saveConfigurationPatch } from '@/views/customer-service/basic-settings/configuration'

const { t } = useI18n('views.robot.robot-config.basic-config.components.suggested-issues')
const { t: commonT } = useI18n('views.customer-service.index')
const customerContext = inject('customerRobotContext', null)
const saving = ref(false)
const saveFailed = ref(false)

const { robotInfo, updateRobotInfo } = inject('robotInfo')
const props = defineProps({
  compact: { type: Boolean, default: false },
  robotList: {
    type: Array,
    default: () => []
  }
})

const settingVisible = ref(false)
const formState = reactive({
  question_guide_mode: 0,
  question_guide_prompt: '',
  question_guide_workflow_key: ''
})

const modeTextMap = computed(() => ({
  0: t('mode_default'),
  1: t('mode_prompt'),
  2: t('mode_workflow')
}))

const currentModeText = computed(() => {
  return modeTextMap.value[Number(robotInfo.question_guide_mode) || 0] || t('mode_default')
})

const descriptionParts = computed(() => t('msg_guide_description', { num: '@@NUM@@' }).split('@@NUM@@'))

const workflowOptions = computed(() => {
  return props.robotList
    .filter((item) => item.has_published == 1)
    .map((item) => ({
      label: item.robot_name,
      value: item.robot_key
    }))
})

const openSettingModal = () => {
  if (props.compact && saving.value) return
  formState.question_guide_mode = Number(robotInfo.question_guide_mode) || 0
  formState.question_guide_prompt = robotInfo.question_guide_prompt || ''
  formState.question_guide_workflow_key = robotInfo.question_guide_workflow_key || ''
  settingVisible.value = true
  saveFailed.value = false
}
const closeSettingModal = () => {
  if (!props.compact || !saving.value) settingVisible.value = false
}

const handleSaveSetting = async () => {
  if (props.compact && saving.value) return
  if (formState.question_guide_mode == 1 && !formState.question_guide_prompt.trim()) {
    return message.error(t('msg_prompt_required'))
  }

  if (formState.question_guide_mode == 2 && !formState.question_guide_workflow_key) {
    return message.error(t('msg_workflow_required'))
  }

  const patch = {
    question_guide_mode: formState.question_guide_mode,
    question_guide_prompt: formState.question_guide_mode == 1 ? formState.question_guide_prompt.trim() : '',
    question_guide_workflow_key: formState.question_guide_mode == 2 ? formState.question_guide_workflow_key : ''
  }
  if (!props.compact) {
    updateRobotInfo(patch)
    settingVisible.value = false
    return
  }
  saving.value = true
  saveFailed.value = false
  try {
    const success = customerContext ? await saveConfigurationPatch(customerContext, patch) : await updateRobotInfo(patch)
    saveFailed.value = success === false
    if (success !== false) settingVisible.value = false
  } catch { saveFailed.value = true }
  finally { saving.value = false }
}

const handleEdit = () => {
  updateRobotInfo({});
}
</script>

<style lang="less" scoped>
.modal-item-box.workspace-generation-form {
  display: flex; flex-direction: column; padding: 8px 19px 8px 8px; gap: 24px; color: #222; font-size: 14px; line-height: 22px;
  .modal-item .label { margin-bottom: 8px; color: #222; }
  :deep(.ant-radio-group) { display: flex; flex-wrap: wrap; gap: 24px; align-items: center; min-height: 26px; }
  :deep(.ant-radio-wrapper) { min-height: 24px; margin: 0; color: #222; line-height: 22px; }
  :deep(.ant-radio) { margin: 0 4px; }
  :deep(.ant-radio + span) { padding-inline-end: 0; }
  :deep(.ant-radio-inner) { border-color: #d8d8d8; }
  :deep(.ant-radio-checked .ant-radio-inner) { border-color: #146ef5; background: #146ef5; }
  :deep(.ant-radio-checked .ant-radio-inner::after) { width: 12px; height: 12px; margin: -6px 0 0 -6px; background: transparent url('@/assets/img/customer-service/basic-settings/recall/radio-selected.svg') center / 12px 12px no-repeat; transform: none; }
  :deep(.ant-select-single .ant-select-selector) { height: 32px; border-color: #d8d8d8; border-radius: 6px; padding: 0 12px; box-shadow: none; }
  :deep(textarea.ant-input) { min-height: 128px; padding: 8px 12px; border-color: #d8d8d8; border-radius: 6px; font-size: 14px; line-height: 22px; resize: vertical; }
}
</style>
