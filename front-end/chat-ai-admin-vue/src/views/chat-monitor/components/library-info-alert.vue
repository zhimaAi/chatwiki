<style lang="less" scoped>
.library-info-alert {
  .close-btn {
    font-size: 16px;
    color: rgba(0, 0, 0, 0.45);
    cursor: pointer;
  }

  .library-info-content {
    .file-list {
      font-size: 0;
      width: 100%;
      overflow: hidden;
      overflow-x: auto;
      white-space: nowrap;

      .file-item {
        display: inline-block;
        line-height: 22px;
        padding: 5px 16px;
        margin-right: 8px;
        font-size: 14px;
        border-radius: 2px;
        color: #595959;
        border: 1px solid #d9d9d9;
        white-space: nowrap;
        cursor: pointer;
        transition: all 0.2s;

        &.active {
          color: #2475fc;
          border: 1px solid #2475fc;
          background-color: #f5f9ff;
        }
      }
    }

    .document-items {
      .document-item {
        padding: 16px;
        margin-top: 8px;
        border-radius: 2px;
        background-color: #ffffff;
      }

      .document-item-header {
        display: flex;
        justify-content: space-between;
        height: 22px;
        line-height: 22px;
      }

      .document-title {
        font-size: 14px;
        font-weight: 600;
        color: #000000;
        padding-right: 8px;
      }

      .document-size {
        font-size: 14px;
        color: #8c8c8c;
      }

      .document-similarity {
        margin-top: 4px;
        font-size: 14px;
        color: #8c8c8c;
        display: flex;
        align-items: center;
        .svg-action{
          margin-right: 4px;
        }
      }

      .document-content {
        line-height: 22px;
        margin-top: 8px;
        font-size: 14px;
        color: #595959;
        white-space: pre-wrap;
      }
      .fragment-img {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        margin-top: 8px;
        img {
          width: 80px;
          height: 80px;
          border-radius: 6px;
          cursor: pointer;
        }
      }
    }
  }
}
</style>

<template>
  <a-drawer
    class="library-info-alert"
    v-model:open="show"
    :title="t('answer_source')"
    placement="right"
    width="746px"
    :closable="false"
    :bodyStyle="{ background: '#F2F4F7' }"
  >
    <template #extra>
      <span class="close-btn" @click="onClose"><CloseOutlined /></span>
    </template>

    <div class="library-info-content">
      <div class="file-list">
        <div
          class="file-item"
          :class="{ active: activeFileId == file.id }"
          v-for="file in fileList"
          :key="file.id"
          @click="chagenFile(file)"
        >
          {{ file.file_name }}
        </div>
      </div>

      <div class="document-items" ref="docBoxRef">
        <div class="document-item" v-for="item in documentList" :key="item.id">
          <div class="document-item-header">
            <div class="left-box">
              <span class="document-title">{{ t('id_prefix') }}{{ item.id }}</span>
              <span class="document-title" v-if="item.title">{{ item.title }}</span>
              <span class="document-size"> {{ t('total_characters', { count: item.word_total }) }} </span>
            </div>
            <div class="right-box"><a @click="toSource">{{ t('view_source_document') }}</a></div>
          </div>

          <div class="document-item-body">
            <div class="document-similarity">
              {{ t('similarity') }} <svg-icon name="similarity" style="font-size: 16px"></svg-icon
              >{{ item.similarity }}
            </div>
            <div class="document-content" v-if="item.question">{{ t('question_prefix') }}{{ item.question }}</div>
            <div class="document-content" v-if="item.answer">{{ t('answer_prefix') }}{{ item.answer }}</div>
            <div class="document-content">
              {{ item.content }}
            </div>
            <div class="fragment-img" v-viewer>
              <img v-for="(item, index) in item.images" :key="index" :src="item" alt="" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </a-drawer>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { CloseOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getAnswerSource } from '@/api/chat/index'
import {useMathJax} from "@/composables/useMathJax.js";

const { t } = useI18n('views.chat-monitor.components.library-info-alert')

const router = useRouter()
const { renderMath } = useMathJax()

const show = ref(false)
const docBoxRef = ref(null)

const fileList = ref([])
const activeFileId = ref(null)

const reset = () => {
  fileList.value = []
  activeFileId.value = null
}

const open = (files, file) => {
  reset()

  fileList.value = files
  activeFileId.value = file.id
  show.value = true

  getDocumentList(file)
}

const chagenFile = (file) => {
  if (file.id == activeFileId.value) {
    return
  }

  activeFileId.value = file.id

  getDocumentList(file)
}

const documentList = ref([])

const getDocumentList = (file) => {
  getAnswerSource({
    message_id: file.message_id,
    file_id: file.id,
    robot_key: file.robot_key,
    openid: file.openid,
  }).then((res) => {
    documentList.value = res.data || []
    renderMath(docBoxRef.value)
  })
}

// 查看原文档
const toSource = () => {
  router.push('/library/preview?id=' + activeFileId.value)
}

const onClose = () => {
  show.value = false
}

defineExpose({
  open
})
</script>
