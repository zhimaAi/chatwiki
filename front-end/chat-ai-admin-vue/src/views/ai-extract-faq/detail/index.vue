<template>
  <div class="main-content-box">
    <cu-scroll
      style="padding-right: 48px"
      :scrollbar="{ minSize: 80, fade: false, interactive: true, scrollbarTrackClickable: true }"
    >
      <div class="breadcrumb-block">
        <a-breadcrumb>
          <a-breadcrumb-item>
            <router-link to="/ai-extract-faq/list"
              ><LeftOutlined /> 文档提取FAQ
            </router-link></a-breadcrumb-item
          >
          <a-breadcrumb-item>{{ file_name }}</a-breadcrumb-item>
        </a-breadcrumb>
      </div>
      <page-alert class="mb16" title="使用说明">
        <div>
          <p>
            您可以查看、编辑提取处的问答对，并直接将问答对导入到知识库中。您也可以将知识库下载成文档，批量导入导问答知识库中。
          </p>
        </div>
      </page-alert>
      <div class="btn-wrapper-box">
        <a-segmented @change="onSearch" v-model:value="is_import" :options="statusOptions">
          <template #label="{ payload = {} }">
            <div>{{ payload.title }}（{{ payload.num }}）</div>
          </template>
        </a-segmented>
        <a-flex :gap="8">
          <a-button
            @click="handleOpenImportModal"
            type="primary"
            :icon="createVNode(UploadOutlined)"
            >导入知识库</a-button
          >
          <a-button @click="handleDel">批量删除</a-button>
          <a-dropdown>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="handleDownload('docx')" key="1">下载为docx</a-menu-item>
                <a-menu-item @click="handleDownload('xlsx')" key="2">下载为xlsx</a-menu-item>
              </a-menu>
            </template>
            <a-button type="primary">
              下载
              <DownOutlined />
            </a-button>
          </a-dropdown>
        </a-flex>
      </div>
      <div class="content-block">
        <SubsectionBox
          ref="subsectionBoxRef"
          :total="total"
          :isLoading="isLoading"
          :paragraphLists="paragraphLists"
          @openEditSubscription="openEditSubscription"
          @handleDelParagraph="getParagraphLists"
        ></SubsectionBox>
        <div class="pagination-box">
          <a-pagination
            v-model:current="paginations.page"
            v-model:page-size="paginations.size"
            :total="total"
            :pageSizeOptions="['100', '200', '500', '1000']"
            show-size-changer
            @change="onShowSizeChange"
          >
          </a-pagination>
        </div>
      </div>
    </cu-scroll>
  </div>
  <EditSubscription ref="editSubscriptionRef" @ok="getParagraphLists" />
  <ImportKnowledgeModal ref="importKnowledgeModalRef" @ok="onSearch" />
</template>

<script setup>
import PageAlert from '@/components/page-alert/page-alert.vue'
import {
  UploadOutlined,
  LeftOutlined,
  ExclamationCircleOutlined,
  DownOutlined
} from '@ant-design/icons-vue'
import { getFAQFileQAList, deleteFAQFileQA } from '@/api/library'
import SubsectionBox from './components/subsection-box.vue'
import EditSubscription from './components/edit-subsection.vue'
import { ref, createVNode } from 'vue'
import ImportKnowledgeModal from '../list/components/import-knowledge-modal.vue'
import { useRoute } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { useUserStore } from '@/stores/modules/user'
const userStore = useUserStore()

const route = useRoute()
const query = route.query
const file_name = query.file_name

const is_import = ref('')
const subsectionBoxRef = ref(null)

const total = ref(0)
const isLoading = ref(false)

const paragraphLists = ref([])

const paginations = ref({
  page: 1,
  size: 100
})

const editSubscriptionRef = ref(null)
const openEditSubscription = (data) => {
  editSubscriptionRef.value.show(JSON.parse(JSON.stringify(data)))
}

const onShowSizeChange = (current, pageSize) => {
  paginations.value.page = current
  paginations.value.size = pageSize
  getParagraphLists()
}

const onSearch = () => {
  paginations.value.page = 1
  getParagraphLists()
}

const getParagraphLists = () => {
  isLoading.value = true
  getFAQFileQAList({
    ...paginations.value,
    id: query.id,
    is_import: is_import.value
  })
    .then((res) => {
      let list = res.data.list || []
      list.forEach((item) => {
        if (item.images) {
          item.images = JSON.parse(item.images)
        }
      })
      paragraphLists.value = list
      total.value = res.data.total
      setStatusOption(res.data.total, res.data.import_total)
      subsectionBoxRef.value.resetSelect()
    })
    .finally(() => {
      isLoading.value = false
    })
}

getParagraphLists()

const importKnowledgeModalRef = ref(null)

const handleOpenImportModal = () => {
  let selectedRowKeys = subsectionBoxRef.value.state.selectedRowKeys
  if (selectedRowKeys.length == 0) {
    return message.error('请选择你要导入的问答')
  }
  importKnowledgeModalRef.value.show({
    id: query.id,
    ids: selectedRowKeys.join(',')
  })
}

const handleDel = () => {
  let selectedRowKeys = subsectionBoxRef.value.state.selectedRowKeys
  if (selectedRowKeys.length == 0) {
    return message.error('请选择你要删除的问答')
  }
  Modal.confirm({
    title: '提示',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确认是否删除选中的分段?',
    onOk() {
      deleteFAQFileQA({ ids: selectedRowKeys.join(',') }).then((res) => {
        message.success('删除成功')
        getParagraphLists()
      })
    },
    onCancel() {}
  })
}

const handleDownload = (type) => {
  let TOKEN = userStore.getToken
  let srcUrl = `/manage/exportFAQFileAllQA?id=${query.id}&token=${TOKEN}&ext=${type}`
  window.location.href = srcUrl
}

const statusOptions = ref([])
const setStatusOption = (total, import_total) => {
  statusOptions.value = [
    {
      value: '',
      payload: {
        num: total,
        title: '全部'
      }
    },
    {
      value: 0,
      payload: {
        num: total - import_total,
        title: '未导入'
      }
    },
    {
      value: 1,
      payload: {
        num: import_total,
        title: '已导入'
      }
    }
  ]
}
</script>

<style lang="less" scoped>
.main-content-box {
  height: 100%;
  overflow: hidden;
  padding: 16px 0 16px 48px;
}
.breadcrumb-block {
  height: 30px;
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}
.mb16 {
  margin-bottom: 16px;
}
.btn-wrapper-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  &::v-deep(.ant-segmented) {
    background: #e4e6eb;
    .ant-segmented-item {
      color: #262626;
    }
    .ant-segmented-item-selected {
      color: #2475fc;
    }
  }
}

.content-block {
  margin-top: 8px;
}
.pagination-box {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
