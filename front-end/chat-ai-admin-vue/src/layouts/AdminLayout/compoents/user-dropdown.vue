<style lang="less" scoped>
.user-dropdown {
  .user-dropdown-link {
    display: flex;
    align-items: center;
    cursor: pointer;

    .user-avatar {
      width: 36px;
      height: 36px;
      border-radius: 50%;
    }

    .user-name {
      max-width: 200px;
      line-height: 22px;
      padding: 0 4px;
      font-size: 14px;
      color: #3a4559;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }
}

.menu-item{
  display: flex;
  align-items: center;
  padding: 2px 8px;
  color: #262626;
  .menu-icon{
    margin-right: 8px;
    font-size: 16px;
  }
}
</style>

<template>
  <div class="user-dropdown" v-if="userInfo">
    <a-dropdown>
      <div class="user-dropdown-link" @click.prevent>
        <img class="user-avatar" :src="avatar" alt="" />
        <!-- <span class="user-name">{{ displayUserName }}</span>
        <svg-icon name="arrow-down" style="font-size: 16px; color: #8c8c8c"></svg-icon> -->
      </div>
      <template #overlay>
        <a-menu>
          <a-menu-item>
            <a class="menu-item" href="javascript:;" @click="toSystem">
              <svg-icon class="menu-icon" name="system-setting"></svg-icon>
              <span class="menu-name">{{ t('layout.navbar.system_manage') }}</span>
            </a>
          </a-menu-item>
          <template v-for="item in navs" :key="item.id">
            <a-menu-item v-if="checkRole(item.permission)">
              <a class="menu-item" href="javascript:;" @click="handleClickNav(item)">
                <svg-icon class="menu-icon" :name="item.icon"></svg-icon>
                <span class="menu-name">{{ item.title }}</span>
              </a>
            </a-menu-item>
          </template>
          <a-menu-item>
            <LocaleDropdown />
          </a-menu-item>
          <a-menu-item>
            <a class="menu-item" href="javascript:;" @click="onLogout">
              <svg-icon class="menu-icon" name="logout"></svg-icon>
              <span class="menu-name">{{ t('layout.navbar.logout') }}</span>
            </a>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import { useUserStore } from '@/stores/modules/user'
import { useCompanyStore } from '@/stores/modules/company'
import { checkRole } from '@/utils/permission'
import LocaleDropdown from './locale-dropdown.vue'

const { t } = useI18n()
const companyStore = useCompanyStore()
const router = useRouter()
const userStore = useUserStore()

const { userInfo, avatar, user_name } = storeToRefs(userStore)

const displayUserName = computed(() => {
  return user_name.value.length > 11 ? user_name.value.substring(0, 11) + '...' : user_name.value
})

const onLogout = () => {
  userStore.logoutConfirm(true)
}

const toSystem = () => {
  window.open(`/#/user/model`, "_blank", "noopener")
}

const baseNavs = computed(() => [
  {
    id: 1,
    key: 'robot',
    label: 'robot',
    title: t('layout.navbar.application'),
    icon: 'nav-robot',
    path: '/robot/list',
    permission: ['RobotManage']
  },
  {
    id: 2,
    key: 'library',
    label: 'library',
    title: t('layout.navbar.knowledge_base'),
    icon: 'nav-library',
    path: '/library/list',
    permission: ['LibraryManage']
  },
  {
    id: 3,
    key: 'PublicLibrary',
    label: 'PublicLibrary',
    title: t('layout.navbar.document'),
    icon: 'nav-doc',
    path: '/public-library/list',
    permission: ['OpenLibDocManage']
  },
  {
    id: 4,
    key: 'library-search',
    label: 'library-search',
    title: t('layout.navbar.search'),
    icon: 'search',
    path: '/library-search/index',
    permission: ['SearchManage']
  },
  {
    id: 5,
    key: 'chat-monitor',
    label: 'chat-monitor',
    title: t('layout.navbar.session'),
    icon: 'nav-chat',
    path: '/chat-monitor/index',
    permission: ['ChatSessionManage']
  },
])

const top_navigate = computed(() => {
  return companyStore.top_navigate
})

const navs = computed(() => {
  const closeList = top_navigate.value.filter(item => !item.open);

  return closeList
    .map(item => baseNavs.value.find(nav => nav.key === item.id))
    .filter(Boolean);
});

const handleClickNav = item => {
  router.push(item.path)
}


</script>
