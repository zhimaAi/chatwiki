<template>
  <a-modal
    v-model:open="visible"
    title="添加协作者"
    width="760px"
    :confirmLoading="confirmLoading"
    @ok="handleOk"
    @cancel="handleCancel"
  >
    <div class="form-box">
      <a-form layout="vertical" :model="formState" :rules="rules" ref="formRef">
        <a-form-item name="permission" label="协作权限">
          <a-radio-group v-model:value="formState.permission">
            <a-radio value="4">可管理</a-radio>
            <a-radio value="2">可编辑</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item name="collaborators" label="协作者">
          <div class="collaborator-list-tip">
            如果协作者没有账号，请先到团队管理中添加账号
            <RouterLink to="/user/manage" target="_blank">去添加</RouterLink>
          </div>

          <div class="collaborator-list" ref="listRef">
            <div class="list-content">
              <div class="collaborator-grid">
                <div
                  class="collaborator-item"
                  v-for="user in dataSource"
                  :key="user.id"
                  :class="{ active: formState.collaborators.includes(user.id) }"
                  @click="onSelectUser(user)"
                >
                  <svg-icon
                    name="check-arrow-filled"
                    class="check-icon"
                    v-show="formState.collaborators.includes(user.id)"
                  ></svg-icon>
                  <div class="user-info">
                    <a-avatar class="user-avatar" :src="user.avatar || defaultAvatar" />
                    <div class="user-details">
                      <div class="user-name">{{ user.user_name }}</div>
                      <div class="user-nickname">{{ user.nick_name }}</div>
                    </div>
                  </div>
                </div>
              </div>
              <div v-if="loading" class="loading-more">加载中...</div>
              <div v-if="noMore" class="no-more">没有更多了</div>
            </div>
          </div>
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { getUserList } from '@/api/manage'
import { saveLibDocPartner } from '@/api/public-library'
import { ref, reactive, inject, nextTick, computed } from 'vue'
import { message } from 'ant-design-vue'
import BScroll from '@better-scroll/core'
import MouseWheel from '@better-scroll/mouse-wheel'
import ScrollBar from '@better-scroll/scroll-bar'
import Pullup from '@better-scroll/pull-up'
import defaultAvatar from '@/assets/img/role_avatar.png'

BScroll.use(Pullup)
BScroll.use(ScrollBar)
BScroll.use(MouseWheel)

const emit = defineEmits(['ok'])
const props = defineProps({
  excludedIds: {
    type: Array,
    default: () => []
  }
})

const libraryState = inject('libraryState', {})

const visible = ref(false)
const formRef = ref(null)
const confirmLoading = ref(false)

// 表单数据
const formState = reactive({
  permission: '2',
  collaborators: []
})

// 表单验证规则
const rules = {
  permission: [{ required: true, message: '请选择协作权限' }],
  collaborators: [{ required: true, message: '请选择协作者', type: 'array', min: 1 }]
}

const onSelectUser = (user) => {
  // 处理用户选择逻辑
  let index = formState.collaborators.indexOf(user.id)
  if (index > -1) {
    formState.collaborators.splice(index, 1)
  } else {
    formState.collaborators.push(user.id)
  }
}

// 模拟用户数据
const userList = ref([])
const dataSource = computed(() => {
  // 过滤掉需要排除的用户
  return userList.value.filter((user) => !props.excludedIds.includes(user.id))
})

const pagination = reactive({
  page: 1,
  size: 20,
  total: 0,
  search: ''
})
const listRef = ref(null)
const loading = ref(false)
const noMore = ref(false)
let scroll = ref(null)

const initScroll = () => {
  if (scroll.value) {
    scroll.value.refresh()
    return
  }

  scroll.value = new BScroll(listRef.value, {
    probeType: 2,
    click: true,
    scrollY: true,
    bounce: false,
    pullUpLoad: true,
    scrollbar: {
      fade: false,
      interactive: true
    },
    mouseWheel: {
      speed: 20,
      invert: false,
      easeTime: 300
    }
  })

  scroll.value.on('pullingUp', pullingUpHandler)
}

const pullingUpHandler = async () => {
  if (loading.value || noMore.value) return

  pagination.page = pagination.page + 1
  loading.value = true

  await getList()

  scroll.value.finishPullUp()

  loading.value = false
}

const getList = async () => {
  // 获取用户列表
  let parmas = {
    page: pagination.page,
    size: pagination.size,
    search: pagination.search
  }

  return getUserList(parmas).then((res) => {
    let lists = res.data.list || []

    if (pagination.page == 1) {
      userList.value = userList.value = lists
    } else {
      userList.value = userList.value.concat(lists)
    }

    pagination.total = +res.data.total

    if (userList.value.length >= pagination.total) {
      noMore.value = true
    }

    nextTick(() => {
      if (scroll.value) {
        scroll.value.refresh()
      }
    })
    return res
  })
}

const handleSave = () => {
  let data = {
    operate_rights: formState.permission,
    user_ids: formState.collaborators.join(','),
    type: 1,
    library_key: libraryState.library_key
  }

  confirmLoading.value = true

  saveLibDocPartner(data)
    .then(() => {
      visible.value = false
      confirmLoading.value = false
      emit('ok', data)
      message.success('添加成功')
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

// 确认按钮处理
const handleOk = async () => {
  try {
    await formRef.value?.validate()
    handleSave()
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 取消按钮处理
const handleCancel = () => {
  formRef.value?.resetFields()
  visible.value = false
}

// 打开弹窗方法
const open = () => {
  visible.value = true

  formRef.value?.resetFields()
  pagination.page = 1
  getList()

  setTimeout(() => {
    initScroll()
  }, 200)
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.form-box {
  padding-top: 16px;
}
.loading-more,
.no-more {
  text-align: center;
  padding: 12px 0;
  color: #999;
  font-size: 14px;
}
.collaborator-list-tip {
  line-height: 22px;
  margin-bottom: 8px;
  color: #595959;
}
.collaborator-list {
  position: relative;
  height: 300px;
  overflow: hidden;
}
.collaborator-grid {
}
.collaborator-item {
  display: inline-block;
  width: 48%;
  margin: 0 1% 1% 0;
  cursor: pointer;
  position: relative;
  border-radius: 6px;
  padding: 16px;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 6px;
    border: 1px solid #e8e8e8;
  }
  &:hover::before {
    border: 1px solid #2475fc;
  }
  &.active::before {
    border: 2px solid #2475fc;
  }
  .check-icon {
    position: absolute;
    right: 0;
    bottom: 0;
    font-size: 24px;
    color: #fff;
  }
}

.user-info {
  display: flex;
  align-items: center;
  padding: 8px;
  border-radius: 6px;
  transition: background-color 0.3s;
}

.user-details {
  margin-left: 18px;
}
.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 8px;
}
.user-name {
  height: 22px;
  line-height: 22px;
  font-size: 14px;
  font-weight: 400;
  color: #262626;
}

.user-nickname {
  height: 20px;
  line-height: 20px;
  font-size: 12px;
  color: #8c8c8c;
}
.user-name,
.user-nickname {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
