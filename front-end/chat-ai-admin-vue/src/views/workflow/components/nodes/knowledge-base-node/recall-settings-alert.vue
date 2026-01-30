<style lang="less" scoped>
.recall-settings-box {
  margin-top: 24px;
  height: 600px;
  padding-right: 16px;
  overflow-y: auto;

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

      .number-slider-box {
        flex: 1;
      }

      .number-input-box {
        margin-left: 20px;
      }
    }

    .retrieval-mode-items {
      .retrieval-mode-item {
        position: relative;
        padding: 16px;
        margin-top: 8px;
        border-radius: 2px;
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
        }

        .recommendation-icon {
          margin-left: 4px;
          font-size: 36px;
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
</style>

<template>
  <a-modal width="600px" v-model:open="show" title="召回设置" @ok="handleSave" okText="确定" cancelText="取消">
    <div class="recall-settings-box">
      <div class="form-box">
        <div class="form-item is-required">
          <div class="form-item-label">
            <span>检索模式</span>
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
                  <span class="title-text">{{ item.title }}</span>
                  <svg-icon
                    class="recommendation-icon"
                    name="recommendation"
                    v-if="item.isRecommendation"
                  ></svg-icon>
                </div>

                <div class="retrieval-mode-desc">
                  {{ item.desc }}
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
            <span>Top K&nbsp;</span>
            <a-tooltip>
              <template #title
                >最多从知识库中召回分段数，最低为1，最高为10。召回分段数越多，消耗的token也会越多。</template
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

        <div class="form-item">
          <div class="form-item-label">
            <span>相似度阈值&nbsp;</span>
            <a-tooltip>
              <template #title>召回时，只会召回相似度大于阈值的文本分段。取值范围：0~1，阈值越大回答的越准确，建议不超过0.9</template>
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
            <span class="setting-title">召回相邻分段&nbsp;</span>
            <a-tooltip :overlayStyle="{ maxWidth: '350px' }">
              <template #title>
                <div style="font-size: 13px;">
                  <p>当在知识库检索到相关分段时，会根据配置拼接相关联的上下文。作为最终内容返回给大模型。</p>
                  <p>开启后，注意分段时不要设置分段重叠长度，否则可能影响最终效果。</p>
                  <p>父子分段类型的知识库，或者分段字数超过3000字时，不会做任何处理。</p>
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
                <span>拼接分段前</span>&nbsp;
                <a-select v-model:value="formState.recall_neighbor_before_num" style="width: 80px;">
                  <a-select-option :value="i - 1" v-for="i in 6" :key="i">{{ i - 1 }}</a-select-option>
                </a-select>
              </div>
              
              <div class="segment-input">
                <span>后</span>&nbsp;
                <a-select v-model:value="formState.recall_neighbor_after_num" style="width: 80px;">
                  <a-select-option :value="i - 1" v-for="i in 6" :key="i">{{ i - 1 }}</a-select-option>
                </a-select>
              </div>
              
              <div class="segment-text">个分段内容</div>
            </div>
          </div>
        </div>

        <div class="form-item">
          <div class="form-item-label">
            <span>Rerank模型</span>
            &nbsp;
            <a-switch
              :checkedValue="1"
              :unCheckedValue="0"
              v-model:checked="formState.rerank_status"
            />
          </div>
          <div class="form-item-body">
            <ModelSelect
              modelType="RERANK"
              v-model:modeName="formState.rerank_use_model"
              v-model:modeId="formState.rerank_model_config_id"
              style="width: 320px"
            />
          </div>
        </div>

        <div class="form-item">
          <div class="form-item-label">
            <a-tooltip title="元数据过滤是使用元数据属性（例如分组，知识创建时间等）来细化和控制系統内相关信息的检索过程。召回时仅会召回满足要求的知识。">
              <span>元数据过滤 <QuestionCircleOutlined/></span>
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
import { message } from 'ant-design-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import MetaFilterBox from "@/views/robot/robot-config/basic-config/components/meta-filter-box.vue";
import {getLibaryMetaSchemaList, getRobotMetaSchemaList} from "@/api/library/index.js";

const emit = defineEmits(['change'])

const retrievalModeList = ref([
  {
    iconName: 'search',
    title: '混合检索',
    value: 1,
    isRecommendation: true,
    desc: '同时执行全文检索和向量检索，使用RRF算法进行排序，从两中查询结果中选择更匹配用户问题的结果'
  },
  {
    iconName: 'file-search',
    title: '向量检索',
    value: 2,
    desc: '将用户提问转成向量之后与知识库分段匹配相似度，返回相似度高的结果'
  },
  {
    iconName: 'comment-search',
    title: '全文检索',
    value: 3,
    desc: '通过分词匹配文档中的词汇，返回包含这些词汇的文本片段'
  }
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

const show = ref(false)
const metaList = ref([])
const libraryIds = ref([])
const robotInfo = ref({})

const open = (data) => {
  formState.rerank_status = data.rerank_status || 0
  formState.rerank_use_model = data.rerank_use_model || undefined
  formState.rerank_model_config_id = data.rerank_model_config_id || ''
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
  let req
  if (libraryIds.value.length) {
    req = getLibaryMetaSchemaList({library_ids: libraryIds.value.toString()})
  } else {
    req = getRobotMetaSchemaList({id: robotInfo.value.id})
  }
  req.then(res => {
    metaList.value = res?.data || []
  })
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
    return message.error('请选择Rerank模型')
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
