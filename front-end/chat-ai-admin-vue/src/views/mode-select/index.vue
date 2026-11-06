<template>
  <main class="mode-select-page">
    <img class="welcome-logo" :src="companyStore.companyInfo?.avatar || logo" alt="ChatWiki" />
    <div v-if="loading" class="page-feedback"><a-spin /><span>{{ t('loading') }}</span></div>
    <div v-else-if="error" class="page-feedback"><a-alert type="error" show-icon :message="t('error_scope')" /><a-button @click="initialize">{{ t('retry') }}</a-button></div>
    <section v-else class="welcome-content">
      <h1>{{ t('welcome_title') }}</h1>
      <SystemModeSetting onboarding />
    </section>
  </main>
</template>

<script setup>
import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCompanyStore } from '@/stores/modules/company'
import { useUserStore } from '@/stores/modules/user'
import { checkRole } from '@/utils/permission'
import { useI18n } from '@/hooks/web/useI18n'
import SystemModeSetting from '@/views/user/enterprise/components/system-mode-setting.vue'
import logo from '@/assets/default_logo.svg'

const { t } = useI18n('views.system-mode')
const router = useRouter()
const companyStore = useCompanyStore()
const userStore = useUserStore()
const loading = ref(true)
const error = ref(false)
let alive = true
async function initialize() {
  loading.value = true
  error.value = false
  try {
    await companyStore.getCompanyInfo()
    if (!alive) return
    if (!userStore.user_is_admin && !checkRole(['CompanyManage'])) { await router.replace('/'); return }
  } catch { if (alive) error.value = true }
  finally { if (alive) loading.value = false }
}
void initialize()
onBeforeUnmount(() => { alive = false })
</script>

<style lang="less" scoped>
.mode-select-page { min-height: 100vh; padding: 40px 32px; background: radial-gradient(ellipse at 50% 45%, var(--app-selected) 0, var(--app-bg-layout) 100%); display: flex; align-items: center; justify-content: center; }
.welcome-logo { position: absolute; top: 32px; left: 40px; width: 32px; height: 32px; object-fit: contain; }
.welcome-content { width: 100%; max-width: 1100px; padding: 60px 0; }
h1 { margin: 0 0 32px; text-align: center; font-size: 26px; line-height: 38px; font-weight: 600; color: var(--app-text-primary); }
.page-feedback { display: flex; align-items: center; gap: 16px; }
@media (max-width: 800px) { .mode-select-page { padding: 80px 20px 24px; } .welcome-content { padding: 0; } h1 { font-size: 22px; } .welcome-logo { top: 24px; left: 24px; } }
</style>
