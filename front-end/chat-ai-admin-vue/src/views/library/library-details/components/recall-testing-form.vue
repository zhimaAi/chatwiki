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
    .hover-label{
      padding: 0 8px;
      width: fit-content;
      border-radius: 6px;
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 4px;
      &:hover{
        background: var(--07, #E4E6EB);
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
      {{ t('alert_tip') }}
    </div>
    <div class="form-box">
      <div class="form-item">
        <div class="form-item-body">
          <a-textarea
            style="height: 76px"
            v-model:value="formState.question"
            :placeholder="t('ph_test_question')"
          />
        </div>
      </div>
      <div class="form-item">
        <a-button :loading="loading" @click="handleRecallTest" type="primary" block>{{ t('btn_test') }}</a-button>
      </div>
      <div class="form-item">
        <div class="form-item-label">
          <div class="hover-label" @click="isHide = !isHide">
            {{ t('label_retrieval_mode') }}
            <DownOutlined v-if="isHide" />
            <UpOutlined v-else />
          </div>
        </div>
        <div class="form-item-body">
          <div class="retrieval-mode-items">
            <div
              class="retrieval-mode-item"
              :class="{ active: formState.search_type == item.value }"
              v-for="item in showRetrievalModeList"
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
              <a-slider class="custom-slider" v-model:value="formState.size" :min="1" :max="500" />
            </div>
            <div class="number-input-box">
              <a-input-number v-model:value="formState.size" :min="1" :max="500" />
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
              <a-input-number v-model:value="formState.similarity" :min="0" :max="1" :step="0.01" />
            </div>
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
          />
        </div>
        <div class="form-item-body">
          <ModelSelect
            modelType="RERANK"
            v-model:modeName="formState.rerank_use_model"
            v-model:modeId="formState.rerank_model_config_id"
            style="width: 320px"
            :placeholder="t('ph_select_rerank_model')"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getModelConfigOption } from '@/api/model/index'
import { reactive, ref, toRaw, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import { QuestionCircleOutlined, DownOutlined, UpOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { libraryRecallTest, getDefaultRrfWeight } from '@/api/library'
import WeightSelect from '@/components/weight-select/index.vue'

const { t } = useI18n('views.library.library-details.components.recall-testing-form')

const route = useRoute()
const loading = ref(false)

const emit = defineEmits(['save', 'load'])

const isHide = ref(true)

const retrievalModeList = ref([
{
    iconName: 'mix-icon',
    title: t('mode_mix_title'),
    value: 1,
    isRecommendation: true,
    desc: t('mode_mix_desc')
  },
  {
    iconName: 'vector-icon',
    title: t('mode_vector_title'),
    value: 2,
    desc: t('mode_vector_desc')
  },
  {
    iconName: 'graph-icon',
    title: t('mode_graph_title'),
    value: 4,
    desc: t('mode_graph_desc')
  },
  {
    iconName: 'search-check-icon',
    title: t('mode_fulltext_title'),
    value: 3,
    desc: t('mode_fulltext_desc')
  },
])



const formState = reactive({
  rerank_status: 0,
  rerank_use_model: undefined,
  rerank_model_config_id: undefined,
  search_type: 1,
  question: '',
  similarity: 0.6,
  size: 5,
  id: route.query.id,
  rrf_weight: {
    vector: 0,
    search: 0,
    graph: 0,
  }
})

const showRetrievalModeList = computed(()=>{
  if(isHide.value){
    return retrievalModeList.value.filter(item => item.value == formState.search_type )
  }else{
    return retrievalModeList.value
  }
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
    return message.error(t('msg_select_rerank_model'))
  }

  show.value = false
  triggerChange()
}

const triggerChange = () => {
  emit('change', toRaw(formState))
}
const handleRecallTest = () => {
  if (!formState.similarity) {
    return message.error(t('msg_input_similarity'))
  }
  if (!formState.size) {
    return message.error(t('msg_input_size'))
  }
  if (!formState.question) {
    return message.error(t('msg_input_question'))
  }
  let parmas = {
    id: formState.id,
    question: formState.question,
    size: formState.size,
    similarity: formState.similarity,
    search_type: formState.search_type,
    rrf_weight: JSON.stringify(formState.rrf_weight),
  }
  if (formState.rerank_status == 1) {
    parmas.rerank_model_config_id = formState.rerank_model_config_id
    parmas.rerank_use_model = formState.rerank_use_model
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

onMounted(() => {
  getDefaultRrfWeight().then((res) => {
    formState.rrf_weight = res.data || {
      vector: 0,
      search: 0,
      graph: 0
    }
  })
})
defineExpose({
  open
})
</script>
