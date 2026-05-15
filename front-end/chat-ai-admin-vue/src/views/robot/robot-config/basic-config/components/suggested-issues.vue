<style lang="less" scoped>
.setting-box {
  position: relative;
  .robot-info-box {
    .robot-prompt {
      display: flex;
      align-items: center;
      line-height: 22px;
      font-size: 14px;
      white-space: pre-wrap;
      word-break: break-all;
      color: #595959;
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
  .switch-item {
    position: absolute;
    right: 16px;
    top: calc(50% - 8px);
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
}
</style>

<template>
  <edit-box class="setting-box" :title="t('title_suggested_questions')" icon-name="suggested-issues">
    <template #extra>
      <span></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt">
        {{ t('msg_guide_description', { num: robotInfo.question_guide_num }) }}
        <a-input-number
          v-model:value="robotInfo.question_guide_num"
          style="width: 80px"
          :precision="0"
          :min="1"
          :max="10"
          @blur="handleEdit"
        />
      </div>
      <div class="guide-action-row">
        <a-button @click="openSettingModal">{{ t('btn_generation_setting') }}</a-button>
        <span>{{ t('label_current_mode') }}{{ currentModeText }}</span>
      </div>
    </div>

    <a-switch
      @change="handleEdit"
      class="switch-item"
      checkedValue="true"
      unCheckedValue="false"
      v-model:checked="robotInfo.enable_question_guide"
      :checked-children="t('btn_on')"
      :un-checked-children="t('btn_off')"
    />

    <div class="modal-box" ref="modalBoxRef">
      <a-modal
        :getContainer="() => $refs.modalBoxRef"
        v-model:open="settingVisible"
        :width="520"
        :title="t('title_generation_mode')"
        @ok="handleSaveSetting"
      >
        <div class="modal-item-box">
          <div class="modal-item">
            <div class="label">{{ t('label_generation_mode') }}</div>
            <a-radio-group v-model:value="formState.question_guide_mode">
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
              show-count
              style="height: 160px"
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
        </div>
      </a-modal>
    </div>
  </edit-box>
</template>
<script setup>
import { ref, reactive, inject, computed } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import EditBox from './edit-box.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.suggested-issues')

const { robotInfo, updateRobotInfo } = inject('robotInfo')
const props = defineProps({
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

const workflowOptions = computed(() => {
  return props.robotList
    .filter((item) => item.has_published == 1)
    .map((item) => ({
      label: item.robot_name,
      value: item.robot_key
    }))
})

const openSettingModal = () => {
  formState.question_guide_mode = Number(robotInfo.question_guide_mode) || 0
  formState.question_guide_prompt = robotInfo.question_guide_prompt || ''
  formState.question_guide_workflow_key = robotInfo.question_guide_workflow_key || ''
  settingVisible.value = true
}

const handleSaveSetting = () => {
  if (formState.question_guide_mode == 1 && !formState.question_guide_prompt.trim()) {
    return message.error(t('msg_prompt_required'))
  }

  if (formState.question_guide_mode == 2 && !formState.question_guide_workflow_key) {
    return message.error(t('msg_workflow_required'))
  }

  updateRobotInfo({
    question_guide_mode: formState.question_guide_mode,
    question_guide_prompt: formState.question_guide_mode == 1 ? formState.question_guide_prompt.trim() : '',
    question_guide_workflow_key: formState.question_guide_mode == 2 ? formState.question_guide_workflow_key : ''
  })
  settingVisible.value = false
}

const handleEdit = () => {
  updateRobotInfo({});
}
</script>
