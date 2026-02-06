<template>
  <edit-box class="setting-box" :title="t('title_model_settings')" icon-name="moxingshezhi" v-model:isEdit="isEdit">
    <template #extra>
      <div class="actions-box">
        <a-button @click="handleEdit(true)" size="small">{{ t('btn_modify') }}</a-button>
      </div>
    </template>
    <div class="setting-info-block">
      <div class="set-item w-100">
        {{ t('label_llm_model') }}：
        <span>{{ getModelName }}</span>
      </div>
      <div class="set-item">
        {{ t('label_temperature') }}：
        <span>{{ formState.temperature }}</span>
      </div>
      <div class="set-item">
        {{ t('label_max_token') }}：
        <span>{{ formState.max_token }}</span>
      </div>
      <div class="set-item">
        {{ t('label_context_pair') }}：
        <span>{{ formState.context_pair }}</span>
      </div>
      <div class="set-item">
        {{ t('label_show_reasoning_process') }}：
        <span>{{ formState.think_switch == 1 ? t('state_on') : t('state_off') }}</span>
      </div>
    </div>
    <a-modal
      v-model:open="isEdit"
      :width="672"
      :title="t('title_model_settings')"
      @cancel="handleCancel"
      @ok="handleSave"
    >
      <div class="form-box">
        <a-row :gutter="[32, 24]">
          <a-col v-bind="grid">
            <div class="form-item">
              <div class="form-item-label mb12">
                <span>{{ t('label_llm_model') }}&nbsp;</span>
              </div>
              <div class="form-item-body">
                <!-- 自定义选择器 -->
                <ModelSelect
                  modelType="LLM"
                  v-model:modeName="formState.use_model"
                  v-model:modeId="formState.model_config_id"
                  @change="handleModelChange"
                  @loaded="onVectorModelLoaded"
                />
              </div>
            </div>
          </a-col>

          <a-col v-bind="grid">
            <div class="form-item" style="display: flex; align-items: center">
              <div class="form-item-label mb12">
                <span>{{ t('label_prompt_role') }}：&nbsp;</span>
              </div>
              <div>
                <a-radio-group v-model:value="formState.prompt_role_type">
                  <a-radio value="0">{{ t('option_system_role') }}</a-radio>
                  <a-radio value="1">{{ t('option_user_role') }}</a-radio>
                </a-radio-group>
              </div>
            </div>
          </a-col>

          <a-col v-bind="grid">
            <div class="form-item">
              <div class="form-item-label">
                <span>{{ t('label_temperature') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>{{ t('tooltip_temperature') }}</template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </div>
              <div class="form-item-body">
                <div class="number-box">
                  <div class="number-slider-box">
                    <a-slider
                      class="custom-slider"
                      v-model:value="formState.temperature"
                      :min="0"
                      :max="2"
                      :step="0.1"
                    />
                  </div>
                  <div class="number-input-box">
                    <a-input-number
                      v-model:value="formState.temperature"
                      :min="0"
                      :max="2"
                      :step="0.1"
                    />
                  </div>
                </div>
              </div>
            </div>
          </a-col>

          <a-col v-bind="grid">
            <div class="form-item">
              <div class="form-item-label">
                <span>{{ t('label_max_token') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>{{ t('tooltip_max_token') }}</template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </div>
              <div class="form-item-body">
                <div class="number-box">
                  <div class="number-slider-box">
                    <a-slider
                      class="custom-slider"
                      v-model:value="formState.max_token"
                      :min="0"
                      :max="100 * 1024"
                    />
                  </div>
                  <div class="number-input-box">
                    <a-input-number
                      v-model:value="formState.max_token"
                      :min="0"
                      :max="100 * 1024"
                    />
                  </div>
                </div>
              </div>
            </div>
          </a-col>

          <a-col v-bind="grid">
            <div class="form-item">
              <div class="form-item-label">
                <span>{{ t('label_context_pair') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>{{ t('tooltip_context_pair') }}</template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </div>
              <div class="form-item-body">
                <div class="number-box">
                  <div class="number-slider-box">
                    <a-slider
                      class="custom-slider"
                      v-model:value="formState.context_pair"
                      :min="0"
                      :max="50"
                    />
                  </div>
                  <div class="number-input-box">
                    <a-input-number v-model:value="formState.context_pair" :min="0" :max="50" />
                  </div>
                </div>
              </div>
            </div>
          </a-col>
          
          <a-col v-bind="grid">
            <div class="form-item justify-between">
              <div class="form-item-label">
                <span>{{ t('label_multimodal_input') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>
                    <span>{{ t('tooltip_multimodal_input') }}</span>
                  </template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </div>
              <CuTooltip :title="t('msg_model_not_support_multimodal')" :disabled="!disableMultimodalInput">
                <a-switch
                  v-model:checked="formState.question_multiple_switch"
                  :checkedValue="1"
                  :unCheckedValue="0"
                  :disabled="disableMultimodalInput"
                />
              </CuTooltip>
            </div>
          </a-col>

          <a-col v-bind="grid" v-if="show_enable_thinking">
            <div class="form-item justify-between">
              <div class="form-item-label">
                <span>{{ t('label_deep_thinking') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>
                    <span>{{ t('tooltip_deep_thinking') }}</span>
                  </template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </div>
              <a-switch
                v-model:checked="formState.enable_thinking"
                :checkedValue="1"
                :unCheckedValue="0"
              />
            </div>
          </a-col>

          <a-col v-bind="grid">
            <div class="form-item justify-between">
              <div class="form-item-label">
                <span>{{ t('label_show_reasoning_process') }}&nbsp;</span>
                <a-tooltip>
                  <template #title>
                    <span>{{ t('tooltip_show_reasoning_process') }}</span>
                  </template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </div>
              <a-switch
                v-model:checked="formState.think_switch"
                :checkedValue="1"
                :unCheckedValue="0"
              />
            </div>
          </a-col>
        </a-row>
      </div>
    </a-modal>
  </edit-box>
</template>

<script setup>
import {
  ref,
  reactive,
  inject,
  toRaw,
  watchEffect,
  computed
} from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import EditBox from './edit-box.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { getModelNameText } from '@/components/model-select/index.js'
import CuTooltip from '@/components/cu-tooltip/index.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.model-settings')

const grid = reactive({ sm: 24, md: 24, lg: 12, xl: 24, xxl: 24 })
// 获取LLM模型
const modelList = ref([])

const isEdit = ref(false)
const { robotInfo, updateRobotInfo, getRobot } = inject('robotInfo')
const formState = reactive({
  use_model: undefined,
  model_config_id: '',
  temperature: 0,
  max_token: 0,
  context_pair: 0,
  think_switch: 1,
  prompt_role_type: '0',
  enable_thinking: 0,
  question_multiple_switch: 1,
})

const disableMultimodalInput = ref(false)

// 处理选择事件
const handleModelChange = (val, data) => {
  const attr = data.attrs

  formState.question_multiple_switch = Number(attr.input_image) || 0;

  if (attr && attr.input_image == 1) {
    disableMultimodalInput.value = false
  } else {
    disableMultimodalInput.value = true
  }

  if (formState.use_model && formState.use_model.toLowerCase().includes('deepseek-r1')) {
    formState.prompt_role_type = '1'
  } else {
    formState.prompt_role_type = '0'
  }
}

const choosable_thinking = ref({})
const onVectorModelLoaded = (list, choosable_thinking_map) => {
  choosable_thinking.value = choosable_thinking_map
}

const show_enable_thinking = computed(() => {
  if (!formState.model_config_id) {
    return false
  }
  let key = formState.model_config_id + '#' + formState.use_model
  return choosable_thinking.value[key]
})

const handleSave = () => {
  let params = { ...toRaw(formState) }
  updateRobotInfo(params)
  isEdit.value = false
}

const handleCancel = () => {
  getRobot(robotInfo.id).then((res) => {
    formState.model_config_id = res.data.model_config_id
    formState.use_model = res.data.use_model
  })
}

const handleEdit = (val) => {
  formState.use_model = robotInfo.use_model
  formState.model_config_id = robotInfo.model_config_id
  formState.temperature = robotInfo.temperature
  formState.max_token = robotInfo.max_token
  formState.context_pair = robotInfo.context_pair
  formState.think_switch = Number(robotInfo.think_switch)
  formState.prompt_role_type = robotInfo.prompt_role_type
  formState.enable_thinking = robotInfo.enable_thinking
  formState.question_multiple_switch = Number(robotInfo.question_multiple_switch) || 0
  isEdit.value = val
}

watchEffect(() => {
  formState.use_model = robotInfo.use_model
  formState.model_config_id = robotInfo.model_config_id
  formState.temperature = robotInfo.temperature
  formState.max_token = robotInfo.max_token
  formState.context_pair = robotInfo.context_pair
  formState.think_switch = Number(robotInfo.think_switch)
  formState.prompt_role_type = robotInfo.prompt_role_type
  formState.enable_thinking = robotInfo.enable_thinking
  formState.question_multiple_switch = Number(robotInfo.question_multiple_switch) || 0
})

const getModelName = computed(() => {
  return getModelNameText(formState.model_config_id, formState.use_model)
})
</script>

<style lang="less" scoped>
.setting-box {
  ::v-deep(.edit-box-body) {
    padding: 0;
  }

  .actions-box {
    display: flex;
    align-items: center;
    line-height: 22px;
    font-size: 14px;
    color: #595959;

    .action-btn {
      cursor: pointer;
    }

    .save-btn {
      color: #2475fc;
    }

    .model-name {
      font-size: 14px;
      line-height: 22px;
      color: #8c8c8c;
    }
  }
}

.set-item-content {
  display: flex;
  align-items: center;
  gap: 6px;
}

.model-icon {
  height: 18px;
  width: 18px;
  object-fit: contain;
  vertical-align: middle;
}

/* 下拉选项对齐优化 */
.ant-select-item-option-content {
  display: flex;
  align-items: center;
  gap: 6px;
}

.setting-info-block {
  padding: 16px;
  padding-top: 0;
  display: flex;
  flex-wrap: wrap;
  gap: 12px 16px;
  color: #595959;
  line-height: 22px;
  .set-item {
    display: flex;
    align-items: center;
  }
  .w-100 {
    width: 100%;
  }
}

.form-box {
  padding: 16px 0 16px 0;

  .justify-between {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .form-item-label {
    line-height: 22px;
    margin-bottom: 4px;
    font-size: 14px;
    color: #262626;

    .question-icon {
      color: #8c8c8c;
    }
  }

  .number-box {
    display: flex;
    align-items: center;

    .number-slider-box {
      flex: 1;
    }

    .number-input-box {
      margin-left: 20px;
    }
  }
}
</style>
