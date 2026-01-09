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
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const permissionStore = usePermissionStore()

const baseMenu = computed(() => [
  {
    name: t('views.user.left-menus.model_management'),
    key: 'model',
    path: '/user/model',
    icon: AppstoreFilled,
    permissionKey: 'ModelManage'
  },
  {
    name: t('views.user.left-menus.token_usage'),
    key: 'usetoken',
    path: '/user/usetoken',
    svg: 'use-token',
    svgActive: 'use-token-active',
    permissionKey: 'TokenManage'
  },
  {
    name: t('views.user.left-menus.team_management'),
    key: 'manage',
    path: '/user/manage',
    svg: 'team-manage-active',
    svgActive: 'team-manage',
    permissionKey: 'TeamManage'
  },
  {
    name: t('views.user.left-menus.account_settings'),
    key: 'account',
    path: '/user/account',
    icon: SettingFilled,
    haspermise: true
  },
  {
    name: t('views.user.left-menus.aliyun_ocr'),
    key: 'enterprise',
    path: '/user/aliocr',
    svg: 'ali-ocr',
    svgActive: 'ali-ocr',
    permissionKey: 'AliyunOCRManage'
  },
  {
    name: t('views.user.left-menus.enterprise_settings'),
    key: 'enterprise',
    path: '/user/enterprise',
    svg: 'enterprise',
    svgActive: 'enterprise-active',
    permissionKey: 'CompanyManage'
  },
  {
    name: t('views.user.left-menus.custom_domain'),
    key: 'domain',
    path: '/user/domain',
    svg: 'network',
    svgActive: 'network',
    permissionKey: 'UserDomainManage'
  },
  {
    name: t('views.user.left-menus.client_download'),
    key: 'clientDownload',
    path: '/user/clientDownload',
    svg: 'client',
    svgActive: 'client',
    permissionKey: 'ClientSideManage'
  },
  {
    name: t('views.user.left-menus.sensitive_word_management'),
    key: 'sensitiveWords',
    path: '/user/sensitive-words',
    svg: 'sensitive-icon',
    svgActive: 'sensitive-icon',
    permissionKey: 'SensitiveWordManage'
  },
  {
    name: t('views.user.left-menus.prompt_template_library'),
    key: 'promptLibrary',
    path: '/user/prompt-library',
    svg: 'prompt-icon',
    svgActive: 'prompt-icon',
    permissionKey: 'PromptTemplateManage'
  },
  {
    name: t('views.user.left-menus.official_account_management'),
    key: 'officialAccount',
    path: '/user/official-account',
    svg: 'wx-app-icon',
    svgActive: 'wx-app-icon',
    permissionKey: 'OfficialAccountMange'
  },
])

const menus = computed(() => {
  let { role_permission, role_type } = permissionStore
  return baseMenu.value.filter((item) => role_type == 1 || item.haspermise || role_permission.includes(item.permissionKey))
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
