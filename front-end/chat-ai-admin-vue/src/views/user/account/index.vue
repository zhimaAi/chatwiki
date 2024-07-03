<style lang="less" scoped>
.user-account {
  width: 100%;
  height: 100%;
  padding: 24px;
  background-color: #fff;
  .page-title {
    line-height: 24px;
    font-size: 16px;
    font-weight: 600;
  }
  .user-avatar-box {
    position: relative;
    width: 120px;
    height: 120px;
    margin: 24px auto 0 auto;
    border-radius: 50%;
    overflow: hidden;
    .user-avatar {
      width: 100%;
      height: 100%;
    }
    .change-avatar {
      position: absolute;
      left: 0;
      right: 0;
      top: 0;
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: rgba(0, 0, 0, 0.45);
      opacity: 0;
      transition: opacity 0.2s;

      &:hover {
        opacity: 1;
      }

      .change-btn {
        width: 40px;
        cursor: pointer;
      }
    }
  }

  .user-info-items {
    .user-info-item {
      display: flex;
      justify-content: space-between;
      line-height: 22px;
      padding-bottom: 16px;
      .item-left {
        display: flex;
      }
      .item-label {
        padding-right: 32px;
        font-size: 14px;
        color: rgb(140, 140, 140);
      }
      .item-value {
        font-size: 14px;
        color: rgb(38, 38, 38);
      }
    }
  }
}
</style>

<template>
  <div class="user-account">
    <div class="page-title">{{ t('views.user.account.accountSettings') }}</div>

    <div class="user-avatar-box">
      <img class="user-avatar" :src="avatar" alt="" />
      <span class="change-avatar" @click="handleSetAvatar">
        <img class="change-btn" src="@/assets/img/user/account/colour-edit.svg" alt="" srcset="" />
      </span>
    </div>

    <div class="user-info-items">
      <div class="user-info-item">
        <div class="item-left">
          <span class="item-label">{{ t('views.user.account.loginAccount') }}</span>
          <span class="item-value">{{ user_name }}</span>
        </div>
        <div class="item-right">
          <!-- <a>修改</a> -->
        </div>
      </div>

      <div class="user-info-item">
        <div class="item-left">
          <span class="item-label">{{ t('views.user.account.accountPassword') }}</span>
          <span class="item-value">**************</span>
        </div>
        <div class="item-right" @click="handleOpenModifyAccountModal">
          <a>{{ t('common.change') }}</a>
        </div>
      </div>
    </div>
    <ModifyAccount @ok="onLogout" ref="modifyAccountRef"></ModifyAccount>
    <ModifyAvatar @ok="handleChangeAvatar" ref="modifyAvatarRef"></ModifyAvatar>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/stores/modules/user'
import { useI18n } from '@/hooks/web/useI18n'
import ModifyAccount from './components/modify-account.vue'
import ModifyAvatar from './components/modify-avatar.vue'

const { t } = useI18n()

const userStore = useUserStore()

const { userInfo, avatar, user_name } = storeToRefs(userStore)

const modifyAccountRef = ref(null)
const handleOpenModifyAccountModal = () => {
  modifyAccountRef.value.open(userInfo.value)
}

const onLogout = () => {
  userStore.reset(true)
}
const modifyAvatarRef = ref(null)
const handleSetAvatar = () => {
  modifyAvatarRef.value.setAvatar(userInfo.value)
}

const handleChangeAvatar = (url) => {
  userStore.setAvatar(url)
}
</script>
