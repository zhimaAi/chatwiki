<template>
  <a-drawer
    :title="null"
    placement="bottom"
    :closable="false"
    :open="open"
    :width="176"
    rootClassName="cookie-auth-drawer"
    :maskClosable="false"
  >
    <div class="cookie-auth-box">
      <div class="cookie-content">
        <div class="cookie-title">{{ t('title_privacy_notice') }}</div>
        <div class="cookie-desc" v-html="cookieDescriptionHtml"></div>
        <div class="btn-block">
          <a-button @click="handleDecline">{{ t('btn_decline') }}</a-button>
          <a-button type="primary" @click="handleAccept">{{ t('btn_accept_all') }}</a-button>
        </div>
      </div>
    </div>
  </a-drawer>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { getCookieTip } from '@/api/user/index'
import { ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.login.cookie-auth')

const open = ref(false)

let LOCAL_COOKIE_KEY = 'cookie_auth_access_token'

const cookieDescriptionHtml = computed(() => {
  return t('msg_cookie_description', {
    cookiePolicy: `<a href="/#/privacy_policy">[${t('link_cookie_policy')}]</a>`,
    privacyPolicy: `<a href="/#/privacy_policy">[${t('link_privacy_policy')}]</a>`
  })
})

const handleAccept = () => {
  localStorage.setItem(LOCAL_COOKIE_KEY, 1)
  open.value = false
}

const handleDecline = () => {
  open.value = false
}

onMounted(() => {
  if (localStorage.getItem(LOCAL_COOKIE_KEY)) {
    return
  }
  getCookieTip({
    position: 'login'
  }).then((res) => {
    open.value = res.data.is_show_cookie_tip
  })
})
</script>

<style>
.cookie-auth-drawer .ant-drawer-content-wrapper {
  height: 180px !important;
}
</style>
<style lang="less" scoped>
.cookie-auth-box {
  width: 100%;
  .cookie-content {
    margin: 0 auto;
    width: 1200px;
  }
  .cookie-title {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }
  .cookie-desc {
    margin-top: 8px;
    color: #595959;
    font-size: 14px;
    line-height: 22px;
  }
  .btn-block {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 24px;
  }
}
</style>
