<template>
  <span class="backup-model-label">
    <span>{{ label }}</span>
    <span v-if="showTip" class="backup-model-tip">
      <ExclamationCircleFilled />
      <span>{{ t('backup_model_disabled') }}</span>
      <template v-if="canManageModels">
        <a v-if="!targetDomain" :href="settingsUrl" target="_blank" rel="noopener noreferrer">{{ t('backup_model_settings') }}</a>
        <button v-else type="button" class="backup-model-settings" @click="openSettings">{{ t('backup_model_settings') }}</button>
      </template>
      <span v-else>{{ t('backup_model_contact_admin') }}</span>
    </span>
  </span>
</template>

<script setup>
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { ExclamationCircleFilled } from '@ant-design/icons-vue'
import { getBackupModelStatus } from '@/api/model/index'
import { useCompanyStore } from '@/stores/modules/company'
import { useUserStore } from '@/stores/modules/user'
import { usePermissionStore } from '@/stores/modules/permission'
import { buildSetTokenRedirectUrl, getCrossDomainTarget } from '@/utils/clawbot-domain'
import { useI18n } from '@/hooks/web/useI18n'

const props = defineProps({
  label: { type: String, required: true },
  active: { type: Boolean, default: true }
})
const { t } = useI18n('components.model-select.model-select')
const router = useRouter()
const companyStore = useCompanyStore()
const userStore = useUserStore()
const permissionStore = usePermissionStore()
const canManageModels = computed(() => permissionStore.role_permission.includes('ModelManage'))
const settingsPath = '/user/model?activeTab=3'
const targetDomain = computed(() => getCrossDomainTarget({
  toPath: '/user/model',
  fromPath: router.currentRoute.value.path,
  fromQuery: router.currentRoute.value.query,
  adminDomain: companyStore.admin_domain,
  agentDomain: companyStore.agent_domain,
  currentOrigin: window.location.origin
}))
const settingsUrl = computed(() => targetDomain.value
  ? `${targetDomain.value}/#${settingsPath}`
  : router.resolve(settingsPath).href)
const openSettings = () => {
  let url = settingsUrl.value
  if (targetDomain.value && userStore.getToken && userStore.userInfo?.user_id) {
    url = buildSetTokenRedirectUrl({
      domain: targetDomain.value,
      redirectUrl: settingsPath,
      token: userStore.getToken,
      exp: userStore.token?.exp,
      ttl: userStore.token?.ttl,
      userId: userStore.userInfo.user_id,
      userName: userStore.userInfo.user_name
    })
  }
  window.open(url, '_blank', 'noopener,noreferrer')
}
const showTip = ref(false)
let requestId = 0

const refresh = async () => {
  if (!props.active) return
  const id = ++requestId
  showTip.value = false
  try {
    const res = await getBackupModelStatus()
    if (id === requestId) showTip.value = res.data.enabled === false
  } catch {
    // 请求层负责错误提示，读取失败时不误报备用模型未开启。
    if (id === requestId) showTip.value = false
  }
}

watch(() => props.active, (active) => {
  showTip.value = false
  if (active) refresh()
  else requestId++
}, { immediate: true })

// 在新标签页完成设置后，返回当前窗口时更新提示。
onMounted(() => window.addEventListener('focus', refresh))
onBeforeUnmount(() => {
  requestId++
  window.removeEventListener('focus', refresh)
})
</script>

<style lang="less" scoped>
.backup-model-label {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 4px 8px;
}
.backup-model-tip {
  display: inline-flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  margin-left: auto;
  color: #8c8c8c;
  font-size: 12px;
  font-weight: 400;
  line-height: 22px;
}
.backup-model-settings {
  padding: 0;
  border: 0;
  background: none;
  color: #2475fc;
  font: inherit;
  cursor: pointer;
}
:global(.backup-model-label-item .ant-form-item-label > label) {
  width: 100%;
}
</style>
