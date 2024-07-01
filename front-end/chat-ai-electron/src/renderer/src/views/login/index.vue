<style lang="less" scoped>
.login-page {
  position: relative;
  height: 100vh;
  width: 100vw;
  background: #f3f6fb;
  overflow: hidden;
  .page-header {
    position: fixed;
    left: 0;
    top: 0;
    right: 0;
    padding: 25px 0 0 16px;

    .logo-box {
      display: flex;
      align-items: center;
      .logo {
        width: 100px;
        height: 24px;
      }
      .logo-text {
        line-height: 22px;
        font-size: 14px;
        color: #595959;
      }
      .line {
        width: 1px;
        height: 12px;
        margin: 0 10px;
        flex-shrink: 0;
        border-radius: 1px;
        background: #d8dde5;
      }
    }
  }

  .page-body {
    position: relative;
    height: 100%;
    width: 100%;
    max-width: 950px;
    margin: 0 auto;
  }

  .login-bg {
    position: absolute;
    height: 334px;
    top: 50%;
    right: 0;
    margin-top: -167px;
  }

  .login-form-wrapper {
    position: absolute;
    width: 382px;
    height: 334px;
    top: 50%;
    right: 93px;
    margin-top: -167px;
    flex-shrink: 0;
    border-radius: 16px;
    background: #fff;
    box-shadow: 0 4px 32px 0 #1e44701f;
    .login-form-box {
      padding: 48px 44px 0 44px;
    }
    .login-form-title {
      line-height: 32px;
      margin-bottom: 28px;
      font-size: 24px;
      font-weight: 600;
      color: #000000;
      text-align: center;
    }
  }
}
.locale-dropdown-wrapper {
  position: fixed;
  bottom: 16px;
  left: 16px;
}
</style>

<template>
  <div class="login-page">
    <div class="page-header">
      <div class="logo-box" style="display: none">
        <img class="logo" src="../../assets//en_logo.svg" alt="" />
        <i class="line"></i>
        <span class="logo-text">企业AI知识库</span>
      </div>
    </div>
    <div class="page-body">
      <img class="login-bg" src="../../assets/img/login/page_bg.png" alt="" />
      <div class="login-form-wrapper">
        <div class="login-form-box">
          <div class="login-form-title">账号登录</div>
          <a-form
            :model="formState"
            name="basic"
            :label-col="{ span: 0 }"
            :wrapper-col="{ span: 24 }"
            layout="vertical"
            autocomplete="off"
            @finish="onFinish"
            @finish-failed="onFinishFailed"
          >
            <a-form-item
              label=""
              name="username"
              :rules="[{ required: true, message: '请输入账号' }]"
            >
              <a-input
                v-model:value="formState.username"
                style="width: 100%"
                placeholder="请输入账号"
              />
            </a-form-item>

            <a-form-item
              label=""
              name="password"
              :rules="[{ required: true, message: '请输入密码' }]"
            >
              <a-input-password v-model:value="formState.password" placeholder="请输入密码" />
            </a-form-item>

            <a-form-item>
              <a-button style="width: 100%" type="primary" html-type="submit">{{
                t('common.login')
              }}</a-button>
            </a-form-item>
          </a-form>
        </div>
      </div>
    </div>

    <div class="locale-dropdown-wrapper"><LocaleDropdown /></div>
  </div>
</template>

<script setup>
import { Modal, message } from 'ant-design-vue'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '../../hooks/web/useI18n'
import { useUserStore } from '@/stores/modules/user'
import LocaleDropdown from './components/locale-dropdown.vue'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()

const formState = reactive({
  username: '',
  password: '',
  remember: true
})

const submitLoading = ref(false)

const onFinish = () => {
  submitLoading.value = true

  userStore
    .login({
      toHome: true,
      username: formState.username,
      password: formState.password
    })
    .then(() => {
      submitLoading.value = false
      router.replace('/chat')
    })
    .catch((err) => {
      submitLoading.value = false

      if (err.data && err.data.type === 'client_side_cannot_login') {
        Modal.warning({
          title: '提示',
          content: err.message
        })
      } else {
        message.error(err.message)
      }
    })
}
const onFinishFailed = () => {
  // console.log('Failed:', errorInfo)
}
</script>
