<template>
  <div class="recall-settings-box">
    <div class="form-heading">
      <img src="@/assets/img/library/detail/recall-settings.svg" alt="" />
      <span>{{ t('title_retrieval_config') }}</span>
    </div>

    <div class="form-item question-item">
      <div class="form-item-label">{{ t('label_test_question') }}</div>
      <div class="question-entry">
        <a-textarea
          v-model:value="formState.question"
          class="question-textarea"
          :placeholder="t('ph_test_question')"
          :auto-size="{ minRows: 2, maxRows: 8 }"
        />
        <div class="question-actions">
          <a-button type="primary" :loading="loading" @click="handleRecallTest">
            {{ t('btn_test') }}
            <span v-if="!loading" class="send-icon">
              <img src="@/assets/img/library/detail/recall-test-arrow.svg" alt="" />
            </span>
          </a-button>
        </div>
      </div>
    </div>

    <div class="form-item retrieval-item">
      <div class="form-item-label">{{ t('label_retrieval_mode') }}</div>
      <div class="retrieval-mode-items">
        <div
          v-for="item in retrievalModeList"
          :key="item.value"
          class="retrieval-mode-item"
          :class="{ active: formState.search_type == item.value }"
          role="radio"
          :aria-checked="formState.search_type == item.value"
          tabindex="0"
          @click="handleSelectRetrievalMode(item.value)"
          @keydown.enter.prevent="handleSelectRetrievalMode(item.value)"
          @keydown.space.prevent="handleSelectRetrievalMode(item.value)"
        >
          <svg-icon :name="item.iconName" class="mode-icon" />
          <span class="mode-title">{{ item.title }}</span>
          <SvgTextTag
            v-if="item.isRecommendation"
            class="recommendation-tag"
            :text="tCommon('recommendation')"
            :width="36"
            :height="21"
            :border-radius="6"
            background-color="#fb363f"
            text-color="#fff"
            :font-size="12"
          />
          <span class="selection-dot">
            <CheckOutlined v-if="formState.search_type == item.value" />
          </span>
        </div>
      </div>
      <div class="retrieval-mode-desc">
        {{ selectedRetrievalMode?.desc }}
      </div>
    </div>

    <div v-if="formState.search_type == 1" class="form-item weight-item">
      <div class="weight-slider-box" :style="weightCssVars">
        <div class="form-label-block">
          <div class="label-title">
            {{ tWeight('label_weight') }}
            <a-tooltip>
              <template #title>
                {{ tWeight('tooltip_weight') }}
              </template>
              <QuestionCircleOutlined />
            </a-tooltip>
          </div>
          <div class="item-list-box">
            <div class="list-item vector">
              <span class="dot"></span>
              <span class="text">{{ tWeight('label_vector') }}：{{ weightFormatter(formState.rrf_weight.vector) }}</span>
            </div>
            <div class="list-item fulltext">
              <span class="dot"></span>
              <span class="text">{{ tWeight('label_fulltext') }}：{{ weightFormatter(formState.rrf_weight.search) }}</span>
            </div>
          </div>
        </div>
        <a-slider v-model:value="weightValue" :tip-formatter="weightFormatter" @change="handleWeightChange" />
      </div>
      <div class="weight-hints">
        <span>{{ t('hint_semantic') }}</span>
        <span>{{ t('hint_keyword') }}</span>
      </div>
    </div>

    <div class="range-fields">
      <div class="form-item">
        <div class="form-item-label">
          <span>{{ t('label_top_k') }}</span>
          <a-tooltip>
            <template #title>{{ t('tooltip_top_k') }}</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="number-box">
          <a-slider v-model:value="formState.size" :min="1" :max="500" />
          <a-input-number v-model:value="formState.size" :min="1" :max="500" />
        </div>
      </div>
      <div v-if="formState.search_type <= 2" class="form-item">
        <div class="form-item-label">
          <span>{{ t('label_similarity_threshold') }}</span>
          <a-tooltip>
            <template #title>{{ t('tooltip_similarity_threshold') }}</template>
            <QuestionCircleOutlined class="question-icon" />
          </a-tooltip>
        </div>
        <div class="number-box">
          <a-slider v-model:value="formState.similarity" :min="0" :max="1" :step="0.01" />
          <a-input-number v-model:value="formState.similarity" :min="0" :max="1" :step="0.01" />
        </div>
      </div>
    </div>

    <div v-if="formState.search_type == 1 || formState.search_type == 3" class="form-item fulltext-item">
      <div class="form-item-label">{{ t('label_full_text_search_mode') }}</div>
      <a-radio-group v-model:value="formState.library_search_type">
        <a-radio value="fullTextSearch">{{ t('full_text_search') }}</a-radio>
        <a-radio value="keywordSearch">{{ t('keyword_match') }}</a-radio>
      </a-radio-group>
    </div>

    <div class="form-item rerank-item">
      <div class="form-item-label">{{ t('label_rerank_model') }}</div>
      <div class="rerank-controls">
        <a-switch
          v-model:checked="formState.rerank_status"
          :checkedValue="1"
          :unCheckedValue="0"
        />
        <div class="rerank-model">
          <ModelSelect
            v-model:modeName="formState.rerank_use_model"
            v-model:modeId="formState.rerank_model_config_id"
            modelType="RERANK"
            :placeholder="t('ph_select_rerank_model')"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import { useStorage } from '@/hooks/web/useStorage'
import { QuestionCircleOutlined, CheckOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { libraryRecallTest, getDefaultRrfWeight } from '@/api/library'
import SvgTextTag from '@/components/icons/SvgTextTag.vue'

const { t } = useI18n('views.library.library-details.components.recall-testing-form')
const { t: tCommon } = useI18n('common')
// 复用公共权重选择组件的文案命名空间，保证中英文案一致
const { t: tWeight } = useI18n('components.weight-select.index')

const route = useRoute()
const loading = ref(false)
const { getStorage, setStorage } = useStorage('localStorage')
const storageKey = `libraryRecallTest:${route.query.id ?? 'default'}`
const emit = defineEmits(['save', 'load', 'error'])

const retrievalModeList = computed(() => [
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
    iconName: 'search-check-icon',
    title: t('mode_fulltext_title'),
    value: 3,
    desc: t('mode_fulltext_desc')
  }
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
    search: 100,
    graph: 0
  },
  library_search_type: 'fullTextSearch'
})

const selectedRetrievalMode = computed(() =>
  retrievalModeList.value.find((mode) => mode.value === formState.search_type)
)

const persistedFields = [
  'question',
  'search_type',
  'size',
  'similarity',
  'rrf_weight',
  'library_search_type',
  'rerank_status',
  'rerank_use_model',
  'rerank_model_config_id'
]
let hasResolvedRrfWeight = false

// 产品要求：召回测试不展示图谱权重，graph 固定为 0，全文权重 = 100 - 向量权重
const normalizeRrfWeight = (weight) => {
  const vector = Number(weight?.vector) || 0
  return { vector, search: 100 - vector, graph: 0 }
}

// 本页面内联两段权重滑块（向量/全文），不复用公共 WeightSelect，避免改动影响其他页面
const weightValue = ref(0)

const handleWeightChange = () => {
  formState.rrf_weight = normalizeRrfWeight({ vector: weightValue.value })
}

watch(
  () => formState.rrf_weight,
  (val) => {
    weightValue.value = Number(val?.vector) || 0
  },
  { immediate: true, deep: true }
)

const weightCssVars = computed(() => ({
  '--background-liner': `linear-gradient(to right, #2475fc ${weightValue.value}%, #03B615 ${100 - weightValue.value}%)`
}))

function weightFormatter(value) {
  if (value <= 0) {
    return 0
  }
  if (value >= 100) {
    return 1
  }
  return (value / 100).toFixed(2)
}

const restoreFormState = () => {
  let storedState
  try {
    storedState = getStorage(storageKey)
  } catch {
    return false
  }
  if (!storedState || typeof storedState !== 'object' || Array.isArray(storedState)) {
    return false
  }

  const hasStoredRrfWeight = Boolean(
    storedState.rrf_weight &&
    typeof storedState.rrf_weight === 'object' &&
    ['vector', 'search', 'graph'].every((key) => typeof storedState.rrf_weight[key] === 'number')
  )

  persistedFields.forEach((field) => {
    const canRestore = field !== 'rrf_weight' || hasStoredRrfWeight
    if (canRestore && Object.prototype.hasOwnProperty.call(storedState, field)) {
      formState[field] = storedState[field]
    }
  })

  if (hasStoredRrfWeight) {
    formState.rrf_weight = normalizeRrfWeight(formState.rrf_weight)
  }

  // 旧缓存可能包含已移除的知识图谱模式，回退到当前可见的默认模式。
  if (!retrievalModeList.value.some((mode) => mode.value === formState.search_type)) {
    formState.search_type = 1
  }

  hasResolvedRrfWeight = hasStoredRrfWeight
  return hasStoredRrfWeight
}

const saveFormState = () => {
  const storedState = {}
  persistedFields.forEach((field) => {
    if (field !== 'rrf_weight' || hasResolvedRrfWeight) {
      storedState[field] = formState[field]
    }
  })
  try {
    setStorage(storageKey, storedState)
  } catch {
    // 浏览器禁用本地存储时，保留当前表单供本次测试使用。
  }
}

watch(
  () => formState.rrf_weight,
  () => {
    hasResolvedRrfWeight = true
  },
  { deep: true }
)
watch(formState, saveFormState, { deep: true })

const handleSelectRetrievalMode = (value) => {
  formState.search_type = value
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
  const params = {
    id: formState.id,
    question: formState.question,
    size: formState.size,
    similarity: formState.similarity,
    search_type: formState.search_type,
    rrf_weight: JSON.stringify(formState.rrf_weight),
    library_search_type: formState.library_search_type
  }
  if (formState.rerank_status == 1) {
    params.rerank_model_config_id = formState.rerank_model_config_id
    params.rerank_use_model = formState.rerank_use_model
  }
  loading.value = true
  emit('load')
  libraryRecallTest(params)
    .then((res) => {
      emit('save', res.data)
    })
    .catch(() => {
      emit('error')
    })
    .finally(() => {
      loading.value = false
    })
}

defineExpose({
  retry: handleRecallTest
})

onMounted(() => {
  const hasStoredRrfWeight = restoreFormState()
  if (!hasStoredRrfWeight) {
    getDefaultRrfWeight().then((res) => {
      if (hasResolvedRrfWeight) return
      hasResolvedRrfWeight = true
      formState.rrf_weight = normalizeRrfWeight(res.data || {
        vector: 0,
        search: 0,
        graph: 0
      })
    })
  }
})
</script>

<style lang="less" scoped>
.recall-settings-box {
  height: 100%;
  overflow-y: auto;
  padding: 0 24px 32px;
  color: #262626;
  scrollbar-width: thin;
  scrollbar-color: #c5cedb transparent;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    border-radius: 6px;
    background: #c5cedb;
  }

  &::-webkit-scrollbar-thumb:hover {
    background: #9eacc0;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  .form-heading {
    display: flex;
    align-items: center;
    gap: 6px;
    height: 24px;
    margin-bottom: 16px;
    font-size: 16px;
    font-weight: 600;

    img {
      width: 16px;
      height: 16px;
    }
  }

  .form-item-label {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-bottom: 8px;
    font-size: 14px;
    line-height: 22px;
  }

  .question-icon {
    color: #8c8c8c;
  }

  .question-item {
    margin-bottom: 40px;
  }

  .question-entry {
    min-height: 100px;
    padding: 10px 10px 8px;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    box-shadow: 0 4px 8px rgb(0 0 0 / 8%);

    :deep(.question-textarea) {
      padding: 0;
      resize: none;
      border: 0;
      box-shadow: none;
      background: transparent;
      font-size: 14px;
    }

    .question-actions {
      display: flex;
      justify-content: flex-end;
      margin-top: 2px;
    }

    :deep(.ant-btn) {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 4px;
      width: 68px;
      height: 28px;
      padding: 0 10px;
      border: 0;
      border-radius: 6px;
      background: #3157e2;
      box-shadow: none;
      color: #fff;
      font-size: 14px;
      font-weight: 500;
      line-height: 22px;
    }

    .send-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      flex: none;
      width: 16px;
      height: 16px;
      overflow: hidden;
    }

    .send-icon img {
      width: 9.33px;
      height: 11.33px;
      transform: rotate(-90deg) scaleX(-1);
    }
  }

  .retrieval-item {
    margin-bottom: 24px;
  }

  .retrieval-mode-items {
    display: flex;
    gap: 8px;
  }

  .retrieval-mode-item {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    min-width: 0;
    height: 38px;
    padding: 0 12px;
    border: 1px solid #d9d9d9;
    border-radius: 8px;
    color: #595959;
    cursor: pointer;

    &.active {
      border-color: #3157e2;
      background: #f5f9ff;
      color: #3157e2;
      font-weight: 600;
    }

    .mode-icon {
      flex: none;
      font-size: 16px;
    }

    .mode-title {
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      font-size: 14px;
    }

    .recommendation-tag {
      flex: none;
    }

    .selection-dot {
      display: flex;
      align-items: center;
      justify-content: center;
      flex: none;
      width: 16px;
      height: 16px;
      margin-left: auto;
      border: 1px solid #d9d9d9;
      border-radius: 50%;
      color: #fff;
      font-size: 10px;
    }

    &.active .selection-dot {
      border-color: #3157e2;
      background: #3157e2;
    }
  }

  .retrieval-mode-desc {
    margin-top: 10px;
    color: #8c8c8c;
    font-size: 14px;
    line-height: 22px;
  }

  .weight-item {
    margin-bottom: 24px;

    .weight-slider-box {
      .form-label-block {
        display: flex;
        align-items: center;
        justify-content: space-between;
        min-height: 24px;

        .label-title {
          display: flex;
          align-items: center;
          gap: 2px;
          margin-right: 8px;
        }

        .item-list-box {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 13px;

          .list-item {
            display: flex;
            align-items: center;
            gap: 4px;
            padding: 3px 6px;
            border-radius: 5px;
            background: #f3f7ff;
            font-size: 12px;
            line-height: 18px;

            & + .list-item {
              background: #f0f9f3;
            }

            .dot {
              width: 6px;
              height: 6px;
              border-radius: 50%;
            }

            &.vector {
              color: #2475fc;

              .dot {
                background: #2475fc;
              }
            }

            &.fulltext {
              color: #03b615;

              .dot {
                background: #03b615;
              }
            }
          }
        }
      }

      :deep(.ant-slider) {
        margin: 8px 0 4px;

        .ant-slider-rail {
          background: var(--background-liner);
        }

        .ant-slider-track {
          background: #2475fc;
        }
      }
    }
  }

  .weight-hints {
    display: flex;
    justify-content: space-between;
    color: #8c8c8c;
    font-size: 12px;
    line-height: 20px;
  }

  .range-fields {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 24px;
    margin-bottom: 24px;
  }

  .number-box {
    display: flex;
    align-items: center;
    gap: 24px;

    :deep(.ant-slider) {
      flex: 1;
      min-width: 0;
      margin: 0;
    }

    :deep(.ant-input-number) {
      flex: none;
      width: 80px;
    }
  }

  .fulltext-item {
    margin-bottom: 24px;
  }

  .rerank-controls {
    display: flex;
    align-items: center;
    gap: 16px;

    .rerank-model {
      width: 220px;
      max-width: calc(100% - 60px);
    }
  }
}

@media (max-width: 1100px) {
  .recall-settings-box {
    height: auto;
    overflow: visible;
  }
}
</style>
