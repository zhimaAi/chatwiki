<template>
  <div class="search-box-warpper">
    <div class="search-box-content">
      <div class="search-set-box">
        <SearchDrawer />
      </div>
      <div class="search-content" v-if="!isSearch">
        <div class="search-label">知识搜索</div>
        <div class="search-tip">快速搜索所选择的知识库中内容</div>
        <div class="search-input-box">
          <!-- <a-input-search
            v-model:value="message"
            placeholder="输入关键词搜索"
            style="width: 800px"
            size="large"
            @search="onSearch"
          /> -->

          <a-input class="search-input" @pressEnter="onSearch" v-model:value="message" allow-clear placeholder="输入关键词搜索">
            <template #suffix>
              <div class="search-icon-box" @click="onSearch">
                <svg-icon class="search-icon" name="search-right"></svg-icon>
              </div>
            </template>
          </a-input>
        </div>
      </div>
      <div class="search-content search-main" v-else>
        <div class="search-input-box">
          <a-input class="search-input" @pressEnter="onSearch" v-model:value="message" allow-clear placeholder="输入关键词搜索">
            <template #suffix>
              <div class="search-icon-box" @click="onSearch">
                <svg-icon class="search-icon" name="search-right"></svg-icon>
              </div>
            </template>
          </a-input>
        </div>
        <list-empty style="margin-top: 100px;" size="250" v-if="!renderedMarkdown && messageObj.finish && !libraryRecall.length">
          <div>
            <p style="color: #262626; font-size: 16px; font-weight: 600; line-height: 24px;">暂无任何内容</p>
          </div>
        </list-empty>
        <div class="search-content-box" v-else>
          <div class="search-tip-box">
            <svg-icon class="nav-icon" name="search-tip"></svg-icon>
            <span class="tip" v-if="!messageObj.finish">智能回答生成中...</span>
            <span class="tip complete" v-else>智能回答生成完毕</span>
          </div>

          <div class="container" v-if="!messageObj.content && !messageObj.finish && !isError">
            <div class="rect-container rect1">
              <div class="rect"></div>
            </div>
            <div class="rect-container rect2">
              <div class="rect"></div>
            </div>
            <div class="rect-container rect3">
              <div class="rect"></div>
            </div>
          </div>

          <!-- 智能回答 -->
          <div class="intelligent-answer" v-else>
            <div v-html="renderedMarkdown" v-if="renderedMarkdown"></div>
            <div v-else>暂无任何内容</div>
          </div>

          <!-- 相关分段 (3) -->
          <div class="section-box">
            <div class="tips">
              <div class="tips-text">以上内容由大模型生成</div>
              <div class="tips-line"></div>
            </div>
            <div v-if="libraryRecall.length" class="section-label">相关分段 <div>({{ libraryRecall.length }})</div>
            </div>
            <div class="section-item" v-for="item in libraryRecall" :key="item.id">
              <div class="section-item-nav">
                <div class="section-item-nav-left">
                  <svg-icon class="section-item-icon" name="document"></svg-icon>
                  <div class="section-item-label" @click="goToFile(item)">
                    {{ item.file_name }}<span v-if="!item.file_name">{{ item.library_name }}-精选</span>
                  </div>
                </div>
                <div class="section-item-nav-right">相似度：{{ formatNumber(item.similarity) }}</div>
              </div>
              <a-tooltip overlayClassName="search-content-tip" placement="top">
                <template #title>{{ item.content }}</template>
                <div class="section-item-content" v-html="highlightKeywords(item.content, message)"></div>
              </a-tooltip>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getModelConfigOption } from '@/api/model/index'
import { ref, watch, reactive, onMounted } from 'vue'
import SearchDrawer from './search-drawer.vue'
import { showErrorMsg } from '@/utils/index'
import { libraryRecallTest } from '@/api/library'
import ListEmpty from './list-empty.vue'
import MarkdownIt from 'markdown-it';
import { duplicateRemoval, removeRepeat } from '@/utils/index'

const emit = defineEmits(['search', 'defaultParmas'])

const props = defineProps({
  librarySearchData: {
    type: Object,
    default: null
  },
  messageObj: {
    type: Object,
    default: null
  },
  library_ids: {
    type: Array,
    default: () => []
  },
  isError: {
    type: Boolean,
    default: false
  }
})

const md = new MarkdownIt({
  html: true,        // 启用 HTML 标签
  linkify: true,     // 自动转换 URL 为链接
  typographer: true // 启用一些排版替换
});
const renderedMarkdown = ref('');
const message = ref('')
const isSearch = ref(false)
const loading = ref(false)
const libraryRecall = ref({})

watch(() => props.messageObj.content, () => {
  renderedMarkdown.value = md.render(props.messageObj.content);
})


watch(() => props.messageObj.error, (val) => {
  if (val) {
    showErrorMsg(val)
  }
})

const onSearch = async () => {

  if (!props.library_ids.length) {
    return showErrorMsg('请在左侧的知识库选择项内至少选择一个知识库')
  }

  if (!message.value) {
    return showErrorMsg('请输入消息内容')
  }

  emit('search', message.value)

  isSearch.value = true
}

function formatNumber(numStr) {
  return numStr.match(/^-?\d+(?:\.\d{0,2})?/)[0];
}

// 高亮关键词函数
const highlightKeywords = (content, keyword) => {
  if (!content || !keyword.trim()) return content
  
  // 转义正则特殊字符
  const escapedKeyword = keyword.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const regex = new RegExp(`(${escapedKeyword})`, 'gi')
  
  return content.replace(regex, '<span class="highlight">$1</span>')
}

// 新窗口跳转到文档
const goToFile = (item) => {
  if(!item.file_name){
    window.open(`#/library/details/categary-manage?id=${item.library_id}`)
    return
  }
  window.open(`#/library/preview?id=${item.file_id}`)
}

const librarySearchDataRef = ref(props.librarySearchData)

const formState = reactive({
  rerank_status: '',
  rerank_use_model: '',
  prompt: '',
  model_config_id: '',
  use_model: '',
  temperature: 0.5,
  max_token: 2000,
  similarity: 0.4,
  search_type: 1,
  summary_switch: 0,
})

const defaultParmas = ref({})

const handleRecallTest = (searchObj) => {
  librarySearchDataRef.value = searchObj
  let parmas = {
    model_config_id: searchObj.model_config_id || formState.model_config_id,
    use_model: searchObj.use_model || formState.use_model,
    temperature: searchObj.temperature || formState.temperature,
    max_token: searchObj.max_token || formState.max_token,
    size: 200,  // 固定值
    similarity: searchObj.similarity || formState.similarity,
    search_type: searchObj.search_type || formState.search_type,
    summary_switch: searchObj.summary_switch || formState.summary_switch,
    question: message.value,
    id: props.library_ids.join(','),
    recall_type: 1,
    rrf_weight: searchObj.rrf_weight,
  }
  if (searchObj.rerank_status) {
    parmas.rerank_status = searchObj.rerank_status
  }
  if (searchObj.rerank_use_model) {
    parmas.rerank_use_model = searchObj.rerank_use_model
  }
  if (searchObj.prompt) {
    parmas.prompt = searchObj.prompt
  }
  if (searchObj.rerank_status == 1) {
    parmas.rerank_model_config_id = searchObj.rerank_model_config_id
  }
  loading.value = true

  defaultParmas.value = Object.assign(parmas)
  emit('defaultParmas', defaultParmas.value)
  libraryRecallTest(parmas)
    .then((res) => {
      libraryRecall.value = res.data
    })
    .catch(() => {
    })
    .finally(() => {
      loading.value = false
    })
}


onMounted(async () => {


})

defineExpose({
  handleRecallTest
})
</script>

<style lang="less" scoped>
.search-box-warpper {
  width: 100%;
  height: 100vh;
  overflow-y: auto;

  .search-box-content {
    position: relative;

    .search-set-box {
      position: absolute;
      right: 16px;
      top: 16px;
    }

    .search-content {
      padding-top: 108px;
      display: flex;
      flex-direction: column;
      align-items: center;

      .search-label {
        color: #262626;
        font-size: 36px;
        font-style: normal;
        font-weight: 600;
        line-height: 44px;
        margin-bottom: 16px;
      }

      .search-tip {
        color: #595959;
        text-align: center;
        font-size: 16px;
        font-style: normal;
        font-weight: 400;
        line-height: 24px;
        margin-bottom: 32px;
      }

      .search-input-box {
        width: 800px;

        .search-input {
          height: 58px;
          border-radius: 12px;

          .search-icon-box {
            display: inline-flex;
            padding: 4px;
            align-items: center;
            gap: 10px;
            border-radius: 6px;

            &:hover {
              background: #E4E6EB;
              cursor: pointer;
            }

            .search-icon {
              font-size: 24px;
            }
          }
        }

        :deep(.ant-input-clear-icon) {
          display: inline-flex;
          padding: 4px;
          align-items: center;
          gap: 10px;
          border-radius: 6px;

          &:hover {
            background: var(--07, #E4E6EB);
          }
        }

        :deep(.anticon-close-circle) {
          font-size: 20px;
        }
      }

      .search-tip-box {
        margin-top: 16px;
        margin-bottom: 12px;
        width: 800px;
        display: flex;
        align-items: center;
        gap: 4px;

        .search-tip {
          font-size: 16px;
        }

        .tip {
          color: #595959;
          font-size: 14px;
          font-style: normal;
          font-weight: 400;
          line-height: 22px;
        }

        .complete {
          color: #6524fc;
        }
      }
    }

    .search-main {
      padding-top: 16px;
    }
  }
}

.container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rect-container {
  overflow: hidden;
  height: 20px;

  &.rect1 {
    width: 393px;
  }

  &.rect2,
  &.rect3 {
    width: 800px;
  }
}

.intelligent-answer {
  width: 800px;
  color: #262626;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 30px;
}

.rect {
  height: 100%;
  width: 100%;
  border-radius: 24px;
  background: linear-gradient(94deg, #00BFFF 2.9%, #C1C4FF 63.43%, #FFF 98.28%);
  transform-origin: left;
  animation: slide 1s infinite;

  .rect2 & {
    animation-delay: 0.1s;
  }

  .rect3 & {
    animation-delay: 0.2s;
  }
}

@keyframes slide {
  from {
    transform: scaleX(0);
  }

  to {
    transform: scaleX(1);
  }
}

.section-box {
  width: 800px;

  .tips {
    margin-top: 14px;
    line-height: 20px;
    display: flex;
    align-items: center;
    color: #8c8c8c;
    font-size: 12px;
    font-style: normal;
    font-weight: 400;

    .tips-text {
      padding-right: 5px;
    }

    .tips-line {
      flex: 1;
      height: 1px;
      border-radius: 1px;
      background: #D9D9D9;
    }
  }

  .section-label {
    display: flex;
    align-items: center;
    align-self: stretch;
    color: #000000;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    padding: 24px 0;
  }

  .section-item {
    cursor: pointer;
    margin-bottom: 24px;
    height: auto;
    align-self: stretch;
    border-radius: 6px;
    border: 1px solid var(--07, #F0F0F0);
    background: #FFF;
    padding: 16px;

    &:hover {
      background: var(--09, #F2F4F7);
    }

    .section-item-nav {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 8px;

      .section-item-nav-left {
        display: flex;
        align-items: center;
        gap: 3px;
      }

      .section-item-icon {
        font-size: 20px;
      }

      .section-item-label {
        color: #262626;
        font-size: 16px;
        font-style: normal;
        font-weight: 600;
        line-height: 24px;
      }

      .section-item-nav-right {
        color: #8c8c8c;
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
      }
    }

    .section-item-content {
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 3;
      overflow: hidden;
      color: #595959;
      text-overflow: ellipsis;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
    }
  }
}

:deep(.highlight) {
  color: red !important;
  padding: 0 2px;
}
</style>

<style lang="less">
.search-content-tip {
  max-width: 800px !important;

  .ant-tooltip-inner {
    max-width: 800px !important;
    width: max-content !important;
    max-height: 300px !important;
    min-height: 30px !important;
    overflow: hidden !important;
    overflow-y: auto !important;

    /* 滚动条样式 */
    &::-webkit-scrollbar {
      width: 4px; /*  设置纵轴（y轴）轴滚动条 */
      height: 4px; /*  设置横轴（x轴）轴滚动条 */
    }
    /* 滚动条滑块（里面小方块） */
    &::-webkit-scrollbar-thumb {
      border-radius: 0px;
      background: transparent;
    }
    /* 滚动条轨道 */
    &::-webkit-scrollbar-track {
      border-radius: 0;
      background: transparent;
    }

    /* hover时显色 */
    &:hover::-webkit-scrollbar-thumb {
      background: rgba(0, 0, 0, 0.2);
    }
    &:hover::-webkit-scrollbar-track {
      background: rgba(0, 0, 0, 0.1);
    }
  }
}
</style>