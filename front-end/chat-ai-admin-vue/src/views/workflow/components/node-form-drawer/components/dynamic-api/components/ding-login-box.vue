<template>
  <div>
    <div v-if="!userInfo" class="not-login-box">
      <div class="tit">请先完成登录</div>
      <a-button type="primary" class="link" @click="showLogin">去登录</a-button>
    </div>
    <div v-else class="login-info-box">
      <img class="avatar" :src="userInfo?.avatarUrl"/>
      <span>{{ userInfo?.nick }}</span>
      <CheckCircleFilled class="icon"/>
    </div>

    <a-modal
      v-model:open="loginVisible"
      :footer="null"
      wrapClassName="qrcode-wrap"
      width="380px"
    >
      <div class="qrcode-box">
        <div class="title">扫码登录</div>
        <div class="content">
          <div class="tip">请使用钉钉扫描下方二维码授权</div>
          <div id="dingLoginCode" class="qrcode"></div>
<!--          <div class="desc">二维码有效期为5分钟，请尽快扫描</div>-->
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import {ref, onMounted, nextTick} from 'vue'
import {message} from 'ant-design-vue'
import {CheckCircleFilled} from '@ant-design/icons-vue'
import {runPlugin} from "@/api/plugins/index.js";
import {jsonDecode} from "@/utils/index.js";

const props = defineProps({
  formState: { type: Object, default: () => null },
})
const emit = defineEmits(['change'])

const loginVisible = ref(false)
const userInfo = ref(jsonDecode(localStorage.getItem('zm:workflow:ding:login:user')))
let scriptEl = null

onMounted(() => {
 setTimeout(() => {
   if (userInfo.value?.unionId) {
     update()
   }
 }, 100)
})

function showLogin() {
  if (!props.formState.dingtalk_app_key || !props.formState.dingtalk_app_secret) return message.error('缺少授权信息')
  loginVisible.value = true
  nextTick(() => {
    loadDingScript(renderLoginComp)
  })
}

function loadDingScript(callback= () => {}) {
  if (!window.DTFrameLogin) {
    scriptEl = document.createElement('script')
    scriptEl.src = 'https://g.alicdn.com/dingding/h5-dingtalk-login/0.21.0/ddlogin.js'
    scriptEl.async = true
    scriptEl.onload = () => callback()
    document.body.appendChild(scriptEl)
  } else {
    callback()
  }
}

function renderLoginComp() {
  window.DTFrameLogin(
    {
      id: 'dingLoginCode',
      width: 220,
      height: 240,
    },
    {
      redirect_uri: encodeURIComponent(window.location.origin),
      client_id: 'ding8xzffz9vztwiqdcu',
      scope: 'openid',
      response_type: 'code',
      prompt: 'consent',
    },
    (loginResult) => {
      console.log('loginResult', loginResult)
      const {redirectUrl, authCode} = loginResult;
      // 也可以在不跳转页面的情况下，使用code进行授权
      getUserInfo(authCode)
    },
    (errorMsg) => {
      // 这里一般需要展示登录失败的具体原因,可以使用toast等轻提示
      console.error(`errorMsg of errorCbk: ${errorMsg}`);
    },
  );
}

function getUserInfo(dingtalk_code) {
  const {dingtalk_app_key, dingtalk_app_secret} = props.formState
  runPlugin({
    name: "dingtalk_ai_table",
    action: "default/exec",
    params: JSON.stringify({
      business: 'GetLoginUserInfo',
      arguments: {
        dingtalk_app_key,
        dingtalk_app_secret,
        dingtalk_code,
      }
    })
  }).then(res => {
    if (res?.data?.unionId) {
      userInfo.value = res?.data
      localStorage.setItem('zm:workflow:ding:login:user', JSON.stringify(userInfo.value))
      update()
    }
  })
}

function update() {
  emit('change', 'operatorId', userInfo.value?.unionId)
}
</script>

<style scoped lang="less">
.not-login-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  border-radius: 6px;
  padding: 16px 16px 0;
  background: #F2F4F7;

  .link {
    width: 168px;
    margin-top: 16px;
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

.login-info-box {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #262626;
  font-size: 14px;

  .avatar {
    width: 16px;
    height: 16px;
    border-radius: 12px;
  }

  .icon {
    color: #21A665;
  }
}

.qrcode-box {
  .title {
    color: #000000;
    font-size: 24px;
    font-weight: 600;
    padding: 32px 24px;
    position: relative;
    &::before {
      content: "";
      position: absolute;
      width: 100%;
      height: 100%;
      left: 0;
      top: -18px;
      background: url("@/assets/img/login-pop-bg.png") no-repeat;
      background-size: 100% 100%;
    }
    &::after {
      content: "扫码登录";
      position: absolute;
      color: #000000;
      font-size: 24px;
      font-weight: 600;
      left: 24px;
      top: 24px;
    }
  }
  .content {
    display: flex;
    align-items: center;
    flex-direction: column;
    color: #262626;
    font-size: 14px;
    padding: 0 24px 24px;

    .qrcode {
      //width: 200px;
      //height: 200px;
    }

    .desc {
      color: #8c8c8c;
    }
  }
}
</style>
<style lang="less">
.qrcode-wrap.ant-modal-wrap {
  .ant-modal-content {
    padding: 0;
  }
}
</style>
