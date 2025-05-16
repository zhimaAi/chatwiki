<template>
  <a-button @click="showDrawer" :icon="h(SettingOutlined)">搜索设置</a-button>
  <a-drawer v-model:open="open" class="custom-class" root-class-name="root-class-name" title="搜索设置" width="472"
    placement="right" @after-open-change="afterOpenChange">
    <div class="prompt-form">
      <div class="prompt-form-label">模型设置</div>

      <div class="prompt-form-item">
        <div class="prompt-form-item-label">
          <span style="color: red;">* </span><span>模型选择</span>
        </div>
        <ModelSelect
          modelType="LLM"
          v-model:modeName="formState.use_model"
          v-model:modeId="formState.model_config_id"
          style="width: 100%;"
          @loaded="onVectorModelLoaded"
          @change="handleChangeModel"
        />
      </div>

      <div class="prompt-form-item">
        <div class="prompt-form-item-label">
          <span style="color: red;">* </span><span>提示词</span>
        </div>
        <a-radio-group v-model:value="formState.prompt_type" @change="onSave">
          <a-radio value="0">默认提示词</a-radio>
          <a-radio value="1">自定义提示词</a-radio>
        </a-radio-group>
        <div class="prompt-form-item-content">
          <div class="prompt-form-item-tip" v-if="formState.prompt_type == '0'">将提交的内容进行智能总结,不要随意发挥</div>
          <a-textarea
            v-else
            @blur="onSave"
            :maxLength="500"
            style="height: 80px;"
            v-model:value="formState.prompt"
            placeholder="请输入自定义提示词"
          />
        </div>
      </div>

      <div class="prompt-form-item">
        <div class="prompt-form-item-label">
          <span>温度</span>
          <a-tooltip>
            <template #title>温度越低，回答越严谨。温度越高，回答越发散。</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="form-item-body">
          <div class="number-box">
            <div class="number-slider-box">
              <a-slider @blur="onSave" class="custom-slider" v-model:value="formState.temperature" :min="0" :max="2"
                :step="0.1" />
            </div>
            <div class="number-input-box">
              <a-input-number @blur="onSave" v-model:value="formState.temperature" :min="0" :max="2" :step="0.1" />
            </div>
          </div>
        </div>
      </div>

      <div class="prompt-form-item">
        <div class="prompt-form-item-label">
          <span>最大token</span>
          <a-tooltip>
            <template #title>
              <div>问题+答案的最大token数，如果出现回答被截断，可调高此值</div>
            </template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="form-item-body">
          <div class="number-box">
            <div class="number-slider-box">
              <a-slider @blur="onSave" class="custom-slider" v-model:value="formState.max_token" :min="0"
                :max="100 * 1024" />
            </div>
            <div class="number-input-box">
              <a-input-number @blur="onSave" v-model:value="formState.max_token" :min="0" :max="100 * 1024" />
            </div>
          </div>
        </div>
      </div>

      <!-- <div class="prompt-form-item">
        <div class="prompt-form-item-label">
          <span>上下文数量</span>
          <a-tooltip>
            <template #title
              >提示词中携带的历史聊天记录轮次。设置为0则不携带聊天记录。最多设置10轮。注意，携带的历史聊天记录越多，消耗的token相应也就越多。</template
            >
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
                :max="10"
              />
            </div>
            <div class="number-input-box">
              <a-input-number v-model:value="formState.context_pair" :min="0" :max="10" />
            </div>
          </div>
        </div>
      </div> -->
    </div>

    <div class="recall-settings-box">
      <div class="form-box prompt-form">
        <div class="prompt-form-label">召回设置</div>

        <div class="form-item is-required">
          <div class="form-item-label">
            <span>检索模式</span>
          </div>
          <div class="form-item-body">
            <div class="retrieval-mode-items">
              <div class="retrieval-mode-item" :class="{ active: formState.search_type == item.value }"
                v-for="item in retrievalModeList" :key="item.value" @click="handleSelectRetrievalMode(item.value)">
                <svg-icon class="check-arrow" name="check-arrow-filled"
                  v-if="formState.search_type == item.value"></svg-icon>

                <div class="retrieval-mode-title">
                  <svg-icon :name="item.iconName" class="title-icon"></svg-icon>
                  <span class="title-text">{{ item.title }}</span>
                  <svg-icon class="recommendation-icon" name="recommendation" v-if="item.isRecommendation"></svg-icon>
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
              <template #title>最多从知识库中召回分段数，最低为1，最高为10。召回分段数越多，消耗的token也会越多。</template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
          </div>
          <div class="form-item-body">
            <div class="number-box">
              <div class="number-slider-box">
                <a-slider @blur="onSave" class="custom-slider" v-model:value="formState.top_k" :min="1" :max="10" />
              </div>
              <div class="number-input-box">
                <a-input-number @blur="onSave" v-model:value="formState.top_k" :min="1" :max="10" />
              </div>
            </div>
          </div>
        </div>

        <div class="form-item" v-if="formState.search_type != 3 && formState.search_type != 4">
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
                <a-slider @blur="onSave" class="custom-slider" v-model:value="formState.similarity" :min="0" :max="1"
                  :step="0.01" />
              </div>
              <div class="number-input-box">
                <a-input-number @blur="onSave" v-model:value="formState.similarity" :min="0" :max="1" :step="0.01" />
              </div>
            </div>
          </div>
        </div>

        <div class="form-item">
          <div class="form-item-label">
            <span>Rerank模型&nbsp;</span>
            <a-tooltip>
              <template #title>召回时，只会召回相似度大于阈值的文本分段。取值范围：0~1，阈值越大回答的越准确，建议不超过0.9</template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>

            <a-switch @change="onSave" style="float: right" :checkedValue="1" :unCheckedValue="0"
              v-model:checked="formState.rerank_status" checked-children="开" un-checked-children="关" />
          </div>
          <div class="form-item-body" v-if="formState.rerank_status == 1">
            <a-select v-model:value="formState.rerank_use_model" placeholder="请选择Rerank模型"
              @change="handleChangeRerankModel" style="width: 100%">
              <a-select-opt-group v-for="item in rerankModelList" :key="item.id">
                <template #label>
                  <span><img class="model-icon" :src="item.icon" alt="" /></span>
                </template>
                <a-select-option :value="val" :rerank_model_config_id="item.id" v-for="val in item.children" :key="val">
                  <span>{{ val }}</span>
                </a-select-option>
              </a-select-opt-group>
            </a-select>
          </div>
        </div>
      </div>
    </div>
  </a-drawer>
</template>
<script lang="ts" setup>
import { getModelConfigOption } from '@/api/model/index'
import { ref, h, reactive, onMounted, toRaw, watch } from 'vue';
import { SettingOutlined } from '@ant-design/icons-vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { saveLibrarySearch } from '@/api/library'
import { message } from 'ant-design-vue'
import ModelSelect from '@/components/model-select/model-select.vue'

const props = defineProps({
  librarySearchData: {
    type: Object,
    default: null
  },
  defaultLibrarySearchData: {
    type: Object,
    default: null
  }
})

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
  }
])

const open = ref(false);

const formState = reactive({
  use_model: '',
  model_config_id: '',
  prompt_type: '0', // 提示词类型 0:默认 1:自定义
  prompt: '', // 提示词 500字符
  temperature: 0.5,
  max_token: 2000,
  context_pair: 0,
  search_type: 1,
  top_k: 5,
  size: 0,
  similarity: 0.4,
  rerank_status: 0,
  rerank_use_model: void 0,
  rerank_model_config_id: void 0,
})

// 获取rerank模型列表
const rerankModelList = ref([])

const getModelListRerank = () => {
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

const handleChangeModel = (val, option) => {
  formState.use_model = option.modelName
  formState.model_config_id = option.modelId
  onSave()
}

const vectorModelList = ref([])
const onVectorModelLoaded = (list) => {
  vectorModelList.value = list
  if (!props.librarySearchData?.model_config_id || !props.librarySearchData?.use_model) {
    if (vectorModelList.value.length > 0) {
      let modelConfig = vectorModelList.value[0]
      if (modelConfig) {
        let model = modelConfig.children[0]
        formState.use_model = model.name
        formState.model_config_id = model.model_config_id
      }
    }
  } else {
    formState.use_model = props.librarySearchData?.use_model + '$$'
    formState.model_config_id = props.librarySearchData?.model_config_id + '$$'
    formState.use_model = formState.use_model.split('$$')[0]
    formState.model_config_id = props.librarySearchData?.model_config_id.split('$$')[0]
  }
}

const handleSelectRetrievalMode = (val) => {
  formState.search_type = val
  onSave()
}

const handleChangeRerankModel = (val, option) => {
  formState.rerank_model_config_id = option.rerank_model_config_id
  onSave()
}

const updateLibraryInfo = (val) => {
  let newState = { ...val }
  newState.max_token = parseFloat(newState.max_token)
  newState.temperature = parseFloat(newState.temperature)
  newState.context_pair = parseFloat(newState.context_pair)
  newState.search_type = parseFloat(newState.search_type)
  newState.rerank_status = parseFloat(newState.rerank_status)
  newState.similarity = parseFloat(newState.similarity)
  newState.top_k = parseFloat(newState.size)
  
  if (newState.rerank_use_model === '') {
    // 这里是因为服务端可能会返回个空字符串，我这里改成undefined才有placeholder
    newState.rerank_use_model = void 0
  }

  Object.keys(newState).forEach((key) => {
    formState[key] = newState[key]
  })
}

const onSave = () => {
  let params = { ...formState }

  if (!params.prompt && formState.prompt_type == '1') {
    return message.error('请输入自定义提示词')
  }

  if (formState.search_type == 3 || formState.search_type == 4) {
    // 当选择“知识图谱检索”和“全文检索”类型时，不显示“相似度阈值”设置项。
    delete params.similarity
  }

  if (formState.prompt_type == '0') {
    delete params.prompt
  }

  params.size = params.top_k
  saveLibrarySearch(params).then((res) => {
    message.success('保存成功')
  })
}

const afterOpenChange = (bool) => {
  // console.log('open', bool);
};

const showDrawer = () => {
  open.value = true;
};

watch(
  () => props.librarySearchData,
  val => {
    if (props.librarySearchData?.id) {
      updateLibraryInfo({ ...toRaw(props.librarySearchData) })
    }
  },
);

onMounted(() => {
  // 获取llm
  getModelListRerank()
})
</script>

<style lang="less" scoped>
.prompt-form {

  .prompt-form-item {
    margin: 24px 0;

    .prompt-form-item-label {
      font-size: 14px;
      color: #262626;
    }
  }

  .prompt-form-label {
    color: #262626;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    margin-bottom: 16px;
  }

  .prompt-form-item-tip {
    line-height: 22px;
    font-size: 14px;
    color: #595959;
  }

  .prompt-form-item-content {
    margin-top: 8px;
  }

  .question-icon {
    color: #8c8c8c;
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

.model-icon {
  height: 18px;
}

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