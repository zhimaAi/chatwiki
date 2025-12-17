<style lang="less" scoped>
.recall-settings-box {
  margin-top: 24px;

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
  }
}

.model-icon {
  height: 18px;
}
</style>

<template>
  <a-modal width="746px" v-model:open="show" title="召回设置" @ok="handleSave">
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
                  <img v-if="item.isRecommendation" style="width: 32px;" src="@/assets/svg/recommendation.svg" alt="">
                </div>

                <div class="retrieval-mode-desc">
                  {{ item.desc }}
                </div>
              </div>
            </div>
          </div>
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
                  :max="10"
                />
              </div>
              <div class="number-input-box">
                <a-input-number v-model:value="formState.top_k" :min="1" :max="10" />
              </div>
            </div>
          </div>
        </div>

        <div class="form-item" v-if="formState.search_type <= 2">
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
            <span>Rerank模型</span>
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
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { reactive, ref, toRaw } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { message } from 'ant-design-vue'

const emit = defineEmits(['change'])

const retrievalModeList = ref([
  {
    iconName: 'mix-icon',
    title: '混合检索',
    value: 1,
    isRecommendation: true,
    desc: '同时执行三种检索模式，使用RRF算法进行排序，从三种查询结果中选择更匹配用户问题的结果。混合检索兼顾语义相似性与逻辑关联性，通过互补优势提升检索的准确性和生成结果的可信度'
  },
  {
    iconName: 'vector-icon',
    title: '向量检索',
    value: 2,
    desc: '将用户提问转成向量之后与知识库分段匹配相似度，返回相似度高的结果。向量检索擅长语义相似性匹配和大规模非结构化数据处理，但缺乏可解释性和精准关系验证'
  },
  {
    iconName: 'graph-icon',
    title: '知识图谱检索',
    value: 4,
    desc: '通过关系推理，检索出与用户问题相关联的知识。知识图谱检索擅长精准的实体关系推理和逻辑验证，但对非结构化文本和语义模糊查询支持较弱'
  },
  {
    iconName: 'search-check-icon',
    title: '全文检索',
    value: 3,
    desc: '通过分词匹配文档中的词汇，返回包含这些词汇的文本片段'
  },
])

const formState = reactive({
  rerank_status: 0,
  rerank_use_model: undefined,
  rerank_model_config_id: undefined,
  search_type: 1
})

const show = ref(false)

const open = (data) => {

  formState.rerank_status = data.rerank_status || 0
  formState.rerank_use_model = data.rerank_use_model || undefined
  formState.rerank_model_config_id = data.rerank_model_config_id > 0 ? data.rerank_model_config_id : ''
  formState.top_k = data.top_k
  formState.similarity = data.similarity
  formState.search_type = data.search_type
  show.value = true
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
