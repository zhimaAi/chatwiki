<template>
  <div class="library-preview-page">
    <div class="document-title-block">
      <LeftOutlined @click="goBack" />
      <span class="title">{{ detailsInfo.file_name }}</span>
      <SegmentationMode :detailsInfo="detailsInfo"></SegmentationMode>
      <a-button @click="openEditSubscription({})" style="margin-left: auto" type="primary">
        <template #icon>
          <PlusOutlined />
        </template>
        <span>添加分段</span>
      </a-button>
    </div>
    <template v-if="pageLoading">
      <div class="page-loading-box">
        <a-spin />
      </div>
    </template>
    <template v-else>
      <template v-if="isTableType || doc_type == 3">
        <!-- 表格类型 不支持预览 -->
        <Empty @openEditSubscription="openEditSubscription" v-if="isEmpty"></Empty>
        <cu-scroll v-else @onScrollEnd="onScrollEnd">
          <div class="content-block">
            <SubsectionBox
              ref="subsectionBoxRef"
              :total="total"
              :paragraphLists="paragraphLists"
              :detailsInfo="detailsInfo"
              @openEditSubscription="openEditSubscription"
              @handleDelParagraph="handleDelParagraph"
            ></SubsectionBox>
          </div>
        </cu-scroll>
      </template>
      <div class="pdf-view-box" v-else>
        <div class="view-content-wrap">
          <cu-scroll class="scroll-box" @onScrollEnd="onScrollEnd">
            <div class="view-content-box">
              <div
                class="list-item"
                @click="handleClickItem(index)"
                v-for="(item, index) in paragraphLists"
                :key="index"
              >
                <div class="content-box" v-if="item.question">{{ item.question }}</div>
                <div class="content-box" v-if="item.answer">{{ item.answer }}</div>
                <div class="content-box" v-html="item.content"></div>
                <div class="fragment-img" v-viewer>
                  <img v-for="(item, index) in item.images" :key="index" :src="item" alt="" />
                </div>
              </div>
            </div>
          </cu-scroll>
        </div>
        <div class="right-contnt-box">
          <Empty @openEditSubscription="openEditSubscription" v-if="isEmpty"></Empty>
          <cu-scroll ref="rightScrollRef" v-else @onScrollEnd="onScrollEnd">
            <div class="content-block">
              <SubsectionBox
                ref="subsectionBoxRef"
                :total="total"
                :detailsInfo="detailsInfo"
                :paragraphLists="paragraphLists"
                @openEditSubscription="openEditSubscription"
                @handleDelParagraph="handleDelParagraph"
              ></SubsectionBox>
            </div>
          </cu-scroll>
        </div>
      </div>
    </template>
    <EditSubscription
      :detailsInfo="detailsInfo"
      @handleEdit="handleEditParagraph"
      ref="editSubscriptionRef"
    ></EditSubscription>
  </div>
</template>

<script setup>
import { ref, computed, defineAsyncComponent, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { LeftOutlined, PlusOutlined } from '@ant-design/icons-vue'
import Empty from './components/empty.vue'
import SubsectionBox from './components/subsection-box.vue'
import { getLibFileInfo, getParagraphList } from '@/api/library'
import CuScroll from '@/components/cu-scroll/cu-scroll.vue'
import EditSubscription from './components/edit-subsection.vue'
import SegmentationMode from './components/segmentation-mode.vue'

const route = useRoute()
const router = useRouter()
const isTableType = computed(() => {
  return detailsInfo.value.is_table_file == '1'
})
const doc_type = computed(() => {
  return detailsInfo.value.doc_type
})
const isEmpty = computed(() => {
  // 1学习中,2学习完成,3文件异常 4待用户切分 待学习
  let status = detailsInfo.value.status
  return !(status == '1' || status == '2') || paragraphLists.value.length == 0
})

const detailsInfo = ref({
  pdf_url: ''
})

const goBack = () => {
  router.back()
}

const rightScrollRef = ref(null)
let lastIndex = -1
const handleClickItem = (index) => {
  let el = document.querySelectorAll('.subsection-box .list-item')[index]
  el.classList.add('flash-border')
  if (lastIndex >= 0) {
    document
      .querySelectorAll('.subsection-box .list-item')
      [lastIndex].classList.remove('flash-border')
  }
  setTimeout(() => {
    lastIndex = -1
    el.classList.remove('flash-border')
  }, 2500)
  lastIndex = index
  rightScrollRef.value.scrollToElement({
    el,
    time: 1000,
  })
}

const pageLoading = ref(true)

const getFileInfo = async () => {
  getLibFileInfo({
    id: route.query.id
  }).then((res) => {
    detailsInfo.value = res.data
    let status = res.data.status
    if (status == 1 || status == 2) {
      getParagraphLists()
    } else {
      pageLoading.value = false
    }
  })
}
getFileInfo()

const paginations = ref({
  page: 1,
  size: 10
})
const paragraphInfo = ref({})
const paragraphLists = ref([])
const total = ref(0)
const listStatusMap = {
  0: '未转换',
  1: '已转换',
  2: '转换异常'
}

let isLoading = true

const getParagraphLists = () => {
  isLoading = true
  getParagraphList({
    file_id: route.query.id,
    ...paginations.value
  }).then((res) => {
    pageLoading.value = false
    isLoading = false
    let data = res.data
    let list = data.list || []
    let isExcelQa = data.info?.is_qa_doc == '1' && data.info?.is_table_file == '1'
    list.forEach((item) => {
      item.status_text = listStatusMap[item.status]
      item.isExcelQa = isExcelQa
    })

    paragraphInfo.value = data.info
    if (paginations.value.page == 1) {
      paragraphLists.value = []
    }
    paragraphLists.value = [...paragraphLists.value, ...list]
    total.value = data.total
  })
}
const onScrollEnd = () => {
  if (isLoading) {
    return
  }
  if (paragraphLists.value.length >= total.value) {
    return
  }
  paginations.value.page++
  getParagraphLists()
}
const handleEditParagraph = (data) => {
  if (!data.id) {
    paginations.value.page = 1
    getParagraphLists()
    return
  }
  // 更新分段内容 无刷更新
  let lists = paragraphLists.value
  let index = lists.findIndex((item) => item.id == data.id)
  if (index > -1) {
    let lastItem = lists[index]
    lastItem.title = data.title
    lastItem.content = data.content
    lastItem.question = data.question
    lastItem.answer = data.answer
    lastItem.answer = data.answer
    lastItem.images = data.images
    lastItem.word_total = data.question.length + data.answer.length + data.content.length
    lists.splice(index, 1, lastItem)
    paragraphLists.value = lists
  }
}
const handleDelParagraph = (id) => {
  // 无刷新 删除列表
  let lists = paragraphLists.value
  let index = lists.findIndex((item) => item.id == id)
  if (index > -1) {
    lists.splice(index, 1)
    paragraphLists.value = lists
    total.value = total.value--
  }
}

const pdfIsLoading = ref(true)

const handleLoad = (e) => {
  pdfIsLoading.value = false
}

const editSubscriptionRef = ref(null)
const openEditSubscription = (data) => {
  editSubscriptionRef.value.showModal(JSON.parse(JSON.stringify(data)))
}
</script>

<style lang="less" scoped>
.library-preview-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  .page-loading-box {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .document-title-block {
    margin-bottom: 24px;
    width: 100%;
    height: 38px;
    border-radius: 2px;
    background: #f2f4f7;
    display: flex;
    align-items: center;
    color: #262626;
    padding: 0 16px;
    .anticon-left {
      cursor: pointer;
    }
    .title {
      font-weight: 600;
      font-size: 14px;
      line-height: 22px;
      margin-left: 8px;
    }
  }

  .pdf-view-box {
    flex: 1;
    display: flex;
    overflow: hidden;
    .view-content-wrap {
      flex: 1;
      overflow: hidden;
      .scroll-box {
        margin-right: 24px;
        border: 1px solid #ccc;
        border-radius: 4px;
      }
      .view-content-box {
        padding: 12px 32px 24px;
        .list-item {
          cursor: pointer;
          padding: 6px 8px;
          border-radius: 4px;
          border: 1px solid #fff;
          transition: all .3s;
          &:hover {
            background: #f7f7f7;
            border: 1px solid #E0E0E0;
          }
        }
        .content-box {
          color: #333;
          font-size: 14px;
          font-weight: 400;
          line-height: 22px;
          margin-top: 4px;
          white-space: pre-wrap;
          word-wrap: break-word;
        }
        .fragment-img {
          display: flex;
          gap: 8px;
          align-items: baseline;
          flex-wrap: wrap;
          img {
            max-width: 100%;
            max-height: 100%;
            width: auto;
            height: fit-content;
          }
        }
      }
    }
    .right-contnt-box {
      flex: 1;
      overflow: hidden;
    }
  }
}
</style>
