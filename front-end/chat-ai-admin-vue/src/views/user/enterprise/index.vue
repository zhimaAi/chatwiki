<style lang="less" scoped>
.enterprise-set {
  width: 100%;
  height: 100%;
  padding: 24px;
  overflow-y: auto;
  color: var(--app-text-primary);
  background-color: var(--app-bg-container);
  .page-title {
    color: var(--app-text-primary);
    line-height: 24px;
    font-size: 16px;
    font-weight: 600;
  }
  .system-mode-section {
    margin-top: 16px;
  }
  .section-title {
    margin: 0 0 12px;
    color: var(--app-text-primary);
    font-size: 14px;
    font-weight: 600;
  }
  .enterprise-box {
    margin-top: 24px;
    .content-label {
      color: var(--app-text-primary);
      font-weight: 500;
    }
    .content-name {
      flex: 1;
      color: var(--app-text-primary);
      font-weight: 600;
      .gray-text {
        color: var(--app-text-secondary);
        font-weight: 400;
      }
    }
    .edit-btn {
      margin-left: auto;
    }
  }
}
.form-box {
  margin-top: 38px;
  min-height: 60px;
}
</style>

<template>
  <div class="enterprise-set">
    <div class="page-title">{{ t('enterpriseSettings') }}</div>
    <section v-if="canManageCompany" class="system-mode-section">
      <h3 class="section-title">{{ modeT('title') }}</h3>
      <SystemModeSetting />
    </section>
    <div class="enterprise-box">
      <a-flex align="center">
        <div class="content-label">{{ t('systemName') }}：</div>
        <div class="content-name">
          <span v-if="name">{{ name }}</span>
          <span v-else class="gray-text">{{ t('notSetTip') }}</span>
        </div>
        <a class="edit-btn" @click="openCompanyModal">{{ t('change') }}</a>
      </a-flex>
    </div>
    <a-divider></a-divider>
    <div class="enterprise-box">
      <a-flex align="center">
        <div class="content-label">{{ t('labelCookieAuth') }}：</div>
        <div class="content-name">
          <span class="gray-text">{{ cookieText }}</span>
        </div>
        <a class="edit-btn" @click="handleOpenCookieModal">{{ t('change') }}</a>
      </a-flex>
    </div>
    <a-divider v-if="showTopNavigationSettings"></a-divider>
    <div v-if="showTopNavigationSettings" class="enterprise-box">
      <a-flex>
        <div class="content-label" style="line-height: 32px">
          {{ t('topNavigationSettings') }}：
        </div>
        <NavSetting @handleGetCompany="handleGetCompany" />
      </a-flex>
    </div>
    <a-modal v-model:open="open" :title="t('SetSystemName')" @ok="handleSetCompany">
      <a-form class="form-box">
        <a-form-item :label="t('systemName')">
          <a-input :maxlength="15" v-model:value="formState.name" :placeholder="t('enterName')" />
        </a-form-item>
      </a-form>
    </a-modal>
    <CookieSetting ref="cookieSettingRef" />
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useCompanyStore } from '@/stores/modules/company'
import { saveCompany } from '@/api/user/index.js'
import { useI18n } from '@/hooks/web/useI18n'
import { checkRole } from '@/utils/permission'
import { SYSTEM_MODES, useSystemModeStore } from '@/stores/modules/system-mode'
import { message } from 'ant-design-vue'
import NavSetting from './components/nav-setting.vue'
import CookieSetting from './components/cookie-setting.vue'
import SystemModeSetting from './components/system-mode-setting.vue'

const { t } = useI18n('views.user.enterprise')
const { t: modeT } = useI18n('views.system-mode')
const canManageCompany = checkRole(['CompanyManage'])
const modeStore = useSystemModeStore()
const showTopNavigationSettings = computed(
  () => modeStore.storageStatus === 'ready' && modeStore.currentMode === SYSTEM_MODES.ADVANCED
)

const companyStore = useCompanyStore()
const { companyInfo } = companyStore

const name = computed(() => {
  return companyStore.name
})
const id = computed(() => {
  return companyStore.id
})
const handleGetCompany = () => {
  companyStore.getCompanyInfo()
}
handleGetCompany()

const open = ref(false)
const formState = reactive({
  name: '',
  id: ''
})
const openCompanyModal = () => {
  formState.name = name.value
  formState.id = id.value
  open.value = true
}
const handleSetCompany = () => {
  saveCompany({
    ...formState
  }).then(() => {
    message.success(t('common.saveSuccess'))
    handleGetCompany()
    let title = document.title.split('Chatwiki')
    document.title = title[0] + 'Chatwiki ' + formState.name
    open.value = false
  })
}

const cookieSettingRef = ref(null)
const handleOpenCookieModal = () => {
  cookieSettingRef.value.show()
}

const cookieText = computed(() => {
  let cookie_tip_positions = companyInfo.cookie_tip_positions
    ? companyInfo.cookie_tip_positions.split(',')
    : []

  let cookie_tip_ip_locations = companyInfo.cookie_tip_ip_locations
    ? companyInfo.cookie_tip_ip_locations.split(',')
    : []

  if (cookie_tip_positions.length == 0 && cookie_tip_ip_locations.length == 0) {
    return t('msgCookieNotSet')
  }

  let cookie_tip_ip_locations_map = {
    domestic: t('labelDomesticIp'),
    overseas: t('labelOverseasIp')
  }

  let cookie_tip_positions_map = {
    home: t('labelHome'),
    login: t('labelLoginPage'),
    webapp: t('labelWebapp')
  }
  let result = []
  cookie_tip_ip_locations.forEach((item) => {
    result.push(cookie_tip_ip_locations_map[item])
  })
  cookie_tip_positions.forEach((item) => {
    result.push(cookie_tip_positions_map[item])
  })
  result = result.filter(Boolean)
  return `${t('msgCookieEnabled')}（${result.join('、')}）`
})
</script>
