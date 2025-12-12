<style lang="less" scoped>
.user-dropdown {
  .user-dropdown-link {
    display: flex;
    align-items: center;
    cursor: pointer;

    .user-avatar {
      width: 24px;
      height: 24px;
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
        <span class="user-name">{{ displayUserName }}</span>
        <svg-icon name="arrow-down" style="font-size: 16px; color: #8c8c8c"></svg-icon>
      </div>
      <template #overlay>
        <a-menu>
          <a-menu-item>
            <a class="menu-item" href="javascript:;" @click="toSystem">
              <svg-icon class="menu-icon" name="system-setting"></svg-icon>
              <span class="menu-name">系统管理</span>
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
            <a class="menu-item" href="javascript:;" @click="onLogout">
              <svg-icon class="menu-icon" name="logout"></svg-icon>
              <span class="menu-name">退出登录</span>
            </a>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>

<script setup>
// import { useRouter } from 'vue-router'
import { computed, } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/modules/user'
import { useCompanyStore } from '@/stores/modules/company'
import { checkRole } from '@/utils/permission'
const companyStore = useCompanyStore()
const router = useRouter()
// const router = useRouter()
const userStore = useUserStore()

const { userInfo, avatar, user_name } = storeToRefs(userStore)

const displayUserName = computed(() => {
  return user_name.value.length > 11 ? user_name.value.substring(0, 11) + '...' : user_name.value
})

const onLogout = () => {
  userStore.logoutConfirm(true)
}

const toSystem = () => {
  // router.push({
  //   path: '/user/model',
  // })
  window.open(`/#/user/model`, "_blank", "noopener") // 建议添加 noopener 防止安全漏洞
}

const baseNavs = [
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
    permission: ['OpenLibDocManage']
  },
  {
    id: 4,
    key: 'library-search',
    label: 'library-search',
    title: '搜索',
    icon: 'search',
    path: '/library-search/index',
    permission: ['SearchManage']
  },
  {
    id: 5,
    key: 'chat-monitor',
    label: 'chat-monitor',
    title: '会话',
    icon: 'nav-chat',
    path: '/chat-monitor/index',
    permission: ['ChatSessionManage']
  },
]

const top_navigate = computed(() => {
  return companyStore.top_navigate
})

const navs = computed(() => {
  const closeList = top_navigate.value.filter(item => !item.open); // 获取所有关闭的菜单项

  return closeList
    .map(item => baseNavs.find(nav => nav.key === item.id)) // 查找匹配的菜单项
    .filter(Boolean); // 过滤掉未找到的菜单项（undefined）
});

const handleClickNav = item => {
  router.push(item.path)
}


</script>
