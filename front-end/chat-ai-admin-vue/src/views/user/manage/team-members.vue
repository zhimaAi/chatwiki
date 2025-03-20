<template>
  <div class="team-members-pages">
    <a-flex justify="space-between">
      <a-button type="primary" @click="handleAdd">
        <template #icon>
          <PlusOutlined />
        </template>
        添加团队成员
      </a-button>
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
                <div class="user-name">{{ record.user_name }}</div>
                <div class="nick-name">{{ record.nick_name }}</div>
              </div>
            </div>
          </template>
        </a-table-column>
        <a-table-column title="角色" data-index="role_name" width="100px">
          <template #default="{ record }">{{ record.role_name }}</template>
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
            <div v-if="record.role_type == 1 || record.role_type == 2" class="list-content">
              <div>全部</div>
            </div>
            <div v-else-if="!record.managed_robot_list" class="list-content">
              <div class="list-highlight" @click="addManage('robot', record)">添加</div>
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
              <div class="list-highlight edit" @click="editManage('robot', record)">
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
            <div v-if="record.role_type == 1 || record.role_type == 2" class="list-content">
              <div>全部</div>
            </div>
            <div v-else-if="!record.managed_library_list" class="list-content">
              <div class="list-highlight" @click="addManage('library', record)">添加</div>
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
              <div class="list-highlight edit" @click="editManage('library', record)">
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
            <div v-if="record.role_type == 1 || record.role_type == 2" class="list-content">
              <div>全部</div>
            </div>
            <div v-else-if="!record.managed_form_list" class="list-content">
              <div class="list-highlight" @click="addManage('form', record)">添加</div>
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
              <div class="list-highlight edit" @click="editManage('form', record)">
                <svg-icon name="edit"></svg-icon>
              </div>
            </div>
          </template>
        </a-table-column>
        <a-table-column title="创建时间" data-index="create_time" width="190px">
          <template #default="{ record }">{{ record.create_time }}</template>
        </a-table-column>

        <a-table-column title="操作" data-index="action" width="230px">
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
    <AddTeamMembers ref="addTeamMembersRef" @ok="getData"></AddTeamMembers>
    <ResetPassword ref="resetPasswordRef" @ok="getData"></ResetPassword>

    <!-- 新增弹出，选择数据 -->
    <SeeModelAlert
      ref="seeModelAlertRef"
      :robotList="robotList"
      :libraryList="libraryList"
      :formList="formList"
      @save="onSave"
    />
  </div>
</template>
<script setup>
import { getRobotList } from '@/api/robot/index.js'
import { getLibraryList } from '@/api/library/index.js'
import { getFormList } from '@/api/database/index.js'
import { saveUserManagedDataList } from '@/api/manage/index.js'
import { ref, reactive, createVNode } from 'vue'
import {
  PlusOutlined,
  ExclamationCircleOutlined,
  QuestionCircleOutlined
} from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import AddTeamMembers from './components/add-team-members.vue'
import ResetPassword from './components/reset-password.vue'
import SeeModelAlert from './components/see-model-alert.vue'
import { getUserList, delUser } from '@/api/manage/index.js'
import dayjs from 'dayjs'
import defaultAvatar from '@/assets/img/role_avatar.png'
import { useUserStore } from '@/stores/modules/user'
const userStore = useUserStore()

const user_id = ref(userStore.userInfo.user_id)
const requestParams = reactive({
  page: 1,
  size: 10,
  total: 0,
  search: ''
})
const robotList = ref([])
const libraryList = ref([])
const formList = ref([])
const activeKey = ref('')
const currentUserId = ref('')
// 查看模型
const seeModelAlertRef = ref(null)

const onTableChange = (pagination) => {
  requestParams.page = pagination.current
  requestParams.size = pagination.pageSize
  getData()
}
const onSearch = () => {
  requestParams.page = 1
  getData()
}
const tableData = ref([])

const formatText = (arr) => {
  let newArr = []
  arr.map((item) => {
    newArr.push(item.name)
  })
  return newArr.join(', ')
}
const getData = () => {
  // 获取用户列表
  let parmas = {
    page: requestParams.page,
    size: requestParams.size,
    search: requestParams.search
  }

  getUserList(parmas).then((res) => {
    let lists = res.data.list
    lists.forEach((item) => {
      item.create_time = dayjs(item.create_time * 1000).format('YYYY-MM-DD HH:mm')
      ;(item.managed_robot_list = item.managed_robot_list
        ? JSON.parse(item.managed_robot_list)
        : ''),
        (item.managed_library_list = item.managed_library_list
          ? JSON.parse(item.managed_library_list)
          : ''),
        (item.managed_form_list = item.managed_form_list ? JSON.parse(item.managed_form_list) : '')
    })
    tableData.value = lists
    requestParams.total = +res.data.total
  })
}
onSearch()
const addTeamMembersRef = ref(null)
const handleAdd = () => {
  // 添加用户
  addTeamMembersRef.value.add()
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

const addManage = async (key, record) => {
  if (key === 'robot') {
    await getList()
    // 打开弹窗
    activeKey.value = key
  } else if (key === 'library') {
    await getLibrary()
    // 打开弹窗
    activeKey.value = key
  } else if (key === 'form') {
    await getForm()
    // 打开弹窗
    activeKey.value = key
  }
  currentUserId.value = record.id
  seeModelAlertRef.value.open(activeKey.value)
}

const editManage = async (key, record) => {
  if (key === 'robot') {
    await getList()
    // 打开弹窗
    activeKey.value = key
  } else if (key === 'library') {
    await getLibrary()
    // 打开弹窗
    activeKey.value = key
  } else if (key === 'form') {
    await getForm()
    // 打开弹窗
    activeKey.value = key
  }
  currentUserId.value = record.id
  seeModelAlertRef.value.open(activeKey.value, 'edit', record)
}

const onSave = (ids) => {
  let params = {
    user_id: currentUserId.value,
    t: activeKey.value,
    id_list: ids.join(',')
  }

  saveUserManagedDataList(params)
    .then((res) => {
      message.success('操作成功')
      getData()
    })
    .catch(() => {})
}
</script>
<style lang="less" scoped>
.team-members-pages {
  background: #fff;
  padding: 24px;
  height: 100%;
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
        .user-name {
          color: #595959;
        }
        .nick-name {
          color: #8c8c8c;
        }
      }
    }
  }
}

.disabled {
  color: rgba(0, 0, 0, 0.25);
}
</style>
