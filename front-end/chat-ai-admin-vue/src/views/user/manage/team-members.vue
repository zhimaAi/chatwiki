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
        <a-table-column title="姓名" data-index="user_name" width="300px">
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
        <a-table-column title="角色" data-index="user_roles" width="190px">
          <template #default="{ record }">{{ record.role }}</template>
        </a-table-column>
        <a-table-column title="最近登录时间" data-index="login_time" width="190px">
          <template #default="{ record }">{{ record.login_time }}</template>
        </a-table-column>
        <a-table-column title="最近登录IP" data-index="IP" width="190px">
          <template #default="{ record }">{{ record.login_ip }}</template>
        </a-table-column>
        <a-table-column title="操作" data-index="action" width="208px">
          <template #default="{ record }">
            <a-flex :gap="16" v-if="record.user_roles != 1">
              <a @click="handleEdit(record)">编辑</a>
              <a @click="handleReSetPassword(record)">重置密码</a>
              <a @click="handleDelete(record)">移除</a>
            </a-flex>
          </template>
        </a-table-column>
      </a-table>
    </div>
    <AddTeamMembers ref="addTeamMembersRef" @ok="getData"></AddTeamMembers>
    <ResetPassword ref="resetPasswordRef" @ok="getData"></ResetPassword>
  </div>
</template>
<script setup>
import { ref, reactive, computed, createVNode } from 'vue'
import { PlusOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import AddTeamMembers from './components/add-team-members.vue'
import ResetPassword from './components/reset-password.vue'
import { getUserList, delUser } from '@/api/manage/index.js'
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/stores/modules/user'
import defaultAvatar from "@/assets/img/role_avatar.png"
const userStore = useUserStore()
const { userInfo } = storeToRefs(userStore)
const user_roles = computed(() => userInfo.value.user_roles)
const requestParams = reactive({
  page: 1,
  size: 10,
  total: 0,
  search: ''
})
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
      item.login_time = item.login_time > 0 ? dayjs(item.login_time * 1000).format('YYYY-MM-DD HH:mm') : '--'
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
</script>
<style lang="less" scoped>
.team-members-pages {
  background: #fff;
  padding: 24px;
  height: 100%;
  .list-box {
    background: #fff;
    margin-top: 8px;
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
</style>
