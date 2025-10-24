<style lang="less" scoped>
.wiki-sidebar-warpper {
  height: 100%;

  .wiki-sidebar-mask {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    width: 100%;
    height: 100%;
    z-index: 999;
    background-color: rgba(0, 0, 0, 0.65);
  }
}
/* 侧边栏样式 */
.wiki-sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  .wiki-sidebar-header {
    padding: 24px 8px;
  }

  .wiki-search-box {
    position: relative;
    padding: 8px;
    background-color: #fff;
    .search-input-box {
      position: relative;
    }
    .search-input {
      width: 100%;
      height: 40px;
      line-height: 40px;
      padding: 0 28px 0 8px;
      border: none;
      border-radius: 6px;
      font-size: 14px;
      box-sizing: border-box;
      background-color: #f2f4f7;
      outline: none;
    }
    .search-icon {
      position: absolute;
      right: 12px;
      width: 16px;
      height: 16px;
      top: 50%;
      transform: translateY(-50%);
      cursor: pointer;
    }
  }

  .sidebar-menus {
    padding: 0 8px 8px;
    .sidebar-menu-item {
      height: 40px;
      border-radius: 6px;
    }
    .sidebar-menu-item:hover {
      background-color: #f2f4f7;
    }
    .sidebar-menu-item .link {
      display: flex;
      align-items: center;
      width: 100%;
      height: 100%;
      padding: 0 16px;
      text-decoration: none;
    }
    .sidebar-menu-item .menu-icon {
      width: 16px;
      height: 16px;
      margin-right: 8px;
    }
    .sidebar-menu-item .default-icon {
      display: block;
    }
    .sidebar-menu-item .active-icon {
      display: none;
    }
    .sidebar-menu-item .menu-name {
      flex: 1;
      font-size: 14px;
      color: #595959;
    }
    .sidebar-menu-item.active {
      background-color: #e6efff;
    }
    .sidebar-menu-item.active .default-icon {
      display: none;
    }
    .sidebar-menu-item.active .active-icon {
      display: block;
    }
    .sidebar-menu-item.active .menu-name {
      color: #2475fc;
    }
  }

  .sidebar-directory {
    flex: 1;
    overflow: hidden;
    overflow-y: auto;
    /* 使用全局定义的滚动条样式 */

    .directory-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      height: 40px;
      padding: 0 8px 0 24px;
    }
    .directory-header .directory-label {
      font-size: 14px;
      color: #262626;
    }
    .directory-header .action-btn {
      cursor: pointer;
    }
    .directory-body {
      padding: 0 8px;
    }
  }
}

@media (max-width: 992px) {
  .wiki-sidebar {
    position: fixed;
    top: 0;
    bottom: 0;
    left: -320px;
    width: 320px;
    height: 100%;
    z-index: 1000;
    background-color: #fff;
    overflow-y: auto;
    transition: left 0.3s ease;
  }

  .wiki-sidebar-open {
    .wiki-sidebar {
      left: 0;
    }
    .wiki-sidebar-mask {
      display: block;
    }
  }
}
</style>

<template>
  <div
    class="wiki-sidebar-warpper"
    id="wikiSidebarWrapper"
    :class="{ 'wiki-sidebar-open': isOpen }"
  >
    <div class="wiki-sidebar-mask" @click="toggleSidebar()"></div>
    <div class="wiki-sidebar">
      <div class="wiki-sidebar-header">
        <WikDropdown :previewKey="previewKey" />
      </div>

      <div class="wiki-search-box">
        <div class="search-input-box">
          <input
            id="sidebar-search-input"
            class="search-input"
            ref="searchInputRef"
            type="text"
            placeholder=""
            v-model="state.keyword"
            @keydown.enter="onSidebarSearch()"
          />
          <img class="search-icon" src="@/assets/img/search.png" @click="onSidebarSearch()" />
        </div>
      </div>

      <div class="sidebar-menus">
        <div class="sidebar-menu-item" :class="{ active: activeKey == libraryInfo.library_key }">
          <router-link
            class="link sidebar-link"
            id="home-link"
            :to="`/home/${libraryInfo.library_key}?${props.previewKey ? 'preview=' + props.previewKey : ''}${token ? '&token=' + token : ''}`"
          >
            <img class="menu-icon default-icon" src="@/assets/img/menu_ai.svg" alt="" />
            <img class="menu-icon active-icon" src="@/assets/img/menu_ai_active.svg" alt="" />
            <span class="menu-name">首页</span>
          </router-link>
        </div>
      </div>

      <div class="sidebar-directory custom-scrollbar">
        <div class="directory-header">
          <span class="directory-label">文档目录</span>
          <div class="action-box">
            <img
              class="action-btn"
              src="@/assets/img/directory_expanded.svg"
              :title="allExpanded ? '收起全部菜单' : '展开全部菜单'"
              @click="handelToggleAllCatalog()"
            />
          </div>
        </div>

        <div class="directory-body">
          <div class="directory-list" id="directory-list">
            <DirectoryMenu
              :activeKey="props.activeKey"
              :previewKey="previewKey"
              :token="token"
              :items="catalog"
              :iconTemplateConfig="iconTemplateConfig"
              :forceExpanded="allExpanded && !!forceToggleTimestamp"
              :forceCollapsed="!allExpanded && !!forceToggleTimestamp"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { message } from 'ant-design-vue'
import { ref, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useOpenDocStore } from '@/stores/open-doc'
import DirectoryMenu from './DirectoryMenu.vue'
import WikDropdown from './WikDropdown.vue'

const props = defineProps({
  activeKey: {
    type: String,
    default: '',
  },
  previewKey: {
    type: String,
    default: '',
  },
})

const router = useRouter()

const searchInputRef = ref(null)
const state = reactive({
  keyword: '',
})

const openDocStore = useOpenDocStore()

const token = computed(() => {
  return openDocStore.token
})

const libraryInfo = computed(() => {
  return openDocStore.libraryInfo
})

const iconTemplateConfig = computed(() => {
  return openDocStore.iconTemplateConfig
})

const catalog = computed(() => {
  return openDocStore.catalog
})

const isOpen = computed(() => {
  return openDocStore.sidebarOpen
})

const allExpanded = computed(() => {
  return openDocStore.allExpanded
})

const forceToggleTimestamp = computed(() => {
  return openDocStore.forceToggleTimestamp
})

const toggleSidebar = () => {
  openDocStore.toggleSidebar()
}

const handelToggleAllCatalog = () => {
  openDocStore.toggleAllCatalog()
}

const onSidebarSearch = () => {
  if (openDocStore.isEditPage) {
    message.warning('编辑模式不支持搜索功能')
    return
  }

  if (props.previewKey) {
    message.warning('预览模式不支持搜索功能')
    return
  }
  // 失去焦点
  searchInputRef.value.blur()
  router.push({
    name: 'open-search',
    params: {
      id: libraryInfo.value.library_key,
    },
    query: {
      v: state.keyword,
    },
  })

  state.keyword = ''
}
</script>
