<template>
  <div class="chatclaw-login-wrapper" :class="['lang-' + lang]">
    <div class="chatclaw-header">
      <div class="chatclaw-header-right">
        <LocaleDropdown />
      </div>
    </div>
    <div class="chatclaw-login-card" v-if="!loginSuccess">
      <div class="chatclaw-logo">
        <span class="chatclaw-logo-text">ChatClaw</span>
      </div>
      <p class="chatclaw-subtitle">{{ t('subtitle') }}</p>

      <a-form
        class="chatclaw-form"
        :model="formState"
        @finish="handleLogin"
      >
        <a-form-item
          name="username"
          :rules="[{ required: true, message: t('pleaseUsername') }]"
        >
          <a-input
            v-model:value="formState.username"
            size="large"
            :placeholder="t('pleaseUsername')"
            autocomplete="off"
          >
            <template #prefix><UserOutlined style="color: rgba(0,0,0,0.25)" /></template>
          </a-input>
        </a-form-item>

        <a-form-item
          name="password"
          :rules="[{ required: true, message: t('pleasePassword') }]"
        >
          <a-input-password
            v-model:value="formState.password"
            size="large"
            :placeholder="t('pleasePassword')"
            autocomplete="off"
          >
            <template #prefix><LockOutlined style="color: rgba(0,0,0,0.25)" /></template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-button
            class="chatclaw-login-btn"
            type="primary"
            html-type="submit"
            size="large"
            block
            :loading="loading"
          >{{ t('login') }}</a-button>
        </a-form-item>
      </a-form>
    </div>

    <div class="chatclaw-success-card" v-else>
      <a-result status="success" :title="t('loginSuccess')" :sub-title="t('openingApp')">
        <template #extra>
          <a-button type="primary" @click="openApp">{{ t('openAppManually') }}</a-button>
        </template>
      </a-result>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { chatClawLoginApi } from '@/api/chatclaw'
import { message } from 'ant-design-vue'
import { useLocaleStore } from '@/stores/modules/locale'
import { useI18n } from '@/hooks/web/useI18n'
import LocaleDropdown from '@/layouts/AdminLayout/compoents/locale-dropdown.vue'

const CHATCLAW_PROTOCOL = 'ChatClaw'

const localeStore = useLocaleStore()
const lang = computed(() => localeStore.getCurrentLocale.lang)
const { t } = useI18n('views.chatclaw.login')

const formState = reactive({
  username: '',
  password: '',
})

const loading = ref(false)
const loginSuccess = ref(false)
const tokenData = ref(null)

// 从 userAgent 中解析操作系统类型与版本，作为兜底
function detectOsFromUserAgent() {
  const ua = navigator.userAgent
  if (/Windows/.test(ua)) {
    const match = ua.match(/Windows NT ([\d.]+)/)
    const ntMap = { '10.0': '10/11', '6.3': '8.1', '6.2': '8', '6.1': '7' }
    const ver = match ? (ntMap[match[1]] || match[1]) : ''
    return { os_type: 'Windows', os_version: ver }
  }
  if (/Mac OS X/.test(ua)) {
    const match = ua.match(/Mac OS X ([\d_]+)/)
    return { os_type: 'macOS', os_version: match ? match[1].replace(/_/g, '.') : '' }
  }
  if (/Linux/.test(ua)) {
    return { os_type: 'Linux', os_version: '' }
  }
  return { os_type: '', os_version: '' }
}

// 优先读取 Electron 端注入的 URL query 参数（?os_type=Windows&os_version=11）
// 注意：项目使用 hash 路由，必须通过 useRoute() 读取 query，
// window.location.search 只能读取 # 之前的参数，在 hash 路由下永远为空
const route = useRoute()
const detected = detectOsFromUserAgent()
const clientInfo = {
  os_type:    String(route.query.os_type    || '') || detected.os_type,
  os_version: String(route.query.os_version || '') || detected.os_version,
}

function buildAppUrl(data) {
  const params = new URLSearchParams({
    token: data.token,
    ttl: String(data.ttl),
    exp: String(data.exp),
    user_id: String(data.user_id),
    user_name: String(data.user_name),
    server_url: window.location.origin,
  })
  return `${CHATCLAW_PROTOCOL}://auth/callback?${params.toString()}`
}

function openApp() {
  if (!tokenData.value) return
  const url = buildAppUrl(tokenData.value)
  window.location.href = url
}

async function handleLogin() {
  loading.value = true
  try {
    const res = await chatClawLoginApi(formState.username, formState.password, clientInfo)
    tokenData.value = res.data
    loginSuccess.value = true

    const url = buildAppUrl(res.data)
    window.location.href = url

    setTimeout(() => {
      message.info(t('manualOpenTip'))
    }, 2000)
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}
</script>

<style lang="less" scoped>
.chatclaw-login-wrapper {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #fff;
}

.chatclaw-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  height: 52px;
  padding: 0 24px;
  background: #fff;
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.06);
}

.chatclaw-header-right {
  display: flex;
  align-items: center;
}

.chatclaw-login-card {
  background: #fff;
  border-radius: 16px;
  padding: 48px 40px;
  width: 420px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
}

.chatclaw-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
}

.chatclaw-logo-text {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  letter-spacing: -0.5px;
}

.chatclaw-subtitle {
  text-align: center;
  color: #8c8c8c;
  font-size: 14px;
  margin-bottom: 32px;
}

.chatclaw-form {
  .chatclaw-login-btn {
    height: 44px;
    font-size: 16px;
    border-radius: 8px;
    border: none;

    &:hover {
      opacity: 0.9;
    }
  }
}

.chatclaw-success-card {
  background: #fff;
  border-radius: 16px;
  padding: 48px 40px;
  width: 480px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
}
</style>
