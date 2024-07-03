<style lang="less" scoped>
.enterprise-set {
  width: 100%;
  height: 100%;
  padding: 24px;
  background-color: #fff;
  .page-title {
    line-height: 24px;
    font-size: 16px;
    font-weight: 600;
  }
  .enterprise-box {
    margin-top: 24px;
    .content-label {
      color: #262626;
      font-weight: 500;
    }
    .content-name {
      flex: 1;
      color: #333;
      font-weight: 600;
      .gray-text {
        color: #8c8c8c;
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
    <div class="page-title">{{ t('views.user.enterprise.enterpriseSettings') }}</div>
    <div class="enterprise-box">
      <a-flex align="center">
        <div class="content-label">{{ t('views.user.enterprise.systemName') }}：</div>
        <div class="content-name">
          <span v-if="name">{{ name }}</span>
          <span v-else class="gray-text">{{ t('views.user.enterprise.notSetTip') }}</span>
        </div>
        <a class="edit-btn" @click="openCompanyModal">{{ t('common.change') }}</a>
      </a-flex>
    </div>
    <a-divider></a-divider>
    <a-modal
      v-model:open="open"
      :title="t('views.user.enterprise.SetSystemName')"
      @ok="handleSetCompany"
    >
      <a-form class="form-box">
        <a-form-item :label="t('views.user.enterprise.systemName')">
          <a-input
            :maxlength="15"
            v-model:value="formState.name"
            :placeholder="t('views.user.enterprise.enterName')"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useCompanyStore } from '@/stores/modules/company'
import { saveCompany } from '@/api/user/index.js'
import { useI18n } from '@/hooks/web/useI18n'
import { message } from 'ant-design-vue'

const { t } = useI18n()

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
  formState.name = name.value;
  formState.id = id.value;
  open.value = true
}
const handleSetCompany = () => {
  saveCompany({
    ...formState
  }).then((res) => {
    message.success(t('common.saveSuccess'))
    handleGetCompany()
    let title = document.title.split('Chatwiki')
    document.title = title[0] + 'Chatwiki ' + formState.name
    open.value = false
  })
}
</script>
