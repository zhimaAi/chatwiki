<template>
  <div class="user-model-page">
    <div class="page-title">
      {{ t('title') }}
      <div class="add-block">
        <a-button type="primary" @click="openAddModal()" :icon="createVNode(PlusOutlined)"
          >{{ t('add_btn') }}</a-button
        >
      </div>
    </div>
    <div class="list-wrapper">
      <cu-scroll>
        <div class="content-box">
          <a-table
            :columns="columns"
            :data-source="tableData"
            :pagination="{
              current: pager.page,
              total: pager.total,
              pageSize: pager.size,
              showQuickJumper: true,
              showSizeChanger: true,
              pageSizeOptions: ['10', '20', '50', '100']
            }"
            @change="onTableChange"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'words_desc'">
                <a-tooltip v-if="record.words_desc.length > 20" :overlayStyle="{'max-width': '500px'}">
                  <template #title>{{ record.words_desc }}</template>
                  <div class="text-over-box">{{ record.words_desc }}</div>
                </a-tooltip>
                <div v-else class="text-over-box">{{ record.words_desc }}</div>
              </template>
              <template v-if="column.key === 'robot_name_desc'">
                <div v-if="record.trigger_type == 0">{{ t('all_robots') }}</div>
                <template v-else>
                  <a-tooltip v-if="record.robot_name_desc.length > 20" :overlayStyle="{'max-width': '500px'}">
                    <template #title>{{ record.robot_name_desc }}</template>
                    <div class="text-over-box">{{ record.robot_name_desc }}</div>
                  </a-tooltip>
                  <div class="text-over-box" v-else>{{ record.robot_name_desc }}</div>
                </template>
              </template>
              <template v-if="column.key === 'status'">
                <a-switch
                  @change="handleChangeStatus(record)"
                  :checked="record.status == 1"
                />
              </template>
              <template v-if="column.key === 'action'">
                <a-flex :gap="12" style="white-space: nowrap">
                  <a @click="openAddModal(record)">{{ t('edit') }}</a>
                  <a @click="handleDel(record)">{{ t('delete') }}</a>
                </a-flex>
              </template>
            </template>
          </a-table>
        </div>
      </cu-scroll>
      <AddaddSensitiveModal @ok="getList" :robotList="robotList" ref="addaddSensitiveModalRef" />
    </div>
  </div>
</template>

<script setup>
import {
  getSensitiveWordsList,
  switchSensitiveWords,
  getRobotList,
  deleteSensitiveWords
} from '@/api/robot'
import { PlusOutlined } from '@ant-design/icons-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { useRoute } from 'vue-router'
import { reactive, ref, createVNode } from 'vue'
import AddaddSensitiveModal from './components/add-sensitive-modal.vue'
import { Modal, message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.sensitive-words.index')
const query = useRoute().query

const tableData = ref([])

const pager = reactive({
  page: 1,
  size: 20,
  total: 0
})

const onTableChange = (pagination) => {
  pager.page = pagination.current
  pager.size = pagination.pageSize
  getList()
}
const getList = () => {
  getSensitiveWordsList({
    ...pager
  }).then((res) => {
    let lists = res.data.list || []
    tableData.value = lists.map((item) => {
      let robot_data = JSON.parse(item.robot_data)
      let robot_name_desc = robot_data.map((item) => item.robot_name).join(';')
      let robot_ids = robot_data.map((item) => item.id)
      let words_desc = item.words.split('\n').join('；')
      return {
        ...item,
        robot_data,
        robot_name_desc,
        robot_ids,
        words_desc
      }
    })
    pager.total = res.data.total
  })
}

getList()

const handleChangeStatus = (record) => {
  switchSensitiveWords({
    id: record.id
  }).then((res) => {
    getList()
  })
}

const handleDel = (record) => {
  Modal.confirm({
    title: t('delete_confirm'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('delete_confirm_content', { words: record.words_desc }),
    okText: t('delete'),
    okType: 'danger',
    cancelText: t('cancel'),
    onOk() {
      deleteSensitiveWords({ id: record.id }).then((res) => {
        message.success(t('delete_success'))
        getList()
      })
    },
    onCancel() {}
  })
}

const addaddSensitiveModalRef = ref(null)
const openAddModal = (data = {}) => {
  addaddSensitiveModalRef.value.show(JSON.parse(JSON.stringify(data)))
}

const robotList = ref([])

getRobotList().then((res) => {
  robotList.value = res.data || []
})

const columns = [
  {
    title: t('sensitive_words'),
    dataIndex: 'words_desc',
    key: 'words_desc',
    ellipsis: true,
    width: 300
  },
  {
    title: t('effective_robots'),
    dataIndex: 'robot_name_desc',
    key: 'robot_name_desc',
    ellipsis: true,
    width: 300
  },
  {
    title: t('enable_status'),
    key: 'status',
    dataIndex: 'status',
    width: 150
  },
  {
    title: t('action'),
    key: 'action',
    width: 120
  }
]
</script>

<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  border-right: 1px solid #fff;
  background-color: #f2f4f7;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  .page-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 24px;
    padding: 24px 24px 16px;
    background-color: #fff;
    color: #000000;
    font-size: 16px;
    font-weight: 600;
  }
  .list-wrapper {
    flex: 1;
    background: #fff;
    overflow: hidden;
  }
  .content-box {
    padding: 0 24px;
    .add-block {
      display: flex;
      justify-content: flex-end;
      margin-bottom: 12px;
    }
  }
}
.text-over-box {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 100%;
}

.ant-tooltip{
  max-width: 500px;
}
</style>
