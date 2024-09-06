<style lang="less" scoped>
.left-menus {
  padding: 8px;

  .menu-item {
    height: 40px;
    padding: 0 16px;
    display: flex;
    align-items: center;
    font-size: 14px;
    color: #595959;
    cursor: pointer;
    transition: all 0.2s;
    &:hover {
      color: #2475fc;
      .menu-icon {
        color: #2475fc;
      }
    }

    .menu-icon {
      margin-right: 4px;
      width: 20px;
      height: 16px;
      color: #a1a7b3;
    }
    .menu-name {
      flex: 1;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }

    &.active {
      color: #2475fc;
      background-color: #e6efff;
      .menu-icon {
        color: #2475fc;
      }
    }

    .svg-action {
      font-size: 16px;
      margin-right: 4px;
      margin-top: 4px;
    }
  }
}
</style>

<template>
  <div class="left-menus">
    <div
      class="menu-item"
      :class="{ active: item.path == activeMenu }"
      v-for="item in menus"
      :key="item.key"
      @click="onChange(item)"
    >
      <span>
        <span class="menu-icon" v-if="item.svg">
          <svg-icon :name="item.svgActive" v-if="item.path == activeMenu"></svg-icon>
          <svg-icon :name="item.svg" v-else></svg-icon>
        </span>

        <component v-if="item.icon" class="menu-icon" :is="item.icon"></component>
      </span>
      <span class="menu-name">{{ item.name }}</span>
    </div>
  </div>
</template>
<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { SettingFilled, AppstoreFilled } from '@ant-design/icons-vue'
import { usePermissionStore } from '@/stores/modules/permission'

const router = useRouter()
const route = useRoute()
// createVNode('div', { style: 'color:red;' }, content),

const permissionStore = usePermissionStore()
const menus = computed(() => {
  let { role_permission } = permissionStore
  let possessedAuthority = []
  
  for (let i = 0; i < role_permission.length; i++) {
    const item = role_permission[i];
    if (item === 'ModelManage') {
      possessedAuthority.push({
        id: 1,
        name: '模型管理',
        key: 'model',
        path: '/user/model',
        icon: AppstoreFilled
      })
    }
    if (item === 'TokenManage') {
      possessedAuthority.push({
        id: 2,
        name: 'Token使用',
        key: 'usetoken',
        path: '/user/usetoken',
        svg: 'use-token',
        svgActive: 'use-token-active'
      })
    }
    if (item === 'TeamManage') {
      possessedAuthority.push({
        id: 3,
        name: '团队管理',
        key: 'manage',
        path: '/user/manage',
        svg: 'team-manage-active',
        svgActive: 'team-manage'
      })
    }
    if (item === 'AccountManage') {
      possessedAuthority.push({
        id: 4,
        name: '账号设置',
        key: 'account',
        path: '/user/account',
        icon: SettingFilled
      })
    }
    if (item === 'CompanyManage') {
      possessedAuthority.push({
        id: 5,
        name: '企业设置',
        key: 'enterprise',
        path: '/user/enterprise',
        svg: 'enterprise',
        svgActive: 'enterprise-active'
      })
    }
    if (item === 'ClientSideManage') {
      possessedAuthority.push({
        id: 6,
        name: '客户端下载',
        key: 'clientDownload',
        path: '/user/clientDownload',
        svg: 'client',
        svgActive: 'client'
      })
    }
  }
  return possessedAuthority.sort((a, b) => a.id - b.id);
})

const activeMenu = computed(() => {
  return route.path
})

const onChange = (item) => {
  if (item.path == activeMenu.value) {
    return
  }

  router.push(item.path)
}
</script>
