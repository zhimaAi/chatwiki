<template>
  <div class="open-doc-page">
    <div class="document-page-header">
      <div class="page-header-content">
        <h1 class="document-title">
          <span class="document-title-icon">{{ docData.doc_icon }}</span>
          <a class="document-title-link">{{ docData.title }}</a>
        </h1>
        <div class="document-last-update">最近更新时间：{{ docData.updatated_at }}</div>
      </div>
    </div>
    <div class="document-page-body">
      <div class="markdown-content">
        <div class="vditor-reset">
          <MdEditor ref="mdEditorRef" />
        </div>
        <div class="document-footer">
          <div class="page-turning" v-if="docData.prev_doc || docData.next_doc">
            <div>
              <a-tooltip v-if="docData.prev_doc">
                <template #title>
                  <span>{{ docData.prev_doc.title }}</span>
                </template>
                <router-link class="prev-page page-link"
                  :to="`/doc/${docData.prev_doc.doc_key}?${previewKey ? 'preview=' + previewKey : ''}${token ? '&token=' + token : ''}`">
                  <svg class="w-icon" aria-hidden="true">
                    <use xlink:href="#icon-jiantou_xiangzuo"></use>
                  </svg>
                  <span>上一篇：</span>
                  <span class="page-title">{{ docData.prev_doc.title }}</span>
                </router-link>
              </a-tooltip>
              <a class="prev-page" href="javascript:;" v-else>
                <svg class="w-icon" aria-hidden="true">
                  <use xlink:href="#icon-jiantou_xiangzuo"></use>
                </svg>
                <span>上一篇：</span>
                <span>无</span>
              </a>
            </div>

            <div>
              <a-tooltip v-if="docData.next_doc">
                <template #title>
                  <span>{{ docData.next_doc.title }}</span>
                </template>
                <router-link class="next-page page-link"
                  :to="`/doc/${docData.next_doc.doc_key}?${previewKey ? 'preview=' + previewKey : ''}${token ? '&token=' + token : ''}`">
                  <span>下一篇：</span>
                  <span class="page-title">{{ docData.next_doc.title }}</span>
                  <svg class="w-icon" aria-hidden="true">
                    <use xlink:href="#icon-jiantou_xiangyou"></use>
                  </svg>
                </router-link>
              </a-tooltip>

              <a class="next-page" href="javascript:;" v-else>
                <span>下一篇：</span>
                <span>无</span>
                <svg class="w-icon" aria-hidden="true">
                  <use xlink:href="#icon-jiantou_xiangyou"></use>
                </svg>
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import dayjs from 'dayjs'
import { ref, computed, onMounted, reactive, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useOpenDocStore } from '@/stores/open-doc'

import MdEditor from './components/md-editor.vue'

const route = useRoute()
const openDocStore = useOpenDocStore()
const mdEditorRef = ref(null)

const docData = reactive({
  title: '',
  body: '',
  prev_doc: null,
  next_doc: null,
  updatated_at: '',
  doc_icon:'',
})

const token = computed(() => {
  return openDocStore.token
})

const previewKey = computed(() => {
  return openDocStore.previewKey
})

const docId = computed(() => {
  return route.params.id
})

const getData = async () => {
  let data = await openDocStore.getDoc(docId.value)

  docData.title = data.title
  docData.body = data.body
  docData.prev_doc = data.prev_doc
  docData.next_doc = data.next_doc
  docData.doc_icon = data.doc_icon
  docData.updatated_at = dayjs(data.update_time*1000).format('YYYY-MM-DD HH:mm:ss')

  mdEditorRef.value.init('preview', data.body)
}

onMounted(() => {
  getData()
})

// 监听路由参数变化，当params.id更新时重新调用getData
watch(docId, () => {
  getData()
})
</script>

<style lang="less" scoped>
.open-doc-page {
  .document-page-header {
    padding: 16px;
    background: #F2F4F7;

    .page-header-content {
      width: 100%;
      max-width: 1000px;
      margin: 0 auto;

      .document-title {
        display: flex;
        align-items: center;

        .document-title-icon {
          width: 24px;
          height: 24px;
          margin-right: 8px;
          font-size: 20px;
        }

        .document-title-link {
          line-height: 32px;
          font-size: 24px;
          font-weight: 600;
          color: #262626;
        }
      }

      .document-last-update {
        line-height: 22px;
        padding-left: 32px;
        margin-top: 8px;
        font-size: 14px;
        font-weight: 400;
        color: #8c8c8c;
      }
    }
  }
  .document-page-body{
    padding: 16px;
  }

  .markdown-content {
    width: 100%;
    max-width: 1000px;
    margin: 0 auto;
  }

  .vditor-reset {
    overflow: hidden;
    line-height: 1.5;
    font-size: 16px;
    word-wrap: break-word;
    word-break: break-word;
  }
}

.document-footer {
  margin-top: 32px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;

  .page-turning {
    display: flex;
    justify-content: space-between;
  }

  .page-turning .prev-page,
  .page-turning .next-page {
    font-size: 14px;
    color: #181818;
  }

  .page-turning .prev-page:hover,
  .page-turning .next-page:hover {
    color: #1366ec !important;
  }

  .page-link {
    display: inline-flex;
    align-items: center;

    .page-title {
      max-width: 300px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

@media (max-width: 992px) {
  .open-doc-page {
    .document-page-body {
      padding: 60px 16px 16px 16px;
    }
  }

  .document-footer {
    .page-turning {
      flex-flow: column;
    }

    .page-turning .w-icon {
      display: none;
    }

    .page-turning .next-page {
      margin-top: 16px;
    }
  }
}
</style>
