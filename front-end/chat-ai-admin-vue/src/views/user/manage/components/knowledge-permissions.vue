<template>
  <div class="main-content-box">
    <a-radio-group v-model:value="object_type" @change="handleChange">
      <a-radio-button value="2">知识库</a-radio-button>
      <a-radio-button value="1">机器人</a-radio-button>
      <a-radio-button value="3">数据库</a-radio-button>
    </a-radio-group>
    <div style="margin-top: 12px">
      <a-button type="primary" @click="addManage">
        <span v-if="object_type == 2">添加知识库权限</span>
        <span v-if="object_type == 1">添加机器人权限</span>
        <span v-if="object_type == 3">添加数据库权限</span>
      </a-button>
    </div>
    <div class="list-content">
      <a-table :columns="columns" :data-source="tableData" :pagination="false">
        <template #headerCell="{ column }">
          <template v-if="column.key === 'name'">
            <span v-if="object_type == 2">知识库</span>
            <span v-if="object_type == 1">机器人</span>
            <span v-if="object_type == 3">数据库</span>
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
                      <span class="title">管理</span>
                    </div>
                    <div class="desc">可编辑，可删除，可查看</div>
                    <div class="check-box" v-if="record.operate_rights == 4">
                      <CheckOutlined />
                    </div>
                  </div>
                  <div class="role-menu" @click="handleChangeStatus(record, 2)">
                    <div class="menu-header">
                      <EditOutlined />
                      <span class="title">编辑</span>
                    </div>
                    <div class="desc">可编辑，可查看</div>
                    <div class="check-box" v-if="record.operate_rights == 2">
                      <CheckOutlined />
                    </div>
                  </div>
                  <div class="role-menu" @click="handleChangeStatus(record, 1)">
                    <div class="menu-header">
                      <EyeOutlined />
                      <span class="title">查看</span>
                    </div>
                    <div class="desc">仅查看</div>
                    <div class="check-box" v-if="record.operate_rights == 1">
                      <CheckOutlined />
                    </div>
                  </div>
                </div>
              </template>
              <div class="hover-btn-box">
                <span v-if="record.operate_rights == 4">管理</span>
                <span v-if="record.operate_rights == 2">编辑</span>
                <span v-if="record.operate_rights == 1">查看</span>
                <DownOutlined />
              </div>
            </a-popover>
          </template>
          <template v-if="column.key === 'action'">
            <a @click="handleDel(record)">删除</a>
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
import { ref, reactive, watch, createVNode } from 'vue'
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

const object_type = ref('2')
const tableData = ref([])

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
  let right_str = operate_rights == 4 ? '管理' : operate_rights == 2 ? '编辑' : '查看'
  Modal.confirm({
    title: '提示?',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确认修改该数据权限为' + '【' + right_str + '】',
    okText: '确认',
    cancelText: '取消',
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
        message.success('修改成功')
        getData()
      })
    },
    onCancel() {}
  })
}

const handleDel = (record) => {
  Modal.confirm({
    title: '提示?',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确认删除该条数据',
    okText: '确认',
    okType: 'danger',
    cancelText: '取消',
    onOk: () => {
      deletePermissionManage({
        identity_id: record.identity_id,
        identity_type: 2,
        object_id: record.object_id,
        object_type: record.object_type
      }).then((res) => {
        message.success('删除成功')
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
    return message.error('请先选择部门')
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

const columns = ref([
  {
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: '权限',
    dataIndex: 'role',
    key: 'role'
  },
  {
    title: '操作',
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
