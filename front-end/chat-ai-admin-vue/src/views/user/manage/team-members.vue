<template>
  <div class="flex-content-box">
    <div class="left-department-box">
      <cu-scroll style="padding-right: 24px">
        <DepartmentTree ref="departmentTreeRef" @select="handleSelect" @refresh="getData" />
      </cu-scroll>
    </div>
    <div class="team-members-pages">
      <cu-scroll style="padding-right: 24px">
        <div class="main-tabs-block">
          <a-segmented v-model:value="tabs" :options="tabOption" @change="handleTabsChange">
            <template #label="{ payload }">
              {{ payload.title }}
            </template>
          </a-segmented>
        </div>
        <template v-if="tabs == 1">
          <a-flex justify="space-between">
            <a-flex :gap="8">
              <a-button type="primary" @click="handleAdd">
                <template #icon>
                  <PlusOutlined />
                </template>
                添加团队成员
              </a-button>
              <a-dropdown>
                <a-button>批量修改权限</a-button>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="handleBatchAdd(1)"> 修改机器人 </a-menu-item>
                    <a-menu-item @click="handleBatchAdd(2)"> 修改知识库 </a-menu-item>
                    <a-menu-item @click="handleBatchAdd(3)"> 修改数据库 </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>

              <a-button @click="handleBatchDepart">批量修改部门</a-button>
            </a-flex>

            <a-input-search
              v-model:value="requestParams.search"
              placeholder="输入成员账号或昵称搜索"
              style="width: 288px"
              @search="onSearch"
            />
          </a-flex>
          <div class="list-box">
            <a-table
              :data-source="tableData"
              :scroll="{ x: 1100 }"
              :loading="loading"
              row-key="id"
              :row-selection="{
                selectedRowKeys: state.selectedRowKeys,
                onChange: onSelectChange,
                getCheckboxProps: getCheckboxProps
              }"
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
              <a-table-column title="成员名称" data-index="user_name" width="190px">
                <template #default="{ record }">
                  <div class="user-box">
                    <img :src="record.avatar || defaultAvatar" alt="" />
                    <div class="name-info">
                      <div class="user-name">
                        <a-tooltip
                          :title="record.user_name"
                          v-if="record.user_name && record.user_name.length > 7"
                        >
                          {{ record.user_name }}
                        </a-tooltip>
                        <template v-else>{{ record.user_name }}</template>
                      </div>
                      <div class="nick-name">
                        <a-tooltip
                          :title="record.nick_name"
                          v-if="record.nick_name && record.nick_name.length > 7"
                        >
                          {{ record.nick_name }}
                        </a-tooltip>
                        <template v-else>{{ record.nick_name }}</template>
                      </div>
                    </div>
                  </div>
                </template>
              </a-table-column>
              <a-table-column title="角色" data-index="role_name" width="100px">
                <template #default="{ record }">{{ record.role_name }}</template>
              </a-table-column>

              <a-table-column data-index="departments" width="200px">
                <template #title> 部门 </template>
                <template #default="{ record }">
                  <div
                    v-if="record.departments && record.departments.length == 0"
                    class="list-content"
                  >
                    --
                  </div>
                  <div v-else class="list-preview list-content">
                    <a-tooltip>
                      <template #title
                        ><div style="display: inline">
                          {{ formatText(record.departments) }}
                        </div></template
                      >
                      <div class="list-item-box">
                        <div class="list-item">{{ formatText(record.departments) }}</div>
                      </div>
                    </a-tooltip>
                  </div>
                </template>
              </a-table-column>

              <a-table-column data-index="managed_robot_list" width="200px">
                <template #title>
                  管理的机器人
                  <a-tooltip>
                    <template #title
                      >所有者和管理员可以创建和管理全部机器人，成员只能管理分配的机器人。</template
                    >
                    <QuestionCircleOutlined />
                  </a-tooltip>
                </template>
                <template #default="{ record }">
                  <div v-if="record.role_type == 1" class="list-content">
                    <div>全部</div>
                  </div>
                  <div v-else-if="!record.managed_robot_list" class="list-content">
                    <div class="list-highlight" @click="addManage(1, record)">添加</div>
                  </div>
                  <div v-else class="list-preview list-content">
                    <a-tooltip>
                      <template #title
                        ><div style="display: inline">
                          {{ formatText(record.managed_robot_list) }}
                        </div></template
                      >
                      <div class="list-item-box">
                        <div class="list-item">{{ formatText(record.managed_robot_list) }}</div>
                      </div>
                    </a-tooltip>
                    <div class="list-highlight edit" @click="editManage(1, record)">
                      <svg-icon name="edit"></svg-icon>
                    </div>
                  </div>
                </template>
              </a-table-column>
              <a-table-column data-index="managed_library_list" width="200px">
                <template #title>
                  管理的知识库
                  <a-tooltip>
                    <template #title
                      >所有者和管理员可以创建和管理全部知识库，成员只能管理分配的知识库。</template
                    >
                    <QuestionCircleOutlined />
                  </a-tooltip>
                </template>
                <template #default="{ record }">
                  <div v-if="record.role_type == 1" class="list-content">
                    <div>全部</div>
                  </div>
                  <div v-else-if="!record.managed_library_list" class="list-content">
                    <div class="list-highlight" @click="addManage(2, record)">添加</div>
                  </div>
                  <div v-else class="list-preview list-content">
                    <a-tooltip>
                      <template #title
                        ><div style="display: inline">
                          {{ formatText(record.managed_library_list) }}
                        </div></template
                      >
                      <div class="list-item-box">
                        <div class="list-item">{{ formatText(record.managed_library_list) }}</div>
                      </div>
                    </a-tooltip>
                    <div class="list-highlight edit" @click="editManage(2, record)">
                      <svg-icon name="edit"></svg-icon>
                    </div>
                  </div>
                </template>
              </a-table-column>
              <a-table-column data-index="managed_form_list" width="200px">
                <template #title>
                  管理的数据库
                  <a-tooltip>
                    <template #title
                      >所有者和管理员可以创建和管理全部数据库，成员只能管理分配的数据库。</template
                    >
                    <QuestionCircleOutlined />
                  </a-tooltip>
                </template>
                <template #default="{ record }">
                  <div v-if="record.role_type == 1" class="list-content">
                    <div>全部</div>
                  </div>
                  <div v-else-if="!record.managed_form_list" class="list-content">
                    <div class="list-highlight" @click="addManage(3, record)">添加</div>
                  </div>
                  <div v-else class="list-preview list-content">
                    <a-tooltip>
                      <template #title
                        ><div style="display: inline">
                          {{ formatText(record.managed_form_list) }}
                        </div></template
                      >
                      <div class="list-item-box">
                        <div class="list-item">{{ formatText(record.managed_form_list) }}</div>
                      </div>
                    </a-tooltip>
                    <div class="list-highlight edit" @click="editManage(3, record)">
                      <svg-icon name="edit"></svg-icon>
                    </div>
                  </div>
                </template>
              </a-table-column>

              <a-table-column title="状态" data-index="expire_status" width="120px">
                <template #default="{ record }">
                  <div class="status-block green" v-if="record.expire_status == 1">
                    <CheckCircleFilled />生效中
                  </div>
                  <div class="status-block yellow" v-if="record.expire_status == 0">
                    <ExclamationCircleFilled />已过期
                  </div>
                </template>
              </a-table-column>
              <a-table-column title="过期时间" data-index="expire_time_desc" width="130px">
                <template #default="{ record }">{{ record.expire_time_desc }}</template>
              </a-table-column>

              <a-table-column title="登录" data-index="login_switch" width="90px" fixed="right">
                <template #default="{ record }">
                  <a-switch
                    @change="handleSwichLogin(record)"
                    :checked="record.login_switch == '1'"
                    :disabled="record.role_type == '1'"
                    checked-children="开"
                    un-checked-children="关"
                  />
                </template>
              </a-table-column>

              <a-table-column title="操作" data-index="action" width="176px" fixed="right">
                <template #default="{ record }">
                  <a-flex :gap="16" v-if="record.role_type == '1'">
                    <span class="disabled">编辑</span>
                    <span class="disabled">重置密码</span>
                    <span class="disabled">删除</span>
                  </a-flex>
                  <a-flex :gap="16" v-else>
                    <a @click="handleEdit(record)">编辑</a>
                    <a @click="handleReSetPassword(record)">重置密码</a>
                    <a @click="handleDelete(record)">删除</a>
                  </a-flex>
                </template>
              </a-table-column>
            </a-table>
          </div>
        </template>
        <template v-else>
          <KnowledgePermissions :treeParmas="treeParmas" />
        </template>
      </cu-scroll>

      <AddTeamMembers ref="addTeamMembersRef" @ok="refreshTreeAndData"></AddTeamMembers>
      <ResetPassword ref="resetPasswordRef" @ok="getData"></ResetPassword>

      <!-- 新增弹出，选择数据 -->
      <SeeModelAlert
        ref="seeModelAlertRef"
        :robotList="robotList"
        :libraryList="libraryList"
        :formList="formList"
        @ok="getData"
      />
      <BatchEditDepartment ref="batchEditDepartmentRef" @ok="refreshTreeAndData" />
    </div>
  </div>
</template>
<script setup>
import { getRobotList } from '@/api/robot/index.js'
import { getLibraryList } from '@/api/library/index.js'
import { getFormList } from '@/api/database/index.js'
import { ref, reactive, createVNode } from 'vue'
import {
  PlusOutlined,
  ExclamationCircleOutlined,
  QuestionCircleOutlined,
  CheckCircleFilled,
  ExclamationCircleFilled
} from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import AddTeamMembers from './components/add-team-members.vue'
import ResetPassword from './components/reset-password.vue'
import SeeModelAlert from './components/see-model-alert.vue'
import { getUserList, delUser, loginSwitch } from '@/api/manage/index.js'
import dayjs from 'dayjs'
import defaultAvatar from '@/assets/img/role_avatar.png'
import DepartmentTree from './components/department-tree.vue'
import KnowledgePermissions from './components/knowledge-permissions.vue'
import BatchEditDepartment from './components/batch-edit-department.vue'

const tabs = ref(1)

const tabOption = [
  {
    value: 1,
    payload: {
      title: '团队成员'
    }
  },
  {
    value: 2,
    payload: {
      title: '知识库权限'
    }
  }
]

const requestParams = reactive({
  page: 1,
  size: 10,
  total: 0,
  search: ''
})
const robotList = ref([])
const libraryList = ref([])
const formList = ref([])
// 查看模型
const seeModelAlertRef = ref(null)

const onTableChange = (pagination) => {
  requestParams.page = pagination.current
  requestParams.size = pagination.pageSize
  getData()
}

// let treeParmas = {}
const treeParmas = reactive({
  id: ''
})
const handleSelect = (data) => {
  Object.assign(treeParmas, { ...data })
  if (tabs.value == 1) {
    onSearch()
  }
}

const onSearch = () => {
  requestParams.page = 1
  getData()
}
const tableData = ref([])

const formatText = (arr) => {
  let newArr = []
  arr.map((item) => {
    newArr.push(item.name || item.department_name)
  })
  return newArr.join(', ')
}
const loading = ref(false)
const getData = () => {
  // 获取用户列表
  let parmas = {
    page: requestParams.page,
    size: requestParams.size,
    search: requestParams.search
  }
  if (treeParmas.id && treeParmas.is_default != 1) {
    parmas.department_id = treeParmas.id
  }
  loading.value = true
  getUserList(parmas)
    .then((res) => {
      let lists = res.data.list
      lists.forEach((item) => {
        item.create_time = dayjs(item.create_time * 1000).format('YYYY-MM-DD HH:mm')

        item.expire_time_desc =
          item.expire_time == '0'
            ? '永久有效'
            : dayjs(item.expire_time * 1000).format('YYYY-MM-DD HH:mm')

        item.expire_status = 0
        if (item.expire_time == '0' || dayjs().unix() < item.expire_time) {
          item.expire_status = 1
        }

        item.managed_robot_list = item.managed_robot_list ? JSON.parse(item.managed_robot_list) : []
        item.managed_library_list = item.managed_library_list
          ? JSON.parse(item.managed_library_list)
          : []
        item.managed_form_list = item.managed_form_list ? JSON.parse(item.managed_form_list) : []

        item.departments = item.departments ? JSON.parse(item.departments) : []
      })
      tableData.value = lists
      state.selectedRowKeys = []
      requestParams.total = +res.data.total
    })
    .finally(() => {
      loading.value = false
    })
}
// onSearch()

const handleTabsChange = () => {
  if (tabs.value == 1) {
    onSearch()
  }
}

const departmentTreeRef = ref(null)
const refreshTreeAndData = () => {
  getData()
  departmentTreeRef.value.getLists()
}

const state = reactive({
  selectedRowKeys: []
})

const onSelectChange = (selectedRowKeys) => {
  state.selectedRowKeys = selectedRowKeys
}

const getCheckboxProps = (record) => {
  return {
    disabled: record.role_type == 1
  }
}

const handleSwichLogin = (record) => {
  loginSwitch({
    user_id: record.id
  }).then((res) => {
    message.success(record.login_switch == 1 ? '关闭成功' : '开启成功')
    getData()
  })
}

const addTeamMembersRef = ref(null)
const handleAdd = () => {
  // 添加用户
  addTeamMembersRef.value.add(treeParmas.id)
}
const handleEdit = (record) => {
  // 编辑用户
  addTeamMembersRef.value.edit(record)
}
const resetPasswordRef = ref(null)
const handleReSetPassword = (record) => {
  // 重置密码
  resetPasswordRef.value.open(record)
}
const handleDelete = (record) => {
  // 删除用户
  Modal.confirm({
    title: '提示?',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确认删除该用户',
    okText: '确认',
    okType: 'danger',
    cancelText: '取消',
    onOk: () => {
      delUser({ id: record.id }).then((res) => {
        message.success('删除成功')
        getData()
      })
    },
    onCancel() {}
  })
}

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

const handleBatchAdd = async (key) => {
  let selectList = []
  state.selectedRowKeys.forEach((item) => {
    let findItem = tableData.value.filter((it) => it.id == item)[0]
    if (findItem.role_type == 1) {
      // 管理员
    } else {
      selectList.push(item)
    }
  })
  state.selectedRowKeys = selectList
  if (state.selectedRowKeys.length == 0) {
    return message.error('请选择要修改的成员')
  }
  if (key == 1) {
    await getList()
  } else if (key == 2) {
    await getLibrary()
  } else if (key == 3) {
    await getForm()
  }
  seeModelAlertRef.value.handleBatchOpen(key, state.selectedRowKeys.join(','))
}

const batchEditDepartmentRef = ref(null)
const handleBatchDepart = () => {
  if (state.selectedRowKeys.length == 0) {
    return message.error('请选择要修改的成员')
  }
  batchEditDepartmentRef.value.show(state.selectedRowKeys.join(','))
}

const addManage = async (key, record) => {
  if (key == 1) {
    await getList()
  } else if (key == 2) {
    await getLibrary()
  } else if (key == 3) {
    await getForm()
  }
  seeModelAlertRef.value.open(key, 'add', record)
}

const editManage = async (key, record) => {
  if (key == 1) {
    await getList()
  } else if (key == 2) {
    await getLibrary()
  } else if (key == 3) {
    await getForm()
  }
  seeModelAlertRef.value.open(key, 'edit', record)
}
</script>
<style lang="less" scoped>
.flex-content-box {
  display: flex;
  height: 100%;
  overflow: hidden;
}
.left-department-box {
  width: 280px;
  height: 100%;
  overflow: hidden;
  padding: 24px;
  padding-right: 0;
  border-right: 1px solid var(--07, #f0f0f0);
}
.main-tabs-block {
  margin-bottom: 16px;
  ::v-deep(.ant-segmented .ant-segmented-item-selected) {
    color: #2475fc;
  }
  ::v-deep(.ant-segmented .ant-segmented-item-label) {
    padding: 0 16px;
  }
}
.team-members-pages {
  background: #fff;
  padding: 24px;
  padding-right: 0;
  padding-bottom: 0;
  flex: 1;
  height: 100%;
  overflow: hidden;
  .list-box {
    background: #fff;
    margin-top: 8px;

    .list-item-box {
      max-width: 200px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .list-preview {
      .list-item {
        display: inline;
      }
    }

    .list-content {
      padding-right: 15px;
      position: relative;

      .list-highlight {
        color: #1677ff;
        cursor: pointer;

        &:hover {
          opacity: 0.7;
        }
      }

      .edit {
        display: none;
        position: absolute;
        right: 0;
        top: 50%;
        margin-top: -10px;
      }

      &:hover .edit {
        display: block;
      }
    }

    .user-box {
      display: flex;
      img {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        margin-right: 8px;
      }
      .name-info {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        font-size: 14px;
        line-height: 22px;
        font-weight: 400;
        overflow: hidden;
        .user-name {
          color: #595959;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
        }
        .nick-name {
          color: #8c8c8c;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
        }
      }
    }
  }
}
.status-block {
  display: flex;
  align-items: center;
  gap: 2px;
  border-radius: 6px;
  padding: 0 6px;
  font-weight: 500;
  width: fit-content;
  &.green {
    background: #cafce4;
    color: #21a665;
  }
  &.yellow {
    background: #fae4dc;
    color: #ed744a;
  }
}
.disabled {
  color: rgba(0, 0, 0, 0.25);
}
</style>
