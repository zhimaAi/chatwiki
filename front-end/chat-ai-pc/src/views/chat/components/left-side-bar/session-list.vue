<template>
  <div class="session-list-box">
    <div class="scroll-box" @scroll="handleScroll">
      <div
        @click="handleOpenChat(item)"
        class="list-item"
        :class="{ active: dialogue_id == item.id }"
        v-for="(item, index) in lists"
        :key="index"
      >
        <svg-icon class="session-icon" name="session-icon"></svg-icon>
        <div class="content-text">{{ item.subject }}</div>
        <div @click.stop>
          <van-popover
            v-model:show="item.showPopover"
            :actions="actions"
            @select="(action) => onSelect(action, item)"
          >
            <template #reference>
              <div class="operate-btn">
                <svg-icon name="point-h"></svg-icon>
              </div>
            </template>
          </van-popover>
        </div>
      </div>
    </div>
    <van-dialog
      :show="showDialog"
      title="重命名"
      @confirm="handleSave"
      @cancel="showDialog = false"
      show-cancel-button
    >
      <div class="input-box">
        <van-field v-model="textValue" placeholder="请输入" />
      </div>
    </van-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { showConfirmDialog, showToast } from 'vant'
import { useChatStore } from '@/stores/modules/chat'
const chatStore = useChatStore()

const emit = defineEmits(['handleOpenChat'])
const showDialog = ref(false)

const lists = computed(() => {
  return chatStore.myChatList || []
})

const dialogue_id = computed(() => {
  return chatStore.dialogue_id
})
const handleOpenChat = (item) => {
  emit('handleOpenChat', item)
}

const actions = [
  { text: '重命名', key: 1 },
  { text: '删除', key: 2 }
]

let currentItem = {}
const textValue = ref('')

const onSelect = (action, item) => {
  let key = action.key
  currentItem = item
  if (key == 1) {
    textValue.value = item.subject
    showDialog.value = true
  }
  if (key == 2) {
    handleDel(item)
  }
}
const handleSave = () => {
  if (textValue.value == '') {
    return showToast('请输入新的名称')
  }
  chatStore
    .editDialogueChat({
      ...currentItem,
      subject: textValue.value
    })
    .then((res) => {
      if (res.res == 0) {
        showToast('修改成功')
        showDialog.value = false
      }
    })
}

const handleDel = (item) => {
  showConfirmDialog({
    title: '删除确认',
    message: '确认删除该记录'
  })
    .then(() => {
      chatStore.delDialogue(item)
    })
    .catch(() => {})
}
const handleScroll = (event) => {
  const { scrollTop, scrollHeight, clientHeight } = event.target
  // 判断是否滚动到底部，这里设置距离底部20px内就算触底
  if (scrollHeight - scrollTop - clientHeight < 50) {
    console.log('滚动触底')
    // 在这里可以触发加载更多数据的逻辑
    chatStore.getMyChatList()
  }
}
</script>

<style lang="less" scoped>
.session-list-box {
  flex: 1;
  overflow: hidden;
  padding-top: 12px;
}

.scroll-box {
  height: 100%;
  overflow-y: auto;
  // 隐藏滚动条但仍可滚动
  -ms-overflow-style: none; /* IE and Edge */
  scrollbar-width: none; /* Firefox */
}

.scroll-box::-webkit-scrollbar {
  width: 0px;
  background: transparent; /* Chrome, Safari, Opera */
}

.list-item {
  padding: 5px 8px;
  border-radius: 6px;
  overflow: hidden;
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  color: #595959;
  margin-top: 4px;

  &:hover {
    background: var(--09, #f2f4f7);
  }
  &.active {
    background: var(--01-, #e5efff);
    color: #2475fc;
  }
  .session-icon {
    font-size: 16px;
    margin-top: 2px;
  }
  .content-text {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    font-size: 14px;
    line-height: 22px;
  }
  .operate-btn {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    color: #595959;
    border-radius: 6px;
    &:hover {
      background: var(--07, #e4e6eb);
    }
  }
}

.input-box {
  border-radius: 6px;
  margin: 24px;
  .van-cell {
    background: #f3f3f3;
    border-radius: 6px;
  }
}
</style>
