<style lang="less" scoped>
.navbar {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #FAFAFA;
  border-radius: 6px;

  .nav-menu {
    display: flex;
    position: relative;
    padding: 9px 16px;
    margin-right: 4px;
    line-height: 22px;
    font-size: 14px;
    font-weight: 700;
    border-radius: 6px;
    color: #262626;
    cursor: pointer;
    transition: all .2s;
    &:hover{
      background: #E4E6EB; 
    }
    &.active {
      color: #fff;
      background: #2475FC;
    }

    .nav-icon {
      margin-right: 8px;
      font-size: 14px; 
    }
  }

  
}
</style>

<template>
  <div class="navbar-wrapper">
    <div class="navbar">
      <template v-for="item in navs">
        <div
          class="nav-menu"
          :class="{ active: item.key === rootPath || item.key === activeMenu }"
          :key="item.key"
          @click="handleClickNav(item)"
          v-if="checkRole(item.permission)"
        >
          <svg-icon class="nav-icon" :name="item.icon"></svg-icon>
          <span class="nav-name">{{ item.title }}</span>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { checkRole } from '@/utils/permission'

const router = useRouter()
const roure = useRoute()

const activeMenu = computed(() => {
  return roure.meta.activeMenu || ''
})

const rootPath = computed(() => {
  return roure.path.split('/')[1]
})

const getActiveMenu = () => {}

const navs = [
  {
    id: 1,
    key: 'robot',
    label: 'robot',
    title: '应用',
    icon: 'nav-robot',
    path: '/robot/list',
    permission: ['RobotManage']
  },
  {
    id: 2,
    key: 'library',
    label: 'library',
    title: '知识库',
    icon: 'nav-library',
    path: '/library/list',
    permission: ['LibraryManage']
  },
  {
    id: 3,
    key: 'PublicLibrary',
    label: 'PublicLibrary',
    title: '文档',
    icon: 'nav-doc',
    path: '/public-library/list',
    permission: ['*:*:*']
  },
  {
    id: 4,
    key: 'library-search',
    label: 'library-search',
    title: '搜索',
    icon: 'search',
    path: '/library-search/index',
    permission: ['LibrarySearch']
  },
  {
    id: 5,
    key: 'chat-monitor',
    label: 'chat-monitor',
    title: '会话',
    icon: 'nav-chat',
    path: '/chat-monitor/index',
    permission: ['*:*:*']
  },
  // {
  //   id: 6,
  //   key: 'user',
  //   label: 'user',
  //   title: '系统管理',
  //   path: '/user/model',
  //   permission: [
  //     'ModelManage',
  //     'TokenManage',
  //     'TeamManage',
  //     'AccountManage',
  //     'CompanyManage',
  //     'ClientSideManage'
  //   ]
  // }
]

const handleClickNav = item => {
  router.push(item.path)
}

watch(
  () => roure.path,
  () => {
    getActiveMenu()
  },
  {
    immediate: true
  }
)
</script>
