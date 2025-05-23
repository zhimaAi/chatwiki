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
      line-height: 22px;
      padding: 0 4px;
      font-size: 14px;
      color: #3a4559;
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
        <span class="user-name">{{ user_name }}</span>
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
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/stores/modules/user'

// const router = useRouter()
const userStore = useUserStore()

const { userInfo, avatar, user_name } = storeToRefs(userStore)

const onLogout = () => {
  userStore.logoutConfirm(true)
}

const toSystem = () => {
  // router.push({
  //   path: '/user/model',
  // })
  window.open(`/#/user/model`, "_blank", "noopener") // 建议添加 noopener 防止安全漏洞
}
</script>
