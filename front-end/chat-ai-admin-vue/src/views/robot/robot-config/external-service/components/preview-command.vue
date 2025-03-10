<template>
  <div class="command-box" ref="buttonsContainer">
    <div v-for="(item, index) in visibleButtons" :key="index" class="fast-item">
      {{ item.short_title }}
    </div>
    <!-- <div v-if="showMoreButton" class="fast-item more">更多 <DownOutlined /></div> -->
  </div>
</template>

<script setup>
import { ref, computed, reactive, nextTick, watch } from 'vue'
// import { DownOutlined } from '@ant-design/icons-vue'
import { useRobotStore } from '@/stores/modules/robot'
import { storeToRefs } from 'pinia'
const robotStore = useRobotStore()
const { quickCommandLists } = storeToRefs(robotStore)

const buttonsContainer = ref(null)

const buttons = computed(() => {
  return quickCommandLists.value.map((item) => {
    return {
      short_title: item.title,
      ...item
    }
  })
})

watch(buttons, (val) => {
  updateButtons()
})

const visibleButtons = ref([])
// const showMoreButton = ref(false)

const updateButtons = () => {
  nextTick(() => {
    // const containerWidth = buttonsContainer.value.clientWidth
    // let totalWidth = 80
    visibleButtons.value = []
    for (let i = 0; i < buttons.value.length; i++) {
      const buttonElement = document.createElement('div')
      buttonElement.classList.add('fast-item')
      buttonElement.textContent = buttons.value[i].short_title
      buttonsContainer.value.appendChild(buttonElement)
      // const buttonWidth = buttonElement.offsetWidth + 34 // 10px 是gap
      buttonsContainer.value.removeChild(buttonElement)

      // if (totalWidth + buttonWidth > containerWidth) {
      //   showMoreButton.value = true
      //   break
      // } else {
        // totalWidth += buttonWidth
        visibleButtons.value.push(buttons.value[i])
        // showMoreButton.value = false
      // }
    }
  })
}
</script>

<style lang="less" scoped>
.command-box {
  position: absolute;
  left: 28px;
  right: 28px;
  height: 45px;
  bottom: 96px;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: nowrap;
  white-space: nowrap;
  overflow-x: auto;
  overflow-y: hidden;

  .fast-item {
    cursor: pointer;
    padding: 5px 12px;
    display: flex;
    align-items: center;
    border-radius: 8px;
    border: 1px solid #d9d9d9;
    background: #fff;
    color: #595959;
    font-size: 14px;
    line-height: normal;
  }
}

/* 滚动条样式 */
.command-box::-webkit-scrollbar {
    width: 4px; /*  设置纵轴（y轴）轴滚动条 */
    height: 4px; /*  设置横轴（x轴）轴滚动条 */
}
/* 滚动条滑块（里面小方块） */
.command-box::-webkit-scrollbar-thumb {
    border-radius: 0px;
    background: transparent;
}
/* 滚动条轨道 */
.command-box::-webkit-scrollbar-track {
    border-radius: 0;
    background: transparent;
}

/* hover时显色 */
.command-box:hover::-webkit-scrollbar-thumb {
    background: rgba(0,0,0,0.2);
}
.command-box:hover::-webkit-scrollbar-track {
    background: rgba(0,0,0,0.1);
}
</style>