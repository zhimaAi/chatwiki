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
      .right-box{
        display: flex;
        align-items: center;
        gap: 8px;
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

.iframe-box{
  width: 100%;
  height: 620px;
}
</style>

<template>
  <a-drawer
    class="library-info-alert"
    v-model:open="show"
    title="答案来源"
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
          <span v-if="file.file_name">{{ file.file_name }}</span>
          <span v-else>{{ file.library_name }}-精选</span>
        </div>
      </div>

      <div class="document-items">
        <div class="document-item" v-for="item in documentList" :key="item.id">
          <div class="document-item-header">
            <div class="left-box">
              <span class="document-title">id：{{ item.id }}</span>
              <span class="document-title" v-if="item.title">{{ item.title }}</span>
              <span class="document-size"> 共{{ item.word_total }}个字符 </span>
            </div>
            <div class="right-box">
              <a @click="viewSourceFile(item)" v-if="item.page_num > 0">预览原文件&gt;</a>
              <a @click="toSource">查看源文档&gt;</a>
            </div>
          </div>

          <div class="document-item-body">
            <div class="document-similarity">
              相似度： <svg-icon name="similarity" style="font-size: 16px"></svg-icon
              >{{ item.similarity }}
            </div>
            <div class="document-content" v-if="item.question">Q：{{ item.question }}</div>
            <div class="document-content" v-if="item.answer">A：{{ item.answer }}</div>
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
  <a-modal v-model:open="viewOpen" :title="viewTitle" :footer="null" :width="740">
    <div class="iframe-box">
      <CuScroll>
        <VuePdfEmbed :source="sourceUrl"  />
      </CuScroll>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, computed} from 'vue'
import { useRouter } from 'vue-router'
import { CloseOutlined } from '@ant-design/icons-vue'
import { getAnswerSource } from '@/api/chat/index'
import VuePdfEmbed from 'vue-pdf-embed'
import { useChatStore } from '@/stores/modules/chat'
const chatStore = useChatStore()
const { robot, user } = chatStore
const router = useRouter()

const show = ref(false)

const fileList = ref([])
const activeFileId = ref(null)

const reset = () => {
  fileList.value = []
  activeFileId.value = null
}

const open = (files, file) => {
  reset()
  fileList.value = files.map(item => {
    let answer_source_data = null;
    if(item.answer_source_data){
      answer_source_data = JSON.parse(item.answer_source_data)
    }
    return {
      ...item,
      answer_source_data,
    }
  })
  activeFileId.value = file.id
  show.value = true
  console.log(files, file)
  if(file.answer_source_data){
    documentList.value = JSON.parse(file.answer_source_data) || []
    return
  }
  getDocumentList(file)
}

const chagenFile = (file) => {
  if (file.id == activeFileId.value) {
    return
  }

  activeFileId.value = file.id

  if(file.answer_source_data){
    documentList.value = file.answer_source_data || []
    return
  }

  getDocumentList(file)
}

const documentList = ref([])

const getDocumentList = (file) => {
  getAnswerSource({
    message_id: file.message_id,
    file_id: file.id,
    robot_key: robot.robot_key,
    openid: robot.openid,
  }).then((res) => {
    documentList.value = res.data || []
  })
}

// 查看原文档
const toSource = () => {
  let currentItem = fileList.value.filter(item => item.id == activeFileId.value)[0] || {}
  if(currentItem.file_name){
    window.open(`#/library/preview?id=${activeFileId.value}`)
  }else{
    window.open(`#/library/details/categary-manage?id=${currentItem.library_id}`)
  }
}

const viewTitle = computed(()=>{
  let currentItem = fileList.value.filter(item => item.id == activeFileId.value)[0] || {}
  if(currentItem.file_name){
    return currentItem.file_name
  }
  return currentItem.library_name + '-精选'
})
const sourceUrl = ref('')
const viewOpen = ref(false)
const viewSourceFile = (item) => {
  console.log(robot)
  viewOpen.value = true;
  sourceUrl.value = '/manage/getLibRawFileOnePage?id=' + activeFileId.value + '&page=' + item.page_num + '&admin_user_id=' + user.admin_user_id
}

const onClose = () => {
  show.value = false
}

defineExpose({
  open
})
</script>
