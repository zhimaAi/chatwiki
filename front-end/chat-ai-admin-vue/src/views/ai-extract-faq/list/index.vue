<template>
  <div class="library-page">
    <PageTabs class="mb-16" :tabs="pageTabs" active="/ai-extract-faq/list">
    </PageTabs>
    <page-alert style="margin-bottom: 16px" :title="t('usage_title')">
      <div>
        <p>
          {{ t('usage_desc_1') }}
        </p>
        <p>
          {{ t('usage_desc_2') }}
        </p>
      </div>
    </page-alert>
    <div class="library-page-body">
      <div class="btn-block" v-if="uploadDocFaq">
        <a-button @click="handleOpenAddModal" type="primary" :icon="createVNode(UploadOutlined)"
          >{{ t('upload_btn') }}</a-button
        >
      </div>
      <div class="list-box">
        <a-table
          :data-source="list"
          :loading="loading"
          :pagination="{
            current: pager.page,
            total: pager.total,
            pageSize: pager.size,
            showQuickJumper: true,
            showSizeChanger: true,
            pageSizeOptions: ['10', '20', '50', '100']
          }"
          @change="onTableChange"
          :scroll="{ x: 800 }"
        >
          <a-table-column key="file_name" data-index="file_name" :title="t('original_doc')" :width="200">
            <template #default="{ record }">
              <a @click="toDetail(record)">{{ record.file_name }}</a>
            </template>
          </a-table-column>
          <a-table-column key="status" :title="t('status')" :width="140">
            <template #default="{ record }">
              <div class="status-block status-bule" v-if="record.status == 0">
                <LoadingOutlined />{{ t('status_queuing') }}
              </div>
              <div class="status-block status-bule" v-if="record.status == 1">
                <ClockCircleFilled />{{ t('status_parsing') }}
              </div>
              <div class="status-block status-bule" v-if="record.status == 2">
                <ClockCircleFilled />{{ t('status_extracting') }}
              </div>
              <div class="status-block status-green" v-if="record.status == 3">
                <CheckCircleFilled />{{ t('status_completed') }}
              </div>
              <div class="status-block status-red" v-if="record.status == 4">
                <CloseCircleFilled />{{ t('status_failed') }}
              </div>
            </template>
          </a-table-column>
          <a-table-column key="chunk_size" :title="t('chunk_method')" :width="200">
            <template #default="{ record }">
              <div v-if="record.chunk_type == 1">{{ t('by_length') }}{{ record.chunk_size }}</div>
              <div v-if="record.chunk_type == 2">
                <div>{{ t('by_separator') }}{{ record.separators_no_desc }}</div>
                <div>{{ t('max_length') }}{{ record.chunk_size }}</div>
              </div>
            </template>
          </a-table-column>
          <a-table-column key="count" :title="t('total_chunks')" :width="110">
            <template #default="{ record }">
              {{ +record.success_count + +record.fail_count }}
            </template>
          </a-table-column>
          <a-table-column key="success_count" :title="t('success_chunks')" :width="140">
            <template #default="{ record }"> {{ record.success_count }}</template>
          </a-table-column>
          <a-table-column key="fail_count" :title="t('failed_chunks')" :width="140">
            <template #default="{ record }">
              <a-flex :gap="12">
                <a @click="handleOpenFailDetail(record)" v-if="record.fail_count > 0">{{
                  record.fail_count
                }}</a>
                <span v-else>{{ record.fail_count }}</span>
                <a v-if="record.status == 3 || record.status == 4">
                  <SyncOutlined v-if="record.fail_count > 0" @click="handleReSync(record)" />
                </a>
              </a-flex>
            </template>
          </a-table-column>
          <a-table-column key="qa_count" :title="t('qa_count')" :width="140">
            <template #default="{ record }">
              <a @click="toDetail(record)">{{ record.qa_count }}</a>
            </template>
          </a-table-column>
          <a-table-column key="create_time_desc" :title="t('upload_time')" :width="150">
            <template #default="{ record }"> {{ record.create_time_desc }} </template>
          </a-table-column>
          <a-table-column key="key7" :title="t('action')" :width="190" fixed="right">
            <template #default="{ record }">
              <a-flex :gap="16" wrap="wrap">
                <a-dropdown v-if="record.status == 3">
                  <template #overlay>
                    <a-menu>
                      <a-menu-item @click="handleDownload('docx', record.id)" key="1"
                        >{{ t('download_as_docx') }}</a-menu-item
                      >
                      <a-menu-item @click="handleDownload('xlsx', record.id)" key="2"
                        >{{ t('download_as_xlsx') }}</a-menu-item
                      >
                    </a-menu>
                  </template>
                  <a-button size="small" type="link" style="padding: 0">{{ t('download') }}</a-button>
                </a-dropdown>
                <a-button v-else disabled size="small" type="link" style="padding: 0"
                  >{{ t('download') }}</a-button
                >

                <a-button
                  @click="handleOpenImportModal(record)"
                  :disabled="record.status != 3"
                  size="small"
                  type="link"
                  style="padding: 0"
                  >{{ t('import_kb') }}</a-button
                >
                <a-button @click="handleDelete(record)" size="small" type="link" style="padding: 0"
                  >{{ t('delete') }}</a-button
                >
              </a-flex>
            </template>
          </a-table-column>
        </a-table>
      </div>
    </div>
  </div>
  <UploadDocument :separatorsOptions="separatorsOptions" ref="uploadDocumentRef" @ok="onSearch" />
  <ImportKnowledgeModal ref="importKnowledgeModalRef" @ok="getList" />
  <FailDetail ref="failDetailRef" />
</template>
<script setup>
import { ref, createVNode, reactive, onUnmounted, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import {
  PlusOutlined,
  LoadingOutlined,
  ClockCircleFilled,
  CheckCircleFilled,
  CloseCircleFilled,
  SyncOutlined,
  ExclamationCircleOutlined,
  UploadOutlined
} from '@ant-design/icons-vue'
import {
  getFAQFileList,
  deleteFAQFile,
  renewFAQFileData,
  getFAQFileInfo,
  getSeparatorsList
} from '@/api/library'
import PageTabs from '@/components/cu-tabs/page-tabs.vue'
import PageAlert from '@/components/page-alert/page-alert.vue'
import UploadDocument from './components/upload-document.vue'
import ImportKnowledgeModal from './components/import-knowledge-modal.vue'
import FailDetail from './components/fail-detail.vue'
import dayjs from 'dayjs'
import { formatSeparatorsNo } from '@/utils/index'
import { useUserStore } from '@/stores/modules/user'
import { usePermissionStore } from '@/stores/modules/permission'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.ai-extract-faq.list.index')
const userStore = useUserStore()
const router = useRouter()

const permissionStore = usePermissionStore()
let { role_permission, role_type } = permissionStore
const uploadDocFaq = computed(() => role_type == 1 || role_permission.includes('UploadDocFaq'))

const pageTabs = ref([
  {
    title: t('tab_library'),
    path: '/library/list'
  },
  {
    title: t('tab_database'),
    path: '/database/list'
  },
  {
    title: t('tab_extract_faq'),
    path: '/ai-extract-faq/list'
  },
  {
    title: t('tab_trigger_stats'),
    path: '/trigger-statics/list'
  }
])

const pager = reactive({
  page: 1,
  size: 20,
  total: 0
})

const list = ref([])

const loading = ref(false)

function formateSeparatorDesc(data) {
  if(!Array.isArray(data)){
    return []
  }
  return data.map((item) => {
    let findItem = separatorsOptions.value.find((it) => it.no === item)
    if (findItem) {
      return findItem.name
    } else {
      return item
    }
  })
}
const getList = () => {
  loading.value = false
  getFAQFileList({
    ...pager
  })
    .then((res) => {
      let datas = res.data.list || []
      console.log(separatorsOptions.value)
      datas = datas.map((item) => {
        let create_time_desc = dayjs(item.create_time * 1000).format('YY-MM-DD HH:mm')
        let separators_no = formatSeparatorsNo(item.separators_no)
        return {
          ...item,
          create_time_desc,
          separators_no_desc: formateSeparatorDesc(separators_no).join('、'),
        }
      })
      pager.total = +res.data.total
      list.value = datas
      startPollingForStatusTwo()
    })
    .finally(() => {
      loading.value = false
    })
}

const onSearch = () => {
  pager.page = 1
  getList()
}

const onTableChange = (pagination) => {
  pager.page = pagination.current
  pager.size = pagination.pageSize
  getList()
}

let pollingIntervals = {}

const startPollingForStatusTwo = () => {
  list.value.forEach((item) => {
    if (item.status == 2) {
      startPolling(item.id)
    }
  })
}
// 开始轮询单个ID的状态
function startPolling(id) {
  // 如果已经有这个ID的定时器，先清除
  if (pollingIntervals[id]) {
    clearInterval(pollingIntervals[id])
  }

  // 设置新的定时器
  pollingIntervals[id] = setInterval(async () => {
    try {
      const response = await getFAQFileInfo({ id })
      const updatedItem = response.data
      updatedItem.create_time_desc = dayjs(updatedItem.create_time * 1000).format('YY-MM-DD HH:mm')
      let separators_no = formatSeparatorsNo(updatedItem.separators_no)
      updatedItem.separators_no_desc = formateSeparatorDesc(separators_no).join('、')
      // 找到列表中对应的项并更新
      const index = list.value.findIndex((item) => item.id == id)
      if (index !== -1) {
        list.value.splice(index, 1, updatedItem)
      }

      // 如果状态不再是2，停止轮询
      if (updatedItem.status != 2) {
        stopPolling(id)
      }
    } catch (error) {
      // 出错时也可以考虑停止轮询或重试
      stopPolling(id)
    }
  }, 5000) // 5秒间隔
}

function stopPolling(id) {
  if (pollingIntervals[id]) {
    clearInterval(pollingIntervals[id])
    delete pollingIntervals[id]
  }
}

const handleDelete = (data) => {
  Modal.confirm({
    title: t('delete_confirm_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('delete_confirm_content', { file_name: data.file_name }),
    okText: t('confirm'),
    cancelText: t('cancel'),
    okType: 'danger',
    onOk() {
      onDelete(data)
    },
    onCancel() {}
  })
}

const onDelete = ({ id }) => {
  deleteFAQFile({ id }).then(() => {
    message.success(t('delete_success'))
    getList()
  })
}

const handleReSync = (record) => {
  Modal.confirm({
    title: t('resync_confirm_title'),
    icon: null,
    okText: t('confirm'),
    cancelText: t('cancel'),
    onOk() {
      renewFAQFileData({
        id: record.id
      }).then((res) => {
        message.success(t('resync_success'))
        getList()
      })
    },
    onCancel() {}
  })
}

const toDetail = (record) => {
  if (record.status == 2 || record.status == 3) {
    router.push({
      path: '/ai-extract-faq/details',
      query: {
        id: record.id,
        file_name: record.file_name
      }
    })
    return
  }
  let tips = t('cannot_view_detail')
  if (record.status == 0) {
    tips = t('status_queuing_tip')
  }
  if (record.status == 1) {
    tips = t('status_parsing_tip')
  }
  if (record.status == 4) {
    tips = t('status_failed_tip')
  }
  message.error(tips)
}

const uploadDocumentRef = ref(null)
const handleOpenAddModal = () => {
  uploadDocumentRef.value.show()
}

const importKnowledgeModalRef = ref(null)

const handleOpenImportModal = (record) => {
  importKnowledgeModalRef.value.show({
    ...record
  })
}

const handleDownload = (type, id) => {
  let TOKEN = userStore.getToken
  let srcUrl = `/manage/exportFAQFileAllQA?id=${id}&token=${TOKEN}&ext=${type}`
  window.location.href = srcUrl
}

const failDetailRef = ref(null)
const handleOpenFailDetail = (record) => {
  failDetailRef.value.show({
    ...record
  })
}

// 分段标识符列表
const separatorsOptions = ref([])

const fetchSeparatorsOptions = async () => {
  getSeparatorsList().then((res) => {
    separatorsOptions.value = res.data || []
  })
}

onMounted(async () => {
  await fetchSeparatorsOptions()
  getList()
})

onUnmounted(() => {
  Object.keys(pollingIntervals).forEach((id) => {
    clearInterval(pollingIntervals[id])
  })
  pollingIntervals = {}
})
</script>

<style lang="less" scoped>
.library-page {
  .list-box {
    margin-top: 8px;
  }
}

.status-block {
  display: flex;
  align-items: center;
  width: fit-content;
  min-width: 100px;
  gap: 3px;
  padding: 0 6px;
  height: 22px;
  border-radius: 6px;
  line-height: 22px;
  font-size: 14px;
  &.status-bule {
    background: #d4e3fc;
    color: #2475fc;
  }
  &.status-green {
    background: #cafce4;
    color: #21a665;
  }
  &.status-red {
    background: #fbddde;
    color: #fb363f;
  }
}
</style>
