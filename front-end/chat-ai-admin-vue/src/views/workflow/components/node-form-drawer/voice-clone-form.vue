<template>
  <NodeFormLayout class="voice-clone-form">
    <NodeFormHeader
      :title="node.node_name"
      :iconName="node.node_icon_name"
      :desc="t('desc_voice_clone')"
    />
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>{{ t('label_input') }}</div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_clone_audio_file') }}</div>
          <div class="option-type">string</div>
        </div>
        <div>
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="formState.tag_map?.file_url || []"
            :defaultValue="formState.file_url"
            ref="atInputRef"
            :placeholder="t('ph_input_file_url')"
            @open="getValueVariableList"
            @change="(text, selectedList) => changeValue('file_url', text, selectedList)"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtInput>
        </div>
        <div class="desc">{{ t('msg_clone_audio') }}</div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_voice_id') }}</div>
          <div class="option-type">string</div>
        </div>
        <div>
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="formState.tag_map?.voice_id || []"
            :defaultValue="formState.voice_id"
            ref="atInputRef"
            :placeholder="t('ph_input_voice_id')"
            @open="getValueVariableList"
            @change="(text, selectedList) => changeValue('voice_id', text, selectedList)"
          >
            <template #option="{ label, payload }">
              <div class="field-list-item">
                <div class="field-label">{{ label }}</div>
                <div class="field-type">{{ payload.typ }}</div>
              </div>
            </template>
          </AtInput>
        </div>
        <div class="desc">{{ t('msg_voice_id_desc') }}</div>
      </div>
      <div class="options-item">
        <div class="options-item-tit flex-between">
          <div class="option-label">{{ t('label_clone_sample_audio') }}</div>
          <a @click="settingsShowChange">{{ t('label_advanced_settings') }} <DownOutlined v-if="!settingsOpen"/><UpOutlined v-else/></a>
        </div>
      </div>

      <template v-if="settingsOpen">
        <div class="options-item">
          <div class="form-box">
            <div class="form-item ">
              <div class="form-tit">
                <a-tooltip :title="t('tip_audio_file')">
                  {{ t('label_audio_file') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <AtInput
                  type="textarea"
                  inputStyle="height: 64px;"
                  :options="variableOptions"
                  :defaultSelectedList="formState.tag_map?.prompt_audio_url || []"
                  :defaultValue="formState.clone_prompt.prompt_audio_url"
                  ref="atInputRef"
                  :placeholder="t('ph_input_file_url')"
                  @open="getValueVariableList"
                  @change="(text, selectedList) => changeValue('prompt_audio_url', text, selectedList, formState.clone_prompt)"
                >
                  <template #option="{ label, payload }">
                    <div class="field-list-item">
                      <div class="field-label">{{ label }}</div>
                      <div class="field-type">{{ payload.typ }}</div>
                    </div>
                  </template>
                </AtInput>
              </div>
            </div>
            <div class="form-item ">
              <div class="form-tit">
                <a-tooltip :title="t('tip_text')">
                  {{ t('label_text') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <AtInput
                  type="textarea"
                  inputStyle="height: 64px;"
                  :options="variableOptions"
                  :defaultSelectedList="formState.tag_map?.prompt_text || []"
                  :defaultValue="formState.clone_prompt.prompt_text"
                  ref="atInputRef"
                  :placeholder="t('ph_input_sample_text')"
                  @open="getValueVariableList"
                  @change="(text, selectedList) => changeValue('prompt_text', text, selectedList, formState.clone_prompt)"
                >
                  <template #option="{ label, payload }">
                    <div class="field-list-item">
                      <div class="field-label">{{ label }}</div>
                      <div class="field-type">{{ payload.typ }}</div>
                    </div>
                  </template>
                </AtInput>
              </div>
            </div>
          </div>
        </div>
        <div class="options-item">
          <div class="options-item-tit">
            <div class="option-label">{{ t('label_clone_preview_params') }}</div>
          </div>
          <div class="form-box">
            <div class="form-item ">
              <div class="form-tit">
                <a-tooltip :title="t('tip_preview_text')">
                  {{ t('label_preview_text') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <AtInput
                  type="textarea"
                  inputStyle="height: 64px;"
                  :options="variableOptions"
                  :defaultSelectedList="formState.tag_map?.text || []"
                  :defaultValue="formState.text"
                  ref="atInputRef"
                  :placeholder="t('ph_input_preview_text')"
                  @open="getValueVariableList"
                  @change="(text, selectedList) => changeValue('text', text, selectedList)"
                >
                  <template #option="{ label, payload }">
                    <div class="field-list-item">
                      <div class="field-label">{{ label }}</div>
                      <div class="field-type">{{ payload.typ }}</div>
                    </div>
                  </template>
                </AtInput>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_small_language_recognition')">
                  {{ t('label_small_language_recognition') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <a-select v-model:value="formState.language_boost" :placeholder="t('ph_select')" @change="update">
                  <a-select-option v-for="(txt, key) in languageMap" :key="key" :value="key">{{ txt }}</a-select-option>
                </a-select>
              </div>
            </div>
          </div>
        </div>
        <div class="options-item">
          <div class="form-box">
            <div class="form-item">
              <div class="form-tit c262626">
                <a-tooltip :title="t('tip_preview_audio_model')">
                  {{ t('label_preview_audio_model') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <ModelSelect
                  style="width: 276px;"
                  model-type="TTS"
                  v-model:modeName="formState.model"
                  v-model:modeId="formState.model_id"
                  v-model:useConfigId="formState.model_config_id"
                  @change="update"
                />
              </div>
            </div>
          </div>
        </div>
        <div class="options-item">
          <div class="form-box">
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_noise_reduction')">
                  {{ t('label_noise_reduction') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <ZmRadioGroup v-model:value="formState.need_noise_reduction" :options="defaultRadioOpt" @change="update"/>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_volume_normalization')">
                  {{ t('label_volume_normalization') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <ZmRadioGroup v-model:value="formState.need_volume_normalization" :options="defaultRadioOpt" @change="update"/>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_audio_rhythm_mark')">
                  {{ t('label_audio_rhythm_mark') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <ZmRadioGroup v-model:value="formState.aigc_watermark" :options="defaultRadioOpt" @change="update"/>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>{{ t('label_output') }}</div>
      </div>
      <div class="options-item">
        <OutputFields :tree-data="outputData"/>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import {ref, onMounted, reactive, inject} from 'vue';
import NodeFormLayout from "@/views/workflow/components/node-form-drawer/node-form-layout.vue";
import ModelSelect from "@/components/model-select/model-select.vue";
import NodeFormHeader from "@/views/workflow/components/node-form-drawer/node-form-header.vue";
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import ZmRadioGroup from "@/components/common/zm-radio-group.vue";
import {QuestionCircleOutlined, DownOutlined, UpOutlined} from '@ant-design/icons-vue';
import {message} from 'ant-design-vue';
import {
  defaultRadioOpt,
  languageMap,
  voiceCloneOutputObj,
} from "@/constants/voice.js";
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";
import {pluginOutputToTree} from "@/constants/plugin.js";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.voice-clone-form')

const props = defineProps({
  lf: {
    type: Object,
    default: null
  },
  node: {
    type: Object,
    default: () => ({})
  }
})

const getNode = inject('getNode')
const setData = inject('setData')
const voiceModalRef = ref(null)
const variableOptions = ref([])
const outputData = ref([])
const voiceSelectOpts = ref([])
const formState = reactive({
  model_id: 0,
  model: '',
  model_config_id: 0,
  file_url: '',
  voice_id: '',
  clone_prompt: {
    prompt_audio_url: '',
    prompt_text: '',
    tag_map: {}
  },
  text: '',
  language_boost: 'auto',
  need_noise_reduction: false,
  need_volume_normalization: false,
  aigc_watermark: false,
  tag_map: {}
})
const settingsOpen = ref(false)

onMounted(() => {
  init()
})

function init() {
  getValueVariableList();
  outputData.value = pluginOutputToTree(voiceCloneOutputObj)
  nodeParamsAssign()
}

function nodeParamsAssign() {
  let nodeParams = JSON.parse(props.node.node_params)
  let arg = nodeParams?.voice_clone?.arguments || {}
  arg.tag_map = arg.tag_map ? arg.tag_map : {}
  Object.assign(formState, arg)
}

function getValueVariableList() {
  variableOptions.value = getNode().getAllParentVariable()
}

function changeValue(field, text, selectedList, _state=null) {
  _state = _state || formState
  _state[field] = text
  formState.tag_map[field] = selectedList
  update()
}

function changeVoiceId(text, selectedList) {
  formState.voice_setting.voice_id = text
  formState.tag_map.voice_setting_voice_id = selectedList
  update()
}

function showVoiceModal() {
  if (!formState.model_config_id) {
    return message.warning(t('msg_select_model_first'))
  }
  voiceModalRef.value.show([formState.voice_setting.voice_id])
}

function voiceSetTypeChange() {
  update()
}

function voiceChange(keys, rows) {
  formState.voice_setting.voice_id = keys[0]
  voiceSelectOpts.value = rows
  update()
}

function settingsShowChange() {
  settingsOpen.value = !settingsOpen.value
}

function update() {
  let nodeParams = JSON.parse(props.node.node_params)
  formState.model_id = Number(formState.model_id || 0)
  Object.assign(nodeParams.voice_clone, {
    model_config_id: formState.model_config_id,
    arguments: formState,
    output: outputData.value
  })
  setData({
    ...props.node,
    node_params: JSON.stringify(nodeParams),
  })
}
</script>

<style scoped lang="less">
@import "./components/node-options";

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.voice-clone-form{

}

.node-options {
  :deep(.mention-input-warpper) {
    height: 33px;
    word-break: break-all;
  }
}

.form-box {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  gap: 4px;

  .form-item {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 4px 0;

    .form-tit {
      display: flex;
      align-items: center;
      justify-content: start;
      flex-shrink: 0;
      word-break: break-all;
      color: #595959;

      .option-type {
        margin-left: 4px;
        display: inline-block;
        white-space: nowrap;
      }
    }

    .form-cont {
      flex: 1;

      :deep(.ant-select) {
        width: 100%;
      }

      :deep(.ant-slider-mark) {
        color: #8c8c8c;
        font-size: 12px;
        font-weight: 400;
        left: 6px;
        top: -16px;
      }

      :deep(.ant-slider-with-marks ) {
        margin-bottom: 0;
      }
    }
  }
}

.required::before {
  content: '*';
  color: #FB363F;
  display: inline-block;
  margin-right: 2px;
}

.ml16 {
  margin-left: 16px;
}
.c262626 {
  color: #262626 !important;
}
</style>
