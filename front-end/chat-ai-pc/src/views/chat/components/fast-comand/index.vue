<style lang="less" scoped>
.fast-comand-container {
  max-width: 738px;
  min-width: 350px;
  margin: 0 auto;
  overflow-x: auto;
  height: 40px;
  margin-bottom: 4px;;
}
.fast-comand-box {
  margin-left: 12px;
  height: 32px;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: nowrap;
}

/* 滚动条样式 */
.fast-comand-container::-webkit-scrollbar {
    width: 4px; /*  设置纵轴（y轴）轴滚动条 */
    height: 4px; /*  设置横轴（x轴）轴滚动条 */
}
/* 滚动条滑块（里面小方块） */
.fast-comand-container::-webkit-scrollbar-thumb {
    cursor: pointer;
    border-radius: 4px;
    background: transparent;
}
/* 滚动条轨道 */
.fast-comand-container::-webkit-scrollbar-track {
    border-radius: 4px;
    background: transparent;
}

/* hover时显色 */
.fast-comand-container:hover::-webkit-scrollbar-thumb {
    background: rgba(0,0,0,0.2);
}
.fast-comand-container:hover::-webkit-scrollbar-track {
    background: rgba(0,0,0,0.1);
}

.fast-item {
  white-space: nowrap;
  cursor: pointer;
  padding: 5px 12px;
  border-radius: 8px;
  border: 1px solid #d9d9d9;
  background: #fff;
  color: #595959;
  font-size: 14px;
  line-height: normal;

  &:hover {
    background: #f0f0f0;
  }
}

.all-list-wrapper {
  padding: 16px 8px 8px 8px;
  max-width: 90vw;
  max-height: 260px;
  overflow: auto;
  .title-block {
    color: #8c8c8c;
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    padding-bottom: 12px;
  }
  .list-item {
    cursor: pointer;
    margin-top: 4px;
    padding: 5px 8px;
    display: flex;
    align-items: center;
    border-radius: 6px;
    background: #fff;
    font-weight: 400;
    font-size: 14px;
    white-space: nowrap;
    .title {
      color: #262626;
      line-height: 22px;
    }
    .line {
      width: 1px;
      height: 16px;
      background: #f0f0f0;
      margin: 0 8px;
    }
    .content {
      flex: 1;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    &:hover {
      background: #f2f4f7;
    }
  }
}
.single-popover-item {
  padding: 16px;
  width: 300px;
  .title-block {
    color: #262626;
    line-height: 22px;
    font-weight: 600;
    font-size: 14px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .content-block {
    margin-top: 4px;
    color: #595959;
    font-size: 14px;
    font-weight: 400;
    line-height: 20px;
  }
}
@media (min-width: 501px) {
  .fast-comand-box {
    margin-left: 0;
  }
  .all-list-wrapper {
    max-width: 738px;
  }
}
</style>

<template>
  <div class="fast-comand-container" v-if="isShowContainer">
    <div class="fast-comand-box" ref="buttonsContainer">
      <div class="fast-comand-item" v-for="(item, index) in visibleButtons" :key="index">
        <div class="fast-item" type="default" @click="handleClickItem(item)">{{ item.elipsis_title }}</div>
      </div>
    </div>
  </div>
</template>

<script setup >
import { useChatStore } from '@/stores/modules/chat'
import { ref, onMounted, reactive, nextTick, watch, computed } from 'vue'
import { getFastCommandList } from '@/api/chat/index.js'
import { Popover } from 'vant'
import { windowWidth } from 'vant/lib/utils'
import { storeToRefs } from 'pinia'
const emit = defineEmits(['send'])
const chatStore = useChatStore()
const { getFastCommand } = chatStore
const { robot } = storeToRefs(chatStore)
const buttonsContainer = ref(null)
const buttons = ref([]);
const visibleButtons = ref([])
const showMoreButton = ref(false)
const isShowContainer = ref(false)
const xOffset = ref(0);
const updateButtons = () => {
  nextTick(() => {
    visibleButtons.value = []
    for (let i = 0; i < buttons.value.length; i++) {
      const buttonElement = document.createElement('div')
      buttonElement.classList.add('fast-item')
      buttonElement.textContent = buttons.value[i].elipsis_title
      buttonsContainer.value.appendChild(buttonElement)
      buttonsContainer.value.removeChild(buttonElement)
      visibleButtons.value.push(buttons.value[i])
      showMoreButton.value = false
    }
  })
}
const comand_list = computed(()=>{
  return robot.value.comand_list || [];
})

watch(comand_list,(val)=>{
  let lists = val || []
  lists.map((item) => {
    item.showPopover = false
    item.elipsis_title = item.title
    return item
  })
  buttons.value = lists
  isShowContainer.value = lists.length > 0
  updateButtons()
})

const showPopover = ref(false)


onMounted(async () => {
  await getFastCommand()
  window.addEventListener('resize', updateButtons)
})

const handleMouseEnter = (item) => {
  item.showPopover = true
}

const handleMouseLeave = (item) => {
  item.showPopover = false
}

const handleClickItem = (item, isShow) => {
  if(isShow){
    showPopover.value = false
  }
  if (item.typ == 1) {
    emit('send', item.content)
  }
  if (item.typ == 2) {
    window.open(item.content)
  }
}
</script>
