<template>
  <div class="recall-testing-box">
    <header class="page-heading">
      <div class="page-tip">
        <InfoCircleFilled />
        <span>{{ tForm('alert_tip') }}</span>
      </div>
    </header>

    <div class="recall-workspace">
      <div class="recall-from-box">
        <RrcallTestingForm ref="formRef" @load="handleLoading" @save="handleRecallTest" @error="handleRecallError" />
      </div>

      <div ref="contRef" class="content-list-box">
        <div class="results-heading">
          <img src="@/assets/img/library/detail/recall-results.svg" alt="" />
          <span>{{ t('title_recall_results') }}</span>
          <span class="results-count">{{ t('text_segment_count', { count: lists.length }) }}</span>
        </div>

        <div v-if="resultStatus === 'loading'" class="loading-results" role="status">
          <div class="loading-status">
            <img src="@/assets/img/library/detail/recall-loading.svg" alt="" />
            <span>{{ t('text_retrieving') }}</span>
          </div>
          <div class="skeleton-list" aria-hidden="true">
            <div v-for="index in 3" :key="index" class="skeleton-card">
              <span class="skeleton-line"></span>
              <span class="skeleton-line"></span>
              <span class="skeleton-line"></span>
              <span class="skeleton-line"></span>
            </div>
          </div>
        </div>
        <div v-else-if="resultStatus === 'empty' || resultStatus === 'error'" class="state-box" role="status">
          <div class="state-icon" :class="{ 'error-icon': resultStatus === 'error' }">
            <img
              :src="resultStatus === 'error' ? recallErrorIcon : recallEmptyIcon"
              alt=""
            />
          </div>
          <h3>{{ resultStatus === 'error' ? t('title_retrieval_error') : t('title_empty_segments') }}</h3>
          <p>
            {{ resultStatus === 'error' ? t('text_retrieval_error') : t('text_empty_segments') }}
          </p>
          <button v-if="resultStatus === 'error'" type="button" class="retry-button" @click="handleRetry">
            {{ t('btn_retry') }}
          </button>
        </div>
        <div v-else class="results-scroll">
          <div v-for="item in lists" :key="item.id" class="list-item">
            <div class="card-meta">
              <div class="meta-line">
                <span>{{ t('label_id') }}{{ item.id }}</span>
                <span class="meta-divider"></span>
                <span class="char-count">{{ t('text_total_chars', { count: item.word_total }) }}</span>
                <span class="similarity">{{ t('label_similarity') }}{{ item.similarity }}</span>
              </div>
              <button type="button" class="source-link" @click="handlePreview(item)">
                <img src="@/assets/img/library/detail/recall-source.svg" alt="" />
                <span>
                  {{ item.file_name || item.library_name }}
                  <template v-if="!item.file_name">{{ t('text_selected') }}</template>
                </span>
              </button>
            </div>

            <div class="card-content">
              <h3 v-if="item.title">{{ item.title }}</h3>
              <template v-if="libraryType == 2">
                <div v-if="item.question" class="content-box">{{ t('label_question') }}{{ item.question }}</div>
                <div v-if="item.similar_questions && item.similar_questions.length" class="content-box similar-questions-box">
                  <span>{{ t('label_similar_questions') }}</span>
                  <ul class="similar-questions-list">
                    <li v-for="(value, index) in item.similar_questions" :key="index">{{ value }}</li>
                  </ul>
                </div>
                <div v-if="item.answer" class="content-box">{{ t('label_answer') }}{{ item.answer }}</div>
              </template>
              <div v-else class="content-box" v-html="item.content"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useLibraryStore } from '@/stores/modules/library'
import { InfoCircleFilled } from '@ant-design/icons-vue'
import RrcallTestingForm from './components/recall-testing-form.vue'
import { useMathJax } from '@/composables/useMathJax.js'
import recallEmptyIcon from '@/assets/img/library/detail/recall-empty.svg'
import recallErrorIcon from '@/assets/img/library/detail/recall-error.svg'

const { t } = useI18n('views.library.library-details.recall-testing')
const { t: tForm } = useI18n('views.library.library-details.components.recall-testing-form')
const libraryStore = useLibraryStore()
const { renderMath } = useMathJax()

const contRef = ref(null)
const formRef = ref(null)
const resultStatus = ref('empty')
const lists = ref([])
const libraryType = computed(() => libraryStore.type)

const handleRecallTest = (data = []) => {
  data.forEach((item) => {
    if (item.similar_questions) {
      item.similar_questions = JSON.parse(item.similar_questions)
    }
  })
  lists.value = data || []
  resultStatus.value = lists.value.length ? 'success' : 'empty'
  if (resultStatus.value === 'success') {
    nextTick(() => renderMath(contRef.value))
  }
}

const handleLoading = () => {
  lists.value = []
  resultStatus.value = 'loading'
}

const handleRecallError = () => {
  lists.value = []
  resultStatus.value = 'error'
}

const handleRetry = () => {
  formRef.value?.retry()
}

const handlePreview = (record) => {
  if (!record.file_name) {
    window.open(`/#/library/details/categary-manage?id=${record.library_id}`)
    return
  }
  window.open(`/#/library/preview?id=${record.file_id}`)
}
</script>

<style lang="less" scoped>
.recall-testing-box {
  display: flex;
  flex-direction: column;
  height: calc(100% + 24px);
  margin: -24px -10px 0 -24px;
  overflow: hidden;
  background: #fff;
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

  .results-scroll,
  .loading-results,
  .card-content {
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
  }

  .page-heading {
    flex: none;
    padding: 24px 24px 0;
    margin-bottom: 14px;

    h1 {
      margin: 0 0 16px;
      color: #000;
      font-size: 20px;
      font-weight: 600;
      line-height: 28px;
    }
  }

  .page-tip {
    display: flex;
    align-items: center;
    gap: 8px;
    min-height: 42px;
    padding: 9px 16px;
    border: 1px solid #99bffd;
    border-radius: 6px;
    background: #e9f1fe;
    color: #3a4559;
    font-size: 14px;
    line-height: 22px;

    :deep(.anticon) {
      flex: none;
      color: #2475fc;
      font-size: 16px;
    }
  }

  .recall-workspace {
    display: grid;
    grid-template-columns: minmax(0, 59%) minmax(0, 41%);
    flex: 1;
    min-height: 0;
  }

  .recall-from-box {
    min-width: 0;
    overflow: hidden;
    border-right: 1px solid #d9d9d9;
  }

  .content-list-box {
    display: flex;
    flex-direction: column;
    min-width: 0;
    min-height: 0;
    padding: 0 24px 24px;
    position: relative;
  }

  .results-heading {
    display: flex;
    align-items: center;
    gap: 6px;
    flex: none;
    height: 24px;
    margin-bottom: 24px;
    color: #262626;
    font-size: 16px;
    font-weight: 600;

    img {
      width: 16px;
      height: 16px;
    }
  }

  .results-count {
    margin-left: auto;
    color: #8c8c8c;
    font-size: 14px;
    font-weight: 400;
  }

  .results-scroll {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
  }

  .list-item {
    display: flex;
    flex-direction: column;
    box-sizing: border-box;
    height: 462px;
    overflow: hidden;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    background: #fff;
    box-shadow: 0 4px 8px rgb(0 0 0 / 8%);

    & + .list-item {
      margin-top: 16px;
    }
  }

  .card-meta {
    flex: none;
    padding: 14px 14px 16px;
    background: #f8faff;
    color: #3a4559;
    font-size: 14px;
    line-height: 22px;
  }

  .meta-line {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 0;
    white-space: nowrap;
  }

  .meta-divider {
    width: 1px;
    height: 12px;
    background: #d8dde5;
  }

  .char-count {
    overflow: hidden;
    color: #7a8699;
    text-overflow: ellipsis;
  }

  .similarity {
    flex: none;
    margin-left: auto;
    padding: 4px 8px;
    border-radius: 5px;
    background: #eaf6ef;
    color: #2c8150;
    font-size: 12px;
    line-height: 20px;
  }

  .source-link {
    display: flex;
    align-items: center;
    gap: 4px;
    min-width: 0;
    width: 100%;
    margin-top: 12px;
    padding: 0;
    border: 0;
    background: transparent;
    color: #164799;
    cursor: pointer;
    font: inherit;
    text-align: left;

    img {
      flex: none;
      width: 16px;
      height: 16px;
    }

    span {
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }

  .card-content {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    padding: 16px 14px;
    color: #595959;
    font-size: 14px;
    line-height: 32px;
    overflow-wrap: break-word;

    .content-box {
      white-space: pre-wrap;
    }

    h3 {
      margin: 0 0 4px;
      color: #262626;
      font-size: 16px;
      font-weight: 600;
      line-height: 24px;
    }

    :deep(p),
    :deep(ol),
    :deep(ul) {
      margin-top: 0;
      margin-bottom: 0;
    }
  }

  .similar-questions-box {
    display: flex;
  }

  .similar-questions-list {
    flex: 1;
    margin: 0;
    padding-left: 0;
    list-style: none;
  }

  .loading-results {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
  }

  .loading-status {
    display: flex;
    align-items: center;
    gap: 11px;
    box-sizing: border-box;
    height: 42px;
    padding: 9px 14px;
    border-radius: 11px;
    background: #edf2fe;
    color: #1e45b8;
    font-size: 13.5px;
    line-height: 21.6px;

    img {
      flex: none;
      width: 24px;
      height: 24px;
      animation: recall-loading-rotate 1s linear infinite;
    }
  }

  .skeleton-list {
    display: flex;
    flex-direction: column;
    gap: 9px;
    margin-top: 11px;
  }

  .skeleton-card {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
    box-sizing: border-box;
    height: 109px;
    padding: 15px;
    border-radius: 11px;
    background: #fafbfc;
  }

  .skeleton-line {
    height: 11px;
    border-radius: 5px;
    background: #e3e7ec;
    opacity: 0.71;

    &:nth-child(1) {
      width: 85%;
      height: 14px;
    }

    &:nth-child(2) {
      width: 100%;
      height: 7px;
    }

    &:nth-child(3) {
      width: 86%;
    }

    &:nth-child(4) {
      width: 68%;
    }
  }

  .state-box {
    display: flex;
    flex-direction: column;
    align-items: center;
    flex: 1;
    padding: 52px 24px;
    text-align: center;

    h3 {
      margin: 0 0 8px;
      color: #262626;
      font-size: 16px;
      font-weight: 600;
      line-height: 24px;
    }

    p {
      max-width: 340px;
      margin: 0;
      color: #8c8c8c;
      font-size: 14px;
      line-height: 22px;
    }
  }

  .state-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    flex: none;
    width: 60px;
    height: 60px;
    margin-bottom: 18px;
    border-radius: 11px;
    background: #f4f6f8;

    &.error-icon {
      background: #fef6e6;
    }

    img {
      width: 26px;
      height: 26px;
    }
  }

  .retry-button {
    box-sizing: border-box;
    height: 28px;
    margin-top: 18px;
    padding: 0 10px;
    border: 1px solid #2475fc;
    border-radius: 6px;
    background: #fff;
    color: #2475fc;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    line-height: 22px;

    &:hover {
      background: #f5f9ff;
    }
  }
}

@keyframes recall-loading-rotate {
  to {
    transform: rotate(360deg);
  }
}

@media (prefers-reduced-motion: reduce) {
  .recall-testing-box .loading-status img {
    animation: none;
  }
}

@media (max-width: 1100px) {
  .recall-testing-box {
    overflow-y: auto;

    .recall-workspace {
      display: block;
    }

    .recall-from-box {
      border-right: 0;
      border-bottom: 1px solid #d9d9d9;
    }

    .content-list-box {
      min-height: 500px;
      padding-top: 24px;
    }
  }
}
</style>
