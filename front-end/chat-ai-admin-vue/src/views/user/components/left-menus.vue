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

let baseMenu = [
  {
    name: '模型管理',
    key: 'model',
    path: '/user/model',
    icon: AppstoreFilled,
    permissionKey: 'ModelManage'
  },
  {
    name: 'Token使用',
    key: 'usetoken',
    path: '/user/usetoken',
    svg: 'use-token',
    svgActive: 'use-token-active',
    permissionKey: 'TokenManage'
  },
  {
    name: '团队管理',
    key: 'manage',
    path: '/user/manage',
    svg: 'team-manage-active',
    svgActive: 'team-manage',
    permissionKey: 'TeamManage'
  },
  {
    name: '账号设置',
    key: 'account',
    path: '/user/account',
    icon: SettingFilled,
    haspermise: true
  },
  {
    name: '阿里云OCR',
    key: 'enterprise',
    path: '/user/aliocr',
    svg: 'ali-ocr',
    svgActive: 'ali-ocr',
    permissionKey: 'AliyunOCRManage'
  },
  {
    name: '企业设置',
    key: 'enterprise',
    path: '/user/enterprise',
    svg: 'enterprise',
    svgActive: 'enterprise-active',
    permissionKey: 'CompanyManage'
  },
  {
    name: '自定义域名',
    key: 'domain',
    path: '/user/domain',
    svg: 'network',
    svgActive: 'network',
    permissionKey: 'UserDomainManage'
  },
  {
    name: '客户端下载',
    key: 'clientDownload',
    path: '/user/clientDownload',
    svg: 'client',
    svgActive: 'client',
    permissionKey: 'ClientSideManage'
  },
  {
    name: '敏感词管理',
    key: 'sensitiveWords',
    path: '/user/sensitive-words',
    svg: 'sensitive-icon',
    svgActive: 'sensitive-icon',
    permissionKey: 'SensitiveWordManage'
  },
  {
    name: '提示词模板库',
    key: 'promptLibrary',
    path: '/user/prompt-library',
    svg: 'prompt-icon',
    svgActive: 'prompt-icon',
    permissionKey: 'PromptTemplateManage'
  },
  {
    name: '公众号管理',
    key: 'officialAccount',
    path: '/user/official-account',
    svg: 'wx-app-icon',
    svgActive: 'wx-app-icon',
    permissionKey: 'OfficialAccountMange'
  },
]

const menus = computed(() => {
  let { role_permission, role_type } = permissionStore
  return baseMenu.filter((item) => role_type == 1 || item.haspermise || role_permission.includes(item.permissionKey))
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
