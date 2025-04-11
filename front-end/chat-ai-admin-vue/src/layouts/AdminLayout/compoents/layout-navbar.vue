<style lang="less" scoped>
.navbar {
  display: flex;
  align-items: center;
  margin-left: 8px;

  .nav-menu {
    position: relative;
    margin: 0 24px;
    cursor: pointer;

    .nav-menu-name {
      display: inline-block;
      line-height: 22px;
      padding: 4px 0 6px 0;
      font-size: 14px;
      font-weight: 400;
      color: #595959;
    }
  }

  .nav-menu.active {
    &::after {
      content: '';
      display: block;
      position: absolute;
      bottom: 0;
      left: 50%;
      margin-left: -20px;
      width: 40px;
      height: 2px;
      background-color: #2475fc;
    }

    .nav-menu-name {
      font-size: 14px;
      font-weight: 600;
      color: #2475fc;
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
          v-if="checkRole(item.permission)"
        >
          <router-link :to="item.path" class="nav-menu-name">{{ item.title }}</router-link>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { checkRole } from '@/utils/permission'

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
    path: '/robot/list',
    permission: ['RobotManage']
  },
  {
    id: 2,
    key: 'library',
    label: 'library',
    title: '知识库',
    path: '/library/list',
    permission: ['LibraryManage']
  },
  {
    id: 3,
    key: 'PublicLibrary',
    label: 'PublicLibrary',
    title: '对外文档',
    path: '/public-library/list',
    permission: ['*:*:*']
  },
  {
    id: 4,
    key: 'database',
    label: 'database',
    title: '数据库',
    path: '/database/list',
    permission: ['FormManage']
  },
  {
    id: 5,
    key: 'chat-monitor',
    label: 'chat-monitor',
    title: '实时会话',
    path: '/chat-monitor/index',
    permission: ['*:*:*']
  },
  {
    id: 6,
    key: 'user',
    label: 'user',
    title: '系统管理',
    path: '/user/model',
    permission: [
      'ModelManage',
      'TokenManage',
      'TeamManage',
      'AccountManage',
      'CompanyManage',
      'ClientSideManage'
    ]
  }
]

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
