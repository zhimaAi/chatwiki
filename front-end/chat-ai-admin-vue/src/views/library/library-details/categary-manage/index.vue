<template>
  <div class="robot-page">
    <div class="search-content-block">
      <div class="selct-star-box">
        <SelectStarBox :hideId="[0]" :startLists="startLists" @change="handleCategoryChange" />
      </div>
      <a-space>
        <a-select
          v-model:value="filterData.status"
          placeholder="请选择"
          style="width: 160px"
          @change="search"
        >
          <a-select-option :value="-1">全部嵌入状态</a-select-option>
          <a-select-option v-for="(i, key) in listStatusMap" :value="key">{{ i }}</a-select-option>
        </a-select>
        <a-button @click="reEmbeddingVectors">重新嵌入向量</a-button>

        <!-- <a-button @click="openEditSubscription({})" type="primary">
          <template #icon>
            <PlusOutlined />
          </template>
          <span>添加分段</span>
        </a-button> -->
      </a-space>
    </div>
    <div class="scroll-box">
      <cu-scroll :scrollbar="{ minSize: 0 }" @onScrollEnd="onScrollEnd">
        <div class="content-block">
          <SubsectionBox
            ref="subsectionBoxRef"
            :isQaLibray="detailsInfo.is_qa_doc"
            :total="total"
            :paragraphLists="paragraphLists"
            @openEditSubscription="openEditSubscription"
            @handleDelParagraph="handleDelParagraph"
            @handleConvert="handleConvert"
            @getStatrList="getStatrList"
          ></SubsectionBox>
        </div>
      </cu-scroll>
    </div>
    <EditSubscription
      :detailsInfo="detailsInfo"
      @handleEdit="handleEditParagraph"
      ref="editSubscriptionRef"
    ></EditSubscription>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import SelectStarBox from '@/views/library/library-preview/components/selct-star-box.vue'
import {
  getLibFileInfo,
  reconstructGraph,
  reconstructCategoryVector,
  getCategoryParagraphList
} from '@/api/library'
import { PlusOutlined, DownOutlined, StarFilled } from '@ant-design/icons-vue'
import SubsectionBox from './components/subsection-box.vue'
import EditSubscription from '@/views/library/library-preview/components/edit-subsection.vue'
const route = useRoute()
const query = route.query

const startLists = ref([])
const paragraphLists = ref([])

const total = ref(0)
const graphStatusMap = {
  0: '待生成',
  4: '生成中',
  2: '生成成功',
  3: '生成失败'
}

const listStatusMap = {
  0: '未转换',
  1: '已转换',
  2: '转换异常',
  3: '转换中...'
}
const detailsInfo = ref({})

const getStatrList = (data) => {
  startLists.value = data
}

const paginations = ref({
  page: 1,
  size: 10
})

const filterData = reactive({
  status: -1,
  category_id: -1,
  library_id: query.id
})

const search = () => {
  paginations.value.page = 1
  total.value = 0
  paragraphLists.value = []
  getParagraphLists()
}
const handleCategoryChange = (id) => {
  filterData.category_id = id
  search()
}
const reEmbeddingVectors = () => {
  reconstructCategoryVector({ id: query.id }).then(() => message.success('操作完成'))
}

const editSubscriptionRef = ref(null)
const openEditSubscription = (data) => {
  editSubscriptionRef.value.showModal(JSON.parse(JSON.stringify(data)))
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
    lastItem.category_id = data.category_id
    lastItem.word_total = data.question.length + data.answer.length + data.content.length
    lastItem.similar_questions = JSON.parse(data.similar_questions)
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
    total.value = --total.value
    console.log(total.value ,'===')
  }
}

const onScrollEnd = () => {
  if (isLoading.value) {
    return
  }
  if (paragraphLists.value.length >= total.value) {
    return
  }
  loadMore()
}

const loadMore = () => {
  paginations.value.page++
  getParagraphLists()
}
const handleConvert = () => {
  updataList()
}

function setIntervalUpStatus() {
  clearInterval(timer)
  // if()
  let list = paragraphLists.value.filter((item) => item.status == 3)
  if (list.length) {
    timer = setInterval(() => {
      updataList()
    }, 5000)
  } else {
    clearInterval(timer)
  }
}

const updataList = () => {
  getCategoryParagraphList({
    page: 1,
    size: paragraphLists.value.length,
    ...filterData
  }).then((res) => {
    let data = res.data
    let list = data.list || []
    list.forEach((item) => {
      item.status_text = listStatusMap[item.status]
      item.graph_status_text = graphStatusMap[item.graph_status]
      if(item.similar_questions){
        item.similar_questions = JSON.parse(item.similar_questions)
      }
    })
    paragraphLists.value = list
    setIntervalUpStatus()
  })
}

let timer = null

const isLoading = ref(false)

const getParagraphLists = () => {
  isLoading.value = true
  getCategoryParagraphList({
    ...paginations.value,
    ...filterData
  }).then((res) => {
    isLoading.value = false
    let data = res.data
    let list = data.list || []
    if (data.info) {
      detailsInfo.value = {
        ...data.info,
        is_qa_doc: data.info.type == 2 ? 1 : 0
      }
    }
    detailsInfo.value = {
      ...data.info,
      is_qa_doc: data.info.type == 2 ? 1 : 0
    }
    list.forEach((item) => {
      item.status_text = listStatusMap[item.status]
      item.graph_status_text = graphStatusMap[item.graph_status]
      if(item.similar_questions){
        item.similar_questions = JSON.parse(item.similar_questions)
      }
    })

    if (paginations.value.page == 1) {
      paragraphLists.value = []
    }
    paragraphLists.value = [...paragraphLists.value, ...list]
    total.value = data.total
    setIntervalUpStatus()
  })
}

getParagraphLists()
</script>

<style lang="less" scoped>
.robot-page {
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.search-content-block {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  justify-content: space-between;
}
.scroll-box {
  flex: 1;
  overflow: hidden;
}

.list-box {
  height: 999px;
}

.empty-box {
  text-align: center;
  height: 100%;
  padding-top: 148px;
  img {
    width: 200px;
    height: 200px;
  }
  .title {
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    color: #262626;
  }
}
</style>
