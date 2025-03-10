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
      padding: 0 16px 0 8px;
      font-size: 14px;
      color: #3a4559;
    }
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
            <a href="javascript:;" @click="onLogout">退出登录</a>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>

<script setup>
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/stores/modules/user'
const userStore = useUserStore()

const { userInfo, avatar, user_name } = storeToRefs(userStore)

const onLogout = () => {
  userStore.logoutConfirm(true)
}
</script>
