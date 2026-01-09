<style lang="less" scoped>
.client-download-page {
  position: relative;
  height: 100%;
  background-color: #fff;

  .page-title {
    line-height: 24px;
    padding: 24px 0 0 24px;
    font-size: 16px;
    font-weight: 600;
    color: #000000;
  }
  .page-body {
    display: flex;
    justify-content: center;
  }
  .app-info {
    margin-top: 49px;
    text-align: center;

    .app-name {
      line-height: 24px;
      margin-bottom: 16px;
      font-size: 16px;
      font-weight: 600;
      color: #262626;
    }

    .banner {
      position: relative;
      .download-btn {
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        width: auto;
        bottom: 24px;
        left: 50%;
        margin-left: -130px;
        line-height: 28px;
        padding: 12px 50px;
        border-radius: 6px;
        font-size: 20px;
        font-weight: 600;
        color: #ffffff;
        background: #2475fc;
        cursor: pointer;
        box-shadow:
          0 3px 14px 2px #0000000d,
          0 8px 10px 1px #0000000f,
          0 5px 5px -3px #0000001a;

        .windows-icon {
          margin-right: 10px;
        }
      }
      .download-btn.disabled {
        box-shadow: none;
        color: rgba(0, 0, 0, 0.25);
        border-color: #d9d9d9;
        background-color: rgba(0, 0, 0, 0.04);
      }
    }
  }
  .login-required-switch {
    position: absolute;
    top: 24px;
    right: 24px;
    display: flex;
    align-items: center;
    font-size: 14px;
    line-height: 22px;
    color: #262626;
    .need-login {
      margin-right: 4px;
    }
  }
}
</style>

<template>
  <div class="client-download-page">
    <div class="login-required-switch" v-if="role_type != '3'">
      <span class="need-login">{{ t('need_login') }}</span>
      <a-switch v-model:checked="loginRequired" @change="handleLoginRequiredSwitch" />
    </div>
    <h3 class="page-title">{{ t('client_download') }}</h3>
    <div class="page-body">
      <div class="app-info">
        <div class="app-name">{{ t('chat_wiki_pc_client') }}</div>
        <div class="banner">
          <img src="../../../assets/img/user/client-download/client-download-pic.png" alt="" />
          <div class="download-btn" @click="handleDownApp" :class="{ disabled: disabled }">
            <WindowsFilled class="windows-icon" /><span>{{ downloadBtnText }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { downloadFile } from '../../../utils/index.js'
import { useI18n } from '@/hooks/web/useI18n'
import { message } from 'ant-design-vue'
import { usePermissionStore } from '@/stores/modules/permission'
import { WindowsFilled } from '@ant-design/icons-vue'
import { getClientSideLoginSwitch, setClientSideLoginSwitch, clientSideDownload } from '@/api/user'

const permissionStore = usePermissionStore()

let { role_type } = permissionStore
const { t } = useI18n('views.user.client-download.index')
const loginRequired = ref(false)

const handleLoginRequiredSwitch = () => {
  setClientSideLoginSwitch({
    client_side_login_switch: loginRequired.value ? 1 : 0
  }).then(() => {
    message.success(t('common.saveSuccess'))
  })
}

const downloadLink = ref('')
const downloadBtnText = ref(t('windows_download'))
const disabled = ref(false)
let timer = null

const getDownloadLink = () => {
  return clientSideDownload({ domain: '' }).then((res) => {
    if (!res.data.file_url) {
      downloadBtnText.value = t('packaging')
      disabled.value = true
      // 3秒后重新获取下载链接
      timer = setTimeout(() => {
        getDownloadLink()
      }, 5000)
    } else {
      disabled.value = false
      downloadBtnText.value = t('windows_download')
      downloadLink.value = res.data.file_url
      clearTimeout(timer)
      timer = null
    }

    return res
  })
}

const handleDownApp = () => {
  if (disabled.value) {
    return
  }

  getDownloadLink().then((res) => {
    if (res.data.file_url) {
      downloadFile('', res.data.file_url)
    } else {
      message.error(t('download_not_ready'))
    }
  })
}

const init = () => {
  getClientSideLoginSwitch().then((res) => {
    loginRequired.value = res.data.client_side_login_switch == 1 ? true : false
  })
}

init()
</script>
