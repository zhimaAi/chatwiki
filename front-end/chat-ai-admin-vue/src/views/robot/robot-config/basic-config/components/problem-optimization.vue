<style lang="less" scoped>
.setting-box {
  position: relative;
  .robot-info-box {
    .robot-prompt {
      line-height: 22px;
      font-size: 14px;
      white-space: pre-wrap;
      word-break: break-all;
      color: #595959;
    }

    .dialog-bg-set {
      margin-top: 16px;
      cursor: pointer;
      display: inline-flex;
      padding: 5px 16px;
      justify-content: center;
      align-items: center;
      gap: 10px;
      border-radius: 6px;
      border: 1px solid #D9D9D9;
      background: #FFF;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
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
        color: #262626;
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
        margin-bottom: 4px;
      }
    }

    .thinking-switch {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-top: 12px;
    }
  }
}
</style>

<template>
  <edit-box class="setting-box" :compact="compact" :title="t('title_problem_optimization')" icon-name="problem-optimization">
    <template #extra>
      <a-tooltip v-if="compact" :title="t('tip_dialogue_background_help')">
        <a-button size="small" @click="openWorkspaceSettings">{{ t('btn_dialogue_background_setting') }}</a-button>
      </a-tooltip>
      <span v-else></span>
    </template>
    <div class="robot-info-box">
      <div class="robot-prompt" :title="t('desc_enable_question_optimize')">
        {{ t('desc_enable_question_optimize') }}
      </div>
      <a-tooltip v-if="!compact">
        <template #title>
          <span>{{ t('tip_dialogue_background_help') }}</span>
        </template>
        <div class="dialog-bg-set" @click="isBgSet = true">
          {{ t('btn_dialogue_background_setting') }}
        </div>
      </a-tooltip>
    </div>

    <template #switch><a-switch
      @change="handleEdit"
      :class="{ 'switch-item': !compact }"
      checkedValue="true"
      unCheckedValue="false"
      v-model:checked="robotInfo.enable_question_optimize"
      :checked-children="t('switch_on')"
      :un-checked-children="t('switch_off')"
    /></template>

    <SettingsDialog v-if="compact" :open="isBgSet" :width="640" :title="t('title_user_question_suggestion_setting')" :ok-text="commonT('confirm')" :cancel-text="commonT('cancel')" :confirm-loading="saving" :closable="!saving" :mask-closable="!saving" :keyboard="!saving" :cancel-button-props="{ disabled: saving }" @ok="saveWorkspaceSettings" @cancel="closeWorkspaceSettings">
      <template #closeIcon><img src="@/assets/img/customer-service/basic-settings/prompt/close.svg" alt="" /></template>
      <div class="workspace-optimization-form" :inert="saving ? '' : null">
        <div class="optimization-field">
          <span id="optimization-model-label">{{ t('label_model') }}</span>
          <a-radio-group v-model:value="modelStatus" aria-labelledby="optimization-model-label" :disabled="saving" @change="handleChangeModel">
            <a-radio value="0">{{ t('radio_follow_robot') }}</a-radio>
            <a-radio value="1">{{ t('radio_specify_model') }}</a-radio>
          </a-radio-group>
        </div>
        <div v-show="modelStatus === '1'" class="optimization-field">
          <span id="optimization-specified-model-label">{{ t('radio_specify_model') }}</span>
          <ModelSelect modelType="LLM" v-model:modeName="formState.optimize_question_use_model" v-model:modeId="formState.optimize_question_model_config_id" aria-labelledby="optimization-specified-model-label" @loaded="onVectorModelLoaded" />
        </div>
        <div v-if="showEnableThinking" class="optimization-field">
          <span id="optimization-thinking-label" :title="t('tooltip_deep_thinking')">{{ t('label_deep_thinking') }}</span>
          <div class="thinking-slider-row">
            <a-slider :value="effortPosition" :min="0" :max="REASONING_EFFORTS.length - 1" :step="1" :disabled="saving" :tip-formatter="value => REASONING_EFFORTS[value]" ariaLabelledByForHandle="optimization-thinking-label" @change="changeEffortPosition" />
            <ReasoningEffortSelect v-model:value="formState.optimize_question_reasoning_effort" @change="syncThinking" />
          </div>
        </div>
        <div class="optimization-field">
          <label for="optimization-background">{{ t('label_dialogue_background') }}</label>
          <a-textarea id="optimization-background" v-model:value="dialogueBackground" :placeholder="t('ph_dialogue_background')" />
        </div>
        <a-alert v-if="saveFailed" type="error" :message="commonT('save_failed')" show-icon />
      </div>
    </SettingsDialog>
    <div v-else class="modal-box" ref="modalBoxRef">
      <a-modal :getContainer="() => $refs.modalBoxRef" v-model:open="isBgSet" :width="472" :title="t('title_user_question_suggestion_setting')" @ok="handleSave">
        <div class="modal-item-box">
          <div class="modal-item">
            <div class="label">{{ t('label_model') }}</div>
            <a-radio-group v-model:value="modelStatus" @change="handleChangeModel">
              <a-radio value="0">{{ t('radio_follow_robot') }}</a-radio>
              <a-radio value="1">{{ t('radio_specify_model') }}</a-radio>
            </a-radio-group>
            <ModelSelect
              v-show="modelStatus == '1'"
              modelType="LLM"
              v-model:modeName="formState.optimize_question_use_model"
              v-model:modeId="formState.optimize_question_model_config_id"
              style="width: 100%; margin-top: 8px;"
              @loaded="onVectorModelLoaded"
            />
            <div class="thinking-switch" v-if="showEnableThinking">
              <span>{{ t('label_deep_thinking') }}</span>
              <a-tooltip :title="t('tooltip_deep_thinking')">
                <QuestionCircleOutlined />
              </a-tooltip>
              <ReasoningEffortSelect
                v-model:value="formState.optimize_question_reasoning_effort"
                @change="(effort) => (formState.optimize_question_enable_thinking = effort === 'none' ? 0 : 1)"
              />
            </div>
          </div>
          <div class="modal-item">
            <div class="label">{{ t('label_dialogue_background') }}</div>
            <a-textarea
              style="height: 120px"
              v-model:value="robotInfo.optimize_question_dialogue_background"
              :placeholder="t('ph_dialogue_background')" />
          </div>
        </div>
      </a-modal>
    </div>
  </edit-box>
</template>
<script setup>
import ReasoningEffortSelect from '@/components/reasoning-effort-select/index.vue'
import { REASONING_EFFORTS, toReasoningEffort } from '@/utils/reasoning-effort'
import { ref, reactive, inject, toRaw, nextTick, onMounted, computed } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import EditBox from './edit-box.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import SettingsDialog from '@/views/customer-service/basic-settings/components/settings-dialog.vue'
import { saveConfigurationPatch } from '@/views/customer-service/basic-settings/configuration'

const { t } = useI18n('views.robot.robot-config.basic-config.components.problem-optimization')
const { t: commonT } = useI18n('views.customer-service.index')
const { t: modelT } = useI18n('components.model-select.model-select')
const props = defineProps({ compact: { type: Boolean, default: false } })
const customerContext = inject('customerRobotContext', null)
const saving = ref(false)
const saveFailed = ref(false)
const dialogueBackground = ref('')

const isBgSet = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')

// optimize_question_model_config_id 接口返回0 modelStatus设置成0，有id modelStatus设置成1
const modelStatus = ref('0')
const formState = reactive({
  optimize_question_use_model: '',
  optimize_question_model_config_id: '',
  optimize_question_enable_thinking: 0,
  optimize_question_reasoning_effort: 'none'
})
const effortPosition = computed(() => {
  const index = REASONING_EFFORTS.indexOf(formState.optimize_question_reasoning_effort)
  return index < 0 ? 3 : index
})
function syncThinking(effort) {
  formState.optimize_question_enable_thinking = effort === 'none' ? 0 : 1
}
function changeEffortPosition(index) {
  formState.optimize_question_reasoning_effort = REASONING_EFFORTS[index]
  syncThinking(REASONING_EFFORTS[index])
}
function openWorkspaceSettings() {
  if (saving.value) return
  modelStatus.value = Number(robotInfo.optimize_question_model_config_id) > 0 ? '1' : '0'
  Object.assign(formState, {
    optimize_question_use_model: robotInfo.optimize_question_use_model || '',
    optimize_question_model_config_id: robotInfo.optimize_question_model_config_id || 0,
    optimize_question_enable_thinking: +robotInfo.optimize_question_enable_thinking || 0,
    optimize_question_reasoning_effort: toReasoningEffort(robotInfo.optimize_question_enable_thinking, robotInfo.optimize_question_reasoning_effort)
  })
  dialogueBackground.value = robotInfo.optimize_question_dialogue_background || ''
  saveFailed.value = false
  isBgSet.value = true
}
function closeWorkspaceSettings() {
  if (!saving.value) isBgSet.value = false
}
async function saveWorkspaceSettings() {
  if (saving.value) return
  if (modelStatus.value === '1' && (!Number(formState.optimize_question_model_config_id) || !formState.optimize_question_use_model)) {
    message.error(modelT('ph_select_model'))
    return
  }
  const patch = { ...toRaw(formState), optimize_question_dialogue_background: dialogueBackground.value }
  if (modelStatus.value === '0') { patch.optimize_question_model_config_id = 0; patch.optimize_question_use_model = '' }
  saving.value = true
  saveFailed.value = false
  try {
    const success = customerContext ? await saveConfigurationPatch(customerContext, patch) : await updateRobotInfo(patch)
    saveFailed.value = success === false
    if (success !== false) isBgSet.value = false
  } catch { saveFailed.value = true }
  finally { saving.value = false }
}

onMounted(() => {
  modelStatus.value = robotInfo.optimize_question_model_config_id != 0 ? '1' : '0'
  formState.optimize_question_use_model = robotInfo.optimize_question_use_model
  formState.optimize_question_model_config_id = robotInfo.optimize_question_model_config_id
  formState.optimize_question_enable_thinking = +robotInfo.optimize_question_enable_thinking || 0
  formState.optimize_question_reasoning_effort = toReasoningEffort(robotInfo.optimize_question_enable_thinking, robotInfo.optimize_question_reasoning_effort)
})

const vectorModelList = ref([])
const choosableThinking = ref({})
const onVectorModelLoaded = (list, thinkingMap) => {
  vectorModelList.value = list
  choosableThinking.value = thinkingMap || {}

  nextTick(() => {
    if (
      modelStatus.value === '1' &&
      (!formState.optimize_question_use_model || !Number(formState.optimize_question_model_config_id))
    ) {
      setDefaultModel()
    }
  })
}
const showEnableThinking = computed(() => {
  const modelConfigId =
    modelStatus.value === '1'
      ? formState.optimize_question_model_config_id
      : robotInfo.model_config_id
  const useModel =
    modelStatus.value === '1' ? formState.optimize_question_use_model : robotInfo.use_model
  return !!choosableThinking.value[modelConfigId + '#' + useModel]
})
const setDefaultModel = () => {
  if (vectorModelList.value.length > 0) {
    // 遍历查找chatwiki模型
    let modelConfig = null
    let model = null

    // 云版默认选中qwen-max
    for (let item of vectorModelList.value) {
      if (item.model_define === 'tongyi') {
        modelConfig = item
        for (let child of modelConfig.children) {
          if (child.name === 'qwen-max') {
            model = child
            break
          }
        }
        break
      }
    }

    if (!modelConfig || (props.compact && !model)) {
      modelConfig = vectorModelList.value[0]
      model = modelConfig.children[0]
    }

    if (modelConfig && model) {
      formState.optimize_question_use_model = model.name
      formState.optimize_question_model_config_id = model.model_config_id
    }
  }
}
const handleSave = () => {
  isBgSet.value = false
  updateRobotInfo({ ...toRaw(formState) })
}

const handleChangeModel = () => {
  if (modelStatus.value !== '1') {
    formState.optimize_question_use_model = ''
    formState.optimize_question_model_config_id = 0
  } else if (
    !formState.optimize_question_use_model ||
    !Number(formState.optimize_question_model_config_id)
  ) {
    setDefaultModel()
  }
}

const handleEdit = () => {
  updateRobotInfo({});
}
</script>

<style lang="less" scoped>
.workspace-optimization-form {
  display: flex; flex-direction: column; gap: 24px; padding: 8px 19px 8px 8px; color: #222; font-size: 14px; line-height: 22px;
  .optimization-field { display: flex; flex-direction: column; gap: 8px; }
  :deep(.ant-radio-group) { display: flex; flex-wrap: wrap; gap: 12px; min-height: 26px; align-items: center; }
  :deep(.ant-radio-wrapper) { min-height: 24px; margin: 0; color: #222; line-height: 22px; }
  :deep(.ant-radio) { margin: 0 4px; }
  :deep(.ant-radio + span) { padding-inline-end: 0; }
  :deep(.ant-radio-inner) { border-color: #d8d8d8; }
  :deep(.ant-radio-checked .ant-radio-inner) { border-color: #146ef5; background: #146ef5; }
  :deep(.ant-radio-checked .ant-radio-inner::after) { width: 12px; height: 12px; margin: -6px 0 0 -6px; background: transparent url('@/assets/img/customer-service/basic-settings/recall/radio-selected.svg') center / 12px 12px no-repeat; transform: none; }
  :deep(.ant-select-single .ant-select-selector) { height: 32px; border-color: #d8d8d8; border-radius: 6px; padding: 0 12px; box-shadow: none; }
  :deep(textarea.ant-input) { height: 172px; min-height: 172px; padding: 8px 12px; border-color: #d8d8d8; border-radius: 6px; font-size: 14px; line-height: 22px; resize: vertical; }
  .thinking-slider-row { display: grid; grid-template-columns: minmax(0, 1fr) 88px; align-items: center; gap: 16px; min-height: 32px; }
  .reasoning-effort-select { width: 88px; }
  :deep(.ant-slider) { height: 16px; margin: 8px; padding: 5px 0; }
  :deep(.ant-slider-rail), :deep(.ant-slider-track) { height: 6px; border-radius: 6px; }
  :deep(.ant-slider-rail) { background: #e8e8e8; }
  :deep(.ant-slider-track) { background: #146ef5; }
  :deep(.ant-slider-handle) { width: 16px; height: 16px; }
  :deep(.ant-slider-handle::after) { width: 12px; height: 12px; inset-inline-start: 2px; inset-block-start: 2px; box-shadow: 0 0 0 2px #146ef5; }
}
</style>
