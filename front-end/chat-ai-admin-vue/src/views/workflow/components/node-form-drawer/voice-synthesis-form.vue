<template>
  <NodeFormLayout class="voice-synthesis-form">
    <NodeFormHeader
      :title="node.node_name"
      :iconName="node.node_icon_name"
      :desc="t('desc_voice_synthesis')"
    />
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>{{ t('label_input') }}</div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_voice_model') }}</div>
        </div>
        <ModelSelect
          model-type="TTS"
          v-model:modeName="formState.use_model"
          v-model:modeId="formState.model_id"
          v-model:useConfigId="formState.model_config_id"
          @change="modelChange"
        />
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_text') }}</div>
          <div class="option-type">string</div>
        </div>
        <div>
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="formState.tag_map?.text || []"
            :defaultValue="formState.text"
            ref="atInputRef"
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
        <div class="desc">{{ t('msg_text_desc') }}</div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_select_voice') }}</div>
        </div>
        <div class="flex-between">
          <ZmRadioGroup v-model:value="formState.voice_setting.voice_id" :options="voiceOptions" @change="update"/>
          <a-tooltip placement="left" :title="t('label_advanced_settings')">
            <span class="setting-icon" @click="settingsShowChange"><SettingOutlined/></span>
          </a-tooltip>
        </div>
      </div>
      <div class="options-item is-required">
        <div class="options-item-tit">
          <div class="option-label">{{ t('label_voice_settings') }}</div>
          <div class="option-type">string</div>
        </div>
        <div class="flex-between">
          <a-select v-model:value="formState.voice_set_type" @change="voiceSetTypeChange">
            <a-select-option :value="1">{{ t('label_select_voice_type') }}</a-select-option>
            <a-select-option :value="2">{{ t('label_set_voice') }}</a-select-option>
          </a-select>
          <a-select
            v-if="formState.voice_set_type == 1"
            v-model:value="formState.voice_setting.voice_id"
            :open="false"
            @dropdownVisibleChange="showVoiceModal"
            style="width: 100%;"
          >
            <a-select-option value="-1">{{ t('ph_select') }}</a-select-option>
            <a-select-option v-for="item in showVoiceOpts" :value="item.value" :key="item.item">
              {{ item.label }}
            </a-select-option>
          </a-select>
          <AtInput
            v-else
            type="textarea"
            inputStyle="height: 33px"
            :options="variableOptions"
            :defaultSelectedList="formState.tag_map?.voice_setting_voice_id || []"
            :defaultValue="formState.voice_setting.voice_id"
            ref="atInputRef"
            :placeholder="t('ph_input_content')"
            @open="getValueVariableList"
            @change="changeVoiceId"
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

      <template v-if="settingsOpen">
        <div class="options-item">
          <div class="options-item-tit">
            <div class="option-label">{{ t('label_voice_settings_label') }}</div>
          </div>
          <div class="form-box">
            <div v-for="field in ['speed', 'vol', 'pitch']" :key="field" class="form-item">
              <div :class="['form-tit', {required: voiceSettingSliderConfig[field]}.required]">
                <a-tooltip :title="voiceSettingSliderConfig[field].tooltip">
                  {{ voiceSettingSliderConfig[field].label }}
                  <QuestionCircleOutlined/>
                  <span class="option-type">{{ voiceSettingSliderConfig[field].type }}</span>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <a-row>
                  <a-col :span="15">
                    <a-slider
                      v-model:value="formState.voice_setting[field]"
                      :step="voiceSettingSliderConfig[field].step"
                      :min="voiceSettingSliderConfig[field].min"
                      :max="voiceSettingSliderConfig[field].max"
                      :marks="voiceSettingSliderConfig[field].marks"
                    />
                  </a-col>
                  <a-col :span="7">
                    <a-input-number
                      v-model:value="formState.voice_setting[field]"
                      :placeholder="t('ph_select')"
                      :min="voiceSettingSliderConfig[field].min"
                      :max="voiceSettingSliderConfig[field].max"
                      class="ml16"
                    />
                  </a-col>
                </a-row>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit required">
                <a-tooltip :title="t('tip_emotion')">
                  {{ t('label_emotion') }}
                  <QuestionCircleOutlined/>
                  <span class="option-type">string</span>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <a-select v-model:value="formState.voice_setting.emotion" :placeholder="t('ph_select')" @change="update">
                  <a-select-option v-for="(txt, key) in emotionMap" :key="key">{{ txt }}</a-select-option>
                </a-select>
              </div>
            </div>
            <div class="form-item ">
              <div class="form-tit required">
                <a-tooltip :title="t('tip_text_normalization')">
                  {{ t('label_text_normalization') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <ZmRadioGroup v-model:value="formState.voice_setting.text_normalization" :options="defaultRadioOpt" @change="update"/>
              </div>
            </div>
          </div>
        </div>
        <div class="options-item">
          <div class="options-item-tit">
            <div class="option-label">{{ t('label_audio_settings') }}</div>
          </div>
          <div class="form-box">
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_sample_rate')">
                  {{ t('label_sample_rate') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
                <span class="option-type">Number</span>
              </div>
              <div class="form-cont">
                <a-select v-model:value="formState.audio_setting.sample_rate" :placeholder="t('ph_select')" @change="update">
                  <a-select-option v-for="item in audioSampleRates" :key="item">{{ item }}</a-select-option>
                </a-select>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_bitrate')">
                  {{ t('label_bitrate') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
                <span class="option-type">Number</span>
              </div>
              <div class="form-cont">
                <a-select v-model:value="formState.audio_setting.bitrate" :placeholder="t('ph_select')" @change="update">
                  <a-select-option v-for="item in audioBitrate" :key="item">{{ item }}</a-select-option>
                </a-select>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_format')">
                  {{ t('label_format') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
                <span class="option-type">string</span>
              </div>
              <div class="form-cont">
                <a-select v-model:value="formState.audio_setting.format" :placeholder="t('ph_select')" @change="update">
                  <a-select-option v-for="item in audioFormats" :key="item">{{ item }}</a-select-option>
                </a-select>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_channel')">
                  {{ t('label_channel') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
                <span class="option-type">Number</span>
              </div>
              <div class="form-cont">
                <a-select v-model:value="formState.audio_setting.channel" :placeholder="t('ph_select')" @change="update">
                  <a-select-option v-for="item in audioChannels" :key="item.value" :value="item.value">
                    {{ item.label }}
                  </a-select-option>
                </a-select>
              </div>
            </div>
          </div>
        </div>
        <div class="options-item">
          <div class="options-item-tit">
            <div class="option-label">{{ t('label_voice_effect_settings') }}</div>
          </div>
          <div class="form-box">
            <div v-for="field in ['pitch', 'intensity', 'timbre']" :key="field" class="form-item">
              <div :class="['form-tit', {required: voiceModifySliderConfig[field]}.required]">
                <a-tooltip :title="voiceModifySliderConfig[field].tooltip">
                  {{voiceModifySliderConfig[field].label}}
                  <QuestionCircleOutlined/>
                  <span class="option-type">{{ voiceModifySliderConfig[field].type }}</span>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <a-row>
                  <a-col :span="15">
                    <a-slider
                      v-model:value="formState.voice_modify[field]"
                      :step="voiceModifySliderConfig[field].step"
                      :min="voiceModifySliderConfig[field].min"
                      :max="voiceModifySliderConfig[field].max"
                      :marks="voiceModifySliderConfig[field].marks"
                    />
                  </a-col>
                  <a-col :span="7">
                    <a-input-number
                      v-model:value="formState.voice_modify[field]"
                      :placeholder="t('ph_select')"
                      :min="voiceModifySliderConfig[field].min"
                      :max="voiceModifySliderConfig[field].max"
                      class="ml16"
                    />
                  </a-col>
                </a-row>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_sound_effect')">
                  {{ t('label_sound_effect') }}
                  <QuestionCircleOutlined/>
                  <span class="option-type">string</span>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <a-select v-model:value="formState.voice_modify.sound_effects" :placeholder="t('ph_select')" @change="update">
                  <a-select-option value="">{{ t('opt_none') }}</a-select-option>
                  <a-select-option v-for="(txt, key) in effectMap" :key="key" :value="key">{{ txt }}</a-select-option>
                </a-select>
              </div>
            </div>
          </div>
        </div>
        <div class="options-item">
          <div class="options-item-tit">
            <div class="option-label">{{ t('label_other_settings') }}</div>
          </div>
          <div class="form-box">
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
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_output_format')">
                  {{ t('label_output_format') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <a-select v-model:value="formState.output_format" :placeholder="t('ph_select')" @change="update">
                  <a-select-option v-for="(txt, key) in formatMap" :key="key" :value="key">{{ txt }}</a-select-option>
                </a-select>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_constant_bitrate_control')">
                  {{ t('label_constant_bitrate_control') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <ZmRadioGroup v-model:value="formState.audio_setting.force_cbr" :options="defaultRadioOpt" @change="update"/>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_subtitle_service')">
                  {{ t('label_subtitle_service') }}
                  <QuestionCircleOutlined/>
                </a-tooltip>
              </div>
              <div class="form-cont">
                <ZmRadioGroup v-model:value="formState.subtitle_enable" :options="defaultRadioOpt" @change="update"/>
              </div>
            </div>
            <div class="form-item">
              <div class="form-tit">
                <a-tooltip :title="t('tip_add_audio_rhythm_mark')">
                  {{ t('label_add_audio_rhythm_mark') }}
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

    <SelectVoiceModal ref="voiceModalRef" :modelConfigId="formState.model_config_id" @change="voiceChange"/>
  </NodeFormLayout>
</template>

<script setup>
import {ref, onMounted, reactive, inject, computed} from 'vue';
import NodeFormLayout from "@/views/workflow/components/node-form-drawer/node-form-layout.vue";
import ModelSelect from "@/components/model-select/model-select.vue";
import NodeFormHeader from "@/views/workflow/components/node-form-drawer/node-form-header.vue";
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import ZmRadioGroup from "@/components/common/zm-radio-group.vue";
import SelectVoiceModal from "@/views/workflow/components/node-form-drawer/components/select-voice-modal.vue";
import {QuestionCircleOutlined, SettingOutlined} from '@ant-design/icons-vue';
import {message} from 'ant-design-vue';
import {
  audioBitrate,
  audioChannels,
  audioFormats,
  audioSampleRates,
  defaultRadioOpt,
  emotionMap,
  languageMap,
  formatMap,
  effectMap,
  voiceModifySliderConfig,
  voiceSettingSliderConfig,
  defaultVoiceOpts,
  voiceOutputObj,
} from "@/constants/voice.js";
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";
import {pluginOutputToTree} from "@/constants/plugin.js";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.workflow.components.node-form-drawer.voice-synthesis-form');

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
  use_model: '',
  model_config_id: 0,
  text: '',
  voice_set_type: 1,
  voice_setting: {
    voice_id: 'male-qn-jingying',
    voice_name: '',
    speed: 1,
    vol: 1,
    pitch: 0,
    emotion: 'happy',
    text_normalization: false,
  },
  audio_setting: {
    sample_rate: 32000,
    bitrate: 128000,
    format: 'mp3',
    channel: 1,
    force_cbr: false,
  },
  voice_modify: {
    pitch: 0,
    intensity: 0,
    timbre: 0,
    sound_effects: ''
  },
  language_boost: 'auto',
  output_format: 'url',
  subtitle_enable: false,
  aigc_watermark: false,
  tag_map: {}
})
const settingsOpen = ref(false)

const voiceOptions = computed(() => {
  const voiceId = formState.voice_setting.voice_id
  const inDefault = defaultVoiceOpts.some(i => i.value == voiceId)
  const customValue = inDefault ? '-1' : voiceSelectOpts.value?.[0]?.voice_id ?? voiceId
  return [
    ...defaultVoiceOpts,
    { label: t('opt_custom'), value: customValue }
  ]
})

const showVoiceOpts = computed(() => {
  return [
    ...defaultVoiceOpts,
    ...voiceSelectOpts.value.map(i => ({label: i.voice_name, value: i.voice_id}))
  ]
})

onMounted(() => {
  init()
})

function init() {
  getValueVariableList();
  outputData.value = pluginOutputToTree(voiceOutputObj)
  nodeParamsAssign()
}

function nodeParamsAssign() {
  let nodeParams = JSON.parse(props.node.node_params)
  let arg = nodeParams?.text_to_audio?.arguments || {}
  arg.tag_map = arg.tag_map ? arg.tag_map : {}
  Object.assign(formState, arg)
  if (formState.voice_setting?.voice_id && defaultVoiceOpts.findIndex(i => i.value == formState.voice_setting?.voice_id) < 0) {
    const {voice_id, voice_name} = formState.voice_setting
    voiceSelectOpts.value = [{voice_id, voice_name}]
  }
}

function getValueVariableList() {
  variableOptions.value = getNode().getAllParentVariable()
}

function changeValue(field, text, selectedList) {
  formState[field] = text
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
    return message.warning(t('msg_select_voice_model_first'))
  }
  voiceModalRef.value.show([formState.voice_setting.voice_id])
}

function modelChange() {
  formState.voice_setting.voice_id = 'male-qn-jingying'
  update()
}

function voiceSetTypeChange() {
  update()
}

function voiceChange(keys, rows) {
  formState.voice_setting.voice_id = keys[0]
  formState.voice_setting.voice_name = rows[0].voice_name
  voiceSelectOpts.value = rows
  update()
}

function settingsShowChange() {
  settingsOpen.value = !settingsOpen.value
}

function update() {
  let nodeParams = JSON.parse(props.node.node_params)
  formState.model_id = Number(formState.model_id || 0)
  Object.assign(nodeParams.text_to_audio, {
    voice_type: 'all',
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

  :deep(.mention-input-warpper) {
    height: 33px;
    word-break: break-all;
  }
}
.voice-synthesis-form{
  .options-item{
    margin-top: 24px !important;
  }
  .options-item-tit{
    margin-bottom: 4px;
  }
}

.setting-icon {
  color: #2475FC;
  padding: 4px 8px;
  border-radius: 24px;
  border: 1px solid #D9D9D9;
  background: #E5EFFF;
  cursor: pointer;
}

.form-box {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  gap: 8px;

  .form-item {
    display: flex;
    flex-direction: column;
    padding: 4px 0;

    .form-tit {
      display: flex;
      align-items: center;
      justify-content: start;
      flex-shrink: 0;
      margin-bottom: 8px;
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
</style>
