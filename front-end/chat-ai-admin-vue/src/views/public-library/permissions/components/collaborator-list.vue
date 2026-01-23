<template>
  <div class="collaborator-list-box">
    <div class="list-tools-box">
      <div class="list-label">
        <svg-icon name="collaborator"></svg-icon>
        <span class="label-text">{{ t('collaborator_permissions') }} ({{ pagination.total }})</span>
      </div>

      <div class="tools-box-right">
        <a-button type="primary" size="small" @click="handleAdd"
          ><PlusOutlined /> {{ t('add_collaborator') }}</a-button
        >
      </div>
    </div>

    <div class="list-box">
      <a-table
        :columns="columns"
        :data-source="dataSource"
        :pagination="pagination"
        @change="onListChange"
      >
        <template #bodyCell="{ column, record }">
          <!-- 姓名列 -->
          <template v-if="column.key === 'name'">
            <div class="name-cell">
              <a-avatar class="avatar" :src="record.avatar || defaultAvatar" />
              <div class="user-info">
                <div class="username">{{ record.user_name }}</div>
                <div class="nickname">{{ record.nick_name }}</div>
              </div>
            </div>
          </template>
          <template v-if="column.key === 'role'">
            <div class="role-cell" v-if="record.showDelete">
              <a-dropdown>
                <a class="dropdown-label" @click.prevent>
                  {{ record.operate_rights_label }}
                  <DownOutlined style="font-size: 12px" />
                </a>
                <template #overlay>
                  <a-menu @click="handleRoleChange($event, record)">
                    <a-menu-item :key="4">{{ t('can_manage') }}</a-menu-item>
                    <a-menu-item :key="2">{{ t('can_edit') }}</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
            <div class="role-cell" v-else>
              <a class="dropdown-label disabled" @click.prevent>
                {{ record.operate_rights_label }}
              </a>
            </div>
          </template>
          <!-- 操作列 -->
          <template v-if="column.key === 'action'">
            <a-space v-if="record.showDelete">
              <a-popconfirm
                :title="t('confirm_remove')"
                :ok-text="t('confirm')"
                :cancel-text="t('cancel')"
                @confirm="handleDel(record)"
              >
                <a href="#">{{ t('remove') }}</a>
              </a-popconfirm>
            </a-space>
            <span v-else>
              <a-button style="padding: 0" type="link" disabled>{{ t('remove') }}</a-button>
            </span>
          </template>
        </template>
      </a-table>
    </div>
    <AddCollaborator :excludedIds="allUserId" ref="addCollaboratorRef" @ok="onAddSuccess" />
  </div>
</template>

<script setup>
import { getLibDocPartnerList, saveLibDocPartner, deleteLibDocPartner } from '@/api/public-library'
import { useUserStore } from '@/stores/modules/user'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, onMounted, inject, computed } from 'vue'
import { PlusOutlined, DownOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import AddCollaborator from './add-collaborator.vue'
import defaultAvatar from '@/assets/img/role_avatar.png'

const { t } = useI18n('views.public-library.permissions.components.collaborator-list')

const userStore = useUserStore()

const libraryState = inject('libraryState', {})
const addCollaboratorRef = ref(null)
const columns = [
  {
    title: t('name'),
    key: 'name',
    width: '40%'
  },
  {
    title: t('role'),
    dataIndex: 'role',
    key: 'role',
    width: '40%'
  },
  {
    title: t('actions'),
    key: 'action',
    width: '20%'
  }
]

const roleMap = {
  4: {
    label: t('can_manage'),
    value: '4'
  },
  2: {
    label: t('can_edit'),
    value: '2'
  }
}

const pagination = reactive({
  current: 1,
  pageSize: 999,
  total: 0
})

const dataSource = ref([])
// 计算所有协作者的用户ID
const allUserId = computed(() => {
  return dataSource.value.map((item) => item.user_id)
})

const handleRoleChange = (event, record) => {
  record.operate_rights = event.key * 1
  record.operate_rights_label = roleMap[record.operate_rights].label || ''
  saveRole(record)
}
const saveRole = (record) => {
  let data = {
    operate_rights: record.operate_rights,
    user_ids: record.user_id,
    type: 2,
    library_key: libraryState.library_key
  }

  saveLibDocPartner(data)
    .then(() => {
      message.success(t('modify_success'))
    })
    .catch(() => {})
}
const handleAdd = () => {
  addCollaboratorRef.value.open()
}

const handleDel = (record) => {
  let data = {
    id: record.id,
    library_key: libraryState.library_key
  }

  deleteLibDocPartner(data)
    .then(() => {
      message.success(t('remove_success'))
      getList()
    })
    .catch(() => {})
}

const onAddSuccess = () => {
  getList()
}

const onListChange = (paginationData) => {
  Object.assign(pagination, paginationData)

  getList()
}

const getList = () => {
  getLibDocPartnerList({
    library_key: libraryState.library_key,
    page: pagination.current,
    size: pagination.pageSize
  }).then((res) => {
    let list = res.data.list || []
    let { user_id } = userStore.userInfo

    pagination.total = res.data.total

    list.forEach((item) => {
      item.showDelete = true

      if (user_id == item.user_id) {
        item.showDelete = false
      }

      if (item.role_type == 1) {
        item.showDelete = false
      }

      item.operate_rights = item.operate_rights * 1
      item.operate_rights_label = t('no_permission')
      if (roleMap[item.operate_rights]) {
        item.operate_rights_label = roleMap[item.operate_rights].label
      }
    })
    dataSource.value = list
  })
}

onMounted(() => {
  getList()
})
</script>

<style lang="less" scoped>
.collaborator-list-box {
  padding: 16px;
  border-radius: 6px;
  background-color: #f2f4f7;
  .list-tools-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
  }
  .list-label {
    display: flex;
    .label-text {
      line-height: 24px;
      padding-left: 4px;
      font-size: 16px;
      font-weight: 600;
      color: #262626;
    }
  }
  .list-box {
    .name-cell {
      display: flex;
      align-items: center;
      .avatar {
        width: 40px;
        height: 40px;
        border-radius: 8px;
      }
      .user-info {
        margin-left: 8px;
        .username {
          font-size: 14px;
          color: #595959;
        }
        .nickname {
          margin-top: 2px;
          font-size: 12px;
          color: #8c8c8c;
        }
      }
    }
    .role-cell {
      .dropdown-label {
        display: inline-block;
        line-height: 22px;
        padding: 1px 8px;
        font-size: 14px;
        font-weight: 400;
        border-radius: 6px;
        color: #595959;
        &:hover {
          background-color: #e4e6eb;
        }
        &.disabled:hover {
          cursor: not-allowed;
          background: none;
        }
      }
    }
  }
}
</style>
