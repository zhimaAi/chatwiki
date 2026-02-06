<template>
  <cu-scroll>
    <div class="content-box">
      <div class="add-btn-block" v-if="currentPermission == 4">
        <a-button type="primary" @click="openAddModal" :icon="h(PlusOutlined)">{{ t('btn_add_collaborator') }}</a-button>
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
                <span v-if="record.operate_rights == 4">{{ t('permission_manage') }}</span>
                <span v-if="record.operate_rights == 2">{{ t('permission_edit') }}</span>
                <span v-if="record.operate_rights == 1">{{ t('permission_view') }}</span>
              </div>
              <template v-else-if="record.role_type == 1 || record.is_creator == 1">
                <div class="hover-btn-box">{{ t('permission_manage') }}</div>
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
                        <span class="title">{{ t('permission_manage') }}</span>
                      </div>
                      <div class="desc">{{ t('permission_manage_desc') }}</div>
                      <div class="check-box" v-if="record.operate_rights == 4">
                        <CheckOutlined />
                      </div>
                    </div>
                    <div class="role-menu" @click="handleChangeStatus(record, 2)">
                      <div class="menu-header">
                        <EditOutlined />
                        <span class="title">{{ t('permission_edit') }}</span>
                      </div>
                      <div class="desc">{{ t('permission_edit_desc') }}</div>
                      <div class="check-box" v-if="record.operate_rights == 2">
                        <CheckOutlined />
                      </div>
                    </div>
                    <div class="role-menu" @click="handleChangeStatus(record, 1)">
                      <div class="menu-header">
                        <EyeOutlined />
                        <span class="title">{{ t('permission_view') }}</span>
                      </div>
                      <div class="desc">{{ t('permission_view_desc') }}</div>
                      <div class="check-box" v-if="record.operate_rights == 1">
                        <CheckOutlined />
                      </div>
                    </div>
                  </div>
                </template>
                <div class="hover-btn-box">
                  <span v-if="record.operate_rights == 4">{{ t('permission_manage') }}</span>
                  <span v-if="record.operate_rights == 2">{{ t('permission_edit') }}</span>
                  <span v-if="record.operate_rights == 1">{{ t('permission_view') }}</span>
                  <DownOutlined />
                </div>
              </a-popover>
            </template>
            <template v-if="column.key === 'action'">
              <template v-if="currentPermission == 4 && record.user_name != user_name">
                <a @click="handleDel(record)" v-if="record.role_type != 1 && record.is_creator != 1"
                  >{{ t('btn_remove') }}</a
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
import { getRobotPermission } from '@/utils/permission'
import { useUserStore } from '@/stores/modules/user'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.role-permission')

const userStore = useUserStore()

const user_name = computed(() => {
  return userStore.user_name
})

import { useRoute } from 'vue-router'
const query = useRoute().query

const currentPermission = ref(getRobotPermission(query.id))

const tableData = ref([])
const queryParams = reactive({
  object_id: query.id,
  object_type: 1
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
      object_type: 1,
      operate_rights: 4
    }
  ])
}

const handleChangeStatus = (record, operate_rights) => {
  let right_str = operate_rights == 4 ? t('permission_manage') : operate_rights == 2 ? t('permission_edit') : t('permission_view')
  Modal.confirm({
    title: t('title_tip'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_change_permission', { permission: right_str }),
    okText: t('btn_confirm'),
    cancelText: t('btn_cancel'),
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
        message.success(t('msg_modify_success'))
        getData()
      })
    },
    onCancel() {}
  })
}

const handleDel = (record) => {
  Modal.confirm({
    title: t('title_tip'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_confirm_delete'),
    okText: t('btn_confirm'),
    okType: 'danger',
    cancelText: t('btn_cancel'),
    onOk: () => {
      deletePermissionManage({
        identity_id: record.identity_id,
        identity_type: record.identity_type,
        object_id: record.object_id,
        object_type: record.object_type
      }).then((res) => {
        message.success(t('msg_delete_success'))
        getData()
      })
    },
    onCancel() {}
  })
}

const columns = ref([
  {
    title: t('label_name'),
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: t('label_permission'),
    dataIndex: 'operate_rights',
    key: 'operate_rights'
  },
  {
    title: t('label_action'),
    dataIndex: 'action',
    key: 'action',
    width: 80
  }
])
</script>

<style lang="less" scoped>
.content-box {
  // padding: 24px;
  padding-right: 24px;
  padding-top: 24px;
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
