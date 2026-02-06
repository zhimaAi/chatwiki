<template>
  <ConfigBox ref="configRef">
    <template v-if="loginStatus" #head-extra>
      <div class="main-box">
        <div class="login-info-box">
          <img class="avatar" :src="loginInfo?.headimgurl"/>
          <div class="info">
            <div class="name">{{ loginInfo?.nickname }}</div>
            <div>{{ t('label_recent_login_time') }}{{ dayjs(loginInfo.login_time * 1000).format('YYYY-MM-DD HH:mm') }}</div>
            <div>{{ t('label_login_duration') }}{{ loginInfo.login_duration_text }}</div>
          </div>
        </div>
        <a class="main-btn" :href="loginUrl" target="_blank">{{ t('btn_change_password_settings') }}</a>
      </div>
    </template>
    <template v-else #body>
      <div class="no-auth-box">
        <img src="@/assets/official-unlogin.png"/>
        <div class="tit">{{ t('title_login_required') }}</div>
        <a-button type="primary" class="link" @click="goLogin">{{ t('btn_go_login') }}</a-button>
        <div class="desc">{{ t('msg_already_logged_in') }}<a-button type="link" @click="refresh" :loading="loading">{{ t('btn_click_refresh') }}</a-button></div>
      </div>
    </template>
  </ConfigBox>
</template>

<script setup>
import {ref} from 'vue';
import dayjs from 'dayjs';
import ConfigBox from "./config-box.vue";
import {useOfficialArticleLogin} from "@/composables/useOfficialArticleLogin.js";
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.explore.plugins.components.official-article-config');

const {
  loginInfo,
  loginStatus,
  loginUrl,
  loading,
  refresh,
  goLogin,
  init,
} = useOfficialArticleLogin(false)
const configRef = ref(null)

function show(info) {
  init()
  configRef.value.show(info)
}


defineExpose({
  show,
})
</script>

<style scoped lang="less">
.main-box {
  border-radius: 6px;
  background: #F2F4F7;
  padding: 16px;
  text-align: center;

  .login-info-box {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
    text-align: left;

    .avatar {
      flex-shrink: 0;
      width: 48px;
      height: 48px;
      border-radius: 12px;
    }

    .info {
      color: #7a8699;
      font-size: 14px;
      display: flex;
      flex-direction: column;
      gap: 2px;

      .name {
        color: #262626;
      }
    }
  }

  .main-btn {
    width: 100%;
    padding: 5px 16px;
    border-radius: 6px;
    border: 1px solid #D9D9D9;
    background: #FFF;
    color: #595959;
    text-align: center;
  }
}

.no-auth-box {
  display: flex;
  flex-direction: column;
  align-items: center;

  .link {
    width: 168px;
    margin-top: 16px;
  }

  img {
    width: 200px;
    height: 200px
  }

  .tit {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
  }

  .desc {
    color: #8c8c8c;
    font-size: 14px;
    font-weight: 400;
    margin: 12px 0;

    :deep(.ant-btn) {
      padding: 0;
    }
  }
}
</style>
