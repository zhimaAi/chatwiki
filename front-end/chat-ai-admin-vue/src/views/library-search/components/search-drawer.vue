<template>
  <a-button v-if="searchSets" @click="showDrawer" :icon="h(SettingOutlined)">{{ t('search_settings') }}</a-button>
  <a-drawer
    v-model:open="open"
    class="custom-class"
    root-class-name="root-class-name"
    :title="t('search_settings')"
    width="472"
    placement="right"
    @after-open-change="afterOpenChange"
  >
    <div class="prompt-form">
      <div class="prompt-form-label">{{ t('model_settings') }}</div>

      <div class="form-item">
        <div class="form-item-label">
          <a-flex align="center">
            <span>{{ t('ai_intelligent_summary') }}</span>
            <a-tooltip>
              <template #title>{{ t('ai_summary_tooltip') }}</template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
            <a-switch
              style="margin-left: 24px"
              @change="onSave"
              v-model:checked="formState.summary_switch"
              :checkedValue="1"
              :unCheckedValue="0"
              :checkedChildren="t('on')"
              :unCheckedChildren="t('off')"
            />
          </a-flex>
        </div>
      </div>
      <div v-show="formState.summary_switch == 1">
        <div class="prompt-form-item">
          <div class="prompt-form-item-label">
            <span style="color: red">* </span><span>{{ t('model_selection') }}</span>
          </div>
          <ModelSelect
            modelType="LLM"
            v-model:modeName="formState.use_model"
            v-model:modeId="formState.model_config_id"
            style="width: 100%"
            @loaded="onVectorModelLoaded"
            @change="handleChangeModel"
          />
        </div>

        <div class="prompt-form-item">
          <div class="prompt-form-item-label">
            <span style="color: red">* </span><span>{{ t('prompt') }}</span>
          </div>
          <a-radio-group v-model:value="formState.prompt_type" @change="handlePromptTypeChange">
            <a-radio value="0">{{ t('default_prompt') }}</a-radio>
            <a-radio value="1">{{ t('custom_prompt') }}</a-radio>
          </a-radio-group>
          <div class="prompt-form-item-content">
            <div class="prompt-form-item-tip" v-if="formState.prompt_type == '0'">
              {{ t('default_prompt_tip') }}
            </div>
            <a-textarea
              v-else
              @blur="onSave"
              :maxLength="500"
              style="height: 80px"
              v-model:value="formState.prompt"
              :placeholder="t('custom_prompt_placeholder')"
            />
          </div>
        </div>

        <div class="prompt-form-item">
          <div class="prompt-form-item-label">
            <span>{{ t('temperature') }}</span>
            <a-tooltip>
              <template #title>{{ t('temperature_tooltip') }}</template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
          </div>
          <div class="form-item-body">
            <div class="number-box">
              <div class="number-slider-box">
                <a-slider
                  @blur="onSave"
                  class="custom-slider"
                  v-model:value="formState.temperature"
                  :min="0"
                  :max="2"
                  :step="0.1"
                />
              </div>
              <div class="number-input-box">
                <a-input-number
                  @blur="onSave"
                  v-model:value="formState.temperature"
                  :min="0"
                  :max="2"
                  :step="0.1"
                />
              </div>
            </div>
          </div>
        </div>

        <div class="prompt-form-item">
          <div class="prompt-form-item-label">
            <span>{{ t('max_token') }}</span>
            <a-tooltip>
              <template #title>
                <div>{{ t('max_token_tooltip') }}</div>
              </template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
          </div>
          <div class="form-item-body">
            <div class="number-box">
              <div class="number-slider-box">
                <a-slider
                  @blur="onSave"
                  class="custom-slider"
                  v-model:value="formState.max_token"
                  :min="0"
                  :max="100 * 1024"
                />
              </div>
              <div class="number-input-box">
                <a-input-number
                  @blur="onSave"
                  v-model:value="formState.max_token"
                  :min="0"
                  :max="100 * 1024"
                />
              </div>
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
        <div class="prompt-form-label">{{ t('recall_settings') }}</div>

        <div class="form-item is-required">
          <div class="form-item-label">
            <span>{{ t('retrieval_mode') }}</span>
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
          <WeightSelect @save="debouncedHandleWeightChange" v-model:rrf_weight="formState.rrf_weight" />
        </div>

        <div class="form-item">
          <div class="form-item-label">
            <span>{{ t('top_k') }}&nbsp;</span>
            <a-tooltip>
              <template #title>{{ t('top_k_tooltip') }}</template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
          </div>
          <div class="form-item-body">
            <div class="number-box">
              <div class="number-slider-box">
                <a-slider
                  @blur="onSave"
                  class="custom-slider"
                  v-model:value="formState.top_k"
                  :min="1"
                  :max="500"
                />
              </div>
              <div class="number-input-box">
                <a-input-number @blur="onSave" v-model:value="formState.top_k" :min="1" :max="500" />
              </div>
            </div>
          </div>
        </div>

        <div class="form-item" v-if="formState.search_type != 3 && formState.search_type != 4">
          <div class="form-item-label">
            <span>{{ t('similarity_threshold') }}&nbsp;</span>
            <a-tooltip>
              <template #title>{{ t('similarity_threshold_tooltip') }}</template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
          </div>
          <div class="form-item-body">
            <div class="number-box">
              <div class="number-slider-box">
                <a-slider
                  @blur="onSave"
                  class="custom-slider"
                  v-model:value="formState.similarity"
                  :min="0"
                  :max="1"
                  :step="0.01"
                />
              </div>
              <div class="number-input-box">
                <a-input-number
                  @blur="onSave"
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
            <span>{{ t('rerank_model') }}&nbsp;</span>
            <a-tooltip>
              <template #title>{{ t('rerank_model_tooltip') }}</template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>

            <a-switch
              @change="onSave"
              style="float: right"
              :checkedValue="1"
              :unCheckedValue="0"
              v-model:checked="formState.rerank_status"
              :checkedChildren="t('on')"
              :unCheckedChildren="t('off')"
            />
          </div>
          <div class="form-item-body" v-if="formState.rerank_status == 1">
            <ModelSelect
              modelType="RERANK"
              v-model:modeName="formState.rerank_use_model"
              v-model:modeId="formState.rerank_model_config_id"
              style="width: 100%"
              :placeholder="t('select_rerank_model')"
              @change="onSave"
            />
          </div>
        </div>
      </div>
    </div>
    <a-modal
      v-model:open="modalOpen"
      @cancel="handleClose"
      :title="t('custom_prompt_modal_title')"
      @ok="handleOk"
    >
      <a-textarea
        v-model:value="promptVal"
        :placeholder="t('please_input')"
        style="min-height: 300px"
        allow-clear
      />
    </a-modal>
  </a-drawer>
</template>
<script lang="ts" setup>
import { ref, h, reactive, onMounted, computed } from 'vue'
import { SettingOutlined } from '@ant-design/icons-vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { saveLibrarySearch } from '@/api/library'
import { message } from 'ant-design-vue'
import WeightSelect from '@/components/weight-select/index.vue'
import { useSearchLiraryStore } from '@/stores/modules/search-lirary'
import ModelSelect from '@/components/model-select/model-select.vue'
import { usePermissionStore } from '@/stores/modules/permission'
import { useI18n } from '@/hooks/web/useI18n'

let { role_permission, role_type } = usePermissionStore()
const { t } = useI18n('views.library-search.components.search-drawer')

const searchSets = computed(() => role_type == 1 || role_permission.includes('SearchSets'))

const searchLiraryStore = useSearchLiraryStore()

const { getLibrarySearchFn } = searchLiraryStore

const retrievalModeList = ref([
  {
    iconName: 'mix-icon',
    title: t('hybrid_retrieval'),
    value: 1,
    isRecommendation: true,
    desc: t('hybrid_retrieval_desc')
  },
  {
    iconName: 'vector-icon',
    title: t('vector_retrieval'),
    value: 2,
    desc: t('vector_retrieval_desc')
  },
  {
    iconName: 'graph-icon',
    title: t('graph_retrieval'),
    value: 4,
    desc: t('graph_retrieval_desc')
  },
  {
    iconName: 'search-check-icon',
    title: t('full_text_retrieval'),
    value: 3,
    desc: t('full_text_retrieval_desc')
  }
])

const open = ref(false)

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
  summary_switch: 0,
  rrf_weight: {}
})

const handleChangeModel = () => {
  onSave()
}

const vectorModelList = ref([])
const onVectorModelLoaded = (list) => {
  vectorModelList.value = list
  if (!formState.model_config_id || !formState.use_model) {
    if (vectorModelList.value.length > 0) {
      let modelConfig = vectorModelList.value[0]
      if (modelConfig) {
        let model = modelConfig.children[0]
        formState.use_model = model.name
        formState.model_config_id = model.model_config_id
      }
    }
  }
}

const handleSelectRetrievalMode = (val) => {
  formState.search_type = val
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
  newState.summary_switch = +newState.summary_switch

  if (newState.rerank_use_model === '') {
    // 这里是因为服务端可能会返回个空字符串，我这里改成undefined才有placeholder
    newState.rerank_use_model = void 0
  }
  newState.rrf_weight = newState.rrf_weight ? JSON.parse(newState.rrf_weight) : newState.rrf_weight
  Object.keys(newState).forEach((key) => {
    formState[key] = newState[key]
  })
}

// 防抖函数实现
const debounce = (func, wait) => {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
};

const handleWeightChange = () => {
  onSave();
}

const debouncedHandleWeightChange = debounce(handleWeightChange, 500); // 500ms 延迟

const onSave = () => {
  let params = { ...formState }

  if (!params.prompt && formState.prompt_type == '1') {
    return message.error(t('enter_custom_prompt'))
  }
  params.rrf_weight = JSON.stringify(formState.rrf_weight)

  if (formState.search_type == 3 || formState.search_type == 4) {
    // 当选择"知识图谱检索"和"全文检索"类型时，不显示"相似度阈值"设置项。
    delete params.similarity
  }

  if (formState.prompt_type == '0') {
    delete params.prompt
  }

  params.size = params.top_k
  saveLibrarySearch(params).then((res) => {
    message.success(t('save_success'))
  })
}

const modalOpen = ref(false)
const promptVal = ref('')
const handlePromptTypeChange = () => {
  if (formState.prompt_type == '1') {
    promptVal.value = formState.prompt
    modalOpen.value = true
  } else {
    onSave()
  }
}

const handleOk = () => {
  if (!promptVal.value) {
    return message.error(t('enter_prompt'))
  } else {
    formState.prompt = promptVal.value
    modalOpen.value = false
    onSave()
  }
}

const handleClose = () => {
  formState.prompt_type = '0'
  modalOpen.value = false
}

const afterOpenChange = (bool) => {
  // console.log('open', bool);
}

const showDrawer = async () => {
  let res = await getLibrarySearchFn()
  updateLibraryInfo(res.data)
  open.value = true
}

onMounted(() => {})
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
    margin-left: 2px;
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
        margin-left: 2px;
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