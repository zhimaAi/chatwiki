<template>
  <div class="team-members-pages">
    <a-flex justify="space-between">
      <a-button type="primary" @click="handleAdd">
        <template #icon>
          <PlusOutlined />
        </template>
        {{ t('add_role_btn') }}
      </a-button>
      <a-input-search
        v-model:value="requestParams.search"
        :placeholder="t('search_placeholder')"
        style="width: 288px"
        @search="onSearch"
      />
    </a-flex>
    <div class="list-box">
      <a-table
        :data-source="tableData"
        :pagination="{
          current: requestParams.page,
          total: requestParams.total,
          pageSize: requestParams.size,
          showQuickJumper: true,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100']
        }"
        @change="onTableChange"
      >
        <a-table-column :title="t('column_role')" data-index="name" width="157px">
          <template #default="{ record }">{{ record.name }}</template>
        </a-table-column>
        <a-table-column :title="t('column_remark')" data-index="mark" width="480px">
          <template #default="{ record }">{{ record.mark }}</template>
        </a-table-column>
        <a-table-column :title="t('column_update_time')" data-index="update_time" width="150px">
          <template #default="{ record }">{{ record.update_time }}</template>
        </a-table-column>
        <a-table-column :title="t('column_operate_name')" data-index="operate_name" width="130px">
          <template #default="{ record }">{{ record.operate_name }}</template>
        </a-table-column>
        <a-table-column :title="t('column_create_name')" data-index="create_name" width="130px">
          <template #default="{ record }">{{ record.create_name }}</template>
        </a-table-column>
        <a-table-column :title="t('column_action')" data-index="action" width="178px">
          <template #default="{ record }">
            <a-flex :gap="16" class="action-box">
              <a-button type="link" @click="handleEdit(record)">{{ t('btn_edit') }}</a-button>
              <a-button type="link" @click="handleDelete(record)" :disabled="record.role_type <= 3 && record.role_type > 0">{{ t('btn_delete') }}</a-button>
            </a-flex>
          </template>
        </a-table-column>
      </a-table>
    </div>
    <AddRole ref="addRoleRef" @ok="getData"></AddRole>
  </div>
</template>
<script setup>
import { ref, reactive, createVNode } from 'vue'
import { PlusOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { getRoleList, delRole } from '@/api/manage/index.js'
import AddRole from './components/add-role.vue'
import dayjs from 'dayjs'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.manage.role-permissions')
const keyword = ref('')
const requestParams = reactive({
  page: 1,
  size: 10,
  total: 0,
  search: ''
})
const tableData = ref([{}])

const onTableChange = (pagination) => {
  requestParams.page = pagination.current
  requestParams.size = pagination.pageSize
  getData()
}
const onSearch = () => {
  requestParams.page = 1
  getData()
}
const getData = () => {
  // 获取用户列表
  let parmas = {
    page: requestParams.page,
    size: requestParams.size,
    search: requestParams.search
  }

  getRoleList(parmas).then((res) => {
    let lists = res.data.list
    lists.forEach((item) => {
      item.update_time = dayjs(item.update_time * 1000).format('YYYY-MM-DD HH:mm')
    })
    tableData.value = lists
    requestParams.total = +res.data.total
  })
}
onSearch()
const addRoleRef = ref(null)
const handleAdd = () => {
  addRoleRef.value.add()
}
const handleEdit = (record) => {
  addRoleRef.value.edit(record)
}
const handleDelete = (record) => {
  // 删除用户
  Modal.confirm({
    title: t('confirm_title'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm_content'),
    okText: t('ok_text'),
    okType: 'danger',
    cancelText: t('cancel_text'),
    onOk: () => {
      delRole({ id: record.id }).then((res) => {
        message.success(t('delete_success'))
        getData()
      })
    },
    onCancel() {}
  })
}
</script>
<style lang="less" scoped>
.team-members-pages {
  background: #fff;
  padding: 24px;
  height: 100%;
  .list-box {
    margin-top: 8px;
  }
}
.action-box .ant-btn-link {
  padding: 0;
}
</style>
