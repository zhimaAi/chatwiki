<template>
  <cu-scroll>
    <div class="content-box">
      <div class="add-btn-block" v-if="currentPermission == 4">
        <a-button type="primary" @click="openAddModal" :icon="h(PlusOutlined)">{{ t('addCollaborator') }}</a-button>
      </div>
      <div class="list-content">
        <a-table :columns="columns" :data-source="tableData" :pagination="false">
          <!-- @change="onTableChange" -->

          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'name'">
              <div class="user-block">
                <a-avatar shape="square" :src="record.avatar" :size="40">
                  <template #icon><UserOutlined v-if="!record.avatar" /></template>
                </a-avatar>
                <div class="name-box">
                  <div class="name">{{ record.user_name }}</div>
                  <div class="remark-text">{{ record.name }}</div>
                </div>
              </div>
            </template>

            <template v-if="column.key === 'operate_rights'">
              <div class="hover-btn-box" v-if="currentPermission == 2 || record.user_name == user_name">
                <span v-if="record.operate_rights == 4">{{ t('manage') }}</span>
                <span v-if="record.operate_rights == 2">{{ t('edit') }}</span>
                <span v-if="record.operate_rights == 1">{{ t('view') }}</span>
              </div>
              <template v-else-if="record.role_type == 1 || record.is_creator == 1">
                <div class="hover-btn-box">{{ t('manage') }}</div>
              </template>
              <a-popover
                v-else
                placement="bottomRight"
                :title="null"
                :overlayInnerStyle="{ padding: '2px' }"
              >
                <template #content>
                  <div class="role-menu-box">
                    <div class="role-menu" @click="handleChangeStatus(record, 4)">
                      <div class="menu-header">
                        <UserOutlined />
                        <span class="title">{{ t('manage') }}</span>
                      </div>
                      <div class="desc">{{ t('manageDesc') }}</div>
                      <div class="check-box" v-if="record.operate_rights == 4">
                        <CheckOutlined />
                      </div>
                    </div>
                    <div class="role-menu" @click="handleChangeStatus(record, 2)">
                      <div class="menu-header">
                        <EditOutlined />
                        <span class="title">{{ t('edit') }}</span>
                      </div>
                      <div class="desc">{{ t('editDesc') }}</div>
                      <div class="check-box" v-if="record.operate_rights == 2">
                        <CheckOutlined />
                      </div>
                    </div>
                    <div class="role-menu" @click="handleChangeStatus(record, 1)">
                      <div class="menu-header">
                        <EyeOutlined />
                        <span class="title">{{ t('view') }}</span>
                      </div>
                      <div class="desc">{{ t('viewDesc') }}</div>
                      <div class="check-box" v-if="record.operate_rights == 1">
                        <CheckOutlined />
                      </div>
                    </div>
                  </div>
                </template>
                <div class="hover-btn-box">
                  <span v-if="record.operate_rights == 4">{{ t('manage') }}</span>
                  <span v-if="record.operate_rights == 2">{{ t('edit') }}</span>
                  <span v-if="record.operate_rights == 1">{{ t('view') }}</span>
                  <DownOutlined />
                </div>
              </a-popover>
            </template>
            <template v-if="column.key === 'action'">
              <template v-if="currentPermission == 4 && record.user_name != user_name">
                <a @click="handleDel(record)" v-if="record.role_type != 1 && record.is_creator != 1"
                  >{{ t('remove') }}</a
                >
              </template>
            </template>
          </template>
        </a-table>
      </div>
    </div>
    <AddCollaborator @ok="getData" ref="addCollaboratorRef" />
  </cu-scroll>
</template>
<script setup>
import { ref, reactive, h, createVNode, computed } from 'vue'
import {
  PlusOutlined,
  UserOutlined,
  DownOutlined,
  EditOutlined,
  EyeOutlined,
  CheckOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'

import {
  savePermissionManage,
  getPartnerManageList,
  deletePermissionManage
} from '@/api/department/index.js'
import { message, Modal } from 'ant-design-vue'
import AddCollaborator from '@/components/add-collaborator/add-collaborator.vue'
import { getLibraryPermission } from '@/utils/permission'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
const query = useRoute().query
const { t } = useI18n('views.library.library-details.knowledge-config.role-permission')

import { useUserStore } from '@/stores/modules/user'
const userStore = useUserStore()

const user_name = computed(() => {
  return userStore.user_name
})

const currentPermission = ref(getLibraryPermission(query.id))

const tableData = ref([])
const queryParams = reactive({
  object_id: query.id,
  object_type: 2
})

// "operate_rights": "4",  // 权限 4:可管理  2:可编辑 1:查看
const getData = () => {
  getPartnerManageList({ ...queryParams }).then((res) => {
    let list = res.data || []
    tableData.value = list
  })
}

getData()

const addCollaboratorRef = ref(null)

const openAddModal = () => {
  addCollaboratorRef.value.show(tableData.value, [
    {
      object_id: +query.id,
      object_type: 2,
      operate_rights: 4
    }
  ])
}

const handleChangeStatus = (record, operate_rights) => {
  let right_str = operate_rights == 4 ? t('manage') : operate_rights == 2 ? t('edit') : t('view')
  Modal.confirm({
    title: t('confirmTitle'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirmModifyPermission', { permission: right_str }),
    okText: t('confirm'),
    cancelText: t('cancel'),
    onOk: () => {
      let data = [
        {
          object_id: +record.object_id,
          object_type: +record.object_type,
          operate_rights: +operate_rights
        }
      ]
      savePermissionManage({
        identity_ids: record.identity_id,
        identity_type: record.identity_type,
        object_array: JSON.stringify(data)
      }).then((res) => {
        message.success(t('modifySuccess'))
        getData()
      })
    },
    onCancel() {}
  })
}

const handleDel = (record) => {
  Modal.confirm({
    title: t('confirmTitle'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirmDelete'),
    okText: t('confirm'),
    okType: 'danger',
    cancelText: t('cancel'),
    onOk: () => {
      deletePermissionManage({
        identity_id: record.identity_id,
        identity_type: record.identity_type,
        object_id: record.object_id,
        object_type: record.object_type
      }).then((res) => {
        message.success(t('deleteSuccess'))
        getData()
      })
    },
    onCancel() {}
  })
}

const columns = ref([
  {
    title: t('name'),
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: t('permission'),
    dataIndex: 'operate_rights',
    key: 'operate_rights'
  },
  {
    title: t('operation'),
    dataIndex: 'action',
    key: 'action',
    width: 80
  }
])
</script>

<style lang="less" scoped>
.content-box {
  padding: 24px;
}
.list-content {
  margin-top: 16px;
}

.user-block {
  display: flex;
  align-items: center;
  gap: 8px;
  .name-box {
    flex: 1;
    color: #595959;
    line-height: 22px;
    font-size: 14px;
  }
  .remark-text {
    margin-top: 2px;
    font-size: 12px;
    color: #8c8c8c;
    line-height: 20px;
  }
}

.hover-btn-box {
  display: flex;
  align-items: center;
  gap: 4px;
  width: fit-content;
  height: 24px;
  font-size: 14px;
  padding: 0 8px;
  border-radius: 6px;
  color: #595959;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
  &:hover {
    background: #e4e6eb;
  }
}

.role-menu-box {
  display: flex;
  flex-direction: column;
  gap: 2px;
  .role-menu {
    padding: 5px 16px;
    border-radius: 6px;
    color: #262626;
    cursor: pointer;
    position: relative;
    &.active {
      color: #2475fc;
    }
    &:hover {
      background: #f2f4f7;
    }
    .check-box {
      position: absolute;
      top: 16px;
      right: 16px;
      color: #2475fc;
    }
    .menu-header {
      display: flex;
      align-items: center;
      gap: 8px;
    }
    .desc {
      color: #8c8c8c;
      margin-top: 4px;
      line-height: 22px;
      padding-left: 24px;
    }
  }
}
</style>
