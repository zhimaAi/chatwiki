<template>
  <div class="open-search-page">
    <div class="search-page-header-wrapper">
      <div class="search-page-header">
        <div class="search-box" id="search-box">
          <input
            class="search-input"
            ref="searchInputRef"
            name="search"
            autocomplete="off"
            placeholder="请输入内容搜索"
            type="text"
            v-model="state.keyword"
            @keydown.enter="handleSearch"
          />
          <span class="search-clear" @click="handleClear" v-if="state.keyword"></span>
          <span class="search-btn" @click="handleSearch">
            <span class="search-icon"></span>
          </span>
        </div>

        <div class="action-box">
          <span class="action-btn">
            <img class="action-btn-icon" src="@/assets/img/open_sidebar.svg" alt="" />
          </span>
        </div>
      </div>
    </div>

    <div class="search-page-body">
      <div
        class="search-result-wrapper"
        v-if="state.searchResult.length > 0 || state.aiResult.text"
      >
        <h3 class="search-result-label">为您找到以下结果:</h3>
        <div class="ai-result-box" v-if="state.showAiResult">
          <h4 class="ai-result-label">以下结果由AI生成</h4>

          <div class="ai-result-text" id="ai-result-text">
            <div v-html="state.aiResult.text"></div>
          </div>

          <div
            class="ai-result-doc-box"
            v-if="state.aiResult.documents.length > 0"
            :class="{ expanded: state.isAiResultDocsExpanded }"
          >
            <h4 class="ai-result-doc-label">
              <img
                class="arrow-icon"
                src="@/assets/img/arrow.svg"
                alt="展开收起相关文档"
                @click="toggleAiResultDocs()"
              />
              <span class="label-text">相关文档 ({{ state.aiResult.documents.length }})</span>
              <span class="line"></span>
            </h4>
            <!-- ai result -->
            <ul
              class="ai-result-doc-items"
              id="ai-result-docs"
              v-show="state.isAiResultDocsExpanded"
            >
              <li
                class="ai-result-doc-item-wrap"
                v-for="doc in state.aiResult.documents"
                :key="doc.doc_key"
              >
                <router-link class="ai-result-doc-item" target="_blank" :to="'/doc/' + doc.doc_key">
                  <div class="doc-title-box">
                    <img class="doc-icon" src="@/assets/img/doc_icon2.svg" alt="" />
                    <h3 class="doc-title">{{ doc.file_name }}</h3>
                  </div>
                  <!-- <p class="doc-content"></p> -->
                </router-link>
              </li>
            </ul>
          </div>
        </div>

        <!-- search result -->
        <ul class="search-result-items" id="search-result-list">
          <li class="search-result-item-wrap">
            <router-link
              class="search-result-item"
              target="_blank"
              :to="`/doc/${item.doc_key}`"
              v-for="item in state.searchResult"
              :key="item.file_id"
            >
              <div class="doc-title-box">
                <i class="doc-icon"></i>
                <h3 class="doc-title" v-html="item.titleHtml"></h3>
              </div>
              <p class="doc-content" v-html="item.contentHtml"></p>
            </router-link>
          </li>
        </ul>
      </div>

      <!-- end -->
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useOpenDocStore } from '@/stores/open-doc'
import { getAiMessage, getSearchResult } from '@/api/open-doc'

const router = useRouter()
const route = useRoute()
const openDocStore = useOpenDocStore()

const v = computed(() => {
  return route.query.v
})

const docId = computed(() => {
  return route.params.id
})

const searchInputRef = ref(null)

const state = reactive({
  keyword: '',
  showAiResult: false,
  aiResult: {
    text: '',
    documents: [],
  },
  searchResult: [],
  isAiResultDocsExpanded: true, // 控制相关文档的展开收起状态
})

const handleClear = () => {
  // 清空输入框
  state.keyword = ''

  // 清空ai结果
  // state.aiResult.text = ''
  // state.aiResult.documents = []
  // state.showAiResult = false
  // 清空搜索结果
  // state.searchResult = []
  // 清空query.v
  // 获得焦点
  searchInputRef.value.focus()
}

const handleSearch = () => {
  if (!state.keyword) {
    return
  }

  state.showAiResult = false
  state.aiResult.text = ''
  state.aiResult.documents = []
  state.searchResult = []

  let query = route.query || {}

  router.replace({
    query: {
      ...query,
      v: state.keyword,
    },
  })

  if (state.keyword == v.value) {
    getAiResult()
    getSearchResultData()
  }
}

const getData = async () => {
  let v = route.query.v
  await openDocStore.getSearch(docId.value, v)
}

let ai = null

const getAiResult = () => {
  if (ai) {
    ai.abort()
    ai = null
  }

  ai = getAiMessage({
    v: state.keyword,
    id: docId.value,
  })

  ai.onMessage = (res) => {
    if (res.event == 'sending') {
      if (!state.showAiResult) {
        state.showAiResult = true
      }
      state.aiResult.text += res.data
    } else if (res.event == 'quote_file') {
      state.aiResult.documents = JSON.parse(res.data)
    }
  }
}

const getSearchResultData = () => {
  getSearchResult({
    v: state.keyword,
    id: docId.value,
  }).then((res) => {
    let list = res.data || []

    list.forEach((item) => {
      if (state.keyword.length) {
        item.titleHtml = item.title.replace(
          new RegExp(state.keyword, 'gi'),
          (match) => `<span class="keyword-text">${match}</span>`,
        )
        item.contentHtml = item.content.replace(
          new RegExp(state.keyword, 'gi'),
          (match) => `<span class="keyword-text">${match}</span>`,
        )
      }
    })

    state.searchResult = list
  })
}

const toggleAiResultDocs = () => {
  // 切换相关文档的展开收起状态
  state.isAiResultDocsExpanded = !state.isAiResultDocsExpanded
}

const init = () => {
  state.keyword = v.value
  state.showAiResult = false
  state.aiResult.text = ''
  state.aiResult.documents = []
  state.searchResult = []
  state.isAiResultDocsExpanded = true // 重置展开状态

  if (v.value) {
    getAiResult()
    getSearchResultData()
  }
}

onMounted(async () => {
  await getData()

  init()
})

// 监听路由参数变化，当params.id更新时重新调用getData
watch(docId, () => {
  init()
})

watch(v, () => {
  init()
})
</script>

<style lang="less" scoped>
.open-search-page {
  padding: 32px 0 0 0;
}
.search-page-header {
  display: flex;
  padding: 8px 16px;
  background: #fff;
}

.search-page-header .action-box {
  display: flex;
  align-items: center;
  margin-left: 16px;
}
.search-page-header .action-box .action-btn {
  width: 40px;
  height: 40px;
  line-height: 40px;
  border-radius: 8px;
  text-align: center;
  transition: all 0.2s;
}
.search-page-header .action-box .action-btn:hover {
  cursor: pointer;
  background-color: #f0f2f5;
}

.search-box {
  position: relative;
  flex: 1;
  border-radius: 12px;
  overflow: hidden;
}
.search-box .search-input {
  width: 100%;
  height: 40px;
  padding: 0 80px 0 24px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 12px;
  color: #1a1a1a;
  border: 2px solid #262626;
  transition: border 0.2s;

  &::placeholder {
    font-weight: 400;
    color: #ccc;
  }
}
.search-box .search-clear {
  position: absolute;
  right: 48px;
  top: 8px;
  width: 24px;
  height: 24px;
  background: url('@/assets/img/search_clear.svg') no-repeat;
  cursor: pointer;
}
.search-box .search-btn {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 40px;
  width: 40px;
  transition: background 0.2s;
}
.search-box .search-icon {
  width: 24px;
  height: 24px;
  background: url('@/assets/img/search_black.svg') no-repeat;
}
.search-box .search-btn:hover {
  cursor: pointer;
}
.search-box .search-input {
  border: 2px solid #2475fc;
}
.search-box .search-input::placeholder {
  color: #ccc;
}
.search-box .search-btn {
  border-radius: 0 12px 12px 0;
  background: #2475fc;
}
.search-box .search-icon {
  background: url('@/assets/img/search_white.svg') no-repeat;
}
.search-box.has-value .search-clear {
  display: block;
}
.search-page-body {
  padding: 16px 16px;
}
.search-page-body .search-result-label {
  line-height: 24px;
  margin-bottom: 8px;
  font-size: 16px;
  font-weight: 400;
  color: rgb(140, 140, 140);
}
.ai-result-box {
  padding: 16px;
  margin-bottom: 32px;
  border-radius: 12px;
  border: 1px solid #2475fc;
  background: #f0f5ff;
}
.ai-result-box .ai-result-label {
  line-height: 24px;
  margin-bottom: 4px;
  font-size: 16px;
  font-weight: 600;
  color: #14161a;
}
.ai-result-box .ai-result-text {
  line-height: 22px;
  font-size: 14px;
  font-weight: 400;
  color: rgb(58, 69, 89);
}
.ai-result-doc-label {
  display: flex;
  align-items: center;
  margin-top: 16px;
  margin-bottom: 8px;
}
.ai-result-doc-label .arrow-icon {
  width: 16px;
  height: 16px;
  margin-right: 4px;
  transition: all 0.2s;
}
.ai-result-doc-label .label-text {
  line-height: 24px;
  font-size: 14px;
  font-weight: 400;
  color: rgb(122, 134, 153);
}
.ai-result-doc-label .line {
  margin-left: 12px;
  flex: 1;
  height: 1px;
  border-radius: 1px;
  background: rgba(216, 221, 230, 1);
}
.ai-result-doc-items .ai-result-doc-item-wrap {
  margin-bottom: 8px;
}
.ai-result-doc-items .ai-result-doc-item-wrap:last-child {
  margin-bottom: 0;
}
.ai-result-doc-items .ai-result-doc-item {
  display: block;
  line-height: 22px;
  padding: 16px;
  border-radius: 6px;
  background-color: #fff;
}
.ai-result-doc-items .doc-title-box {
  display: flex;
  align-items: center;
}
.ai-result-doc-items .doc-title-box .doc-icon {
  margin-right: 4px;
}
.ai-result-doc-items .doc-title-box .doc-title {
  line-height: 22px;
  margin-bottom: 0;
  font-size: 14px;
  font-weight: 400;
  color: #262626;
}
.ai-result-doc-items .doc-content {
  height: 44px;
  line-height: 22px;
  margin-top: 4px;
  font-size: 14px;
  font-weight: 400;
  color: #8c8c8c;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}
/* search  result */
.search-result-items .search-result-item-wrap {
  margin-bottom: 32px;
}
.search-result-items .search-result-item-wrap:last-child {
  margin-bottom: 0;
}
.search-result-items .search-result-item {
  display: block;
  background-color: #fff;
}
.search-result-items .doc-title-box {
  position: relative;
  display: flex;
  align-items: center;
  margin-bottom: 4px;
}
.search-result-items .doc-title-box::before {
  display: block;
  content: '';
  width: 16px;
  height: 16px;
  margin-right: 4px;
  background: url('../../../assets/img/doc_icon2.svg') 0 0 no-repeat;
}
.search-result-items .doc-title-box .doc-title {
  line-height: 24px;
  margin-bottom: 0;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}
.search-result-items .doc-content {
  line-height: 22px;
  padding-left: 20px;
  font-size: 14px;
  font-weight: 400;
  color: #595959;
  display: -webkit-box;
  -webkit-line-clamp: 6;
  line-clamp: 6;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}
.search-result-items .search-result-item :deep(.keyword-text) {
  color: #fb363f;
}
.ai-result-doc-box .ai-result-doc-items {
  max-height: 256px;
  overflow: hidden;
}
.ai-result-doc-box.expanded .ai-result-doc-items {
  max-height: 100% !important;
  overflow: hidden;
}
.ai-result-doc-box.expanded .ai-result-doc-label .arrow-icon {
  transform: rotate(-180deg);
}

@media (min-width: 992px) {
  .search-result-wrapper,
  .search-page-header {
    width: 100%;
    max-width: 1000px;
    margin: 0 auto;
  }
  .search-page-header-wrapper {
    margin: 0 24px;
  }
  .search-page-header {
    position: relative;
    padding: 0;
  }
  .search-page-header .action-box {
    display: none;
  }
  .search-page-header .search-box {
    width: 100%;
  }
  .search-page-header .search-box .search-input {
    height: 58px;
    padding: 0 106px 0 24px;
  }
  .search-page-header .search-box .search-btn {
    width: 58px;
    height: 58px;
  }
  .search-page-header .search-box .search-clear {
    top: 17px;
    right: 74px;
  }
  .search-page-body {
    flex: 1;
    padding: 0 24px;
    margin: 24px 0;
    overflow-y: auto;
  }
  .ai-result-box {
    padding: 24px 24px 16px 24px;
  }
  .ai-result-box .ai-result-label {
    margin-bottom: 8px;
  }
  .ai-result-doc-items {
    display: flex;
    flex-flow: row wrap;
    margin: 0 -8px;
    overflow: hidden;
  }

  .ai-result-doc-items .ai-result-doc-item-wrap {
    width: 50%;
    padding: 8px;
    margin-bottom: 0;
    box-sizing: border-box;
  }

  .search-result-items .search-result-item-wrap {
    margin-bottom: 24px;
  }
  .search-result-items .search-result-item {
    padding: 16px;
    border-radius: 6px;
    transition: all 0.2s;
  }
  .search-result-items .search-result-item:hover {
    background-color: #e4e6eb;
  }
  .search-result-items .doc-title-box {
    margin-bottom: 8px;
  }
  .ai-result-doc-box .ai-result-doc-items {
    max-height: 148px;
    overflow: hidden;
  }
  .ai-result-doc-box.expanded .ai-result-doc-items {
    max-height: 100% !important;
    overflow: hidden;
  }
}
</style>
