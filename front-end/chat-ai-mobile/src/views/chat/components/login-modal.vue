<template>
  <!-- 触发按钮 -->
  <!-- <van-button type="primary" @click="showLogin = true">打开登录弹窗</van-button> -->

  <!-- 登录弹窗 -->
  <van-popup
    v-model:show="showLogin"
    round
    position="center"
    class="login-popup"
    :close-on-click-overlay="false"
    :closeable="false"
  >
    <div class="login-container">
      <h2 class="title">{{ t('title_welcome') }}</h2>

      <van-form @submit="handleLogin">
        <div class="form-item">
          <svg-icon name="user" class="input-icon" />
          <van-field
            v-model="username"
            type="text"
            :placeholder="t('ph_input_account')"
            class="custom-input"
            :rules="[{ required: true, message: t('msg_input_account') }]"
          />
        </div>

        <div class="form-item">
          <svg-icon name="password" class="input-icon" />
          <van-field
            v-model="password"
            type="password"
            :placeholder="t('ph_input_password')"
            class="custom-input"
            :rules="[{ required: true, message: t('msg_input_password') }]"
          />
        </div>

        <div class="form-item">
          <div class="info-box">
            <span>{{ t('msg_agree_terms') }}&nbsp;</span>
            <span class="info-link" @click="onGoLink('https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/kcs5ogf88ola88gk?source=aHR0cHM6Ly9jbG91ZC5jaGF0d2lraS5jb20vIy9sb2dpbj9jb2RlPS9saWJyYXJ5L2xpc3Q=')">{{ t('label_service_agreement') }}</span>
            <span>&nbsp;{{ t('label_and') }}&nbsp;</span>
            <span class="info-link" @click="onGoLink('https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/wktycu3clg16pcv6?source=aHR0cHM6Ly9jbG91ZC5jaGF0d2lraS5jb20vIy9sb2dpbj9jb2RlPS9saWJyYXJ5L2xpc3Q=')">{{ t('label_privacy_agreement') }}</span>
          </div>
        </div>

        <van-button
          block
          type="primary"
          native-type="submit"
          class="login-btn"
        >
          {{ t('btn_login') }}
        </van-button>
      </van-form>
    </div>
  </van-popup>
</template>

<script setup>
import { ref } from 'vue';
import { Popup, Button, Field, Form } from 'vant';
import { useUserStore } from '@/stores/modules/user'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.chat.components.login-modal')
const userStore = useUserStore()

const showLogin = ref(false);
const username = ref('');
const password = ref('');
const submitLoading = ref(false)

const handleLogin = () => {
  // 这里添加实际登录逻辑
  submitLoading.value = true

  userStore
  .login({
    username: username.value,
    password: password.value
  })
  .then(() => {
    submitLoading.value = false
    showLogin.value = false
  })
};

const show = () => {
  showLogin.value = true
}

const onGoLink = (url) => {
  window.open(url)
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.login-popup {
  width: 440px;
  height: 378px;
}
.login-container {
  padding: 46px 60px;
  height: 100%;
  box-sizing: border-box;
  background: #fff;
  border-radius: 16px;
}

@media only screen and (max-width: 768px) {
  .login-container {
    padding: 22px 20px;
  }

  .login-popup {
    height: 328px;
  }
}

.title {
  text-align: center;
  color: #000000;
  font-size: 24px;
  font-style: normal;
  font-weight: 600;
  line-height: 32px;
  margin-bottom: 40px;
}

.form-item {
  position: relative;
  margin-bottom: 20px;

  .input-icon {
    position: absolute;
    left: 10px;
    top: 50%;
    transform: translateY(-50%);
    color: #999;
    font-size: 16px;
    z-index: 1;
  }

  .info-box {
    color: #8c8c8c;
    font-size: 12px;
    font-style: normal;
    font-weight: 400;
    line-height: 20px;

    .info-link {
      cursor: pointer;
      color: #2475fc;
    }
  }
}

.custom-input {
  width: 100%;
  padding: 8px 12px 8px 28px;
  align-items: center;
  height: 40px;
  border-radius: 6px;
  margin: 0 auto;
  transition: all 0.3s;
  border: 1px solid #D9D9D9;
  background: #FFF;
  outline:none;
}

.login-btn {
  width: 100%;
  height: 40px;
  border-radius: 6px;
  margin: 20px auto 0;
  background: #2475FC;
  border: none;
  color: #fff;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.3s;

  &:hover {
    opacity: 0.8;
  }
}

.login-btn:active {
  opacity: 0.9;
  transform: scale(0.98);
}

.van-popup {
  border-radius: 16px;
  overflow: hidden;
}

.van-popup__close-icon {
  color: #999;
  font-size: 18px;
  top: 12px;
  right: 12px;
}
</style>