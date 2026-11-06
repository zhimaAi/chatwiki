<style lang="less" scoped>
.recall-settings-box {
  &.async-submit { max-height: calc(100vh - 220px); overflow-y: auto; }
  margin-top: 24px;
  height: 650px;
  overflow-x: auto;
  padding-right: 16px;

  .form-box {
    .form-item {
      margin-bottom: 24px;
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

    .is-required {
      .form-item-label::before {
        content: '*';
        padding-right: 2px;
        font-size: 14px;
        color: #fb363f;
      }
    }

    .number-box {
      display: flex;
      align-items: center;
      width: 60%;
      .number-slider-box {
        flex: 1;
      }

      .number-input-box {
        margin-left: 20px;
      }
    }

    .retrieval-mode-items {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      .retrieval-mode-item {
        width: calc(50% - 8px);
        position: relative;
        padding: 16px;
        border-radius: 6px;
        border: 1px solid #d9d9d9;
        cursor: pointer;
      }

      .retrieval-mode-title {
        display: flex;
        align-items: center;
        line-height: 22px;
        margin-bottom: 4px;
        color: #262626;

        .title-icon {
          margin-right: 4px;
          font-size: 16px;
        }

        .title-text {
          font-size: 14px;
          font-weight: 600;
          margin-right: 4px;
        }
      }

      .retrieval-mode-desc {
        min-height: 44px;
        line-height: 22px;
        font-size: 14px;
        color: #595959;
      }

      .check-arrow {
        display: none;
      }
    }

    .retrieval-mode-item.active {
      border: 2px solid #2475fc;

      .check-arrow {
        position: absolute;
        display: block;
        right: -1px;
        bottom: -1px;
        width: 24px;
        height: 24px;
        font-size: 24px;
        color: #fff;
      }

      .retrieval-mode-title {
        color: #2475fc;
      }
    }

    .segment-controls{
      display: flex;
      align-items: center;
      gap: 4px;
      padding-top: 4px;
    }
  }
}

.model-icon {
  height: 18px;
}

.meta-box {
  border-radius: 6px;
  background: #F2F4F7;
  padding: 16px;
  margin-top: 8px;
}
.recall-settings-box.workspace-recall-settings {
  height: auto;
  max-height: none;
  margin-top: 0;
  padding: 8px 19px 8px 8px;
  overflow: visible;
  color: #222;

  .form-box {
    .form-item { margin-bottom: 24px; &:last-child { margin-bottom: 0; } }
    .form-item-label { display: flex; align-items: center; min-height: 22px; margin-bottom: 8px; color: #222; font-size: 14px; line-height: 22px; }
    .is-required .form-item-label::before { display: none; }
    .form-item-label :deep(.anticon-question-circle) { display: none; }
    .retrieval-mode-items { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); grid-template-rows: minmax(154px, auto) minmax(130px, auto); gap: 16px; }
    .retrieval-mode-item { width: auto; min-width: 0; padding: 16px; border: 1px solid #d8d8d8; border-radius: 8px; background: #fff; &.active { border: 1px solid #146ef5; background: #f1f6fe; } &:focus-visible { outline: 2px solid #146ef5; outline-offset: 2px; } }
    .retrieval-mode-title { gap: 8px; min-height: 24px; margin-bottom: 8px; color: #222; .title-text { margin: 0; font-weight: 400; line-height: 22px; } }
    .retrieval-mode-item.active .retrieval-mode-title { color: #222; }
    .retrieval-mode-desc { min-height: 0; margin-left: 32px; color: #5a5a5a; line-height: 22px; overflow-wrap: anywhere; }
    .mode-radio { position: relative; width: 24px; height: 24px; flex-shrink: 0; &::before { content: ''; position: absolute; top: 4px; left: 4px; width: 16px; height: 16px; box-sizing: border-box; border: 1px solid #d8d8d8; border-radius: 50%; background: #fff; } img { position: absolute; top: 6px; left: 6px; } }
    .active .mode-radio::before { border-color: #146ef5; background: #146ef5; }
    .recommendation { padding: 2px 6px; border-radius: 4px; background: #2475fc; color: #fff; font-size: 12px; line-height: 20px; }
    .number-box { gap: 16px; width: 100%; min-height: 32px; .number-input-box { flex-shrink: 0; margin-left: 0; } }
    .number-input-box :deep(.ant-input-number) { width: 88px; }
    .number-slider-box :deep(.ant-slider) { margin: 10px 8px; .ant-slider-rail { background: #e8e8e8; } .ant-slider-track { background: #146ef5; } .ant-slider-handle::after { box-shadow: 0 0 0 2px #146ef5; } }
    :deep(.ant-input-number) { height: 32px; border-radius: 6px; border-color: #d8d8d8; .ant-input-number-input { height: 30px; padding: 0 12px; } }
    .segment-controls { flex-wrap: wrap; gap: 8px; padding-top: 0; font-size: 14px; line-height: 22px; }
    .segment-input { display: inline-flex; align-items: center; white-space: nowrap; }
    .toggle-setting .form-item-label { margin-bottom: 0; }
    .toggle-setting.expanded .form-item-label { margin-bottom: 8px; }
    .form-item-label > :deep(.ant-switch) { min-width: 36px; width: 36px; height: 20px; margin-left: auto; line-height: 20px; .ant-switch-handle { width: 16px; height: 16px; top: 2px; } .ant-switch-inner { display: none; } &.ant-switch-checked { background: #146ef5; .ant-switch-handle { inset-inline-start: calc(100% - 18px); } } }
    .fulltext-options { display: flex; align-items: center; flex-wrap: wrap; gap: 12px; min-height: 26px; }
    .fulltext-options :deep(.ant-radio-wrapper) { margin-inline-end: 0; }
    .fulltext-options :deep(.ant-radio-inner) { border-color: #d8d8d8; }
    .fulltext-options :deep(.ant-radio-checked .ant-radio-inner) { border-color: #146ef5; background: #146ef5; &::after { top: 1px; left: 1px; width: 12px; height: 12px; margin: 0; border: 0; border-radius: 0; transform: none; background: url('@/assets/img/customer-service/basic-settings/recall/radio-selected.svg') center no-repeat; } }
  }
  .meta-box { max-width: 100%; overflow-x: auto; border-radius: 8px; background: #f5f5f5; }
  @media (max-width: 680px) { padding: 8px; .form-box .retrieval-mode-items { grid-template-columns: 1fr; grid-template-rows: auto; } .form-box .retrieval-mode-item { min-height: 130px; } }
}
</style>

<template>
  <component :is="isWorkspace ? SettingsDialog : AntModal" :width="isWorkspace ? 880 : '786px'" v-model:open="show" :title="t('title_recall_settings')" :ok-text="isWorkspace ? workspaceT('confirm') : undefined" :cancel-text="isWorkspace ? workspaceT('cancel') : undefined" :confirm-loading="saving" :closable="!saving" :mask-closable="!saving" :keyboard="!saving" :cancel-button-props="{ disabled: saving }" @ok="handleSave" @cancel="handleCancel">
    <template v-if="isWorkspace" #closeIcon><img src="@/assets/img/customer-service/basic-settings/recall/close.svg" alt="" /></template>
    <div class="recall-settings-box" :class="{ 'async-submit': !!props.submit, 'workspace-recall-settings': isWorkspace }" :inert="saving ? '' : null">
      <div class="form-box">
        <div class="form-item is-required">
          <div class="form-item-label">
            <span>{{ t('label_retrieval_mode') }}</span>
          </div>
          <div class="form-item-body">
            <div class="retrieval-mode-items" :role="isWorkspace ? 'radiogroup' : undefined" :aria-label="isWorkspace ? t('label_retrieval_mode') : undefined">
              <div
                class="retrieval-mode-item"
                :class="{ active: formState.search_type == item.value }"
                v-for="(item, index) in retrievalModeList"
                :key="item.value"
                :ref="element => modeElements[index] = element"
                :role="isWorkspace ? 'radio' : undefined"
                :aria-checked="isWorkspace ? formState.search_type == item.value : undefined"
                :tabindex="isWorkspace ? formState.search_type == item.value ? 0 : -1 : undefined"
                @keydown="event => handleModeKeydown(event, index)"
                @click="handleSelectRetrievalMode(item.value)"
              >
                <svg-icon
                  class="check-arrow"
                  name="check-arrow-filled"
                  v-if="!isWorkspace && formState.search_type == item.value"
                ></svg-icon>

                <div class="retrieval-mode-title">
                  <span v-if="isWorkspace" class="mode-radio" aria-hidden="true"><img v-if="formState.search_type == item.value" src="@/assets/img/customer-service/basic-settings/recall/radio-selected.svg" alt="" /></span>
                  <svg-icon v-else :name="item.iconName" class="title-icon"></svg-icon>
                  <span class="title-text">{{ t(item.title) }}</span>
                  <span v-if="isWorkspace && item.isRecommendation" class="recommendation">{{ tCommon('recommendation') }}</span>
                  <SvgTextTag :text="tCommon('recommendation')" v-if="!isWorkspace && item.isRecommendation" />
                </div>

                <div class="retrieval-mode-desc">
                  {{ t(item.desc) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="form-item weight-setting" v-if="formState.search_type == 1">
          <WeightSelect v-model:rrf_weight="formState.rrf_weight" />
        </div>

        <div class="form-item">
          <div class="form-item-label" :title="isWorkspace ? t('tooltip_top_k') : undefined">
            <span>{{ t('label_top_k') }}&nbsp;</span>
            <a-tooltip>
              <template #title
                >{{ t('tooltip_top_k') }}</template
              >
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
          </div>
          <div class="form-item-body">
            <div class="number-box">
              <div class="number-slider-box">
                <a-slider
                  class="custom-slider"
                  v-model:value="formState.top_k"
                  :min="1"
                  :max="500"
                />
              </div>
              <div class="number-input-box">
                <a-input-number v-model:value="formState.top_k" :min="1" :max="500" :precision="isWorkspace ? 0 : undefined" :controls="!isWorkspace" :aria-label="isWorkspace ? t('label_top_k') : undefined" />
              </div>
            </div>
          </div>
        </div>

        <div class="form-item" v-if="showSimilarity">
          <div class="form-item-label">
            <span>{{ t('label_similarity_threshold') }}&nbsp;</span>
            <a-tooltip>
              <template #title>{{ t('tooltip_similarity_threshold') }}</template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
          </div>
          <div class="form-item-body">
            <div class="number-box">
              <div class="number-slider-box">
                <a-slider
                  class="custom-slider"
                  v-model:value="formState.similarity"
                  :min="0"
                  :max="1"
                  :step="0.01"
                />
              </div>
              <div class="number-input-box">
                <a-input-number
                  v-model:value="formState.similarity"
                  :min="0"
                  :max="1"
                  :step="0.01"
                  :controls="!isWorkspace"
                  :aria-label="isWorkspace ? t('label_similarity_threshold') : undefined"
                />
              </div>
            </div>
          </div>
        </div>
        <div class="form-item" v-if="formState.search_type == 1 || formState.search_type == 3">
          <div class="form-item-label">
            {{ t('label_full_text_search_mode') }}
          </div>
          <div class="form-item-body">
            <a-radio-group class="fulltext-options" v-model:value="formState.library_search_type">
              <a-radio value="fullTextSearch">{{ t('full_text_search') }}</a-radio>
              <a-radio value="keywordSearch">{{ t('keyword_match') }}</a-radio>
            </a-radio-group>
          </div>
        </div>
        <div class="form-item">
          <div class="form-item-label">
            <span class="setting-title">{{ t('label_recall_neighbor_segments') }}&nbsp;</span>
            <a-tooltip :overlayStyle="{ maxWidth: '350px' }">
              <template #title>
                <div style="font-size: 13px;">
                  <p>{{ t('tooltip_recall_neighbor_desc_1') }}</p>
                  <p>{{ t('tooltip_recall_neighbor_desc_2') }}</p>
                  <p>{{ t('tooltip_recall_neighbor_desc_3') }}</p>
                  <p>{{ t('tooltip_recall_neighbor_desc_4') }}</p>
                  <p>{{ t('tooltip_recall_neighbor_desc_5') }}</p>
                </div>
              </template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
            &nbsp;
          </div>

          <div class="form-item-body">
            <div class="segment-controls">
              <div class="segment-input">
                <span>{{ t('label_top_k_ranking') }}</span>&nbsp;
                <a-input-number v-if="isWorkspace" v-model:value="formState.recall_neighbor_top_k" :min="0" :max="10" :precision="0" :controls="false" :aria-label="t('label_top_k_ranking')" style="width: 72px" />
                <a-select v-else v-model:value="formState.recall_neighbor_top_k" style="width: 80px;">
                  <a-select-option :value="i - 1" v-for="i in 11" :key="i">{{ i - 1 }}</a-select-option>
                </a-select>
                <span>&nbsp;{{ t('label_of_segments') }}，</span>
              </div>

              <div class="segment-input">
                <span>{{ t('label_auto_concat_before') }}</span>&nbsp;
                <a-input-number v-if="isWorkspace" v-model:value="formState.recall_neighbor_before_num" :min="0" :max="5" :precision="0" :controls="false" :aria-label="t('label_auto_concat_before')" style="width: 72px" />
                <a-select v-else v-model:value="formState.recall_neighbor_before_num" style="width: 80px;">
                  <a-select-option :value="i - 1" v-for="i in 6" :key="i">{{ i - 1 }}</a-select-option>
                </a-select>
              </div>

              <div class="segment-input">
                <span>{{ t('label_after') }}</span>&nbsp;
                <a-input-number v-if="isWorkspace" v-model:value="formState.recall_neighbor_after_num" :min="0" :max="5" :precision="0" :controls="false" :aria-label="t('label_after')" style="width: 72px" />
                <a-select v-else v-model:value="formState.recall_neighbor_after_num" style="width: 80px;">
                  <a-select-option :value="i - 1" v-for="i in 6" :key="i">{{ i - 1 }}</a-select-option>
                </a-select>
              </div>

              <div class="segment-text">{{ t('label_segments_content') }}</div>
            </div>
          </div>
        </div>

        <div class="form-item toggle-setting" :class="{ expanded: formState.rerank_status == 1 }">
          <div class="form-item-label">
            <span>{{ t('label_rerank_model') }}</span>
            &nbsp;
            <a-switch
              :checkedValue="1"
              :unCheckedValue="0"
              v-model:checked="formState.rerank_status"
              @change="handleRerankStatusChange"
            />
          </div>
          <div class="form-item-body" v-if="formState.rerank_status == 1">
            <ModelSelect
              modelType="RERANK"
              v-model:modeName="formState.rerank_use_model"
              v-model:modeId="formState.rerank_model_config_id"
              :style="{ width: isWorkspace ? '100%' : '60%' }"
            />
          </div>
        </div>

        <div class="form-item toggle-setting" :class="{ expanded: formState.meta_search_switch == 1 }">
          <div class="form-item-label">
            <a-tooltip :title="t('tooltip_metadata_filter')">
              <span>{{ t('label_metadata_filter') }} <QuestionCircleOutlined/></span>
            </a-tooltip>
            &nbsp;
            <a-switch
              :checkedValue="1"
              :unCheckedValue="0"
              v-model:checked="formState.meta_search_switch"
            />
          </div>
          <div class="form-item-body" v-if="formState.meta_search_switch == 1">
            <MetaFilterBox
              v-model:rule="formState.meta_search_condition_list"
              v-model:type="formState.meta_search_type"
              ref="metaFilterRef"
              class="meta-box"
              :variable-options="variableOptions"
              :meta-data="metaList"/>
          </div>
        </div>
      </div>
    </div>
  </component>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { reactive, ref, toRaw, computed } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import WeightSelect from '@/components/weight-select/index.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { message, Modal as AntModal } from 'ant-design-vue'
import {getRobotMetaSchemaList} from "@/api/library/index.js";
import MetaFilterBox from "@/views/robot/robot-config/basic-config/components/meta-filter-box.vue";
import SvgTextTag from '@/components/icons/SvgTextTag.vue'
import {getChatVariables} from "@/api/robot/index.js";
import dayjs from 'dayjs'
import SettingsDialog from '@/views/customer-service/basic-settings/components/settings-dialog.vue'

const { t } = useI18n('views.robot.robot-config.basic-config.components.associated-knowledge-base.recall-settings-alert')
const { t: tCommon } = useI18n('common')
const { t: workspaceT } = useI18n('views.customer-service.index')

const emit = defineEmits(['change'])
const props = defineProps({ submit: { type: Function, default: null }, variant: { type: String, default: 'legacy' } })
const isWorkspace = computed(() => props.variant === 'workspace')
const modeElements = ref([])
const saving = ref(false)

const retrievalModeList = ref([
  {
    iconName: 'mix-icon',
    title: 'title_mixed_retrieval',
    value: 1,
    isRecommendation: true,
    desc: 'desc_mixed_retrieval'
  },
  {
    iconName: 'vector-icon',
    title: 'title_vector_retrieval',
    value: 2,
    desc: 'desc_vector_retrieval'
  },
  {
    iconName: 'graph-icon',
    title: 'title_graph_retrieval',
    value: 4,
    desc: 'desc_graph_retrieval'
  },
  {
    iconName: 'search-check-icon',
    title: 'title_full_text_retrieval',
    value: 3,
    desc: 'desc_full_text_retrieval'
  },
])

const formState = reactive({
  rerank_status: 0,
  rerank_use_model: undefined,
  rerank_model_config_id: undefined,
  search_type: 1,
  meta_search_switch: 0,
  meta_search_type: 1,
  meta_search_condition_list: "",
  rrf_weight: {},
  recall_neighbor_switch: false,
  recall_neighbor_top_k: 5,
  recall_neighbor_before_num: 1,
  recall_neighbor_after_num: 1,
  library_search_type: 'fullTextSearch'
})

const showSimilarity = computed(() => {
  const searchType = Number(formState.search_type)
  const fullTextMode = formState.library_search_type || 'fullTextSearch'
  return searchType === 1 || searchType === 2 ||
    (searchType === 3 && fullTextMode === 'fullTextSearch')
})

const metaFilterRef = ref(null)
const show = ref(false)
const robotInfo = ref({})
const metaList = ref([])
const variableOptions = ref([])

const open = (data, r=null) => {
  robotInfo.value = r
  getMetaList()
  getChatVars()
  formState.rerank_status = data.rerank_status || 0
  formState.rerank_use_model = data.rerank_use_model || undefined
  formState.rerank_model_config_id = data.rerank_model_config_id > 0 ? data.rerank_model_config_id : ''
  formState.top_k = data.top_k
  formState.similarity = data.similarity
  formState.search_type = data.search_type
  formState.meta_search_switch = Number(data.meta_search_switch)
  formState.meta_search_type = Number(data.meta_search_type)
  formState.meta_search_condition_list = data.meta_search_condition_list
  formState.rrf_weight = data.rrf_weight
  formState.recall_neighbor_switch = data.recall_neighbor_switch
  formState.recall_neighbor_top_k = isWorkspace.value ? data.recall_neighbor_top_k ?? 5 : data.recall_neighbor_top_k
  formState.recall_neighbor_before_num = isWorkspace.value ? data.recall_neighbor_before_num ?? 1 : data.recall_neighbor_before_num
  formState.recall_neighbor_after_num = isWorkspace.value ? data.recall_neighbor_after_num ?? 1 : data.recall_neighbor_after_num
  formState.library_search_type = data.library_search_type || 'fullTextSearch'
  show.value = true
}


const getMetaList = () => {
  getRobotMetaSchemaList({id: robotInfo.value.id}).then(res => {
    metaList.value = res?.data || []
  })
}

const getChatVars = () => {
  getChatVariables({
    robot_key: robotInfo.value.robot_key
  }).then(res => {
    const typeMap = {
      input_string: 'string',
      input_number: 'number',
      input_array: 'array<object>'
    }
    const now = dayjs().unix()
    let _vars = res?.data || []
    if (_vars.length > 0) {
      variableOptions.value = [
        {
          label: t('label_robot_variables'),
          value: 'chat_variable',
          node_id: now,
          node_type: '1',
          typ: 'node',
          loop_parent_key: '',
          children: _vars.map(item => {
            const key = item.variable_key
            return {
              label: item.variable_name,
              value: `【chat_variable:${key}】`,
              original_value: `chat_variable:${key}`,
              node_id: now,
              node_name: t('label_robot_variables'),
              node_type: '1',
              text: item.variable_name,
              key: key,
              id: now,
              typ: typeMap[item.variable_type] || 'string',
              children: []
            }
          })
        }
      ]
    }
  })
}

const handleRerankStatusChange = (val) => {
  if(val == 0){
    formState.rerank_model_config_id = ''
    formState.rerank_use_model = '';
  }
}

const handleSelectRetrievalMode = (val) => {
  if (saving.value) return
  formState.search_type = val
}
const handleModeKeydown = (event, index) => {
  if (!isWorkspace.value || saving.value) return
  let next = { ArrowRight: (index + 1) % 4, ArrowLeft: (index + 3) % 4, ArrowDown: (index + 2) % 4, ArrowUp: (index + 2) % 4, Home: 0, End: 3 }[event.key]
  if (['Enter', ' '].includes(event.key)) next = index
  if (next === undefined) return
  event.preventDefault()
  handleSelectRetrievalMode(retrievalModeList.value[next].value)
  modeElements.value[next]?.focus()
}
const handleCancel = () => { if (!saving.value) show.value = false }

const checkRerank = () => {
  if (formState.rerank_status == 1 && !formState.rerank_model_config_id) {
    return true
  }

  return false
}

const handleSave = async () => {
  if (saving.value) return
  if (isWorkspace.value) {
    const valid = (value, min, max, integer = true) => value !== null && value !== '' && Number.isFinite(Number(value)) && Number(value) >= min && Number(value) <= max && (!integer || Number.isInteger(Number(value)))
    if (!valid(formState.top_k, 1, 500) ||
      (formState.search_type <= 2 && !valid(formState.similarity, 0, 1, false)) ||
      !valid(formState.recall_neighbor_top_k, 0, 10) || !valid(formState.recall_neighbor_before_num, 0, 5) || !valid(formState.recall_neighbor_after_num, 0, 5)) {
      return message.error(workspaceT('recall_invalid_parameters'))
    }
  }
  if (checkRerank()) {
    return message.error(t('msg_please_select_rerank_model'))
  }
  if (formState.meta_search_switch == 1 && !metaFilterRef.value?.verify()) {
    return
  }

  if (props.submit) {
    saving.value = true
    try {
      if (await props.submit(JSON.parse(JSON.stringify(toRaw(formState))))) show.value = false
    } catch (err) {
      message.error(err?.message || String(err))
    } finally { saving.value = false }
    return
  }
  show.value = false
  triggerChange()
}

const triggerChange = () => {
  emit('change', toRaw(formState))
}


defineExpose({
  open
})
</script>
