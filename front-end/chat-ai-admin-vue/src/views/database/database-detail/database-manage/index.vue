<template>
  <div class="loading-box" v-if="pageLoading">
    <a-spin></a-spin>
  </div>
  <div class="empty-box" v-if="!pageLoading && !column.length">
    <img src="@/assets/img/library/preview/empty.png" alt="" />
    <div>{{ t('empty_no_fields') }}</div>
  </div>
  <div class="field-manage-page" v-if="!pageLoading && column.length">
    <div class="page-title">
      {{ t('title') }}
      <a-flex :gap="16">
        <a-button type="primary" @click="handleAddData()">
          <template #icon>
            <PlusOutlined />
          </template>
          {{ t('btn_add_data') }}
        </a-button>
        <a-button @click="handleImportData()">
          <template #icon>
            <DownloadOutlined />
          </template>
          {{ t('btn_import_data') }}
        </a-button>
        <a-dropdown>
          <template #overlay>
            <a-menu>
              <a-menu-item key="1"><div @click="handleOpenExportModal">{{ t('menu_export_data') }}</div></a-menu-item>
              <a-menu-item key="2"><div @click="handleClearData">{{ t('menu_clear_data') }}</div></a-menu-item>
            </a-menu>
          </template>
          <a-button>
            {{ t('btn_more_actions') }}
            <DownOutlined />
          </a-button>
        </a-dropdown>
      </a-flex>
    </div>
    <a-tabs
      class="tabs-wrapper"
      v-model:activeKey="activeKey"
      type="editable-card"
      @edit="onEditTabs"
      @change="handleChangeTab"
    >
      <a-tab-pane key="" :closable="false">
        <template #tab>{{ t('tab_all', { count: allCount }) }}</template>
      </a-tab-pane>
      <a-tab-pane v-for="pane in panes" :key="pane.id" :closable="false">
        <template #tab>
          <span v-if="pane.name.length > 8">
            <a-tooltip>
              <template #title>{{ pane.name }}</template>
              {{ pane.name.slice(0, 8) + '...' }}
            </a-tooltip>
          </span>
          <span v-else>{{ pane.name }} </span>
          ({{ pane.entry_count }})</template
        >
      </a-tab-pane>
      <template #addIcon> <PlusOutlined /> {{ t('tab_add_category') }} </template>
      <template #rightExtra>
        <span class="setting-btn" @click="handleOpenManageModal">
          <SettingOutlined />
          {{ t('btn_manage_category') }}
        </span>
      </template>
    </a-tabs>
    <div class="table-wrapper customize-scroll-style">
      <a-table
        sticky
        :data-source="data"
        :loading="loading"
        :scroll="{ x: 1000 }"
        :pagination="{
          current: queryParams.page,
          total: queryParams.total,
          pageSize: queryParams.size,
          showQuickJumper: true,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100']
        }"
        @change="onTableChange"
      >
        <a-table-column key="id" data-index="id" :title="t('column_id_title')" :width="60">
        </a-table-column>
        <a-table-column
          v-for="item in column"
          :key="item.name"
          :data-index="item.name"
          :width="165"
        >
          <template #title>
            <span v-if="item.name.length > 12">
              <a-tooltip>
                <template #title>{{ item.name }}</template>
                {{ item.name.slice(0, 12) + '...' }}
              </a-tooltip>
            </span>
            <span v-else>{{ item.name }} </span>
            <a-tooltip>
              <template #title>{{ item.description }}</template>
              <QuestionCircleOutlined style="margin-left: 2px" />
            </a-tooltip>
          </template>
          <template #default="{ record }">
            <span v-if="record[item.name].length > 10">
              <a-tooltip>
                <template #title>{{ record[item.name] }}</template>
                {{ record[item.name].slice(0, 9) + '...' }}
              </a-tooltip>
            </span>
            <span v-else>{{ record[item.name] || '-' }} </span>
          </template>
        </a-table-column>
        <a-table-column key="action" :title="t('column_action_title')" :width="120" fixed="right">
          <template #default="{ record }">
            <span>
              <a @click="handleAddData(record)">{{ t('action_edit') }}</a>
              <a-divider type="vertical" />
              <a @click="onDelete(record)">{{ t('action_delete') }}</a>
            </span>
          </template>
        </a-table-column>
      </a-table>
    </div>

    <AddDataModal @ok="handleManageOkBack" :column="column" ref="addDataModalRef"></AddDataModal>
    <AddFilrerModal
      @ok="handleManageOkBack"
      :column="column"
      ref="addFilrerModalRef"
    ></AddFilrerModal>
    <ExportModal ref="exportModalRef"></ExportModal>
    <FilterManageModal
      @ok="getData"
      @change="getSortLists"
      :column="column"
      ref="filterManageModalRef"
    ></FilterManageModal>
    <ImportDataModal @ok="onSearch" :column="column" ref="importDataModalRef" />
  </div>
</template>

<script setup>
import {
  PlusOutlined,
  QuestionCircleOutlined,
  DownOutlined,
  SettingOutlined,
  ExclamationCircleOutlined,
  DownloadOutlined
} from '@ant-design/icons-vue'
import { ref, reactive, createVNode, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import AddDataModal from './components/add-data-modal.vue'
import ExportModal from './components/export-modal.vue'
import AddFilrerModal from './components/add-filter-modal.vue'
import FilterManageModal from './components/filter-manage-modal.vue'
import { useDatabaseStore } from '@/stores/modules/database'
import ImportDataModal from './components/import-data-modal.vue'
import {
  getFormFieldList,
  getFormEntryList,
  delFormEntry,
  emptyFormEntry,
  getFormFilterList
} from '@/api/database'

const { t } = useI18n('views.database.database-detail.database-manage.index')

const databaseStore = useDatabaseStore()
const allCount = computed(() => {
  return databaseStore.databaseInfo.entry_count
})
const rotue = useRoute()
const query = rotue.query

const data = ref([])
const column = ref([])
const pageLoading = ref(true)

const queryParams = reactive({
  form_id: query.form_id,
  page: 1,
  size: 10,
  total: 0
})

const loading = ref(false)
const onTableChange = (pagination) => {
  queryParams.page = pagination.current
  queryParams.size = pagination.pageSize
  getData()
}
const onSearch = () => {
  queryParams.page = 1
  getData()
}
const getData = () => {
  loading.value = true
  getFormEntryList({ ...queryParams, filter_id: activeKey.value })
    .then((res) => {
      data.value = res.data.list || []
      queryParams.total = +res.data.total || 0
      databaseStore.getDatabaseInfo({ id: query.form_id })
    })
    .finally(() => {
      loading.value = false
    })
}

getFormFieldList({ form_id: query.form_id }).then((res) => {
  column.value = res.data
  pageLoading.value = false
  getData()
})

const panes = ref([])

const activeKey = ref('')
const getSortLists = () => {
  getFormFilterList({ form_id: query.form_id }).then((res) => {
    let lists = res.data || []
    panes.value = lists.filter((item) => item.enabled == 'true')
    if (activeKey.value) {
      // 判断选中的该分类是否被禁用 被禁用的话更换选中分类
      let activeItem = lists.filter((item) => item.id == activeKey.value)
      if ((activeItem.length && activeItem[0].enabled == 'false') || !activeItem.length) {
        activeKey.value = ''
        getData()
      }
    }
  })
}
getSortLists()

const handleManageOkBack = () => {
  getSortLists()
  getData()
}

const handleChangeTab = () => {
  queryParams.page = 1
  getData()
}

const addFilrerModalRef = ref(null)
const onEditTabs = (targetKey, action) => {
  if (action === 'add') {
    // add()
    addFilrerModalRef.value.show()
  }
}

const onDelete = (record) => {
  Modal.confirm({
    title: t('modal_delete_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('modal_delete_content'),
    okText: t('modal_delete_ok'),
    okType: 'danger',
    cancelText: t('modal_delete_cancel'),
    onOk() {
      delFormEntry({ id: record.id }).then((res) => {
        message.success(t('msg_delete_success'))
        handleManageOkBack()
      })
    },
    onCancel() {}
  })
}
const handleClearData = () => {
  Modal.confirm({
    title: t('modal_clear_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('modal_clear_content'),
    okText: t('modal_clear_ok'),
    okType: 'danger',
    cancelText: t('modal_clear_cancel'),
    onOk() {
      emptyFormEntry({ form_id: query.form_id }).then((res) => {
        message.success(t('msg_clear_success'))
        handleManageOkBack()
      })
    },
    onCancel() {}
  })
}
const addDataModalRef = ref(null)
const handleAddData = (data = {}) => {
  addDataModalRef.value.show(JSON.parse(JSON.stringify(data)))
}

const importDataModalRef = ref(null)
const handleImportData = (data = {}) => {
  importDataModalRef.value.show()
}

const exportModalRef = ref(null)
const handleOpenExportModal = () => {
  exportModalRef.value.show()
}

const filterManageModalRef = ref(null)
const handleOpenManageModal = () => {
  filterManageModalRef.value.show()
}
</script>

<style lang="less" scoped>
.loading-box {
  padding-top: 300px;
  text-align: center;
}
.empty-box {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  padding-top: 100px;
  color: #8c8c8c;
  img {
    width: 300px;
  }
}
.field-manage-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  .page-title {
    display: flex;
    align-items: center;
    background-color: #fff;
    color: #000000;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    justify-content: space-between;
    padding-right: 24px;
  }
  .table-wrapper {
    margin-top: 10px;
    padding-right: 24px;
    flex: 1;
    overflow: auto;
  }
  .tabs-wrapper {
    margin-top: 24px;
    padding-right: 24px;
  }
  .setting-btn {
    color: #333;
    cursor: pointer;
  }
  ::v-deep(.ant-tabs-nav-add) {
    min-width: 100px !important;
  }
  ::v-deep(.ant-tabs-extra-content) {
    margin-left: 8px;
  }
  ::v-deep(.ant-table-sticky-scroll) {
    opacity: 0;
  }
}
</style>
