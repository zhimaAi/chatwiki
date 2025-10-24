<style lang="less" scoped>
.wiki-dropdown {
  position: relative;
  .current-wiki {
    display: flex;
    align-items: center;
    padding: 4px 8px;
    border-radius: 6px;
    transition: all 0.2s;
    cursor: pointer;

    &:hover {
      background: #e4e6eb;
    }
    .wiki-logo {
      width: 32px;
      height: 32px;
      margin-right: 8px;
      border-radius: 4px;
    }
    .wiki-name {
      flex: 1;
      height: 24px;
      line-height: 24px;
      margin-bottom: 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      color: #262626;
      font-size: 16px;
      font-weight: 600;
    }

    .action-down {
      margin-left: 8px ;
      font-size: 14px;
      color: #595959;
    }
  }

  .wike-dropdown-menu {
    position: absolute;
    top: 40px;
    left: 0;
    width: 324px;
    max-height: 400px;
    padding: 2px;
    border-radius: 12px;
    background: #fff;
    z-index: 100;
    overflow-y: auto;
    box-shadow:
      0 6px 30px 5px rgba(0, 0, 0, 0.05),
      0 16px 24px 2px rgba(0, 0, 0, 0.04),
      0 8px 10px -5px rgba(0, 0, 0, 0.08);
  }

  .wike-menu-item {
    display: flex;
    align-items: center;
    width: 100%;
    height: 40px;
    padding: 0 8px;
    margin-bottom: 4px;
    border-radius: 6px;
    transition: all 0.2s;
    cursor: pointer;
    &:last-child {
      margin-bottom: 0;
    }

    &:hover {
      background: #e4e6eb;
    }

    .wiki-logo {
      width: 32px;
      height: 32px;
      margin-right: 8px;
      border-radius: 4px;
    }
    .wiki-name {
      flex: 1;
      height: 24px;
      line-height: 24px;
      margin-bottom: 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      color: #262626;
      font-size: 16px;
      font-weight: 600;
    }
    .check-icon {
      display: none;
      margin-left: 8px;
      font-size: 14px;
      color: #2475fc;
    }
  }
  .wike-menu-item.active {
    background: #e6efff;
    .check-icon {
      display: block;
    }
  }

  .wike-dropdown-wrapper{
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    z-index: 99;
  }
}
</style>

<template>
  <div class="wiki-dropdown">
    <div class="current-wiki" @click="toggleDropdown" ref="currentWikiRef">
      <img class="wiki-logo" :src="props.libraryInfo.avatar" v-if="props.libraryInfo.avatar" />
      <img class="wiki-logo" src="@/assets/img/t.svg" alt="" v-else />
      <h3 class="wiki-name">{{ props.libraryInfo.library_name }}</h3>
      <span class="action-down" v-if="props.libraryList.length > 1">
        <DownOutlined />
      </span>
    </div>
    <div class="wike-dropdown-wrapper" v-if="isDropdownOpen"></div>
    <div class="wike-dropdown-menu custom-scrollbar" v-if="isDropdownOpen" ref="dropdownMenuRef">
      <div
        class="wike-menu-item"
        :class="{ active: props.libraryInfo.libraryKey === item.library_key }"
        v-for="item in props.libraryList"
        :key="item.library_key"
        @click="changeLibrary(item)"
      >
        <img class="wiki-logo" :src="item.avatar" v-if="item.avatar" />
        <img class="wiki-logo" src="@/assets/img/t.svg" alt="" v-else />
        <h3 class="wiki-name">{{ item.library_name }}</h3>
        <span class="check-icon"><CheckOutlined /></span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { DownOutlined, CheckOutlined } from '@ant-design/icons-vue'
import { ref, onMounted, onUnmounted } from 'vue'

const emit = defineEmits(['change'])

const props = defineProps({
  libraryList: {
    type: Array,
    default: () => []
  },
  libraryInfo: {
    type: Object,
    default: () => {
      return {
        library_key: '',
        library_name: '',
        library_intro: '',
        avatar: '',
      }
    }
  }
})

// 控制下拉菜单显示状态
const isDropdownOpen = ref(false)
const currentWikiRef = ref(null)
const dropdownMenuRef = ref(null)

// 切换下拉菜单显示/隐藏
const toggleDropdown = (event) => {
  if (props.libraryList.length == 1) {
    return
  }
  // 阻止事件冒泡，避免触发document的click事件
  event.stopPropagation()
  isDropdownOpen.value = !isDropdownOpen.value
}

// 点击外部区域关闭下拉菜单
const handleClickOutside = (event) => {
  // 判断点击是否在当前组件之外
  if (
    currentWikiRef.value &&
    !currentWikiRef.value.contains(event.target) &&
    dropdownMenuRef.value &&
    !dropdownMenuRef.value.contains(event.target)
  ) {
    isDropdownOpen.value = false
  }
}

// 切换选中的知识库
const changeLibrary = (item) => {
  // 切换后关闭下拉菜单
  isDropdownOpen.value = false
  emit('change', item)
}

// 组件挂载时添加点击事件监听
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

// 组件卸载时移除点击事件监听
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
