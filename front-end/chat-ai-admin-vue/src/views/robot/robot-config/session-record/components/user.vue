<template>
  <div class="user-content" ref="scrollUserBoxRef" @scroll="onScroll">
    <div class="empty-box" v-if="userLists.length === 0">
      <img src="@/assets/img/library/detail/empty.png" alt="" />
      <div class="title">暂无结果，请重试</div>
    </div>
    <div
      v-else
      class="user-item"
      :class="{ active: item.openid == currentItem.openid }"
      @click="search(item)"
      v-for="item in userLists"
      :key="item.session_id"
    >
      <div class="user-item-left">
        <img class="user-img" :src="item.avatar" alt="" />
      </div>
      <div class="user-item-right">
        <div class="user-name-box">
          <div class="user-name">{{ item.name }}</div>
          <div v-ftime="'MM-DD HH:mm'" class="user-timer">{{ item.last_chat_time }}</div>
        </div>
        <div class="user-info">{{ item.last_chat_message }}</div>
        <div class="user-source-box">
          <div class="user-source-title">来自：</div>
          <div class="user-source-content">{{ formatSource(item) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { onMounted, ref, watch } from 'vue'

const scrollUserBoxRef = ref(null)
const props = defineProps({
  userList: {
    type: Array,
    default: null
  },
  channelItem: {
    type: Array,
    default: null
  },
})

const emit = defineEmits(['userClick', 'userScroll', 'userScrollStart', 'userScrollEnd'])

const userLists = ref([])
const channelItems = ref([])
const currentItem = ref({})

const search = (item) => {
  currentItem.value = item
  emit('userClick', item)
}

const scrollOption = {
  scrollTop: 0,
  scrollHeight: 0,
  clientHeight: 0,
  scrollStartDiff: 60,
  scrollEndDiff: 60,
  scrollDirection: '' // 滚动方向
}

let scrollEventTimer = null // 滚动条防抖
let onScrollEventLock = false // 时间触发锁
// 监听滚动条滚动
function onScroll(e) {
  if (onScrollEventLock) {
    return
  }

  if (scrollEventTimer) {
    clearTimeout(scrollEventTimer)
    scrollEventTimer = null
  }

  scrollEventTimer = setTimeout(() => {
    if (scrollOption.scrollTop - e.target.scrollTop > 0) {
      scrollOption.scrollDirection = 'up'
    }

    if (scrollOption.scrollTop - e.target.scrollTop < 0) {
      scrollOption.scrollDirection = 'down'
    }

    scrollOption.scrollTop = e.target.scrollTop
    scrollOption.scrollHeight = e.target.scrollHeight
    scrollOption.clientHeight = e.target.clientHeight

    emit('userScroll', { ...scrollOption })
    // 触顶
    let isAtTop = Math.abs(scrollOption.scrollTop) <= scrollOption.scrollStartDiff + 100

    if (isAtTop && scrollOption.scrollDirection === 'up') {
      onScrollStart()
    }
    // 触底
    let isAtBottom =
      Math.abs(scrollOption.scrollHeight - scrollOption.scrollTop - scrollOption.clientHeight) <=
      scrollOption.scrollEndDiff

    if (isAtBottom && scrollOption.scrollDirection === 'down') {
      onScrollEnd()
    }
  }, 50)
}

function onScrollStart() {
  // 如果消息列表为空可能是断线重连等逻辑手动清空了消息列表造成的抖动，此时不触发事件
  if (props.userList.length == 0) {
    return
  }
  emit('userScrollStart', {
    msg: props.userList[0]
  })
}

function onScrollEnd() {
  // 如果消息列表为空可能是断线重连等逻辑手动清空了消息列表造成的抖动，此时不触发事件
  if (props.userList.length == 0) {
    return
  }
  emit('userScrollEnd', {
    msg: props.userList[props.userList.length - 1]
  })
}

watch(
    () => props.userList,
    (val) => {
      userLists.value = val
      if (userLists.value.length > 0 && userLists.value.length <= 20) {
        search(userLists.value[0])
        currentItem.value = userLists.value[0]
      }
    },
    {
      immediate: true
    }
)

watch(
    () => props.channelItem,
    (val) => {
      channelItems.value = val
    },
    {
      immediate: true
    }
)

const formatSource = (item) => {
  for (let i = 0; i < channelItems.value.length; i++) {
    if (item.app_type == channelItems.value[i].app_type
        && (item.app_id || '') == channelItems.value[i].app_id) {
      return channelItems.value[i].app_name
    }
  }
}

onMounted(() => {
  userLists.value = props.userList
  channelItems.value = props.channelItem
})
</script>
<style lang="less" scoped>
.user-content {
  width: 320px;
  height: 100%;
  overflow-x: hidden;
  overflow-y: auto;
  background-color: white;

  .user-item {
    width: 100%;
    display: flex;
    padding: 12px 16px;
    align-items: center;
    justify-content: flex-start;
    border-radius: 6px;
    background: white;
    cursor: pointer;

    .user-item-left {
      margin-right: 8px;

      .user-img {
        width: 48px;
        height: 48px;
        border-radius: 32px;
        border: 1px solid var(--6, #0000000f);
      }
    }

    .user-item-right {
      width: 100%;
      height: 100%;

      .user-name-box {
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .user-name {
          display: -webkit-box;
          -webkit-box-orient: vertical;
          -webkit-line-clamp: 1;
          flex: 1 0 0;
          overflow: hidden;
          color: #262626;
          text-overflow: ellipsis;
          font-family: 'PingFang SC';
          font-size: 14px;
          font-style: normal;
          font-weight: 400;
          line-height: 22px;
        }

        .user-timer {
          color: #8c8c8c;
          text-align: right;
          font-family: 'PingFang SC';
          font-size: 14px;
          font-style: normal;
          font-weight: 400;
          line-height: 22px;
        }
      }

      .user-info {
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 1;
        align-self: stretch;
        overflow: hidden;
        color: #595959;
        text-overflow: ellipsis;
        font-family: 'PingFang SC';
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
      }

      .user-source-box {
        width: 100%;
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 1;
        overflow: hidden;
        color: #8c8c8c;
        text-overflow: ellipsis;
        font-family: 'PingFang SC';
        font-size: 12px;
        font-style: normal;
        font-weight: 400;
        line-height: 20px;

        .user-source-title {
          display: inline;
        }

        .user-source-content {
          display: inline;
        }
      }
    }
  }

  .user-item:hover {
    background-color: #F2F4F7;
  }

  .active {
    background-color: #e6efff;
  }

  .active:hover {
    background-color: #e6efff;
  }
}

.empty-box {
  text-align: center;
  height: 100%;
  padding-top: 148px;
  img {
    width: 200px;
    height: 200px;
  }
  .title {
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    color: #262626;
  }
}

/* 滚动条样式 */
.user-content::-webkit-scrollbar {
  width: 4px; /*  设置纵轴（y轴）轴滚动条 */
  height: 4px; /*  设置横轴（x轴）轴滚动条 */
}
/* 滚动条滑块（里面小方块） */
.user-content::-webkit-scrollbar-thumb {
  border-radius: 0px;
  background: transparent;
}
/* 滚动条轨道 */
.user-content::-webkit-scrollbar-track {
  border-radius: 0;
  background: transparent;
}

/* hover时显色 */
.user-content:hover::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
}
.user-content:hover::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
}
</style>
