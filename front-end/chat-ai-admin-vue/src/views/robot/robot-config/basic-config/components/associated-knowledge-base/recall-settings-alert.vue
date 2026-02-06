<style lang="less" scoped>
.recall-settings-box {
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
</style>

<template>
  <a-modal width="786px" v-model:open="show" :title="t('title_recall_settings')" @ok="handleSave">
    <div class="recall-settings-box">
      <div class="form-box">
        <div class="form-item is-required">
          <div class="form-item-label">
            <span>{{ t('label_retrieval_mode') }}</span>
          </div>
          <div class="form-item-body">
            <div class="retrieval-mode-items">
              <div
                class="retrieval-mode-item"
                :class="{ active: formState.search_type == item.value }"
                v-for="item in retrievalModeList"
                :key="item.value"
                @click="handleSelectRetrievalMode(item.value)"
              >
                <svg-icon
                  class="check-arrow"
                  name="check-arrow-filled"
                  v-if="formState.search_type == item.value"
                ></svg-icon>

                <div class="retrieval-mode-title">
                  <svg-icon :name="item.iconName" class="title-icon"></svg-icon>
                  <span class="title-text">{{ t(item.title) }}</span>
                  <img v-if="item.isRecommendation" style="width: 32px;" src="@/assets/svg/recommendation.svg" alt="">
                </div>

                <div class="retrieval-mode-desc">
                  {{ t(item.desc) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="form-item" v-if="formState.search_type == 1">
          <WeightSelect v-model:rrf_weight="formState.rrf_weight" />
        </div>

        <div class="form-item">
          <div class="form-item-label">
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
                <a-input-number v-model:value="formState.top_k" :min="1" :max="500" />
              </div>
            </div>
          </div>
        </div>

        <div class="form-item" v-if="formState.search_type <= 2">
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
                />
              </div>
            </div>
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
                </div>
              </template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
            &nbsp;
            <a-switch v-model:checked="formState.recall_neighbor_switch" />
          </div>
          
          <div class="form-item-body">
            <div class="segment-controls">
              <div class="segment-input">
                <span>{{ t('label_concat_before') }}</span>&nbsp;
                <a-select v-model:value="formState.recall_neighbor_before_num" style="width: 80px;">
                  <a-select-option :value="i - 1" v-for="i in 6" :key="i">{{ i - 1 }}</a-select-option>
                </a-select>
              </div>
              
              <div class="segment-input">
                <span>{{ t('label_after') }}</span>&nbsp;
                <a-select v-model:value="formState.recall_neighbor_after_num" style="width: 80px;">
                  <a-select-option :value="i - 1" v-for="i in 6" :key="i">{{ i - 1 }}</a-select-option>
                </a-select>
              </div>
              
              <div class="segment-text">{{ t('label_segments_content') }}</div>
            </div>
          </div>
        </div>

        <div class="form-item">
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
              style="width: 60%"
            />
          </div>
        </div>

        <div class="form-item">
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
              :meta-data="metaList"/>
          </div>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { reactive, ref, toRaw } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import WeightSelect from '@/components/weight-select/index.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { message } from 'ant-design-vue'
import {getRobotMetaSchemaList} from "@/api/library/index.js";
import MetaFilterBox from "@/views/robot/robot-config/basic-config/components/meta-filter-box.vue";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.associated-knowledge-base.recall-settings-alert')

const emit = defineEmits(['change'])

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
  recall_neighbor_before_num: 1,
  recall_neighbor_after_num: 1,
})

const metaFilterRef = ref(null)
const show = ref(false)
const robotInfo = ref({})
const metaList = ref([])

const open = (data, r=null) => {
  robotInfo.value = r
  getMetaList()
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
  formState.recall_neighbor_before_num = data.recall_neighbor_before_num
  formState.recall_neighbor_after_num = data.recall_neighbor_after_num

  show.value = true
}


const getMetaList = () => {
  getRobotMetaSchemaList({id: robotInfo.value.id}).then(res => {
    metaList.value = res?.data || []
  })
}

const handleRerankStatusChange = (val) => {
  if(val == 0){
    formState.rerank_model_config_id = ''
    formState.rerank_use_model = '';
  }
}

const handleSelectRetrievalMode = (val) => {
  formState.search_type = val
}

const checkRerank = () => {
  if (formState.rerank_status == 1 && !formState.rerank_model_config_id) {
    return true
  }

  return false
}

const handleSave = () => {
  if (checkRerank()) {
    return message.error(t('msg_please_select_rerank_model'))
  }
  if (formState.meta_search_switch == 1 && !metaFilterRef.value.verify()) {
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
