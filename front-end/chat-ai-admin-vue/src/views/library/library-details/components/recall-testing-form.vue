<style lang="less" scoped>
.recall-settings-box {
  padding: 0 16px;
  .alert-box {
    padding: 9px 16px;
    border: 1px solid #99bffd;
    background: #e9f1fe;
    border-radius: 2px;
    color: #3a4559;
    font-size: 14px;
    line-height: 22px;
    font-weight: 400;
    margin-bottom: 16px;
  }
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
  }
}

.model-icon {
  height: 18px;
}
</style>

<template>
  <div class="recall-settings-box">
    <div class="alert-box">
      用户提问时，会将召回的文档分段传递给大模型，大模型会根据传递的分段回答用户提问。您可以通过召回测试，输入问题查看匹配的文档分段和相似度
    </div>
    <div class="form-box">
      <div class="form-item">
        <div class="form-item-body">
          <a-textarea
            style="height: 76px"
            v-model:value="formState.question"
            placeholder="请输入测试问题"
          />
        </div>
      </div>
      <div class="form-item">
        <a-button :loading="loading" @click="handleRecallTest" type="primary" block>测试</a-button>
      </div>
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
              <a-slider class="custom-slider" v-model:value="formState.size" :min="1" :max="10" />
            </div>
            <div class="number-input-box">
              <a-input-number v-model:value="formState.size" :min="1" :max="10" />
            </div>
          </div>
        </div>
      </div>

      <div class="form-item">
        <div class="form-item-label">
          <span>相似度阈值&nbsp;</span>
          <a-tooltip>
            <template #title>召回时，只会召回相似度大于阈值的文本分段。</template>
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
              <a-input-number v-model:value="formState.similarity" :min="0" :max="1" :step="0.01" />
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
          />
        </div>
        <div class="form-item-body">
          <a-select
            v-model:value="formState.rerank_use_model"
            placeholder="请选择Rerank模型"
            @change="handleChangeRerankModel"
            style="width: 320px"
          >
            <a-select-opt-group v-for="item in rerankModelList" :key="item.id">
              <template #label>
                <span><img class="model-icon" :src="item.icon" alt="" /></span>
              </template>
              <a-select-option
                :value="val"
                :rerank_model_config_id="item.id"
                v-for="val in item.children"
                :key="val"
              >
                <span>{{ val }}</span>
              </a-select-option>
            </a-select-opt-group>
          </a-select>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getModelConfigOption } from '@/api/model/index'
import { reactive, ref, toRaw } from 'vue'
import { useRoute } from 'vue-router'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { libraryRecallTest } from '@/api/library'
const route = useRoute()
const loading = ref(false)

const emit = defineEmits(['save', 'load'])

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
  question: '',
  similarity: 0.5,
  size: 5,
  id: route.query.id
})

const handleChangeRerankModel = (val, option) => {
  formState.rerank_model_config_id = option.rerank_model_config_id
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
const handleRecallTest = () => {
  if (!formState.similarity) {
    return message.error('请输入相似度')
  }
  if (!formState.size) {
    return message.error('请输入分段数')
  }
  if (!formState.question) {
    return message.error('请输入测试问题')
  }
  let parmas = {
    id: formState.id,
    question: formState.question,
    size: formState.size,
    similarity: formState.similarity,
    search_type: formState.search_type
  }
  if (formState.rerank_status == 1) {
    parmas.rerank_model_config_id = formState.rerank_model_config_id
  }
  loading.value = true
  emit('load');
  libraryRecallTest(parmas)
    .then((res) => {
      emit('save', res.data)
    })
    .catch(() => {
      emit('save', [])
    })
    .finally(() => {
      loading.value = false
    })
}
// 获取rerank模型列表
const rerankModelList = ref([])

const getList = () => {
  getModelConfigOption({
    model_type: 'RERANK'
  }).then((res) => {
    let list = res.data || []

    rerankModelList.value = list.map((item) => {
      return {
        id: item.model_config.id,
        name: item.model_info.model_name,
        icon: item.model_info.model_icon_url,
        children: item.model_info.rerank_model_list
      }
    })
  })
}
getList()
defineExpose({
  open
})
</script>
