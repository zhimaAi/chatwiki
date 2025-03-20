<style lang="less" scoped>
.navbar {
  display: flex;
  justify-content: center;
  align-items: center;

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
      <div
        class="nav-menu"
        :class="{ active: item.key === rootPath || item.key === activeMenu }"
        v-for="item in items"
        :key="item.key"
      >
        <router-link :to="item.path" class="nav-menu-name">{{ item.title }}</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { usePermissionStore } from '@/stores/modules/permission'

// front-end\chat-ai-admin-vue\src\utils\permission.js
const roure = useRoute()

const activeMenu = computed(() => {
  return roure.meta.activeMenu || ''
})

const rootPath = computed(() => {
  return roure.path.split('/')[1]
})
const permissionStore = usePermissionStore()
const getActiveMenu = () => {}

const items = computed(() => {
  const SystemManageChildren = [
    'ModelManage',
    'TokenManage',
    'TeamManage',
    'AccountManage',
    'CompanyManage',
    'ClientSideManage'
  ]
  let flag = false // 控制添加 "系统管理" 菜单的，如果添加过就不进行循环和添加。
  let { role_permission } = permissionStore
  let possessedAuthority = []

  for (let i = 0; i < role_permission.length; i++) {
    const item = role_permission[i]
    if (item === 'RobotManage') {
      possessedAuthority.push({
        id: 1,
        key: 'robot',
        label: 'robot',
        title: '应用',
        path: '/robot/list'
      })
    }
    if (item === 'LibraryManage') {
      possessedAuthority.push({
        id: 2,
        key: 'library',
        label: 'library',
        title: '知识库',
        path: '/library/list'
      })
    }

    if (item === 'OpenDoc') {
      // 插入对外文档
      possessedAuthority.push({
        id: 3,
        key: 'PublicLibrary',
        label: 'PublicLibrary',
        title: '对外文档',
        path: '/public-library/list'
      })
    }

    if (item === 'FormManage') {
      possessedAuthority.push({
        id: 4,
        key: 'database',
        label: 'database',
        title: '数据库',
        path: '/database/list'
      })
    }
    if (!flag) {
      for (let j = 0; j < SystemManageChildren.length; j++) {
        // 作用是看系统管理里面有没有子权限，有一个则显示系统管理，否则不显示系统管理菜单。
        const child = SystemManageChildren[j]
        if (child === item) {
          possessedAuthority.push({
            id: 5,
            key: 'user',
            label: 'user',
            title: '系统管理',
            path: '/user/model'
          })
          flag = true
          break
        }
      }
    }
  }
  return possessedAuthority.sort((a, b) => a.id - b.id)
})

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
