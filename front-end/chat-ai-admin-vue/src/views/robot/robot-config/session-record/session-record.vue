<template>
  <div class="team-members-pages">
    <a-flex justify="flex-start" class="screen-box">
      <div class="set-model">
        <div class="label set-model-label">渠道：</div>
        <div class="set-model-body">
          <a-select
            v-model:value="requestParams.app_type"
            placeholder="全部渠道"
            @change="handleChangeModel"
            :style="{ width: '200px' }"
          >
            <a-select-option
              v-for="item in channelItem"
              :key="item.app_type"
              :value="item.app_type"
            >
              <span>{{ item.app_name }}</span>
            </a-select-option>
          </a-select>
        </div>
      </div>

      <div class="set-date">
        <div class="label set-date-label">
          <span>日期：</span>
        </div>
        <div class="set-date-body">
          <DateSelect @dateChange="onDateChange" :datekey="datekey"></DateSelect>
        </div>
      </div>

      <div class="set-name">
        <a-input-search
          v-model:value="requestParams.name"
          placeholder="请输入用户名称搜索"
          style="width: 200px"
          @search="onSearch"
        />
      </div>

      <div class="set-reset">
        <a-button @click="onReset">重置</a-button>
      </div>
    </a-flex>
    <div class="list-box">
      <div class="user-box">
        <User
          :userList="userList"
          @userScrollStart="onUserScrollStart"
          @userScrollEnd="onUserScrollEnd"
          @userClick="onuserClick"
        ></User>
      </div>
      <div class="records-box">
        <MessageList
          ref="messageListRef"
          :isEmpty="isEmpty"
          :messages="messageList"
          :robotInfo="robotInfo"
          :sessionSource="sessionSource"
          @scrollStart="onScrollStart"
          @scrollEnd="onScrollEnd"
        />
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, reactive, onMounted, watch, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import DateSelect from './components/date.vue'
import { storeToRefs } from 'pinia'
import User from './components/user.vue'
import MessageList from './components/message-list.vue'
import { useRobotStore } from '@/stores/modules/robot'
import { useChatStore } from '@/stores/modules/chat'

const isEmpty = ref(false)
const route = useRoute()
const query = route.query
const robotStore = useRobotStore()
const { robotInfo } = robotStore
const chatStore = useChatStore()
const { messageList } = storeToRefs(chatStore)

// 是否允许自动滚动到底部
let isAllowedScrollToBottom = true
const { createMsg, onGetChatMessage, getRecordList, getChannelList } = chatStore
const messageListRef = ref(null)
// master环境只有2个 yun_master环境有5个
const channelItem = ref([])
const userList = ref([])
const sessionSource = ref('')
const datekey = ref('1')

const requestParams = reactive({
  robot_id: query.id, // 机器人ID
  app_type: 'all', // 应用类型:从来源接口获取
  app_id: '', // 应用类型:从来源接口获取-云版才有此参数
  start_time: '', // 开始时间-时间戳
  end_time: '', // 结束时间-时间戳
  name: '', // 客户名称检索
  page: 1, // 页码
  size: 20 // 每页数
})

const has_more = ref(false)

const onDateChange = (date) => {
  requestParams.start_time = date.start_time
  requestParams.end_time = date.end_time
  onSearch()
}

const onReset = () => {
  // 重置
  requestParams.app_type = 'all'
  requestParams.name = ''
  requestParams.page = 1

  // 初始化子组件
  datekey.value = '1' + '-' + Math.random()
}

const onSearch = () => {
  requestParams.page = 1
  getRecordLists()
}

const handleChangeModel = (val) => {
  requestParams.app_type = val
  onSearch()
}

// 监听滚动到底部
const onScrollEnd = () => {
  // console.log('滚动到底部')
}

// 点击用户
const onuserClick = (item) => {
  sessionSource.value = item.app_type
  openNewChat(item)
}

// 用户列表滚动到顶部
const onUserScrollStart = async () => {}

// 用户列表滚动到底部
const onUserScrollEnd = () => {
  if (!has_more.value) {
    return
  }

  requestParams.page++
  getRecordLists()
}

// 滚动到顶部
const onScrollStart = async () => {
  isAllowedScrollToBottom = false

  let msgId = messageList.value[0].uid

  let res = await onGetChatMessage()

  if (res) {
    scrollToMessageById(msgId)
  }
}

const messageListScrollToBottom = () => {
  if (messageListRef.value && isAllowedScrollToBottom) {
    messageListRef.value.scrollToBottom()
  }
}

const scrollToMessageById = (id) => {
  if (messageListRef.value) {
    messageListRef.value.scrollToMessage(id)
  }
}

const resetScroll = () => {
  if (messageListRef.value) {
    messageListRef.value.resetScroll()
  }
}

// 监听 messageList 改变
const onMessageListChange = () => {
  if (!isAllowedScrollToBottom) {
    return
  }
}

let scrollEventTimer = null // 滚动条防抖
const openNewChat = async (item) => {
  isAllowedScrollToBottom = true

  let query = route.query || {}

  let data = {
    robot_key: query.robot_key,
    avatar: item.avatar,
    name: item.name,
    nickname: item.name,
    is_background: 1,
    openid: item.openid,
    dialogue_id: item.dialogue_id
  }

  resetScroll()

  await createMsg(data)
  
  if (scrollEventTimer) {
    clearTimeout(scrollEventTimer)
    scrollEventTimer = null
  }
  // 异步处理避免重复调用
  scrollEventTimer = setTimeout(async () => {
    // 获取对话记录
    let res = await onGetChatMessage()
    if (res) {
      messageListScrollToBottom()
    }
  }, 100)
}

const getChannelLists = async () => {
  const res = await getChannelList()
  channelItem.value = [...[{ app_type: 'all', app_name: '全部渠道' }], ...res.data]
}

const getRecordLists = async () => {
  let params = {
    robot_id: requestParams.robot_id, // 机器人ID
    // app_id: requestParams.app_id, // 应用类型:从来源接口获取-云版才有此参数
    start_time: requestParams.start_time, // 开始时间-时间戳
    end_time: requestParams.end_time, // 结束时间-时间戳
    name: requestParams.name, // 客户名称检索
    page: requestParams.page, // 页码
    size: requestParams.size // 每页数
  }
  if (requestParams.app_type !== 'all') {
    params.app_type = requestParams.app_type // 应用类型:从来源接口获取
  }
  let res = await getRecordList(params)

  let list = res.data.list || []
  if (requestParams.page == 1) {
    userList.value = []
  }
  userList.value = [...userList.value, ...list]
  isEmpty.value = userList.value.length == 0
  has_more.value = res.data.has_more
}

watch(
  () => messageList.value,
  () => {
    onMessageListChange()
  }
)
onMounted(async () => {
  // 获取来源列表
  getChannelLists()
  // 获取会话记录列表
  // await getRobot(query.id)
})

onUnmounted(() => {})
</script>
<style lang="less" scoped>
.team-members-pages {
  background: #fff;
  padding: 0 24px 24px;
  height: 100%;

  .screen-box {
    gap: 24px;
  }

  .list-box {
    margin-top: 16px;
    display: flex;
    flex-wrap: nowrap;
    gap: 24px;
    width: 100%;
    height: calc(100% - 25px);

    .user-box {
      height: 100%;
    }

    .records-box {
      flex: 1;
      flex-shrink: 0;
      border-radius: 6px;
      background: #f0f2f5;
    }
  }
}

.label {
  color: #262626;
  font-family: 'PingFang SC';
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
}

.set-model {
  display: flex;
  align-items: center;

  .set-model-body {
    .set-model-select {
      display: flex;
      padding: 4px 12px;
      align-items: flex-start;
      align-self: stretch;
      border-radius: 2px;
      border: 1px solid var(--06, #d9d9d9);
      background: var(--Neutral-1, #fff);
    }
  }
}

.set-date {
  display: flex;
  align-items: center;
}
</style>
