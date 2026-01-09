<template>
  <div class="main-content-box">
    <a-radio-group v-model:value="object_type" @change="handleChange">
      <a-radio-button value="2">{{ t('knowledge_base') }}</a-radio-button>
      <a-radio-button value="1">{{ t('robot') }}</a-radio-button>
      <a-radio-button value="3">{{ t('database') }}</a-radio-button>
    </a-radio-group>
    <div style="margin-top: 12px">
      <a-button type="primary" @click="addManage">
        {{ addButtonText }}
      </a-button>
    </div>
    <div class="list-content">
      <a-table :columns="columns" :data-source="tableData" :pagination="false">
        <template #headerCell="{ column }">
          <template v-if="column.key === 'name'">
            {{ typeName }}
          </template>
        </template>
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div class="user-block">{{ record.name }}</div>
          </template>

          <template v-if="column.key === 'role'">
            <a-popover
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
                    <div class="desc">{{ t('manage_desc') }}</div>
                    <div class="check-box" v-if="record.operate_rights == 4">
                      <CheckOutlined />
                    </div>
                  </div>
                  <div class="role-menu" @click="handleChangeStatus(record, 2)">
                    <div class="menu-header">
                      <EditOutlined />
                      <span class="title">{{ t('edit') }}</span>
                    </div>
                    <div class="desc">{{ t('edit_desc') }}</div>
                    <div class="check-box" v-if="record.operate_rights == 2">
                      <CheckOutlined />
                    </div>
                  </div>
                  <div class="role-menu" @click="handleChangeStatus(record, 1)">
                    <div class="menu-header">
                      <EyeOutlined />
                      <span class="title">{{ t('view') }}</span>
                    </div>
                    <div class="desc">{{ t('view_desc') }}</div>
                    <div class="check-box" v-if="record.operate_rights == 1">
                      <CheckOutlined />
                    </div>
                  </div>
                </div>
              </template>
              <div class="hover-btn-box">
                {{ t(roleText[record.operate_rights]) }}
                <DownOutlined />
              </div>
            </a-popover>
          </template>
          <template v-if="column.key === 'action'">
            <a @click="handleDel(record)">{{ t('delete') }}</a>
          </template>
        </template>
      </a-table>
    </div>
  </div>
  <!-- 新增弹出，选择数据 -->
  <SeeModelAlert
    ref="seeModelAlertRef"
    :robotList="robotList"
    :libraryList="libraryList"
    :formList="formList"
    @ok="getData"
  />
</template>

<script setup>
import { ref, reactive, watch, createVNode, computed } from 'vue'
import {
  PlusOutlined,
  UserOutlined,
  DownOutlined,
  EditOutlined,
  EyeOutlined,
  CheckOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'
import SeeModelAlert from './see-model-alert.vue'
import { getRobotList } from '@/api/robot/index.js'
import { getLibraryList } from '@/api/library/index.js'
import { getFormList } from '@/api/database/index.js'
import {
  getPermissionManageList,
  savePermissionManage,
  deletePermissionManage
} from '@/api/department/index.js'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.manage.components.knowledge-permissions')

const object_type = ref('2')
const tableData = ref([])

// 动态类型名称
const typeName = computed(() => {
  const typeMap = {
    '2': 'knowledge_base',
    '1': 'robot',
    '3': 'database'
  }
  return t(typeMap[object_type.value] || 'knowledge_base')
})

const addButtonText = computed(() => {
  const addMap = {
    '2': 'add_knowledge_base_permission',
    '1': 'add_robot_permission',
    '3': 'add_database_permission'
  }
  return t(addMap[object_type.value] || 'add_knowledge_base_permission')
})

const roleText = computed(() => {
  const roleMap = {
    '4': 'manage',
    '2': 'edit',
    '1': 'view'
  }
  return roleMap
})

const roleDesc = computed(() => {
  const descMap = {
    '4': 'manage_desc',
    '2': 'edit_desc',
    '1': 'view_desc'
  }
  return descMap
})

const props = defineProps({
  treeParmas: {
    type: Object,
    default: () => {
      return {}
    }
  }
})

watch(
  () => props.treeParmas,
  (value) => {
    getData()
  },
  {
    deep: true
  }
)

const getData = () => {
  getPermissionManageList({
    object_type: object_type.value,
    identity_type: props.treeParmas.is_user ? 1 : 2,
    identity_id: props.treeParmas.id
  }).then((res) => {
    let data = res.data || []
    tableData.value = data
  })
}

getData()

const handleChange = () => {
  getData()
}

const handleChangeStatus = (record, operate_rights) => {
  const roleKey = roleText[operate_rights]
  const roleName = t(roleKey)
  Modal.confirm({
    title: t('prompt'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm_modify_permission', { role: roleName }),
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
        identity_type: 2,
        object_array: JSON.stringify(data)
      }).then((res) => {
        message.success(t('modify_success'))
        getData()
      })
    },
    onCancel() {}
  })
}

const handleDel = (record) => {
  Modal.confirm({
    title: t('prompt'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm_delete_data'),
    okText: t('confirm'),
    okType: 'danger',
    cancelText: t('cancel'),
    onOk: () => {
      deletePermissionManage({
        identity_id: record.identity_id,
        identity_type: 2,
        object_id: record.object_id,
        object_type: record.object_type
      }).then((res) => {
        message.success(t('delete_success'))
        getData()
      })
    },
    onCancel() {}
  })
}

const robotList = ref([])
const libraryList = ref([])
const formList = ref([])

// 获取机器人列表
const getList = async () => {
  await getRobotList()
    .then((res) => {
      robotList.value = res.data
    })
    .catch(() => {})
}

// 获取知识库列表
const getLibrary = async () => {
  await getLibraryList({ type: '' })
    .then((res) => {
      libraryList.value = res.data
    })
    .catch(() => {})
}

// 获取数据库列表
const getForm = async () => {
  await getFormList()
    .then((res) => {
      formList.value = res.data
    })
    .catch(() => {})
}
const seeModelAlertRef = ref(null)
const addManage = async () => {
  if (!props.treeParmas.id) {
    return message.error(t('select_department_first'))
  }
  let key = object_type.value
  if (key == 1) {
    await getList()
  } else if (key == 2) {
    await getLibrary()
  } else if (key == 3) {
    await getForm()
  }

  seeModelAlertRef.value.handleDepartmentOpen(key, props.treeParmas, tableData.value)
}

const columns = computed(() => [
  {
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: t('permission'),
    dataIndex: 'role',
    key: 'role'
  },
  {
    title: t('action'),
    dataIndex: 'action',
    key: 'action',
    width: 80
  }
])
</script>

<style lang="less" scoped>
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
